package protoconfig

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"buf.build/go/protoyaml"
	"github.com/birdayz/protobuf-ecosystem/protoiter"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
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

	for v := range protoiter.Fields(cfg.ProtoReflect()) {
		leaf := v.Index(-1)
		parent := v.Index(-2)

		// Works only if FieldAccess
		var fd protoreflect.FieldDescriptor
		if leaf.Step.Kind() == protopath.FieldAccessStep {
			fd = leaf.Step.FieldDescriptor()
		} else if leaf.Step.Kind() == protopath.ListIndexStep {
			fd = parent.Step.FieldDescriptor()
		} else if leaf.Step.Kind() == protopath.MapIndexStep {
			fd = parent.Step.FieldDescriptor()
		} else {
			continue
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
				v, err := stringToProtoValue(fd, envVal)
				if err != nil {
					return *new(T), fmt.Errorf("failed to convert list value: %w", err)
				}
				parent.Value.Message().Set(fd, v)
				continue
			}

			var highestIndex *int
			for _, subField := range subFields() {
				trimmed := strings.Split(strings.TrimPrefix(subField, envKey+"_"), "_")[0]
				if number, err := strconv.Atoi(trimmed); err == nil {
					if highestIndex == nil || number > *highestIndex {
						highestIndex = &number
					}
				}
			}
			if highestIndex != nil {
				size := *highestIndex + 1 // Size is index + 1
				l := parent.Value.Message().Mutable(fd).List()

				// For every message index found in env vars, that is a message beyond the current size, initialize a list entry with default value.
				// Protorange will call us for this field.
				for i := l.Len(); i < size; i++ {
					if fd.Kind() == protoreflect.MessageKind {
						md, err := protoregistry.GlobalTypes.FindMessageByName(fd.Message().FullName())
						if err != nil {
							return *new(T), err
						}
						l.Append(protoreflect.ValueOfMessage(md.New()))
					} else if fd.HasDefault() {
						l.Append(protoreflect.ValueOf(fd.Default()))
					}
				}
			}

		} else if fd.IsMap() {
			// --- Map

			if envVal, ok := os.LookupEnv(envKey); ok {
				mp, err := stringToProtoValue(fd, envVal)
				if err != nil {
					return *new(T), fmt.Errorf("failed to parse map value: %w", err)
				}
				parent.Value.Message().Set(fd, mp)
			} else {
				// TODO: this still needs to be implemented.
				// Support override if nested env vars are set.
				// Initializing map fields should be enough; as with list, protorange would subsequently visit these and the normal "direct" access code would kick in.
				// If primitive, we can just take the entire suffix as the map key
				// If message, split off field names of the sub-message. take the prefix before the field-name-suffix; that is the map key.
				// Check if nested..
			}
		} else {
			// -- Ordinary Field - no map/list.
			envVal, ok := os.LookupEnv(envKey)
			if ok && envVal != "" {
				val, err := stringToProtoValue(fd, envVal)
				if err != nil {
					return *new(T), err
				}
				parent.Value.Message().Set(fd, val)
				continue
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
					return *new(T), err
				}
				parent.Value.Message().Set(fd, protoreflect.ValueOfMessage(md.New()))
			}
		}
	}

	return cfg.(T), err
}
