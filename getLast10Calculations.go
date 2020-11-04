package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

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
		calculations[i]["userKey"] = nil
	}
	return retErr
}

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
