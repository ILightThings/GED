package html

import (
	"bytes"
	"database/sql"
	"github.com/ilightthings/GED/mysql"
	"github.com/ilightthings/GED/typelib"
	"html/template"
	"log"
)

func GenerateImportPage(db *sql.DB) []byte {
	var htmtBuffer bytes.Buffer
	pageData, err := GenerateGeneral(db)
	if err != nil {
		log.Fatal(err)
	}
	var template = template.Must(template.ParseFiles("html/header.html", "html/import.html", "html/footer.html"))
	err = template.ExecuteTemplate(&htmtBuffer, "import", pageData)
	if err != nil {
		log.Fatal(err)
	}

	return htmtBuffer.Bytes()

}

func GenerateCredsTable(db *sql.DB) []byte {
	page, err := GenerateGeneral(db)
	if err != nil {
		log.Fatal(err)
	}
	page.CredEntries, err = mysql.GetCredTableSQLEntries(db)
	if err != nil {
		log.Fatal(err)
	}
	var htmtBuffer bytes.Buffer

	//
	//
	//tabletempl, err := template.ParseFiles("html/credtable.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//tabletempl.Execute(&htmtBuffer, page)

	var template = template.Must(template.ParseFiles("html/header.html", "html/credtable.html", "html/footer.html"))
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

func GenerateHostUpdate(db *sql.DB, id int) ([]byte, error) {
	var html bytes.Buffer
	updateTempl, err := template.ParseFiles("html/updateHost.html")
	if err != nil {
		return html.Bytes(), err
	}
	data, err := mysql.GetHost(db, id)
	if err != nil {
		return html.Bytes(), err
	}

	updateTempl.Execute(&html, data)
	return html.Bytes(), nil

}

func GenerateSettingsPage(db *sql.DB) []byte {
	var htmtBuffer bytes.Buffer
	pageData, _ := GenerateGeneral(db)
	var template = template.Must(template.ParseFiles("html/header.html", "html/setting.html", "html/footer.html"))
	err := template.ExecuteTemplate(&htmtBuffer, "settings", pageData)
	if err != nil {
		log.Fatal(err)
	}
	return htmtBuffer.Bytes()
}

func GenerateHostTable(db *sql.DB) ([]byte, error) {
	page, err := GenerateGeneral(db)
	if err != nil {
		return nil, err
	}
	page.HostEntries, err = mysql.GetHostList(db)
	if err != nil {
		return nil, err
	}
	var htmtBuffer bytes.Buffer

	//
	//
	//tabletempl, err := template.ParseFiles("html/credtable.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//tabletempl.Execute(&htmtBuffer, page)

	var template = template.Must(template.ParseFiles("html/header.html", "html/hosttable.html", "html/footer.html"))
	err = template.ExecuteTemplate(&htmtBuffer, "table", page)
	if err != nil {
		log.Fatal(err)
	}

	return htmtBuffer.Bytes(), nil

}

//Generate the command bar and command list from SQL list
func GenerateGeneral(db *sql.DB) (typelib.PageEntries, error) {
	var pageData typelib.PageEntries
	var err error
	pageData.CommanderBar, err = mysql.GetCommandBarEntry(db)
	if err != nil {
		return pageData, err
	}

	data, err := mysql.GetCommandLib(db)
	pageData.CommandList = data.ListOfCommands
	if err != nil {
		return pageData, err
	}

	return pageData, nil

}
