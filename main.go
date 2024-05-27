package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type result struct {
	data interface{}
	api  string
}

type Brasilapi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type Viacep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
	cep := "01153000"
	if len(os.Args) > 1 {
		cep = os.Args[1]
	}

	chRes := make(chan result)
	var brasilapi Brasilapi
	var viacep Viacep
	go getAPI("https://brasilapi.com.br/api/cep/v1/"+cep, &brasilapi, chRes)
	go getAPI("http://viacep.com.br/ws/"+cep+"/json/", &viacep, chRes)

	select {
	case res := <-chRes:
		fmt.Printf("API: %s\nResult: %+v\n", res.api, res.data)
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}

}

func getAPI(api string, st interface{}, chRes chan result) {
	resp, err := http.Get(api)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	var data []byte
	data, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(data, st)
	if err != nil {
		fmt.Println(err)
		return
	}
	chRes <- result{
		data: st,
		api:  api,
	}
}
