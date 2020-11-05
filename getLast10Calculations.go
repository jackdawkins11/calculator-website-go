package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

/*
	Finds the username associated with the given
	primary key. Returns
		username (string) the username
		(bool) whether there was an error accessing the username, including
			if it is an invalid primary key
*/
func getUsername(primaryKey int) (string, bool) {
	username := ""
	err := db.QueryRow("SELECT Username FROM Users WHERE PrimaryKey = ?", primaryKey).Scan(&username)
	errBool := false
	if err != nil {
		fmt.Println(err)
		errBool = true
	}
	return username, errBool
}

/*
	Gets the most recent 10 calculations in the database.
	Returns
		([]map[string]interface{}) an array of maps. Each map is accessed with
			the db field specified by a string and stores the data from the db.
			This is nil if there was an error.
*/
func getCalculations() []map[string]interface{} {
	stmt, err := db.Prepare("select X, Op, Y, Val, Date, UserKey from Calculations order by Date desc limit 10")
	if err != nil {
		fmt.Println("getCalculations(): ", err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		fmt.Println("getCalculations(): ", err)
		return nil
	}
	calculations := make([]map[string]interface{}, 0)
	defer rows.Close()
	for rows.Next() {
		var x, op, y, val, date string
		var userKey int
		err := rows.Scan(&x, &op, &y, &val, &date, &userKey)
		if err != nil {
			fmt.Println("getCalculations(): ", err)
			return nil
		}
		m := map[string]interface{}{
			"X":       x,
			"Op":      op,
			"Y":       y,
			"Val":     val,
			"Date":    date,
			"UserKey": userKey,
		}
		calculations = append(calculations, m)
		//fmt.Println(x, " ", op, " ", y, " = ", val, " ", date, " key ", userKey)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("getCalculations(): ", err)
		return nil
	}
	return calculations
}

/*
	For each map in the given slice of maps, it adds a Username field and
	removes the UserKey field. Returns
		(bool) whether there was an error
*/
func addUsernameToCalculations(calculations []map[string]interface{}) bool {
	retErr := false
	for i := 0; i < len(calculations); i++ {
		username, err := getUsername(calculations[i]["UserKey"].(int))
		if err {
			retErr = true
			fmt.Println("Error getting primary key of user")
		} else {
			calculations[i]["Username"] = username
		}
		calculations[i]["UserKey"] = nil
	}
	return retErr
}

/*
	Handles the request to getLast10Calculations.
	Returns json containing
		error (bool) whether there was an error
		calculations ([]map[string]interface{}) array of maps. Each map's keys are a field
			from the database and values are the values from the database
*/
func getLast10Calculations(w http.ResponseWriter, r *http.Request) {

	calculations := getCalculations()

	if calculations == nil {
		responseData := map[string]interface{}{
			"error": true,
		}
		writeResponse(w, responseData)
		return
	}

	errBool := addUsernameToCalculations(calculations)

	responseData := map[string]interface{}{
		"error":        errBool,
		"calculations": calculations,
	}

	writeResponse(w, responseData)
}
