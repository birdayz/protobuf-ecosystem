package protoconfig

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"buf.build/go/protoyaml"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protorange"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// RangeAllFields is a decorator of protorange.Range, that also includes unset primitive fields.
// No ordering of fields within a message is guaranteed.
func RangeAllFields(m protoreflect.Message, f func(protopath.Values) error) error {
	return protorange.Range(m, func(v protopath.Values) error {
		if err := f(v); err != nil {
			return err
		}
		leaf := v.Index(-1)
		// In addition to the "normal" calls to the provided function f, also call it for unset primitive fields.
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
	// if err := recurse(nil, cfg.ProtoReflect()); err != nil {
	// 	return cfg.(T), err
	// }

	err = RangeAllFields(cfg.ProtoReflect(), func(v protopath.Values) error {
		leaf := v.Index(-1)
		parent := v.Index(-2)

		// We're only interested in primitive fields.
		if leaf.Step.Kind() == protopath.FieldAccessStep && leaf.Step.FieldDescriptor().Kind() == protoreflect.MessageKind {
			return nil
		}

		// Lookup env var for this field.
		// It is the uppercase concatenated field path.
		// For map, the map key is used.
		// For list, the index is used.
		envKey := pathToEnvKey(v.Path)

		// Get override
		envVal, ok := os.LookupEnv(envKey)
		if !ok || envVal == "" {
			return nil
		}

		strVal, err := stringToValue(leaf.Step.FieldDescriptor(), envVal)
		if err != nil {
			return err
		}

		// Store the redacted string back into the message.
		switch leaf.Step.Kind() {
		case protopath.FieldAccessStep:
			m := parent.Value.Message()
			fd := leaf.Step.FieldDescriptor()
			m.Set(fd, strVal)
		case protopath.ListIndexStep:
			ls := parent.Value.List()
			i := leaf.Step.ListIndex()
			ls.Set(i, strVal)
		case protopath.MapIndexStep:
			ms := parent.Value.Map()
			k := leaf.Step.MapIndex()
			ms.Set(k, strVal)
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

func stringToValue(fd protoreflect.FieldDescriptor, strVal string) (protoreflect.Value, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		if strVal != "true" && strVal != "false" {
			return protoreflect.Value{}, fmt.Errorf("could not convert string %s to bool, unsupported value", strVal)
		}
		return protoreflect.ValueOfBool(strVal == "true"), nil
	case protoreflect.EnumKind:
		enumVal := fd.Enum().Values().ByName(protoreflect.Name(strVal))
		if enumVal == nil {
			return protoreflect.Value{}, fmt.Errorf("could not convert string %s to enum: enum value does not exist", strVal)
		}
		return protoreflect.ValueOfEnum(enumVal.Number()), nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		v, err := strconv.ParseInt(strVal, 10, 32)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("could not convert string %s to int32: %w", strVal, err)
		}
		return protoreflect.ValueOfInt32(int32(v)), nil
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		v, err := strconv.ParseUint(strVal, 10, 32)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("could not convert string %s to uint32: %w", strVal, err)
		}
		return protoreflect.ValueOfUint32(uint32(v)), nil
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		v, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("could not convert string %s to int64: %w", strVal, err)
		}
		return protoreflect.ValueOfInt64(int64(v)), nil
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		v, err := strconv.ParseUint(strVal, 10, 64)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("could not convert string %s to uint64: %w", strVal, err)
		}
		return protoreflect.ValueOfUint64(uint64(v)), nil
	case protoreflect.FloatKind:
		v, err := strconv.ParseFloat(strVal, 32)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("could not convert string %s to float: %w", strVal, err)
		}
		return protoreflect.ValueOfFloat32(float32(v)), nil
	case protoreflect.DoubleKind:
		v, err := strconv.ParseFloat(strVal, 64)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("could not convert string %s to float: %w", strVal, err)
		}
		return protoreflect.ValueOfFloat64(float64(v)), nil
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(strVal), nil
	case protoreflect.BytesKind:
		v, err := base64.StdEncoding.DecodeString(strVal)
		if err != nil {
			return protoreflect.Value{}, fmt.Errorf("could not base64 decode string %s: %w", strVal, err)
		}
		return protoreflect.ValueOfBytes(v), nil
	}
	return protoreflect.Value{}, nil
}
