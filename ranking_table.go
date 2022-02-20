package goranking

import (
	"log"
	"sort"
	"sync"
)

type RankingTable struct {
	mux       sync.Mutex
	beginRank Rank
	numRecord uint64
	records   map[uint64]*RankingRecord
	scores    []uint64
}

func NewRankingTable() *RankingTable {

	return &RankingTable{

		records: map[uint64]*RankingRecord{},
		scores:  make([]uint64, 0),
	}
}

func (table *RankingTable) Leave(score uint64, userID string) {

	table.mux.Lock()
	record, ok := table.records[score]
	table.mux.Unlock()

	needUpdate := false
	if !ok {
		log.Fatal("record not found")
	}

	if record.Leave(userID) {

		table.mux.Lock()
		table.numRecord--
		table.mux.Unlock()

		needUpdate = true
	}

	if needUpdate {

		table.UpdateRecordRank()
	}
}

func (table *RankingTable) Join(score uint64, userID string) Rank {

	table.mux.Lock()

	record, ok := table.records[score]

	needUpdate := false

	//fmt.Println("begin join table:", table.beginRank, userID, "score:", score, "num:", table.numRecord)

	if !ok {

		record = NewRankingRecord(table)

		table.records[score] = record

		table.scores = append(table.scores, score)

		sort.Slice(table.scores, func(i, j int) bool { return table.scores[i] < table.scores[j] })

		needUpdate = true
	}
	table.numRecord++

	table.mux.Unlock()

	//fmt.Println("\tafter join table:", table.beginRank, userID, "score:", score, "num:", table.numRecord)

	resultRank := table.beginRank + record.Join(userID)

	if needUpdate || resultRank > 0 {

		table.UpdateRecordRank()
	}

	return resultRank
}

func (table *RankingTable) Get(lastScore uint64, userID string) Rank {

	table.mux.Lock()
	record, ok := table.records[lastScore]
	table.mux.Unlock()

	if ok {
		recordRank := record.Get(userID)

		resultRank := table.beginRank + recordRank
		//fmt.Printf("table get: userID:%s last:%d begin:%d recordRank:%d result:%d\n", userID, lastScore, table.beginRank, recordRank, resultRank)
		return resultRank
	}
	return 0
}

func (table *RankingTable) UpdateRecordRank() {

	currRank := Rank(0)
	table.mux.Lock()
	table.numRecord = 0

	for _, score := range table.scores {

		record := table.records[score]
		record.beginRank = currRank

		//fmt.Printf("update table record: score:%d begin:%d numuser:%d cur:%d\n", score, record.beginRank, record.numUser, currRank)

		table.numRecord += record.numUser
		currRank += Rank(record.numUser)
	}

	table.mux.Unlock()

}

func (table *RankingTable) PrintDebug(level int) {
	prefix := ""
	for i := 0; i < level; i++ {
		prefix += "\t"
	}
	table.mux.Lock()
	for _, score := range table.scores {
		record := table.records[score]
		record.PrintDebug(score, level+1)
	}
	table.mux.Unlock()
}
