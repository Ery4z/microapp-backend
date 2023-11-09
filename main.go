package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/labstack/echo/v4"
)

var db *sql.DB

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/data", postData)

	var err error
	db, err = initSqlite3()
	if err != nil {
		log.Fatal("Failed instanciating the db object: " + err.Error())
	}

	e.Logger.Fatal(e.Start(":1323"))

}

func postData(c echo.Context) error {
	sensorId := c.FormValue("sensorId") // Unique Id of the sender
	groupId := c.FormValue("groupId")   // Unique Id of the sender
	dataType := c.FormValue("dataType") // string || int || float || date
	dataUnit := c.FormValue("dataUnit") // Unit of the data if relevant
	dataInfo := c.FormValue("dataInfo") // Information about the data
	data := c.FormValue("data")

	err := insertDataSqlite3(db, sensorId, groupId, dataType, dataUnit, dataInfo, data)
	if err != nil {
		log.Println(err)
	}

	return c.String(http.StatusOK, "Sensor: "+sensorId+" ("+dataType+") | "+data+dataUnit+" | "+dataInfo)

}

func initSqlite3() (*sql.DB, error) {
	dbPath := "db.sqlite"

	db, err := sql.Open("sqlite", dbPath)
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
	defer rows.Close()

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
	exist := false
	for rows.Next() {
		var (
			cid          int
			name         string
			dataType     string
			notNull      int
			defaultValue sql.NullString
			primaryKey   int
		)

		err := rows.Scan(&cid, &name, &dataType, &notNull, &defaultValue, &primaryKey)
		if err != nil {
			log.Fatal(err)
		}

		// Check if the column name matches the one you want to verify,
		// and if the data type matches your expected type.
		if name == "data" && !(dataType == dtype) {
			log.Fatal("The table " + groupId + " cannot handle only " + dataType + " not " + dtype)
		}
		exist = true
	}

	if !exist {

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

	insertStatement := fmt.Sprintf("INSERT INTO %s (sensorId,dataUnit,dataInfo,data,time) VALUES (?,?,?,?,?)", groupId)

	result, err := db.Exec(insertStatement, sensorId, dataUnit, dataInfo, data, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	rowsAffected, err := result.RowsAffected()

	fmt.Printf("Inserted %d rows.\n", rowsAffected)

	if err != nil {
		log.Fatal(err)
	}

	return nil

}
