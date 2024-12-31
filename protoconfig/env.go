package protoconfig

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
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
	default:
		return nil, fmt.Errorf("invalid field kind: %v", fd.Kind())
	}
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
	}
	return protoreflect.Value{}, nil
}
