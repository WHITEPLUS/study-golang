package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NozomiSugiyama/study-golang/chap1/controller"
	"github.com/NozomiSugiyama/study-golang/chap1/model"
)

func main() {
	http.HandleFunc("/users", usersHandler)

	http.ListenAndServe(":8080", nil)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		users, err := controller.ListUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		switch r.URL.Query().Get("format") {
		case "html":
			w.Header().Set("Content-Type", "text/html")
			for _, value := range users {
				fmt.Fprint(w, toStringFromUser(value) + "\n")
			}
		case "plain":
			w.Header().Set("Content-Type", "text/plain")
			for _, value := range users {
				fmt.Fprint(w, toStringFromUser(value) + "\n")
			}
		// json
		default:
			res, err := json.Marshal(users)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(res))
		}
	case "POST":
		decoder := json.NewDecoder(r.Body)
		var user model.UserCreate
		err := decoder.Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newUser, err := controller.CreateUser(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switch r.URL.Query().Get("format") {
		case "html":
			w.Header().Set("Content-Type", "text/html")
			fmt.Fprint(w, toStringFromUser(newUser) + "\n")
		case "plain":
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprint(w, toStringFromUser(newUser) + "\n")
		// json
		default:
			res, err := json.Marshal(newUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(res))
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func toStringFromUser(user model.User) string {
	return strconv.Itoa(user.ID)+","+user.Name+","+user.Email+","+user.Birthday+","+user.PhoneNumber
}
