package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/tapvanvn/goranking"
	"github.com/tapvanvn/goutil"
)

const MAX = 100

var system = goranking.NewRankingSystem(10)

var mux sync.Mutex
var lastScore map[uint64]uint64 = map[uint64]uint64{}

func RandomScore() {

	mux.Lock()
	user := rand.Int63n(MAX)
	score := rand.Int63n(MAX)

	if last, ok := lastScore[uint64(user)]; ok {

		system.PutScore(strconv.FormatInt(user, 10), last, uint64(score))

	} else {

		system.PutScore(strconv.FormatInt(user, 10), 0, uint64(score))
	}

	lastScore[uint64(user)] = uint64(score)
	mux.Unlock()
}
func Debug() {
	system.PrintDebug()
}
func main() {

	port := "9000"

	goutil.Schedule(RandomScore, 5*time.Millisecond)
	goutil.Schedule(Debug, 1*time.Second)

	http.HandleFunc("/score/", func(w http.ResponseWriter, r *http.Request) {

		userString := r.URL.Path[7:]

		if userID, err := strconv.ParseInt(userString, 10, 64); err == nil {
			mux.Lock()
			last, ok := lastScore[uint64(userID)]
			defer mux.Unlock()
			if ok {
				rank := system.GetScore(userString, last)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(strconv.FormatUint(uint64(rank), 10)))
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	})

	fmt.Println("run on port:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
