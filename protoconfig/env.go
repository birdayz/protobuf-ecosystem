package protoconfig

import (
	"fmt"
	"strconv"
	"strings"

	"buf.build/go/protoyaml"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// stringToMap transforms a protojson string into a protoreflect.Value.
// Since protojson and protoyaml don't allow direct unmarshaling of non-message types, we create a copy of the
// parent message type, unmarshal into it, and return the relevant field.
func stringToProtoValue(fd protoreflect.FieldDescriptor, strVal string) (protoreflect.Value, error) {
	md, err := protoregistry.GlobalTypes.FindMessageByName(fd.Parent().FullName())
	if err != nil {
		return protoreflect.Value{}, err
	}

	fake := md.New()

	if err := protoyaml.Unmarshal([]byte(fmt.Sprintf(`{"%s":%s}`, fd.Name(), strVal)), fake.Interface()); err != nil {
		return protoreflect.Value{}, err
	}

	return fake.Get(fd), nil

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
