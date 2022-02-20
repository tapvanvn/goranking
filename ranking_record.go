package goranking

import (
	"fmt"
	"sync"
)

type RankingRecord struct {
	table     *RankingTable
	beginRank Rank
	numUser   uint64
	users     map[string]uint64
	mux       sync.Mutex
}

func NewRankingRecord(rankingTable *RankingTable) *RankingRecord {

	return &RankingRecord{
		table: rankingTable,
		users: map[string]uint64{},
	}
}

func (record *RankingRecord) Leave(userID string) bool {

	record.mux.Lock()
	defer record.mux.Unlock()
	if _, ok := record.users[userID]; ok {

		delete(record.users, userID)
		record.numUser--
		var i uint64 = 0
		//fmt.Printf("leave user:%s  numafter:%d\n", userID, record.numUser)
		for user, _ := range record.users {
			record.users[user] = i
			i++
		}
		return true
	}
	return false
}

func (record *RankingRecord) Join(userID string) Rank {

	rank := Rank(0)

	record.mux.Lock()

	if testRank, ok := record.users[userID]; !ok {

		record.users[userID] = record.numUser
		rank = record.beginRank + Rank(record.numUser)

		record.numUser++

	} else {
		rank = Rank(testRank)
	}
	record.mux.Unlock()

	return rank
}

func (record *RankingRecord) Get(userID string) Rank {

	record.mux.Lock()
	defer record.mux.Unlock()

	if rank, ok := record.users[userID]; ok {

		return record.beginRank + Rank(rank)
	}
	return 0
}
func (record *RankingRecord) PrintDebug(score uint64, level int) {

	prefix := ""

	for i := 0; i < level; i++ {

		prefix += "\t"
	}

	record.mux.Lock()
	defer record.mux.Unlock()

	fmt.Printf("%sscore:%d begin:%d numUser:%d\n", prefix, score, record.beginRank, len(record.users))

	for user, rank := range record.users {

		fmt.Printf("\t%suser:%s rank:%d\n", prefix, user, record.beginRank+Rank(rank))
	}
}
