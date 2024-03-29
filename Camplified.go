package main

import (
	"encoding/json"
	"fmt"
	"github.com/biter777/countries"
	tm "github.com/buger/goterm"
	"github.com/tkanos/gonfig"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Configuration struct {
	CLIENT_ID     string
	CLIENT_SECRET string
}

func config() Configuration {
	config := Configuration{}
	err := gonfig.GetConf("config.json", &config)

	if err != nil {
		f, _ := os.Create("config.json")
		s, _ := json.Marshal(Configuration{})
		_, _ = f.Write(s)
		fmt.Println("Looks you are the first time to use this program!")
		fmt.Println("Please fill the config.json file with your API v2 ID and Secret.")
		fmt.Println("Here's the quick link for you: https://osu.ppy.sh/home/account/edit#new-oauth-application")
		fmt.Println("DO NOT give your Secret Key to anyone but you.")
		fmt.Println("Be careful about giving out your Secret Key when asking anyone for help.")
		fmt.Println("You can press any key to exit and reopen this.")
		_, _ = fmt.Scanln()
		os.Exit(0)
	}

	return config
}

func getAccessToken() string {
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

func subScoreCalculation(score int) string {
	var subScore = strconv.Itoa(score)
	if math.Abs(float64(score)) >= 1000 {
		return "Away"
	}
	if score > 0 {
		subScore = "+" + subScore
	}
	return subScore
}

func GetLocation() (string, string, string, error) {
	url := "https://osu.ppy.sh/cdn-cgi/trace/"
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	lines := strings.Split(string(body), "\n")
	loc := ""
	colo := ""
	warp := ""

	for _, line := range lines {
		parts := strings.Split(line, "=")
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			if key == "loc" {
				loc = value
			} else if key == "colo" {
				colo = value
			} else if key == "warp" {
				warp = value
			}
		}
	}
	country := countries.ByName(loc).String()

	return country, colo, warp, nil
}

func main() {
	tm.Clear()

	accessToken := getAccessToken()
	if accessToken == "" {
		fmt.Println("Error: Can't get access_token, please check your config.json file.")
		fmt.Println("Please fill the config.json file with your API v2 ID and Secret.")
		fmt.Println("Here's the quick link for you: https://osu.ppy.sh/home/account/edit#new-oauth-application")
		fmt.Println("DO NOT give your Secret Key to anyone but you.")
		fmt.Println("Be careful about giving out your Secret Key when asking anyone for help.")
		fmt.Println("You can press any key to exit and reopen this.")
		_, _ = fmt.Scanln()
		os.Exit(0)
	}

	fmt.Println("Connected to Bancho Coordination Server:")
	fmt.Println()

	loc, colo, warp, err := GetLocation()
	if err != nil {
		fmt.Println("Cloudflare can't get your datacenter location.")
		fmt.Println("This is not meaning you can't submit score, please make sure in osu! client.")
	} else {
		switch warp {
		case "off":
			fmt.Println("You are not using Cloudflare WARP; Cloudflare identified you in", loc)
		case "on":
			fmt.Println("You are using Cloudflare WARP; Cloudflare identified you in", loc)
		case "plus":
			fmt.Println("You are using Cloudflare WARP+; Cloudflare identified you in", loc)
		}
		fmt.Println("Bancho connected you to Cloudflare", colo, "Edge Server.")
	}

	fmt.Println()

	kernet32 := syscall.NewLazyDLL("kernel32.dll")
	Beep := kernet32.NewProc("Beep")

	for {
		fmt.Print("Input bid to start: ")
		var bid string
		_, _ = fmt.Scanln(&bid)

		var data = getBeatmapInfo(bid, accessToken)
		bsData := data["beatmapset"].(map[string]interface{})
		var title = bsData["artist"].(string) + " - " + bsData["title"].(string) + " [" + data["version"].(string) + "]"
		var length = data["total_length"].(float64)

		tm.MoveCursor(1, 1)
		_, _ = tm.Println(title)

		tm.Flush()

		var caughtTime time.Time
		for {
			data = getBeatmapInfo(bid, accessToken)
			if data == nil {
				tm.MoveCursor(1, 1)
				_, _ = tm.Println("Error: Invalid bid")
				tm.Flush()
				break
			}
			status := data["status"].(string)

			tm.MoveCursor(4, 2)
			_, _ = tm.Println(status, time.Now())
			tm.MoveCursor(1, 3)
			_, _ = tm.Println("Map length in second(s):", length/1.5, "("+strconv.Itoa(int(length))+" without DT)")

			tm.Flush()

			if status == "ranked" {
				caughtTime = time.Now()
				Beep.Call(1000, 500)
				break
			}
		}

		tm.MoveCursor(1, 4)

		bsData = data["beatmapset"].(map[string]interface{})
		rankedTime, _ := time.Parse("2006-01-02T15:04:05Z", bsData["ranked_date"].(string))
		submittable := rankedTime.Add(time.Duration(length/1.5) * time.Second)

		_, _ = tm.Printf("Ranked on %s (+%d)\n", rankedTime.Format("2006-01-02 15:04:05"), caughtTime.Sub(rankedTime).Milliseconds())
		_, _ = tm.Printf("Submittable at %s\n", submittable.Format("2006-01-02 15:04:05"))

		tm.Flush()
		duration := submittable.Sub(time.Now())

		for duration > 0 {
			tm.MoveCursor(1, 7)
			_, _ = tm.Println("Submission accepted in", tm.Color(duration.Round(time.Second).String(), tm.YELLOW), "\t")
			tm.Flush()
			time.Sleep(time.Second)
			duration = submittable.Sub(time.Now())
		}

		tm.MoveCursor(1, 7)
		_, _ = tm.Println(tm.Color("Map accepted the score submission.", tm.GREEN))
		tm.Println()
		tm.Flush()

		time.Sleep(time.Second * 5)

		_, _ = tm.Printf("Rank\tSubmission_ID\tDiff\tSec\tRank\tPlayer\n")

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
				createdTime, _ := time.Parse("2006-01-02T15:04:05Z", score["created_at"].(string))
				subTime := createdTime.Sub(submittable).Seconds()
				var subSeconds string
				if subTime > 30 {
					subSeconds = "Away"
				} else {
					subSeconds = "+" + strconv.Itoa(int(subTime)) + "s"
				}
				_, _ = tm.Printf("#%d\t%.0f\t(%s)\t%s\t%s\t%s\n", i+1, score["id"], subScoreCalculation(subScore), subSeconds, score["rank"], user["username"])
			}

			tm.Println()
			tm.Flush()
			break
		}
	}
}
