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
	pageData, err := GenerateHeaderFooterCmdBar(db)
	if err != nil {
		log.Fatal(err)
	}
	var template = template.Must(template.ParseFS(HTML, "header.html", "import.html", "footer.html"))
	err = template.ExecuteTemplate(&htmtBuffer, "import", pageData)
	if err != nil {
		log.Fatal(err)
	}

	return htmtBuffer.Bytes()

}

func GenerateTableCreds(db *sql.DB) []byte {
	page, err := GenerateHeaderFooterCmdBar(db)
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
	//tabletempl, err := template.ParseFiles("html/table_cred.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//tabletempl.Execute(&htmtBuffer, page)

	var template = template.Must(template.ParseFS(HTML, "header.html", "table_cred.html", "footer.html"))
	err = template.ExecuteTemplate(&htmtBuffer, "table", page)
	if err != nil {
		log.Fatal(err)
	}

	return htmtBuffer.Bytes()

}

func GenerateUpdateFormCreds(db *sql.DB, id int) ([]byte, error) {
	var html bytes.Buffer
	updateTempl, err := template.ParseFS(HTML, "header.html", "updateform_cred.html", "footer.html")
	if err != nil {
		return html.Bytes(), err
	}
	pagedata, err := GenerateHeaderFooterCmdBar(db)
	if err != nil {
		return nil, err
	}
	pagedata.CredUpdate, err = mysql.GetCred(db, id)
	if err != nil {
		return html.Bytes(), err
	}

	err = updateTempl.ExecuteTemplate(&html, "updateCred", pagedata)
	if err != nil {
		return nil, err
	}

	return html.Bytes(), nil

}

func GenerateUpdateFormHost(db *sql.DB, id int) ([]byte, error) {
	var html bytes.Buffer
	pagedata, err := GenerateHeaderFooterCmdBar(db)
	//updateTempl, err := template.ParseFiles("html/updateform_host.html")
	updateTempl, err := template.ParseFS(HTML, "updateform_host.html", "header.html", "footer.html")
	if err != nil {
		return html.Bytes(), err
	}
	pagedata.HostUpdate, err = mysql.GetHost(db, id)
	if err != nil {
		return html.Bytes(), err
	}

	err = updateTempl.ExecuteTemplate(&html, "updateHost", pagedata)
	if err != nil {
		return nil, err
	}
	return html.Bytes(), nil

}

func GenerateSettingsPage(db *sql.DB) []byte {
	var htmtBuffer bytes.Buffer
	pageData, _ := GenerateHeaderFooterCmdBar(db)
	var template = template.Must(template.ParseFS(HTML, "header.html", "setting.html", "footer.html"))
	err := template.ExecuteTemplate(&htmtBuffer, "settings", pageData)
	if err != nil {
		log.Fatal(err)
	}
	return htmtBuffer.Bytes()
}

func GenerateTableHosts(db *sql.DB) ([]byte, error) {
	page, err := GenerateHeaderFooterCmdBar(db)
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
	//tabletempl, err := template.ParseFiles("html/table_cred.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//tabletempl.Execute(&htmtBuffer, page)

	var template = template.Must(template.ParseFS(HTML, "header.html", "table_host.html", "footer.html"))
	err = template.ExecuteTemplate(&htmtBuffer, "table", page)
	if err != nil {
		log.Fatal(err)
	}

	return htmtBuffer.Bytes(), nil

}

func GenerateTableCommands(db *sql.DB) ([]byte, error) {
	// GenerateHeaderFooterCmdBar already get the commands from table so we can reuse it.
	page, err := GenerateHeaderFooterCmdBar(db)
	if err != nil {
		return nil, err
	}
	var htmtBuffer bytes.Buffer

	//
	//
	//tabletempl, err := template.ParseFiles("html/table_cred.html")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//tabletempl.Execute(&htmtBuffer, page)

	var template = template.Must(template.ParseFS(HTML, "header.html", "table_commands.html", "footer.html"))
	err = template.ExecuteTemplate(&htmtBuffer, "table", page)
	if err != nil {
		log.Fatal(err)
	}

	return htmtBuffer.Bytes(), nil

}

//Generate the command bar and command list from SQL list
func GenerateHeaderFooterCmdBar(db *sql.DB) (typelib.PageEntries, error) {
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
