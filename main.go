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

// Daily 日付ごとの件数
type Daily struct {
	Date  string
	Count int
}

// Record JOIN後のレコード
type Record struct {
	Date     string
	Tested   int
	Positive int
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

		records = append(records, record)
	}

	if err := json.NewEncoder(w).Encode(records); err != nil {
		log.Println(err)
		return
	}
}

func loadDailyCountCSV(filename string) ([]Daily, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
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
