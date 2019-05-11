package gamedata

import (
	"sync"
)

type MinerList struct {
	MinerMap map[string]int32 `json:"miners"` // [address]id
	mux      sync.Mutex
}

func (miners *MinerList) Add(addr string, id int32) {
	miners.mux.Lock()
	defer miners.mux.Unlock()
	miners.MinerMap[addr] = id
}
