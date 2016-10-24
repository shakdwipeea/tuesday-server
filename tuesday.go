package main

import (
	"log"
	"gopkg.in/redis.v5"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
)

var client *redis.Client

type HTTPResponse struct {
	TuesId string `json:"tues_id"`
}

type HTTPErrorResponse struct {
	Message string `json:"message"`
}

func newRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	_, err := client.Ping().Result()
	return client, err
}

func handleNewTuesId(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tuesID := GetNextSeq(client)
	if tuesID == "" {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(&HTTPErrorResponse{
			Message: "No tuesid available",
		})
	}

	json.NewEncoder(w).Encode(&HTTPResponse{
		TuesId: tuesID,
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

	log.Println(http.ListenAndServe(":9090", router))
}
