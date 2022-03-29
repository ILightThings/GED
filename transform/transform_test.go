package transform

import (
	"github.com/ilightthings/GED/typelib"
	"testing"
)

func TestCrackMapExecOut(t *testing.T) {
	var sampleuser typelib.CommandBar
	sampleuser.Domain = "master.com"
	sampleuser.Password = "password123"
	sampleuser.User = "masterchief"
	sampleuser.Host = "192.168.1.1"

	expectedResult := "crackmapexec smb 192.168.1.1 -u masterchief -p password123 -d master.com"
	actualResult, err := CrackMapExecOut(sampleuser)
	if err != nil {
		t.Errorf(err.Error())
	}

	if actualResult != expectedResult {
		t.Errorf("Expected: %s\nActual: %s\n", expectedResult, actualResult)
	}
}

func TestCustomCommand(t *testing.T) {
	var sampleuser typelib.CommandBar
	sampleuser.Domain = "master.com"
	sampleuser.Password = "password123"
	sampleuser.User = "masterchief"
	sampleuser.Host = "192.168.1.1"
	customcommand := "smbmap -u ##USER## -p ##PASSWORD## -d ##DOMAIN## -H ##HOST##"
	expectOutPut := "smbmap -u masterchief -p password123 -d master.com -H 192.168.1.1"
	result := CustomCommand(customcommand, sampleuser)
	if result != expectOutPut {
		t.Errorf("Expected %s\nGot: %s\n", expectOutPut, result)
	}

}
