package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
	firebase "github.com/wuman/firebase-server-sdk-go"
)

var db *sql.DB

type TuesIDResponse struct {
	TuesId string `json:"tues_id"`
}

type HTTPResponse struct {
	Message string `json:"message"`
}

func handlePhone(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//user signup information
	var reqBody User
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		SendErrorResponse(400, err.Error(), w)
		return
	}

	user, err := getUser(db, reqBody.Phone)
	if err == sql.ErrNoRows {
		// user does not exist so create the user and send otp
		uid, err := createUser(reqBody)
		if err != nil {
			SendErrorResponse(500, err.Error(), w)
			return
		}

		user.Uid = uid
		user.Verified = false
	}

	if err != nil && err != sql.ErrNoRows {
		SendErrorResponse(500, err.Error(), w)
		return
	}

	token, err := createSignInToken(strconv.Itoa(user.Uid))
	if err != nil {
		SendErrorResponse(500, err.Error(), w)
		return
	}

	user.Token = token
	user.Otp = ""
	json.NewEncoder(w).Encode(&user)
}

func createSignInToken(uid string) (string, error) {
	log.Println("uid is " + uid)
	auth, _ := firebase.GetAuth()
	token, err := auth.CreateCustomToken(uid, nil)
	return token, err
}

func createUser(user User) (int, error) {
	otp, err := genOtp()
	if err != nil {
		return -1, err
	}

	user.Otp = otp
	user.Verified = false

	id, err := saveUser(db, user)
	if err != nil {
		return -1, err
	}

	go sendOtp(user.Otp, user.Phone)
	return id, nil
}

func handleOtpVerification(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var reqBody User
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&HTTPResponse{
			Message: "Bad request",
		})
		return
	}

	user, err := getUser(db, reqBody.Phone)
	if err != nil {
		SendErrorResponse(500, err.Error(), w)
		return
	}

	if user.Otp != reqBody.Otp {
		SendErrorResponse(400, "Otp incorrect", w)
		return
	}

	user.Verified = true
	err = verifyUser(db, user.Phone, user.Verified)
	if err != nil {
		SendErrorResponse(500, err.Error(), w)
		return
	}

	token, err := createSignInToken(strconv.Itoa(user.Uid))
	if err != nil {
		SendErrorResponse(500, err.Error(), w)
		return
	}

	user.Token = token
	json.NewEncoder(w).Encode(&user)
}

func handleProfile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	phone := r.URL.Query().Get("phone")

	user, err := getUser(db, phone)
	if err != nil {
		SendErrorResponse(500, err.Error(), w)
		return
	}

	user.Otp = ""

	json.NewEncoder(w).Encode(&user)
}

// Sign up route
func handleUpdateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	log.Println("Incoming request for user update")

	var reqBody User
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(&HTTPResponse{
			Message: "Bad request",
		})
		return
	}

	err = updateUser(db, reqBody)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(&HTTPResponse{
			Message: err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(&HTTPResponse{
		Message: "Details saved",
	})
}

func main() {
	var err error

	// Connect to mysql
	db, err = sql.Open("mysql", "root:@/tuesday")
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Create schema
	err = createSchema(db)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("Schema created")

	firebase.InitializeApp(&firebase.Options{
		ServiceAccountPath: "firebase-credentials.json",
	})

	router := httprouter.New()

	router.POST("/phone", handlePhone)
	router.POST("/register", handleUpdateUser)
	router.POST("/verify", handleOtpVerification)
	router.GET("/profile", handleProfile)

	err = http.ListenAndServe(":9090", router)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
