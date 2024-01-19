# nanoc

nanoc is a compiler for NanoPack that compiles NanoPack schemas to code that enables portable data
serialization/deserialization in NanoPack.

## Requirements

- Go >= 1.20
- Formatter installed for the corresponding language:
    - C++: `clang-format`
    - Swift: `swift-format`
    - JavaScript/TypeScript: `prettier` will be used via `npx`, so `npm` needs to be available
