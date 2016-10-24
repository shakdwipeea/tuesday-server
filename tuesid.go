package main

import (
	"gopkg.in/redis.v5"
	"errors"
	"strings"
	"strconv"
)

const (
	TUESID_KEY = "tuesid_key"
)

// SetupDB for Identifier generation
// For details on discussions of identifier generation  refer https://github.com/shakdwipeea/tuesday-android/issues/14
//
// For the first 45k people the unique id wiil be
//	3 pairs of consecutive digits/alphabets
//	Eg AB ST 67
//	   34 GH 12
//	   ZA BC QR
//
// Here, we have to keep track of all the alphabets and numbers appearing in each position.
//
// So the db organization may look something like
// {
//	"numbersGenerated": []
// }
//
//
// make sure no collisions for 45k numbers, then we can remove the check and asynchronously notify user for any errors
// since we have established no chances for error.
func SetupDB(client *redis.Client) error {
	return client.Del(TUESID_KEY).Err()
}

const chars = "abcdefghijklmnopqrstuvwxyz"

func nextChar(currentChar string) (string, error) {
	switch currentChar {
	case "y":
		return "za", nil
	case "z":
		return "ab", nil
	default:
		index := strings.Index(chars, currentChar)
		if index == -1 {
			return "", errors.New("Not an alphabet")
		}

		return string([]byte{
			chars[index + 1],
			chars[index + 2],
		}), nil
	}
}

func nextInt(num int) (string, error) {
	switch num {
	case 8:
		return "90", nil
	case 9:
		return "01", nil
	default:
		return strconv.Itoa(num + 1) + strconv.Itoa(num + 2), nil
	}
}


func genTuesPool() ([]string, error) {
	var pool []string

	for _, char := range chars {
		seq, err := nextChar(string(char))
		if err != nil {
			return pool, err
		}
		pool = append(pool, seq)
	}

	for i := 0; i < 10; i++ {
		seq, err := nextInt(i)
		if err != nil {
			return pool, err
		}
		pool = append(pool, seq)
	}

	return pool, nil
}

func GenCombination(client *redis.Client) error {
	pool, err := genTuesPool()
	if err != nil {
		return err
	}

	for _, i := range pool {
		for _, j := range pool {
			for _, k := range pool {
				err = SaveTuesId(client, i + j + k)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func GetNextSeq(client *redis.Client) string {
	return client.SPop(TUESID_KEY).Val()
}

func SaveTuesId(client *redis.Client, tuesID string) error {
	if checkPresence(client, tuesID) {
		return errors.New("Key already generated")
	}

	return client.SAdd(TUESID_KEY, tuesID).Err()
}

func checkPresence(client *redis.Client, tuesID string) bool {
	return client.SIsMember(TUESID_KEY, tuesID).Val()
}