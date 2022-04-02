package inital

import _ "embed"

//go:embed output_command_reference.json
var Commands []byte

//go:embed input_command_parse.json
var InputCommandJson []byte
