package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)


const workerCount = 3
const urlsFilePath = "websites.txt"

func readUrls() ([]string, error) {
	data, err := os.ReadFile(urlsFilePath)
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(data)), nil
}


func chechWebSite(url string, results chan string) {
	
	client := http.Client{Timeout: 5* time.Second}

	response, err:=client.Get(url)
	if err != nil{
		results <- fmt.Sprintf("%s - down", url)
		return
	}
	defer response.Body.Close()
	ok:=response.StatusCode == http.StatusOK

	if ok {
		results <-  fmt.Sprintf("%s - ok", url)
	}else{
		results <-  fmt.Sprintf("%s - down", url)
	}
}


func Scrape(id int, jobs, results chan string, wg * sync.WaitGroup){
	defer wg.Done()
	for job:=range jobs{
		fmt.Printf("Worker %d started working on %s\n", id, job)
		chechWebSite(job, results)
	}
}


func main(){

	
	urls,err :=readUrls()
	if err!=nil{
		log.Fatal(err)
	}


	jobs:=make(chan string)
	results:=make(chan string)
	var wg sync.WaitGroup

	for w:=1; w<=workerCount; w++{
		wg.Add(1)
		go Scrape(w, jobs, results, &wg)
	}

	
	go func(){
		wg.Wait()
		close(results)
	}()
	
	go func(){
		for result:=range results{
			fmt.Println(result)
		}
	}()

	for _,url:=range urls{
		jobs<-url
	}

	close(jobs)
	

}