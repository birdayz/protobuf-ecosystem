package protoconfig

import (
	"os"
	"testing"
	"time"

	. "github.com/birdayz/protobuf-ecosystem/pkg/pbgomega"
	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	t.Setenv("STRING_FIELD", "string-from-env")
	t.Setenv("BOOL_FIELD", "true")
	t.Setenv("MESSAGE_FIELD_BOOL_FIELD", "true")
	t.Setenv("ENUM_FIELD", "EXAMPLE_ENUM_EXAMPLE_VAL")
	t.Setenv("INT32_FIELD", "100")
	t.Setenv("SINT32_FIELD", "101")
	t.Setenv("OPTIONAL_INT32_FIELD", "90")
	t.Setenv("UINT32_FIELD", "102")
	t.Setenv("INT64_FIELD", "103")
	t.Setenv("SINT64_FIELD", "104")
	t.Setenv("UINT64_FIELD", "105")
	t.Setenv("SFIXED32_FIELD", "106")
	t.Setenv("FIXED32_FIELD", "107")
	t.Setenv("FLOAT_FIELD", "108.50")
	t.Setenv("SFIXED64_FIELD", "109")
	t.Setenv("FIXED64_FIELD", "110")
	t.Setenv("DOUBLE_FIELD", "111.50")
	t.Setenv("BYTES_FIELD", "c29tZS1yYXctYnl0ZXM=")

	// Nested same type
	t.Setenv("MESSAGE_FIELD_STRING_FIELD", "string-from-env")
	t.Setenv("MESSAGE_FIELD_BOOL_FIELD", "true")
	t.Setenv("MESSAGE_FIELD_MESSAGE_FIELD_BOOL_FIELD", "true")
	t.Setenv("MESSAGE_FIELD_ENUM_FIELD", "EXAMPLE_ENUM_EXAMPLE_VAL")
	t.Setenv("MESSAGE_FIELD_INT32_FIELD", "100")
	t.Setenv("MESSAGE_FIELD_SINT32_FIELD", "101")
	t.Setenv("MESSAGE_FIELD_OPTIONAL_INT32_FIELD", "90")
	t.Setenv("MESSAGE_FIELD_UINT32_FIELD", "102")
	t.Setenv("MESSAGE_FIELD_INT64_FIELD", "103")
	t.Setenv("MESSAGE_FIELD_SINT64_FIELD", "104")
	t.Setenv("MESSAGE_FIELD_UINT64_FIELD", "105")
	t.Setenv("MESSAGE_FIELD_SFIXED32_FIELD", "106")
	t.Setenv("MESSAGE_FIELD_FIXED32_FIELD", "107")
	t.Setenv("MESSAGE_FIELD_FLOAT_FIELD", "108.50")
	t.Setenv("MESSAGE_FIELD_SFIXED64_FIELD", "109")
	t.Setenv("MESSAGE_FIELD_FIXED64_FIELD", "110")
	t.Setenv("MESSAGE_FIELD_DOUBLE_FIELD", "111.50")
	t.Setenv("MESSAGE_FIELD_BYTES_FIELD", "c29tZS1yYXctYnl0ZXM=")

	// Other type
	t.Setenv("NESTED_MESSAGE_FIELD_STRING_FIELD", "nested-string-from-env")

	// List of messages
	t.Setenv("REPEATED_NESTED_MESSAGE_0_STRING_FIELD", "overridden")
	t.Setenv("REPEATED_NESTED_MESSAGE_1_STRING_FIELD", "overridden2")

	// Map string to message
	t.Setenv("STRING_TO_MAP_MY_MAP_KEY_STRING_FIELD", "some-string-field-very-nested")

	t.Setenv("LIST_OF_INTS", "[1,2,3]")
	t.Setenv("LIST_OF_STRINGS", `["first","second","third"]`)
	t.Setenv("LIST_OF_ENUMS", `["EXAMPLE_ENUM_EXAMPLE_VAL","EXAMPLE_ENUM_EXAMPLE_VAL"]`)

	// WKT
	t.Setenv("TIMESTAMP", `2000-02-03T17:00:53.123Z`)
	t.Setenv("TIMESTAMPS", `["2000-02-02T17:00:53.123Z","2000-02-03T17:00:53.123Z"]`)

	// NOT YET SUPPORTED
	// t.Setenv("TIMESTAMPS_0", `2000-02-02T17:00:53.123Z`)
	// t.Setenv("TIMESTAMPS_1", `2000-02-03T17:00:53.123Z`)

	t.Setenv("OVERRIDDEN_BY_ENV", `{"string_field":"string-field-val"}`)

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
		RepeatedNestedMessage: []*protoconfigv1.Nested2{
			{
				StringField:      "abc",
				NotUpdatedViaEnv: "abc1",
			},
			{
				StringField:      "abc2",
				NotUpdatedViaEnv: "abc2",
			},
		},
		StringToMap: map[string]*protoconfigv1.Nested2{
			"my_map_key": &protoconfigv1.Nested2{
				StringField:      "some-string-field-very-nested",
				NotUpdatedViaEnv: "some-string-field-very-nested-too",
			},
		},
		Timestamps: []*timestamppb.Timestamp{
			timestamppb.Now(),
		},
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
		RepeatedNestedMessage: []*protoconfigv1.Nested2{
			{
				StringField:      "overridden",
				NotUpdatedViaEnv: "abc1",
			},
			{
				StringField:      "overridden2",
				NotUpdatedViaEnv: "abc2",
			},
		},
		StringToMap: map[string]*protoconfigv1.Nested2{
			"my_map_key": &protoconfigv1.Nested2{
				StringField:      "some-string-field-very-nested",
				NotUpdatedViaEnv: "some-string-field-very-nested-too",
			},
		},
		ListOfInts:    []int32{1, 2, 3},
		ListOfStrings: []string{"first", "second", "third"},
		ListOfEnums: []protoconfigv1.Test_ExampleEnum{
			protoconfigv1.Test_EXAMPLE_ENUM_EXAMPLE_VAL,
			protoconfigv1.Test_EXAMPLE_ENUM_EXAMPLE_VAL,
		},
		Timestamp: timestamppb.New(time.Date(2000, time.February, 3, 17, 0, 53, 123*1000*1000, time.UTC)),
		Timestamps: []*timestamppb.Timestamp{
			timestamppb.New(time.Date(2000, time.February, 2, 17, 0, 53, 123*1000*1000, time.UTC)),
			timestamppb.New(time.Date(2000, time.February, 3, 17, 0, 53, 123*1000*1000, time.UTC)),
		},
		OverriddenByEnv: &protoconfigv1.Nested2{
			StringField: "string-field-val",
		},
	}))

	// TODO test for:
	// Nested message in nested message, configured by env
	// message in list of messages, configured by env
	// list of messages in list of messages, configured by env

}
