package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"src/bet"
	"src/coins"
	"src/connexion"
	"src/match"
	"src/user"
	"src/utils"
	"time"
)

func handleUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	log.Printf("User \tip %s\n", r.RemoteAddr)
	if err != nil {
		handleProblem(w, r)
	}
	switch r.Method {
	case "GET":
		user.GetUser(w, r)
	case "POST":
		user.AddUser(w, r)
	case "DELETE":
		user.DeleteUser(w, r)
	default:
		handleProblem(w, r)
	}
}

func handleBet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	log.Printf("Bet \tip %s\n", r.RemoteAddr)
	if err != nil {
		handleProblem(w, r)
	}
	switch r.Method {
	case "GET":
		bet.GetBet(w, r)
	case "POST":
		bet.AddBet(w, r)
	case "DELETE":
		bet.DeleteBet(w, r)
	default:
		handleProblem(w, r)
	}
}

func handleMatch(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	log.Printf("Match \tip %s\n", r.RemoteAddr)
	if err != nil {
		handleProblem(w, r)
	}
	switch r.Method {
	case "GET":
		match.GetMatch(w, r)
	default:
		handleProblem(w, r)
	}
}

func handleConnexion(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	log.Printf("Connexion \tip %s\n", r.RemoteAddr)
	if err != nil {
		handleProblem(w, r)
	}
	switch r.Method {
	case "POST":
		connexion.Connect(w, r)
	case "DELETE":
		connexion.Disconnect(w, r)
	default:
		handleProblem(w, r)
	}
}

func handleCoins(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	log.Printf("Coins \tip %s\n", r.RemoteAddr)
	if err != nil {
		handleProblem(w, r)
	}
	switch r.Method {
	case "POST":
		coins.Generate(w, r)
	default:
		handleProblem(w, r)
	}
}

//HandleHome
func handleHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("Welcome\n")
	//http.Redirect(w, r, "http://projet-pc3r.eba-d6ekfsap.eu-west-3.elasticbeanstalk.com/home/", 301)
	utils.SendResponse(w, http.StatusOK, `{"message":"hello world!"}`)
}

func handleProblem(w http.ResponseWriter, _ *http.Request) {
	log.Printf("A problem appear\n")
	utils.SendResponse(w, http.StatusInternalServerError, `{"message":"problem"}`)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "5000"
		http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("web"))))
		//http.HandleFunc("/", handleHome)
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

			http.Redirect(w, r, "/home", http.StatusSeeOther)
		})
	}

	updateComingMatches()
	updateResultMatchesAndBet()

	http.HandleFunc("/user", handleUser)
	http.HandleFunc("/bet", handleBet)
	http.HandleFunc("/match", handleMatch)
	http.HandleFunc("/connexion", handleConnexion)
	http.HandleFunc("/coins", handleCoins)

	log.Printf("Listening on port %s\n\n", port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		panic(err.Error())
	}

	/*
		err = f.Close()
		if err != nil {
			panic(err.Error())
		}*/
}

//list files directory
func listFiles() {
	var files []string

	root := "."
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		fmt.Println(file)
	}
}

func updateComingMatches() {
	ticker := time.NewTicker(24 * time.Hour)
	match.LoadComingMatchFor2Week()
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				match.LoadComingMatchFor2Week()
			}
		}
	}()
}

func updateResultMatchesAndBet() {
	ticker := time.NewTicker(1 * time.Hour)
	match.LoadResultMatchFor3Hours()
	bet.UpdateResult1Hour()
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				match.LoadResultMatchFor3Hours()
				bet.UpdateResult1Hour()
			}
		}
	}()
}
