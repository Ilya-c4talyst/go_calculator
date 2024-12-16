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
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var request Request
	var response Response
	var responseError ErrorResponse

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error while using calculator:", err)

		responseError.Error = "InternalServerError"
		jsonData, err := json.Marshal(responseError)
		if err == nil {
			w.Write(jsonData)
		}

		return
	}

	result, err := calculator.Calc(request.Expression)

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Println("Error while using calculator:", err)

		responseError.Error = err.Error()
		jsonData, err := json.Marshal(responseError)
		if err == nil {
			w.Write(jsonData)
		}

		return
	}

	response.Result = result
	jsonData, err := json.Marshal(response)

	if err != nil {
		log.Println("Error while marshal JSON", err)
		http.Error(w, "InternalServerError", http.StatusInternalServerError)
	}

	w.Write(jsonData)
}

func (a *Application) RunServer() error {

	http.HandleFunc("/api/v1/calculate", CalcHandler)
	err := http.ListenAndServe(":"+a.port, nil)

	if err != nil {
		log.Println("Error while running server")
	} else {
		log.Println("Server started")
	}
	return err
}
