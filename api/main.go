package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type config struct {
	YoutubeApiKey string `json:youtubeApiKey`
}
type apiResponse struct {
	Kind string `json:kind`
}

func main() {
	var youtubeListId string
	var config config

	// 引数チェック
	if len(os.Args) != 2 {
		fmt.Println("Invalid args")
		return
	} else {
		youtubeListId = os.Args[1]
	}

	getPlayListAPIUrl := "https://www.googleapis.com/youtube/v3/playlistItems"
	// 設定ファイル読み込み
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		fmt.Println(err.Error())
	}
	resp, err := http.Get(getPlayListAPIUrl + "?key=" + config.YoutubeApiKey + "&part=snippet&playlistId=" + youtubeListId)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if resp.StatusCode != 200 {
		fmt.Println("ERROR: StatusCode:" + strconv.Itoa(resp.StatusCode))
		return
	}

	//fmt.Println("url:" + getPlayListAPIUrl + "?key=" + config.YoutubeApiKey)

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var x apiResponse

	err = json.Unmarshal(byteArray, &x)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("(%%#v) %#v\n", x)

	//fmt.Println(string(byteArray))
}

// https://www.googleapis.com/youtube/v3/playlistItems?key=AIzaSyBsIWp3vAWFF6_UfX_4F5P5Pb7KCMcSaiM&part=snippet&playlistId=PLRdiaanKAFQl5ERDgJHx2ZRKCcIl-I8fz&maxResults=50
