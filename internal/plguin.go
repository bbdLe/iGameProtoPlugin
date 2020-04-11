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
	// comment map
	pathMap := make(map[string]*descriptor.SourceCodeInfo_Location)
	for _, loc := range file.GetSourceCodeInfo().GetLocation() {
		key := ""

		for _, v := range loc.GetPath() {
			key = fmt.Sprintf("%s|%d", key, v)
		}

		pathMap[key] = loc
	}

	// message map
	msgMap := make(map[string]*descriptor.DescriptorProto)
	for _, msg := range file.GetMessageType() {
		msgMap[msg.GetName()] = msg
	}

	// CS宏解析出信息
	for i, enumEntry := range file.GetEnumType() {
		if *enumEntry.Name != "CsMsgId" {
			continue
		}

		for j, v := range enumEntry.GetValue() {
			key := fmt.Sprintf("|5|%d|2|%d", i, j)
			local, ok := pathMap[key]
			if !ok {
				continue
			}

			comment := local.GetTrailingComments()
			comment = strings.TrimLeft(comment, "\r\n ")
			comment = strings.TrimRight(comment, "\r\n ")
			if len(comment) == 0 {
				continue
			}

			kvMap := util.Comment2Map(comment)
			if len(kvMap) == 0 {
				continue
			}

			// 获取msg
			msgName, ok := kvMap["message"]
			if !ok {
				log.Printf("msg is empty, macro is %s, kvMap is %v", v.GetName(), kvMap)
				continue
			}

			msgMeta, ok := msgMap[msgName]
			if !ok {
				log.Println("can't find msg : ", msgName)
			}

			log.Printf("macro is %s, msg meta is : %v", msgName, msgMeta.GetName())
		}
	}

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