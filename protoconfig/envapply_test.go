package protoconfig

import (
	"strings"
	"testing"
	"time"

	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func entries(pairs ...string) []envEntry {
	out := make([]envEntry, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		key := pairs[i]
		val := pairs[i+1]
		// Split by "__" to match the default delimiter.
		segs := strings.Split(key, "__")
		out = append(out, envEntry{segments: segs, value: val, rawKey: key})
	}
	return out
}

// --- basic scalars ---

func TestApplyEnv_Scalar(t *testing.T) {
	cfg := &protoconfigv1.Test{StringField: "old"}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"STRING_FIELD", `"new"`,
		"INT32_FIELD", "42",
		"BOOL_FIELD", "true",
	), true)
	if err != nil {
		t.Fatalf("applyEnv: %v", err)
	}
	want := &protoconfigv1.Test{StringField: "new", Int32Field: 42, BoolField: true}
	if diff := cmp.Diff(want, cfg, protocmp.Transform()); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

func TestApplyEnv_EmptyValueSkipped(t *testing.T) {
	cfg := &protoconfigv1.Test{StringField: "keep"}
	err := applyEnv(cfg.ProtoReflect(), entries("STRING_FIELD", ""), true)
	if err != nil {
		t.Fatalf("applyEnv: %v", err)
	}
	if cfg.StringField != "keep" {
		t.Errorf("got %q want keep", cfg.StringField)
	}
}

// --- unknown field handling ---

func TestApplyEnv_UnknownStrict(t *testing.T) {
	cfg := &protoconfigv1.Test{}
	err := applyEnv(cfg.ProtoReflect(), entries("UNKNOWN_FIELD", "x"), true)
	if err == nil {
		t.Errorf("want error, got nil")
	}
}

func TestApplyEnv_UnknownLenient(t *testing.T) {
	cfg := &protoconfigv1.Test{}
	err := applyEnv(cfg.ProtoReflect(), entries("UNKNOWN_FIELD", "x"), false)
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
}

// --- nested message ---

func TestApplyEnv_NestedMessageAutoInstantiate(t *testing.T) {
	cfg := &protoconfigv1.Test{} // NestedMessageField is nil
	err := applyEnv(cfg.ProtoReflect(), entries(
		"NESTED_MESSAGE_FIELD__STRING_FIELD", `"x"`,
	), true)
	if err != nil {
		t.Fatalf("applyEnv: %v", err)
	}
	if cfg.NestedMessageField == nil || cfg.NestedMessageField.StringField != "x" {
		t.Errorf("got %+v", cfg.NestedMessageField)
	}
}

func TestApplyEnv_WholeMessage(t *testing.T) {
	cfg := &protoconfigv1.Test{
		NestedMessageField: &protoconfigv1.Nested{StringField: "old", NotUpdatedViaEnv: "keep"},
	}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"NESTED_MESSAGE_FIELD", `{"string_field":"new"}`,
	), true)
	if err != nil {
		t.Fatalf("applyEnv: %v", err)
	}
	// Whole-message replaces — NotUpdatedViaEnv is cleared.
	if cfg.NestedMessageField.StringField != "new" || cfg.NestedMessageField.NotUpdatedViaEnv != "" {
		t.Errorf("got %+v", cfg.NestedMessageField)
	}
}

// --- lists ---

func TestApplyEnv_ListWholeReplaces(t *testing.T) {
	cfg := &protoconfigv1.Test{ListOfInts: []int32{10, 20, 30}}
	err := applyEnv(cfg.ProtoReflect(), entries("LIST_OF_INTS", "[1,2,3]"), true)
	if err != nil {
		t.Fatalf("applyEnv: %v", err)
	}
	got := cfg.ListOfInts
	want := []int32{1, 2, 3}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

func TestApplyEnv_ListIndexedOverlay(t *testing.T) {
	cfg := &protoconfigv1.Test{
		RepeatedNestedMessage: []*protoconfigv1.Nested2{
			{StringField: "a", NotUpdatedViaEnv: "keep-a"},
			{StringField: "b", NotUpdatedViaEnv: "keep-b"},
			{StringField: "c", NotUpdatedViaEnv: "keep-c"},
		},
	}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"REPEATED_NESTED_MESSAGE__1__STRING_FIELD", `"override"`,
	), true)
	if err != nil {
		t.Fatalf("applyEnv: %v", err)
	}
	want := []*protoconfigv1.Nested2{
		{StringField: "a", NotUpdatedViaEnv: "keep-a"},
		{StringField: "override", NotUpdatedViaEnv: "keep-b"},
		{StringField: "c", NotUpdatedViaEnv: "keep-c"},
	}
	if diff := cmp.Diff(want, cfg.RepeatedNestedMessage, protocmp.Transform()); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

func TestApplyEnv_ListOutOfBounds(t *testing.T) {
	cfg := &protoconfigv1.Test{
		RepeatedNestedMessage: []*protoconfigv1.Nested2{{StringField: "a"}},
	}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"REPEATED_NESTED_MESSAGE__5__STRING_FIELD", `"x"`,
	), true)
	if err == nil {
		t.Errorf("want OOB error, got nil")
	}
}

func TestApplyEnv_ListCollision(t *testing.T) {
	cfg := &protoconfigv1.Test{}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"LIST_OF_INTS", "[1,2,3]",
		"LIST_OF_INTS__0", "5",
	), true)
	if err == nil {
		t.Errorf("want collision error, got nil")
	}
}

// --- maps ---

func TestApplyEnv_MapWholeReplaces(t *testing.T) {
	cfg := &protoconfigv1.Test{
		PrimitiveMap: map[string]string{"old-key": "old-val"},
	}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"PRIMITIVE_MAP", `{"k":"v"}`,
	), true)
	if err != nil {
		t.Fatalf("applyEnv: %v", err)
	}
	want := map[string]string{"k": "v"}
	if diff := cmp.Diff(want, cfg.PrimitiveMap); diff != "" {
		t.Errorf("-want +got:\n%s", diff)
	}
}

func TestApplyEnv_MapPerKeyUpdate(t *testing.T) {
	cfg := &protoconfigv1.Test{
		StringToMap: map[string]*protoconfigv1.Nested2{
			"existing": {StringField: "old", NotUpdatedViaEnv: "keep"},
		},
	}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"STRING_TO_MAP__existing__STRING_FIELD", `"new"`,
	), true)
	if err != nil {
		t.Fatalf("applyEnv: %v", err)
	}
	entry := cfg.StringToMap["existing"]
	if entry.StringField != "new" || entry.NotUpdatedViaEnv != "keep" {
		t.Errorf("got %+v", entry)
	}
}

func TestApplyEnv_MapAddKey(t *testing.T) {
	cfg := &protoconfigv1.Test{
		StringToMap: map[string]*protoconfigv1.Nested2{
			"existing": {StringField: "old"},
		},
	}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"STRING_TO_MAP__new__STRING_FIELD", `"fresh"`,
	), true)
	if err != nil {
		t.Fatalf("applyEnv: %v", err)
	}
	if cfg.StringToMap["new"].StringField != "fresh" {
		t.Errorf("got %+v", cfg.StringToMap["new"])
	}
	if cfg.StringToMap["existing"].StringField != "old" {
		t.Errorf("existing mutated: %+v", cfg.StringToMap["existing"])
	}
}

func TestApplyEnv_MapCollision(t *testing.T) {
	cfg := &protoconfigv1.Test{}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"PRIMITIVE_MAP", `{"k":"v"}`,
		"PRIMITIVE_MAP__other", `"x"`,
	), true)
	if err == nil {
		t.Errorf("want collision error, got nil")
	}
}

// --- WKTs ---

func TestApplyEnv_WKTWholeValue(t *testing.T) {
	cfg := &protoconfigv1.Test{}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"TIMESTAMP", `"2024-02-03T17:00:00Z"`,
	), true)
	if err != nil {
		t.Fatalf("applyEnv: %v", err)
	}
	want := timestamppb.New(time.Date(2024, 2, 3, 17, 0, 0, 0, time.UTC))
	if !cfg.Timestamp.AsTime().Equal(want.AsTime()) {
		t.Errorf("got %v want %v", cfg.Timestamp.AsTime(), want.AsTime())
	}
}

func TestApplyEnv_WKTDeepAddressingRejected(t *testing.T) {
	cfg := &protoconfigv1.Test{}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"TIMESTAMP__SECONDS", "1",
	), true)
	if err == nil {
		t.Errorf("want WKT-deep-addressing error, got nil")
	}
}

// --- self-ref ---

func TestApplyEnv_SelfRefDoesNotInfiniteRecurse(t *testing.T) {
	// Test.message_field = Test. If we have an env var targeting
	// message_field.message_field.string_field, we must instantiate both layers
	// but not recurse further.
	cfg := &protoconfigv1.Test{}
	err := applyEnv(cfg.ProtoReflect(), entries(
		"MESSAGE_FIELD__MESSAGE_FIELD__STRING_FIELD", `"deep"`,
	), true)
	if err != nil {
		t.Fatalf("applyEnv: %v", err)
	}
	if cfg.MessageField == nil || cfg.MessageField.MessageField == nil {
		t.Fatalf("got %+v", cfg)
	}
	if cfg.MessageField.MessageField.StringField != "deep" {
		t.Errorf("got %q", cfg.MessageField.MessageField.StringField)
	}
}

// --- tree construction errors ---

func TestBuildEnvTree_DuplicateKey(t *testing.T) {
	_, err := buildEnvTree(entries(
		"A__B", "1",
		"A__B", "2",
	))
	if err == nil {
		t.Errorf("want duplicate key error")
	}
}
