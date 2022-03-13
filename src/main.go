package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Content string `json:"Content"`
	Summary string `json:"Summary"`
}

func topsecret(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var post Article
	json.Unmarshal(reqBody, &post)

	json.NewEncoder(w).Encode(post)

	newData, err := json.Marshal(post)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(string(newData))
	}
}

func topsecretx(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	//if err != nil {
	//	respondWithError(w, http.StatusBadRequest, err.Error())
	//	return
	//}

	//reqBody, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	//}

	fmt.Println("Hello, World!")
	fmt.Println("Vars:", json.NewDecoder(r.Body))
	respondWithJSON(w, http.StatusOK, params)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/topsecret", topsecret).Methods("POST")

	log.Printf("Listening...")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
