package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/data", postData)
	e.Logger.Fatal(e.Start(":1323"))

}

func postData(c echo.Context) error {
	sensorId := c.FormValue("sensorId") // Unique Id of the sender
	groupId := c.FormValue("groupId")   // Unique Id of the sender
	dataType := c.FormValue("dataType") // string || int || float || date
	dataUnit := c.FormValue("dataUnit") // Unit of the data if relevant
	dataInfo := c.FormValue("dataInfo") // Information about the data
	data := c.FormValue("data")

	return c.String(http.StatusOK, "Sensor: "+sensorId+" ("+dataType+") | "+data+dataUnit+" | "+dataInfo)

}

func initSqlite3() (*sql.DB, error) {
	dir, err := os.MkdirTemp("", "test-")
	if err != nil {
		return nil, err
	}

	defer os.RemoveAll(dir)

	fn := filepath.Join(dir, "db")

	db, err := sql.Open("sqlite", fn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func insertDataSqlite3(db *sql.DB, sensorId string, groupId string, dataType string, dataUnit string, dataInfo string, data string) error {

	// Check if the table exist
	query := fmt.Sprintf("PRAGMA table_info(%s);", groupId)

	rows, err := db.Query(query)
	if err != nil {
		return err
	}

	if !rows.Next() {
		var dtype string

		switch dataType {
		case "string":
			dtype = "TEXT"
		case "text":
			dtype = "TEXT"
		case "float":
			dtype = "REAL"
		case "int":
			dtype = "INTEGER"
		}

		query = fmt.Sprintf(`CREATE TABLE %s  (
			ID INTEGER PRIMARY KEY,
			sensorId TEXT,
			dataUnit TEXT,
			dataInfo TEXT,
			data %s,
			time TIME
			);`, groupId, dtype)
		_, err = db.Exec(query)
		if err != nil {
			return err
		}
	}

	// insert the data

	return nil

}
