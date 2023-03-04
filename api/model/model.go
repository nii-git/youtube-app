package model

type Config struct {
	YoutubeApiKey string `json:"youtubeApiKey"`
}
type ApiPlayListItemListResponse struct {
	Etag          string  `json:"etag"`
	NextPageToken string  `json:"nextPageToken"`
	Items         []Items `json:"items"`
}

type Items struct {
	Etag    string  `json:"etag"`
	Id      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type Snippet struct {
	Title      string     `json:"title"`
	ResourceId ResourceId `json:"resourceId"`
}

type ResourceId struct {
	Kind    string `json:"kind"`
	VideoId string `json:"videoId"`
}
