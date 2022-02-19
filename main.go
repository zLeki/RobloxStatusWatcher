package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)
var (
	stalkUserID = "439403718"
	RobloSecurity = ""
	client http.Client
	online bool
)


type UserInfo struct {
	IsOnline     bool        `json:"IsOnline"`
	LastLocation string      `json:"LastLocation"`
	PlaceID      interface{} `json:"PlaceId"`
	VisitorID    int         `json:"VisitorId"`
}


func getUserStatus() (UserInfo, error) {
	resp, err := http.NewRequest("GET", "https://api.roblox.com/users/"+stalkUserID+"/onlinestatus/", nil)
	resp.AddCookie(&http.Cookie{Name: ".ROBLOSECURITY", Value: RobloSecurity})
	if err != nil {
		return UserInfo{}, err
	}
	req, err := client.Do(resp)
	if err != nil {
		dataBytes, _ := ioutil.ReadAll(req.Body)
		log.Println("error decoding response", err, string(dataBytes))
		return UserInfo{}, err
	}
	var userInfo UserInfo
	err = json.NewDecoder(req.Body).Decode(&userInfo)
	if err != nil {
		dataBytes, _ := ioutil.ReadAll(req.Body)
		log.Println("error decoding response", err, string(dataBytes))
		return UserInfo{}, err
	}
	if userInfo.IsOnline {
		if userInfo.PlaceID != nil {
			return userInfo, nil
		} else {
			return UserInfo{}, errors.New("User is online but not in a game refreshing in 60 seconds")
		}
	}else{
		return UserInfo{}, errors.New("User is not online refreshing in 60 seconds")
	}
}

func main() {
	for {
		gameID, err := getUserStatus()
		if err != nil {
			log.Println("error while fetching game status: waiting 15 seconds ", err)
			time.Sleep(15 * time.Second)

		}else{
			if gameID.PlaceID != nil {
				var iAreaId = int(gameID.PlaceID.(float64))
				s:=getUsernameFromUserID(gameID.VisitorID)
				log.Println(s+" is in game:", iAreaId)
				if online != true {
					http.HandleFunc("/", WebPage)
					log.Println("Enlarge: https://robloxstalker.leki2.repl.co/?placeid=" + strconv.Itoa(iAreaId) + "&username=" + s + "&placename=" + strings.Replace(strings.Replace(strings.Replace(strings.Replace(gameID.LastLocation, " ", "%20", -1), "[", "%5B", -1), "]", "%5D", -1), "!", "%21", -1))
					go func() {
						err := http.ListenAndServe(":8081", nil)
						if err != nil {
							log.Fatal(err)
						}
					}()
					log.Println("Server started")
					online = true
				}else{
					log.Println("Already online user is in game:", iAreaId)
				}

			}
			time.Sleep(10 * time.Second)
		}}
}
func WebPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func getUsernameFromUserID(userID int) string {
	type Struct struct {
		DisplayName string `json:"Name"`
	}
	resp, err := http.Get("https://users.roblox.com/v1/users/" + strconv.Itoa(userID))
	if err != nil {
		log.Println(err, resp.Body)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	var s Struct
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		log.Println(err)
		return "unable to fetch"
	}
	return s.DisplayName

}
