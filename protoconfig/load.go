// Package protoconfig loads a protobuf message as layered configuration:
// defaults → YAML → env vars. The proto descriptor is the sole source of
// truth for types, keys, and cardinality.
//
// Quick start:
//
//	cfg, err := protoconfig.Load(
//	    &myapp.Config{LogLevel: "info"},
//	    protoconfig.FromYAMLFile("config.yaml"),
//	    protoconfig.FromEnv("MYAPP"),
//	)
//
// See DESIGN.md for the full semantics.
package protoconfig

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"buf.build/go/protoyaml"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Error is returned from Load when a specific layer fails. Always use
// errors.As to introspect.
type Error struct {
	// Layer identifies which source triggered the failure, e.g. "yaml:path",
	// "env:PREFIX", or "defaults".
	Layer string
	// Inner is the underlying error.
	Inner error
}

func (e *Error) Error() string {
	return fmt.Sprintf("protoconfig: %s: %v", e.Layer, e.Inner)
}

func (e *Error) Unwrap() error { return e.Inner }

// Load applies the given option/source list to a clone of defaults and
// returns the result. Defaults are never mutated.
func Load[T proto.Message](defaults T, opts ...LoadOption) (T, error) {
	var zero T
	cfg := proto.Clone(defaults).(T)

	lc := &loadConfig{}
	for _, o := range opts {
		o.applyLoad(lc)
	}

	ctx := &loadContext{msg: cfg.ProtoReflect(), strict: lc.strict}
	for _, src := range lc.sources {
		if err := src.fn(ctx); err != nil {
			return zero, &Error{Layer: src.label, Inner: err}
		}
	}
	return cfg, nil
}

// FromYAMLFile loads a YAML file as a layer. Lists and maps replace per-field;
// messages and map values merge recursively (see DESIGN.md §merge).
func FromYAMLFile(path string, opts ...FileOption) LoadOption {
	fo := &fileOpts{}
	for _, o := range opts {
		o.applyFile(fo)
	}
	return funcOption(func(c *loadConfig) {
		c.sources = append(c.sources, sourceFn{
			label: "yaml:" + path,
			fn: func(ctx *loadContext) error {
				b, err := os.ReadFile(path)
				if err != nil {
					if fo.optional && errors.Is(err, os.ErrNotExist) {
						return nil
					}
					return err
				}
				msg := ctx.msg.(protoreflect.Message)
				// Parse into a fresh message of the same type so we can merge
				// with our own overlay semantics (not proto.Merge).
				fresh := msg.New()
				if err := safeProtoyamlUnmarshal(protoyaml.UnmarshalOptions{DiscardUnknown: !ctx.strict}, b, fresh.Interface()); err != nil {
					return err
				}
				mergeFrom(msg, fresh)
				return nil
			},
		})
	})
}

// FromEnv reads environment variables matching the given prefix and applies
// them as the next layer. See DESIGN.md §env-key-scheme for the encoding.
// An empty prefix means every env var is a candidate (filtered later by
// whether the key segments match proto fields).
func FromEnv(prefix string, opts ...EnvOption) LoadOption {
	eo := defaultEnvOpts()
	for _, o := range opts {
		o.applyEnv(eo)
	}
	return funcOption(func(c *loadConfig) {
		c.sources = append(c.sources, sourceFn{
			label: "env:" + prefix,
			fn: func(ctx *loadContext) error {
				msg := ctx.msg.(protoreflect.Message)
				entries, err := collectEnvEntries(prefix, eo)
				if err != nil {
					return err
				}
				return applyEnv(msg, entries, ctx.strict)
			},
		})
	})
}

// collectEnvEntries scans the environment via eo.environ(), filters by
// prefix, runs the optional transform, splits on the delimiter, and returns
// the parsed entries.
func collectEnvEntries(prefix string, eo *envOpts) ([]envEntry, error) {
	var out []envEntry
	for _, kv := range eo.environ() {
		rawKey, rawVal, ok := strings.Cut(kv, "=")
		if !ok {
			continue
		}
		segments, matched := splitEnvKey(prefix, eo.delim, rawKey)
		if !matched {
			continue
		}
		if eo.transform != nil {
			newKey, newVal, keep := eo.transform(rawKey, rawVal)
			if !keep {
				continue
			}
			// If key changed, re-split.
			if newKey != rawKey {
				segments, matched = splitEnvKey(prefix, eo.delim, newKey)
				if !matched {
					continue
				}
				rawKey = newKey
			}
			rawVal = newVal
		}
		if len(segments) == 0 {
			// prefix-only var has nothing to apply
			continue
		}
		out = append(out, envEntry{
			segments: segments,
			value:    rawVal,
			rawKey:   rawKey,
		})
	}
	return out, nil
}
