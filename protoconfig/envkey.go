package protoconfig

import (
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const defaultDelim = "__"

// pathSegments returns the env-key segments for a protopath.Path.
// Root is omitted. Field names are uppercased. Map keys keep their
// original stringified form. List indices are decimal.
func pathSegments(path protopath.Path) ([]string, error) {
	var out []string
	for _, step := range path {
		switch step.Kind() {
		case protopath.RootStep:
			continue
		case protopath.FieldAccessStep:
			out = append(out, strings.ToUpper(string(step.FieldDescriptor().Name())))
		case protopath.ListIndexStep:
			out = append(out, strconv.Itoa(step.ListIndex()))
		case protopath.MapIndexStep:
			seg := step.MapIndex().String()
			out = append(out, seg)
		default:
			return nil, fmt.Errorf("unsupported path step kind: %v", step.Kind())
		}
	}
	return out, nil
}

// envKey joins path segments with the delimiter, prepending prefix if set.
// Segments must not be empty, must not contain the delimiter, and must not
// begin or end with any character from the delimiter (otherwise adjacent
// segments could merge into the delimiter boundary and create an ambiguous
// split).
func envKey(prefix, delim string, segments []string) (string, error) {
	if delim == "" {
		return "", fmt.Errorf("delimiter must not be empty")
	}
	for _, s := range segments {
		if err := validateSegment(s, delim); err != nil {
			return "", err
		}
	}
	if prefix != "" {
		if err := validateSegment(prefix, delim); err != nil {
			return "", fmt.Errorf("prefix: %w", err)
		}
	}
	joined := strings.Join(segments, delim)
	if prefix == "" {
		return joined, nil
	}
	if joined == "" {
		return prefix, nil
	}
	return prefix + delim + joined, nil
}

// validateSegment enforces the no-ambiguity rules for a single segment.
func validateSegment(s, delim string) error {
	if s == "" {
		return fmt.Errorf("empty segment")
	}
	if strings.Contains(s, delim) {
		return fmt.Errorf("segment %q contains delimiter %q", s, delim)
	}
	first := s[0]
	last := s[len(s)-1]
	for i := 0; i < len(delim); i++ {
		if first == delim[i] {
			return fmt.Errorf("segment %q starts with delimiter character %q", s, delim[i])
		}
		if last == delim[i] {
			return fmt.Errorf("segment %q ends with delimiter character %q", s, delim[i])
		}
	}
	return nil
}

// splitEnvKey strips the prefix and splits on the delimiter.
// Returns (segments, matched). matched is false if the prefix
// is set and the key does not start with prefix+delim.
func splitEnvKey(prefix, delim, key string) ([]string, bool) {
	if prefix != "" {
		if !strings.HasPrefix(key, prefix) {
			return nil, false
		}
		rest := key[len(prefix):]
		if rest == "" {
			return []string{}, true
		}
		if !strings.HasPrefix(rest, delim) {
			return nil, false
		}
		rest = rest[len(delim):]
		if rest == "" {
			return []string{}, true
		}
		return strings.Split(rest, delim), true
	}
	if key == "" {
		return []string{}, true
	}
	return strings.Split(key, delim), true
}

// segmentMatchesField returns true if the segment names the given field
// (case-insensitive compare against field's proto name).
func segmentMatchesField(segment string, fd protoreflect.FieldDescriptor) bool {
	return strings.EqualFold(segment, string(fd.Name()))
}

// findField looks up a field on the message descriptor by the given segment.
// Case-insensitive on snake_case field name.
func findField(md protoreflect.MessageDescriptor, segment string) protoreflect.FieldDescriptor {
	fields := md.Fields()
	for i := 0; i < fields.Len(); i++ {
		fd := fields.Get(i)
		if segmentMatchesField(segment, fd) {
			return fd
		}
	}
	return nil
}
