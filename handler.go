package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
)

// Here you'll move the postData, createGroup, updateGroup, getGroup, and listGroups functions.
// Remember to update the function signatures if needed to make them compatible with this file.
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

func getSensors(c echo.Context) error {
	groupId := c.Param("groupId")
	fmt.Println("Got req")

	// This is a basic security measure to avoid SQL injection by ensuring the groupId is alphanumeric.
	// You should enforce stricter checks depending on your application requirements.
	if !isAlphanumeric(groupId) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid group ID")
	}

	sensors, err := getLatestSensorData(groupId)
	if err != nil {
		log.Printf("Failed to get sensors for group %s: %v", groupId, err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get sensor data")
	}
	return c.JSON(http.StatusOK, sensors)
}

// isAlphanumeric checks if a string contains only alphanumeric characters.
func isAlphanumeric(str string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(str)
}
