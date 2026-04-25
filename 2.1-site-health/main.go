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

var URLS_FILE_PATH = "websites.txt"

func readUrls() ([]string, error) {
	data, err := os.ReadFile(URLS_FILE_PATH)
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(data)), nil
}


func chechWebSite(url string, status_map *Statuses) {
	
	client := http.Client{Timeout: 5* time.Second}

	response, err:=client.Get(url)
	if err != nil{
		status_map.statuses <- fmt.Sprintf("%s - down", url)
		return
	}
	defer response.Body.Close()
	ok:=response.StatusCode == http.StatusOK

	if ok {
		status_map.statuses <-  fmt.Sprintf("%s - ok", url)
	}else{
		status_map.statuses <-  fmt.Sprintf("%s - down", url)
	}
}

type Statuses struct {
	statuses chan string
}

func (statuses *Statuses) Print(){
	for status := range statuses.statuses {
		fmt.Println(status)
	}
}


func main(){
	urls,err:=readUrls()

	if err!=nil{
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	wsStatus := Statuses{statuses:make(chan string)}


	for _,url:=range urls{
		wg.Add(1)
		go func(u string){
			defer wg.Done()
			chechWebSite(u, &wsStatus)
		}(url)
	}

	go func(){
		wg.Wait()
		close(wsStatus.statuses)
	}()
	wsStatus.Print()
	

}