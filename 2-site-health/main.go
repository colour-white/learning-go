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

func (s *Statuses) update(url string, ok bool){
	
		s.mu.Lock()
		defer s.mu.Unlock()
		s.statuses[url] = ok
}

func chechWebSite(url string, status_map *Statuses, wg *sync.WaitGroup) {
	
	defer wg.Done()
	client := http.Client{Timeout: 5* time.Second}

	response, err:=client.Get(url)
	if err != nil{
		status_map.update(url, false)
		return
	}
	defer response.Body.Close()
	ok:=response.StatusCode == http.StatusOK
	status_map.update(url, ok)

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