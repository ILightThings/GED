package mydata

// STOP USING THIS YOU FUCK
import (
	"encoding/json"
	"errors"
	"fmt"
)

type UserCreds struct {
	Username         string
	Password         string
	Hash             string
	Domain           string
	CommandReference string
	CommandPattern   string
	UsedAgainst      []string
	Parsed           bool
}

type Hosts struct {
	IP            string
	HostName      string
	Domain        string
	ImportHistory []string
}

func (u *UserCreds) StringCreds() string {
	return fmt.Sprintf("User: \"%s\", Domain: \"%s\", Password: \"%s\", Hash: \"%s\", Command Pattern: \"%s\"", u.Username, u.Domain, u.Password, u.Hash, u.CommandPattern)
}

//func (u *UserCreds) VerifyValid() bool {}

func (u *UserCreds) toJSON() ([]byte, error) {
	b, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func (u *UserCreds) Verify() error {
	value := 0
	if u.Username != "" {
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
