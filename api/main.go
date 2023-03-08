package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"youtube-app/api/model"

	"golang.org/x/exp/slices"
)

func main() {
	var youtubeListId string
	var config model.Config

	// 引数チェック
	if len(os.Args) != 2 {
		fmt.Println("Invalid args. Usage: go run main.go {PlaylistId}")
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
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		fmt.Println(err.Error())
		return
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

		byteArray, err := io.ReadAll(resp.Body)
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

		if len(apiPlayListItemListResponse.Items) == 0 {
			fmt.Println("ERROR: Response count is 0")
			return
		}

		var itemsListParameter string

		for i := 0; i < len(apiPlayListItemListResponse.Items); i++ {
			if i == len(apiPlayListItemListResponse.Items)-1 {
				itemsListParameter = itemsListParameter + apiPlayListItemListResponse.Items[i].Snippet.ResourceId.VideoId
			} else {
				itemsListParameter = itemsListParameter + apiPlayListItemListResponse.Items[i].Snippet.ResourceId.VideoId + ","
			}
		}

		resp, err = http.Get(getVideoAPIUrl + "?key=" + config.YoutubeApiKey + "&part=statistics&id=" + itemsListParameter)

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

		byteArray, err = io.ReadAll(resp.Body)
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

		result, err := mergeAPIResponseToResult(apiPlayListItemListResponse.Items, apivideoListResponseWithStatistics.Items)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		resultList = append(resultList, result...)
		if apiPlayListItemListResponse.NextPageToken == "" {
			break
		} else {
			nextPageToken = apiPlayListItemListResponse.NextPageToken
		}
	}

	// sort
	sort.SliceStable(resultList, func(i, j int) bool { return resultList[i].ViewCount > resultList[j].ViewCount })
	fmt.Printf("(%%#v) %#v\n", resultList)
}

func mergeAPIResponseToResult(listItems []model.ListItems, videoItems []model.VideosItems) ([]model.ResultVideo, error) {
	var result []model.ResultVideo
	for i := 0; i < len(listItems); i++ {
		// search videoItems index
		idx := slices.IndexFunc(videoItems, func(item model.VideosItems) bool { return item.VideoId == listItems[i].Snippet.ResourceId.VideoId })
		if idx == -1 {
			return []model.ResultVideo{}, errors.New("videoNotFoundError")
		}

		// convert
		viewCount, err := strconv.Atoi(videoItems[idx].Statistics.ViewCount)
		if err != nil {
			return []model.ResultVideo{}, err
		}
		likeCount, err := strconv.Atoi(videoItems[idx].Statistics.LikeCount)
		if err != nil {
			return []model.ResultVideo{}, err
		}
		commentCount, err := strconv.Atoi(videoItems[idx].Statistics.CommentCount)
		if err != nil {
			return []model.ResultVideo{}, err
		}
		res := model.ResultVideo{
			VideoId:      listItems[i].Snippet.ResourceId.VideoId,
			Title:        listItems[i].Snippet.Title,
			ViewCount:    viewCount,
			LikeCount:    likeCount,
			CommentCount: commentCount,
		}
		result = append(result, res)
	}
	return result, nil
}
