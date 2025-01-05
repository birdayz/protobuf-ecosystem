package protoconfig

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"buf.build/go/protoyaml"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// RangeAllFields is a decorator of protorange.Range, that also includes unset primitive fields.
// No ordering of fields within a message is guaranteed.
func RangeAllFields(m protoreflect.Message, f func(protopath.Values) error) error {
	return protorange.Range(m, func(v protopath.Values) error {
		if err := f(v); err != nil {
			return err
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
						if err := f(protopath.Values{
							Path:   append(v.Path, protopath.FieldAccess(fd)),
							Values: append(v.Values, fd.Default()),
						}); err != nil {
							return err
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
					if err := f(protopath.Values{
						Path:   append(v.Path, protopath.FieldAccess(fd)),
						Values: append(v.Values, fd.Default()),
					}); err != nil {
						return err
					}
				}
			}
		}

		// TODO: inject for messages in maps.
		return nil
	})
}

func Load[T proto.Message](path string, defaults T) (T, error) {
	// Clone defaults, we don't want to suprise callers by modifying their
	// "static" defaults variable.
	cfg := proto.Clone(defaults)

	// 1. Load YAML file.
	bytez, err := os.ReadFile(path)
	if err != nil {
		return defaults, err
	}

	fromFile := defaults.ProtoReflect().New().Interface()
	err = protoyaml.UnmarshalOptions{
		DiscardUnknown: true, // Lenient parsing for forwards compatibility.
	}.Unmarshal(bytez, fromFile)
	if err != nil {
		return *new(T), fmt.Errorf("failed to unmarshal YAML configuration file %s: %w", path, err)
	}

	// Merge YAML content on top of defaults.
	proto.Merge(cfg, fromFile)

	// 2. Load additional env vars on top.

	err = RangeAllFields(cfg.ProtoReflect(), func(v protopath.Values) error {
		leaf := v.Index(-1)
		parent := v.Index(-2)

		// Works only if FieldAccess
		var fd protoreflect.FieldDescriptor
		if leaf.Step.Kind() == protopath.FieldAccessStep {
			fd = leaf.Step.FieldDescriptor()
		} else if leaf.Step.Kind() == protopath.ListIndexStep {
			fd = parent.Step.FieldDescriptor()
		} else {
			return nil
		}

		envKey := pathToEnvKey(v.Path)

		subFields := func() []string {
			e := os.Environ()

			var sub []string
			for _, key := range e {
				key = strings.Split(key, "=")[0]

				if strings.HasPrefix(key, envKey+"_") {
					sub = append(sub, key)
				}
			}
			return sub
		}

		if fd.IsList() {
			if envVal, ok := os.LookupEnv(envKey); ok {
				// Found exact match, that is supposed to provide the entire value for this field.
				l := parent.Value.Message().Mutable(fd).List()
				l.Truncate(0)

				values, err := stringToValues(fd, envVal)
				if err != nil {
					return fmt.Errorf("failed to convert list value: %w", err)
				}

				for _, item := range values {
					l.Append(item)
				}

				return nil
			}

			// Check all envs below "this one" (trim prefix).
			// _1, _2,
			var highestInt *int
			for _, subField := range subFields() {
				trimmed := strings.Split(strings.TrimPrefix(subField, envKey+"_"), "_")[0]
				if number, err := strconv.Atoi(trimmed); err == nil {
					if highestInt == nil || number > *highestInt {
						highestInt = &number
					}
				}
			}
			if highestInt != nil {
				size := *highestInt + 1 // Size is index + 1
				l := parent.Value.Message().Mutable(fd).List()

				for i := l.Len(); i < size; i++ {
					if fd.Kind() == protoreflect.MessageKind {
						md, err := protoregistry.GlobalTypes.FindMessageByName(fd.Message().FullName())
						if err != nil {
							return err
						}
						l.Append(protoreflect.ValueOfMessage(md.New()))
					} else if fd.HasDefault() {
						l.Append(protoreflect.ValueOf(fd.Default()))
					}
				}
			}

		} else if fd.IsMap() {
			// --- Map
			if envVal, ok := os.LookupEnv(envKey); ok && envVal != "" {
				panic("Map is unsupported")
			}
		} else {
			// -- Ordinary Field - no map/list.
			envVal, ok := os.LookupEnv(envKey)
			if ok && envVal != "" {
				val, err := stringToValue(fd, envVal)
				if err != nil {
					return err
				}
				parent.Value.Message().Set(fd, val)
				return nil
			}

			// This is a message field (default value is "not present").
			// If not present, set it - but only if an env var wants to set fields of this message.
			// We don't want to set nested messages to a value by default, because
			// this can lead to infinite recursion.
			// By initializing this field, protorange.Range will consider this field, and call
			// our range function for it.
			if fd.Kind() == protoreflect.MessageKind && !leaf.Value.IsValid() && len(subFields()) > 0 {
				md, err := protoregistry.GlobalTypes.FindMessageByName(fd.Message().FullName())
				if err != nil {
					return err
				}
				parent.Value.Message().Set(fd, protoreflect.ValueOfMessage(md.New()))
			}

		}

		return nil
	})

	return cfg.(T), err
}

func pathToEnvKey(path protopath.Path) string {
	var value strings.Builder

	for _, step := range path {
		switch step.Kind() {
		case protopath.RootStep:
		// Do nothing
		case protopath.FieldAccessStep:
			// Skip _ prefix if we haven't written a field yet.
			if value.Len() > 0 {
				_, _ = value.WriteString("_")
			}
			_, _ = value.WriteString(strings.ToUpper(string(step.FieldDescriptor().Name())))
		case protopath.MapIndexStep:
			// Always add _ separator, because there's always a precending field access.
			// Map indexing can't be the topmost item.
			_, _ = value.WriteString("_")
			_, _ = value.WriteString(strings.ToUpper(step.MapIndex().String()))
		case protopath.ListIndexStep:
			// Always add _ separator, because there's always a precending field access.
			// Map indexing can't be the topmost item.
			_, _ = value.WriteString("_")
			_, _ = value.WriteString(strings.ToUpper(strconv.FormatInt(int64(step.ListIndex()), 10)))
		}
	}
	result := value.String()
	return result
}
