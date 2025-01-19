package protoiter

import (
	"iter"

	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type opts struct {
	excludeUnpopulated bool
}

func ExcludeUnpopulated() func(*opts) {
	return func(o *opts) {
		o.excludeUnpopulated = true
	}
}

// Fields return an iterator that iterates over all fields.
func Fields(m protoreflect.Message, options ...func(*opts) opts) iter.Seq[protopath.Values] {
	var opts opts
	for _, opt := range options {
		opt(&opts)
	}
	return func(yield func(protopath.Values) bool) {
		_ = protorange.Range(m, func(v protopath.Values) error {
			if !yield(v) {
				return protorange.Terminate
			}

			if !opts.excludeUnpopulated {
				// In addition to the "normal" calls to the provided function f, also call it for unset primitive fields.
				leaf := v.Index(-1)

				// Inject for items in lists.
				if leaf.Step.Kind() == protopath.ListIndexStep {
					parent := v.Index(-2)
					// If item in the list is a message, inject its fields!
					if parent.Step.FieldDescriptor().Kind() == protoreflect.MessageKind {
						msg := parent.Value.List().Get(leaf.Step.ListIndex()).Message()
						for i := 0; i < msg.Descriptor().Fields().Len(); i++ {
							fd := msg.Descriptor().Fields().Get(i)
							if !msg.Has(fd) {
								if !yield(protopath.Values{
									Path:   append(v.Path, protopath.FieldAccess(fd)),
									Values: append(v.Values, fd.Default()),
								}) {
									return protorange.Terminate
								}
							}

						}
					}
					return nil
				}

				// Inject for items in maps.
				if leaf.Step.Kind() == protopath.MapIndexStep {
					parent := v.Index(-2)
					// If item in the list is a message, inject its fields!
					if parent.Step.FieldDescriptor().MapValue().Kind() == protoreflect.MessageKind {
						msg := parent.Value.Map().Get(leaf.Step.MapIndex()).Message()
						for i := 0; i < msg.Descriptor().Fields().Len(); i++ {
							fd := msg.Descriptor().Fields().Get(i)
							if !msg.Has(fd) {
								if !yield(protopath.Values{
									Path:   append(v.Path, protopath.FieldAccess(fd)),
									Values: append(v.Values, fd.Default()),
								}) {
									return protorange.Terminate
								}
							}

						}
					}
					return nil
				}

				// Inject for fields in messages.
				if (leaf.Step.Kind() == protopath.RootStep) ||
					(leaf.Step.Kind() == protopath.FieldAccessStep && leaf.Step.FieldDescriptor().Kind() == protoreflect.MessageKind &&
						!leaf.Step.FieldDescriptor().IsMap() && !leaf.Step.FieldDescriptor().IsList()) {
					msg := leaf.Value.Message()

					for i := 0; i < msg.Descriptor().Fields().Len(); i++ {
						fd := msg.Descriptor().Fields().Get(i)
						if !msg.Has(fd) {
							if !yield(protopath.Values{
								Path:   append(v.Path, protopath.FieldAccess(fd)),
								Values: append(v.Values, fd.Default()),
							}) {
								return protorange.Terminate
							}
						}
					}
				}
			}

			return nil
		})

	}
}
