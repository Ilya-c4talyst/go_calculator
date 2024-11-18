package application

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Ilya-c4talyst/go_calculator/pkg/calculator"
)

type Application struct {
	port string
}

func NewApp(port string) *Application {
	return &Application{port: port}
}

func (a *Application) LocalRun() error {
	log.Println("Local Run")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		inputString := scanner.Text()
		inputString = strings.TrimSpace(inputString)
		if inputString == "" || inputString == "enough" {
			log.Println("calculator said bye-bye")
			return nil
		}
		result, err := calculator.Calc(inputString)
		if err != nil {
			log.Println(err, "calculator said bye-bye")
		} else {
			log.Println(result)
		}
	}
}

type Request struct {
	Input string
}

type Response struct {
	Output float64
	Error  string
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var request Request
	var response Response

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		response.Error = err.Error()
		jsonData, err := json.Marshal(response)
		if err == nil {
			w.Write(jsonData)
		}

		log.Println("Error while using calculator", err)
		return
	}
	result, err := calculator.Calc(request.Input)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		response.Error = err.Error()
		jsonData, err := json.Marshal(response)
		if err == nil {
			w.Write(jsonData)
		}

		log.Println("Error while using calculator", err)
		return

	}
	response.Output = result
	response.Error = ""

	jsonData, err := json.Marshal(response)

	if err != nil {
		log.Println("Error while marshal JSON", err)
		http.Error(w, "Error", http.StatusBadGateway)
	}

	w.Write(jsonData)
}

func (a *Application) RunServer() error {

	http.HandleFunc("/calc/", CalcHandler)
	err := http.ListenAndServe(":"+a.port, nil)

	if err != nil {
		log.Println("Error while running server")
	} else {
		log.Println("Server started")
	}
	return err
}
