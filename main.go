package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Note struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

var noteStore = make(map[string]Note)
var id int

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	var notes []Note
	for _, v := range noteStore {
		notes = append(notes, v)
	}

	w.Header().Set("Content-type", "application/json")
	j, err := json.Marshal(notes)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}
	note.CreatedAt = time.Now()
	id++
	k := strconv.Itoa(id)
	noteStore[k] = note

	w.Header().Set("Content-type", "application/json")
	j, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(j)

}

func PutNoteHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	k := vars["id"]
	var noteUpdate Note
	err := json.NewDecoder(r.Body).Decode(&noteUpdate)
	if err != nil{
		panic(err)
	}
	if note, ok := noteStore[k]; ok{
		noteUpdate.CreatedAt = note.CreatedAt
		delete(noteStore, k)
		noteStore[k] = noteUpdate
	} else{
		log.Printf("No encontramos el id %s", k)
	}
	w.WriteHeader(http.StatusNoContent)

}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	id := vars["id"]
	if _, ok := noteStore[id]; ok{
		delete(noteStore, id)
	}else{
		log.Printf("No encontramos el id %s", id)
	}
	w.WriteHeader(http.StatusNoContent)
}



func main() {

	router := mux.NewRouter().StrictSlash(false)

	router.HandleFunc("/api/notes", GetNoteHandler).Methods("GET")
	router.HandleFunc("/api/notes", PostNoteHandler).Methods("POST")
	router.HandleFunc("/api/notes/{id}", PutNoteHandler).Methods("PUT")
	router.HandleFunc("/api/notes/{id}", DeleteNoteHandler).Methods("DELETE")

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Print("Listening 8080")
	log.Fatal(server.ListenAndServe())

}
