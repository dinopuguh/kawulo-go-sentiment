package services

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"time"
)

var (
	BaseUrl   = "https://translate.yandex.net/api/v1.5/tr.json/translate"
	transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
	}
	res     *http.Response
	retries int
)

type TranslateResponse struct {
	Code int32    `json:"code"`
	Lang string   `json:"lang"`
	Text []string `json:"text"`
}

func GetSourceLang(srcLang string) string {
	langs := map[string]string{
		"in":   "id",
		"zhCN": "zh",
		"zhTW": "zh",
		"iw":   "he",
		"aeAE": "ar",
		"enAU": "en",
		"enCA": "en",
		"enHK": "en",
		"enIN": "en",
		"enIE": "en",
		"enMY": "en",
		"enNZ": "en",
		"enPH": "en",
		"enSG": "en",
		"enZA": "en",
		"enUK": "en",
		"frBE": "fr",
		"frCA": "fr",
		"frCH": "fr",
		"deAT": "de",
		"itCH": "it",
		"ptPT": "pt",
		"esAR": "es",
		"esCO": "es",
		"esMX": "es",
		"esPE": "es",
		"esVE": "es",
		"esCL": "es",
	}

	return langs[srcLang]
}

func TranslateReview(text string, srcLang string, APIKey string) string {
	targetLang := "en"
	lang := targetLang

	if srcLang != "" {
		tripadvisorLang := GetSourceLang(srcLang)
		if tripadvisorLang != "" {
			srcLang = tripadvisorLang
		}
		lang = srcLang + "-" + targetLang
	}

	client := &http.Client{
		Transport: transport,
	}

	var data TranslateResponse
	retries = 3

	req, err := http.NewRequest("POST", BaseUrl, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	q := req.URL.Query()
	q.Add("key", APIKey)
	q.Add("lang", lang)
	q.Add("text", text)
	req.URL.RawQuery = q.Encode()

	req.Close = true

	for retries > 0 {
		res, err = client.Do(req)

		if err != nil {
			retries -= 1
			log.Println("Retrying...")
		} else {
			defer res.Body.Close()
			break
		}
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		log.Fatal(err.Error())
	}

	return data.Text[0]
}
