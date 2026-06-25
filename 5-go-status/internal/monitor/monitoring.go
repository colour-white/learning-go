package monitor

import (
	"database/sql"
	"fmt"
	"go-status/internal/models"
	"time"
)


func Monitor(t*models.Target, db* sql.DB){
	for {
		
		probe:=checkWebSite(t)
		// models.InsertProbe(db)
		fmt.Println(probe)
		time.Sleep(time.Duration(t.Interval_sec))
	}
}