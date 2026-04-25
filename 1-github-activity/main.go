package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)


func read_username()(string, error){

	fmt.Println("Enter github user name: ")

	var username string

	_, err := fmt.Scanln(&username)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("You've entered: ", username)
	return username, nil

}


type GitUserResponse struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	PublicRepos int  `json:"public_repos"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

func get_github_info(username string) (GitUserResponse, error){

	resp, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s", username))

	if err != nil{
		// return nil, err
		log.Fatal(err)
	}


	decoder:= json.NewDecoder(resp.Body)

	var body GitUserResponse

	err=decoder.Decode(&body)



	if resp.StatusCode > 299 {
		log.Fatal("Response failed with status code %d and \nbody %s\n", resp.StatusCode, body)
	}
	if err!=nil{
		log.Fatal(err)
	}

	fmt.Println("Fetched successfully.")

	
	return body, nil

}


func main(){

	username,_:=read_username()

	fmt.Println(username)

	info,err := get_github_info(username)

	if err != nil{
		fmt.Println("Could not fetch user info:", err)
	}




	fmt.Println(info)

}