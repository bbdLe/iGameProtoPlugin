package internal

import (
	"bytes"
	"fmt"
	"github.com/bbdLe/iGameProtoPlugin/internal/util"
	"log"
	"strings"
	"text/template"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

const (
	tmpl = `
	comm.RegMessageMeta(&comm.MessageMeta{
		MsgId: int({{.Id}}),
		Type: reflect.TypeOf((*{{.Type}})(nil)).Elem(),
		Codec: codec.MustGetCodec("gogopb"),
	})
`
)

type MsgEntry struct {
	Type string
	Id string
}

var (
	pathMap = make(map[string]*descriptor.SourceCodeInfo_Location)
)

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
	// comment map
	for _, loc := range file.GetSourceCodeInfo().GetLocation() {
		key := ""

		for _, v := range loc.GetPath() {
			key = fmt.Sprintf("%s|%d", key, v)
		}

		pathMap[key] = loc
	}

	p.P("func init() {\n")
	for i, msg := range file.MessageType {
		p.genMessageCode(i, msg)
	}
	p.P("\n}")
}

func (p *gameproto) genMessageCode(index int, msg *descriptor.DescriptorProto) {
	key := fmt.Sprintf("|4|%d", index)

	comment := ""
	if v, ok := pathMap[key]; ok {
		comment = v.GetLeadingComments()
		comment = strings.TrimSpace(comment)
		comment = strings.TrimSpace(comment)
	}

	if len(comment) == 0 {
		return
	}
	comment = strings.TrimPrefix(comment, "[")
	comment = strings.TrimSuffix(comment, "]")
	m := util.Comment2Map(comment)
	if len(m) == 0  {
		log.Printf("msg(%s) without id", msg.GetName())
		return
	}
	id, ok := m["id"]
	if !ok {
		log.Printf("msg(%s) without id", msg.GetName())
		return
	}
	id = fmt.Sprintf("ProtoID_%s", id)

	msgEntry := MsgEntry{
		Type: msg.GetName(),
		Id: id,
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