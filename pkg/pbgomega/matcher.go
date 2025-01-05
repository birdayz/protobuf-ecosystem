package pbgomega

import (
	"fmt"

	"github.com/google/go-cmp/cmp"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

// EqualProto succeeds if actual proto matches the passed-in proto.
func EqualProto(message proto.Message) types.GomegaMatcher {
	return &equalProtoMatcher{expected: message}
}

type equalProtoMatcher struct {
	expected proto.Message
}

func (matcher *equalProtoMatcher) Match(actual any) (bool, error) {
	message, ok := actual.(proto.Message)
	if !ok {
		return false, fmt.Errorf("EqualProto matcher expects a proto message.  Got:\n%s", format.Object(actual, 1))
	}

	if actual == nil && matcher.expected == nil {
		return false, fmt.Errorf("Refusing to compare <nil> to <nil>.\nBe explicit and use BeNil() instead.  This is to avoid mistakes where both sides of an assertion are erroneously uninitialized")
	}

	return cmp.Equal(message, matcher.expected, protocmp.Transform()), nil
}

func (matcher *equalProtoMatcher) FailureMessage(actual any) string {
	diff := cmp.Diff(actual, matcher.expected, protocmp.Transform())
	return "Mismatch.\n-: present, but not expected\n+: expected, but not present(-actual +expected):\n" + diff
}

func (matcher *equalProtoMatcher) NegatedFailureMessage(actual any) string {
	diff := cmp.Diff(actual, matcher.expected, protocmp.Transform())
	return "Mismatch.\n-: present, but not expected\n+: expected, but not present(-actual +expected):\n" + diff
}
