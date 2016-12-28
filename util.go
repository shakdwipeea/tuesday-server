package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

func SendErrorResponse(statusCode int, message string, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&HTTPResponse{
		Message: message,
	})
}

func genOtp() (string, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	num := 100000 + rand.Intn(899999)
	return strconv.Itoa(num), nil
}

func sendOtp(otp string, phone string) {
	sendMessage(phone, otp+" is your otp")
}

func sendMessage(phone string, message string) {
	params := url.Values{}
	params.Add("authkey", "81123A0Lic9Q63l5505c468")
	params.Add("mobiles", phone)
	params.Add("message", message)
	params.Add("sender", "TUEOTP")
	params.Add("route", "4")

	reqURL := "https://control.msg91.com/api/sendhttp.php?" + params.Encode()

	res, err := http.Get(reqURL)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer res.Body.Close()
	_, err = io.Copy(os.Stdout, res.Body)
	if err != nil {
		log.Println(err.Error())
	}
}
