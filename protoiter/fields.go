package protoiter

import (
	"iter"

	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Fields return an iterator that iterates over all fields.
func Fields(m protoreflect.Message) iter.Seq[protopath.Values] {
	return func(f func(protopath.Values) bool) {
		_ = protorange.Range(m, func(v protopath.Values) error {
			if !f(v) {
				return protorange.Terminate
			}

			// In addition to the "normal" calls to the provided function f, also call it for unset primitive fields.
			leaf := v.Index(-1)

			// Inject for messages in lists.
			if leaf.Step.Kind() == protopath.ListIndexStep {

				parent := v.Index(-2)
				if parent.Step.FieldDescriptor().Kind() == protoreflect.MessageKind {
					msg := parent.Value.List().Get(leaf.Step.ListIndex()).Message()
					for i := 0; i < msg.Descriptor().Fields().Len(); i++ {
						fd := msg.Descriptor().Fields().Get(i)
						if !msg.Has(fd) {
							if !f(protopath.Values{
								Path:   append(v.Path, protopath.FieldAccess(fd)),
								Values: append(v.Values, fd.Default()),
							}) {
								return protorange.Terminate
							}
						}

					}
				}
			}

			// Inject for messages in message-fields.
			if (leaf.Step.Kind() == protopath.RootStep) || (leaf.Step.Kind() == protopath.FieldAccessStep && leaf.Step.FieldDescriptor().Kind() == protoreflect.MessageKind && !leaf.Step.FieldDescriptor().IsMap() && !leaf.Step.FieldDescriptor().IsList()) {
				msg := leaf.Value.Message()

				for i := 0; i < msg.Descriptor().Fields().Len(); i++ {
					fd := msg.Descriptor().Fields().Get(i)
					if !msg.Has(fd) {
						if !f(protopath.Values{
							Path:   append(v.Path, protopath.FieldAccess(fd)),
							Values: append(v.Values, fd.Default()),
						}) {
							return protorange.Terminate
						}
					}
				}
			}

			// TODO: inject for messages in maps.
			return nil
		})

	}
}
