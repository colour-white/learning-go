package monitor

import (
	"context"
	"database/sql"
	"fmt"
	"go-status/internal/models"
	"sync"
	"time"
)

type Manager struct {
	ActiveWorkers map[int]context.CancelFunc
	mu            sync.Mutex // protect the map from concurrent API requests
	DB            *sql.DB
	RootCtx       context.Context
}

func (m *Manager) StartTarget(target *models.Target) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// create a child context for this target
	targetCtx, cancel := context.WithCancel(m.RootCtx)
	m.ActiveWorkers[target.Id] = cancel

	go func() {
		ticker := time.NewTicker(time.Duration(target.Interval_sec) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-targetCtx.Done():
				fmt.Printf("Worker for %s shutting down...\n", target.Url)
				return
			case <-ticker.C:
				fmt.Printf("Probing %s...\n", target.Url)
				probe := checkWebSite(target)
				probe, err := models.InsertProbe(m.DB, probe)
				if err != nil {
					fmt.Println(err.Error())
					return
				}
			}
		}

	}()
}

func (m *Manager) StopTarget(targetId int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cancelF, ok := m.ActiveWorkers[targetId]
	if ok {
		cancelF()
		delete(m.ActiveWorkers, targetId)
	}
}

func (m *Manager) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for targetId, cancel := range m.ActiveWorkers {
		cancel()
		delete(m.ActiveWorkers, targetId)
	}
}
