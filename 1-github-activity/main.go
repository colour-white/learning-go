package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"os"
)

type Config struct{
	GIT_TOKEN string
}

func read_config() Config {
	env, err := godotenv.Read(".env")
	if err!=nil{
		log.Fatal("Could not load config: ", err)
	}
	return Config{GIT_TOKEN: env["GIT_TOKEN"]}
}

var config = read_config()

func read_username()(string, error){

	fmt.Println("Enter github user name: ")

	var username string

	_, err := fmt.Scanln(&username)

	if err != nil {
		log.Fatal(err)
	}
	return username, nil
}


type Printer interface {
	TextPrint()
	JSONPrint()

}

type GitUserInfoResponse struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	Bio string `json:"bio"`
	PublicRepos int  `json:"public_repos"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

func (gituserinfo GitUserInfoResponse) TextPrint(){
	fmt.Println("Login: ", gituserinfo.Login)
	fmt.Println("Bio: ", gituserinfo.Bio)
	fmt.Println("Public repos: ", gituserinfo.PublicRepos)
	fmt.Println("Followers: ", gituserinfo.Followers)
}

func (gituserinfo GitUserInfoResponse) JSONPrint(){
	encoder:= json.NewEncoder(os.Stdout)
	err:=encoder.Encode(gituserinfo)
	if err != nil {
		log.Fatal(err)
	}
}

type GitRepo struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Url string `json:"url"`

}

type GitEvent struct {
	Id string `json:"id"`
	Type string `json:"type"`
	Repo GitRepo `json:"repo"`
	CreatedAt string `json:"created_at"`
}

func (gitevent GitEvent) TextPrint(){
	switch gitevent.Type{
		case "PushEvent":{
			fmt.Printf("%s: Pushed to %s\n", gitevent.CreatedAt, gitevent.Repo.Name)
		}
		case "ReleaseEvent":{
			fmt.Printf("%s: Created a release at %s\n", gitevent.CreatedAt, gitevent.Repo.Name)
		}
		default :{
			fmt.Printf("%s: Did something!\n", gitevent.CreatedAt)
		}
	}
}

func (gitevent GitEvent) JSONPrint(){
	encoder:= json.NewEncoder(os.Stdout)
	err:=encoder.Encode(gitevent)
	if err != nil {
		log.Fatal(err)
	}
}

func get_github_info(username string) (GitUserInfoResponse, error){

	url:= fmt.Sprintf("https://api.github.com/users/%s", username)

	request,_:=http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "Bearer "+config.GIT_TOKEN)
	resp, err := http.DefaultClient.Do(request)

	if err != nil{
		log.Fatal(err)
	}

	var body GitUserInfoResponse
	decoder:= json.NewDecoder(resp.Body)
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

func get_github_event_history(username string) ([]GitEvent, error){
	url:=fmt.Sprintf("https://api.github.com/users/%s/events/public", username)

	request,_:=http.NewRequest("GET", url, nil)
	request.Header.Add("Authorization", "Bearer "+config.GIT_TOKEN)
	resp, err := http.DefaultClient.Do(request)

	if err != nil {
		log.Fatal(err)
	}

	var body []GitEvent
	decoder:=json.NewDecoder(resp.Body)
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
	info,err := get_github_info(username)
	if err != nil{
		fmt.Println("Could not fetch user info:", err)
	}
	info.TextPrint()
	info.JSONPrint()

	events, err:= get_github_event_history(username)
	if err != nil{
		fmt.Println("Could not fetch user history:", err)
	}
	for _, e := range events {
		e.TextPrint()
		e.JSONPrint()
	}


}