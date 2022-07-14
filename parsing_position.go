package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type people struct {
	Number int `json:"number"`
}

func parse(w http.ResponseWriter, r *http.Request) {
	url := "https://search.wb.ru/exactmatch/ru/common/v4/search?appType=1&couponsGeo=12,7,3,6,18,22,21&curr=rub&dest=-1075831,-79374,-367666,-2133466&emp=0&lang=ru&locale=ru&pricemarginCoeff=1.0&reg=1&regions=68,64,83,4,38,80,33,70,82,86,30,69,22,66,31,40,1,48&resultset=catalog&sort=popular&spp=19"
	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}
	counter := 1
	target_id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	query := r.URL.Query().Get("query")
	fmt.Println(query)
	for page := 1; page < 50; page++ {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		q := req.URL.Query()
		q.Add("page", fmt.Sprintf("%d", page))
		q.Add("query",  query)
		req.URL.RawQuery = q.Encode()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(req.URL.String())

		req.Header.Set("User-Agent", "spacecount-tutorial")

		res, getErr := spaceClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		var dat map[string]interface{}
		json.Unmarshal(body, &dat)
		products := dat["data"].(map[string]interface{})["products"].([]interface{})

		for product := range products {
			id := int(products[product].(map[string]interface{})["id"].(float64))
			if id != target_id {
				counter += 1
			} else {
				fmt.Println("Нашли нужный товар", counter)
				w.Write([]byte(fmt.Sprintf("Нашли нужный товар %d", counter)))
			}
		}
	}
}

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", parse)

	http.ListenAndServe(":8080", mux)
	start := time.Now()
	elapsed := time.Since(start)
	log.Printf("%s took %s", "name", elapsed)
}
