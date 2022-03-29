package typelib

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type PageEntries struct {
	CredEntries  []CredEntry
	CommanderBar CommandBar
	CommandList  []CommandBuild
}

type CommandBar struct {
	User     string
	Domain   string
	Password string
	Hash     string
	Host     string
	Command  string
}

type CommandLibrary struct {
	ListOfCommands []CommandBuild
}
type CommandBuild struct {
	Command string
	Example string
	Display string
}

type CredEntry struct {
	User     string `json:"user" form:"user"`
	Domain   string `json:"domain" form:"domain"`
	Password string `json:"password" form:"password"`
	Hash     string `json:"hash" form:"hash"`
	ID       int    `json:"ID_int" form:"ID_int"`
	IDString string `json:"ID" form:"ID"`

	CommandReference string
	CommandPattern   string
	UsedAgainst      []string
	Parsed           bool
}

func (u *CredEntry) StringCreds() string {
	return fmt.Sprintf("User: \"%s\", Domain: \"%s\", Password: \"%s\", Hash: \"%s\", Command Pattern: \"%s\"", u.User, u.Domain, u.Password, u.Hash, u.CommandPattern)
}

func (u *CredEntry) Verify() error {
	value := 0
	if u.User != "" {
		value = value + 2
	}
	if u.Domain != "" {
		value = value + 1
	}
	if u.Password != "" || u.Hash != "" {
		value = value + 1
	}

	if value >= 1 {
		return nil
	} else {
		return errors.New("Empty Entry")
	}
}

func (c *CommandBar) Prepare() error {
	if c.Host == "" {
		c.Host = "127.0.0.1"
	}

	if c.Password == "" && c.Hash == "" {
		return errors.New("no password nor hash passed")
	}

	if c.Password != "" {
		c.Hash = ""
	}

	if c.User == "" {
		return errors.New("no user passed")
	}

	return nil
}

//TODO build a CommandBuild HTML page and CommandBuild Builder similar to user update
func (customCommand *CommandBuild) BuildCommand(bar CommandBar) string {
	start := customCommand.Command
	start = strings.ReplaceAll(start, "##USER##", bar.User)
	start = strings.ReplaceAll(start, "##PASSWORD##", bar.Password)
	start = strings.ReplaceAll(start, "##HASH##", bar.Hash)
	start = strings.ReplaceAll(start, "##DOMAIN##", bar.Domain)
	start = strings.ReplaceAll(start, "##HOST##", bar.Host)
	return start
}

func (comlib *CommandLibrary) ImportFromJson(cmdJson []byte) error {
	var commandLib []CommandBuild
	json.Unmarshal([]byte(cmdJson), &commandLib)
	comlib.ListOfCommands = commandLib
	return nil

}
