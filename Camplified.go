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

	// get config from config.json
	config := Configuration{}
	err := gonfig.GetConf("config.json", &config)

	if err != nil {
		f, _ := os.Create("config.json")
		s, _ := json.Marshal(Configuration{})
		_, _ = f.Write(s)
		fmt.Println("Looks you are the first time to use this program, please fill the config.json file.")
		fmt.Println("You can press any key to exit and reopen this.")
		_, _ = fmt.Scanln()
		os.Exit(0)
	}

	return config
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

func getBeatmapScores(bid string, accessToken string) []map[string]interface{} {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://osu.ppy.sh/api/v2/beatmaps/"+bid+"/scores", nil)
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
	}
	scores := data["scores"].([]interface{})
	var scData []map[string]interface{}
	for _, score := range scores {
		scData = append(scData, score.(map[string]interface{}))
	}
	return scData
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
				caughtTime := time.Now()

				tm.MoveCursor(1, 4)
				tm.Flush() // why flush here???

				length := data["total_length"].(float64) / 1.5
				rankedTime, _ := time.Parse("2006-01-02T15:04:05Z", bsData["ranked_date"].(string))
				submittable := rankedTime.Add(time.Duration(length) * time.Second)

				_, _ = fmt.Printf("Ranked on %s (+%d)\n", rankedTime.Format("2006-01-02 15:04:05"), caughtTime.Sub(rankedTime).Milliseconds())
				_, _ = fmt.Printf("Submittable at %s\n", submittable.Format("2006-01-02 15:04:05"))

				tm.Flush()
				duration := submittable.Sub(time.Now())

				for duration > 0 {
					tm.MoveCursor(1, 7)
					tm.Flush()
					fmt.Print("Submittable in ")
					fmt.Print(tm.Color(duration.Round(time.Second).String(), tm.YELLOW))
					tm.Flush()
					time.Sleep(time.Second)
					duration = submittable.Sub(time.Now())
				}

				fmt.Println()
				tm.Flush()

				time.Sleep(10)

				for {
					scData := getBeatmapScores(bid, accessToken)

					if len(scData) == 0 {
						continue
					}

					for i, score := range scData {
						if i+1 == 9 {
							break
						}
						user := score["user"].(map[string]interface{})
						var subScore = 0
						if i != 0 {
							subScore = int(score["id"].(float64) - scData[i-1]["id"].(float64))
						}
						fmt.Printf("%.0f\t(+%d)\t#%d\t%s\t%s\t\n", score["id"], subScore, i+1, score["rank"], user["username"])
					}
					break
				}

				fmt.Println()
				tm.Flush()
				break
			}
		}
	}
}
