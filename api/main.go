package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"youtube-app/api/model"
)

func main() {
	var youtubeListId string
	var config model.Config
	var apiPlayListItemListResponse model.ApiPlayListItemListResponse

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

	err = json.Unmarshal(byteArray, &apiPlayListItemListResponse)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("(%%#v) %#v\n", apiPlayListItemListResponse)

	//fmt.Println(string(byteArray))
}
