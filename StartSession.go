package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

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
