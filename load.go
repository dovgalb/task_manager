// TODO load.go исключительно для локального теста рпс, позже удалить
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const url = "http://localhost:8082/user"

// RequestPayload описывает тело запроса
type RequestPayload struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func sendRequest(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(200 * time.Millisecond)
	payload := RequestPayload{
		Login:    "dovgalb",
		Password: "qwerty",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Ошибка сериализации JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Ответ код:", resp.StatusCode)
}

func main() {
	start := time.Now()
	var wg sync.WaitGroup

	for range 5 {
		time.Sleep(time.Second * 1)
		for range 100 {
			wg.Add(1)
			go sendRequest(&wg)
		}
		fmt.Println("цикл\n")
	}
	wg.Wait()

	fmt.Println(time.Since(start))
}
