package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {

	start := time.Now()
	// Код для измерения
	var wg sync.WaitGroup
	duration := time.Since(start)
	counter := 0
	commonReq := 1000
	for i := 0; i < commonReq; i++ {
		wg.Add(1)
		go func() {

			defer wg.Done()

			content := "{\"#\",8f8a65b4-94c4-4794-b3e2-800d18d503ca,151:80baa4bf015829f711e9995968a3b691}"

			url := "http://localhost:8082/client"
			data := []byte(content)

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
			if err != nil {
				log.Fatal("Error reading request. ", err)
			}

			req.Header.Set("Content-Type", "application/text")
			req.Header.Set("Host", "localhost")

			client := &http.Client{Timeout: time.Second * 10}

			resp, err := client.Do(req)
			if err != nil {
				log.Fatal("Error reading response. ", err)
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Error reading body. ", err)
			}
			if resp.StatusCode != 200 {
				fmt.Printf("%s\n", body)
			}

		}()

	}
	wg.Wait()

	fmt.Println(duration)
	fmt.Println("Lost package ", counter)
}
