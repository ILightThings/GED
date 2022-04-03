package parseinput

import (
	"database/sql"
	"github.com/ilightthings/GED/mysql"
	"github.com/ilightthings/GED/typelib"
	"regexp"
	"strings"
)

// Parses the command input string, using the command_parse table in the database
func ExtractCommandData(commandString string, db *sql.DB) (typelib.CredEntry, typelib.HostEntry, error) {
	var credentry typelib.CredEntry
	var hostentry typelib.HostEntry
	commandsDB, err := mysql.RetreiveParseTable(db)
	if err != nil {
		return credentry, hostentry, err
	}

	for _, y := range commandsDB.Array {
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
