package mysql

import (
	"github.com/ilightthings/GED/typelib"
	"math/rand"
	"testing"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func genRandomCred() typelib.CommandBar {
	randomcred := typelib.CommandBar{}
	randomcred.User = RandStringRunes(32)
	randomcred.Password = RandStringRunes(16)
	randomcred.Domain = RandStringRunes(8)
	randomcred.Command = RandStringRunes(20)
	return randomcred
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestSetCredBarEntry(t *testing.T) {
	sqliteDatabase := OpenDatabase()
	person := genRandomCred()
	err := SetCredBarEntry(sqliteDatabase, person)
	if err != nil {
		t.Error(err)
	}

	//data := genRandomCred()

}

func TestGetCommandBarEntry(t *testing.T) {
	sqliteDatabase := OpenDatabase()
	entry, err := GetCommandBarEntry(sqliteDatabase)
	if err != nil {
		t.Error(err)
	}

	if len(entry.User) != 32 {
		t.Errorf("username is not 32 chars long. Username: %s\n", entry.User)
	}
}
