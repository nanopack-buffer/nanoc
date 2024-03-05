# nanoc

nanoc is a compiler for NanoPack that compiles NanoPack schemas to code that enables portable data
serialization/deserialization in NanoPack.

## Requirements

- Go >= 1.20
- Formatter installed for the corresponding language:
    - C++: `clang-format`
    - Swift: `swift-format`
    - JavaScript/TypeScript: `prettier` will be used via `npx`, so `npm` needs to be available

## Installation

`nanoc` is available as pre-built binaries on the [GitHub Releases page](https://github.com/poly-gui/nanoc/releases).

### Compiling the Compiler

```
# Assuming the cwd is in the project directory where go.mod is
go install nanoc/cmd/nanoc
```

The above compiles the source code, produces an executable called `nanoc`, and then moves (installs) it to `GOBIN`, or `GOPATH/bin` if `GOBIN` is not set. **Make sure either is in PATH**.

## Usage

`nanoc` has the following arguments:
- `--language` specifies the language of the generated code. Can be `c++`, `swift`, or `ts`.
- `--factory-out` (optional) specifies the folder in which the message factory code file should be put.
- The list of relative paths to the NanoPack schemas that should be compiled.

### Example

Compiling `person.yml` and `car.yml` in the current  directory to C++:

```
nanoc --language=c++ ./person.yml ./car.yml
```

Compiling `person.yml` and `car.yml` in the current directory to TypeScript, and generate the message factory in the current directory:

```
nanoc --language=ts --factory-out=. ./person.yml ./car.yml
```
