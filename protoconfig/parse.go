package protoconfig

import (
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

	proto.Merge(cfg, fromFile)

	// 2. Load additional env vars on top.
	if err := recurse(cfg.ProtoReflect()); err != nil {
		return cfg.(T), err
	}

	return cfg.(T), nil
}

func recurse(m protoreflect.Message) error {

	for i := 0; i < m.Descriptor().Fields().Len(); i++ {
		fd := m.Descriptor().Fields().Get(i)

		opts := proto.GetExtension(fd.Options().(*descriptorpb.FieldOptions), protoconfigv1.E_Options).(*protoconfigv1.FieldOptions)
		if opts == nil || opts.Env == nil {
			continue
		}

		envVal, ok := os.LookupEnv(*opts.Env)
		if !ok || envVal == "" {
			continue
		}

		switch fd.Kind() {
		case protoreflect.BoolKind:
			if envVal != "true" && envVal != "false" {
				return fmt.Errorf("could not convert env var %s with value %s to bool, unsupported value", *opts.Env, envVal)
			}
			m.Set(fd, protoreflect.ValueOfBool(envVal == "true"))
		case protoreflect.EnumKind:
			enumVal := fd.Enum().Values().ByName(protoreflect.Name(envVal))
			if enumVal == nil {
				return fmt.Errorf("could not convert env var %s with value %s, because enum value does not exist", *opts.Env, envVal)
			}
			m.Set(fd, protoreflect.ValueOfEnum(enumVal.Number()))
		case protoreflect.StringKind:
			m.Set(fd, protoreflect.ValueOfString(envVal))
		case protoreflect.Int32Kind, protoreflect.Sint32Kind:
			v, err := strconv.ParseInt(envVal, 10, 32)
			if err != nil {
				return fmt.Errorf("could not convert env var %s with value %s to int32: %w", *opts.Env, envVal, err)
			}
			m.Set(fd, protoreflect.ValueOfInt32(int32(v)))
		}

	}

	return nil
}
