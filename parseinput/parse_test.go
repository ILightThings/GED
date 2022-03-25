package parseinput

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseData(t *testing.T) {

	testcmd := "bandadasdd  vecktor.facebook/narration:'aaaaaahhhhhh'@10.0.0.1 -hashes :8846F7EAEE8FB117AD06BDD830B7586C"
	credstruct := ImpacketInput(testcmd)

	expectedHash := "8846F7EAEE8FB117AD06BDD830B7586C"
	expectedDomain := "vecktor.facebook"
	exepectedUser := "narration"
	expectedPassword := "aaaaaahhhhhh"

	if credstruct.Domain != expectedDomain {
		t.Errorf("Expecting Domain %s, got %s", expectedDomain, credstruct.Domain)
	}
	if credstruct.User != exepectedUser {
		t.Errorf("Expecting Username %s, got %s", exepectedUser, credstruct.User)
	}
	if credstruct.Password != expectedPassword {
		t.Errorf("Expecting Password %s, got %s", expectedPassword, credstruct.Password)
	}
	if credstruct.Hash != expectedHash {
		t.Errorf("Expecting Hash %s, got %s", expectedHash, credstruct.Hash)
	}

}

func TestIdentifyCMEline(t *testing.T) {
	data := ``

	newdata := strings.Split(data, "\n")

	result := IdentifyCMEline(newdata)
	if len(result.CredEntries) == 0 {
		t.Error("Something should have back, but didnt.")
	} else {
		for x := range result.CredEntries {
			fmt.Printf("User: %s, Domain: %s, Pass: %s, Hash: %s \n",
				result.CredEntries[x].User,
				result.CredEntries[x].Domain,
				result.CredEntries[x].Password,
				result.CredEntries[x].Hash)
		}
	}
}
