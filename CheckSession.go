package main

import (
	"fmt"
	"net/http"
)

/*
	Handles the request for checking if there is a session.
	Checks if the session variable PrimaryKey is set. Returns
	json response containing:
		hasSession (bool)
	Will return an internal server error in certain error
	situations
*/

func CheckSession(w http.ResponseWriter, r *http.Request) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Read the session variable.
	primaryKey := session.Values["PrimaryKey"]

	fmt.Println(fmt.Sprintf("Detected session -- PrimaryKey=%v", primaryKey))

	responseData := map[string]interface{}{
		"hasSession": primaryKey != nil,
	}

	writeResponse(w, responseData)
}
