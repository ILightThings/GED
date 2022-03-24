package typelib

import (
	"errors"
	"fmt"
)

type PageEntries struct {
	CredEntries []CredEntry
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
