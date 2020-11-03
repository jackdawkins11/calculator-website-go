package main

import (
	"fmt"
	"net/http"
)

func CheckSession(w http.ResponseWriter, r *http.Request) {
	// Get a session. Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Read the signed in session variable.
	primaryKey := session.Values["PrimaryKey"]

	fmt.Println(fmt.Sprintf("Detected session -- PrimaryKey=%v", primaryKey))

	responseData := map[string]interface{}{
		"hasSession": primaryKey != nil,
	}

	sendResponse(w, responseData)
}
