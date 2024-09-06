package cxxgen

import (
	"nanoc/internal/datatype"
	"nanoc/internal/npschema"
)

type messageHeaderFileTemplateInfo struct {
	Namespace         string
	MessageName       string
	HasParentMessage  bool
	ParentMessageName string
	TypeID            datatype.TypeID
	IncludeGuardName  string
	LibraryImports    []string
	RelativeImports   []string
	IsInherited       bool

	FieldGetters          []string
	FieldDeclarationLines []string
	ConstructorParameters []string
}

type messageImplFileTemplateInfo struct {
	Namespace         string
	HeaderName        string
	MessageName       string
	HasParentMessage  bool
	ParentMessageName string
	RelativeImports   []string

	ConstructorParameters []string
	SuperConstructorArgs  []string
	FieldInitializers     []string
	FieldGetters          []string

	ReadPtrStart           int
	FieldReadCodeFragments []string

	HeaderSize              int
	InitialWriteBufferSize  int
	FieldWriteCodeFragments []string
}

type childMessageFactoryHeaderFileTemplateInfo struct {
	Namespace           string
	IncludeGuardName    string
	MessageName         string
	MessageHeaderName   string
	FactoryFunctionName string
}

type childMessageFactoryImplFileTemplateInfo struct {
	Namespace               string
	Schema                  *npschema.Message
	ChildMessageImportPaths []string
	HeaderName              string
	FactoryFunctionName     string
}

type messageFactoryHeaderFileTemplateInfo struct {
	Namespace string
}

type messageFactoryImplFileTemplateInfo struct {
	Namespace          string
	MessageImportPaths []string
	MessageSchemas     []*npschema.Message
}

type enumHeaderFileInfo struct {
	Namespace        string
	Schema           *npschema.Enum
	BackingTypeName  string
	IsIntType        bool
	MemberNames      []string
	IncludeGuardName string
}

type serviceHeaderFileInfo struct {
	RelativeImports  []string
	LibraryImports   []string
	Namespace        string
	IncludeGuardName string
	Schema           *npschema.Service
}

type serviceImplFileInfo struct {
	Namespace  string
	HeaderName string
	Schema     *npschema.Service
}

const (
	templateNameMessageHeaderFile             = "CxxMessageHeaderFile"
	templateNameMessageImplFile               = "CxxMessageImplFile"
	templateNameMessageFactoryHeaderFile      = "CxxMessageFactoryHeaderFile"
	templateNameMessageFactoryImplFile        = "CxxMessageFactoryImplFile"
	templateNameChildMessageFactoryHeaderFile = "CxxChildMessageFactoryHeaderFile"
	templateNameChildMessageFactoryImplFile   = "CxxChildMessageFactoryImplFile"
	templateNameEnumHeaderFile                = "CxxEnumHeaderFile"
	templateNameServiceHeaderFile             = "CxxServiceHeaderFile"
	templateNameServiceImplFile               = "CxxServiceImplFile"
)

const (
	extHeaderFile = ".np.hxx"
	extImplFile   = ".np.cxx"
)

const cxxSymbolMemberOf = "::"

const fileNameMessageFactory = "nanopack_message_factory"

const messageHeaderFile = `// AUTOMATICALLY GENERATED BY NANOC

#ifndef {{.IncludeGuardName}}
#define {{.IncludeGuardName}}

{{- range .LibraryImports}}
#include <{{.}}>
{{- end}}
{{- if not .HasParentMessage}}
#include <nanopack/message.hxx>
{{- end}}
#include <nanopack/nanopack.hxx>
#include <nanopack/reader.hxx>

{{range .RelativeImports}}
#include "{{.}}"
{{- end}}

{{if .Namespace}}namespace {{.Namespace}} { {{- end}}

struct {{.MessageName}} : {{if .HasParentMessage}}{{.ParentMessageName}}{{else}}NanoPack::Message{{end}} {
  static constexpr NanoPack::TypeId TYPE_ID = {{.TypeID}};

  {{range .FieldDeclarationLines}}{{.}}{{end}}

  {{.MessageName}}() = default;

  {{$l := len .ConstructorParameters}}{{if eq $l 1}}explicit {{end}}
  {{- .MessageName}}({{range $i, $v := .ConstructorParameters}}{{if $i}}, {{end}}{{$v}}{{end}});

  size_t read_from(NanoPack::Reader &reader);

  {{range .FieldGetters}}
  {{.}}
  {{- end}}

  size_t write_to(NanoPack::Writer &writer, int offset) const override;

  [[nodiscard]] NanoPack::TypeId type_id() const override;

  [[nodiscard]] size_t header_size() const override;
};

{{if .Namespace}}} // namespace {{.Namespace}}{{end}}

#endif
`

const messageImplFile = `// AUTOMATICALLY GENERATED BY NANOC

#include <nanopack/reader.hxx>
#include <nanopack/writer.hxx>

{{range .RelativeImports}}
#include "{{.}}"
{{- end}}

#include "{{.HeaderName}}"

{{if .Namespace}}{{.Namespace}}::{{end}}{{.MessageName}}::{{.MessageName}}({{join .ConstructorParameters ", "}}) :
	{{if .HasParentMessage}}{{.ParentMessageName}}({{join .SuperConstructorArgs ", "}}), {{end}}
	{{- join .FieldInitializers ", "}} {}

size_t {{if .Namespace}}{{.Namespace}}::{{end}}{{.MessageName}}::read_from(NanoPack::Reader &reader) {
	uint8_t *buf = reader.buffer;
	int ptr = {{.ReadPtrStart}};

	{{range .FieldReadCodeFragments}}
	{{.}}

	{{end}}

	return ptr;
}

{{range .FieldGetters}}
{{.}}
{{- end}}

NanoPack::TypeId {{if .Namespace}}{{.Namespace}}::{{end}}{{.MessageName}}::type_id() const {
  return TYPE_ID;
}

size_t {{if .Namespace}}{{.Namespace}}::{{end}}{{.MessageName}}::header_size() const {
  return {{.HeaderSize}};
}

size_t {{if .Namespace}}{{.Namespace}}::{{end}}{{.MessageName}}::write_to(NanoPack::Writer &writer, int offset) const {
	const size_t writer_size_before = writer.size();

	writer.reserve_header({{.HeaderSize}});

	writer.write_type_id(TYPE_ID, offset);

	{{range .FieldWriteCodeFragments}}
	{{.}}

	{{end}}

	return writer.size() - writer_size_before;
}
`

const childMessageFactoryHeaderFile = `// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#ifndef {{.IncludeGuardName}}
#define {{.IncludeGuardName}}

#include <nanopack/reader.hxx>
#include <memory>
#include "{{.MessageHeaderName}}"

{{if .Namespace}}namespace {{.Namespace}} { {{- end}}

std::unique_ptr<{{.MessageName}}> {{.FactoryFunctionName}}(NanoPack::Reader &reader, size_t &bytes_read);

{{if .Namespace}}} // namespace {{.Namespace}}{{end}}

#endif
`

const childMessageFactoryImplFile = `// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#include "{{.HeaderName}}"
{{range .ChildMessageImportPaths}}
#include "{{.}}"
{{- end}}

std::unique_ptr<{{if .Namespace}}{{.Namespace}}::{{end}}{{.Schema.Name}}> {{if .Namespace}}{{.Namespace}}::{{end}}{{.FactoryFunctionName}}(NanoPack::Reader &reader, size_t &bytes_read) {
  switch (reader.read_type_id()) {
  case {{.Schema.TypeID}}: {
	auto ptr = std::make_unique<{{.Schema.Name}}>();
	bytes_read = ptr->read_from(reader);
	return ptr;
  }
  {{range .Schema.ChildMessages -}}
  case {{.TypeID}}: {
	auto ptr = std::make_unique<{{.Name}}>();
	bytes_read = ptr->read_from(reader);
	return ptr;
  }
  {{end -}}
  default: return nullptr;
  }
}
`

const messageFactoryHeaderFile = `// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#ifndef NANOPACK_MESSAGE_FACTORY_HXX
#define NANOPACK_MESSAGE_FACTORY_HXX

#include <nanopack/message.hxx>
#include <nanopack/reader.hxx>
#include <memory>

{{if .Namespace}}namespace {{.Namespace}} { {{- end}}

std::unique_ptr<NanoPack::Message> make_nanopack_message(NanoPack::Reader &reader);
std::unique_ptr<NanoPack::Message> make_nanopack_message(NanoPack::Reader &reader, size_t &bytes_read);

{{if .Namespace}}} // namespace {{.Namespace}}{{end}}

#endif
`

const messageFactoryImplFile = `// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#include "nanopack_message_factory.np.hxx"

{{range .MessageImportPaths}}
#include "{{.}}"
{{- end}}

std::unique_ptr<NanoPack::Message> {{if .Namespace}}{{.Namespace}}::{{end}}make_nanopack_message(NanoPack::Reader &reader) {
  size_t _;
  return make_nanopack_message(reader, _);
}

std::unique_ptr<NanoPack::Message> {{if .Namespace}}{{.Namespace}}::{{end}}make_nanopack_message(NanoPack::Reader &reader, size_t &bytes_read) {
  switch (reader.read_type_id()) {
  {{range .MessageSchemas -}}
  case {{.TypeID}}: {
	auto ptr = std::make_unique<{{.Name}}>();
	bytes_read = ptr->read_from(reader);
	return ptr;
  }
  {{- end}}
  default: return nullptr;
  }
}
`

const enumHeaderFile = `// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

#ifndef {{.IncludeGuardName}}
#define {{.IncludeGuardName}}

#include <array>
#include <stdexcept>
{{- if eq .BackingTypeName "std::string_view"}}
#include <string_view>
#include <unordered_map>
{{else if .IsIntType}}
#include <cstdint>
{{- end}}

{{if .Namespace}}namespace {{.Namespace}} { {{- end}}

class {{.Schema.Name}} {
public:
  enum {{.Schema.Name}}Member {
  	{{- range .MemberNames}}
    {{.}},
    {{- end}}
  };

private:
  constexpr static std::array<{{.BackingTypeName}}, {{len .MemberNames}}> values = { {{range .Schema.Members}}{{.ValueLiteral}}, {{end}} };
  {{- if eq .BackingTypeName "std::string_view" -}}
  inline static std::unordered_map<std::string_view, {{.Schema.Name}}Member> lookup{
    {{- $m := .MemberNames -}}
    {{range $i, $e := .Schema.Members}}{{if $i}}, {{end}}{ {{$e.ValueLiteral}}, {{index $m $i}} }{{end}}
  };
  {{- end -}}

  {{.Schema.Name}}Member enum_value;
  {{.BackingTypeName}} _value;

public:
  {{.Schema.Name}}() = default;

  {{if eq .BackingTypeName "std::string_view" -}}
  explicit {{.Schema.Name}}(const {{.BackingTypeName}} &value) : enum_value(lookup.find(value)->second), _value(values[enum_value]) {}
  {{- else -}}
  explicit {{.Schema.Name}}(const {{.BackingTypeName}} &value) {
    switch (value) {
      {{range .Schema.Members}}case {{.ValueLiteral}}:
        enum_value = {{.Name}};
        break;{{end}}
      default: throw std::runtime_error("invalid value for enum {{.Schema.Name}}");
    }
    _value = values[enum_value];
  }
  {{- end}}

  constexpr {{.Schema.Name}}({{.Schema.Name}}Member member) : enum_value(member), _value(values[member]) {}

  [[nodiscard]] constexpr const {{.BackingTypeName}} &value() const { return _value; }

  constexpr operator {{.Schema.Name}}Member() const { return enum_value; }

  explicit operator bool() const = delete;
};

{{if .Namespace}}} // namespace {{.Namespace}}{{end}}

#endif
`

const serviceHeaderFile = `#ifndef {{.IncludeGuardName}}
#define {{.IncludeGuardName}}

#include <unordered_map>
#include <string_view>
#include <future>
#include <nanopack/rpc.hxx>
{{- range .LibraryImports}}
#include <{{.}}>
{{- end}}

{{- range .RelativeImports}}
#include "{{.}}"
{{- end}}

{{if .Namespace}}namespace {{.Namespace}} { {{- end}}

class {{.Schema.Name}}ServiceServer : public NanoPack::RpcServer {
	std::unordered_map<std::string_view, MethodCallResult ({{.Schema.Name}}ServiceServer::*)(uint8_t *, size_t, NanoPack::MessageId)> handlers;

	MethodCallResult on_method_call(const std::string_view &method, uint8_t *request_data, size_t offset, NanoPack::MessageId msg_id) override;

	{{- range .Schema.DeclaredFunctions}}
	MethodCallResult _{{snake .Name}}(uint8_t *request_data, size_t offset, NanoPack::MessageId msg_id);
	virtual {{if .ReturnType}}{{typeDeclaration .ReturnType}}{{else}}void{{end}} {{snake .Name}}({{range $i, $param := .Parameters}}{{if $i}}, {{end}}{{parameterDeclaration $param}}{{end}}) = 0;

	{{end}}

  public:
    {{.Schema.Name}}ServiceServer();
};

class {{.Schema.Name}}ServiceClient : public NanoPack::RpcClient {
  public:
    using NanoPack::RpcClient::RpcClient;

	{{- range .Schema.DeclaredFunctions}}
	std::future<{{if .ReturnType}}{{typeDeclaration .ReturnType}}{{else}}void{{end}}> {{snake .Name}}({{range $i, $param := .Parameters}}{{if $i}}, {{end}}{{parameterDeclaration $param}}{{end}});

	{{end}}
};

{{if .Namespace -}} } // namespace {{.Namespace}}{{end}}

#endif
`

const serviceImplFile = `#include <exception>
#include <string>
#include <nanopack/reader.hxx>
#include <nanopack/writer.hxx>

#include "{{.HeaderName}}"

{{if .Namespace}}{{.Namespace}}::{{end}}{{.Schema.Name}}ServiceServer::{{.Schema.Name}}ServiceServer() : NanoPack::RpcServer(), handlers() {
	handlers.reserve({{len .Schema.DeclaredFunctions}});
	{{- range .Schema.DeclaredFunctions}}
	handlers.emplace("{{.Name}}", &{{$.Schema.Name}}ServiceServer::_{{snake .Name}});
	{{- end}}
}

NanoPack::RpcServer::MethodCallResult {{if .Namespace}}{{.Namespace}}::{{end}}{{.Schema.Name}}ServiceServer::on_method_call(const std::string_view &method, uint8_t *request_data, size_t offset, NanoPack::MessageId msg_id) {
	const auto handler = handlers.find(method);
	if (handler == handlers.end()) {
		throw std::invalid_argument("Unknown method " + std::string(method) + " called on {{.Schema.Name}}Service.");
	}
	return (this->*(handler->second))(request_data, offset, msg_id);
}

{{- range .Schema.DeclaredFunctions}}
NanoPack::RpcServer::MethodCallResult {{if $.Namespace}}{{$.Namespace}}::{{end}}{{$.Schema.Name}}ServiceServer::_{{snake .Name}}(uint8_t *request_data, size_t offset, NanoPack::MessageId msg_id) {
	NanoPack::Reader reader(request_data);
	size_t ptr = offset;
	uint8_t *buf = reader.buffer;
	{{generateReadParamCode .}}
	{{if .ReturnType -}}
	{{typeDeclaration .ReturnType}} result =
		{{- if isTriviallyCopyable .ReturnType -}}
		{{snake .Name}}({{range $i, $param := .Parameters}}{{if $i}}, {{end}}{{rvalue $param}}{{end}});
		{{- else -}}
		std::move({{snake .Name}}({{range $i, $param := .Parameters}}{{if $i}}, {{end}}{{rvalue $param}}{{end}}));
		{{- end -}}
	{{- else -}}
	{{snake .Name}}({{range $i, $param := .Parameters}}{{if $i}}, {{end}}{{rvalue $param}}{{end}});
	{{end}}
	NanoPack::Writer writer({{if and .ReturnType (gt .ReturnType.ByteSize 0) }}6 + {{.ReturnType.ByteSize}}{{else}}6{{end}});
	writer.append_uint8(NanoPack::RpcMessageType::Response);
	writer.append_uint32(msg_id);
	writer.append_uint8(0);
	{{if .ReturnType}}{{generateWriteResultCode .}}{{end}}
	return {writer.into_data(), writer.size()};
}

{{end}}

{{- range .Schema.DeclaredFunctions}}
std::future<{{if .ReturnType}}{{typeDeclaration .ReturnType}}{{else}}void{{end}}> {{if $.Namespace}}{{$.Namespace}}::{{end}}{{$.Schema.Name}}ServiceClient::{{snake .Name}}({{range $i, $param := .Parameters}}{{if $i}}, {{end}}{{parameterDeclaration $param}}{{end}}) {
	NanoPack::Writer writer(9 + {{stringByteSize .Name}}{{if gt .ParametersByteSize 0}} + {{.ParametersByteSize}}{{end}});
	const auto msg_id = new_message_id();
	writer.append_uint8(NanoPack::RpcMessageType::Request);
	writer.append_uint32(msg_id);
	writer.append_uint32({{stringByteSize .Name}});
	writer.append_string_view("{{.Name}}");
	{{generateWriteParamCode .}}
	return std::async(
		[this](uint32_t msg_id, uint8_t *req_data, size_t req_size) {
			auto res_data = send_request_data_async(msg_id, req_data, req_size).get();
			NanoPack::Reader reader(res_data);
			size_t ptr = 0;
			uint8_t err_flag;
			reader.read_uint8(ptr++, err_flag);
			if (err_flag == 1) {
				throw std::runtime_error("RPC on {{$.Schema.Name}}::{{.Name}} failed.");
			}
			free(req_data);
			{{if .ReturnType}}
			uint8_t *buf = reader.buffer;
			{{generateReadResultCode .}}
			return result;
			{{- end}}
		},
		msg_id, writer.into_data(), writer.size());
}

{{end}}
`
