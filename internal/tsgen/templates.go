package tsgen

import "nanoc/internal/npschema"

type messageClassTemplateInfo struct {
	Schema          *npschema.Message
	ExternalImports []string

	ConstructorParameters []string
	ConstructorArgs       []string
	SuperConstructorArgs  []string

	FieldReadCodeFragments  []string
	FieldWriteCodeFragments []string
}

type messageClassFactoryTemplateInfo struct {
	Schema             *npschema.Message
	MessageClassImport string
	MessageImports     []string
}

type messageFactoryTemplateInfo struct {
	Schemas        []*npschema.Message
	MessageImports []string
}

type enumTemplateInfo struct {
	Schema             *npschema.Enum
	MemberDeclarations []string
}

const (
	extImport = ".np.js"

	extTsFile = ".np.ts"
)

const (
	templateNameMessageClass = "TsMessageClass"

	templateNameMessageClassFactory = "TsMessageClassFactory"

	templateNameMessageFactory = "TsMessageFactory"

	templateNameEnum = "TsEnum"
)

const fileNameMessageFactoryFile = "message-factory"

const messageClassTemplate = `// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import { NanoBufReader, NanoBufWriter{{if not .Schema.HasParentMessage}}, type NanoPackMessage{{end}} } from "nanopack";

{{range .ExternalImports}}
{{.}}
{{- end}}

class {{.Schema.Name}} {{if .Schema.HasParentMessage}}extends {{.Schema.ParentMessage.Name}}{{else}}implements NanoPackMessage{{end}} {
  public static TYPE_ID = {{.Schema.TypeID}};

  public {{if .Schema.HasParentMessage}}override {{end}}readonly typeId: number = {{.Schema.TypeID}};

  public {{if .Schema.HasParentMessage}}override {{end}}readonly headerSize: number = {{.Schema.HeaderSize}};

  {{if .Schema.HasParentMessage -}}
  constructor({{join .ConstructorParameters ", "}}) {
    super({{join .SuperConstructorArgs ", "}})
  }
  {{else}}
  constructor({{join .ConstructorParameters ", "}}) {}
  {{- end}}

  public static fromBytes(bytes: Uint8Array): { bytesRead: number, result: {{.Schema.Name}} } | null {
    const reader = new NanoBufReader(bytes);
    return {{.Schema.Name}}.fromReader(reader);
  }

  public static fromReader(reader: NanoBufReader, offset = 0): { bytesRead: number, result: {{.Schema.Name}} } | null {
    let ptr = offset + {{.Schema.HeaderSize}};

    {{range .FieldReadCodeFragments}}
    {{.}}

    {{end}}

    return { bytesRead: ptr - offset, result: new {{.Schema.Name}}({{join .ConstructorArgs ", "}}) };
  }

  {{if .Schema.HasParentMessage}}override {{end}}public writeTo(writer: NanoBufWriter, offset = 0): number {
    let bytesWritten = {{.Schema.HeaderSize}};

    writer.writeTypeId({{.Schema.TypeID}}, offset);

    {{range .FieldWriteCodeFragments}}
    {{.}}

    {{end}}

    return bytesWritten;
  }

  {{if .Schema.HasParentMessage}}override {{end}}public bytes(): Uint8Array {
    const writer = new NanoBufWriter({{.Schema.HeaderSize}});
    this.writeTo(writer)
    return writer.bytes;
  }
}

export { {{.Schema.Name}} };
`

const messageClassFactoryTemplate = `// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import { NanoBufReader } from "nanopack";

import { {{.Schema.Name}} } from "{{.MessageClassImport}}";
{{range .MessageImports -}}
{{.}}
{{- end}}

function make{{.Schema.Name}}(bytes: Uint8Array) {
  const reader = new NanoBufReader(bytes);
  switch (reader.readTypeId()) {
  case {{.Schema.TypeID}}: return {{.Schema.Name}}.fromReader(reader);
  {{- range .Schema.ChildMessages}}
  case {{.TypeID}}: return {{.Name}}.fromReader(reader);
  {{- end}}
  default: return null;
  }
}

export { make{{.Schema.Name}} } ;
`

const messageFactoryTemplate = `// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

import { NanoBufReader, type NanoPackMessage } from "nanopack";
{{range .MessageImports}}
{{.}}
{{- end}}

function makeNanoPackMessage(bytes: Uint8Array): { bytesRead: number, result: NanoPackMessage } | null {
  const reader = new NanoBufReader(bytes);
  switch (reader.readTypeId()) {
  {{range .Schemas}}
  case {{.TypeID}}: return {{.Name}}.fromReader(reader);
  {{- end}}
  default: return null;
  }
}

export { makeNanoPackMessage }
`

const enumTemplate = `// AUTOMATICALLY GENERATED BY NANOPACK. DO NOT MODIFY BY HAND.

const {{.Schema.Name}} = {
  {{range .MemberDeclarations}}
  {{.}},
  {{- end}}
} as const;

type T{{.Schema.Name}} = typeof {{.Schema.Name}}[keyof typeof {{.Schema.Name}}];

export { {{.Schema.Name}} };
export type { T{{.Schema.Name}} };
`
