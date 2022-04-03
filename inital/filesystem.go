package inital

import (
	_ "embed"
	"encoding/json"
	"github.com/ilightthings/GED/typelib"
)

//go:embed output_command_reference.json
var Commands []byte

//go:embed input_command_parse.json
var InputCommandJson []byte

func GenerateCommandTemplate() (typelib.CommandParseDB, error) {
	// TODO store in SQL
	var commandsDB typelib.CommandParseDB

	err := json.Unmarshal(InputCommandJson, &commandsDB)
	if err != nil {
		return commandsDB, err
	}
	return commandsDB, nil
}
