package match

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"src/database"
	"src/utils"
	"strconv"
	"strings"
	"time"
)

var gen int

type Match struct {
	Id        int       `json:"id"`
	Sport     string    `json:"sport"`
	League    string    `json:"league"`
	EquipeA   string    `json:"equipeA"`
	EquipeB   string    `json:"equipeB"`
	Cote      float32   `json:"cote"`
	Statut    string    `json:"statut"`
	Vainqueur string    `json:"vainqueur"`
	Date      time.Time `json:"date"`
}

func GetMatch(w http.ResponseWriter, r *http.Request) {
	req := r.FormValue("req")
	fmt.Printf("La requete est la suivante : %v\n\n", req)
	idSession := r.FormValue("idSession")
	fmt.Printf("Utilisateur connecté : %v\n\n", idSession)

	login := utils.IsConnectedIdSession(idSession)
	if login == "" {
		fmt.Printf("L'id de la session est le suivant :%v\n", idSession)
		utils.SendResponse(w, http.StatusForbidden, `{"message": "<Match> : user not connected"}`)
		return
	}

	db := database.Connect()
	if db == nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "<Match> : problem with connection to database"}`)
		return
	}

	var res *sql.Rows
	var err error
	if req == "" {
		res, err = db.Query("Select * From `Matchs`;")
		//res, err = db.Query("Select * From `Matchs` where statut='not_started' and equipeA<>'' and equipeB<>'' order by date DESC ;")
	} else {
		res, err = db.Query("Select * From `Matchs` where (sport=? or league=? or equipeA=? or equipeB=?) order by date DESC;", req, req, req, req)
	}

	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with searching database"}`)
		return
	}
	err = db.Close()
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem with closing database"}`)
		return
	}
	resultat := make([]Match, 0)
	for res.Next() {
		m := Match{}
		var date string
		err := res.Scan(&m.Sport, &m.League, &m.EquipeA, &m.EquipeB, &date, &m.Id, &m.Cote, &m.Statut, &m.Vainqueur)
		m.Date, _ = time.Parse("2006-01-02 15:04:05", date)
		fmt.Printf("Affichage du match %v, sport:%v, ligue:%v, equipe A:%v, equipe B:%v, date:%v, cote:%v, statut:%v, vainqueur:%v, date:%v", m.Id, m.Sport, m.League, m.EquipeA, m.EquipeB, m.Date, m.Cote, m.Statut, m.Vainqueur, m.Date)
		if err != nil {
			utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem reading result request"}`)
			return
		}
		resultat = append(resultat, m)
		fmt.Printf("%v\n", resultat)
	}
	resultJSON, err := json.Marshal(resultat)
	if err != nil {
		utils.SendResponse(w, http.StatusInternalServerError, `{"message": "problem creation of JSON"}`)
		return
	}
	utils.SendResponse(w, http.StatusOK, `{"message": "coming matches", "result":`+string(resultJSON)+"}")

}

//Ne pas appeler : LoadAllPastMatch
func _() {

	req := "https://api.pandascore.co/fifa/matches/past?token=x86LDA2MmKcz_PDNNXdKkBiT04kocdn_AYk_XJk1ckH7vBKBAhE"

	resp, _ := http.Get(req + "&page[size]=100")
	JSONMatch2SQL(resp)

	test := resp.Header.Get("Link")
	res := strings.Split(test, ",")
	last := ""
	for _, v := range res {
		if strings.Contains(v, "last") {
			last = strings.Split(v, ";")[0][2 : len(strings.Split(v, ";")[0])-1]
		}
	}

	u, err := url.Parse(last)
	if err != nil {
		panic(err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		panic(err)
	}
	max, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		panic(err)
	}
	fmt.Println(max)
	for i := 2; i < max+1; i++ {
		s := req + "&page[size]=100&page[number]=" + strconv.Itoa(i)
		//fmt.Println(s)
		resp, _ := http.Get(s)
		JSONMatch2SQL(resp)
	}
}

func LoadComingMatchFor2Week() {

	req := "https://api.pandascore.co/fifa/matches/upcoming?token=x86LDA2MmKcz_PDNNXdKkBiT04kocdn_AYk_XJk1ckH7vBKBAhE"
	t := time.Now()
	req += "&range[begin_at]=" + strings.Split(t.Format("2006-01-02T15:04:05-0700"), "+")[0] + "," + strings.Split(t.Add(time.Hour*24*7*2).Format("2006-01-02T15:04:05-0700"), "+")[0]
	s := req + "&page[size]=100"
	//fmt.Println(s)
	resp, _ := http.Get(s)
	JSONMatch2SQL(resp)

	test := resp.Header.Get("Link")
	res := strings.Split(test, ",")
	last := ""
	for _, v := range res {
		if strings.Contains(v, "last") {
			last = strings.Split(v, ";")[0][2 : len(strings.Split(v, ";")[0])-1]
		}
	}

	u, err := url.Parse(last)
	if err != nil {
		panic(err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		panic(err)
	}
	max, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		max = 0
	}

	for i := 2; i < max+1; i++ {
		s := req + "&page[size]=100&page[number]=" + strconv.Itoa(i)
		//fmt.Println(s)
		resp, _ := http.Get(s)
		go JSONMatch2SQL(resp)
	}
}

func JSONMatch2SQL(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	var data utils.MatchJSON // TopTracks
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	addMulipleMatch(data)
}

func addMulipleMatch(data utils.MatchJSON) {
	for _, v := range data {
		//time.Sleep(150*time.Millisecond)
		//fmt.Println(v)
		if len(v.Opponents) == 2 {
			addMatch(v.Videogame.Name, v.League.Name, v.Opponents[0].Opponent.Name, v.Opponents[1].Opponent.Name, v.Status, v.Winner.Name, v.OriginalScheduledAt)
		} else {
			fmt.Println("<addMultiplematch> le nombre d'adversaire n'est pas égal à deux.")
			addMatch(v.Videogame.Name, v.League.Name, "", "", v.Status, "", v.OriginalScheduledAt)
		}
	}
}

func addMatch(sport string, league string, equipeA string, equipeB string, statut string, winner string, date time.Time) {
	//cote := calculCote(equipeA, equipeA)
	cote := 1
	db := database.Connect()
	fmt.Printf("<AddMatch> Update `Matchs` set equipeA=%v , equipeB=%v , vainqueur=%v , statut=%v where sport=%v and league=%v and equipeA='' and equipeB='' and date=%v ;\n", equipeA, equipeB, winner, statut, sport, league, date)
	r, err := db.Exec("Update `Matchs` set equipeA=? , equipeB=? , vainqueur=? , statut=? , cote=? where sport=? and league=? and equipeA='' and equipeB='' and date=? ;", equipeA, equipeB, winner, statut, cote, sport, league, date)
	//fmt.Println(err)
	if err == nil {
		nbRows, err2 := r.RowsAffected()
		if err2 != nil || nbRows != 1 {
			//fmt.Println(err.Error())
			_, err := db.Exec("Insert into `Matchs` (sport, league, equipeA, equipeB, id, cote,statut, vainqueur, date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);", sport, league, equipeA, equipeB, gen, cote, statut, winner, date)
			gen++
			if err != nil {

				if !strings.Contains(err.Error(), "Duplicate") {
					panic(err.Error())
				}

			}
		}
	}
	err = db.Close()
	if err != nil {
		panic(err)
	}
}

func LoadResultMatchFor3Hours() {
	req := "https://api.pandascore.co/fifa/matches/past?token=x86LDA2MmKcz_PDNNXdKkBiT04kocdn_AYk_XJk1ckH7vBKBAhE"
	t := time.Now()
	req += "&range[end_at]=" + strings.Split(t.Add(-3*time.Hour).Format("2006-01-02T15:04:05-0700"), "+")[0] + "," + strings.Split(t.Format("2006-01-02T15:04:05-0700"), "+")[0]
	s := req + "&page[size]=100"
	//fmt.Println(s)
	resp, _ := http.Get(s)
	JSONMatchUpdate(resp)

	test := resp.Header.Get("Link")
	res := strings.Split(test, ",")
	last := ""
	for _, v := range res {
		if strings.Contains(v, "last") {
			last = strings.Split(v, ";")[0][2 : len(strings.Split(v, ";")[0])-1]
		}
	}

	u, err := url.Parse(last)
	if err != nil {
		panic(err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		panic(err)
	}
	max, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		max = 0
	}
	for i := 2; i < max+1; i++ {
		s := req + "&page[size]=100&page[number]=" + strconv.Itoa(i)
		//fmt.Println(s)
		resp, _ := http.Get(s)
		JSONMatchUpdate(resp)
	}

}

func JSONMatchUpdate(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	var data utils.MatchJSON // TopTracks
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}
	updateMulipleMatch(data)
}

func updateMulipleMatch(data utils.MatchJSON) {
	fmt.Println("Mis a jour des matchs")
	for _, v := range data {
		if len(v.Opponents) == 2 {
			fmt.Printf("<updateMultupleMatch> Sport :%v , League :%v ,  OppentsA : %v, OpponentsB : %v, vainqueur : %v, date : %v, statut : %v\n", v.Videogame.Name, v.League.Name, v.Opponents[0].Opponent.Name, v.Opponents[1].Opponent.Name, v.Winner.Name, v.ScheduledAt, v.Status)
			updateMatch(v.Videogame.Name, v.League.Name, v.Opponents[0].Opponent.Name, v.Opponents[1].Opponent.Name, v.Winner.Name, v.OriginalScheduledAt, v.Status)
		}
	}
}

func updateMatch(sport string, league string, equipeA string, equipeB string, winner string, date time.Time, statut string) {
	db := database.Connect()

	_, err := db.Exec("Update `projet-pc3r`.`Matchs` SET `vainqueur`=? , `statut`=? where sport=? and league=? and equipeA=? and equipeB=? and `date`=? and statut='not_started';", winner, statut, sport, league, equipeA, equipeB, date)
	if err != nil {
		//fmt.Println(err.Error())
		panic(err.Error())
	}
	err = db.Close()
	if err != nil {
		return
	}
}

func WinnerIdMatch(idMatch int) string {
	db := database.Connect()
	if db == nil {
		return ""
	}
	m := Match{}
	err := db.QueryRow("Select * From `Matchs` where id=?;", idMatch).Scan(&m.Id, &m.Sport, &m.League, &m.EquipeA, &m.EquipeB, &m.Cote, &m.Statut, &m.Vainqueur, &m.Date)
	if err != nil {
		panic(err.Error())
	}
	err = db.Close()
	if err != nil {
		return ""
	}
	return m.Vainqueur
}

//Calcul Cote
func calculCote(equipeA string, equipeB string) float32 {
	if equipeA == "" || equipeB == "" {
		return 1
	}
	nbMatchTotal := nbMatchTotal(equipeA, equipeB)
	if nbMatchTotal == -1 {
		return 1
	}
	nbMatchGagneA := nbMatchGagne(equipeA, equipeB)
	if nbMatchGagneA == -1 {
		return 1
	}
	nbMatchGagne5DerniersA := nbMatchGagne5Derniers(equipeA, equipeB)
	if nbMatchGagne5DerniersA == -1 {
		return 1
	}

	pourcentageVictoireTotale := (float32(nbMatchGagneA)/float32(nbMatchTotal))/2 + (float32(nbMatchGagne5DerniersA)/5)/2
	return 100 / pourcentageVictoireTotale
}

func nbMatchGagne(equipeA string, equipeB string) int {
	db := database.Connect()
	if db == nil {
		return -1
	}
	res := 0
	err := db.QueryRow("Select Count(*) From  `projet-pc3r`.`Matchs` where  (equipeA=? or equipeA=?) and (equipeB=? or equipeB=?) and vainqueur=?;", equipeA, equipeB, equipeA, equipeB, equipeA).Scan(&res)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	err = db.Close()
	if err != nil {
		return -1
	}

	return res
}

func nbMatchTotal(equipeA string, equipeB string) int {
	db := database.Connect()
	if db == nil {
		return -1
	}
	res := 0
	err := db.QueryRow("Select Count(*) From  `projet-pc3r`.`Matchs` where (equipeA=? or equipeA=?) and (equipeB=? or equipeB=?);", equipeA, equipeB, equipeA, equipeB).Scan(&res)
	if err != nil {
		//fmt.Println(err.Error())
		panic(err.Error())
	}
	err = db.Close()
	if err != nil {
		return -1
	}
	return res
}

func nbMatchGagne5Derniers(equipeA string, equipeB string) int {
	db := database.Connect()
	if db == nil {
		return -1
	}
	res := 0
	err := db.QueryRow("Select Count(*) From (Select * From `Matchs` where (equipeA=? or equipeA=?) and (equipeB=? or equipeB=?) order by date DESC LIMIT 5 ) as `M*` where vainqueur=?;", equipeA, equipeB, equipeA, equipeB, equipeA).Scan(&res)
	if err != nil {
		//fmt.Println(err.Error())
		panic(err.Error())
	}
	err = db.Close()
	if err != nil {
		return -1
	}

	return res
}
