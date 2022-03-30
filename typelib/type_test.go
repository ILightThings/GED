package typelib

import (
	"os"
	"testing"
)

func TestCommandLibrary_ImportFromJson(t *testing.T) {
	var testlib CommandLibrary
	jsonfile, err := os.ReadFile("inital/commands.json")
	if err != nil || jsonfile == nil {
		t.Error(err.Error())
	}
	err = testlib.ImportFromJson(jsonfile)
	if err != nil {
		t.Error(err.Error())
	}
	if len(testlib.ListOfCommands) != 7 {
		t.Errorf("Execpect 7,\nGot %d", len(testlib.ListOfCommands))
	}

}
