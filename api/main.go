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

	// 引数チェック
	if len(os.Args) != 2 {
		fmt.Println("Invalid args")
		return
	} else {
		youtubeListId = os.Args[1]
	}

	getPlayListAPIUrl := "https://www.googleapis.com/youtube/v3/playlistItems"
	getVideoAPIUrl := "https://www.googleapis.com/youtube/v3/videos"
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

	var nextPageToken string
	var resultList []model.ResultVideo

	for {
		var apiPlayListItemListResponse model.ApiPlayListItemListResponseWithItems
		resp, err := http.Get(getPlayListAPIUrl + "?key=" + config.YoutubeApiKey + "&part=snippet&playlistId=" + youtubeListId + "&pageToken=" + nextPageToken + "&maxResults=50")
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if resp.StatusCode != 200 {
			fmt.Println("ERROR: StatusCode:" + strconv.Itoa(resp.StatusCode))
			return
		}

		//fmt.Println(getPlayListAPIUrl + "?key=" + config.YoutubeApiKey + "&part=snippet&playlistId=" + youtubeListId + "&pageToken=" + nextPageToken)

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

		//fmt.Printf("(%%#v) %#v\n", apiPlayListItemListResponse)

		for i := 0; i < len(apiPlayListItemListResponse.Items); i++ {
			fmt.Println(getVideoAPIUrl + "?key=" + config.YoutubeApiKey + "&part=statistics&id=" + apiPlayListItemListResponse.Items[i].Snippet.ResourceId.VideoId)
			resp, err := http.Get(getVideoAPIUrl + "?key=" + config.YoutubeApiKey + "&part=statistics&id=" + apiPlayListItemListResponse.Items[i].Snippet.ResourceId.VideoId)
			if resp != nil {
				defer resp.Body.Close()
			}
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if resp.StatusCode != 200 {
				fmt.Println("ERROR: StatusCode:" + strconv.Itoa(resp.StatusCode))
				return
			}

			byteArray, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			var apivideoListResponseWithStatistics model.ApivideoListResponseWithStatistics

			err = json.Unmarshal(byteArray, &apivideoListResponseWithStatistics)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			result, err := convertResponseToResult(apiPlayListItemListResponse.Items[i].Snippet, apivideoListResponseWithStatistics)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			resultList = append(resultList, result)
			//fmt.Printf("(%%#v) %#v\n", apivideoListResponseWithStatistics)
		}
		if apiPlayListItemListResponse.NextPageToken == "" {
			break
		} else {
			nextPageToken = apiPlayListItemListResponse.NextPageToken
		}
	}
	fmt.Printf("(%%#v) %#v\n", resultList)
}

func convertResponseToResult(snippet model.Snippet, video model.ApivideoListResponseWithStatistics) (res model.ResultVideo, err error) {
	viewCount, err := strconv.Atoi(video.Items[0].Statistics.ViewCount)
	if err != nil {
		return model.ResultVideo{}, err
	}
	likeCount, err := strconv.Atoi(video.Items[0].Statistics.LikeCount)
	if err != nil {
		return model.ResultVideo{}, err
	}
	commentCount, err := strconv.Atoi(video.Items[0].Statistics.CommentCount)
	res = model.ResultVideo{
		VideoId:      snippet.ResourceId.VideoId,
		Title:        snippet.Title,
		ViewCount:    viewCount,
		LikeCount:    likeCount,
		CommentCount: commentCount,
	}
	return res, err
}
