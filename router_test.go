package jaakitrouter

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func Test(t *testing.T) {
	go func() {
		router := &Router{DebugLog: true}

		router.GET("/", func(w http.ResponseWriter, r *http.Request, params NamedParams) {
			w.Write([]byte("INDEX"))
		})

		router.GET("/winner/@name", func(w http.ResponseWriter, r *http.Request, params NamedParams) {
			w.WriteHeader(200)
			w.Header().Set("content-type", "application/json")

			json.NewEncoder(w).Encode(params)
		})

		router.Logger("Starting Server on :8081")
		log.Fatal(http.ListenAndServe(":8081", router))
	}()

	resp, err := http.Get("http://localhost:8081/")
	if err != nil {
		t.Error("Couldn't do request")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error("StatusCode should have been 200")
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("Error decoding resp.Body")
	}
	if string(data) != "INDEX" {
		t.Error("Wanted to see the INDEX")
	}

	respTwo, err := http.Get("http://localhost:8081/winner/jaaki")
	if err != nil {
		t.Error("Couldn't do request")
	}
	defer respTwo.Body.Close()
	if respTwo.StatusCode != 200 {
		t.Error("StatusCode should have been 200")
	}

	expected := make(map[string]interface{})

	json.NewDecoder(respTwo.Body).Decode(&expected)
	if expected["winner"] != "winner" {
		t.Error("Should be winner")
	}
	if expected["name"] != "jaaki" {
		t.Error("Name should be jaaki")
	}

}
