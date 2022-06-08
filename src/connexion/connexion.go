package connexion

import (
	"fmt"
	"math/rand"
	"net/http"
	"src/database"
	"src/utils"
	"time"
)

func Connect(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("login : %v, password : %v\n", r.FormValue("login"), r.FormValue("password"))
	login := r.FormValue("login")
	password := r.FormValue("password")
	idSession := utils.IsConnectedLogin(login)
	if idSession != "" {
		utils.SendResponse(w, http.StatusOK, `{"message":"user still connected", "idSession":"`+idSession+`"}`)
		return
	}
	if utils.IsUser(login, password) {
		idSession := addConnexion(login)
		if idSession != "" {
			utils.SendResponse(w, http.StatusOK, `{"message":"user connected", "idSession":"`+idSession+`"}`)
		} else {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message":"problem with database"}`)
		}

	} else {
		utils.SendResponse(w, http.StatusForbidden, `{"message":"unknown user"}`)
	}
}

func Disconnect(w http.ResponseWriter, r *http.Request) {
	idSession := r.FormValue("idSession")

	if utils.IsConnectedIdSession(idSession) != "" && idSession != "" {
		if utils.RemoveConnection(idSession) {
			fmt.Printf("L'utilisateur avec la session %v s'est deconnecté du site", idSession)
			utils.SendResponse(w, http.StatusOK, `{"message":"user disconnected"}`)
		} else {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message":"a problem appeared"}`)
		}
	} else {
		utils.SendResponse(w, http.StatusForbidden, `{"message":"user was not connected"}`)
	}
}

func addConnexion(login string) string {
	db := database.Connect()
	if db == nil {
		return ""
	}
	idSession := randSeq()
	res, err := db.Exec("INSERT INTO Connexion (login, idSession, date) VALUES (?, ?, ?);", login, idSession, time.Now().UTC())
	if err != nil {
		return ""
	}
	err = db.Close()
	if err != nil {
		return ""
	}

	r, _ := res.RowsAffected()
	if r == 0 {
		return ""
	} else {
		return idSession
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")

func randSeq() string {
	b := make([]rune, 16)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
