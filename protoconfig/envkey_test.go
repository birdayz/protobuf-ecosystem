package protoconfig

import (
	"fmt"
	"testing"

	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
	"google.golang.org/protobuf/reflect/protopath"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func TestEnvKey(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		delim    string
		segments []string
		want     string
		wantErr  bool
	}{
		{name: "no prefix, single segment", delim: "__", segments: []string{"STRING_FIELD"}, want: "STRING_FIELD"},
		{name: "no prefix, two segments", delim: "__", segments: []string{"MSG", "FIELD"}, want: "MSG__FIELD"},
		{name: "prefix, empty segments", prefix: "APP", delim: "__", segments: nil, want: "APP"},
		{name: "prefix, single segment", prefix: "APP", delim: "__", segments: []string{"STRING_FIELD"}, want: "APP__STRING_FIELD"},
		{name: "prefix, nested", prefix: "APP", delim: "__", segments: []string{"MSG", "FIELD"}, want: "APP__MSG__FIELD"},
		{name: "list index", delim: "__", segments: []string{"LIST", "0", "FIELD"}, want: "LIST__0__FIELD"},
		{name: "map key preserves case", delim: "__", segments: []string{"MAP", "my_key", "FIELD"}, want: "MAP__my_key__FIELD"},
		{name: "segment contains delimiter", delim: "__", segments: []string{"BAD__KEY"}, wantErr: true},
		{name: "segment starts with delim char", delim: "__", segments: []string{"_LEADING"}, wantErr: true},
		{name: "segment ends with delim char", delim: "__", segments: []string{"TRAILING_"}, wantErr: true},
		{name: "empty delim", delim: "", segments: []string{"X"}, wantErr: true},
		{name: "empty segment", delim: "__", segments: []string{""}, wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := envKey(tc.prefix, tc.delim, tc.segments)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("want error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q want %q", got, tc.want)
			}
		})
	}
}

func TestSplitEnvKey(t *testing.T) {
	tests := []struct {
		name        string
		prefix      string
		delim       string
		key         string
		wantSegs    []string
		wantMatched bool
	}{
		{name: "no prefix, empty", delim: "__", key: "", wantSegs: []string{}, wantMatched: true},
		{name: "no prefix, single", delim: "__", key: "STRING_FIELD", wantSegs: []string{"STRING_FIELD"}, wantMatched: true},
		{name: "no prefix, nested", delim: "__", key: "MSG__FIELD", wantSegs: []string{"MSG", "FIELD"}, wantMatched: true},
		{name: "prefix-only key", prefix: "APP", delim: "__", key: "APP", wantSegs: []string{}, wantMatched: true},
		{name: "prefix plus separator only", prefix: "APP", delim: "__", key: "APP__", wantSegs: []string{}, wantMatched: true},
		{name: "prefix match", prefix: "APP", delim: "__", key: "APP__FIELD", wantSegs: []string{"FIELD"}, wantMatched: true},
		{name: "prefix no delim after", prefix: "APP", delim: "__", key: "APPLE", wantMatched: false},
		{name: "prefix mismatch", prefix: "APP", delim: "__", key: "OTHER__FIELD", wantMatched: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotSegs, gotMatched := splitEnvKey(tc.prefix, tc.delim, tc.key)
			if gotMatched != tc.wantMatched {
				t.Fatalf("matched: got %v want %v", gotMatched, tc.wantMatched)
			}
			if !tc.wantMatched {
				return
			}
			if len(gotSegs) != len(tc.wantSegs) {
				t.Fatalf("len(segs): got %d want %d (%v vs %v)", len(gotSegs), len(tc.wantSegs), gotSegs, tc.wantSegs)
			}
			for i := range gotSegs {
				if gotSegs[i] != tc.wantSegs[i] {
					t.Errorf("seg %d: got %q want %q", i, gotSegs[i], tc.wantSegs[i])
				}
			}
		})
	}
}

func TestPathSegments(t *testing.T) {
	testDesc := (&protoconfigv1.Test{}).ProtoReflect().Descriptor()
	nested2Desc := (&protoconfigv1.Nested2{}).ProtoReflect().Descriptor()

	fdStringField := testDesc.Fields().ByName("string_field")
	fdRepNested := testDesc.Fields().ByName("repeated_nested_message")
	fdStringToMap := testDesc.Fields().ByName("string_to_map")
	fdNested2String := nested2Desc.Fields().ByName("string_field")

	tests := []struct {
		name string
		path protopath.Path
		want []string
	}{
		{
			name: "root only",
			path: protopath.Path{protopath.Root(testDesc)},
			want: nil,
		},
		{
			name: "single field",
			path: protopath.Path{protopath.Root(testDesc), protopath.FieldAccess(fdStringField)},
			want: []string{"STRING_FIELD"},
		},
		{
			name: "field with list index with nested field",
			path: protopath.Path{
				protopath.Root(testDesc),
				protopath.FieldAccess(fdRepNested),
				protopath.ListIndex(2),
				protopath.FieldAccess(fdNested2String),
			},
			want: []string{"REPEATED_NESTED_MESSAGE", "2", "STRING_FIELD"},
		},
		{
			name: "field with map index with nested field",
			path: protopath.Path{
				protopath.Root(testDesc),
				protopath.FieldAccess(fdStringToMap),
				protopath.MapIndex(protoreflect.ValueOfString("my_key").MapKey()),
				protopath.FieldAccess(fdNested2String),
			},
			want: []string{"STRING_TO_MAP", "my_key", "STRING_FIELD"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := pathSegments(tc.path)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fmt.Sprint(got) != fmt.Sprint(tc.want) {
				t.Errorf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestFindField(t *testing.T) {
	md := (&protoconfigv1.Test{}).ProtoReflect().Descriptor()

	tests := []struct {
		segment string
		want    string // field name, empty if no match
	}{
		{segment: "STRING_FIELD", want: "string_field"},
		{segment: "string_field", want: "string_field"},
		{segment: "String_Field", want: "string_field"},
		{segment: "STRINGFIELD", want: ""}, // no underscore, no match
		{segment: "UNKNOWN", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.segment, func(t *testing.T) {
			fd := findField(md, tc.segment)
			if tc.want == "" {
				if fd != nil {
					t.Errorf("want nil, got %s", fd.Name())
				}
				return
			}
			if fd == nil {
				t.Fatalf("want %s, got nil", tc.want)
			}
			if string(fd.Name()) != tc.want {
				t.Errorf("got %s want %s", fd.Name(), tc.want)
			}
		})
	}
}

// FuzzEnvKeyRoundTrip asserts that segments → envKey → splitEnvKey yields
// the original segments, for any valid input (no delimiter in segments).
func FuzzEnvKeyRoundTrip(f *testing.F) {
	f.Add("APP", "__", "A", "B", "C")
	f.Add("", "__", "ONLY", "", "")
	f.Add("PREFIX", "_X_", "a", "b", "c")

	f.Fuzz(func(t *testing.T, prefix, delim, s1, s2, s3 string) {
		var segs []string
		for _, s := range []string{s1, s2, s3} {
			if s == "" {
				continue
			}
			segs = append(segs, s)
		}
		if len(segs) == 0 {
			return
		}

		// Generate a key. If envKey rejects the input, skip — validation is
		// part of the contract, we only test round-trip for accepted inputs.
		key, err := envKey(prefix, delim, segs)
		if err != nil {
			return
		}

		got, matched := splitEnvKey(prefix, delim, key)
		if !matched {
			t.Fatalf("splitEnvKey did not match key %q (prefix=%q delim=%q)", key, prefix, delim)
		}
		if len(got) != len(segs) {
			t.Fatalf("len got=%d want=%d (segs=%v, key=%q, got=%v)", len(got), len(segs), segs, key, got)
		}
		for i := range got {
			if got[i] != segs[i] {
				t.Fatalf("segment %d: got %q want %q (key=%q)", i, got[i], segs[i], key)
			}
		}
	})
}
