package protoconfig

import (
	"fmt"
	"sort"
	"strconv"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// envEntry represents one stripped env var (prefix already removed, split
// into segments by the delimiter).
type envEntry struct {
	segments []string
	value    string
	rawKey   string // original env var name, for error messages
}

// applyEnv overlays the given env entries onto msg per D1/D2 semantics.
// If strict is true, entries whose segments don't map to any field error; if
// false, they are silently ignored.
func applyEnv(msg protoreflect.Message, entries []envEntry, strict bool) error {
	root, err := buildEnvTree(entries)
	if err != nil {
		return err
	}
	return applyMessage(msg, root, strict, "")
}

// envTree is a trie built from env entries. A node can either hold a leaf
// value (whole-value env) or descend into children (indexed / per-field env).
// Holding both at the same node is a collision and is rejected at build time.
type envTree struct {
	children map[string]*envTree
	hasLeaf  bool
	leaf     string
	rawKey   string // for collision error messages
}

func newEnvTree() *envTree {
	return &envTree{children: map[string]*envTree{}}
}

func buildEnvTree(entries []envEntry) (*envTree, error) {
	root := newEnvTree()
	for _, e := range entries {
		if len(e.segments) == 0 {
			// Prefix-only env var (no path). Nothing to apply.
			continue
		}
		cur := root
		for i, seg := range e.segments {
			child := cur.children[seg]
			if child == nil {
				child = newEnvTree()
				cur.children[seg] = child
			}
			if i == len(e.segments)-1 {
				if child.hasLeaf {
					return nil, fmt.Errorf("duplicate env keys: %s and %s", child.rawKey, e.rawKey)
				}
				if len(child.children) > 0 {
					return nil, fmt.Errorf("env key %s conflicts with sub-keys", e.rawKey)
				}
				child.hasLeaf = true
				child.leaf = e.value
				child.rawKey = e.rawKey
			} else {
				if child.hasLeaf {
					return nil, fmt.Errorf("env key %s conflicts with whole-value form %s", e.rawKey, child.rawKey)
				}
				cur = child
			}
		}
	}
	return root, nil
}

// applyMessage walks the node's children as field names on msg.
func applyMessage(msg protoreflect.Message, node *envTree, strict bool, pathCtx string) error {
	// Oneof conflict detection: at most one arm per oneof may be set.
	seenOneof := map[protoreflect.Name]protoreflect.FieldDescriptor{}
	for _, seg := range sortedKeys(node.children) {
		child := node.children[seg]
		fd := findField(msg.Descriptor(), seg)
		if fd == nil {
			if strict {
				return fmt.Errorf("unknown env key %s (no field %s on %s)", firstRawKey(child), seg, msg.Descriptor().FullName())
			}
			continue
		}
		if oneof := fd.ContainingOneof(); oneof != nil && !oneof.IsSynthetic() {
			if prior, ok := seenOneof[oneof.Name()]; ok {
				return fmt.Errorf("oneof %s: cannot set both %s and %s (via %s)",
					oneof.Name(), prior.Name(), fd.Name(), firstRawKey(child))
			}
			seenOneof[oneof.Name()] = fd
		}
		if err := applyField(msg, fd, child, strict, joinPath(pathCtx, string(fd.Name()))); err != nil {
			return err
		}
	}
	return nil
}

// applyField dispatches based on the field's shape.
func applyField(msg protoreflect.Message, fd protoreflect.FieldDescriptor, node *envTree, strict bool, pathCtx string) error {
	switch {
	case fd.IsList():
		return applyListField(msg, fd, node, strict, pathCtx)
	case fd.IsMap():
		return applyMapField(msg, fd, node, strict, pathCtx)
	case fd.Kind() == protoreflect.MessageKind, fd.Kind() == protoreflect.GroupKind:
		return applySingularMessage(msg, fd, node, strict, pathCtx)
	default:
		return applySingularScalar(msg, fd, node, pathCtx)
	}
}

func applySingularScalar(msg protoreflect.Message, fd protoreflect.FieldDescriptor, node *envTree, pathCtx string) error {
	if !node.hasLeaf {
		return fmt.Errorf("env key for %s: scalar cannot have sub-keys", pathCtx)
	}
	if node.leaf == "" {
		// Empty value: skip per DESIGN.md (env "" means unset).
		return nil
	}
	v, err := parseScalar(fd.Kind(), fd.Enum(), node.leaf)
	if err != nil {
		return fmt.Errorf("env key %s (%s): %w", node.rawKey, pathCtx, err)
	}
	msg.Set(fd, v)
	return nil
}

func applySingularMessage(msg protoreflect.Message, fd protoreflect.FieldDescriptor, node *envTree, strict bool, pathCtx string) error {
	md := fd.Message()
	if node.hasLeaf {
		// Whole-value form. Skip empty.
		if node.leaf == "" {
			return nil
		}
		v, err := parseMessage(md, node.leaf)
		if err != nil {
			return fmt.Errorf("env key %s (%s): %w", node.rawKey, pathCtx, err)
		}
		msg.Set(fd, v)
		return nil
	}
	// No leaf: recurse into sub-fields.
	if isWKT(md.FullName()) {
		return fmt.Errorf("env key %s: WKT %s does not support deep addressing in v1", firstRawKey(node), md.FullName())
	}
	// msg.Mutable lazily creates the nested message if nil.
	sub := msg.Mutable(fd).Message()
	return applyMessage(sub, node, strict, pathCtx)
}

func applyListField(msg protoreflect.Message, fd protoreflect.FieldDescriptor, node *envTree, strict bool, pathCtx string) error {
	if node.hasLeaf {
		if node.leaf == "" {
			return nil
		}
		v, err := parseWholeList(fd, node.leaf)
		if err != nil {
			return fmt.Errorf("env key %s (%s): %w", node.rawKey, pathCtx, err)
		}
		// Copy into the mutable list — we can't Set a read-only value from
		// Get() onto another message.
		dst := msg.Mutable(fd).List()
		dst.Truncate(0)
		src := v.List()
		for i := 0; i < src.Len(); i++ {
			dst.Append(cloneListElem(fd, src.Get(i)))
		}
		return nil
	}
	// Indexed form. D2: in-bounds overlay; out-of-bounds = error.
	list := msg.Mutable(fd).List()
	for _, seg := range sortedKeys(node.children) {
		child := node.children[seg]
		idx, err := strconv.Atoi(seg)
		if err != nil {
			return fmt.Errorf("env key %s: expected list index, got %q", firstRawKey(child), seg)
		}
		if idx < 0 || idx >= list.Len() {
			return fmt.Errorf("env key %s: list index %d out of bounds (list has %d elements); use whole-value form to extend", firstRawKey(child), idx, list.Len())
		}
		subPath := fmt.Sprintf("%s[%d]", pathCtx, idx)
		if err := applyListElem(list, fd, idx, child, strict, subPath); err != nil {
			return err
		}
	}
	return nil
}

func applyListElem(list protoreflect.List, fd protoreflect.FieldDescriptor, idx int, node *envTree, strict bool, pathCtx string) error {
	if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
		md := fd.Message()
		if node.hasLeaf {
			if node.leaf == "" {
				return nil
			}
			v, err := parseMessage(md, node.leaf)
			if err != nil {
				return fmt.Errorf("env key %s (%s): %w", node.rawKey, pathCtx, err)
			}
			list.Set(idx, v)
			return nil
		}
		if isWKT(md.FullName()) {
			return fmt.Errorf("env key %s: WKT %s does not support deep addressing in v1", firstRawKey(node), md.FullName())
		}
		// Recurse into the element. list.Get returns a value; for messages we
		// need a mutable handle. The descriptor docs require a copy-modify-set.
		existing := list.Get(idx)
		clone := proto.Clone(existing.Message().Interface()).ProtoReflect()
		if err := applyMessage(clone, node, strict, pathCtx); err != nil {
			return err
		}
		list.Set(idx, protoreflect.ValueOfMessage(clone))
		return nil
	}
	// Scalar list element.
	if !node.hasLeaf {
		return fmt.Errorf("env key for %s: scalar list element cannot have sub-keys", pathCtx)
	}
	if node.leaf == "" {
		return nil
	}
	v, err := parseScalar(fd.Kind(), fd.Enum(), node.leaf)
	if err != nil {
		return fmt.Errorf("env key %s (%s): %w", node.rawKey, pathCtx, err)
	}
	list.Set(idx, v)
	return nil
}

func applyMapField(msg protoreflect.Message, fd protoreflect.FieldDescriptor, node *envTree, strict bool, pathCtx string) error {
	if node.hasLeaf {
		if node.leaf == "" {
			return nil
		}
		v, err := parseWholeMap(fd, node.leaf)
		if err != nil {
			return fmt.Errorf("env key %s (%s): %w", node.rawKey, pathCtx, err)
		}
		// Copy into the mutable map — a read-only value from Get() can't be
		// Set onto another message.
		dst := msg.Mutable(fd).Map()
		// Clear existing entries.
		var keys []protoreflect.MapKey
		dst.Range(func(k protoreflect.MapKey, _ protoreflect.Value) bool {
			keys = append(keys, k)
			return true
		})
		for _, k := range keys {
			dst.Clear(k)
		}
		v.Map().Range(func(k protoreflect.MapKey, mv protoreflect.Value) bool {
			dst.Set(k, cloneMapValue(fd, mv))
			return true
		})
		return nil
	}
	m := msg.Mutable(fd).Map()
	keyFd := fd.MapKey()
	for _, seg := range sortedKeys(node.children) {
		child := node.children[seg]
		keyVal, err := parseScalar(keyFd.Kind(), nil, seg)
		if err != nil {
			return fmt.Errorf("env key %s: invalid map key %q: %w", firstRawKey(child), seg, err)
		}
		key := keyVal.MapKey()
		subPath := fmt.Sprintf("%s[%s]", pathCtx, seg)
		if err := applyMapEntry(m, fd, key, child, strict, subPath); err != nil {
			return err
		}
	}
	return nil
}

func applyMapEntry(m protoreflect.Map, fd protoreflect.FieldDescriptor, key protoreflect.MapKey, node *envTree, strict bool, pathCtx string) error {
	valFd := fd.MapValue()
	if valFd.Kind() == protoreflect.MessageKind || valFd.Kind() == protoreflect.GroupKind {
		md := valFd.Message()
		if node.hasLeaf {
			if node.leaf == "" {
				return nil
			}
			v, err := parseMessage(md, node.leaf)
			if err != nil {
				return fmt.Errorf("env key %s (%s): %w", node.rawKey, pathCtx, err)
			}
			m.Set(key, v)
			return nil
		}
		if isWKT(md.FullName()) {
			return fmt.Errorf("env key %s: WKT %s does not support deep addressing in v1", firstRawKey(node), md.FullName())
		}
		// Copy-modify-set pattern. For maps, Get returns by value.
		var clone protoreflect.Message
		if m.Has(key) {
			clone = proto.Clone(m.Get(key).Message().Interface()).ProtoReflect()
		} else {
			clone = m.NewValue().Message()
		}
		if err := applyMessage(clone, node, strict, pathCtx); err != nil {
			return err
		}
		m.Set(key, protoreflect.ValueOfMessage(clone))
		return nil
	}
	// Scalar map value.
	if !node.hasLeaf {
		return fmt.Errorf("env key for %s: scalar map value cannot have sub-keys", pathCtx)
	}
	if node.leaf == "" {
		return nil
	}
	v, err := parseScalar(valFd.Kind(), valFd.Enum(), node.leaf)
	if err != nil {
		return fmt.Errorf("env key %s (%s): %w", node.rawKey, pathCtx, err)
	}
	m.Set(key, v)
	return nil
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// firstRawKey returns any raw key under node for error messages.
func firstRawKey(node *envTree) string {
	if node.hasLeaf {
		return node.rawKey
	}
	for _, seg := range sortedKeys(node.children) {
		if k := firstRawKey(node.children[seg]); k != "" {
			return k
		}
	}
	return ""
}

func joinPath(ctx, seg string) string {
	if ctx == "" {
		return seg
	}
	return ctx + "." + seg
}
