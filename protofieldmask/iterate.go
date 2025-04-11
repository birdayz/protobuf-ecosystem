package protofieldmask

import (
	"iter"
	"slices"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Matcher defines how we find the matching field on the other side.
type Matcher func(fd protoreflect.FieldDescriptor, otherMessage protoreflect.MessageDescriptor) protoreflect.FieldDescriptor

func MatcherByNumber(fd protoreflect.FieldDescriptor, otherMessage protoreflect.MessageDescriptor) protoreflect.FieldDescriptor {
	return otherMessage.Fields().ByNumber(fd.Number())
}

func MatcherByName(fd protoreflect.FieldDescriptor, otherMessage protoreflect.MessageDescriptor) protoreflect.FieldDescriptor {
	return otherMessage.Fields().ByName(fd.Name())
}

func MatcherByJSONName(fd protoreflect.FieldDescriptor, otherMessage protoreflect.MessageDescriptor) protoreflect.FieldDescriptor {
	return otherMessage.Fields().ByJSONName(fd.JSONName())
}

func MatcherByTextName(fd protoreflect.FieldDescriptor, otherMessage protoreflect.MessageDescriptor) protoreflect.FieldDescriptor {
	return otherMessage.Fields().ByTextName(fd.TextName())
}

// IterFieldDescriptors returns matching pairs of FielDescriptors within two messages.
// "Matching" can be based on Number, Name, JSONName, TextName.
func IterFieldDescriptor(leftMessage protoreflect.MessageDescriptor) iter.Seq[protoreflect.FieldDescriptor] {
	return func(yield func(protoreflect.FieldDescriptor) bool) {
		fdsLeft := leftMessage.Fields()

		// 1) First iterate over left FieldDescriptors.
		// For every item, find the equivalent in right FieldDescriptors.
		for i := 0; i < fdsLeft.Len(); i++ {
			var l protoreflect.FieldDescriptor
			l = fdsLeft.Get(i)

			if !yield(l) {
				return
			}
		}
	}
}

func WrapStaticNil(it iter.Seq[protoreflect.FieldDescriptor]) iter.Seq2[protoreflect.FieldDescriptor, protoreflect.FieldDescriptor] {
	return func(yield func(protoreflect.FieldDescriptor, protoreflect.FieldDescriptor) bool) {
		for item := range it {
			if !yield(item, nil) {
				return
			}
		}
	}
}

// IterFieldDescriptors returns matching pairs of FielDescriptors within two messages.
// "Matching" can be based on Number, Name, JSONName, TextName.
func IterFieldDescriptors(leftMessage, rightMessage protoreflect.MessageDescriptor, matchBy Matcher) iter.Seq2[protoreflect.FieldDescriptor, protoreflect.FieldDescriptor] {
	return func(yield func(protoreflect.FieldDescriptor, protoreflect.FieldDescriptor) bool) {
		fdsLeft := leftMessage.Fields()
		fdsRight := rightMessage.Fields()

		visitedNumbersOnRight := []protowire.Number{}

		// 1) First iterate over left FieldDescriptors.
		// For every item, find the equivalent in right FieldDescriptors.
		for i := 0; i < fdsLeft.Len(); i++ {
			var l, r protoreflect.FieldDescriptor
			l = fdsLeft.Get(i)

			r = matchBy(l, rightMessage)
			if r != nil {
				idx := r.Number()
				if idx != 0 {
					visitedNumbersOnRight = append(visitedNumbersOnRight, idx)
				}
			}

			if !yield(l, r) {
				return
			}
		}

		// 2) After finding all left + matching partner on the right,
		// find the left-overs:
		// Items that are in right, but have no partner in left.
		// Skip items already visted in step 1).
		for i := 0; i < fdsRight.Len(); i++ {
			var l, r protoreflect.FieldDescriptor
			r = fdsRight.Get(i)

			l = matchBy(r, leftMessage)

			if slices.Contains(visitedNumbersOnRight, r.Number()) {
				continue
			}

			if !yield(l, r) {
				return
			}
		}

		// We visted indexes 0-fdsLeft.Len already.
		// Now iterate over all of the right side, but skip the ones visited already (bbecause triggered by left)

		// Iterate over left.
		// Find matching on right
	}
}
