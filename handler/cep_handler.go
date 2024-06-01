package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	apicep = iota
	viacep
)

type Cep struct {
	Cep        string `json:cep`
	Logradouro string `json:rua`
	Cidade     string `json:cidade`
	Bairro     string `json:bairro`
	UF         string `json:uf`
	Tipo       string `json:"tipo"`
}

type CepResponse struct {
	Cep    string
	Street string
	State  string
	City   string
	Tipo   string
}

type ApicepStruct struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

type ViacepStruct struct {
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

func GetCepValue(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	apicepUrl := "https://cdn.apicep.com/file/apicep/" + cep + ".json"
	viacepUrl := "http://viacep.com.br/ws/" + cep + "/json/"

	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ch1 := make(chan Cep)
	ch2 := make(chan Cep)

	go func() {
		req, err := http.Get(apicepUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao buscar determnado cep: %s, erro: %v", cep, err)
			ch1 <- Cep{}
			return
		}
		defer req.Body.Close()

		res, err := io.ReadAll(req.Body)

		var apiResponse ApicepStruct
		err = json.Unmarshal(res, &apiResponse)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao buscar determnado cep: %s, erro: %v", cep, err)
		}

		response := Cep{
			Cep:        apiResponse.Code,
			Logradouro: apiResponse.Address,
			Cidade:     apiResponse.City,
			Bairro:     apiResponse.District,
			UF:         apiResponse.State,
			Tipo:       "apicep",
		}

		ch1 <- response
	}()

	go func() {
		req, err := http.Get(viacepUrl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao buscar cep: %s, erro: %v", cep, err)
		}

		defer req.Body.Close()
		res, err := io.ReadAll(req.Body)

		var apiResponse ViacepStruct
		err = json.Unmarshal(res, &apiResponse)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao buscar cep: %s, erro: %v", cep, err)
		}

		response := Cep{
			Cep:        apiResponse.Cep,
			Logradouro: apiResponse.Logradouro,
			Cidade:     apiResponse.Localidade,
			Bairro:     apiResponse.Bairro,
			UF:         apiResponse.Uf,
			Tipo:       "viacep",
		}
		fmt.Println(response)

		ch2 <- response
	}()

	tipo := ""
	var response CepResponse

	select {
	case msg1 := <-ch1:
		tipo = "apicep"

		response.Cep = msg1.Cep
		response.Street = msg1.Logradouro
		response.State = msg1.UF
		response.City = msg1.Cidade
		response.Tipo = tipo

	case msg2 := <-ch2:
		tipo = "viacep"
		response.Cep = msg2.Cep
		response.Street = msg2.Logradouro
		response.State = msg2.UF
		response.City = msg2.Cidade
		response.Tipo = tipo

	case <-time.After(time.Second * 1):
		fmt.Printf("timeout \n")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// response := CepResponse{
	// 	Cep:  cep,
	// 	Tipo: tipo,
	// }

	json.NewEncoder(w).Encode(response)
}
