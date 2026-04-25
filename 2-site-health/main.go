package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var URLS_FILE_PATH = "websites.txt"

func readUrls() ([]string, error) {
	data, err := os.ReadFile(URLS_FILE_PATH)
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(data)), nil
}

func chechWebSite(url string) bool{
	
	response, err:=http.Get(url)

	if err != nil || response.StatusCode != 200{
		return false
	}
	return true

}


func main(){
	urls,err:=readUrls()

	if err!=nil{
		log.Fatal(err)
	}


	for _,url:=range urls{
		if chechWebSite(url){
			fmt.Printf("%s - ok\n", url)
		}else{
			fmt.Printf("%s - down\n", url)
		}
	}

}