package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const csvBaseURL = "https://www.mhlw.go.jp/content/"

// Daily 日付ごとの件数
type Daily struct {
	Date  string
	Count int
}

// NI225 日経225インデックス
type NI225 struct {
	Date  string
	Price float64
}

// Record JOIN後のレコード
type Record struct {
	Date     string
	Tested   int
	Positive int
	NI225    float64
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)

	if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		log.Panicln(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	tested, err := loadDailyCountCSV("pcr_tested_daily.csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	positive, err := loadDailyCountCSV("pcr_positive_daily.csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	ni225, err := loadNI225()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var records []Record
	for _, t := range tested {
		record := Record{
			Date:   t.Date,
			Tested: t.Count,
		}

		// Look up the positive count
		for _, p := range positive {
			if p.Date == record.Date {
				record.Positive = p.Count
			}
		}
		// Look up the Nikkei225 price
		for _, n := range ni225 {
			if n.Date == record.Date {
				record.NI225 = n.Price
			}
		}

		records = append(records, record)
	}

	if err := json.NewEncoder(w).Encode(records); err != nil {
		log.Println(err)
		return
	}
}

func loadDailyCountCSV(filename string) ([]Daily, error) {
	resp, err := http.Get(csvBaseURL + filename)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	// Drop a header record
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var ds []Daily
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return ds, nil
		}
		if err != nil {
			return nil, err
		}
		count, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}
		ds = append(ds, Daily{
			Date:  record[0],
			Count: count,
		})
	}
}

func loadNI225() ([]NI225, error) {
	file, err := os.Open("ni225.csv")
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(file)
	// Drop a header record
	if _, err := reader.Read(); err != nil {
		return nil, err
	}
	var ns []NI225
	for {
		record, err := reader.Read()
		if err == io.EOF {
			return ns, nil
		}
		if err != nil {
			return nil, err
		}
		price, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}
		ns = append(ns, NI225{
			Date:  record[0],
			Price: price,
		})
	}
}
