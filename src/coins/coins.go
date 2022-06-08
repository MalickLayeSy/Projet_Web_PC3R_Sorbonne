package coins

import (
	"fmt"
	"net/http"
	"src/database"
	"src/utils"
)

func Generate(w http.ResponseWriter, r *http.Request) {
	idSession := r.FormValue("idSession")
	montant := r.FormValue("montant")

	login := utils.IsConnectedIdSession(idSession)
	if login == "" {
		fmt.Printf("L'id de la session est le suivant :%v\n", idSession)
		utils.SendResponse(w, http.StatusForbidden, `{"message": "<coins.go> (1) user not connected"}`)
		return
	}

	db := database.Connect()
	if db == nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with connection database"}`)
		return
	}

	res, err := db.Exec("Update Utilisateur SET cagnotte = cagnotte + ? where login=? ;", montant, login)

	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with connection database"}`)
		return
	}

	row, err := res.RowsAffected()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with connection database"}`)
		return
	}

	if row != 1 {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with request"}`)
		return
	}
	utils.SendResponse(w, http.StatusOK, `{"message": "user has now more coins"}`)

}
