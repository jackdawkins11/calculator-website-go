package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

//Tries to insert the given calculation
//into the database. Returns the number of affected
//rows and boolean that indicates if there was an error
func insertCalculation(x, op, y, val, date string, userKey int) (int, bool) {
	sqlQuery := "INSERT INTO Calculations(X, Op, Y, Val, Date, UserKey) VALUES(?, ?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		fmt.Println(fmt.Sprintf("db.Prepare( %v ) failed with %v", sqlQuery, err))
		return 0, true
	}
	res, err := stmt.Exec(x, op, y, val, date, userKey)
	if err != nil {
		fmt.Println(fmt.Sprintf("stmt.Exec for '%v' failed with %v", sqlQuery, err))
		return 0, true
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		fmt.Println(fmt.Sprintf("result.RowsAffected() for '%v' failed with %v", sqlQuery, err))
		return 0, true
	}
	return int(rowCnt), false
}

func AddCalculation(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	x := r.PostFormValue("x")
	op := r.PostFormValue("op")
	y := r.PostFormValue("y")
	date := r.PostFormValue("date")
	val := r.PostFormValue("val")

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	primaryKey := session.Values["PrimaryKey"].(int)

	rowCnt, errBool := insertCalculation(x, op, y, val, date, primaryKey)

	responseData := map[string]interface{}{
		"error": errBool || rowCnt != 1,
	}
	sendResponse(w, responseData)
}
