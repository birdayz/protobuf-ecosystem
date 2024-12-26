package protoconfig

import (
	"os"
	"testing"

	"github.com/birdayz/protobuf-ecosystem/pkg/pbgomega"
	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	. "github.com/onsi/gomega"
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
	Expect(cfg).To(pbgomega.EqualProto(&protoconfigv1.Test{
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
	yml := `test: fromConfig`
	tf := t.TempDir() + "test.yaml"

	err := os.WriteFile(tf, []byte(yml), 0600)
	Expect(err).ToNot(HaveOccurred())

	cfg, err := Load(tf, &protoconfigv1.Test{
		StringField: "default-string",
		Int32Field:  1,
		Sint32Field: 1,
	})
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).To(pbgomega.EqualProto(&protoconfigv1.Test{
		StringField: "string-from-env",
		BoolField:   true,
		EnumField:   protoconfigv1.Test_EXAMPLE_ENUM_EXAMPLE_VAL,
		Int32Field:  100,
		Sint32Field: 101,
	}))
}
