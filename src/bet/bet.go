package bet

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"src/database"
	"src/match"
	"src/user"
	"src/utils"
	"strconv"
	"time"
)

var genP int = 0

type Bet struct {
	Id             int       `json:"id"`
	IdMatch        int       `json:"idMatch"`
	EquipeGagnante string    `json:"equipeGagnante"`
	Cote           float32   `json:"cote"`
	Montant        float32   `json:"montant"`
	Login          string    `json:"login"`
	Resultat       string    `json:"resultat"`
	Date           time.Time `json:"date"`
}

func GetBet(w http.ResponseWriter, r *http.Request) {
	idSession := r.FormValue("idSession")
	statutParis := r.FormValue("statutParis")

	fmt.Printf("La valeur de statusParis est %v\n", statutParis)

	login := utils.IsConnectedIdSession(idSession)
	if login == "" {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "<bet.go> : user not connected"}`)
		return
	}

	db := database.Connect()
	if db == nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"}`)
		return
	}

	var res *sql.Rows
	var err error
	if statutParis == "coming" {
		res, err = db.Query("Select * From `projet-pc3r`.`Pari` where login=? and resultat='coming';", login)
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"}`)
			return
		}
	} else {
		res, err = db.Query("Select * From `projet-pc3r`.`Pari` where login=? and resultat<>'coming';", login)
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"}`)
			return
		}
	}

	err = db.Close()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"}`)
		return
	}
	resultat := make([]Bet, 0)
	for res.Next() {
		b := Bet{}

		err := res.Scan(&b.IdMatch, &b.EquipeGagnante, &b.Id, &b.Cote, &b.Montant, &b.Login, &b.Resultat, &b.Date)
		fmt.Printf("id=%v, idMatch=%v, winner=%v, cote=%v, montant=%v, login=%v, resultat=%v, date=%v", b.IdMatch, b.Id, b.EquipeGagnante, b.Cote, b.Montant, b.Login, b.Resultat, b.Date)
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem reading result request"}`)
			return
		}
		resultat = append(resultat, b)
	}
	resultJSON, err := json.Marshal(resultat)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem creation of JSON"}`)
		return
	}
	utils.SendResponse(w, http.StatusOK, `{"message": "request effected", "result":`+string(resultJSON)+"}")
}

func AddBet(w http.ResponseWriter, r *http.Request) {
	idSession := r.FormValue("idSession")
	idMatch := r.FormValue("idMatch")
	equipeGagnante := r.FormValue("equipeGagnante")
	coteStr := r.FormValue("cote")
	montantStr := r.FormValue("montant")

	fmt.Printf("Cote : %v, Montant : %v, Vainqueur : %v, idMatch : %v, idSession : %v\n", coteStr, montantStr, equipeGagnante, idMatch, idSession)

	cote, err := strconv.ParseFloat(coteStr, 32)

	if err != nil {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "wrong value for cote"}`)
	}
	montant64, err := strconv.ParseFloat(montantStr, 32)
	if err != nil {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "wrong value for montant"}`)
	}

	montant := float32(montant64)

	login := utils.IsConnectedIdSession(idSession)

	if login == "" {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "<Mbet.go> : user not connected"}`)
		return
	}

	montantCompte := user.GetAccountMoney(login)

	if montantCompte < montant {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "not enough coin"}`)
		return
	}

	testInsert := addBetSql(idMatch, equipeGagnante, float32(cote), montant, login)

	if !testInsert {
		utils.SendResponse(w, http.StatusInternalServerError, `{"<addBet> message": "problem with database"}`)
	} else {
		user.AlterMoney(login, -montant)
		utils.SendResponse(w, http.StatusOK, `{"message":"New bet created"}`)
	}

}

func DeleteBet(w http.ResponseWriter, r *http.Request) {
	idPari := r.FormValue("idPari")
	idSession := r.FormValue("idSession")
	login := utils.IsConnectedIdSession(idSession)

	if login == "" {
		utils.SendResponse(w, http.StatusForbidden, `{"message": "user not connected"}`)
		return
	}

	testInsert := removeBetSQL(idPari, login)

	if !testInsert {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with database"}`)
	} else {
		utils.SendResponse(w, http.StatusOK, `{"message":"New user created"}`)
	}
}

func removeBetSQL(pari, login string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	res, err := db.Exec("Delete from `projet-pc3r`.Pari where `projet-pc3r`.Pari.id=? and `projet-pc3r`.Pari.login=?;", pari, login)
	if err != nil {
		return false
	}
	test, err := res.RowsAffected()
	if err != nil || test != 1 {
		return false
	}
	return true
}

func addBetSql(idMatch, equipeGagnante string, cote float32, montant float32, login string) bool {
	db := database.Connect()
	if db == nil {
		return false
	}
	res, err := db.Exec("Insert into `projet-pc3r`.Pari(`projet-pc3r`.Pari.idmatch, `projet-pc3r`.Pari.equipegagnante, `projet-pc3r`.Pari.cote, `projet-pc3r`.Pari.montant, `projet-pc3r`.Pari.login, `projet-pc3r`.Pari.date, `projet-pc3r`.Pari.id) Values(?, ?, ?, ?, ?, ?, ?) ;", idMatch, equipeGagnante, cote, montant, login, time.Now(), genP)
	genP++
	if err != nil {
		return false
	}
	test, err := res.RowsAffected()
	if err != nil || test != 1 {
		return false
	}
	return true
}

func UpdateResult1Hour() {
	db := database.Connect()
	if db == nil {
		panic(errors.New("problem database connection"))
	}
	res, err := db.Query("Select * from  Pari as P where resultat='coming' and EXISTS( Select * From `Matchs` where id=P.idMatch and statut='finished');")
	if err != nil {
		panic(err.Error())
	}
	for res.Next() {
		b := Bet{}
		err := res.Scan(&b.Id, &b.IdMatch, &b.EquipeGagnante, &b.Cote, &b.Montant, &b.Login, &b.Resultat, &b.Date)

		if err != nil {
			return
		}

		var r sql.Result
		win := match.WinnerIdMatch(b.IdMatch)
		if b.EquipeGagnante == win {
			user.AlterMoney(b.Login, b.Montant*b.Cote)
			r, err = db.Exec("Update Pari set resultat='win' where id=?", b.Id)
		} else {
			r, err = db.Exec("Update Pari set resultat='loose' where id=?", b.Id)
		}
		if err != nil {
			return
		}
		res, err := r.RowsAffected()
		if res != 1 {
			return
		}
	}

}
