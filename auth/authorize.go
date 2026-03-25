package auth

/*
import (
	"log"
	"database/sql"
	"errors"
	"net/http"

	"websocket-server/db"
)

type AuthResult struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Id string `json:"id"`
	Name string `json:"name"`
	Roles string `json:"roles"`
	Privileges string `json:"privileges"`
}

\/**
 * Authorize - Query database for user
 *\/
func Authorize(email string, password string) (*AuthResult, error) {

    result := &AuthResult{
		Success: false,
		Message: "",
		Id: "",
		Name: "",
		Roles: "",
		Privileges: "",
    }


		conn := db.DB_GetConnection()

		if conn != nil {

			query, err := db.LoadSQL("authorize.sql")

			if err != nil {
				log.Fatal(err)
			}

			rows, err := conn.Query(query, password, email)

			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return 
			}

			password_match := false

			for rows.Next() {
				_ = rows.Scan(&result.Id, &result.Name, &password_match, &result.Roles, &result.Privileges)

				if !password_match {
					http.Error(w, "Wrong password", http.StatusUnauthorized)
					return
				}

			}

			if result.Id == "" {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

		} else {
			http.Error(w, "Database unavailable", http.StatusInternalServerError)
			return nil, err
		}
}

*/