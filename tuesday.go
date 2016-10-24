package tuesday

import (
	"log"
	"gopkg.in/redis.v5"
	"strconv"
	"strings"
	"errors"
)


func newRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

	_, err := client.Ping().Result()
	return client, err
}

func nextChar(currentChar string) (string, error) {
	chars := "abcdefghijklmnopqrstuvwxyz"

	if currentChar == "z" {
		return "a", nil
	} else {
		index := strings.Index(chars, currentChar)
		if index == -1 {
			return "", errors.New("Not an alphabet")
		}

		return string(chars[index + 1]), nil
	}
}

func nextInt(currentInt string) (string, error) {
	num, err := strconv.Atoi(currentInt)
	if err != nil {
		return "", err
	}

	if num == 9 {
		return "0", nil
	} else {
		return strconv.Itoa(num + 1), nil
	}
}

func main() {
	// Open database
	client, err := newRedisClient()
	if err != nil {
		log.Fatalln(err.Error())
	}

	err = SetupDB(client)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
