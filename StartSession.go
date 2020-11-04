package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

/*
	Checks if a user with the given username and password
	is in the database. Returns
		(bool) whether the user is authenticated
		(bool) whether there was an error accessing the db
*/
func authenticateUser(username, password string) (bool, bool) {
	var isAuthenticated bool
	err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM Users WHERE username = ? AND password = ?", username, password).Scan(&isAuthenticated)
	errBool := false
	if err != nil {
		fmt.Println("Error in authenticateUser: ", err)
		errBool = true
	}
	return isAuthenticated, errBool
}

/*
	Gets the PrimaryKey of the given username. Returns:
		(int) the primary key
		(bool) whether there was an error accessing the primary key. this also gets set
			when the username is invalid
*/
func getPrimaryKey(username string) (int, bool) {
	primaryKey := -1
	err := db.QueryRow("SELECT primaryKey FROM Users WHERE username = ?", username).Scan(&primaryKey)
	errBool := false
	if err != nil {
		fmt.Println(err)
		errBool = true
	}
	return primaryKey, errBool
}

/*
	Handles the request to StartSession.
	Requires username and password POST parameters in the request.
	Returns json containing
		error (bool)
		hasSession (bool)
	Will return an internal server error in certain situations
*/

func StartSession(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	isAuthenticated, errBool := authenticateUser(username, password)

	//TODO switch from internal server error to returning error as boolean in json
	if !errBool && isAuthenticated {
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		primaryKey, errBool := getPrimaryKey(username)
		if errBool {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("accessed primary key ", primaryKey)
		session.Values["PrimaryKey"] = primaryKey
		if session.Save(r, w) != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	responseData := map[string]interface{}{
		"error":      errBool,
		"hasSession": isAuthenticated,
	}

	sendResponse(w, responseData)
}
