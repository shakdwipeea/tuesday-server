package main

import (
	"encoding/json"
	"log"
	"net/http"

	"database/sql"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

var db *sql.DB

type TuesIDResponse struct {
	TuesId string `json:"tues_id"`
}

type HTTPResponse struct {
	Message string `json:"message"`
}

func handleNewTuesId(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	log.Println("Incoming request for new tues id")

	tuesID, err := GetNextSeq(db)
	if err != nil || tuesID == "" {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(&HTTPResponse{
			Message: "No tuesid available",
		})
		return
	}

	log.Println("New tuesid generated " + tuesID)
	json.NewEncoder(w).Encode(&TuesIDResponse{
		TuesId: tuesID,
	})
}

func handleNewUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("Incoming request for new user registration")

	var reqBody User
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&HTTPResponse{
			Message: "Bad request",
		})
		return
	}

	err = saveUser(db, reqBody)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(&HTTPResponse{
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(&HTTPResponse{
		Message: "Successfully indexed your name",
	})
}

func handleSearch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("Incoming request for search")

	type SearchResponse struct {
		Results []string `json:"results"`
	}

	prefix := strings.ToLower(r.URL.Query().Get("key"))

	users, err := getUsers(db, prefix)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(&HTTPResponse{
			Message: err.Error(),
		})
		return
	}

	log.Println("Response sent")
	json.NewEncoder(w).Encode(&users)
}

func main() {
	var err error

	// Connect to mysql
	db, err = sql.Open("mysql", "root:morning_star@/tuesday")
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Create schema
	err = createSchema(db)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Schema created")

	err = GenCombination(db)
	if err != nil && err != keyPresentError {
		log.Fatalln(err.Error())
	}
	log.Println("Keys generated successfully.")

	router := httprouter.New()

	router.GET("/tuesid", handleNewTuesId)
	router.POST("/register", handleNewUser)
	router.GET("/search", handleSearch)

	err = http.ListenAndServe(":9090", router)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
