package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nineteenseventy/minichat/server/util"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	conn := util.GetDatabase()
	rows, err := conn.Query(r.Context(), "SELECT id, name FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var result struct {
		Users []User `json:"users"`
	}

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result.Users = append(result.Users, User{ID: id, Name: name})
	}

	util.JSONResponse(w, result)
}

func UserRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/users", userHandler)
	return r
}
