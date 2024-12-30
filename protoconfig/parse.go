package protoconfig

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

	"buf.build/go/protoyaml"
	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

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
	if err := recurse(nil, cfg.ProtoReflect()); err != nil {
		return cfg.(T), err
	}

	return cfg.(T), nil
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

func recurse(path []protoreflect.Name, m protoreflect.Message) error {
	for i := 0; i < m.Descriptor().Fields().Len(); i++ {
		fd := m.Descriptor().Fields().Get(i)

		fieldPath := append(path, fd.Name())
		fmt.Println("p", fieldPath)

		// TODO

		// If is-repeated

		// If is-map

		// --> Do the below for all items.
		// --> For map, list, if of message type, only existing entries will be handled
		// --> For primitives, JSON unmarshal will be used.
		// --> We might want to have a protojson-parse for value.

		// Recurse if: is message, list+message, map+message-value.
		if fd.Kind() == protoreflect.MessageKind {

			if m.Has(fd) {
				if fd.IsMap() && fd.MapValue().Kind() == protoreflect.MessageKind {
					mp := m.Get(fd).Map()
					mp.Range(func(mk protoreflect.MapKey, v protoreflect.Value) bool {
						p := append(fieldPath, protoreflect.Name(mk.String()))
						if err := recurse(p, v.Message()); err != nil {
							panic(err) // TODO
							// return fmt.Errorf("failed to recurse into field %s: %w", fd.Name(), err)
						}
						return true
					})
				} else if fd.IsList() {
					l := m.Get(fd).List()

					for i := 0; i < l.Len(); i++ {
						p := append(fieldPath, protoreflect.Name(fmt.Sprintf("%d", i)))
						if err := recurse(p, l.Get(i).Message()); err != nil {
							return fmt.Errorf("failed to recurse into field %s: %w", fd.Name(), err)
						}
					}
				} else {
					if err := recurse(append(path, fd.Name()), m.Get(fd).Message()); err != nil {
						return fmt.Errorf("failed to recurse into field %s: %w", fd.Name(), err)
					}
				}
			}
		}

		// TODO: now, we've got values we want to directly read from env vars.

		opts := proto.GetExtension(fd.Options().(*descriptorpb.FieldOptions), protoconfigv1.E_Options).(*protoconfigv1.FieldOptions)
		if opts == nil || opts.Env == nil {
			continue
		}

		envVal, ok := os.LookupEnv(*opts.Env)
		if !ok || envVal == "" {
			continue
		}

		valToSet, err := stringToValue(fd, envVal)
		if err != nil {
			return err
		}

		m.Set(fd, valToSet)

	}

	return nil
}
