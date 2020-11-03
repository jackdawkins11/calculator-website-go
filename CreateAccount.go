package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

//Tries to insert the given username and password
//into the database. Returns the number of affected
//rows and boolean that indicates if there was an error
func insertUser(username, password string) (int, bool) {
	sqlQuery := "INSERT INTO Users(username, password) VALUES(?, ?)"
	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Println(fmt.Sprintf("db.Prepare( %v ) failed with %v", sqlQuery, err))
		return 0, true
	}
	res, err := stmt.Exec(username, password)
	if err != nil {
		log.Println(fmt.Sprintf("stmt.Exec for '%v' failed with %v", sqlQuery, err))
		return 0, true
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Println(fmt.Sprintf("result.RowsAffected() for '%v' failed with %v", sqlQuery, err))
		return 0, true
	}
	return int(rowCnt), false
}

//Returns whether there is an accont with the given
//username and if there was an error
func userExists(username string) (bool, bool) {
	var exists bool
	err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM Users WHERE username = ?", username).Scan(&exists)
	errBool := false
	if err != nil {
		fmt.Println("Error in authenticateUser: ", err)
		errBool = true
	}
	return exists, errBool
}

func checkUsernameAndPassword(username, password string) (bool, string) {
	if len(username) < 8 {
		return false, "Username must be at least 8 characters long"
	}
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}
	matched, err := regexp.Match(`[0-9]`, []byte(password))
	if err != nil {
		return false, "There was an error"
	}
	if !matched {
		return false, "Password must contain a digit"
	}
	matched, err = regexp.Match(`[!@#$%^&*()]`, []byte(password))
	if err != nil {
		return false, "There was an error"
	}
	if !matched {
		return false, "Password must contain a character from !@#$%^&*()"
	}
	matched, err = regexp.Match(`[a-z]`, []byte(password))
	if err != nil {
		return false, "There was an error"
	}
	if !matched {
		return false, "Password must contain a lowercase letter"
	}
	matched, err = regexp.Match(`[A-Z]`, []byte(password))
	if err != nil {
		return false, "There was an error"
	}
	if !matched {
		return false, "Password must contain a uppercase letter"
	}
	return true, ""
}

//Handles requests for creating accounts.
//Required x-www-urlencoded params: username, password
//Tries to make an account with the given username and password.
//Will return an http internal server error
//or json.
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	goodUsernameAndPassword, message := checkUsernameAndPassword(username, password)

	if !goodUsernameAndPassword {
		responseData := map[string]interface{}{
			"error":          false,
			"createdAccount": false,
			"message":        message,
		}
		sendResponse(w, responseData)
		return
	}

	alreadyInUse, errBool := userExists(username)

	if !errBool && alreadyInUse {
		responseData := map[string]interface{}{
			"error":          false,
			"createdAccount": false,
			"message":        "That username is taken",
		}
		sendResponse(w, responseData)
		return
	} else if errBool {
		responseData := map[string]interface{}{
			"error":          true,
			"createdAccount": false,
			"message":        "",
		}
		sendResponse(w, responseData)
	}

	rowCnt, errBool := insertUser(username, password)

	if errBool {
		responseData := map[string]interface{}{
			"error":          true,
			"createdAccount": false,
			"message":        "Error inserting into database",
		}
		sendResponse(w, responseData)
		return
	} else if !errBool && rowCnt != 1 {
		responseData := map[string]interface{}{
			"error":          true,
			"createdAccount": false,
			"message":        "Error inserting into db: number of affected rows wrong",
		}
		sendResponse(w, responseData)
	} else if !errBool && rowCnt == 1 {
		responseData := map[string]interface{}{
			"error":          false,
			"createdAccount": true,
			"message":        "",
		}
		sendResponse(w, responseData)
	}

}
