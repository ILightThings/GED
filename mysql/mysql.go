package mysql

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ilightthings/GED/inital"
	"github.com/ilightthings/GED/typelib"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

const DELETE_HOST = `DELETE FROM hosts WHERE idHost = (?)`
const DELETE_CRED = `DELETE FROM creds WHERE idCred = (?)`
const DELETE_COMMAND = `DELETE FROM commands WHERE idCommand = (?)`

func FreshInstall(db *sql.DB) error {
	createCredTable(db)
	createHostsTable(db)
	createCommandBar(db)
	createCommandTable(db)
	createCommandParseTable(db)
	var begin typelib.CommandLibrary
	begin.ImportFromJson(inital.CommandOutTemplateBytes)
	err := InsertCommandIntoLib(db, begin)
	if err != nil {
		return err
	}
	initalParseDB, err := inital.GenerateCommandTemplate()
	if err != nil {
		return err
	}
	err = InsertCommandParseTable(db, initalParseDB)
	if err != nil {
		return err
	}
	return nil
}

// TODO Pass argument for location of new database
func MakeNewDatabase() *sql.DB {
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	sqliteDatabase := OpenDatabase()

	err = FreshInstall(sqliteDatabase)
	if err != nil {
		log.Fatal(err)
	}

	return sqliteDatabase
}

// TODO Pass argument for location of database
func OpenDatabase() *sql.DB {
	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db")
	if err != nil {
		log.Fatal(err)
	}
	return sqliteDatabase
}

func createCredTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE creds (
		"idCred" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"username" TEXT  NOT NULL,
		"password" TEXT,
		"domain" TEXT,
		"hash" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create Creds table...")
	statement, err := db.Prepare(createStudentTableSQL) // Pregolanpare SQL Statement
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("Cred table created")
}

func createCommandTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE commands (
		"idCommand" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"displayString" TEXT  NOT NULL,
		"templateString" TEXT NOT NULL,
		"exampleString" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Creating Command Table")
	statement, err := db.Prepare(createStudentTableSQL) // Pregolanpare SQL Statement
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("Created Command Table")
}

func createCommandBar(db *sql.DB) {
	// commandID,user, domain, password,hash,host,command
	createCommandBarSQL := `CREATE TABLE commandBar (
		"commandID" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"user" TEXT  NOT NULL,
		"domain" TEXT,
		"password" TEXT,
		"hash" TEXT,
		"host" TEXT,
		"command" TEXT	
	  );` // SQL Statement for Create Table

	log.Println("Creating CommandBar Table")
	statement, err := db.Prepare(createCommandBarSQL) // Pregolanpare SQL Statement
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("Created Command Table")
}

func createHostsTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE hosts (
		"idHost" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"ip" TEXT  NOT NULL,
		"hostname" TEXT,
		"fqdn" TEXT,
		"usersAdmin" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create hosts table...")
	statement, err := db.Prepare(createStudentTableSQL) // Pregolanpare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("Hosts table created")
}

func createCommandParseTable(db *sql.DB) {
	createCommandTableSQL := `CREATE TABLE command_parse (
    	"idParse" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"cmdDisplay" Text,		
		"alias" BLOB,
		"parseType" TEXT,
		"commandMatch" BLOB	
	  );`

	statement, err := db.Prepare(createCommandTableSQL) // Pregolanpare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
}

func InsertCommandParseTable(db *sql.DB, cmdObject typelib.CommandParseDB) error {
	for x := range cmdObject.Array {

		sqlQuery := `INSERT INTO command_parse(cmdDisplay, alias, parseType, commandMatch) VALUES (?,?,?,?)`
		statement, err := db.Prepare(sqlQuery)

		if err != nil {
			return err
		}
		aliasblob, err := json.Marshal(cmdObject.Array[x].Alias)
		if err != nil {
			return err
		}
		cmdMatchBlob, err := json.Marshal(cmdObject.Array[x].CommandMatch)
		if err != nil {
			return err
		}
		_, err = statement.Exec(cmdObject.Array[x].CommandName, aliasblob, cmdObject.Array[x].ParseType, cmdMatchBlob)
		if err != nil {
			return err
		}
		statement.Close()

	}
	return nil
}

func RetreiveParseTable(db *sql.DB) (typelib.CommandParseDB, error) {
	var ParseTable typelib.CommandParseDB
	sqlQuery := `SELECT cmdDisplay,alias,parseType,commandMatch  FROM command_parse`
	statement, err := db.Query(sqlQuery)
	if err != nil {
		return ParseTable, err
	}
	defer statement.Close()

	i := 0
	for statement.Next() {
		i++
		var newParse typelib.CommandTemplate
		var cmdDisplay string
		var alias []byte
		var parseType string
		var commandMatch []byte
		statement.Scan(&cmdDisplay, &alias, &parseType, &commandMatch)
		newParse.CommandName = cmdDisplay
		newParse.ParseType = parseType
		err := json.Unmarshal(alias, &newParse.Alias)
		if err != nil {
			return ParseTable, err
		}
		err = json.Unmarshal(commandMatch, &newParse.CommandMatch)
		if err != nil {
			return ParseTable, err
		}
		ParseTable.Array = append(ParseTable.Array, newParse)

	}
	return ParseTable, nil

}

func InsertCreds(db *sql.DB, entry typelib.CredEntry) string {
	insertStudentSQL := `INSERT INTO creds(username, domain, password,hash) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		return err.Error()
	}
	_, err = statement.Exec(entry.User, entry.Domain, entry.Password, entry.Hash)
	if err != nil {
		return err.Error()
	}
	return ""
}

func DeleteEntry(db *sql.DB, table string, ID int) error {
	var deleteQuery string
	switch table {
	case "host":
		deleteQuery = DELETE_HOST
	case "cred":
		deleteQuery = DELETE_CRED
	case "command":
		deleteQuery = DELETE_COMMAND
	}

	statement, err := db.Prepare(deleteQuery)
	if err != nil {
		return err
	}
	_, err = statement.Exec(ID)
	if err != nil {
		return err
	}
	return nil
}

func GetCred(db *sql.DB, ID int) (typelib.CredEntry, error) {
	var entry typelib.CredEntry
	query := `SELECT idCred,username,domain,password,hash FROM creds WHERE idCred = ?`
	query2, err := db.Prepare(query)
	statement, err := query2.Query(ID)
	defer statement.Close()
	if err != nil {
		return entry, err
	}
	// Make sure entry is not zero
	i := 0
	for statement.Next() {
		i++
		var id int
		var user string
		var domain string
		var password string
		var hash string
		statement.Scan(&id, &user, &domain, &password, &hash)
		entry.ID = id
		entry.User = user
		entry.Password = password
		entry.Hash = hash
		entry.Domain = domain
		return entry, nil
	}
	if i == 0 {
		return entry, errors.New("no entries found")
	}
	return entry, errors.New("could not build cred object")

}

func UpdateCred(db *sql.DB, entry typelib.CredEntry) error {
	// TODO implement history
	updateCommand := `UPDATE creds SET username = ?,domain = ?,password = ?,hash = ? WHERE idCred = ?`
	updateStatement, err := db.Prepare(updateCommand)
	if err != nil {
		return err
	}
	_, err = updateStatement.Exec(entry.User, entry.Domain, entry.Password, entry.Hash, entry.ID)
	if err != nil {
		return err
	}
	return nil
}

func SetCredBarEntry(db *sql.DB, commmandBarObject typelib.CommandBar) error {
	query := `SELECT * FROM commandbar WHERE commandID = 1`
	statement, err := db.Prepare(query)
	doQuery, err := statement.Query()
	if err != nil {
		return err
	}
	defer doQuery.Close()

	//Determine if there is an existing entry, if not create one and return
	x := 0
	for doQuery.Next() {
		x++
	}
	if x < 1 {
		insertSQL := `INSERT INTO commandbar(commandID,user, domain, password,hash,host,command ) VALUES (?,?, ?, ?, ?,?,?)`
		statement, err := db.Prepare(insertSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			return err
		}
		defer statement.Close()
		_, err = statement.Exec(
			1,
			commmandBarObject.User,
			commmandBarObject.Domain,
			commmandBarObject.Password,
			commmandBarObject.Hash,
			commmandBarObject.Host,
			commmandBarObject.Command,
		)
		if err != nil {
			return err
		} else {
			return nil
		}
	} else {

		//updateSQL := `UPDATE commandbar SET  `

		// Update Exisitng Entry
		// Note to self, make sure strings passed are in single quotes
		// https://stackoverflow.com/questions/67608290/how-to-update-sqlite-using-go-without-other-libraries

		updateStatment := `UPDATE commandbar SET user = ?,domain = ?,password = ?,hash = ?,host = ?,command = ? WHERE commandID = 1`
		statement, err := db.Prepare(updateStatment)
		if err != nil {
			return err
		}
		_, err1 := statement.Exec(
			commmandBarObject.User,
			commmandBarObject.Domain,
			commmandBarObject.Password,
			commmandBarObject.Hash,
			commmandBarObject.Host,
			commmandBarObject.Command,
		)
		if err1 != nil {
			return err1
		} else {
			return nil
		}

	}
}

func GetCommandBarEntry(db *sql.DB) (typelib.CommandBar, error) {
	// TODO when history is implemented, order by descending, get top 1
	var entry typelib.CommandBar
	query := `SELECT user,domain,password,hash,host,command FROM commandBar WHERE commandID = 1`
	statement, err := db.Query(query)
	defer statement.Close()
	if err != nil {
		return entry, err
	}
	// Make sure entry is not zero
	i := 0
	for statement.Next() {
		i++
		var user string
		var domain string
		var password string
		var hash string
		var host string
		var command string
		statement.Scan(&user, &domain, &password, &hash, &host, &command)
		entry.User = user
		entry.Domain = domain
		entry.Password = password
		entry.Hash = hash
		entry.Host = host
		entry.Command = command

		return entry, nil
	}
	if i != 1 {
		return entry, nil
	}
	return entry, errors.New("could not build commandBar object")

}

func GetCredTableSQLEntries(db *sql.DB) ([]typelib.CredEntry, error) {
	var entries []typelib.CredEntry
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

		entries = append(entries, cred)

	}
	return entries, nil
}

func InsertCommandIntoLib(db *sql.DB, commandLib typelib.CommandLibrary) error {
	for x := range commandLib.ListOfCommands {
		insertSQL := `INSERT INTO commands(templateString,displayString,exampleString) VALUES (?,?,?)`
		statement, err := db.Prepare(insertSQL) // Prepare statement.
		defer statement.Close()
		// This is good to avoid SQL injections
		if err != nil {
			return err
		}

		_, err = statement.Exec(commandLib.ListOfCommands[x].Template,
			commandLib.ListOfCommands[x].Short,
			commandLib.ListOfCommands[x].Example)
		if err != nil {
			return err
		}

	}
	return nil
}

//TODO, delete single entry

func GetCommandLib(db *sql.DB) (typelib.CommandLibrary, error) {
	var entries typelib.CommandLibrary
	rows, err := db.Query("SELECT idCommand,templateString,displayString,exampleString FROM commands")
	defer rows.Close()
	if err != nil {
		return typelib.CommandLibrary{}, err
	}

	for rows.Next() {
		var cmd typelib.CommandBuild
		var ID string
		var command string
		var display string
		var example string
		rows.Scan(&ID, &command, &display, &example)

		cmd.ID = ID
		cmd.Template = command
		cmd.Short = display
		cmd.Example = example

		entries.ListOfCommands = append(entries.ListOfCommands, cmd)

	}
	return entries, nil
}

func InsertHost(db *sql.DB, hostEntry typelib.HostEntry) error {
	log.Println("Inserting Creds record ...")
	insertHostSQL := `INSERT INTO hosts(ip, hostname, fqdn,usersAdmin) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertHostSQL) // Prepare statement.
	defer statement.Close()
	// This is good to avoid SQL injections
	if err != nil {
		return err
	}
	_, err = statement.Exec(hostEntry.IP, hostEntry.Hostname, hostEntry.FQDN, hostEntry.Admins)
	if err != nil {
		return err
	}
	return nil
}

func GetHostList(db *sql.DB) ([]typelib.HostEntry, error) {
	var hostlists []typelib.HostEntry
	rows, err := db.Query("SELECT idHost,ip,fqdn,hostname,usersAdmin FROM hosts")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var host typelib.HostEntry
		var id int
		var ip string
		var fqdn string
		var hostname string
		var userAdmin string
		rows.Scan(&id, &ip, &fqdn, &hostname, &userAdmin)

		host.ID = id
		host.IP = ip
		host.FQDN = fqdn
		host.Hostname = hostname
		host.Admins = userAdmin

		hostlists = append(hostlists, host)

	}
	return hostlists, nil
}

func GetHost(db *sql.DB, ID int) (typelib.HostEntry, error) {
	var entry typelib.HostEntry
	query := `SELECT idHost,IP,FQDN,Hostname FROM hosts WHERE idHost = (?)`
	statement, err := db.Prepare(query)
	defer statement.Close()
	rows, err := statement.Query(ID)
	defer rows.Close()

	if err != nil {
		return entry, err
	}
	// Make sure entry is not zero
	i := 0
	for rows.Next() {
		i++
		var idHost int
		var IP string
		var FQDN string
		var Hostname string
		rows.Scan(&idHost, &IP, &FQDN, &Hostname)
		entry.ID = idHost
		entry.IP = IP
		entry.FQDN = FQDN
		entry.Hostname = Hostname

		return entry, nil
	}
	if i == 0 {
		return entry, errors.New("no entries found")
	}
	return entry, errors.New("could not build cred object")

}

func UpdateHost(db *sql.DB, entry typelib.HostEntry) error {
	// TODO implement history
	updateCommand := `UPDATE hosts SET IP = ?,FQDN = ?,Hostname = ? WHERE idHost = ?`
	updateStatement, err := db.Prepare(updateCommand)
	defer updateStatement.Close()
	if err != nil {
		return err
	}
	_, err = updateStatement.Exec(entry.IP, entry.FQDN, entry.Hostname, entry.ID)

	if err != nil {
		return err
	}
	return nil
}

func AddBlankHost(db *sql.DB) (entryNumber int, err error) {
	insertHostSQL := `INSERT INTO hosts(ip, hostname, fqdn,usersAdmin) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertHostSQL) // Prepare statement.
	defer statement.Close()
	// This is good to avoid SQL injections
	if err != nil {
		return 0, err
	}
	id, err := statement.Exec("", "", "", "")
	if err != nil {
		return 0, err
	}
	num, _ := id.LastInsertId()

	return int(num), nil

}

func GetCMD(db *sql.DB, ID int) (typelib.CommandBuild, error) {
	var entry typelib.CommandBuild
	query := `SELECT idCommand,displayString,templateString,exampleString FROM commands WHERE idCommand = (?)`
	statement, err := db.Prepare(query)
	defer statement.Close()
	rows, err := statement.Query(ID)
	defer rows.Close()

	if err != nil {
		return entry, err
	}
	// Make sure entry is not zero
	i := 0
	for rows.Next() {
		i++
		var idCommand string
		var displayString string
		var templateString string
		var exampleString string
		rows.Scan(&idCommand, &displayString, &templateString, &exampleString)
		entry.ID = idCommand
		entry.Template = templateString
		entry.Example = exampleString
		entry.Short = displayString

		return entry, nil
	}
	if i == 0 {
		return entry, errors.New("no entries found")
	}
	return entry, errors.New("could not build command object")

}

func UpdateCMD(db *sql.DB, entry typelib.CommandBuild) error {
	// TODO implement history
	updateCommand := `UPDATE commands SET displayString = ?,templateString = ?,exampleString = ? WHERE idCommand = ?`
	updateStatement, err := db.Prepare(updateCommand)
	defer updateStatement.Close()
	if err != nil {
		return err
	}
	_, err = updateStatement.Exec(entry.Short, entry.Template, entry.Example, entry.ID)

	if err != nil {
		return err
	}
	return nil
}

//Command
//Example
//Display
