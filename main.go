package main

import (
	// "bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {

	start := time.Now()
	// Код для измерения
	duration := time.Since(start)
	counter := 0
	commonReq := 500
	for i := 0; i < commonReq; i++ {

		url := "https://portal.optic-center.ru/rest/475/bkittz3zhrr86z1o/crm.contact.list?filter[ORIGIN_ID]=3a33dfda-b969-11e3-81c1-00269e587e49"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal("Error reading request. ", err)
		}

		req.Header.Set("Content-Type", "application/text")
		req.Header.Set("Host", "localhost")

		client := &http.Client{Timeout: time.Minute * 1}

		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error reading response. ", err)
			return
		}

		counter++

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading body. ", err)
		}
		if resp.StatusCode != 200 {
			fmt.Printf("%s\n", body)
		}

		log.Printf("counter %d, status %d", counter, resp.StatusCode)

		if counter%100 == 0 {
			log.Println("Sleep")
			time.Sleep(time.Second * 10)
			log.Println("Wake")
		}
	}

	fmt.Println(duration)
	// fmt.Println("Lost package ", counter)
}
