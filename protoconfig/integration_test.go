package protoconfig

import (
	"strings"
	"testing"
	"time"

	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/utils/ptr"
)

// --- the big end-to-end everything test ---

func TestIntegration_AllScalarKinds_ViaEnv(t *testing.T) {
	defaults := &protoconfigv1.Test{}
	cfg, err := Load(defaults,
		FromEnv("APP", EnvironFunc(func() []string {
			return envList(
				"APP__BOOL_FIELD", "true",
				"APP__ENUM_FIELD", "EXAMPLE_ENUM_EXAMPLE_VAL",
				"APP__INT32_FIELD", "-100",
				"APP__SINT32_FIELD", "-200",
				"APP__UINT32_FIELD", "300",
				"APP__INT64_FIELD", `"-9000000000"`,
				"APP__SINT64_FIELD", `"-8000000000"`,
				"APP__UINT64_FIELD", `"9000000000"`,
				"APP__SFIXED32_FIELD", "-400",
				"APP__FIXED32_FIELD", "500",
				"APP__FLOAT_FIELD", "1.5",
				"APP__SFIXED64_FIELD", `"-7000000000"`,
				"APP__FIXED64_FIELD", `"6000000000"`,
				"APP__DOUBLE_FIELD", "2.25",
				"APP__STRING_FIELD", `"hello"`,
				"APP__BYTES_FIELD", `"c29tZS1ieXRlcw=="`,
				"APP__OPTIONAL_INT32_FIELD", "42",
			)
		})),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	want := &protoconfigv1.Test{
		BoolField:          true,
		EnumField:          protoconfigv1.Test_EXAMPLE_ENUM_EXAMPLE_VAL,
		Int32Field:         -100,
		Sint32Field:        -200,
		Uint32Field:        300,
		Int64Field:         -9000000000,
		Sint64Field:        -8000000000,
		Uint64Field:        9000000000,
		Sfixed32Field:      -400,
		Fixed32Field:       500,
		FloatField:         1.5,
		Sfixed64Field:      -7000000000,
		Fixed64Field:       6000000000,
		DoubleField:        2.25,
		StringField:        "hello",
		BytesField:         []byte("some-bytes"),
		OptionalInt32Field: ptr.To(int32(42)),
	}
	if diff := cmp.Diff(want, cfg, protocmp.Transform()); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

// --- list of enums ---

func TestIntegration_ListOfEnums(t *testing.T) {
	cfg, err := Load(&protoconfigv1.Test{},
		FromEnv("APP", EnvironFunc(func() []string {
			return envList(
				"APP__LIST_OF_ENUMS", `["EXAMPLE_ENUM_EXAMPLE_VAL","EXAMPLE_ENUM_UNSPECIFIED"]`,
			)
		})),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := []protoconfigv1.Test_ExampleEnum{
		protoconfigv1.Test_EXAMPLE_ENUM_EXAMPLE_VAL,
		protoconfigv1.Test_EXAMPLE_ENUM_UNSPECIFIED,
	}
	if diff := cmp.Diff(want, cfg.ListOfEnums); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

// --- WKTs ---

func TestIntegration_WKTsViaEnv(t *testing.T) {
	cfg, err := Load(&protoconfigv1.Test{},
		FromEnv("APP", EnvironFunc(func() []string {
			return envList(
				"APP__TIMESTAMP", `"2024-02-03T17:00:00Z"`,
				"APP__DURATION", `"1.5s"`,
				"APP__STR_WRAPPER", `"wrapped"`,
				"APP__INT32_WRAPPER", `42`,
				"APP__BOOL_WRAPPER", `true`,
				"APP__FIELD_MASK", `"a.b,c.d"`,
				"APP__STRUCT_FIELD", `{"nested": {"x": 1}, "arr": [1, 2]}`,
			)
		})),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	wantTS := time.Date(2024, 2, 3, 17, 0, 0, 0, time.UTC)
	if !cfg.Timestamp.AsTime().Equal(wantTS) {
		t.Errorf("timestamp: got %v want %v", cfg.Timestamp.AsTime(), wantTS)
	}

	wantDur := durationpb.New(1500 * time.Millisecond)
	if cfg.Duration.AsDuration() != wantDur.AsDuration() {
		t.Errorf("duration: got %v want %v", cfg.Duration.AsDuration(), wantDur.AsDuration())
	}

	if cfg.StrWrapper.GetValue() != "wrapped" {
		t.Errorf("str_wrapper: got %q", cfg.StrWrapper.GetValue())
	}
	if cfg.Int32Wrapper.GetValue() != 42 {
		t.Errorf("int32_wrapper: got %v", cfg.Int32Wrapper.GetValue())
	}
	if cfg.BoolWrapper.GetValue() != true {
		t.Errorf("bool_wrapper: got %v", cfg.BoolWrapper.GetValue())
	}

	wantMask := &fieldmaskpb.FieldMask{Paths: []string{"a.b", "c.d"}}
	if diff := cmp.Diff(wantMask, cfg.FieldMask, protocmp.Transform()); diff != "" {
		t.Errorf("field_mask: -want +got:\n%s", diff)
	}

	if cfg.StructField == nil {
		t.Fatal("struct_field nil")
	}
	nested, ok := cfg.StructField.Fields["nested"].Kind.(*structpb.Value_StructValue)
	if !ok {
		t.Fatalf("struct_field.nested: got %T", cfg.StructField.Fields["nested"].Kind)
	}
	if nested.StructValue.Fields["x"].GetNumberValue() != 1 {
		t.Errorf("struct_field.nested.x: got %v", nested.StructValue.Fields["x"].GetNumberValue())
	}
}

func TestIntegration_WrapperInitiallyNilViaEnv(t *testing.T) {
	// str_wrapper is nil by default. Env sets value. Outcome: wrapper is
	// instantiated and holds the value.
	cfg, err := Load(&protoconfigv1.Test{},
		FromEnv("APP", EnvironFunc(func() []string {
			return envList("APP__STR_WRAPPER", `"ok"`)
		})),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.StrWrapper.GetValue() != "ok" {
		t.Errorf("got %q", cfg.StrWrapper.GetValue())
	}
}

// --- wrapper preserves existing value when yaml sets a different field ---

func TestIntegration_WrapperPreservedAcrossYAML(t *testing.T) {
	tf := writeFile(t, "string_field: from-yaml\n")
	defaults := &protoconfigv1.Test{
		StrWrapper: wrapperspb.String("default"),
	}
	cfg, err := Load(defaults, FromYAMLFile(tf))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.StrWrapper.GetValue() != "default" {
		t.Errorf("got %q want default", cfg.StrWrapper.GetValue())
	}
}

// --- non-string-keyed maps ---

func TestIntegration_IntKeyedMapViaEnv(t *testing.T) {
	cfg, err := Load(&protoconfigv1.Test{},
		FromEnv("APP", EnvironFunc(func() []string {
			return envList(
				"APP__INT_KEYED_MAP__1", `"one"`,
				"APP__INT_KEYED_MAP__42", `"forty-two"`,
			)
		})),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := map[int32]string{1: "one", 42: "forty-two"}
	if diff := cmp.Diff(want, cfg.IntKeyedMap); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

func TestIntegration_BoolKeyedMapViaEnv(t *testing.T) {
	cfg, err := Load(&protoconfigv1.Test{},
		FromEnv("APP", EnvironFunc(func() []string {
			return envList(
				"APP__BOOL_KEYED_MAP__true", `"on"`,
				"APP__BOOL_KEYED_MAP__false", `"off"`,
			)
		})),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	want := map[bool]string{true: "on", false: "off"}
	if diff := cmp.Diff(want, cfg.BoolKeyedMap); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

// --- deeply nested ---

func TestIntegration_DeeplyNested(t *testing.T) {
	cfg, err := Load(&protoconfigv1.Test{},
		FromEnv("APP", EnvironFunc(func() []string {
			return envList(
				"APP__NESTED_WITH_NESTED__NESTED_NESTED__DEEPLY_NESTED_STRING", `"deep"`,
			)
		})),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	got := cfg.NestedWithNested.GetNestedNested().GetDeeplyNestedString()
	if got != "deep" {
		t.Errorf("got %q", got)
	}
}

// --- oneof ---

func TestIntegration_OneofViaEnv(t *testing.T) {
	cfg, err := Load(&protoconfigv1.Test{},
		FromEnv("APP", EnvironFunc(func() []string {
			return envList("APP__CHOICE_STRING", `"picked"`)
		})),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	gotStr, ok := cfg.Choice.(*protoconfigv1.Test_ChoiceString)
	if !ok {
		t.Fatalf("got %T want *Test_ChoiceString", cfg.Choice)
	}
	if gotStr.ChoiceString != "picked" {
		t.Errorf("got %q", gotStr.ChoiceString)
	}
}

func TestIntegration_OneofMultipleArmsError(t *testing.T) {
	_, err := Load(&protoconfigv1.Test{},
		FromEnv("APP", EnvironFunc(func() []string {
			return envList(
				"APP__CHOICE_STRING", `"a"`,
				"APP__CHOICE_INT", `1`,
			)
		})),
		Strict(true),
	)
	if err == nil {
		t.Fatal("want error for multi-arm oneof")
	}
	if !strings.Contains(err.Error(), "oneof") {
		t.Errorf("got %v; want mention of oneof", err)
	}
}

func TestIntegration_OneofEnvOverridesYAMLArm(t *testing.T) {
	tf := writeFile(t, "choice_string: from-yaml\n")
	cfg, err := Load(&protoconfigv1.Test{},
		FromYAMLFile(tf),
		FromEnv("APP", EnvironFunc(func() []string {
			return envList("APP__CHOICE_INT", `99`)
		})),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	gotInt, ok := cfg.Choice.(*protoconfigv1.Test_ChoiceInt)
	if !ok {
		t.Fatalf("got %T want *Test_ChoiceInt", cfg.Choice)
	}
	if gotInt.ChoiceInt != 99 {
		t.Errorf("got %d", gotInt.ChoiceInt)
	}
}

func TestIntegration_OneofMessageArm(t *testing.T) {
	cfg, err := Load(&protoconfigv1.Test{},
		FromEnv("APP", EnvironFunc(func() []string {
			return envList(
				"APP__CHOICE_MSG__STRING_FIELD", `"inner"`,
			)
		})),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	msg, ok := cfg.Choice.(*protoconfigv1.Test_ChoiceMsg)
	if !ok {
		t.Fatalf("got %T", cfg.Choice)
	}
	if msg.ChoiceMsg.StringField != "inner" {
		t.Errorf("got %q", msg.ChoiceMsg.StringField)
	}
}

// --- self-referencing message ---

func TestIntegration_SelfRefNoEnvStaysNil(t *testing.T) {
	cfg, err := Load(&protoconfigv1.Test{StringField: "default"})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.MessageField != nil {
		t.Errorf("MessageField should be nil, got %+v", cfg.MessageField)
	}
}

// --- YAML + env compose ---

func TestIntegration_YAMLPlusEnvCompose(t *testing.T) {
	tf := writeFile(t, `
string_field: from-yaml
int32_field: 100
repeated_nested_message:
  - string_field: "y0"
    not_updated_via_env: "y0-keep"
  - string_field: "y1"
    not_updated_via_env: "y1-keep"
`)

	cfg, err := Load(&protoconfigv1.Test{StringField: "default", Int64Field: 7},
		FromYAMLFile(tf),
		FromEnv("APP", EnvironFunc(func() []string {
			return envList(
				"APP__STRING_FIELD", `"from-env"`,
				"APP__REPEATED_NESTED_MESSAGE__1__STRING_FIELD", `"env-overlay"`,
			)
		})),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	// Defaults: Int64Field=7 (kept).
	// YAML: StringField=from-yaml, Int32Field=100, repeated_nested_message=[y0, y1].
	// Env:  StringField=from-env (wins), index 1 overlay.
	want := &protoconfigv1.Test{
		StringField: "from-env",
		Int32Field:  100,
		Int64Field:  7,
		RepeatedNestedMessage: []*protoconfigv1.Nested2{
			{StringField: "y0", NotUpdatedViaEnv: "y0-keep"},
			{StringField: "env-overlay", NotUpdatedViaEnv: "y1-keep"},
		},
	}
	if diff := cmp.Diff(want, cfg, protocmp.Transform()); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

// --- parsing errors produce typed error ---

func TestIntegration_ParseErrorWrapped(t *testing.T) {
	_, err := Load(&protoconfigv1.Test{},
		FromEnv("APP", EnvironFunc(func() []string {
			return envList("APP__ENUM_FIELD", "NO_SUCH_ENUM")
		})),
	)
	if err == nil {
		t.Fatal("want error")
	}
	var le *Error
	if !asLoadError(err, &le) {
		t.Fatalf("not *Error: %v", err)
	}
	if !strings.HasPrefix(le.Layer, "env:") {
		t.Errorf("Layer=%s", le.Layer)
	}
}

func asLoadError(err error, target **Error) bool {
	for err != nil {
		if le, ok := err.(*Error); ok {
			*target = le
			return true
		}
		type unwrapper interface{ Unwrap() error }
		u, ok := err.(unwrapper)
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}

// Prevent unused-import warnings — _ = timestamppb if the test gets trimmed.
var _ = timestamppb.Now
