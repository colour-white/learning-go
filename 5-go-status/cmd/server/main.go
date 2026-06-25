package main

import (
	"fmt"
	"go-status/internal/models"
	"go-status/internal/storage"
	"time"
)

func main() {

	db, err:=storage.InitDatabase("sqlite.db")

	if err!=nil{
		panic(err)
	}
	t:= &models.Target{Url:"google.com", Interval_sec: 10, Contact_info: "mail.com", Is_active: true, Created_at: time.Now()}
	t,err = models.InsertTarget(db, t)

	fmt.Println(t)


	var targets []models.Target
	targets, err = models.SelectAllTargets(db)

	fmt.Println(targets)


	
}
