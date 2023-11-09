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

// Define a new struct for the group information.
type Group struct {
	GroupID     string `json:"groupId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

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

	// Initialize the groups table.
	initGroupsTable()

	// Define new endpoints.
	e.POST("/groups", createGroup)
	e.PUT("/groups/:id", updateGroup)
	e.GET("/groups/:id", getGroup)
	e.GET("/groups", listGroups)

	e.Logger.Fatal(e.Start(":1323"))

}

func postData(c echo.Context) error {
	sensorId := c.FormValue("sensorId") // Unique Id of the sender
	groupId := c.FormValue("groupId")   // Unique Id of the group
	dataType := c.FormValue("dataType") // string || int || float || date
	dataUnit := c.FormValue("dataUnit") // Unit of the data if relevant
	dataInfo := c.FormValue("dataInfo") // Information about the data
	data := c.FormValue("data")

	// Check if the group exists in the database
	if !groupExists(groupId) {
		// If not, insert the group with a placeholder name and description
		err := insertGroup(groupId, "Default Name", "Default Description")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert group: "+err.Error())
		}
	}

	// Insert the data into the database
	err := insertDataSqlite3(db, sensorId, groupId, dataType, dataUnit, dataInfo, data)
	if err != nil {
		log.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to insert data: "+err.Error())
	}

	return c.String(http.StatusOK, "Sensor: "+sensorId+" ("+dataType+") | "+data+dataUnit+" | "+dataInfo)
}

// groupExists checks if a group with the given ID already exists in the database.
func groupExists(groupId string) bool {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM groups WHERE groupId = ?)", groupId).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Error checking if group exists:", err)
		return false
	}
	return exists
}

// insertGroup inserts a new group with the given ID, name, and description into the database.
func insertGroup(groupId, name, description string) error {
	_, err := db.Exec("INSERT INTO groups (groupId, name, description) VALUES (?, ?, ?)", groupId, name, description)
	return err
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

	err := verifyTableExistance(groupId, dtype)
	if err != nil {
		log.Println(err.Error())
	}

	err = verifyTableSchema(groupId, "data", dtype)
	if err != nil {
		log.Println(err.Error())
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

func verifyTableSchema(tableName string, columnToVerify string, typeWanted string) error {
	// Check if the table exist
	query := fmt.Sprintf("PRAGMA table_info(%s);", tableName)

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

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
			return err
		}

		// Check if the column name matches the one you want to verify,
		// and if the data type matches your expected type.
		if name == columnToVerify && !(dataType == typeWanted) {
			return fmt.Errorf("The table " + tableName + " cannot handle only " + dataType + " not " + typeWanted)
		}
	}
	return nil
}

func verifyTableExistance(tableName string, dataType string) error {
	query := fmt.Sprintf("PRAGMA table_info(%s);", tableName)

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		query = fmt.Sprintf(`CREATE TABLE %s  (
			ID INTEGER PRIMARY KEY,
			sensorId TEXT,
			dataUnit TEXT,
			dataInfo TEXT,
			data %s,
			time TIME
			);`, tableName, dataType)
		_, err = db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func initGroupsTable() {
	// Create the groups table if it doesn't exist.
	createTableSQL := `CREATE TABLE IF NOT EXISTS groups (
        groupId TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        description TEXT
    );`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating groups table: %s", err)
	}
}

func createGroup(c echo.Context) error {
	// Parse the request body to get the group data.
	var group Group
	if err := c.Bind(&group); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Insert the group data into the database.
	_, err := db.Exec("INSERT INTO groups (groupId, name, description) VALUES (?, ?, ?)",
		group.GroupID, group.Name, group.Description)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, group)
}

func updateGroup(c echo.Context) error {
	// Get the group ID from the path.
	groupId := c.Param("id")

	// Parse the request body for the updated data.
	var group Group
	if err := c.Bind(&group); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Update the group data in the database.
	_, err := db.Exec("UPDATE groups SET name = ?, description = ? WHERE groupId = ?",
		group.Name, group.Description, groupId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, group)
}

func getGroup(c echo.Context) error {
	// Get the group ID from the path.
	groupId := c.Param("id")

	// Retrieve the group data from the database.
	var group Group
	row := db.QueryRow("SELECT groupId, name, description FROM groups WHERE groupId = ?", groupId)
	if err := row.Scan(&group.GroupID, &group.Name, &group.Description); err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Group not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, group)
}

func listGroups(c echo.Context) error {
	// Query for all group data in the database.
	rows, err := db.Query("SELECT groupId, name, description FROM groups")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	// Parse the rows into a list of groups.
	var groups []Group
	for rows.Next() {
		var group Group
		if err := rows.Scan(&group.GroupID, &group.Name, &group.Description); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		groups = append(groups, group)
	}

	return c.JSON(http.StatusOK, groups)
}
