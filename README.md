# covid19nikkei225

新型コロナウイルス感染症のPRC検査実施人数および陽性者数の日付ごとの件数と日経平均株価を結合した結果をJSON形式で返すWebアプリケーションです。

## Start

```shell
go run main.go
```

With auto-reload:

```shell
leaf -x 'go run main.go'
```

# Datasource

[オープンデータ｜厚生労働省](https://www.mhlw.go.jp/stf/covid-19/open-data.html)
[Nikkei 225](https://docs.google.com/spreadsheets/d/1wsokKD7g5FRZk1Un0DJILLactzvTC5dVQ-PIK6E3F0A/edit?usp=sharing)

# Deploy to Google App Engine

```shell
gcloud app deploy
```