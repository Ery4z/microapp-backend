package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Here you'll move the verifyTableSchema, verifyTableExistence, insertDataSqlite3, groupExists, and insertGroup functions.
// Adjust the functions as needed for this file.
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

// getLatestSensorData retrieves the latest data for each sensor in the given group table.
func getLatestSensorData(groupId string) ([]SensorData, error) {
	// Ensure the groupId is safe to concatenate by checking against a pattern
	// that allows only alphanumeric characters (and underscore if needed).
	if !isAlphanumeric(groupId) {
		return nil, fmt.Errorf("invalid group ID")
	}

	// Dynamically create the query string using the safe groupId.
	query := fmt.Sprintf(`
    SELECT sensorId, dataUnit, dataInfo, data, MAX(time) as timestamp
    FROM '%s'
    GROUP BY sensorId
    ORDER BY timestamp DESC`, groupId) // Safe to concatenate after isAlphanumeric check

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []SensorData
	for rows.Next() {
		var sd SensorData
		if err := rows.Scan(&sd.SensorID, &sd.DataUnit, &sd.DataInfo, &sd.Data, &sd.Timestamp); err != nil {
			return nil, err
		}
		sensors = append(sensors, sd)
	}

	return sensors, nil
}

type SensorData struct {
	SensorID  string `json:"sensorId"`
	Data      string `json:"data"`
	DataUnit  string `json:"dataUnit"`
	DataInfo  string `json:"dataInfo"`
	Timestamp string `json:"timestamp"`
}
