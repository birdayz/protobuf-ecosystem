# Protofieldmask

Or, protoiter?


## Why
- Iterate over _ALL fields. Even if not set, or different protos, mismatching types, ...
- Write advanced fieldmask merging
- For ^, make it possible to easily mix in specific logic, eg. error on IMMUTABLE googleapis annotation
- Support fieldmask, by also show presence / absence of message,map, etc, so it can be initialized / reset to zero / ... if required by fieldmask.
- Fieldmask can be impelemented on top
- Make it easy to write a protodefault package.
- Write advanced and better protocmp / assertion lib, with rich info about type etc
- Maybe implement field manager like in SSA ?
- Field types to not need to be IDENTICAL, but being EQUIVALENT / COMPATIBLE (wire, either binary or JSON, or based on field names?? option to use number, field name, or json) is good enough
- Maybe fast. bench how good it is. Traverse the tree only once. Maybe add advanced prefix support w/ a fieldmask provider interface.
- Can be used to check if two types are "fieldmask compatible" ? buf breaking is too strict. Different nested struct, that is compatible, is considered breaking, even if wire compatible.
- Protovalidate skip checks if not in fieldmask?
- Write - on top - extremelt simpe Merge function, that takes two protos. Maybe with opts
- Move this to protoiter

## TODO
- Special types/pitfalls: Map, List, Any.
- Handle one message being empty
- Fix infinite recursion. Do not enter empty messages?
