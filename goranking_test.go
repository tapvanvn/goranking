package goranking_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/tapvanvn/goranking"
)

func Test(t *testing.T) {

	system := goranking.NewRankingSystem(10)

	for i := 1; i < 100; i++ {
		if i > 15 && i < 20 {
			continue
		}
		rank := system.PutScore(strconv.Itoa(i), 0, uint64(i))
		fmt.Println("puyt rank:", i, rank)
	}

	time.Sleep(2 * time.Second)
	system.PrintDebug()

	for i := 10; i < 20; i++ {
		rank := system.GetScore(strconv.Itoa(i), uint64(i))
		fmt.Println("rank:", i, rank)
	}

	time.Sleep(2 * time.Second)
	for i := 1; i < 30; i++ {
		rank := system.GetScore(strconv.Itoa(i), uint64(i))
		fmt.Println("rank:", i, rank)
	}
}
