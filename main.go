package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ilightthings/GED/html"
	"github.com/ilightthings/GED/mysql"
	"github.com/ilightthings/GED/parseinput"
	"github.com/ilightthings/GED/transform"
	"github.com/ilightthings/GED/typelib"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
)

func main() {

	//sqliteDatabase := mysql.MakeNewDatabase()
	sqliteDatabase := mysql.OpenDatabase()

	r := gin.Default()
	r.LoadHTMLFiles("html/import.html")
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "string", html.GenerateImportPage(sqliteDatabase))
	})

	r.POST("/import-command", func(c *gin.Context) {
		command := c.PostForm("cred-import")
		responce := parseinput.ParseData(command)
		err := responce.Verify()
		// Empty or invalid Entry
		if err != nil {
			c.String(500, "Empty or Invalid entry. Requires at least one field.")
		} else {
			sqlres := mysql.InsertCreds(sqliteDatabase, responce.User, responce.Domain, responce.Password, responce.Hash)

			// SQL Error issue
			if sqlres != "" {
				c.String(999, sqlres)
			} else {
				c.String(200, responce.StringCreds())
			}
		}
	})

	r.POST("/import-blob", func(c *gin.Context) {
		commandblob := strings.Split(c.PostForm("command-blob"), "\n")
		credArray := parseinput.IdentifyCMEline(commandblob)
		if len(credArray.CredEntries) < 0 {
			c.String(500, "No entries found using regex.")
			return
		}
		var entries []string
		for _, y := range credArray.CredEntries {
			err := y.Verify()
			if err == nil {
				mysql.InsertCreds(sqliteDatabase, y.User, y.Domain, y.Password, y.Hash)
				entries = append(entries, y.StringCreds())
			}
		}
		c.String(200, "Added Creds: \n"+strings.Join(entries, "\n"))

	})

	r.GET("/deletecred/:id", func(c *gin.Context) {
		idofcred, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not int")
			return
		}
		err = mysql.DeleteCred(sqliteDatabase, idofcred)
		if err != nil {
			c.String(500, "Could not Delete: "+err.Error())
		} else {
			c.Redirect(301, "/creds")
			return
		}

	})

	r.GET("/creds", func(c *gin.Context) {
		data := html.GenerateCredsTable(sqliteDatabase)
		c.Data(200, "string", data)

	})

	r.GET("/cme/:id", func(c *gin.Context) {
		idofcred, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not a number")
			return
		}
		target, err := c.Cookie("target")
		if err != nil {
			target = "replaceme"
		}

		c.String(200, transform.CMEOUT(sqliteDatabase, idofcred, target))

	})

	r.GET("/imp/:id", func(c *gin.Context) {
		idofcred, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not a number")
			return
		}
		target, err := c.Cookie("target")
		if err != nil {
			target = "replaceme"
		}

		c.String(200, transform.IMPOUT(sqliteDatabase, idofcred, target))

	})

	r.GET("/updateID/:id", func(c *gin.Context) {

		ID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not a number")
			return
		}

		responce, err := html.GenerateCredUpdate(sqliteDatabase, ID)
		if err != nil {
			c.String(500, err.Error())
		}
		c.Data(200, "string", responce)
	})

	r.POST("/updateID/:id", func(c *gin.Context) {
		var updateCredObj typelib.CredEntry
		ID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not a number")
			return
		}
		err = c.ShouldBind(&updateCredObj)
		if err != nil {
			c.String(500, "Could not bind JSON\n"+err.Error())
			return
		}
		updateCredObj.ID = ID
		err = mysql.UpdateCred(sqliteDatabase, updateCredObj)
		if err != nil {
			c.String(500, "Could not update Database Entry "+err.Error())
			return
		} else {
			c.String(200, "Entry Updated")
			return
		}

	})

	r.GET("/settings", func(c *gin.Context) {
		data := html.GenerateSettingsPage()
		c.Data(200, "strings", data)

	})

	r.GET("/creds/json/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not a number")
			return
		} else {
			credsEntry, err := mysql.GetCred(sqliteDatabase, id)
			if err != nil {
				c.String(500, "Error with SQL:"+err.Error())
				return
			}
			c.JSON(200, credsEntry)
		}

	})
	r.GET("/SetCred/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not a number")
			return
		} else {
			currentBar, err := mysql.GetCommandBarEntry(sqliteDatabase)
			if err != nil {
				c.String(500, "Error with SQL:"+err.Error())
				return
			}

			credsEntry, err := mysql.GetCred(sqliteDatabase, id)
			if err != nil {
				c.String(500, "Error with SQL:"+err.Error())
				return
			}

			currentBar.User = credsEntry.User
			currentBar.Domain = credsEntry.Domain
			currentBar.Password = credsEntry.Password
			currentBar.Hash = credsEntry.Hash

			err = mysql.SetCredBarEntry(sqliteDatabase, currentBar)
			if err != nil {
				c.String(500, "Error with SQL:"+err.Error())
				return
			}

			c.String(200, "It worked!")

		}

	})

	r.POST("/creds/execute", func(c *gin.Context) {
		var passedCommands typelib.CommandBar
		var newCommand typelib.CommandBuild
		passedCommands.User = c.PostForm("footer-user")
		passedCommands.Password = c.PostForm("footer-password")
		passedCommands.Domain = c.PostForm("footer-domain")
		passedCommands.Hash = c.PostForm("footer-hash")
		passedCommands.Host = c.PostForm("footer-host")
		passedCommands.Command = c.PostForm("command")
		newCommand.Command = c.PostForm("command")

		err := mysql.SetCredBarEntry(sqliteDatabase, passedCommands)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		result := newCommand.BuildCommand(passedCommands)
		c.String(200, result)
		return

	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
