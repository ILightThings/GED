package parseinput

import (
	"testing"
)

func TestExtractCommandData(t *testing.T) {
	testuser := "mike"
	testpassword := ""
	testdomain := ""
	testhost := ""

	commandToTest := ""
	cred, host, err := ExtractCommandData(commandToTest)
	if err != nil {
		t.Error(err)
	}

	if cred.User != testuser {
		t.Errorf("Expected: %s, Got: %s\n", testuser, cred.User)
	}
	if cred.Password != testpassword {
		t.Errorf("Expected: %s, Got: %s\n", testpassword, cred.Password)
	}
	if cred.User != testuser {
		t.Errorf("Expected: %s, Got: %s\n", testdomain, cred.Domain)
	}
	if host.IP != testhost {
		t.Errorf("Expected: %s, Got: %s\n", testhost, host.IP)
	}
}
