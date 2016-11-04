package main

import (
	"gopkg.in/redis.v5"
	"strings"
)

func generatePrefixes(name string) []string {
	var prefixes []string

	for i := 0; i <= len(name); i++ {
		prefixes = append(prefixes, name[0:i])
	}

	return prefixes
}

func savePrefixes(client *redis.Client, name string) error {
	name = strings.ToLower(name)
	prefixes := generatePrefixes(name)

	for _, prefix := range prefixes {
		if err := client.SAdd(prefix, name).Err(); err != nil {
			return err
		}
	}

	return nil
}
