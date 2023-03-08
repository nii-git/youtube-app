

<div align="center">

# Youtube Playlist Analyze 

![](https://img.shields.io/badge/version-1.0-blue.svg)
![](https://img.shields.io/badge/golang-1.20-green.svg)
![Last commit](https://img.shields.io/github/last-commit/nii-git/youtube-app)

Youtube Playlist Analyze は再生リストから動画のデータを取得することができます。

[Getting started](#getting-started) / [Run the program](#run-the-program)

</div>

## Getting started
### Create Youtube Data API Token
[YouTube Data API の概要](https://developers.google.com/youtube/v3/getting-started?hl=ja)の`作業を始める前に`を参考に、Youtube API認証情報を取得してください。  

取得したAPI Keyを`youtube-playlist-analyze/config.json`内の`youtubeApiKey`にセットしてください。

```json
    "youtubeApiKey": "INSERT_YOUR_KEY"
```

### Get YoutubeListId
Youtubeにアクセスし、プレイリストのIDを取得します。  
Yotube再生リストURLの`list`パラメーターの値です。

下記の例では、プレイリストのIDは`samplePlaylistId`になります。
```
https://www.youtube.com/watch?v={videoid}&list={samplePlaylistId} 
```

### clone the program
```sh
git clone https://github.com/nii-git/youtube-app.git
cd youtube-app/youtube-playlist-analyze
go run main.go {youtube_list_id}
```

## Run the program
```sh
go run main.go {youtube_list_id}
```

結果は`youtube-playlist-analyze/result/YYYYMMDD_YOUTUBELISTID.csv`に出力されます。 
