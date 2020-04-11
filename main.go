package main

import (
	"github.com/gogo/protobuf/vanity/command"

	_ "github.com/bbdLe/iGameProtoPlugin/internal"
)

func main() {
	command.Write(command.Generate(command.Read()))
}
