package transform

import (
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/ilightthings/GED/typelib"
	"strings"
)

//TODO, Change the passed parameter to typelib.CommandBar. No more SQL calls from this sections please.
func CMEOUT(db *sql.DB, id int, target string) string {
	search := fmt.Sprintf("SELECT username,domain,password,hash FROM creds WHERE IDCRED IS %d", id)
	rows, err := db.Query(search)
	defer rows.Close()
	v := validator.New()

	if err != nil {
		return err.Error()
	}

	for rows.Next() {
		cmd := fmt.Sprintf("crackmapexec smb %s ", target)
		var user string
		var domain string
		var password string
		var hash string
		rows.Scan(&user, &domain, &password, &hash)
		if user != "" {
			cmd = cmd + fmt.Sprintf("-u '%s' ", user)
		}
		if domain != "" {
			err := v.Var(domain, "ipv4")
			if err != nil {
				cmd = cmd + fmt.Sprintf("-d '%s' ", domain)
			}
		}
		if password != "" {
			cmd = cmd + fmt.Sprintf("-p '%s' ", password)
		} else if hash != "" {
			cmd = cmd + fmt.Sprintf("-H '%s' ", hash)
		}
		return cmd

	}
	return "Well this should not have happend"
}

func IMPOUT(db *sql.DB, id int, target string) string {
	search := fmt.Sprintf("SELECT username,domain,password,hash FROM creds WHERE IDCRED IS %d", id)
	rows, err := db.Query(search)
	defer rows.Close()
	if err != nil {
		return err.Error()
	}

	for rows.Next() {
		var cmd string
		var user string
		var domain string
		var password string
		var hash string
		rows.Scan(&user, &domain, &password, &hash)
		v := validator.New()

		if domain != "" {
			err := v.Var(domain, "ipv4")
			if err != nil {

				cmd = cmd + domain + "/"
			} else {
				cmd = cmd + "./"
			}
		}

		if user != "" {
			cmd = cmd + user
		}

		if password != "" {
			cmd = cmd + ":'" + password + "'"
		} else if hash != "" {
			cmd = cmd + " -hashes :" + hash
		}
		return cmd

	}
	return "How did you get here....."

}

func CrackMapExecOut(bar typelib.CommandBar) (string, error) {
	err := bar.Prepare()
	if err != nil {
		return "", err
	}
	var userout string
	var domainout string
	var authout string
	var hostout string

	hostout = bar.Host

	userout = fmt.Sprintf("-u %s", bar.User)

	if bar.Password != "" {
		authout = fmt.Sprintf("-p '%s'", bar.Password)
	} else {
		authout = fmt.Sprintf("-H %s", bar.Hash)
	}

	if bar.Domain != "" {
		domainout = fmt.Sprintf("-d %s", bar.Domain)
	}

	final := fmt.Sprintf("crackmapexec smb %s %s %s %s", hostout, userout, authout, domainout)
	return final, nil
}

func CustomCommand(commandTemplate string, bar typelib.CommandBar) string {
	commandTemplate = strings.ReplaceAll(commandTemplate, "##USER##", bar.User)
	commandTemplate = strings.ReplaceAll(commandTemplate, "##PASSWORD##", bar.Password)
	commandTemplate = strings.ReplaceAll(commandTemplate, "##HASH##", bar.Hash)
	commandTemplate = strings.ReplaceAll(commandTemplate, "##DOMAIN##", bar.Domain)
	commandTemplate = strings.ReplaceAll(commandTemplate, "##HOST##", bar.Host)
	return commandTemplate
}
