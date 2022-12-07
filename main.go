package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/JPGSilva/GO_CRUD/model"
)

func main() {

	service, err := regras.NewService("regras.json")
	if err != nil {
		fmt.Printf("Erro ao Criar Serviço de Aluno: %s\n", err.Error())
	}

	http.HandleFunc("/aluno/", func(resposta http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			path := strings.TrimPrefix(req.URL.Path, "/aluno/")
			if path == "" {

				resposta.WriteHeader(http.StatusOK)
				resposta.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(resposta).Encode(service.List())

				if err != nil {
					http.Error(resposta, "Erro ao Listar Alunos", http.StatusInternalServerError)
					return
				}
			} else {
				idAluno, err := strconv.Atoi(path)
				if err != nil {
					http.Error(resposta, "Por Favor informe apenas Números Inteiros", http.StatusBadRequest)
					return
				}
				aluno, err := service.GetByID(idAluno)
				if err != nil {
					http.Error(resposta, err.Error(), http.StatusNotFound)
					return
				}
				resposta.WriteHeader(http.StatusOK)
				resposta.Header().Set("Content-Type", "application/json")
				err = json.NewEncoder(resposta).Encode(aluno)
				if err != nil {
					http.Error(resposta, "Não Foi Possível Realizar a Busca por Esse Aluno", http.StatusInternalServerError)
					return
				}
			}
			return
		}

		if req.Method == "POST" {
			var aluno model.Aluno
			err := json.NewDecoder(req.Body).Decode(&aluno)
			if err != nil {
				fmt.Printf("Por Favor Insira no Formato Json {}: %s\n", err.Error())
				http.Error(resposta, "Não foi Possivel Criar Esse Aluno", http.StatusBadRequest)
				return
			}
			if aluno.ID <= 0 {
				http.Error(resposta, "O Id do Aluno deve ser Inteiro e Positivo", http.StatusBadRequest)
				return
			}

			err = service.Create(aluno)
			if err != nil {
				fmt.Printf("Erro ao Criar Aluno: %s\n", err.Error())
				http.Error(resposta, "Não Foi Possivel Criar o Aluno", http.StatusInternalServerError)
				return
			}
			resposta.WriteHeader(http.StatusCreated)
			return
		}

		if req.Method == "PUT" {
			var aluno model.Aluno
			err := json.NewDecoder(req.Body).Decode(&aluno)
			if err != nil {
				fmt.Printf("Por Favor Insira no Formato Json {}: %s\n", err.Error())
				http.Error(resposta, "Erro ao Tentar Atualizar Aluno :)", http.StatusBadRequest)
				return
			}
			if aluno.ID <= 0 {
				http.Error(resposta, "O Id do Aluno deve ser Inteiro e Positivo", http.StatusBadRequest)
				return
			}

			err = service.Update(aluno)
			if err != nil {
				fmt.Printf("Error ao Tentar Atualizar Informações do Aluno: %s\n", err.Error())
				http.Error(resposta, "Erro ao Tentar Atualizar Aluno :(", http.StatusInternalServerError)
				return
			}
			resposta.WriteHeader(http.StatusOK)
			return
		}
	})

	http.ListenAndServe(":8087", nil)
}
