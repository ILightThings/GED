package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/ilightthings/GED/typelib"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func FreshInstall(db *sql.DB) {
	createCredTable(db)
	createHostsTable(db)
}

// TODO Pass argument for location of new database
func MakeNewDatabase() *sql.DB {
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()

	sqliteDatabase := OpenDatabase()

	FreshInstall(sqliteDatabase)
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

	log.Println("Create student table...")
	statement, err := db.Prepare(createStudentTableSQL) // Pregolanpare SQL Statement
	if err != nil {
		fmt.Println(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("Cred table created")
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

func InsertCreds(db *sql.DB, username string, domain string, password string, hash string) string {
	log.Println("Inserting Creds record ...")
	insertStudentSQL := `INSERT INTO creds(username, domain, password,hash) VALUES (?, ?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		return err.Error()
	}
	_, err = statement.Exec(username, domain, password, hash)
	if err != nil {
		return err.Error()
	}
	return ""
}

func DeleteCred(db *sql.DB, ID int) error {
	deleteStudentSQL := `DELETE FROM creds WHERE idCred = (?)`
	statement, err := db.Prepare(deleteStudentSQL)
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
	query := fmt.Sprintf("SELECT idCred,username,domain,password,hash FROM creds WHERE idCred = %d", ID)
	statement, err := db.Query(query)
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
	updatesmt := `UPDATE creds SET
username = ?,
domain = ?,
password = ?,
hash = ?
WHERE idCred = ?
`
	_, err := db.Exec(updatesmt, entry.User, entry.Domain, entry.Password, entry.Hash, entry.ID)

	if err != nil {
		return err
	}
	return nil
}
