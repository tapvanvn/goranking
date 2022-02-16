package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
)

const MAX = 100_000_000

func main() {
	for {
		userID := rand.Int63n(MAX)
		url := fmt.Sprintf("http://127.0.0.1:9000/score/%d", userID)

		if resp, err := http.Get(url); err == nil {

			if resp.StatusCode == http.StatusOK {

				if data, err := ioutil.ReadAll(resp.Body); err == nil {

					fmt.Println("user:", userID, "rank:", string(data))
				}
			}
			resp.Close = true

		} else {

			fmt.Println(err.Error())
		}
	}
}
