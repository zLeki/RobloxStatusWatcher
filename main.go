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
	stalkUserID = "293025549" // userid
	RobloSecurity = ""
	client http.Client
	online bool
	playerStatus string
)
type SMSRequestBody struct {
	From      string `json:"from"`
	Text      string `json:"text"`
	To        string `json:"to"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}
func init() {
	err := os.Setenv("NEXMO_API_KEY", "")
	if err != nil {
		return
	}
	err = os.Setenv("NEXMO_API_SECRET", "")
	if err != nil {
		return
	}

}
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
			if playerStatus != "InGame" {
				sendSms("Robux Stalker: "+getUsernameFromUserID(userInfo.VisitorID)+"is in game! Enlarge: https://robloxstalker.leki2.repl.co/?placeid=" + strconv.Itoa(int(userInfo.PlaceID.(float64))) + "&username=" + getUsernameFromUserID(userInfo.VisitorID) + "&placename=" + strings.Replace(strings.Replace(strings.Replace(strings.Replace(userInfo.LastLocation, " ", "%20", -1), "[", "%5B", -1), "]", "%5D", -1), "!", "%21", -1), "InGame")
			}
			log.Println("In game")
			return userInfo, nil
		} else {
			log.Println(userInfo)
			if playerStatus != "Online" {
				sendSms("Player is online", "Online")
			}
			return UserInfo{}, errors.New("User is online but not in a game refreshing in 60 seconds")
		}
	}else{
		if playerStatus != "Offline" {
			sendSms("Player is offline", "Offline")
		}

		return UserInfo{}, errors.New("User is not online refreshing in 60 seconds")
	}
}
func sendSms(content, status string) {
	body := SMSRequestBody{
		APIKey:    os.Getenv("NEXMO_API_KEY"),
		APISecret: "",
		To:        "",
		From:      "",
		Text:      content,
	}
	log.Println(body.APISecret)
	smsBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post("https://rest.nexmo.com/sms/json", "application/json", bytes.NewBuffer(smsBody))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		log.Println("SMS sent")

	} else {
		log.Println("SMS failed")
	}
	playerStatus = status
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
					log.Println("Enlarge: http://localhost:8081/?placeid=" + strconv.Itoa(iAreaId) + "&username=" + s + "&userID="+strconv.Itoa(gameID.VisitorID)+"&placename=" + strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(gameID.LastLocation, " ", "-", -1), "[", "-", -1), "]", "-", -1), "!", "-", -1), "Playing", "", -1))
					go func() {
						err := http.ListenAndServe(":8081", nil)
						if err != nil {
							log.Fatal(err)
						}
					}()
					log.Println("Server started")
					online = true
				}else{
					log.Println("Enlarge: http://localhost:8081/?placeid=" + strconv.Itoa(iAreaId) + "&username=" + s + "&userID="+strconv.Itoa(gameID.VisitorID)+"&placename=" + strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(gameID.LastLocation, " ", "-", -1), "[", "-", -1), "]", "-", -1), "!", "-", -1), "Playing", "", -1))
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
