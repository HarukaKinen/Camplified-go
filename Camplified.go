package main

import (
	"encoding/json"
	"fmt"
	tm "github.com/buger/goterm"
	"github.com/tkanos/gonfig"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Configuration struct {
	CLIENT_ID     string
	CLIENT_SECRET string
}

func config() Configuration {

	err := gonfig.GetConf("config.json", &Configuration{})
	if err != nil {
		f, _ := os.Create("config.json")
		s, _ := json.Marshal(Configuration{})
		_, _ = f.Write(s)
		fmt.Println("Looks you are the first time to use this program, please fill the config.json file.")
		fmt.Println("You can press any key to exit and reopen this.")
		_, _ = fmt.Scanln()
		os.Exit(0)
	}

	return Configuration{}
}

func getAccessToken() string {
	// get access_token from osu api v2

	config := config()
	params := url.Values{}
	params.Add("client_id", config.CLIENT_ID)
	params.Add("client_secret", config.CLIENT_SECRET)
	params.Add("grant_type", "client_credentials")
	params.Add("scope", "public")
	resp, err := http.PostForm("https://osu.ppy.sh/oauth/token", params)

	if err != nil {
		return ""
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return ""
	}
	return data["access_token"].(string)
}

func getBeatmapInfo(bid string, accessToken string) map[string]interface{} {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://osu.ppy.sh/api/v2/beatmaps/"+bid, nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		return nil
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	} else {
		return data
	}
}

func main() {
	tm.Clear()

	accessToken := getAccessToken()
	if accessToken == "" {
		fmt.Println("Error: Can't get access_token, please check your config.json file.")
		_, _ = fmt.Scanln()
		os.Exit(0)
	}
	fmt.Println("Get access token successful.")

	for {
		fmt.Print("Input bid to start: ")
		var bid string
		_, _ = fmt.Scanln(&bid)

		tm.Clear()
		for {
			data := getBeatmapInfo(bid, accessToken)
			if data == nil {
				tm.MoveCursor(1, 1)
				_, _ = tm.Println("Error: Invalid bid")
				tm.Flush()
				break
			}
			status := data["status"].(string)
			bsData := data["beatmapset"].(map[string]interface{})
			title := bsData["artist"].(string) + " - " + bsData["title"].(string) + " [" + data["version"].(string) + "]"

			tm.MoveCursor(1, 1)
			_, _ = tm.Println(title)

			tm.MoveCursor(4, 2)
			_, _ = tm.Println(status, time.Now())

			tm.Flush()

			if status == "ranked" {
				tm.MoveCursor(1, 4)
				_, _ = tm.Println("Caught ranked at", time.Now().Format("2006-01-02 15:04:05"))
				tm.Flush()
				break
			}
		}
	}
}
