package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var URLS_FILE_PATH = "websites.txt"

func readUrls() ([]string, error) {
	data, err := os.ReadFile(URLS_FILE_PATH)
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(data)), nil
}

func chechWebSite(url string, status_map *Statuses, wg *sync.WaitGroup) {
	
	defer wg.Done()

	response, err:=http.Get(url)

	status_map.mu.Lock()
	defer status_map.mu.Unlock()

	if err != nil || response.StatusCode != 200{
		status_map.statuses[url] = false
		return
	}
	status_map.statuses[url] = true

}

type Statuses struct {
	mu sync.Mutex
	statuses map[string]bool
}

func (statuses *Statuses) Print(){
	for url,status := range statuses.statuses {
		if status{
			fmt.Printf("%s - ok\n", url)
		}else{
			fmt.Printf("%s - down\n", url)
		}
	}
}


func main(){
	urls,err:=readUrls()

	if err!=nil{
		log.Fatal(err)
	}

	wsStatus := Statuses{statuses:make(map[string]bool)}

	var wg sync.WaitGroup


	for _,url:=range urls{
		wg.Add(1)
		go chechWebSite(url, &wsStatus, &wg)
	}
	wg.Wait()

	wsStatus.Print()
	

}