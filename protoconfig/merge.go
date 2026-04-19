package protoconfig

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// mergeFrom overlays src onto dst in place. Semantics:
//   - scalars: src replaces dst if present in src
//   - messages (singular): recursive merge
//   - lists (repeated): src replaces dst entirely if src is non-empty
//   - maps: merge by key — src keys overwrite dst's; other dst keys preserved
//
// WKTs are treated as opaque messages: the value is swapped wholesale, not
// recursed into. This matches DESIGN.md §merge/overlay.
//
// Both messages must share the same descriptor type (not enforced here — the
// caller is responsible; usually enforced by generic type).
func mergeFrom(dst, src protoreflect.Message) {
	src.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch {
		case fd.IsList():
			dstList := dst.Mutable(fd).List()
			// Replace: clear and copy src entries.
			dstList.Truncate(0)
			srcList := v.List()
			for i := 0; i < srcList.Len(); i++ {
				dstList.Append(cloneListElem(fd, srcList.Get(i)))
			}
		case fd.IsMap():
			dstMap := dst.Mutable(fd).Map()
			v.Map().Range(func(k protoreflect.MapKey, mv protoreflect.Value) bool {
				dstMap.Set(k, cloneMapValue(fd, mv))
				return true
			})
		case fd.Kind() == protoreflect.MessageKind, fd.Kind() == protoreflect.GroupKind:
			if isWKT(fd.Message().FullName()) {
				// WKT: replace opaque value.
				dst.Set(fd, cloneMessageValue(v))
				break
			}
			mergeFrom(dst.Mutable(fd).Message(), v.Message())
		default:
			dst.Set(fd, v)
		}
		return true
	})
}

// cloneListElem returns a deep copy of a list element so dst doesn't share
// memory with src. For primitives the protoreflect.Value is already safe; for
// messages we must proto.Clone.
func cloneListElem(fd protoreflect.FieldDescriptor, v protoreflect.Value) protoreflect.Value {
	if fd.Kind() != protoreflect.MessageKind && fd.Kind() != protoreflect.GroupKind {
		return v
	}
	return cloneMessageValue(v)
}

// cloneMapValue returns a deep copy of a map value.
func cloneMapValue(mapFd protoreflect.FieldDescriptor, v protoreflect.Value) protoreflect.Value {
	valFd := mapFd.MapValue()
	if valFd.Kind() != protoreflect.MessageKind && valFd.Kind() != protoreflect.GroupKind {
		return v
	}
	return cloneMessageValue(v)
}

func cloneMessageValue(v protoreflect.Value) protoreflect.Value {
	return protoreflect.ValueOfMessage(proto.Clone(v.Message().Interface()).ProtoReflect())
}

// isWKT returns true for messages whose behavior under env/YAML is "opaque
// leaf": replace wholesale rather than recursively merge.
func isWKT(name protoreflect.FullName) bool {
	switch name {
	case "google.protobuf.Timestamp",
		"google.protobuf.Duration",
		"google.protobuf.FieldMask",
		"google.protobuf.Any",
		"google.protobuf.Empty",
		"google.protobuf.Struct",
		"google.protobuf.Value",
		"google.protobuf.ListValue",
		"google.protobuf.NullValue",
		"google.protobuf.BoolValue",
		"google.protobuf.StringValue",
		"google.protobuf.BytesValue",
		"google.protobuf.Int32Value",
		"google.protobuf.Int64Value",
		"google.protobuf.UInt32Value",
		"google.protobuf.UInt64Value",
		"google.protobuf.FloatValue",
		"google.protobuf.DoubleValue":
		return true
	}
	return false
}
