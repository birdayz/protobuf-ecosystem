package protoconfig

import (
	"os"
	"testing"

	. "github.com/birdayz/protobuf-ecosystem/pkg/pbgomega"
	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	. "github.com/onsi/gomega"
	"k8s.io/utils/ptr"
)

func TestLoadSimple(t *testing.T) {
	RegisterTestingT(t)
	yml := `string_field: abc`
	tf := t.TempDir() + "test.yaml"

	err := os.WriteFile(tf, []byte(yml), 0600)
	Expect(err).ToNot(HaveOccurred())

	cfg, err := Load(tf, &protoconfigv1.Test{
		StringField: "previous",
	})
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).To(EqualProto(&protoconfigv1.Test{
		StringField: "abc",
	}))

}

func TestLoadWithEnvVarOverride(t *testing.T) {
	RegisterTestingT(t)

	t.Setenv("STRING_FIELD_VAL", "string-from-env")
	t.Setenv("BOOL_FIELD_VAL", "true")
	t.Setenv("ENUM_FIELD_VAL", "EXAMPLE_ENUM_EXAMPLE_VAL")
	t.Setenv("INT32_FIELD_VAL", "100")
	t.Setenv("SINT32_FIELD_VAL", "101")
	t.Setenv("OPTIONAL_INT32_FIELD_VAL", "90")
	t.Setenv("NESTED_STRING_FIELD_VAL", "nested-string-from-env")
	t.Setenv("UINT32_FIELD_VAL", "102")
	t.Setenv("INT64_FIELD_VAL", "103")
	t.Setenv("SINT64_FIELD_VAL", "104")
	t.Setenv("UINT64_FIELD_VAL", "105")
	t.Setenv("SFIXED32_FIELD_VAL", "106")
	t.Setenv("FIXED32_FIELD_VAL", "107")
	t.Setenv("FLOAT_FIELD_VAL", "108.50")
	t.Setenv("SFIXED64_FIELD_VAL", "109")
	t.Setenv("FIXED64_FIELD_VAL", "110")
	t.Setenv("DOUBLE_FIELD_VAL", "111.50")
	t.Setenv("BYTES_FIELD_VAL", "c29tZS1yYXctYnl0ZXM=")
	yml :=
		`
  nested_message_field:
    not_updated_via_env: some-val-from-yaml
    string_field: also-from-yaml-but-overridden
  `
	tf := t.TempDir() + "test.yaml"

	err := os.WriteFile(tf, []byte(yml), 0600)
	Expect(err).ToNot(HaveOccurred())

	cfg, err := Load(tf, &protoconfigv1.Test{
		StringField:        "default-string",
		Int32Field:         1,
		Sint32Field:        1,
		MessageField:       &protoconfigv1.Test{},
		NestedMessageField: &protoconfigv1.Nested{},
	})
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).To(EqualProto(&protoconfigv1.Test{
		StringField:        "string-from-env",
		BoolField:          true,
		EnumField:          protoconfigv1.Test_EXAMPLE_ENUM_EXAMPLE_VAL,
		Int32Field:         100,
		Sint32Field:        101,
		OptionalInt32Field: ptr.To(int32(90)),
		Uint32Field:        102,
		Int64Field:         103,
		Sint64Field:        104,
		Uint64Field:        105,
		Sfixed32Field:      106,
		Fixed32Field:       107,
		FloatField:         108.5,
		Sfixed64Field:      109,
		Fixed64Field:       110,
		DoubleField:        111.5,
		BytesField:         []byte("some-raw-bytes"),
		MessageField: &protoconfigv1.Test{
			StringField:        "string-from-env",
			BoolField:          true,
			EnumField:          protoconfigv1.Test_EXAMPLE_ENUM_EXAMPLE_VAL,
			Int32Field:         100,
			Sint32Field:        101,
			OptionalInt32Field: ptr.To(int32(90)),
			Uint32Field:        102,
			Int64Field:         103,
			Sint64Field:        104,
			Uint64Field:        105,
			Sfixed32Field:      106,
			Fixed32Field:       107,
			FloatField:         108.5,
			Sfixed64Field:      109,
			Fixed64Field:       110,
			DoubleField:        111.5,
			BytesField:         []byte("some-raw-bytes"),
		},
		NestedMessageField: &protoconfigv1.Nested{
			StringField:      "nested-string-from-env",
			NotUpdatedViaEnv: "some-val-from-yaml",
		},
	}))
}
