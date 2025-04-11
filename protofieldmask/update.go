package protofieldmask

import (
	"iter"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Compare(l, r proto.Message) iter.Seq2[protopath.Values, protopath.Values] {
	left := l.ProtoReflect()
	right := r.ProtoReflect()
	return func(yield func(l, r protopath.Values) bool) {
		l := protopath.Values{
			Path:   []protopath.Step{protopath.Root(left.Descriptor())},
			Values: []protoreflect.Value{protoreflect.ValueOfMessage(left)},
		}
		r := protopath.Values{
			Path:   []protopath.Step{protopath.Root(right.Descriptor())},
			Values: []protoreflect.Value{protoreflect.ValueOfMessage(right)},
		}
		processValues(l, r, yield)
	}
}

func processValues(left, right protopath.Values, yield func(l, r protopath.Values) bool) {
	// If both primitive, just return them. Done.
	leftLeaf := left.Index(-1)
	rightLeaf := right.Index(-1)

	// Determine how we're going to iterate over this.
	// Either iterate over both, and match left / right.
	// If that is not possible, because the value/fd is empty on one side, just present nil for the missing side (like an outer join basically).
	var iterator iter.Seq2[protoreflect.FieldDescriptor, protoreflect.FieldDescriptor]

	if leftLeaf.Step.Kind() == protopath.RootStep && rightLeaf.Step.Kind() == protopath.RootStep {
		iterator = IterFieldDescriptors(
			leftLeaf.Value.Message().Descriptor(),
			rightLeaf.Value.Message().Descriptor(),
			MatcherByNumber)
	} else {
		fdLeft := leftLeaf.Step.FieldDescriptor()
		fdRight := rightLeaf.Step.FieldDescriptor()
		if fdLeft.Kind() == protoreflect.MessageKind && fdRight.Kind() == protoreflect.MessageKind {
			iterator = IterFieldDescriptors(
				leftLeaf.Value.Message().Descriptor(),
				rightLeaf.Value.Message().Descriptor(),
				MatcherByNumber)
		} else if fdLeft.Kind() == protoreflect.MessageKind && fdRight.Kind() != protoreflect.MessageKind {
			iterator = WrapStaticNil(IterFieldDescriptor(leftLeaf.Value.Message().Descriptor()))
		} else {
			iterator = WrapStaticNil(IterFieldDescriptor(rightLeaf.Value.Message().Descriptor()))
		}
	}

	for fdLeft, fdRight := range iterator {

		// Vars for this iteration
		var (
			valueLeft, valueRight protoreflect.Value
			pathLeft, pathRight   protopath.Step
		)

		if fdLeft != nil {
			leftMessage := left.Index(-1).Value.Message()
			valueLeft = leftMessage.Get(fdLeft)
			pathLeft = protopath.FieldAccess(fdLeft)
		}

		if fdRight != nil {
			rightMessage := right.Index(-1).Value.Message()
			valueRight = rightMessage.Get(fdRight)
			pathRight = protopath.FieldAccess(fdRight)
		}

		// Always send, even if message
		yld := yield(
			protopath.Values{
				Path:   append(left.Path, pathLeft),
				Values: append(left.Values, valueLeft),
			}, protopath.Values{
				Path:   append(right.Path, pathRight),
				Values: append(right.Values, valueRight),
			},
		)

		// Only recurse, if at least on of the two is recurse-able (for now: Message)

		if (fdLeft != nil && fdLeft.Kind() == protoreflect.MessageKind && !fdLeft.IsMap() && !fdLeft.IsList() && valueLeft.IsValid()) ||
			(fdRight != nil && fdRight.Kind() == protoreflect.MessageKind && !fdRight.IsMap() && !fdRight.IsList() && valueRight.IsValid()) {
			processValues(protopath.Values{
				Path:   append(left.Path, pathLeft),
				Values: append(left.Values, valueLeft),
			}, protopath.Values{
				Path:   append(right.Path, pathRight),
				Values: append(right.Values, valueRight),
			}, yield)
		}

		if !yld {
			return
		}
	}
}
