package internal

import (
	"bytes"
	"log"
	"text/template"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

const (
	tmpl = `
	comm.RegMessageMeta(&comm.MessageMeta{
		MsgId: int(util.StringHash("proto.{{.MsgType}}")),
		Type: reflect.TypeOf((*{{.MsgType}})(nil)).Elem(),
		Codec: codec.MustGetCodec("gogopb"),
	})
`
)

type MsgEntry struct {
	MsgType string
}

type gameproto struct{ *generator.Generator }

func (p *gameproto) Name() string                { return "gameproto" }
func (p *gameproto) Init(g *generator.Generator) { p.Generator = g }

func (p *gameproto) GenerateImports(file *generator.FileDescriptor) {
	if len(file.MessageType) == 0 {
		return
	}

	p.P("\"github.com/bbdLe/iGame/comm\"")
	p.P("\"github.com/bbdLe/iGame/comm/codec\"")
	p.P("\"github.com/bbdLe/iGame/comm/util\"")
	p.P("_ \"github.com/bbdLe/iGame/comm/codec/gogopb\"")
	p.P("\"reflect\"")
}

func (p *gameproto) Generate(file *generator.FileDescriptor) {
	p.P("func init() {\n")
	for _, msg := range file.MessageType {
		p.genMessageCode(msg)
	}
	p.P("\n}")
}

func (p *gameproto) genMessageCode(msg *descriptor.DescriptorProto) {
	msgEntry := MsgEntry{
		MsgType: *msg.Name,
	}

	var buf bytes.Buffer
	t := template.Must(template.New("").Parse(tmpl))
	err := t.Execute(&buf, msgEntry)
	if err != nil {
		log.Fatal(err)
	}

	p.P(buf.String())
}

func init() {
	generator.RegisterPlugin(new(gameproto))
}