package regras

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"https://github.com/JPGSilva/GO_CRUD"
)

type Service struct {
	dbFilePath string
	alunos model.ListaDeAlunos
}

// Criação do Arquivo para Salvar as Informações
func NewService(dbFilePath string) (Service, error) {
	_, err := os.Stat(dbFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = createEmptyFile(dbFilePath)
			if err != nil {
				return Service{}, err
			}
			return Service{
				dbFilePath: dbFilePath,
				alunos:     model.ListaDeAlunos{},
			}, nil
		} else {
			return Service{}, err
		}
	}

	jsonFile, err := os.Open(dbFilePath)
	if err != nil {
		return Service{}, fmt.Errorf("Erro Encontrado no Arquivo que possui todas as Informações: %s", err.Error())
	}

	jsonFileContentByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return Service{}, fmt.Errorf("Erro Ao Ler arquivo processado: %s", err.Error())
	}

	var totalAlunos model.ListaDeAlunos
	json.Unmarshal(jsonFileContentByte, &totalAlunos)

	return Service{
		dbFilePath: dbFilePath,
		alunos:     totalAlunos,
	}, nil
}

func (s Service) saveFile() error {
	infoAlunosJSON, err := json.Marshal(s.alunos)
	if err != nil {
		return fmt.Errorf("Erro ao tentar codificar aluno como JSON: %s", err.Error())
	}
	return ioutil.WriteFile(s.dbFilePath, infoAlunosJSON, 0755)
}

func (s *Service) AddAluno(alunoModel model.Aluno) error {
	s.alunos.ListaDeAlunos = append(s.alunos.ListaDeAlunos, alunoModel)
	return s.saveFile()
}

func (s *Service) Create(alunoModel model.Aluno) error {
	if s.exists(alunoModel) {
		return fmt.Errorf("Já Existe um Aluno cadastrado com essas Informções")
	}

	err := s.AddAluno(alunoModel)
	if err != nil {
		return fmt.Errorf("Não Foi Possível adicionar as Informações no Arquivo: %s", err.Error())
	}

	return nil
}

func (s Service) exists(alunoEnti model.Aluno) bool {
	for _, alunoInfo := range s.alunos.ListaDeAlunos {
		if alunoInfo.ID == alunoEnti.ID {
			return true
		}
	}
	return false
}

func (s Service) List() model.ListaDeAlunos {
	return s.alunos
}

func (s Service) GetByID(personID int) (model.Aluno, error) {
	for _, alunoInfo := range s.alunos.ListaDeAlunos {
		if alunoInfo.ID == personID {
			return alunoInfo, nil
		}
	}
	return model.Aluno{}, fmt.Errorf("Aluno não Existe")
}

func (s *Service) Update(alunoEnti model.Aluno) error {
	var indexToUpdate int = -1
	for index, alunoInfo := range s.alunos.ListaDeAlunos {
		if int(alunoInfo.NotaAluno) == int(alunoEnti.NotaAluno) {
			indexToUpdate = index
			break
		}
	}
	if indexToUpdate < 0 {
		return fmt.Errorf("Nota Inválida")
	}

	s.alunos.ListaDeAlunos[indexToUpdate] = alunoEnti
	return s.saveFile()
}

func createEmptyFile(dbFilePath string) error {
	var aluno model.ListaDeAlunos = model.ListaDeAlunos{
		ListaDeAlunos: []model.Aluno{},
	}
	ListAlunosJSON, err := json.Marshal(aluno)
	if err != nil {
		return fmt.Errorf("Erro ao Processar os Dados via Json: %s", err.Error())
	}

	err = ioutil.WriteFile(dbFilePath, ListAlunosJSON, 0755)
	if err != nil {
		return fmt.Errorf("Erro ao Gravar Arquivo: %s", err.Error())
	}

	return nil
}
