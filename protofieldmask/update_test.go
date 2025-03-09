package protofieldmask

import (
	"fmt"
	"testing"

	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Use protorange without recursion.

// Proceed "by hand", once both old and new finished with the step.

func TestSimple(t *testing.T) {
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

	// var nok wrapperspb.StringValue
	// nok.Value = "wrong"

	// var m simplev1.ExampleTable
	// m.Id = uuid.NewString()
	// m.SomeData = "test123"
	// m.Timestamp = timestamppb.Now()
	//
	// x := &m

	//
	for l, r := range Compare(left.ProtoReflect(), right.ProtoReflect()) {
		fmt.Println(l.Path, r.Path, l.Index(-1).Value, r.Index(-1).Value)
		if r.Index(-1).Step.FieldDescriptor().Name() == "string_field" {
			// This makes total sense, but is quite difficult to use if not familiar with protopath.
			// Maybe add a helper that allows setting/replacing the value.
			r.Index(-2).Value.Message().Set(r.Index(-1).Step.FieldDescriptor(), protoreflect.ValueOfString("replacement"))
		}
	}

	jzon, err := protojson.MarshalOptions{
		Multiline:         true,
		EmitDefaultValues: true,
	}.Marshal(&right)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jzon))

}
