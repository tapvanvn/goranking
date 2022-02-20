package goranking

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

type Rank uint64

type RankingSystem struct {
	tableSize       int //the size of each table
	tables          map[uint64]*RankingTable
	tableIDs        []uint64
	tableMux        sync.Mutex
	updateMux       sync.Mutex //sync the update loop
	lastUpdateScore uint64     // 0 mean no update, only store the smallest score
	maxRank         Rank
}

func (sys *RankingSystem) PutScore(userID string, oldScore uint64, score uint64) Rank {

	tableID := score / uint64(sys.tableSize)

	if oldScore > 0 {

		oldTableID := oldScore / uint64(sys.tableSize)

		sys.tableMux.Lock()
		oldTable, ok := sys.tables[oldTableID]
		sys.tableMux.Unlock()

		if !ok {

			log.Fatal("last score not found")
		}

		oldTable.Leave(oldScore, userID)

	}
	if score == 0 {

		return Rank(0)
	}
	sys.tableMux.Lock()

	table, ok := sys.tables[tableID]

	if !ok {

		table = NewRankingTable()
		sys.tables[tableID] = table
		sys.tableIDs = append(sys.tableIDs, tableID)
		sort.Slice(sys.tableIDs, func(i, j int) bool { return sys.tableIDs[i] < sys.tableIDs[j] })
	}
	sys.tableMux.Unlock()

	rank := table.Join(score, userID)

	if rank > sys.maxRank {

		sys.maxRank = rank
	}

	sys.updateMux.Lock()
	if score < sys.lastUpdateScore || sys.lastUpdateScore == 0 {
		sys.lastUpdateScore = score
	}
	sys.updateMux.Unlock()

	return rank
}

func (sys *RankingSystem) GetScore(userID string, lastScore uint64) Rank {

	tableID := lastScore / uint64(sys.tableSize)
	sys.tableMux.Lock()
	table, ok := sys.tables[tableID]
	sys.tableMux.Unlock()

	if ok {
		//fmt.Println("get:", userID, lastScore, tableID)
		rank := table.Get(lastScore, userID)
		if rank > 0 {

			if rank > sys.maxRank {

				return 1
			}
			return sys.maxRank - rank
		}
	}
	return Rank(0)
}

func (sys *RankingSystem) run() {

	for {

		sys.updateMux.Lock()
		lastUpdateScore := sys.lastUpdateScore
		sys.lastUpdateScore = 0
		sys.updateMux.Unlock()

		if lastUpdateScore == 0 {

			time.Sleep(25 * time.Millisecond)
			continue
		}

		sys.tableMux.Lock()

		fmt.Println("---BEGIN UPDATE---")

		lastRank := Rank(0)

		for _, tableID := range sys.tableIDs {

			table := sys.tables[tableID]

			table.beginRank = lastRank

			lastRank += Rank(table.numRecord)
		}

		fmt.Println("---END UPDATE---")
		sys.tableMux.Unlock()

	}
}

func (sys *RankingSystem) PrintDebug() {

	sys.tableMux.Lock()

	for _, tableID := range sys.tableIDs {

		table := sys.tables[tableID]
		fmt.Println("table:", tableID, "begin:", table.beginRank, "num:", table.numRecord)
		table.PrintDebug(1)
	}
	sys.tableMux.Unlock()
}

func NewRankingSystem(tableSize int) *RankingSystem {

	sys := &RankingSystem{

		tableSize: tableSize,
		tables:    map[uint64]*RankingTable{},
	}
	go sys.run()
	return sys
}
