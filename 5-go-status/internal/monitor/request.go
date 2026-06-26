package monitor

import (
	"go-status/internal/models"
	"net/http"
	"time"
)

func checkWebSite(t *models.Target) *models.Probe {
	client := http.Client{Timeout: 5 * time.Second}
	start := time.Now()
	response, err := client.Get(t.Url)
	duration := time.Since(start)
	if err != nil || duration.Seconds() > 5 {
		return &models.Probe{
			Target_id:   t.Id,
			Err_msg:     err.Error(),
			Timestamp:   time.Now(),
			Status_code: 500, // do not overthink it, yet
			Latency_ms:  int(duration.Milliseconds()),
		}
	}
	defer response.Body.Close()
	return &models.Probe{
		Target_id:   t.Id,
		Timestamp:   time.Now(),
		Status_code: response.StatusCode,
		Latency_ms:  int(duration.Milliseconds()),
	}
}
