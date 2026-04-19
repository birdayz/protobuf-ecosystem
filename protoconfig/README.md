# protoconfig

Typed, layered configuration for Go apps — backed by a protobuf message.

```go
cfg, err := protoconfig.Load(
    &myapp.Config{LogLevel: "info"},              // defaults
    protoconfig.FromYAMLFile("config.yaml"),      // layered on top
    protoconfig.FromEnv("MYAPP"),                 // env wins
)
```

The proto is your schema. Types, enums, optional/required, lists, maps, nested messages, oneofs, and well-known types all come from the descriptor — no struct tags, no mapstructure decode hooks, no second source of truth.

## Why

- **One schema.** Your `.proto` describes the config shape. Use it in services, clients, docs, and [`protovalidate`](https://github.com/bufbuild/protovalidate) the same way.
- **Predictable layering.** Defaults → YAML → env. Each layer is additive with precise, documented merge rules. No surprises from silent appends or prefix collisions.
- **Typed, per-layer errors.** `*protoconfig.Error` tells you which layer failed and why.
- **Test-friendly.** `EnvironFunc` swaps `os.Environ` with a fake list. No `t.Setenv` juggling.

## Install

```
go get github.com/birdayz/protobuf-ecosystem/protoconfig
```

Requires Go 1.23+.

## 30-second tour

```proto
// config.proto
syntax = "proto3";
package myapp;

message Config {
  string log_level = 1;
  int32 port = 2;
  repeated string allowed_origins = 3;
  map<string, string> labels = 4;
  google.protobuf.Timestamp started_at = 5;
}
```

```go
import (
    "github.com/birdayz/protobuf-ecosystem/protoconfig"
    "myapp/gen/myapp"
)

func main() {
    cfg, err := protoconfig.Load(
        &myapp.Config{LogLevel: "info", Port: 8080},
        protoconfig.FromYAMLFile("config.yaml", protoconfig.Optional()),
        protoconfig.FromEnv("MYAPP"),
    )
    // cfg is *myapp.Config — typed. Use it directly.
}
```

```yaml
# config.yaml
port: 9000
allowed_origins: ["https://example.com"]
labels:
  team: platform
```

```bash
export MYAPP__LOG_LEVEL='"debug"'
export MYAPP__LABELS__region='"us-east"'   # adds a label; team stays
export MYAPP__STARTED_AT='"2024-02-03T17:00:00Z"'  # RFC3339 for Timestamp
```

## API

```go
func Load[T proto.Message](defaults T, opts ...LoadOption) (T, error)
```

`LoadOption` is one interface both source constructors and behavior toggles satisfy. Order matters: sources later in the list override earlier ones.

### Sources

| Constructor | What it reads |
|---|---|
| `FromYAMLFile(path, ...FileOption)` | A YAML file. Parsed with `protoyaml` (superset of `protojson`). |
| `FromEnv(prefix, ...EnvOption)` | Environment variables matching `<prefix>__<path>`. |

### Behavior toggles

| Option | Effect |
|---|---|
| `Strict(bool)` | When true, unknown YAML fields and unknown prefixed env keys are errors. Default lenient (matches protobuf wire-compat norms). |

### Source options

| Option | Effect |
|---|---|
| `Optional()` | `FromYAMLFile`: soft-fail if the file is missing. Parse errors still fatal. |
| `EnvironFunc(fn)` | `FromEnv`: replaces `os.Environ`. Essential for tests. |
| `EnvDelimiter(s)` | `FromEnv`: replaces the default `__` delimiter. |
| `EnvTransformFunc(fn)` | `FromEnv`: rewrite or drop (k, v) pairs. Escape hatch for legacy naming. |

## Env var scheme

Path segments are joined by `__` (double underscore). Field names are uppercased; single `_` inside field names is preserved.

| Proto path | Env var |
|---|---|
| `log_level` | `LOG_LEVEL` |
| `database.host` | `DATABASE__HOST` |
| `allowed_origins` (whole list) | `ALLOWED_ORIGINS='["a","b"]'` |
| `allowed_origins[0]` | `ALLOWED_ORIGINS__0='"a"'` |
| `labels["team"]` | `LABELS__team='"platform"'` |
| `started_at` (Timestamp) | `STARTED_AT='"2024-02-03T17:00:00Z"'` |

With a prefix: `FromEnv("MYAPP")` expects `MYAPP__LOG_LEVEL`, `MYAPP__DATABASE__HOST`, etc.

### Values are protoyaml fragments

Env values go through `protoyaml` (with a `protojson` fallback for WKT special forms). Whatever you'd write in the YAML config file, you can write as an env value:

```
MYAPP__PORT=9000                          # bare number
MYAPP__LOG_LEVEL='"debug"'                # quoted string
MYAPP__ENABLED=true                       # bool (strict: only true/false)
MYAPP__LABELS='{"team":"platform"}'       # inline JSON object
MYAPP__ALLOWED_ORIGINS='["a","b"]'        # inline JSON array
MYAPP__TIMEOUT='"1.5s"'                   # Duration
MYAPP__STARTED_AT='"2024-02-03T17:00:00Z"'  # Timestamp
MYAPP__MASK='"a,b.c"'                     # FieldMask
```

Gotchas:
- **Quote large int64/uint64** (`'"9223372036854775807"'`) to avoid JSON float precision loss.
- **Float specials** are quoted strings: `'"NaN"'`, `'"Infinity"'`, `'"-Infinity"'`.
- **Bool is strict**: only `true`/`false`. Not `1`, `yes`, `True`.
- **Bytes** are base64-encoded strings.
- **Enum names are case-sensitive**: `"EXAMPLE_VAL"` not `"example_val"`. Numbers also work: `1`.
- **Empty env (`VAR=""`) is skipped** — not treated as "clear the default."

## Merge semantics

### YAML over defaults

| Field kind | What happens |
|---|---|
| Scalar | YAML replaces default |
| Message (sub-struct) | recursive merge |
| List | **YAML replaces whole list** (differs from `proto.Merge`) |
| Map | merge by key — YAML keys overwrite, others preserved |

### Env over previous layer

Env has two forms per list/map: **whole-value** (replaces) and **addressed** (overlay).

| Form | Meaning |
|---|---|
| `LIST='[…]'` | replace the whole list |
| `LIST__0=…` | overlay element 0 (must exist — out of bounds = error) |
| `MAP='{…}'` | replace the whole map |
| `MAP__key=…` | add or update one key |
| `MSG__field=…` | recurse into a sub-field (nil messages auto-instantiate) |

Setting both whole-value and addressed forms for the same list or map is an error — pick one.

### WKTs

Well-known types are **leaves** for env addressing — you set them as a single value:

```
MY_TIMESTAMP='"2024-02-03T17:00:00Z"'
MY_DURATION='"1.5s"'
MY_FIELD_MASK='"a.b,c.d"'
MY_STR_WRAPPER='"hello"'           # Wrappers take the inner value directly
MY_STRUCT='{"foo": {"x": 1}}'      # Struct/Value/ListValue: any JSON blob
```

(Deep env addressing into `Struct`/`ListValue` and full `Any` support are planned for v2.)

## Recipes

### Tests without touching os.Environ

```go
cfg, err := protoconfig.Load(&myapp.Config{},
    protoconfig.FromEnv("MYAPP",
        protoconfig.EnvironFunc(func() []string {
            return []string{
                `MYAPP__LOG_LEVEL="debug"`,
                "MYAPP__PORT=9090",
            }
        }),
    ),
)
```

### Optional config file

```go
protoconfig.FromYAMLFile("/etc/myapp.yaml", protoconfig.Optional())
```

Missing file → no error. Malformed file → fatal error.

### Strict mode in CI, lenient in prod

```go
strict := os.Getenv("CI") != ""
cfg, err := protoconfig.Load(defaults,
    protoconfig.FromYAMLFile("config.yaml"),
    protoconfig.FromEnv("MYAPP"),
    protoconfig.Strict(strict),
)
```

### Adapting legacy env names

```go
protoconfig.FromEnv("MYAPP",
    protoconfig.EnvTransformFunc(func(k, v string) (string, string, bool) {
        // Rename MYAPP__OLD_NAME -> MYAPP__NEW_NAME
        if k == "MYAPP__OLD_NAME" {
            return "MYAPP__NEW_NAME", v, true
        }
        return k, v, true
    }),
)
```

### Validate after loading

```go
import "github.com/bufbuild/protovalidate-go"

cfg, err := protoconfig.Load(defaults, ...)
if err != nil {
    log.Fatal(err)
}
if err := protovalidate.Validate(cfg); err != nil {
    log.Fatalf("config failed validation: %v", err)
}
```

### Error inspection

```go
cfg, err := protoconfig.Load(...)
var lerr *protoconfig.Error
if errors.As(err, &lerr) {
    log.Printf("layer %s failed: %v", lerr.Layer, lerr.Inner)
}
```

## FAQ

**Does it modify my `defaults` argument?** No. `defaults` is cloned before any mutation.

**Can a nil nested message auto-instantiate?** Yes, but only if env vars target sub-fields of it. This prevents infinite recursion on self-referencing types like `Config.child = Config`.

**What if two env keys conflict?** Setting `LIST` and `LIST__0` for the same field → error. Setting `MY_ONEOF_A` and `MY_ONEOF_B` for different arms of the same oneof → error.

**Can I use proto field JSON names (camelCase)?** Env keys are path-derived from the proto field names (snake_case, uppercased). YAML accepts both snake_case and camelCase (protojson convention).

**Does unknown env var handling affect unrelated vars?** No — only env vars matching the prefix are considered. Without a prefix, only env keys that map to real proto fields are touched; everything else is ignored (even in strict mode at the top level — strictness only applies *within* your namespace).

**Can I watch the file for changes?** Not yet. Reload is on the v2 list — today, re-call `Load`.

## Status

v1. API shape is stable. `DESIGN.md` pins the decisions; `TODO.md` tracks what's done and what's parked.

## License

Same as the parent repository.
