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

		if leaf.Step.Kind() != protopath.FieldAccessStep {
			return nil
		}

		// Works only if FieldAccess
		var fd protoreflect.FieldDescriptor
		if leaf.Step.Kind() == protopath.FieldAccessStep {
			fd = leaf.Step.FieldDescriptor()
		}

		envKey := pathToEnvKey(v.Path)

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
		} else if fd.IsMap() {
			if envVal, ok := os.LookupEnv(envKey); ok && envVal != "" {
				panic("Map is unsupported")
			}
			// --- Map
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

			// This is a message field.
			if !leaf.Value.IsValid() {
				e := os.Environ()

				var sub []string
				for _, s := range e {
					s = strings.Split(s, "=")[0]
					if strings.HasPrefix(s, envKey+"_") {
						sub = append(sub, s)
					}
				}

				if len(sub) > 0 {
					md, err := protoregistry.GlobalTypes.FindMessageByName(fd.Message().FullName())
					if err != nil {
						return err
					}
					parent.Value.Message().Set(fd, protoreflect.ValueOfMessage(md.New()))
				}
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
