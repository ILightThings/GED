package html

import (
	"github.com/ilightthings/GED/mysql"
	"testing"
)

func TestGeneratePage(t *testing.T) {
	db := mysql.OpenDatabase()
	GeneratePage(db)
}
