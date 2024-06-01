# desafio-multithreading


### Rodando o projeto

<!-- Precisará rodar duas aplicações:
-   server.go
-   client.go -->

Execute o comando na raiz do projeto:

```
go mod tidy

```

Execute o comando na pasta ./server:
```
go run main.go

```




Faça uma requisição com o valor do cep que deseja consultar
```
curl http://127.0.0.1:8000/consulta-cep/xxxxxxxx

```

Deverá retornar informações do cep consultado junto com a api que retornou o valor mais rápido

