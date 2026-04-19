# protoconfig design

Load a protobuf message as layered configuration. The proto descriptor is the only source of truth for types, keys, cardinality, and value shape.

## Mental model

```
defaults (the zero proto you pass in)
   ↓ overridden by
YAML file(s)
   ↓ overridden by
env vars
   ↓
final cfg
```

Each layer overrides the previous. The final value has whatever the highest-precedence layer set, with everything else falling through.

## API

```go
cfg, err := protoconfig.Load(
    &myapp.Config{LogLevel: "info"},            // defaults (typed)
    protoconfig.FromYAMLFile("config.yaml"),
    protoconfig.FromEnv("MYAPP"),
    protoconfig.Strict(true),
)
```

- `Load[T proto.Message](defaults T, opts ...LoadOption) (T, error)` — the only entry point.
- Source constructors (`FromYAMLFile`, `FromEnv`, ...) and behavior toggles (`Strict`, ...) share one `LoadOption` interface.
- Sources apply in call order; later overrides earlier. Defaults are always first (implicit).

### Source options

`FromEnv(prefix string, opts ...EnvOption)` takes:
- `EnvironFunc(func() []string)` — override `os.Environ` (tests, sandboxing).
- `EnvDelimiter(string)` — default `"__"`.
- `TransformFunc(func(k, v string) (string, any))` — escape hatch for non-standard schemes.

`FromYAMLFile(path string, opts ...FileOption)` takes:
- `OptionalFile(bool)` — soft-fail on missing file.

### Errors

All failures wrap in `*LoadError`:
```go
type LoadError struct {
    Layer string  // "yaml", "env", "defaults"
    Key   string  // env var name or YAML path
    Path  string  // proto field path (human-readable)
    Inner error
}
```

Always `errors.As` to access fields.

## Env key scheme

**Path segments joined by `__` (configurable).** Single `_` is preserved inside field names so `string_field` stays readable.

| Proto path | Env key |
|---|---|
| `string_field` | `STRING_FIELD` |
| `message_field.string_field` | `MESSAGE_FIELD__STRING_FIELD` |
| `repeated_nested_message[0].string_field` | `REPEATED_NESTED_MESSAGE__0__STRING_FIELD` |
| `string_to_map["my_key"].string_field` | `STRING_TO_MAP__my_key__STRING_FIELD` |
| `list_of_ints[0]` | `LIST_OF_INTS__0` |
| `primitive_map["some-key"]` | `PRIMITIVE_MAP__some-key` |

With a prefix (`FromEnv("MYAPP")`): every key is prepended with `MYAPP__`.

**Bijection rule:** `pathToEnv(path) ↔ envToPath(env)` must round-trip for every valid path. Enforced by fuzz test.

**Map keys** preserve their original case and are literal. Map keys containing the delimiter (`__`) are an error.

**Field matching** is case-insensitive for field segments (so `STRING_FIELD`, `string_field`, `String_Field` all match). Map keys and enum values are case-sensitive (protojson convention).

## Value parsing

Every env value, regardless of kind, goes through the **protoyaml wrapper trick**: the value is wrapped into `{"<field_name>": <env_value>}` and fed to protoyaml via a synthetic parent message. We extract the parsed field.

This gives us, for free:
- Scalars (numbers, strings, bools, null)
- Enums by name (`EXAMPLE_VAL`) or by number (`1`)
- Bytes as base64
- Timestamps (RFC3339)
- Durations (`"1.5s"`)
- Wrappers — value set directly (no `value` field needed)
- FieldMask (comma-separated paths)
- Struct / Value / ListValue — any JSON/YAML blob
- Messages / lists / maps as JSON or YAML

Quote large int64/uint64 (`"9007199254740993"`) to avoid JSON precision loss.
Use `"NaN"`, `"Infinity"`, `"-Infinity"` for float specials.

## Merge / overlay semantics

### YAML over defaults

| Field kind | Semantics |
|---|---|
| Scalar | replace on set |
| Message (sub-struct) | recursive merge per sub-field |
| List | **replace whole list** (not `proto.Merge`'s append) |
| Map | merge by key (YAML keys replace; other keys preserved) |

Rationale: YAML lists can't syntactically address a single index, so the only coherent rule is "you wrote the list, that's the list." Maps naturally express per-key overrides. Messages naturally express per-field overrides.

### Env over result

| Env form | Meaning |
|---|---|
| `FIELD=<value>` | whole-value via protoyaml (replaces) |
| `LIST='[…]'` | whole-list via protoyaml (replaces) |
| `LIST__0=…`, `LIST__1=…` | per-index **overlay** onto existing list |
| `MAP='{…}'` | whole-map via protoyaml (replaces) |
| `MAP__key__subfield=…` | add-or-update single key |
| `MSG__field=…` | recurse into message (instantiate if nil, only if sub-env exists) |

**Out-of-bounds list index is an error.** To extend a list, use the whole-value form.

**Collision error:** for the same list or map, setting both whole-value and indexed/keyed env keys is an error (D1).

**Nil message auto-instantiation:** a nil message field is only instantiated if env sub-keys target it. This prevents infinite recursion on self-referencing types like `Test.message_field = Test`.

### WKTs

All WKTs are treated as **leaves for env addressing** in v1. Whole-value protoyaml parsing only; no sub-field env keys.

- `Timestamp`, `Duration`, `FieldMask`: ergonomic whole-value forms.
- `Wrappers`: inner value directly (`MY_STRINGVALUE=hello`, not `...__VALUE=hello`).
- `Struct`, `Value`, `ListValue`: whole-value protoyaml blob. Deep env addressing is v2.
- `Any`: whole-value only; requires concrete type registered in `protoregistry.GlobalTypes`. Low-prio TODO.

## Strictness

Single `Strict(bool)` option, default **lenient** (matches protobuf convention — unknown wire fields are ignored).

| Mode | YAML unknown field | Env unknown key (prefix matches) |
|---|---|---|
| lenient (default) | ignored | ignored |
| strict | error | error |

## Defaults are cloned

The `defaults` argument is cloned before any mutation. Caller's proto is never modified.

## Oneof

Setting any sub-field of an oneof arm selects that arm. Setting sub-fields of multiple arms (at any layer) is an error.

Cross-layer: higher-precedence layer wins. If YAML sets arm A and env sets arm B, env wins and A is cleared.

## Non-goals (v1)

- File watching / hot reload — needs a persistent `Loader[T]` handle; add when someone needs it.
- Flag source — trivial to add as another `From...` layer.
- Deep env addressing into Struct/Value/ListValue.
- `Any` with type discovery.
- Per-field `FieldOptions.env` annotations. Path-derived keys only.

## What doesn't live here

- **Validation.** Use `protovalidate` on the returned `cfg` when you need it. Keeping validation out of `Load` lets callers choose whether to fail startup or return a warning.
- **Secrets.** Env vars are fine for non-secrets; for secrets use a secret manager upstream and inject via env or YAML at deploy time.
