package protoconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
)

// FuzzLoad feeds arbitrary YAML + env to Load and asserts it never panics.
// Errors are fine; panics are not.
func FuzzLoad(f *testing.F) {
	// Seeds: realistic-ish inputs.
	f.Add([]byte("string_field: hi\n"), "APP__INT32_FIELD", "5")
	f.Add([]byte("list_of_ints: [1, 2, 3]\n"), "APP__STRING_FIELD", `"x"`)
	f.Add([]byte("nested_message_field:\n  string_field: deep\n"), "APP__BOOL_FIELD", "true")
	f.Add([]byte("primitive_map: {a: 1}\n"), "APP__PRIMITIVE_MAP__b", `"2"`)
	f.Add([]byte(""), "", "")
	f.Add([]byte("garbage!@#$%"), "APP", "APP")

	tmpDir := f.TempDir()
	f.Fuzz(func(t *testing.T, yaml []byte, envKey, envVal string) {
		// Filter obviously destructive inputs (NUL in path, etc.).
		if strings.ContainsAny(envKey, "\x00") {
			return
		}
		path := filepath.Join(tmpDir, "c.yaml")
		if err := os.WriteFile(path, yaml, 0600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		environ := func() []string {
			if envKey == "" {
				return nil
			}
			return []string{envKey + "=" + envVal}
		}
		// Run Load. Don't assert on the result, just ensure no panic.
		_, _ = Load(&protoconfigv1.Test{},
			FromYAMLFile(path, Optional()),
			FromEnv("APP", EnvironFunc(environ)),
		)
	})
}

// FuzzLoadStrict runs in strict mode.
func FuzzLoadStrict(f *testing.F) {
	f.Add([]byte("string_field: hi\n"), "APP__INT32_FIELD", "5")
	f.Add([]byte("list_of_ints: [1, 2, 3]\n"), "APP__STRING_FIELD", `"x"`)

	tmpDir := f.TempDir()
	f.Fuzz(func(t *testing.T, yaml []byte, envKey, envVal string) {
		if strings.ContainsAny(envKey, "\x00") {
			return
		}
		path := filepath.Join(tmpDir, "c.yaml")
		if err := os.WriteFile(path, yaml, 0600); err != nil {
			t.Fatalf("WriteFile: %v", err)
		}
		environ := func() []string {
			if envKey == "" {
				return nil
			}
			return []string{envKey + "=" + envVal}
		}
		_, _ = Load(&protoconfigv1.Test{},
			FromYAMLFile(path, Optional()),
			FromEnv("APP", EnvironFunc(environ)),
			Strict(true),
		)
	})
}

// FuzzEnvEntriesApply feeds random env entries directly to applyEnv, bypassing
// YAML and prefix filtering. Catches panics in the tree/apply logic.
func FuzzEnvEntriesApply(f *testing.F) {
	f.Add("STRING_FIELD", `"hi"`)
	f.Add("LIST_OF_INTS__0", "1")
	f.Add("REPEATED_NESTED_MESSAGE__0__STRING_FIELD", `"x"`)
	f.Add("STRING_TO_MAP__key__STRING_FIELD", `"v"`)
	f.Add("NESTED_WITH_NESTED__NESTED_NESTED__DEEPLY_NESTED_STRING", `"y"`)
	f.Add("TIMESTAMP", `"2024-01-01T00:00:00Z"`)

	f.Fuzz(func(t *testing.T, key, val string) {
		if strings.ContainsAny(key, "\x00") {
			return
		}
		segs := strings.Split(key, "__")
		cfg := &protoconfigv1.Test{}
		_ = applyEnv(cfg.ProtoReflect(), []envEntry{{segments: segs, value: val, rawKey: key}}, false)
	})
}
