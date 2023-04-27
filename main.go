package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/avast/retry-go"
)

const (
	apiKey      = "A9fcb0683b9a9a8d43539e4f6fc397ff"
	latMarkyMoo = "43.8563707"
	lonMarkyMoo = "-79.3376825"
	latLon      = "42.9832406"
	lonLon      = "-81.243372"
	baseUrl     = "http://api.openweathermap.org/data/2.5/weather"
)

type ApiData struct {
	Weather []struct {
		Description string `json:"description"`
	}
	Main struct {
		Temp float32 `json:"temp"`
	}
}

func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func httpRequest(client *http.Client, method string, url string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func main() {
	url := fmt.Sprintf("%s?lat=%s&lon=%s&appid=%s&units=metric", baseUrl, latMarkyMoo, lonMarkyMoo, apiKey)
	c := httpClient()
	var data ApiData

	err := retry.Do(
		func() error {
			res, err := httpRequest(c, "GET", url)
			if err != nil {
				fmt.Println("something went wrong with fetch")
				return err
			}
			err = json.Unmarshal(res, &data)
			if err != nil {
				fmt.Println("something went wrong")
				return err
			}

			return nil
		},
	)

	if err != nil {
		fmt.Println("something went wrong")
		return
	}

	temp := fmt.Sprintf("%v° %s", int(data.Main.Temp), data.Weather[0].Description)
	fmt.Println(temp)
}
