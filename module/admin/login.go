package admin

import (
	"ekhoes-server/auth"
	"ekhoes-server/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

/**
 * POST /login
 * -H "x-user-agent: Radar/1.0.0" -H "x-platform: Android" -d '{ email: "admin@hal9k.net", password: "admin" }'
 */
func Login(w http.ResponseWriter, r *http.Request) {

	//isGuest := r.URL.Query().Has("guest")
	nosession := r.URL.Query().Has("nosession") // Used by cli

	var credentials auth.Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)

	//fmt.Println(credentials)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authRes := &AuthResult{}

	authRes, err = Authorize(credentials.Email, credentials.Password)

	if err != nil {
		utils.LogErr(thisModule, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !authRes.Success {
		http.Error(w, authRes.Message, http.StatusUnauthorized)
		return
	}

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

	token, err := auth.GenerateJWT(sessionId, authRes.Id, credentials.Email, authRes.Name, authRes.Roles, authRes.Privileges, expiresAt)

	if err != nil {
		log.Println(err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	hostname, _ := os.Hostname()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"token":"%s", "name":"%s", "id":"%s", "hostname":"%s" }`, token, authRes.Name, authRes.Id, hostname)))

	//fmt.Println(token)

	utils.Log(thisModule, "%s successfully authenticated\n", user.Email)
}
