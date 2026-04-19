package protoconfig

import (
	"bytes"
	"math"
	"testing"
	"time"

	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// --- scalars ---

func TestParseScalar_Bool(t *testing.T) {
	// protoyaml is strict: only JSON-shape booleans.
	cases := []struct {
		in   string
		want bool
	}{
		{"true", true},
		{"false", false},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			v, err := parseScalar(protoreflect.BoolKind, nil, c.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v.Bool() != c.want {
				t.Errorf("got %v want %v", v.Bool(), c.want)
			}
		})
	}
}

func TestParseScalar_BoolInvalid(t *testing.T) {
	// Shell-style variants not accepted. Users must write true/false.
	for _, bad := range []string{"True", "1", "0", "yes", "maybe", "2", "yesss"} {
		if _, err := parseScalar(protoreflect.BoolKind, nil, bad); err == nil {
			t.Errorf("want error for %q, got nil", bad)
		}
	}
}

func TestParseScalar_Int32(t *testing.T) {
	cases := []struct {
		in   string
		want int32
	}{
		{"0", 0},
		{"1", 1},
		{"-1", -1},
		{"2147483647", math.MaxInt32},
		{"-2147483648", math.MinInt32},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			v, err := parseScalar(protoreflect.Int32Kind, nil, c.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if int32(v.Int()) != c.want {
				t.Errorf("got %v want %v", v.Int(), c.want)
			}
		})
	}
}

func TestParseScalar_Int32Overflow(t *testing.T) {
	for _, bad := range []string{"2147483648", "-2147483649", "99999999999"} {
		if _, err := parseScalar(protoreflect.Int32Kind, nil, bad); err == nil {
			t.Errorf("want overflow error for %q, got nil", bad)
		}
	}
}

func TestParseScalar_Uint32NegativeError(t *testing.T) {
	if _, err := parseScalar(protoreflect.Uint32Kind, nil, "-1"); err == nil {
		t.Errorf("want error for negative uint, got nil")
	}
}

func TestParseScalar_Int64Large(t *testing.T) {
	// int64 precision must be preserved. Quote the value to bypass JSON
	// precision loss.
	v, err := parseScalar(protoreflect.Int64Kind, nil, `"9223372036854775807"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Int() != math.MaxInt64 {
		t.Errorf("got %v want %v", v.Int(), int64(math.MaxInt64))
	}
}

func TestParseScalar_Float(t *testing.T) {
	v, err := parseScalar(protoreflect.FloatKind, nil, "3.25")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Float() != 3.25 {
		t.Errorf("got %v want 3.25", v.Float())
	}
}

func TestParseScalar_FloatSpecials(t *testing.T) {
	cases := map[string]float64{
		`"NaN"`:       math.NaN(),
		`"Infinity"`:  math.Inf(1),
		`"-Infinity"`: math.Inf(-1),
	}
	for in, want := range cases {
		t.Run(in, func(t *testing.T) {
			v, err := parseScalar(protoreflect.DoubleKind, nil, in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := v.Float()
			if math.IsNaN(want) {
				if !math.IsNaN(got) {
					t.Errorf("got %v want NaN", got)
				}
				return
			}
			if got != want {
				t.Errorf("got %v want %v", got, want)
			}
		})
	}
}

func TestParseScalar_String(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{`"hello"`, "hello"},
		{"hello", "hello"}, // bare YAML string
		{`"line1\nline2"`, "line1\nline2"},
		{`""`, ""},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			v, err := parseScalar(protoreflect.StringKind, nil, c.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if v.String() != c.want {
				t.Errorf("got %q want %q", v.String(), c.want)
			}
		})
	}
}

func TestParseScalar_Bytes(t *testing.T) {
	// "some-raw-bytes" base64-encoded.
	v, err := parseScalar(protoreflect.BytesKind, nil, `"c29tZS1yYXctYnl0ZXM="`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(v.Bytes(), []byte("some-raw-bytes")) {
		t.Errorf("got %q want some-raw-bytes", v.Bytes())
	}
}

func TestParseScalar_BytesInvalidBase64(t *testing.T) {
	if _, err := parseScalar(protoreflect.BytesKind, nil, `"not base64 @"`); err == nil {
		t.Errorf("want error for invalid base64, got nil")
	}
}

func TestParseScalar_Enum(t *testing.T) {
	enum := (&protoconfigv1.Test{}).ProtoReflect().Descriptor().Fields().ByName("enum_field").Enum()

	t.Run("by name", func(t *testing.T) {
		v, err := parseScalar(protoreflect.EnumKind, enum, "EXAMPLE_ENUM_EXAMPLE_VAL")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Enum() != 1 {
			t.Errorf("got %v want 1", v.Enum())
		}
	})

	t.Run("by number", func(t *testing.T) {
		v, err := parseScalar(protoreflect.EnumKind, enum, "1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if v.Enum() != 1 {
			t.Errorf("got %v want 1", v.Enum())
		}
	})

	t.Run("unknown name", func(t *testing.T) {
		if _, err := parseScalar(protoreflect.EnumKind, enum, "NO_SUCH_VALUE"); err == nil {
			t.Errorf("want error, got nil")
		}
	})

	t.Run("wrong case", func(t *testing.T) {
		if _, err := parseScalar(protoreflect.EnumKind, enum, "example_enum_example_val"); err == nil {
			t.Errorf("want error for lowercased enum name, got nil (protojson is case-sensitive)")
		}
	})
}

// --- messages ---

func TestParseMessage_Timestamp(t *testing.T) {
	// Timestamp is a WKT, protojson accepts RFC3339.
	ts := timestamppb.New(time.Date(2024, 2, 3, 17, 0, 0, 0, time.UTC))
	md := ts.ProtoReflect().Descriptor()
	v, err := parseMessage(md, `"2024-02-03T17:00:00Z"`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := v.Message().Interface().(*timestamppb.Timestamp).AsTime()
	if !got.Equal(ts.AsTime()) {
		t.Errorf("got %v want %v", got, ts.AsTime())
	}
}

func TestParseMessage_Nested(t *testing.T) {
	md := (&protoconfigv1.Nested2{}).ProtoReflect().Descriptor()
	v, err := parseMessage(md, `{"string_field": "hello", "not_updated_via_env": "keep"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	msg := v.Message().Interface().(*protoconfigv1.Nested2)
	if msg.StringField != "hello" || msg.NotUpdatedViaEnv != "keep" {
		t.Errorf("got %+v", msg)
	}
}

// --- whole-list / whole-map ---

func TestParseWholeList_Primitive(t *testing.T) {
	fd := (&protoconfigv1.Test{}).ProtoReflect().Descriptor().Fields().ByName("list_of_ints")
	v, err := parseWholeList(fd, `[1, 2, 3]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list := v.List()
	if list.Len() != 3 {
		t.Fatalf("len=%d want 3", list.Len())
	}
	for i, want := range []int32{1, 2, 3} {
		if int32(list.Get(i).Int()) != want {
			t.Errorf("list[%d]: got %v want %v", i, list.Get(i).Int(), want)
		}
	}
}

func TestParseWholeList_Message(t *testing.T) {
	fd := (&protoconfigv1.Test{}).ProtoReflect().Descriptor().Fields().ByName("repeated_nested_message")
	v, err := parseWholeList(fd, `[{"string_field":"a"},{"string_field":"b"}]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list := v.List()
	if list.Len() != 2 {
		t.Fatalf("len=%d want 2", list.Len())
	}
	got0 := list.Get(0).Message().Interface().(*protoconfigv1.Nested2).StringField
	if got0 != "a" {
		t.Errorf("list[0].string_field = %q want a", got0)
	}
}

func TestParseWholeMap_Primitive(t *testing.T) {
	fd := (&protoconfigv1.Test{}).ProtoReflect().Descriptor().Fields().ByName("primitive_map")
	v, err := parseWholeMap(fd, `{"key1":"val1","key2":"val2"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := v.Map()
	if m.Len() != 2 {
		t.Fatalf("len=%d want 2", m.Len())
	}
	got := m.Get(protoreflect.ValueOfString("key1").MapKey()).String()
	if got != "val1" {
		t.Errorf("map[key1]=%q want val1", got)
	}
}

func TestParseWholeMap_Message(t *testing.T) {
	fd := (&protoconfigv1.Test{}).ProtoReflect().Descriptor().Fields().ByName("string_to_map")
	v, err := parseWholeMap(fd, `{"k":{"string_field":"x"}}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := v.Map()
	got := m.Get(protoreflect.ValueOfString("k").MapKey()).Message().Interface().(*protoconfigv1.Nested2)
	if got.StringField != "x" {
		t.Errorf("got %q want x", got.StringField)
	}
}

// --- not a scalar ---

func TestParseScalar_MessageKindRejected(t *testing.T) {
	if _, err := parseScalar(protoreflect.MessageKind, nil, `{}`); err == nil {
		t.Errorf("want error for MessageKind in parseScalar")
	}
}
