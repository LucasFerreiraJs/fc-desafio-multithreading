package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type CepResponse struct {
	Cep  string
	Tipo int
}

// func GetCepValue(w http.ResponseWriter,r &http.Request) {
// 	cep := chi.URLParam(r, "cep")
// 	if cep == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}

//		w.Header().Set("Content-Type", "application/json")
//	}
func GetCepValue(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")

	if cep == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("cep recebido: %s \n", cep)
	fmt.Println()

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		msg := "consulta 01"
		ch1 <- msg
	}()

	go func() {
		msg := "consulta 02"
		time.Sleep(time.Second * 3)
		ch2 <- msg
	}()
	tipo := 0
	select {
	case msg1 := <-ch1:
		fmt.Printf("msg 01 %s \n", msg1)
		tipo = 1
	case msg2 := <-ch2:
		fmt.Printf("msg 02 %s \n", msg2)
		tipo = 2
	case <-time.After(time.Second * 2):
		fmt.Printf("timeout \n")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := CepResponse{
		Cep:  cep,
		Tipo: tipo,
	}

	json.NewEncoder(w).Encode(response)
}
