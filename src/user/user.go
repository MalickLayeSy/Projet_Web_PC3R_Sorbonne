package user

import (
	"fmt"
	"net/http"
	"src/database"
	"src/utils"
	"time"
)

type User struct {
	id                int
	login             string
	password          string
	mail              string
	cagnotte          float32
	derniereConnexion time.Time
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	//Recuperation des parametres de la requete http
	idSession := req.FormValue("idSession")

	login := utils.IsConnectedIdSession(idSession)
	//verif connexion
	if login == "" {
		utils.SendResponse(res, http.StatusForbidden, `{"message":"user was not connected"}`)
		return
	}

	var user = searchUser(login)
	if user == nil {
		utils.SendResponse(res, http.StatusInternalServerError, `{"message":"a problem appeared"}`)
		return
	}

	if user.login != "" {
		cagnotte := fmt.Sprintf("%f", user.cagnotte)
		utils.SendResponse(res, http.StatusOK, `{"message":"user found", "login":"`+user.login+`", "mail":"`+user.mail+`", "cagnotte":"`+cagnotte+`"}`)
	} else {
		utils.SendResponse(res, http.StatusForbidden, `{"message":"problem login user don't exist"}`)
	}

}

func AddUser(res http.ResponseWriter, req *http.Request) {
	mail := req.FormValue("mail")
	login := req.FormValue("login")
	password := req.FormValue("password")
	fmt.Printf("login : %v, password : %v mail : %v\n", login, password, mail)
	if acceptLogin(login) {
		if insertUser(login, password, mail) {
			utils.SendResponse(res, http.StatusOK, `{"message":"New user created"}`)
		} else {
			utils.SendResponse(res, http.StatusInternalServerError, `{"message":"A problem appeared"}`)
		}
	} else {
		utils.SendResponse(res, http.StatusForbidden, `{"message":"Problem login already exist"}`)
	}
}

func DeleteUser(res http.ResponseWriter, req *http.Request) {
	login := req.FormValue("login")
	password := req.FormValue("password")
	idSession := req.FormValue("idSession")
	if utils.IsConnectedIdSession(idSession) != "" && idSession != "" {
		if utils.IsUser(login, password) {
			if removeUser(login, password) {
				utils.SendResponse(res, http.StatusOK, `{"message":"deleted user"}`)
			} else {
				utils.SendResponse(res, http.StatusInternalServerError, `{"message":"problem with database"}`)
			}
		} else {
			utils.SendResponse(res, http.StatusForbidden, `{"message":"Error : wrong login or password"}`)
		}
	}
}

func acceptLogin(login string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	var count int
	err := db.QueryRow("Select count(*) From Utilisateur where login=?;", login).Scan(&count)
	if err != nil {
		return false
	}
	err = db.Close()
	if err != nil {
		return false
	}
	if count == 0 {
		return true
	}
	return false
}

func searchUser(login string) *User {
	db := database.Connect()
	if db == nil {
		return nil
	}

	u := User{}

	err := db.QueryRow("Select login, mail, cagnotte From Utilisateur where login=?;", login).Scan(&u.login, &u.mail, &u.cagnotte)
	if err != nil {
		return nil
	}
	err = db.Close()
	if err != nil {
		return nil
	}
	return &u
}

func insertUser(login string, password string, mail string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}

	res, err := db.Exec("INSERT INTO Utilisateur(login, password, mail, cagnotte) VALUES (?, ?, ?, 100);", login, password, mail)
	err = db.Close()
	if err != nil {
		return false
	}

	r, err := res.RowsAffected()
	if r == 0 || err != nil {
		return false
	} else {
		return true
	}
}

func removeUser(login string, password string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	res, err := db.Exec("Delete from Utilisateur where login=? and password=?;", login, password)
	if err != nil {
		return false
	}
	err = db.Close()
	if err != nil {
		return false
	}

	r, _ := res.RowsAffected()
	if r == 1 {
		return true
	}
	return false
}

func AlterMoney(loginUser string, amount float32) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	r, err := db.Exec("Update Utilisateur set cagnotte=cagnotte+? where login=?", amount, loginUser)
	if err != nil {
		return false
	}
	row, err := r.RowsAffected()
	if err != nil || row != 1 {
		return false
	}
	return true
}

func GetAccountMoney(login string) float32 {
	u := searchUser(login)
	if u == nil {
		return -1
	}
	return u.cagnotte
}
