package protofieldmask

import (
	"iter"

	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Compare(left, right protoreflect.Message) iter.Seq2[protopath.Values, protopath.Values] {
	if left.Descriptor().FullName() != right.Descriptor().FullName() {
		panic("different types are not yet supported")
	}
	return func(yield func(l, r protopath.Values) bool) {
		l := protopath.Values{
			Path:   []protopath.Step{protopath.Root(left.Descriptor())},
			Values: []protoreflect.Value{protoreflect.ValueOfMessage(left)},
		}
		r := protopath.Values{
			Path:   []protopath.Step{protopath.Root(right.Descriptor())},
			Values: []protoreflect.Value{protoreflect.ValueOfMessage(right)},
		}

		iterMessage(l, r, yield)

	}
}

func iterMessage(left, right protopath.Values, yield func(l, r protopath.Values) bool) {
	l := left.Index(-1).Value.Message()
	r := right.Index(-1).Value.Message()

	fdsLeft := l.Descriptor().Fields()
	fdsRight := r.Descriptor().Fields()

	// TODO support different number of fields. Eg. Resource may have more fields than ResourceCreate or ResourceUpdate.
	n := max(fdsLeft.Len(), fdsRight.Len())

	for i := 0; i < n; i++ {
		fdLeft := fdsLeft.Get(i)
		fdRight := fdsRight.Get(i)

		leftStep := protopath.FieldAccess(fdLeft)
		rightStep := protopath.FieldAccess(fdRight)

		leftValue := l.Get(fdLeft)
		rightValue := r.Get(fdRight)

		// Ensure these are identical.
		if leftStep.FieldDescriptor().Kind() != rightStep.FieldDescriptor().Kind() {
			panic("not ok")
		}

		// If it's message, recurse into it, instead of returning it directly.
		// TODO: we can additionally emit an event for the field itself
		if leftStep.FieldDescriptor().Kind() == protoreflect.MessageKind &&
			!fdLeft.IsMap() && !fdRight.IsMap() &&
			!fdLeft.IsList() && !fdRight.IsList() &&
			leftValue.Message().IsValid() && rightValue.Message().IsValid() {
			iterMessage(protopath.Values{
				Path:   append(left.Path, protopath.FieldAccess(fdLeft)),
				Values: append(left.Values, leftValue),
			}, protopath.Values{
				Path:   append(right.Path, protopath.FieldAccess(fdLeft)),
				Values: append(right.Values, rightValue),
			}, yield)
			continue
		}

		// TODO: Fill up for any mismatches:
		// Different number of array items, map items
		// For message, don't do anything; field indexes must point to identical types

		// TODO: filter out fields not covered by fieldmask

		if !yield(
			protopath.Values{
				Path:   append(left.Path, protopath.FieldAccess(fdLeft)),
				Values: append(left.Values, leftValue),
			}, protopath.Values{
				Path:   append(right.Path, protopath.FieldAccess(fdLeft)),
				Values: append(right.Values, rightValue),
			},
		) {
			return
		}

	}

}
