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

	loc, _ = time.LoadLocation("Europe/Ljubljana")
)

const isoDateFormat = "2006-01-02"

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
	Keys []ExposureNotificationExportKey `json:"keys"`
}

type ExposureNotificationExportKey struct {
	KeyData                    string `json:"key_data"`
	TransmissionRiskLevel      int    `json:"transmission_risk_level"`
	RollingStartIntervalNumber int    `json:"rolling_start_interval_number"`
	RollingPeriod              int    `json:"rolling_period"`
}

type DailyKeyCount struct {
	Date                string                          `json:"date" csv:"date"`
	NewKeysCount        int                             `json:"new_key_count" csv:"new_key_count"`
	KeysTotal           int                             `json:"total_keys" csv:"total_keys"`
	NewKeysInLast14Days int                             `json:"new_keys_in_last_14_days" csv:"new_keys_in_last_14_days"`
	NonExpiredKeys      int                             `json:"non_expired_keys" csv:"non_expired_keys"`
	Keys                []ExposureNotificationExportKey `json:"-" csv:"-"`
}


func getDailyNewKeyCount(date string) ([]ExposureNotificationExportKey, error) {
	fileName := fmt.Sprintf("%s/%s.json", *path, date)
	export, err := readExportJSON(fileName)
	if err != nil {
		return nil, err
	}

	return export.Keys, nil
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
	startDate := time.Date(2020, 8, 27, 0, 0, 0, 0, loc)
	dailyKeyCounts := make([]DailyKeyCount, 0)

	newKeysInLast14days := ring.New(14)
	for i := 0; i < newKeysInLast14days.Len(); i++ {
		newKeysInLast14days.Value = make([]ExposureNotificationExportKey, 0)
		newKeysInLast14days = newKeysInLast14days.Next()
	}

	date := startDate
	runningTotal := 0
	for {
		dateIso := date.Format(isoDateFormat)

		fmt.Println("Counting keys on:", dateIso)

		dailyKeys, err := getDailyNewKeyCount(dateIso)
		if err != nil {

		}

		n := len(dailyKeys)
		runningTotal += n

		newKeysInLast14days.Value = dailyKeys
		sum := 0
		nonExpiredKeys := 0
		twoWeeksAgo := date.AddDate(0, 0, -14)
		newKeysInLast14days.Do(func(p interface{}) {
			keys := p.([]ExposureNotificationExportKey)
			sum += len(keys)

			for _, k := range keys {
				if k.getStart().After(twoWeeksAgo) {
					nonExpiredKeys++
					// fmt.Println("active:", k.getIsoDate(), k.getStart(), k)
				}
			}
			fmt.Println(dateIso, twoWeeksAgo, ":", nonExpiredKeys, "of", sum)
		})

		dailyKeyCounts = append(dailyKeyCounts, DailyKeyCount{
			Date:                dateIso,
			NewKeysCount:        n,
			KeysTotal:           runningTotal,
			NewKeysInLast14Days: sum,
			NonExpiredKeys:      nonExpiredKeys,
		})

		date = date.AddDate(0, 0, 1)
		newKeysInLast14days = newKeysInLast14days.Next()
		if date.After(time.Now().In(loc)) {
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

func (k *ExposureNotificationExportKey) getStart() time.Time {
	return time.Unix(int64(k.RollingStartIntervalNumber)*600, 0) // 10-minute slot since unix epoch
}

func (k *ExposureNotificationExportKey) getIsoDate() string {
	return k.getStart().Format(isoDateFormat)
}

func main() {
	flag.Parse()

	dailyKeyCounts := getDailyKeyCounts()
	writeJSON(dailyKeyCounts, *path+"/keycount.json")
	writeCSV(dailyKeyCounts, *path+"/keycount.csv")
	writeChart(dailyKeyCounts, *path+"/keycount.png")
}
