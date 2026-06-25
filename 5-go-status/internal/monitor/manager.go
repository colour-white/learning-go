package monitor

import (
	"context"
	"fmt"
	"go-status/internal/models"
	"sync"
	"time"
	"database/sql"
)

type Manager struct{
	activeWorkers map[int]context.CancelFunc
	mu sync.Mutex // protect the map from concurrent API requests
	db *sql.DB
}

func (m *Manager) StartTarget(ctx context.Context, target *models.Target) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// create a child context for this target
	targetCtx, cancel:=context.WithCancel(ctx)
	m.activeWorkers[target.Id] = cancel

	go func (){
		ticker := time.NewTicker(time.Duration(target.Interval_sec) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-targetCtx.Done():
				fmt.Printf("Worler for %s shutting down...\n", target.Url)
				return
			case <-ticker.C:
				fmt.Printf("Probing %s...\n", target.Url)
				probe:=checkWebSite(target)
				probe,err :=models.InsertProbe(m.db, probe)
				if err!=nil{
					fmt.Println(err.Error())
					return
				}
			}
		}

	}()
}

func (m *Manager) StopTarget(targetId int){
	m.mu.Lock()
	defer m.mu.Unlock()
	cancelF, ok := m.activeWorkers[targetId]
	if ok {
		cancelF()
		delete(m.activeWorkers,targetId)
	}
}