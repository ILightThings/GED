package inital

import (
	_ "embed"
	"encoding/json"
	"github.com/ilightthings/GED/typelib"
)

//go:embed output_command_reference.json
var CommandOutTemplateBytes []byte

//go:embed input_command_parse.json
var InputeParseBytes []byte

func GenerateCommandTemplate() (typelib.CommandParseDB, error) {
	var commandsDB typelib.CommandParseDB

	err := json.Unmarshal(InputeParseBytes, &commandsDB)
	if err != nil {
		return commandsDB, err
	}
	return commandsDB, nil
}

//TODO add catagroyies to commands for searching features later.
