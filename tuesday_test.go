package main

import (
	"testing"
	"net/http"
	"encoding/json"
)

func BenchmarkTuesday(t *testing.B)  {
	var netClient = &http.Client{}

	for i := 0; i < 100; i++ {
		response, err := netClient.Get("http://localhost:9090/tuesid")
		if err != nil {
			t.Fatal(err)
		}

		var body *TuesIDResponse
		json.NewDecoder(response.Body).Decode(&body)

		t.Log(body.TuesId)
	}
}
