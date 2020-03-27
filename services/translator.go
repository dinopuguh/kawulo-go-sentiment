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

func GetSourceLang(src_lang string) string {
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

	return langs[src_lang]
}

func TranslateReview(text string, src_lang string, api_key string) string {
	target_lang := "en"
	lang := target_lang

	if src_lang != "" {
		tripadvisor_lang := GetSourceLang(src_lang)
		if tripadvisor_lang != "" {
			src_lang = tripadvisor_lang
		}
		lang = src_lang + "-" + target_lang
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
	q.Add("key", api_key)
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
