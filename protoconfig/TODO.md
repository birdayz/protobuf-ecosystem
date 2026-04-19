# protoconfig — rewrite plan

## Mission

Load a proto message as layered config: **defaults → YAML → env**. Proto descriptor is the schema-of-truth (types, enums, optional/required, repeated, map, nested, oneof, WKT). Predictable, descriptor-driven, no magic.

## Decisions (locked)

| # | Decision |
|---|---|
| D1 | Env list: both `LIST='[…]'` and `LIST__0=…`; collision = error; values parsed via protoyaml (superset of protojson) |
| D2 | Indexed env = in-bounds overlay (list) / add-or-update (map); list out-of-bounds index = error; whole-value form replaces |
| D3 | Delimiter `__` between path segments; single `_` preserved inside field names |
| D4 | `EnvPrefix` optional (koanf-shape); `TransformFunc`, `EnvironFunc` available as options |
| D5 | Descriptor-walking merger for YAML-over-defaults: scalars & lists replace, messages & maps merge per sub-key. Drop `proto.Merge`. |
| D6 | Delete `FieldOptions.env` annotation — one consistent scheme, no per-field overrides |
| D7 | Single `Strict(bool)` option governs YAML+env unknown-key handling; **default lenient** (matches protobuf norm of ignoring unknown fields) |
| D8 | All WKTs supported via whole-value protoyaml in v1. Deep env addressing into Struct/Value/ListValue = v2. Any = low-prio (requires type registration). |
| D9 | Clean rewrite, no deprecation shim |
| D10 | Stacked loader: `Load[T](defaults, opts ...LoadOption) (T, error)`. Sources and behavior toggles share `LoadOption`. |

## API sketch

```go
cfg, err := protoconfig.Load(
    &myapp.Config{LogLevel: "info"},             // defaults
    protoconfig.FromYAMLFile("config.yaml"),
    protoconfig.FromEnv("MYAPP",
        protoconfig.EnvironFunc(customEnviron),  // test-friendly
    ),
    protoconfig.Strict(true),                    // global
)
```

Errors: `LoadError{Layer, Key, Path, Inner}` so users know which layer blew up.

## Env key scheme

- Path segment separator: `__`
- Field names: snake_case preserved, uppercased (e.g. `string_field` → `STRING_FIELD`)
- List index: decimal, as its own segment (`LIST__0`)
- Map key: stringified (may preserve case; document), as its own segment (`MAP__my_key`)
- Map keys containing `__` → error

Round-trip property: path → env → path must be bijective. Enforced by fuzz test.

## Build order

- [x] 1. Write `DESIGN.md` pinning D1–D10 with worked examples (doubles as README)
- [x] 2. `envkey.go` + round-trip fuzz test — path ↔ env key bijection
- [x] 3. `scalar.go` — env string → `protoreflect.Value` via protoyaml `{"field":<val>}` trick; table test per scalar kind + WKT leaf
- [x] 4. `merge.go` — descriptor-walking YAML merger (per D5); tests for every kind
- [x] 5. `envapply.go` — scan env, match, apply overlay semantics (D1/D2)
- [x] 6. `load.go` + `options.go` — orchestration, `LoadOption`, `FromYAMLFile`, `FromEnv`, `Strict`, `EnvironFunc`
- [x] 7. Delete `FieldOptions.env` from `options.proto` and regenerate
- [x] 8. Integration test matrix — every field kind × every layer combo
- [x] 9. Fuzz `Load` — random env + YAML, never panic

## Status: v1 shipped

- 84% line coverage.
- 4 fuzzers, ~11M exec cumulative, no panics.
- Defensive `recover()` around `protoyaml.Unmarshal` (upstream panics on malformed input).

## Test matrix (in scope for v1)

1. envkey round-trip — every field kind + fuzz
2. Scalar matrix — 14+ scalar kinds × {defaults, YAML, env, all-three} layering
3. Numeric edges — max, min, overflow, negative-in-unsigned, NaN, Inf, exponent, leading-zero
4. Bool — `true/false/True/1/0/yes/no` (pin protoyaml's behavior)
5. Bytes — valid base64, invalid base64, empty (skip rule)
6. String — empty, newlines, YAML specials, unicode, whitespace
7. Enum — by name, by number, unknown → error, case sensitivity pinned
8. Repeated — whole-value + indexed, replace vs overlay, out-of-bounds error, size 0
9. Maps — string/int/bool keys; nested message values; key containing `__` → error; add-or-update
10. Nested recursion — self-ref stays nil without sub-env; deep nesting
11. Oneof — each arm, two-arm conflict → error, env-beats-YAML arm switch
12. WKTs — Timestamp, Duration, all 9 Wrappers, FieldMask, Struct, Value, ListValue (whole-value)
13. YAML replace-merge — default list + YAML list → YAML list (NOT append)
14. Precedence — all three layers, env wins
15. Errors — file missing, YAML invalid, env wrong type, strict mode unknown key, typed `LoadError`
16. Defaults immutability — caller's `*proto.Message` not mutated

## Out of scope for v1

- Deep env addressing into `Struct` / `Value` / `ListValue` (fast-follow)
- `google.protobuf.Any` with type discovery (low-prio; whole-value works if registered)
- File watching / hot reload (needs persistent handle type)
- Flag source (`FromFlags`)
- Multi-file merge (can compose via multiple `FromYAMLFile` calls already)
- Extracting the descriptor merger into a standalone package (defer until another caller needs it)

## Open questions parked for later

- Reintroduce `FieldOptions.env_name` as opt-in alias/opt-out? Only if someone needs migration from legacy env-var names.
- Should `Strict(true)` also enable fuzzy-match "did-you-mean" hints in error messages?
- Hot-reload: `New[T](defaults, ...) *Loader[T]` with `.Watch(cb)` and `.Current() T`.
