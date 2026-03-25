package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"websocket-server/db"
)

type Credentials struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Agent      string `json:"agent"`
	Platform   string `json:"platform"`
	Model      string `json:"model"`
	DeviceName string `json:"deviceName"`
	DeviceType string `json:"deviceType"`
}

/**
 * POST /login
 * -H "x-user-agent: Radar/1.0.0" -H "x-platform: Android" -d '{ email: "admin@hal9k.net", password: "admin" }'
 */
func Login(w http.ResponseWriter, r *http.Request) {

	id, name, roles, privileges := "", "", "", ""

	isGuest := r.URL.Query().Has("guest")
	nosession := r.URL.Query().Has("nosession") // Used by cli

	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)

	//fmt.Println(credentials)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if isGuest {
		id = "dummyGuestId"
		name = credentials.Name
	} else {

		conn := db.DB_GetConnection()

		if conn != nil {

			query, err := db.LoadSQL("authorize.sql")

			if err != nil {
				log.Fatal(err)
			}

			query = strings.ReplaceAll(query, "{{DB_SCHEMA}}", os.Getenv("DB_SCHEMA"))

			rows, err := conn.Query(query, credentials.Password, credentials.Email)

			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			password_match := false

			for rows.Next() {
				_ = rows.Scan(&id, &name, &password_match, &roles, &privileges)

				if !password_match {
					http.Error(w, "Wrong password", http.StatusUnauthorized)
					return
				}

			}

			if id == "" {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

		} else {
			http.Error(w, "Database unavailable", http.StatusInternalServerError)
			return
		}
	}

	// User authenticated or guest

	sessionId := ""

	user := User{
		Id:    id,
		Name:  name,
		Email: credentials.Email,
	}

	if !nosession {

		ip := r.RemoteAddr
		status := "idle"
		updated := time.Now() //.Format(time.RFC3339)

		sess := Session{
			User:       user,
			Agent:      credentials.Agent,
			Platform:   credentials.Platform,
			Model:      credentials.Model,
			DeviceName: credentials.DeviceName,
			DeviceType: credentials.DeviceType,
			Ip:         ip,
			Status:     status,
			Updated:    updated,
		}

		sessionId, err = CreateSession(db.RedisGetConnection(), sess)

		if err != nil {
			log.Println(err)
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			return
		}
	}

	var expiresAt *time.Time = nil

	if isGuest {
		t := time.Now().Add(24 * time.Hour)
		expiresAt = &t
	}

	token, err := generateJWT(sessionId, id, credentials.Email, name, roles, privileges, expiresAt)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"token":"%s", "name":"%s", "id":"%s", "hostname":"%s" }`, token, name, id, hostname)))

	//fmt.Println(token)

	if isGuest {
		log.Printf("Guest %s entered\n", user.Name)
	} else {
		log.Printf("User %s successfully authenticated\n", user.Name)
	}

}
