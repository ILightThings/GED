package parseinput

import (
	"fmt"
	"github.com/ilightthings/GED/typelib"
	"regexp"
	"strings"
)

func ParseData(command string) typelib.CredEntry {
	commandParts := strings.Split(command, " ")
	var newcred typelib.CredEntry
	newcred.CommandReference = command

	//Parse impacket command
	// "impacket-wmiexec  vecktor.facebook/narration:'aaaaaahhhhhh'@10.0.0.1 -hashes :asdasdasdasdas"
	if strings.Contains(commandParts[0], "impacket") {
		newcred = ImpacketInput(command)
		return newcred
	}

	//Parse CrackMapExec command
	if strings.Contains(commandParts[0], "crackmapexec") || strings.Contains(commandParts[0], "cme") {
		newcred = CrackMapExecInput(command)
		return newcred
	}

	if commandParts[0] == "SMB" {
		newcred = CrackmapExecOutput(command)
		return newcred

	}

	return typelib.CredEntry{}
}

func CrackMapExecInput(command string) typelib.CredEntry {
	commandParts := strings.Split(command, " ")
	var newcred typelib.CredEntry
	newcred.CommandPattern = "Crackmapexec Input"
	for x := range commandParts {
		if commandParts[x] == "-u" {
			newcred.User = commandParts[x+1]
		}
		if commandParts[x] == "-p" {
			newcred.Password = commandParts[x+1]
		}
		if commandParts[x] == "-H" {
			newcred.Hash = commandParts[x+1]
		}
		if commandParts[x] == "-d" {
			newcred.Domain = commandParts[x+1]
		}
	}
	return newcred
}

func CrackmapExecOutput(command string) typelib.CredEntry {
	commandParts := strings.Split(command, " ")
	var newcred typelib.CredEntry
	whitepaceFix1 := whitespacefix(command)
	commandParts = strings.Split(whitepaceFix1, " ")

	// Parse CrackMapExec output

	for x := range commandParts {
		// CME Goodusername and password check
		if strings.Contains(commandParts[x], ":") && strings.Contains(commandParts[x], "\\") {
			newcred.CommandPattern = "Crackmapexec Account Response"
			userpassdomain := commandParts[x]

			// place.local\username:password\withslash
			newcred.Domain = strings.Split(userpassdomain, "\\")[0]
			// place.local     username:password      withslash
			userpass := strings.Join(strings.Split(userpassdomain, "\\")[1:], "\\")

			//username:password\withslash
			newcred.User = strings.Split(userpass, ":")[0]
			hashorpass := strings.Join(strings.Split(userpass, ":")[1:], "")

			ntlmregex := regexp.MustCompile("[A-Fa-f0-9]{32}")
			if ntlmregex.MatchString(hashorpass) {
				newcred.Hash = hashorpass
			} else {
				newcred.Password = hashorpass
			}

			//  DESKTOP-4JIB526\ilightthings:8846F7EAEE8FB117AD06BDD830B7586C
		}

		// SMB         10.0.0.1       445    PRDDC03   XXXXXX-J012$:12853:aad3b435b51404eeaad3b435b51404ee:aaaab24fc345016718dad3a719012061:::
		if strings.Count(commandParts[x], ":") == 6 {
			newcred.CommandPattern = "Crackmapexec SAM Dump"
			samHash := strings.Split(commandParts[x], ":")
			newcred.User = samHash[0]
			newcred.Hash = samHash[3]
			return newcred
		}

	}

	return newcred

}

func ImpacketInput(command string) typelib.CredEntry {
	commandParts := strings.Split(command, " ")
	var newcred typelib.CredEntry
	newcred.CommandPattern = "Impacket"
	for x := range commandParts {
		if commandParts[x] == "-hashes" {
			hashinput := strings.Split(commandParts[x+1], ":")
			if len(hashinput) != 2 {
				newcred.Hash = hashinput[0]
			} else {
				newcred.Hash = hashinput[1]
			}

		}
		if strings.Contains(commandParts[x], "@") && strings.Contains(commandParts[x], "/") {
			// maybe this is the password part
			parts := strings.Split(commandParts[x], "/")
			newcred.Domain = parts[0]
			if strings.Contains(parts[1], ":") {
				usernamepassword := strings.Split(parts[1], ":")
				newcred.User = usernamepassword[0]

				// Find password string that have escape sequence with ' character
				if strings.Count(usernamepassword[1], "'") == 2 {
					escapedStrings := strings.Split(usernamepassword[1], "'")
					newcred.Password = escapedStrings[1]
				} else {
					newcred.Password = strings.Split(usernamepassword[1], "@")[0]
				}

			} else {
				newcred.User = strings.Split(parts[1], "@")[0]
			}

		}

	}
	return newcred
}
func whitespacefix(input string) string {
	regexWhiteSpace := regexp.MustCompile(`\s{2,}`)
	whitepaceFix := regexWhiteSpace.ReplaceAllString(input, " ")
	return whitepaceFix
}

func IdentifyBlob(input string) {
	data := strings.Split(input, "\n")

	if strings.Contains(data[0], "cme") || strings.Contains(data[0], "crackmapexec") {
		//process
	}
}

func IdentifyCMEline(data []string) typelib.PageEntries {
	var creds typelib.PageEntries
	RegexLibary := map[string]string{}
	RegexLibary["SecretsDump/NTDS Dump/SamDump"] = `([A-zÀ-ú0-9.]{1,256})(?::[[:digit:]]{1,5}:)(?:[A-Fa-f0-9]{32})(?::)([A-Fa-f0-9]{32})(?::{3})`
	RegexLibary["Lsassy Hash"] = `(?:\s)([A-zÀ-ú0-9.]{1,256}\\[A-zÀ-ú0-9]{1,256})(?:\s)([A-Fa-f0-9]{32})`
	RegexLibary["Lsassy Password"] = `(?:\s)([A-zÀ-ú0-9.]{1,256}\\[A-zÀ-ú0-9]{1,256})(?:\s)(\S{1,256})`
	RegexLibary["Crackmapexec Hash Input"] = `(?:(\[\S\])\s)([A-zÀ-ú0-9.]{1,256}\\[A-zÀ-ú0-9]{1,256})(?::)([A-Fa-f0-9]{32})`
	RegexLibary["Crackmapexec Password Input"] = `(?:\[\S\]\s)([A-zÀ-ú0-9.]{1,256}\\[A-zÀ-ú0-9]{1,256})(?::)(\S{1,256})`
	for x := range data {
		for _, y := range RegexLibary {
			var newcred typelib.CredEntry
			re := regexp.MustCompile(y)
			result_slice := re.FindStringSubmatch(data[x])
			if len(result_slice) < 2 {
				continue
			}

			if strings.Contains(result_slice[1], "\\") {
				UserDomain := strings.Split(result_slice[1], "\\")
				newcred.Domain = UserDomain[0]
				newcred.User = UserDomain[1]
			} else {
				ipRegex := regexp.MustCompile(`(?:\s)([0-2]{0,1}[0-9]{0,1}[0-9]{1}\.[0-2]{0,1}[0-9]{0,1}[0-9]{1}\.[0-2]{0,1}[0-9]{0,1}[0-9]{1}\.[0-2]{0,1}[0-9]{0,1}[0-9]{1})`)
				ipDomain := ipRegex.FindStringSubmatch(data[x])
				if len(ipDomain) != 0 {
					//Right Most IP Address, typicall the crackmapexec target
					newcred.Domain = ipDomain[len(ipDomain)-1]
				}
				newcred.User = result_slice[1]
			}

			hashReg := regexp.MustCompile(`[A-Fa-f0-9]{32}`)
			if hashReg.FindString(result_slice[2]) != "" {
				newcred.Hash = result_slice[2]
			} else {
				newcred.Password = result_slice[2]
			}
			creds.CredEntries = append(creds.CredEntries, newcred)
		}
	}

	return creds
	//case "lsassy":
	//	lsassyParse(data)
	//	//REGEX NTML ([A-Za-z0-9\\.]{3,}\s[A-Za-z0-9]{32}\n)
	//	//REGEX Maybe All ([A-Za-z0-9.]{3,256}\\[A-Za-z0-9.]{3,256}\s[\S]{3,256}\n)
	// //GOLANG Password/Hash of Domain (?:\s)([A-Za-z0-9]{1,256}\\[A-Za-z0-9]{1,256}\s[A-Za-z0-9]{1,256})(?:\n)
	// //GOLANG Hash Only (?:\s)([A-Za-z0-9]{1,256}\\[A-Za-z0-9]{1,256}\s[A-Fa-f0-9]{32})(?:\n)

	//case "--sam":
	//	// REGEX ([A-Za-z0-9]{1,256}:[0-9]{3,5}:[A-Fa-f0-9]{32}:[A-Fa-f0-9]{32}:::)
	//	// REGEX (?<=\s)([A-Za-z0-9]{1,256}:[0-9]{3,5}:[A-Fa-f0-9]{32}:[A-Fa-f0-9]{32}:::) - After a space,
	// GOLANG REGEX (?:\s)([A-Za-z0-9]{1,256}\s[A-Fa-f0-9]{32})(?:\n) --Local Accounts only. 32 character hash Will not get anything with a slash before it.
	//
	//// TODO --sam output parse
	//
	//case "--lsa":
	//	// TODO --lsa output parse
	//
	//case "--ntds":
	//	// REGEX ([A-Za-z0-9\.]{1,256}\\[A-Za-z0-9]{1,256}:[0-9]{3,5}:[A-Fa-f0-9]{32}:[A-Fa-f0-9]{32}:::\n)
	//	// TODO --ntds output parse
	//
	//case "nanodump":
	//	// TODO nanodump output parse
	//
	//case "procdump":
	//	// TODO procdump output parse
	//
	//}

	//Secrets Dump, Golang, Capture Groups (Name,LM,NT)
	//evil.corp\mike:1189:aad3b435b51404eeaad3b435b51404ee:70896e37c98a78a9adb86932aa64a2bf:::
	//([A-zÀ-ú0-9.]{1,256})(?::[[:digit:]]{1,5}:)([A-Fa-f0-9]{32})(?::)([A-Fa-f0-9]{32})(?::{3})

	// domain.com\user hash (output dominan/user, hash) [Supports accented chars] [Space character needed at start]
	// light.com\admin 89776771f6a491847a848063f042960b
	// (?:\s)([A-zÀ-ú0-9.]{1,256}\\[A-zÀ-ú0-9]{1,256})(?:\s)([A-Fa-f0-9]{32})

	//domain.com\user hash/password (lsassy out) (output dominan/user, hash/password) [Supports accented chars and spec chars but not spaces] [Space character needed at start]
	// light.com\admin PASSWORD123
	//(?:\s)([A-zÀ-ú0-9.]{1,256}\\[A-zÀ-ú0-9]{1,256})(?:\s)(\S{1,256})

	//CME Tested Username and HASH (domain/user, HASH) [ Characters [+] or [-] or [*] is needed at start ]
	// [+] domain.com\vumetric:89776771f6a491847a848063f042960b
	// (?:(\[\S\])\s)([A-zÀ-ú0-9.]{1,256}\\[A-zÀ-ú0-9]{1,256})(?::)([A-Fa-f0-9]{32})

	//CME Tested Username and Password (domain/user, password) [ Characters [+] or [-] or [*] is needed at start ]
	// [+] domain.com\vumetric:djg&GZW&X8PAi8gk
	// (?:(\[\S\])\s)([A-zÀ-ú0-9.]{1,256}\\[A-zÀ-ú0-9]{1,256})(?::)(\S{1,256})

}

func lsassyParse(data []string) {
	fmt.Println()
}
