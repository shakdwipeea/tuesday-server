package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/redis.v5"
	"strings"
)

var client *redis.Client

type TuesIDResponse struct {
	TuesId string `json:"tues_id"`
}

type HTTPResponse struct {
	Message string `json:"message"`
}

func newRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	return client, err
}

func handleNewTuesId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tuesID := GetNextSeq(client)
	if tuesID == "" {
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
	var reqBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&HTTPResponse{
			Message: "Bad request",
		})
		return
	}

	name := reqBody["name"]
	if name == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&HTTPResponse{
			Message: "Bad request",
		})
		return
	}

	err = savePrefixes(client, name)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(&HTTPResponse{
			Message: "Could not index",
		})
		return
	}

	json.NewEncoder(w).Encode(&HTTPResponse{
		Message: "Successfully indexed your name",
	})
}

func handleSearch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	type SearchResponse struct {
		Results []string `json:"results"`
	}

	prefix := strings.ToLower(r.URL.Query().Get("key"))

	names, err := client.SMembers(prefix).Result()
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(&HTTPResponse{
			Message: "Could not index",
		})
		return
	}

	json.NewEncoder(w).Encode(&SearchResponse{
		Results: names,
	})
}

func main() {
	var err error

	// Open database
	client, err = newRedisClient()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = SetupDB(client)
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = GenCombination(client)
	if err != nil {
		log.Fatalln(err.Error())
	}

	router := httprouter.New()

	router.GET("/tuesid", handleNewTuesId)
	router.POST("/register", handleNewUser)
	router.GET("/search", handleSearch);

	log.Println(http.ListenAndServe(":9090", router))
}
