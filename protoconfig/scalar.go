package protoconfig

import (
	"fmt"
	"strconv"

	"buf.build/go/protoyaml"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// parseScalar parses envVal for a non-message scalar kind using protoyaml
// semantics — the same rules that apply to a YAML config file.
// enum is required when kind == EnumKind.
func parseScalar(kind protoreflect.Kind, enum protoreflect.EnumDescriptor, envVal string) (protoreflect.Value, error) {
	var wrapper proto.Message
	switch kind {
	case protoreflect.BoolKind:
		wrapper = &wrapperspb.BoolValue{}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		wrapper = &wrapperspb.Int32Value{}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		wrapper = &wrapperspb.UInt32Value{}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		wrapper = &wrapperspb.Int64Value{}
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		wrapper = &wrapperspb.UInt64Value{}
	case protoreflect.FloatKind:
		wrapper = &wrapperspb.FloatValue{}
	case protoreflect.DoubleKind:
		wrapper = &wrapperspb.DoubleValue{}
	case protoreflect.StringKind:
		wrapper = &wrapperspb.StringValue{}
	case protoreflect.BytesKind:
		wrapper = &wrapperspb.BytesValue{}
	case protoreflect.EnumKind:
		if enum == nil {
			return protoreflect.Value{}, fmt.Errorf("parseScalar: EnumKind requires enum descriptor")
		}
		return parseEnum(enum, envVal)
	case protoreflect.MessageKind, protoreflect.GroupKind:
		return protoreflect.Value{}, fmt.Errorf("parseScalar: %v is not a scalar kind, use parseMessage", kind)
	default:
		return protoreflect.Value{}, fmt.Errorf("parseScalar: unsupported kind %v", kind)
	}
	if err := safeProtoyamlUnmarshal(protoyaml.UnmarshalOptions{}, []byte(envVal), wrapper); err != nil {
		return protoreflect.Value{}, err
	}
	return valueFromWrapper(wrapper), nil
}

func valueFromWrapper(w proto.Message) protoreflect.Value {
	switch v := w.(type) {
	case *wrapperspb.BoolValue:
		return protoreflect.ValueOfBool(v.Value)
	case *wrapperspb.Int32Value:
		return protoreflect.ValueOfInt32(v.Value)
	case *wrapperspb.UInt32Value:
		return protoreflect.ValueOfUint32(v.Value)
	case *wrapperspb.Int64Value:
		return protoreflect.ValueOfInt64(v.Value)
	case *wrapperspb.UInt64Value:
		return protoreflect.ValueOfUint64(v.Value)
	case *wrapperspb.FloatValue:
		return protoreflect.ValueOfFloat32(v.Value)
	case *wrapperspb.DoubleValue:
		return protoreflect.ValueOfFloat64(v.Value)
	case *wrapperspb.StringValue:
		return protoreflect.ValueOfString(v.Value)
	case *wrapperspb.BytesValue:
		return protoreflect.ValueOfBytes(v.Value)
	}
	panic(fmt.Sprintf("valueFromWrapper: unreachable type %T", w))
}

// parseEnum accepts either the enum value name or its decimal number.
func parseEnum(enum protoreflect.EnumDescriptor, envVal string) (protoreflect.Value, error) {
	if n, err := strconv.ParseInt(envVal, 10, 32); err == nil {
		return protoreflect.ValueOfEnum(protoreflect.EnumNumber(n)), nil
	}
	ev := enum.Values().ByName(protoreflect.Name(envVal))
	if ev == nil {
		return protoreflect.Value{}, fmt.Errorf("unknown enum value %q for %s", envVal, enum.FullName())
	}
	return protoreflect.ValueOfEnum(ev.Number()), nil
}

// parseMessage parses envVal as a whole message value. Tries protojson first
// (which handles WKT special forms like FieldMask-as-string and
// Duration/Timestamp strings exactly), then falls back to protoyaml (which
// accepts richer YAML syntax for regular messages). Requires md's full name
// to be registered in protoregistry.GlobalTypes.
func parseMessage(md protoreflect.MessageDescriptor, envVal string) (protoreflect.Value, error) {
	mt, err := protoregistry.GlobalTypes.FindMessageByName(md.FullName())
	if err != nil {
		return protoreflect.Value{}, fmt.Errorf("message type %s not registered: %w", md.FullName(), err)
	}
	msg := mt.New()
	if err := safeProtojsonUnmarshal([]byte(envVal), msg.Interface()); err == nil {
		return protoreflect.ValueOfMessage(msg), nil
	}
	msg = mt.New()
	if err := safeProtoyamlUnmarshal(protoyaml.UnmarshalOptions{}, []byte(envVal), msg.Interface()); err != nil {
		return protoreflect.Value{}, err
	}
	return protoreflect.ValueOfMessage(msg), nil
}

// safeProtoyamlUnmarshal wraps protoyaml.Unmarshal with recover. Upstream
// protoyaml has known panics on certain malformed inputs (e.g. bare CR +
// scalar). Convert them to errors so callers never crash.
func safeProtoyamlUnmarshal(opts protoyaml.UnmarshalOptions, b []byte, m proto.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("protoyaml panic: %v", r)
		}
	}()
	return opts.Unmarshal(b, m)
}

func safeProtojsonUnmarshal(b []byte, m proto.Message) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("protojson panic: %v", r)
		}
	}()
	return protojson.Unmarshal(b, m)
}

// parseWholeList parses envVal as the entire value of a repeated field.
// listFd must be a repeated field. Its parent message must be registered.
func parseWholeList(listFd protoreflect.FieldDescriptor, envVal string) (protoreflect.Value, error) {
	if !listFd.IsList() {
		return protoreflect.Value{}, fmt.Errorf("parseWholeList: field %s is not repeated", listFd.FullName())
	}
	return parseViaParent(listFd, envVal)
}

// parseWholeMap parses envVal as the entire value of a map field.
func parseWholeMap(mapFd protoreflect.FieldDescriptor, envVal string) (protoreflect.Value, error) {
	if !mapFd.IsMap() {
		return protoreflect.Value{}, fmt.Errorf("parseWholeMap: field %s is not a map", mapFd.FullName())
	}
	return parseViaParent(mapFd, envVal)
}

// parseViaParent synthesizes {"<fd.Name()>": <envVal>} into a new instance of
// fd.Parent() and extracts the parsed field.
func parseViaParent(fd protoreflect.FieldDescriptor, envVal string) (protoreflect.Value, error) {
	parent, ok := fd.Parent().(protoreflect.MessageDescriptor)
	if !ok {
		return protoreflect.Value{}, fmt.Errorf("field %s has non-message parent", fd.FullName())
	}
	mt, err := protoregistry.GlobalTypes.FindMessageByName(parent.FullName())
	if err != nil {
		return protoreflect.Value{}, fmt.Errorf("parent type %s not registered: %w", parent.FullName(), err)
	}
	host := mt.New()
	body := []byte(fmt.Sprintf(`{%q: %s}`, fd.Name(), envVal))
	if err := safeProtoyamlUnmarshal(protoyaml.UnmarshalOptions{}, body, host.Interface()); err != nil {
		return protoreflect.Value{}, err
	}
	return host.Get(fd), nil
}
