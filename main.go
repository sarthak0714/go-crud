package main

import (
	"math/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type Player struct{
	Id string `json:"id"`
	Name string `json:"name"`
	Gun *Gun `json:"gun"`
}

type Gun struct{
	Name string `json:"name"`
	Damage int `json:"damage"`
}

var players []Player

func root(w http.ResponseWriter,r *http.Request){
	log.Printf("Method: %s, URL: %s", r.Method, r.URL)
	io.WriteString(w,"root page")
}

func getPlayers(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	log.Printf("%s %s", r.URL, r.Method)
	json.NewEncoder(w).Encode(players)
}

func getPlayer(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	params := mux.Vars(r)
	for _,player := range players{
		if player.Id == params["id"]{
			log.Printf("Method: %s, URL: %s", r.Method, r.URL)
			json.NewEncoder(w).Encode(player)
			break
		}
	}
}

func updateGun(w http.ResponseWriter, r* http.Request){
	params:=mux.Vars(r)
	w.Header().Set("Content-Type","application/json")
	var player Player
	_ = json.NewDecoder(r.Body).Decode(&player)
	for idx,i:=range players{
		if i.Id == params["id"]{
			players[idx].Gun = player.Gun
			log.Printf("Method: %s, URL: %s", r.Method, r.URL)
			json.NewEncoder(w).Encode(players[idx])
			break
		}
	}
}

func deletePlayer(w http.ResponseWriter, r* http.Request){
	params := mux.Vars(r)
	w.Header().Set("Content-Type","application/json")
	for idx,i := range players{
		if i.Id ==params["id"]{
			players = append(players[:idx],players[idx+1:]...)
			break
		}
	}
	log.Printf("Method: %s, URL: %s", r.Method, r.URL)
	json.NewEncoder(w).Encode(players)
}

func addPlayer(w http.ResponseWriter, r* http.Request){
	var player Player
	_=json.NewDecoder(r.Body).Decode(&player)
	player.Id=strconv.Itoa(rand.Intn(1000))
	players = append(players, player)
	log.Printf("Method: %s, URL: %s", r.Method, r.URL)
	json.NewEncoder(w).Encode(players)
}

func main(){
	r:= mux.NewRouter()
	players = append(players, Player{Id:"123",Name:"Shrihari",Gun:&Gun{Name: "ak47",Damage: 69}})
	players = append(players, Player{Id:"124",Name:"Surya",Gun:&Gun{Name: "uzi",Damage: 39}})
	r.HandleFunc("/",root)
	r.HandleFunc("/players",getPlayers).Methods("GET")
	r.HandleFunc("/player/{id}",getPlayer).Methods("GET")
	r.HandleFunc("/player/{id}",updateGun).Methods("PUT")
	r.HandleFunc("/player/{id}",deletePlayer).Methods("DELETE")
	r.HandleFunc("/players",addPlayer).Methods("POST")
	fmt.Print("Starting server at 8080\n")
	err := http.ListenAndServe(":8080",r)
	if err!=nil {
		log.Fatal("Server active at port 8080\n")
	}
}