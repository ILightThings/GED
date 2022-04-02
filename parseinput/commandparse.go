package parseinput

import (
	"encoding/json"
	"github.com/ilightthings/GED/inital"
	"github.com/ilightthings/GED/typelib"
	"regexp"
	"strings"
)

type CommandTemplate []struct {
	CommandName  string     `json:"CommandName"`
	Alias        []string   `json:"Alias"`
	ParseType    string     `json:"ParseType"`
	CommandMatch [][]string `json:"CommandMatch"`
}

func GenerateCommandTemplate() (CommandTemplate, error) {
	// TODO store in SQL
	var commandsDB CommandTemplate

	err := json.Unmarshal(inital.InputCommandJson, &commandsDB)
	if err != nil {
		return nil, err
	}
	return commandsDB, nil
}

func ExtractCommandData(commandString string) (typelib.CredEntry, typelib.HostEntry, error) {
	var credentry typelib.CredEntry
	var hostentry typelib.HostEntry
	commandsDB, err := GenerateCommandTemplate()
	if err != nil {
		return credentry, hostentry, err
	}

	for _, y := range commandsDB {
		for z := range y.Alias {
			credentry.CommandPattern = y.CommandName
			if strings.Contains(commandString, y.Alias[z]) {
				if y.ParseType == "args" {

					commandsplit := strings.Split(commandString, " ")
					for x := range y.CommandMatch {
						for xx := range commandsplit {
							if commandsplit[xx] == y.CommandMatch[x][0] {
								switch y.CommandMatch[x][1] {
								case "host":
									hostentry.IP = commandsplit[xx+1]
								case "user":
									credentry.User = commandsplit[xx+1]
								case "domain":
									credentry.Domain = commandsplit[xx+1]
								case "password":
									credentry.Password = commandsplit[xx+1]
								case "hash":
									credentry.Hash = commandsplit[xx+1]

								}
							}
						}

					}
				} else if y.ParseType == "regex" {
					for x := range y.CommandMatch {
						regexString := regexp.MustCompile(y.CommandMatch[x][0])
						regexarray := regexString.FindStringSubmatch(commandString)
						if len(regexarray) < 2 {
							continue
						} else {
							switch y.CommandMatch[x][1] {
							case "host":
								hostentry.IP = regexarray[1]
							case "user":
								credentry.User = regexarray[1]
							case "domain":
								credentry.Domain = regexarray[1]
							case "password":
								credentry.Password = regexarray[1]
							case "hash":
								credentry.Hash = regexarray[1]

							}
						}
					}
				}
			}
		}
	}

	return credentry, hostentry, nil
}
