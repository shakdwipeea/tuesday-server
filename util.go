package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func SendErrorResponse(statusCode int, message string, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(&HTTPResponse{
		Message: message,
	})
}

func genOtp() (string, error) {
	len := 8

	b := make([]byte, len)

	_, err := io.ReadFull(rand.Reader, b)
	return base64.StdEncoding.EncodeToString(b)[:6], err
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
