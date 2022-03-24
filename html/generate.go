package html

import (
	"bytes"
	"database/sql"
	"github.com/ilightthings/GED/mysql"
	"github.com/ilightthings/GED/typelib"
	"html/template"
	"log"
)

func GenerateImportPage() []byte {
	var htmtBuffer bytes.Buffer
	var template = template.Must(template.ParseFiles("html/header.html", "html/import.html", "html/footer.html"))
	err := template.ExecuteTemplate(&htmtBuffer, "import", nil)
	if err != nil {
		log.Fatal(err)
	}

	return htmtBuffer.Bytes()

}

func GenerateCredsTable(db *sql.DB) []byte {
	var page typelib.PageEntries
	var htmtBuffer bytes.Buffer

	rows, err := db.Query("SELECT idCred,username,domain,password,hash FROM creds")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var cred typelib.CredEntry
		var id int
		var us string
		var do string
		var pa string
		var ha string
		rows.Scan(&id, &us, &do, &pa, &ha)

		cred.ID = id
		cred.User = us
		cred.Domain = do
		cred.Password = pa
		cred.Hash = ha

		page.CredEntries = append(page.CredEntries, cred)

	}

	//
	//
	//tabletempl, err := template.ParseFiles("html/table.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//tabletempl.Execute(&htmtBuffer, page)

	var template = template.Must(template.ParseFiles("html/header.html", "html/table.html", "html/footer.html"))
	err = template.ExecuteTemplate(&htmtBuffer, "table", page)
	if err != nil {
		log.Fatal(err)
	}

	return htmtBuffer.Bytes()

}

func GenerateCredUpdate(db *sql.DB, id int) ([]byte, error) {
	var html bytes.Buffer
	updateTempl, err := template.ParseFiles("html/updateCred.html")
	if err != nil {
		return html.Bytes(), err
	}
	data, err := mysql.GetCred(db, id)
	if err != nil {
		return html.Bytes(), err
	}

	updateTempl.Execute(&html, data)
	return html.Bytes(), nil

}

func GenerateSettingsPage() []byte {
	var htmtBuffer bytes.Buffer
	var template = template.Must(template.ParseFiles("html/header.html", "html/setting.html", "html/footer.html"))
	err := template.ExecuteTemplate(&htmtBuffer, "settings", nil)
	if err != nil {
		log.Fatal(err)
	}
	return htmtBuffer.Bytes()
}
