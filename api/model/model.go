package model

type Config struct {
	YoutubeApiKey string `json:"youtubeApiKey"`
}
type ApiPlayListItemListResponseWithItems struct {
	Etag          string      `json:"etag"`
	NextPageToken string      `json:"nextPageToken"`
	Items         []ListItems `json:"items"`
}

type ListItems struct {
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

type ApivideoListResponseWithStatistics struct {
	Etag  string        `json:"etag"`
	Items []VideosItems `json:"items"`
}

type VideosItems struct {
	Etag       string     `json:"etag"`
	Statistics Statistics `json:"statistics"`
}

type Statistics struct {
	ViewCount string `json:"viewCount"`
	LikeCount string `json:"likeCount"`
	//FavoriteCount string `json:"favoriteCount"`
	CommentCount string `json:"commentCount"`
}

type ResultVideo struct {
	VideoId      string
	Title        string
	ViewCount    int
	LikeCount    int
	CommentCount int
}
