package protoconfig

import (
	"testing"
	"time"

	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/go-cmp/cmp"
)

func mustMerge(t *testing.T, dst, src proto.Message) proto.Message {
	t.Helper()
	mergeFrom(dst.ProtoReflect(), src.ProtoReflect())
	return dst
}

func TestMerge_ScalarReplaces(t *testing.T) {
	dst := &protoconfigv1.Test{StringField: "old", Int32Field: 1}
	src := &protoconfigv1.Test{StringField: "new"}

	mustMerge(t, dst, src)
	want := &protoconfigv1.Test{StringField: "new", Int32Field: 1}
	if diff := cmp.Diff(want, dst, protocmp.Transform()); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

func TestMerge_MessageRecurses(t *testing.T) {
	dst := &protoconfigv1.Test{
		NestedMessageField: &protoconfigv1.Nested{
			StringField:      "keep-if-not-set",
			NotUpdatedViaEnv: "keep",
		},
	}
	src := &protoconfigv1.Test{
		NestedMessageField: &protoconfigv1.Nested{
			StringField: "override",
		},
	}

	mustMerge(t, dst, src)
	want := &protoconfigv1.Test{
		NestedMessageField: &protoconfigv1.Nested{
			StringField:      "override",
			NotUpdatedViaEnv: "keep",
		},
	}
	if diff := cmp.Diff(want, dst, protocmp.Transform()); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

func TestMerge_ListReplaces(t *testing.T) {
	// This is the key deviation from proto.Merge (which appends).
	dst := &protoconfigv1.Test{
		ListOfInts: []int32{1, 2, 3},
	}
	src := &protoconfigv1.Test{
		ListOfInts: []int32{9},
	}

	mustMerge(t, dst, src)
	want := &protoconfigv1.Test{ListOfInts: []int32{9}}
	if diff := cmp.Diff(want, dst, protocmp.Transform()); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

func TestMerge_ListMessageReplaces(t *testing.T) {
	dst := &protoconfigv1.Test{
		RepeatedNestedMessage: []*protoconfigv1.Nested2{
			{StringField: "a"},
			{StringField: "b"},
			{StringField: "c"},
		},
	}
	src := &protoconfigv1.Test{
		RepeatedNestedMessage: []*protoconfigv1.Nested2{
			{StringField: "only"},
		},
	}

	mustMerge(t, dst, src)
	if len(dst.RepeatedNestedMessage) != 1 {
		t.Fatalf("len=%d want 1", len(dst.RepeatedNestedMessage))
	}
	if dst.RepeatedNestedMessage[0].StringField != "only" {
		t.Errorf("got %q want only", dst.RepeatedNestedMessage[0].StringField)
	}
	// Verify clone: mutate src, dst should not change.
	src.RepeatedNestedMessage[0].StringField = "mutated"
	if dst.RepeatedNestedMessage[0].StringField != "only" {
		t.Errorf("dst aliased to src (got %q)", dst.RepeatedNestedMessage[0].StringField)
	}
}

func TestMerge_MapMergesPerKey(t *testing.T) {
	dst := &protoconfigv1.Test{
		PrimitiveMap: map[string]string{"a": "1", "b": "2"},
	}
	src := &protoconfigv1.Test{
		PrimitiveMap: map[string]string{"b": "99", "c": "3"},
	}

	mustMerge(t, dst, src)
	want := &protoconfigv1.Test{
		PrimitiveMap: map[string]string{"a": "1", "b": "99", "c": "3"},
	}
	if diff := cmp.Diff(want, dst, protocmp.Transform()); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

func TestMerge_MapMessageValuesNotShared(t *testing.T) {
	dst := &protoconfigv1.Test{
		StringToMap: map[string]*protoconfigv1.Nested2{
			"keep": {StringField: "old"},
		},
	}
	src := &protoconfigv1.Test{
		StringToMap: map[string]*protoconfigv1.Nested2{
			"new": {StringField: "v"},
		},
	}

	mustMerge(t, dst, src)
	if dst.StringToMap["new"].StringField != "v" {
		t.Fatalf("got %q", dst.StringToMap["new"].StringField)
	}
	// Clone check.
	src.StringToMap["new"].StringField = "mutated"
	if dst.StringToMap["new"].StringField != "v" {
		t.Errorf("dst aliased to src")
	}
}

func TestMerge_WKTReplaces(t *testing.T) {
	// Timestamp is a WKT: opaque swap, not recursive merge.
	dst := &protoconfigv1.Test{
		Timestamp: timestamppb.New(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
	}
	src := &protoconfigv1.Test{
		Timestamp: timestamppb.New(time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)),
	}

	mustMerge(t, dst, src)
	got := dst.Timestamp.AsTime()
	want := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func TestMerge_SrcFieldAbsent_DstUnchanged(t *testing.T) {
	dst := &protoconfigv1.Test{StringField: "keep"}
	src := &protoconfigv1.Test{} // nothing set

	mustMerge(t, dst, src)
	if dst.StringField != "keep" {
		t.Errorf("got %q want keep", dst.StringField)
	}
}

func TestMerge_NestedMessageAutoInstantiates(t *testing.T) {
	dst := &protoconfigv1.Test{} // NestedMessageField is nil
	src := &protoconfigv1.Test{
		NestedMessageField: &protoconfigv1.Nested{StringField: "set"},
	}

	mustMerge(t, dst, src)
	if dst.NestedMessageField == nil {
		t.Fatalf("nil after merge")
	}
	if dst.NestedMessageField.StringField != "set" {
		t.Errorf("got %q want set", dst.NestedMessageField.StringField)
	}
}
