package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
)

func main() {
	viaCEP := make(chan string)
	brasilAPI := make(chan string)
	go func() {
		defer close(viaCEP)
		resp, err := http.Get("https://viacep.com.br/ws/01153000/json/")
		if err != nil {
			log.Error(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		}
		viaCEP <- string(body)
	}()
	go func() {
		defer close(brasilAPI)
		resp, err := http.Get("https://brasilapi.com.br/api/cep/v1/01153000")
		if err != nil {
			log.Error(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		}
		brasilAPI <- string(body)
	}()
	select {
	case data := <-viaCEP:
		fmt.Println("ViaCEP:", data)
	case data := <-brasilAPI:
		fmt.Println("BrasilAPI:", data)
	case <-time.After(1 * time.Second):
		log.Error("timeout")
	}
}
