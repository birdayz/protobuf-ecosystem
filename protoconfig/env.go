package protoconfig

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/yaml.v3"
)

func convertPrimitive[T interface {
	bool | int32 | int64 | uint32 | uint64 | float32 | float64 | string
}](strVal string) ([]protoreflect.Value, error) {
	var result []protoreflect.Value
	dec := json.NewDecoder(strings.NewReader(strVal))
	var items []T
	if err := dec.Decode(&items); err != nil {
		return nil, err
	}
	for _, item := range items {
		result = append(result, protoreflect.ValueOf(item))
	}

	return result, nil
}

// stringToValues converts a string that contains a JSON array into a slice of
// protoreflect.Value, considering the desired field descriptor.
func stringToValues(fd protoreflect.FieldDescriptor, strVal string) ([]protoreflect.Value, error) {
	switch fd.Kind() {
	case protoreflect.BoolKind:
		return convertPrimitive[bool](strVal)
	case protoreflect.EnumKind:
		// Special
		var items []string
		dec := json.NewDecoder(strings.NewReader(strVal))
		if err := dec.Decode(&items); err != nil {
			return nil, err
		}

		var results []protoreflect.Value
		for _, item := range items {
			results = append(results, protoreflect.ValueOfEnum(fd.Enum().Values().ByName(protoreflect.Name(item)).Number()))
		}
		return results, nil
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return convertPrimitive[int32](strVal)
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return convertPrimitive[uint32](strVal)
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return convertPrimitive[int64](strVal)
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return convertPrimitive[uint64](strVal)
	case protoreflect.FloatKind:
		return convertPrimitive[float32](strVal)
	case protoreflect.DoubleKind:
		return convertPrimitive[float64](strVal)
	case protoreflect.StringKind:
		return convertPrimitive[string](strVal)
	case protoreflect.BytesKind:
		var items []string
		dec := json.NewDecoder(strings.NewReader(strVal))
		if err := dec.Decode(&items); err != nil {
			return nil, err
		}

		var results []protoreflect.Value
		for _, item := range items {
			v, err := base64.StdEncoding.DecodeString(item)
			if err != nil {
				return nil, err
			}
			results = append(results, protoreflect.ValueOfBytes(v))
		}
		return results, nil
	case protoreflect.MessageKind: // Handle WKT
		switch fd.Message().FullName() {
		case "google.protobuf.Timestamp":
			var items []string
			dec := yaml.NewDecoder(strings.NewReader(strVal))
			if err := dec.Decode(&items); err != nil {
				return nil, err
			}

			var results []protoreflect.Value
			for _, item := range items {
				t, err := time.Parse(time.RFC3339Nano, item)
				if err != nil {
					return nil, fmt.Errorf("could not parse timestamps as RFC3339Nano: %w", err)
				}
				results = append(results, protoreflect.ValueOfMessage(timestamppb.New(t).ProtoReflect()))
			}
			return results, nil
		default:
			var result []protoreflect.Value
			md, err := protoregistry.GlobalTypes.FindMessageByName(fd.Message().FullName())
			if err != nil {
				return nil, err
			}

			var items []json.RawMessage
			dec := json.NewDecoder(strings.NewReader(strVal))
			if err := dec.Decode(&items); err != nil {
				return nil, err
			}

			for _, item := range items {
				msg := md.New()
				if err := protojson.Unmarshal(item, msg.Interface()); err != nil {
					return nil, err
				}
				result = append(result, protoreflect.ValueOfMessage(msg))
			}

			return result, nil
		}
	}
	return nil, fmt.Errorf("invalid field kind: %v", fd.Kind())
}

// stringToValue converts a string containing primitive values into their respective protoreflect.Value.
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
	case protoreflect.MessageKind: // Handle WKT
		switch fd.Message().FullName() {
		case "google.protobuf.Timestamp":
			t, err := time.Parse(time.RFC3339Nano, strVal)
			if err != nil {
				return protoreflect.Value{}, fmt.Errorf("could not parse timestamp as RFC3339Nano: %w", err)
			}
			return protoreflect.ValueOfMessage(timestamppb.New(t).ProtoReflect()), nil

		default:
			md, err := protoregistry.GlobalTypes.FindMessageByName(fd.Message().FullName())
			if err != nil {
				return protoreflect.Value{}, err
			}
			msg := md.New()
			if err := protojson.Unmarshal([]byte(strVal), msg.Interface()); err != nil {
				return protoreflect.Value{}, err
			}

			return protoreflect.ValueOfMessage(msg), nil

		}

		// switch fd.Message().FullName() {
		// case "google.protobuf.Timestamp":
		// 	t, err := time.Parse(time.RFC3339Nano, strVal)
		// 	if err != nil {
		// 		return protoreflect.Value{}, fmt.Errorf("could not parse timestamp as RFC3339Nano: %w", err)
		// 	}
		// 	return protoreflect.ValueOfMessage(timestamppb.New(t).ProtoReflect()), nil
		// }

	}
	return protoreflect.Value{}, fmt.Errorf("unsupported kind: %v", fd.Kind())
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
