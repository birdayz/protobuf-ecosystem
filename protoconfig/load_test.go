package protoconfig

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func writeFile(t *testing.T, body string) string {
	t.Helper()
	tf := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(tf, []byte(body), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return tf
}

func envList(pairs ...string) []string {
	out := make([]string, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, pairs[i]+"="+pairs[i+1])
	}
	return out
}

// --- basic layering ---

func TestLoad_DefaultsOnly(t *testing.T) {
	defaults := &protoconfigv1.Test{StringField: "default"}
	cfg, err := Load(defaults)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.StringField != "default" {
		t.Errorf("got %q", cfg.StringField)
	}
}

func TestLoad_YAMLOverDefaults(t *testing.T) {
	tf := writeFile(t, "string_field: from-yaml\n")
	defaults := &protoconfigv1.Test{StringField: "default", Int32Field: 7}

	cfg, err := Load(defaults, FromYAMLFile(tf))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := &protoconfigv1.Test{StringField: "from-yaml", Int32Field: 7}
	if diff := cmp.Diff(want, cfg, protocmp.Transform()); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

func TestLoad_EnvOverDefaults(t *testing.T) {
	defaults := &protoconfigv1.Test{StringField: "default"}
	cfg, err := Load(defaults,
		FromEnv("APP", EnvironFunc(func() []string {
			return envList("APP__STRING_FIELD", `"from-env"`)
		})),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.StringField != "from-env" {
		t.Errorf("got %q", cfg.StringField)
	}
}

func TestLoad_Precedence_EnvBeatsYAML(t *testing.T) {
	tf := writeFile(t, "string_field: from-yaml\nint32_field: 1\n")
	defaults := &protoconfigv1.Test{StringField: "default"}

	cfg, err := Load(defaults,
		FromYAMLFile(tf),
		FromEnv("APP", EnvironFunc(func() []string {
			return envList("APP__STRING_FIELD", `"from-env"`)
		})),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	// YAML set Int32Field=1; env overrode StringField only.
	want := &protoconfigv1.Test{StringField: "from-env", Int32Field: 1}
	if diff := cmp.Diff(want, cfg, protocmp.Transform()); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

// --- defaults immutability ---

func TestLoad_DefaultsNotMutated(t *testing.T) {
	tf := writeFile(t, "string_field: from-yaml\n")
	defaults := &protoconfigv1.Test{
		StringField: "default",
		NestedMessageField: &protoconfigv1.Nested{
			StringField: "nested-default",
		},
		ListOfInts: []int32{1, 2, 3},
	}
	snap := proto.Clone(defaults)

	_, err := Load(defaults,
		FromYAMLFile(tf),
		FromEnv("APP", EnvironFunc(func() []string {
			return envList(
				"APP__LIST_OF_INTS", "[9,8,7]",
				"APP__NESTED_MESSAGE_FIELD__STRING_FIELD", `"mutated"`,
			)
		})),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if diff := cmp.Diff(snap, defaults, protocmp.Transform()); diff != "" {
		t.Errorf("defaults mutated! -before +after:\n%s", diff)
	}
}

// --- strict flag ---

func TestLoad_StrictYAMLUnknownField(t *testing.T) {
	tf := writeFile(t, "totally_unknown_field: 1\n")
	defaults := &protoconfigv1.Test{}

	// Strict errors.
	_, err := Load(defaults, FromYAMLFile(tf), Strict(true))
	if err == nil {
		t.Errorf("want error in strict mode")
	}

	// Lenient ignores.
	_, err = Load(defaults, FromYAMLFile(tf))
	if err != nil {
		t.Errorf("lenient should not error, got %v", err)
	}
}

func TestLoad_StrictEnvUnknownKey(t *testing.T) {
	defaults := &protoconfigv1.Test{}

	environ := func() []string {
		return envList("APP__NO_SUCH_FIELD", "x")
	}
	_, err := Load(defaults,
		FromEnv("APP", EnvironFunc(environ)),
		Strict(true),
	)
	if err == nil {
		t.Errorf("want strict error for unknown env")
	}

	_, err = Load(defaults,
		FromEnv("APP", EnvironFunc(environ)),
	)
	if err != nil {
		t.Errorf("lenient should not error, got %v", err)
	}
}

// --- error wrapping ---

func TestLoad_ErrorHasLayer(t *testing.T) {
	defaults := &protoconfigv1.Test{}
	_, err := Load(defaults,
		FromEnv("APP", EnvironFunc(func() []string {
			return envList("APP__INT32_FIELD", "not-an-int")
		})),
	)
	if err == nil {
		t.Fatal("want error")
	}
	var le *Error
	if !errors.As(err, &le) {
		t.Fatalf("want *Error, got %T", err)
	}
	if !strings.HasPrefix(le.Layer, "env:") {
		t.Errorf("Layer=%q want env:...", le.Layer)
	}
}

// --- optional file ---

func TestLoad_OptionalMissingFile(t *testing.T) {
	defaults := &protoconfigv1.Test{StringField: "default"}
	cfg, err := Load(defaults,
		FromYAMLFile("/nonexistent/path/to/config.yaml", Optional()),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.StringField != "default" {
		t.Errorf("got %q", cfg.StringField)
	}
}

func TestLoad_RequiredMissingFile(t *testing.T) {
	defaults := &protoconfigv1.Test{}
	_, err := Load(defaults, FromYAMLFile("/nonexistent/path/config.yaml"))
	if err == nil {
		t.Errorf("want error for missing file")
	}
}

// --- YAML replace semantics (D5) ---

func TestLoad_YAMLListReplaces(t *testing.T) {
	tf := writeFile(t, "list_of_ints: [4, 5]\n")
	defaults := &protoconfigv1.Test{ListOfInts: []int32{1, 2, 3}}

	cfg, err := Load(defaults, FromYAMLFile(tf))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := []int32{4, 5}
	if diff := cmp.Diff(want, cfg.ListOfInts); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

func TestLoad_YAMLMapMerges(t *testing.T) {
	tf := writeFile(t, `
primitive_map:
  b: overridden
  c: new-key
`)
	defaults := &protoconfigv1.Test{PrimitiveMap: map[string]string{"a": "1", "b": "2"}}

	cfg, err := Load(defaults, FromYAMLFile(tf))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := map[string]string{"a": "1", "b": "overridden", "c": "new-key"}
	if diff := cmp.Diff(want, cfg.PrimitiveMap); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

// --- WKT integration ---

func TestLoad_TimestampViaEnv(t *testing.T) {
	defaults := &protoconfigv1.Test{
		Timestamp: timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
	}
	cfg, err := Load(defaults,
		FromEnv("APP", EnvironFunc(func() []string {
			return envList("APP__TIMESTAMP", `"2024-06-01T12:00:00Z"`)
		})),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	if !cfg.Timestamp.AsTime().Equal(want) {
		t.Errorf("got %v want %v", cfg.Timestamp.AsTime(), want)
	}
}

// --- transform func ---

func TestLoad_EnvTransformFunc(t *testing.T) {
	defaults := &protoconfigv1.Test{}
	cfg, err := Load(defaults,
		FromEnv("APP",
			EnvironFunc(func() []string {
				return envList(
					"APP__STRING_FIELD", `"keep"`,
					"APP__DROP_ME", "x", // dropped by transform
				)
			}),
			EnvTransformFunc(func(k, v string) (string, string, bool) {
				if strings.Contains(k, "DROP_ME") {
					return "", "", false
				}
				return k, v, true
			}),
		),
		Strict(true), // would otherwise error on DROP_ME
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.StringField != "keep" {
		t.Errorf("got %q", cfg.StringField)
	}
}
