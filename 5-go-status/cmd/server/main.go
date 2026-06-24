package main

import (
	"fmt"
	"go-status/internal/models"
	"go-status/internal/storage"
)

func main() {

	db, err:=storage.InitDatabase("sqlite.db")

	if err!=nil{
		panic(err)
	}

	t,err:=models.InsertTarget(db, "google.com", 10, "mail.com", true)

	fmt.Println(t)


	var targets []models.Target
	targets, err = models.SelectAllTargets(db)

	fmt.Println(targets)


	
}
