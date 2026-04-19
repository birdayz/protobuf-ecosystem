package protoconfig

import "os"

// LoadOption configures a call to Load. Both source constructors (e.g.
// FromYAMLFile, FromEnv) and behavior toggles (e.g. Strict) implement this
// interface.
type LoadOption interface {
	applyLoad(*loadConfig)
}

type loadConfig struct {
	strict  bool
	sources []sourceFn
}

// sourceFn applies one layer over the current config. Invoked in registration
// order (later sources overlay earlier ones).
type sourceFn struct {
	label string // "yaml:<path>", "env:<prefix>", etc. — for error wrapping
	fn    func(ctx *loadContext) error
}

// loadContext carries per-call state to sources so they don't each need to
// know the full loadConfig.
type loadContext struct {
	// The mutable proto reflect handle we're overlaying into.
	msg any // protoreflect.Message; kept as any to avoid import cycle
	// strict, set from LoadOption list.
	strict bool
}

// funcOption adapts a closure to LoadOption.
type funcOption func(*loadConfig)

func (f funcOption) applyLoad(c *loadConfig) { f(c) }

// Strict controls how unknown fields are handled across YAML and env sources.
// Default is lenient (ignore unknown) — matching protobuf wire-compat norms.
// When true, unknown YAML fields or unknown prefixed env keys are errors.
func Strict(v bool) LoadOption {
	return funcOption(func(c *loadConfig) { c.strict = v })
}

// --- env source ---

// EnvOption configures FromEnv.
type EnvOption interface {
	applyEnv(*envOpts)
}

type envOpts struct {
	environ   func() []string
	delim     string
	transform func(k, v string) (string, string, bool) // (newKey, newVal, keep)
}

func defaultEnvOpts() *envOpts {
	return &envOpts{
		environ: os.Environ,
		delim:   defaultDelim,
	}
}

type envOptFn func(*envOpts)

func (f envOptFn) applyEnv(o *envOpts) { f(o) }

// EnvironFunc overrides the source of env vars. Lets tests inject a fixed set
// without touching os.Environ.
func EnvironFunc(fn func() []string) EnvOption {
	return envOptFn(func(o *envOpts) { o.environ = fn })
}

// EnvDelimiter overrides the path segment delimiter. Default is "__".
// Changing this is rarely necessary and must satisfy the constraints in
// DESIGN.md §env-key-scheme (no segment may start/end with a delimiter
// character).
func EnvDelimiter(d string) EnvOption {
	return envOptFn(func(o *envOpts) { o.delim = d })
}

// EnvTransformFunc runs on each (rawKey, rawVal) after prefix filtering.
// Return keep=false to drop the entry. Return a different key/value to
// rewrite it. Useful for legacy naming adapters.
func EnvTransformFunc(fn func(k, v string) (string, string, bool)) EnvOption {
	return envOptFn(func(o *envOpts) { o.transform = fn })
}

// --- file source ---

// FileOption configures FromYAMLFile.
type FileOption interface {
	applyFile(*fileOpts)
}

type fileOpts struct {
	optional bool
}

type fileOptFn func(*fileOpts)

func (f fileOptFn) applyFile(o *fileOpts) { f(o) }

// Optional marks a YAML source as soft: if the file is missing, Load proceeds
// without error. Parse errors are always fatal.
func Optional() FileOption {
	return fileOptFn(func(o *fileOpts) { o.optional = true })
}
