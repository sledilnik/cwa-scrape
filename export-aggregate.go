package main

import (
	"container/ring"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

var (
	path = flag.String("path", "data/SI", "path to the export json files")
)

// ExposureNotificationExport was auto-generated using https://mholt.github.io/json-to-go/
type ExposureNotificationExport struct {
	StartTimestamp int    `json:"start_timestamp"`
	EndTimestamp   int    `json:"end_timestamp"`
	Region         string `json:"region"`
	BatchNum       int    `json:"batch_num"`
	BatchSize      int    `json:"batch_size"`
	SignatureInfos []struct {
		VerificationKeyVersion string `json:"verification_key_version"`
		VerificationKeyID      string `json:"verification_key_id"`
		SignatureAlgorithm     string `json:"signature_algorithm"`
	} `json:"signature_infos"`
	Keys []struct {
		KeyData                    string `json:"key_data"`
		TransmissionRiskLevel      int    `json:"transmission_risk_level"`
		RollingStartIntervalNumber int    `json:"rolling_start_interval_number"`
		RollingPeriod              int    `json:"rolling_period"`
	} `json:"keys"`
}

type DailyKeyCount struct {
	Date             string `json:"date" csv:"date"`
	NewKeysCount     int    `json:"new_key_cout" csv:"new_key_cout"`
	KeysInLast14Days int    `json:"keys_in_last_14_days" csv:"keys_in_last_14_days"`
}

// InitialDailyNewKeyCounts populated with initial reconstructed data from before production and scraping
var InitialDailyNewKeyCounts = map[string]int{
	"2020-08-10": 1,
	"2020-08-12": 9,
	"2020-08-13": 1,
}

func getDailyNewKeyCount(date string) (int, error) {
	fileName := fmt.Sprintf("%s/%s.json", *path, date)
	export, err := readExportJSON(fileName)
	if err != nil {
		if c, ok := InitialDailyNewKeyCounts[date]; ok {
			return c, nil
		}

		return 0, err
	}

	return len(export.Keys), nil
}

func readExportJSON(fileName string) (*ExposureNotificationExport, error) {

	blob, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	data := ExposureNotificationExport{}

	err = json.Unmarshal([]byte(blob), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func getDailyKeyCounts() []DailyKeyCount {
	loc, err := time.LoadLocation("Europe/Ljubljana")
	if err != nil {
		panic(err)
	}
	startDate := time.Date(2020, 8, 10, 0, 0, 0, 0, loc)
	dailyKeyCounts := make([]DailyKeyCount, 0)

	newKeysInLast14days := ring.New(14)
	for i := 0; i < newKeysInLast14days.Len(); i++ {
		newKeysInLast14days.Value = 0
		newKeysInLast14days = newKeysInLast14days.Next()
	}

	date := startDate
	for {
		dateIso := date.Format("2006-01-02")

		fmt.Println("Counting keys on:", dateIso)

		n, err := getDailyNewKeyCount(dateIso)
		if err != nil {

		}

		newKeysInLast14days.Value = n
		sum := 0
		newKeysInLast14days.Do(func(p interface{}) {
			sum += p.(int)
		})

		dailyKeyCounts = append(dailyKeyCounts, DailyKeyCount{
			Date:             dateIso,
			NewKeysCount:     n,
			KeysInLast14Days: sum,
		})

		date = date.Local().AddDate(0, 0, 1)
		newKeysInLast14days = newKeysInLast14days.Next()
		if date.After(time.Now()) {
			break
		}
	}
	return dailyKeyCounts
}

func writeJSON(data interface{}, fileName string) {
	jsonBlob, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(fileName, jsonBlob, 0644)
}

func writeCSV(data interface{}, fileName string) {
	csvFile, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	err = gocsv.MarshalFile(data, csvFile)
	if err != nil {
		panic(err)
	}
}

func main() {
	dailyKeyCounts := getDailyKeyCounts()
	writeJSON(dailyKeyCounts, *path+"/keycount.json")
	writeCSV(dailyKeyCounts, *path+"/keycount.csv")
}
