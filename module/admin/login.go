package admin

import (
	"ekhoes-server/auth"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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

	isGuest := r.URL.Query().Has("guest")
	nosession := r.URL.Query().Has("nosession") // Used by cli

	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)

	//fmt.Println(credentials)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authRes := &AuthResult{
		Success:    false,
		Message:    "",
		Id:         "",
		Name:       "",
		Roles:      "",
		Privileges: "",
	}

	if isGuest {
		authRes.Id = "dummyGuestId"
		authRes.Name = credentials.Name
	} else {

		authRes, err = Authorize(credentials.Email, credentials.Password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !authRes.Success {
			http.Error(w, authRes.Message, http.StatusUnauthorized)
			return
		}
	}

	// User authenticated or guest

	sessionId := ""

	user := auth.User{
		Id:    authRes.Id,
		Name:  authRes.Name,
		Email: credentials.Email,
	}

	if !nosession {

		ip := r.RemoteAddr
		status := "idle"

		sess := auth.Session{
			User:       user,
			Agent:      credentials.Agent,
			Platform:   credentials.Platform,
			Model:      credentials.Model,
			DeviceName: credentials.DeviceName,
			DeviceType: credentials.DeviceType,
			Ip:         ip,
			Status:     status,
			Updated:    time.Now().UTC(),
		}

		sessionId, err = auth.Create("ek", sess)

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

	token, err := auth.GenerateJWT(sessionId, authRes.Id, credentials.Email, authRes.Name, authRes.Roles, authRes.Privileges, expiresAt)

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
	w.Write([]byte(fmt.Sprintf(`{"token":"%s", "name":"%s", "id":"%s", "hostname":"%s" }`, token, authRes.Name, authRes.Id, hostname)))

	//fmt.Println(token)

	if isGuest {
		log.Printf("Guest %s entered\n", user.Name)
	} else {
		log.Printf("User %s successfully authenticated\n", user.Name)
	}
}
