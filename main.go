package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ilightthings/GED/html"
	"github.com/ilightthings/GED/mysql"
	"github.com/ilightthings/GED/parseinput"
	"github.com/ilightthings/GED/typelib"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
	"strings"
)

func main() {

	var sqliteDatabase *sql.DB
	sqliteDatabase = initalizeSQL()

	r := gin.Default()
	r.LoadHTMLFiles("html/import.html")
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "string", html.GenerateImportPage(sqliteDatabase))
	})

	r.POST("/import-command", func(c *gin.Context) {
		command := c.PostForm("cred-import")
		newcred, newhost, err3 := parseinput.ExtractCommandData(command, sqliteDatabase)
		if err3 != nil {
			c.String(500, "Error with host insert:\n "+err3.Error())
			return
		}
		err := newcred.Verify()
		err2 := newhost.Verify()

		// We don't care too much if there is no host
		if err2 == nil {
			err = mysql.InsertHost(sqliteDatabase, newhost)
			if err != nil {
				c.String(500, "Error with host insert:\n "+err.Error())
				return
			}
		}
		// Empty or invalid Cred Entry
		if err != nil {
			c.String(500, "Empty or Invalid entry. Requires at least one field.")
		} else {
			sqlres := mysql.InsertCreds(sqliteDatabase, newcred)

			// SQL Error issue
			if sqlres != "" {
				c.String(999, sqlres)
			} else {
				responce := fmt.Sprintf("Credentials:\n%s\n\nHosts:\n%s", newcred.StringCreds(), newhost.StringHost())
				c.String(200, responce)
			}
		}
	})

	r.POST("/import-blob", func(c *gin.Context) {
		commandblob := strings.Split(c.PostForm("command-blob"), "\n")
		credArray, hostarray := parseinput.IdentifyCMEline(commandblob)
		if len(credArray) < 0 {
			c.String(500, "No entries found using regex.")
			return
		}
		var credentries []string
		for _, y := range credArray {
			err := y.Verify()
			if err == nil {
				mysql.InsertCreds(sqliteDatabase, y)
				credentries = append(credentries, y.StringCreds())
			}
		}
		var hostentries []string
		for _, y := range hostarray {
			mysql.InsertHost(sqliteDatabase, y)
			hostentries = append(hostentries, y.StringHost())
		}

		c.String(200, "Added Creds: \n"+strings.Join(credentries, "\n")+"Host Entries: \n"+strings.Join(hostentries, "\n"))

	})

	r.GET("/delete/:table/:id", func(c *gin.Context) {
		idofcred, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not int")
			return
		}
		err = mysql.DeleteEntry(sqliteDatabase, (c.Param("table")), idofcred)
		if err != nil {
			c.String(500, "Could not Delete: "+err.Error())
		} else {
			c.Redirect(301, "/"+c.Param("table"))
			return
		}

	})

	r.GET("/cred", func(c *gin.Context) {
		data := html.GenerateTableCreds(sqliteDatabase)
		c.Data(200, "string", data)

	})

	r.GET("/updateID/:id", func(c *gin.Context) {

		ID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not a number")
			return
		}

		responce, err := html.GenerateUpdateFormCreds(sqliteDatabase, ID)
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
			c.String(200, "Entry Updated.")
			return
		}

	})

	r.GET("/settings", func(c *gin.Context) {
		data := html.GenerateSettingsPage(sqliteDatabase)
		c.Data(200, "strings", data)

	})

	r.GET("/cred/json/:id", func(c *gin.Context) {
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

	r.POST("/cred/execute", func(c *gin.Context) {
		var passedCommands typelib.CommandBar
		var newCommand typelib.CommandBuild
		passedCommands.User = c.PostForm("footer-user")
		passedCommands.Password = c.PostForm("footer-password")
		passedCommands.Domain = c.PostForm("footer-domain")
		passedCommands.Hash = c.PostForm("footer-hash")
		passedCommands.Host = c.PostForm("footer-host")
		passedCommands.Command = c.PostForm("command")
		newCommand.Template = c.PostForm("command")

		err := mysql.SetCredBarEntry(sqliteDatabase, passedCommands)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		result := newCommand.BuildCommand(passedCommands)
		c.String(200, result)
		return

	})

	r.GET("/host", func(c *gin.Context) {
		data, err := html.GenerateTableHosts(sqliteDatabase)
		if err != nil {
			c.String(500, "Error Generating Host Table:\n"+err.Error())
		}
		c.Data(200, "string", data)
	})

	r.GET("/addHost", func(c *gin.Context) {
		ID, err := mysql.AddBlankHost(sqliteDatabase)
		if err != nil {
			c.String(500, err.Error())
		}

		responce, err := html.GenerateUpdateFormHost(sqliteDatabase, ID)
		if err != nil {
			c.String(500, err.Error())
		}
		c.Data(200, "string", responce)

	})

	r.GET("/setHost/:id", func(c *gin.Context) {
		hostString := c.Param("id")
		currentBar, err := mysql.GetCommandBarEntry(sqliteDatabase)
		if err != nil {
			c.String(500, "Error with SQL:"+err.Error())
			return
		}
		currentBar.Host = hostString

		err = mysql.SetCredBarEntry(sqliteDatabase, currentBar)
		if err != nil {
			c.String(500, "Error with SQL:"+err.Error())
			return
		}

		c.String(200, "It worked!")

	})

	r.GET("/updateHost/:id", func(c *gin.Context) {

		ID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not a number")
			return
		}

		responce, err := html.GenerateUpdateFormHost(sqliteDatabase, ID)
		if err != nil {
			c.String(500, err.Error())
			return
		}
		c.Data(200, "string", responce)
	})

	r.GET("/command", func(c *gin.Context) {
		data, err := html.GenerateTableCommands(sqliteDatabase)
		if err != nil {
			c.String(500, "Error Generating Host Table:\n"+err.Error())
		}
		c.Data(200, "string", data)
	})

	r.POST("/updateHost/:id", func(c *gin.Context) {
		var updateHostObj typelib.HostEntry
		ID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not a number")
			return
		}
		err = c.ShouldBind(&updateHostObj)
		if err != nil {
			c.String(500, "Could not bind JSON\n"+err.Error())
			return
		}
		updateHostObj.ID = ID
		err = mysql.UpdateHost(sqliteDatabase, updateHostObj)
		if err != nil {
			c.String(500, "Could not update Database Entry "+err.Error())
			return
		} else {
			c.String(200, "Entry Updated")
			return
		}
	})

	r.GET("/export-entireDB", func(c *gin.Context) {
		downloadData, err := html.GenerateExportZip(sqliteDatabase)
		if err != nil {
			c.String(500, "Could not Generate Download "+err.Error())
			return
		} else {
			//Force browser download
			c.Header("Content-Disposition", "attachment; filename=database.zip")
			//Browser download or preview
			c.Header("Content-Disposition", "inline;filename=database.zip")
			c.Header("Content-Transfer-Encoding", "binary")
			c.Header("Cache-Control", "no-cache")
			c.Data(200, "application/octet-stream", downloadData)
			return
		}
	})

	r.GET("/updateCommand/:id", func(c *gin.Context) {

		ID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.String(500, "ID is not a number")
			return
		}

		responce, err := html.GenerateUpdateCommand(sqliteDatabase, ID)
		if err != nil {
			c.String(500, err.Error())
		}
		c.Data(200, "string", responce)
	})

	r.POST("/updateCommand/:id", func(c *gin.Context) {
		var updateCmdObj typelib.CommandBuild
		var err error
		if err != nil {
			c.String(500, "ID is not a number")
			return
		}
		err = c.ShouldBind(&updateCmdObj)
		if err != nil {
			c.String(500, "Could not bind JSON\n"+err.Error())
			return
		}

		updateCmdObj.Template = strings.ReplaceAll(updateCmdObj.Template, "\\n", "\n")
		updateCmdObj.Example = strings.ReplaceAll(updateCmdObj.Example, "\\n", "\n")

		err = mysql.UpdateCMD(sqliteDatabase, updateCmdObj)
		if err != nil {
			c.String(500, "Could not update Database Entry "+err.Error())
			return
		} else {
			c.String(200, "Entry Updated")
			return
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func initalizeSQL() *sql.DB {
	if sqldbfile, err := os.Open("sqlite-database.db"); err != nil {
		return mysql.MakeNewDatabase()
	} else {
		sqldbfile.Close()
		return mysql.OpenDatabase()
	}
}
