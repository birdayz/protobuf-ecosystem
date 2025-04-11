package protofieldmask

import (
	"fmt"
	"testing"

	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	. "github.com/onsi/gomega"

	"google.golang.org/protobuf/reflect/protopath"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

// Use protorange without recursion.

// Proceed "by hand", once both old and new finished with the step.

func TestDifferentTypes(t *testing.T) {
	left := MyResource{
		Test: "old",
		Nested: &MyResource_Nested{
			NestedString: "nezted-str-orig",
		},
	}

	right := MyResourceUpdate{
		Test: "new",
		Nested: &MyResourceUpdate_Nested{
			NestedString: "new-nested-string",
		},
	}

	for l, r := range Compare(&left, &right) {
		fmt.Println("===================================================")
		fmt.Println("Path", l.Path.String())
		fmt.Println("L Step", l.Index(-1).Step, l.Index(-1).Step.FieldDescriptor())
		fmt.Println("R step", r.Index(-1).Step, r.Index(-1).Step.FieldDescriptor())
	}

}

func TestSimple(t *testing.T) {
	RegisterTestingT(t)
	var left protoconfigv1.Test
	left.StringField = "abc"
	left.NestedMessageField = &protoconfigv1.Nested{
		StringField: "nested-string-left",
	}

	var right protoconfigv1.Test
	right.StringField = "def"
	right.NestedMessageField = &protoconfigv1.Nested{
		StringField: "nested-string-right",
	}

	tests := []struct {
		name string
		lVal protoreflect.Value
		rVal protopath.Values
	}{
		{
			lVal: protoreflect.Value{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for l, r := range Compare(&left, &right) {
				fmt.Println(l.Index(-1).Value.Interface())
				fmt.Println(r.Index(-1).Value.Interface())
			}

		})
	}

}
