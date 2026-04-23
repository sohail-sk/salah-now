package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func fetchFromAPI(lat, lng float64) (map[string]string, error) {
	url := fmt.Sprintf(
		"https://api.aladhan.com/v1/timings?latitude=%f&longitude=%f&method=1",
		lat, lng,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	data := result["data"].(map[string]interface{})
	timings := data["timings"].(map[string]interface{})

	return map[string]string{
		"fajr":    timings["Fajr"].(string),
		"dhuhr":   timings["Dhuhr"].(string),
		"asr":     timings["Asr"].(string),
		"maghrib": timings["Maghrib"].(string),
		"isha":    timings["Isha"].(string),
	}, nil
}

func GetPrayerTimesCached(lat, lng float64) (map[string]string, error) {

	key := fmt.Sprintf("%.4f_%.4f_%s", lat, lng, todayKey())

	// 1. Check cache
	if data, found := GetCache(key); found {
		return data, nil
	}

	// 2. Fetch from API
	data, err := fetchFromAPI(lat, lng)
	if err != nil {
		return nil, err
	}

	// 3. Store in cache
	SetCache(key, data)

	return data, nil
}

func todayKey() string {
	return time.Now().Format("2006-01-02")
}

func GetCurrentAndNext(times map[string]string) (string, string) {
	now := time.Now()

	parse := func(t string) time.Time {
		parsed, _ := time.Parse("15:04", t)
		return time.Date(now.Year(), now.Month(), now.Day(),
			parsed.Hour(), parsed.Minute(), 0, 0, now.Location())
	}

	fajr := parse(times["fajr"])
	dhuhr := parse(times["dhuhr"])
	asr := parse(times["asr"])
	maghrib := parse(times["maghrib"])
	isha := parse(times["isha"])

	switch {
	case now.Before(fajr):
		return "Isha", "Fajr"
	case now.Before(dhuhr):
		return "Fajr", "Dhuhr"
	case now.Before(asr):
		return "Dhuhr", "Asr"
	case now.Before(maghrib):
		return "Asr", "Maghrib"
	case now.Before(isha):
		return "Maghrib", "Isha"
	default:
		return "Isha", "Fajr"
	}
}
