package protoconfig

import (
	"errors"
	"testing"

	protoconfigv1 "github.com/birdayz/protobuf-ecosystem/protoconfig/proto/gen/go/protoconfig/v1"
)

func TestEnvDelimiter_Custom(t *testing.T) {
	// Use "::" as a custom delimiter.
	cfg, err := Load(&protoconfigv1.Test{},
		FromEnv("APP",
			EnvDelimiter("::"),
			EnvironFunc(func() []string {
				return []string{`APP::STRING_FIELD="x"`}
			}),
		),
		Strict(true),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.StringField != "x" {
		t.Errorf("got %q", cfg.StringField)
	}
}

func TestError_UnwrapExposesInner(t *testing.T) {
	inner := errors.New("boom")
	e := &Error{Layer: "env:APP", Inner: inner}
	if errors.Unwrap(e) != inner {
		t.Errorf("Unwrap didn't return inner")
	}
	if !errors.Is(e, inner) {
		t.Errorf("errors.Is failed to find inner")
	}
}

func TestFromEnv_NoPrefix_AllEnvConsidered(t *testing.T) {
	// Without prefix, matching env vars apply. Unknown vars ignored by default.
	cfg, err := Load(&protoconfigv1.Test{},
		FromEnv("", EnvironFunc(func() []string {
			return []string{
				`STRING_FIELD="ok"`,
				"PATH=/usr/bin",   // unrelated — silently ignored (lenient default)
				"HOME=/home/user", // unrelated — silently ignored
			}
		})),
	)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.StringField != "ok" {
		t.Errorf("got %q", cfg.StringField)
	}
}
