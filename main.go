package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Person struct {
	Err string
}

func (person *Person) Error() string {
	return person.Err
}

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	endpoint := Endpoints{}

	http.HandleFunc("/create", endpoint.HandleCreateFile)
	http.HandleFunc("/update", endpoint.HandleUpdateFile)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println(err)
	}
}

type CreateFileRequest struct {
	FileName string `json:"file_name"`
	Payload  string `json:"payload"`
}

type Endpoints struct {
}

func (endpoint *Endpoints) HandleUpdateFile(w http.ResponseWriter, r *http.Request) {
	request := CreateFileRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Failed to update file"))
		return
	}
	fileErr := os.WriteFile(request.FileName, []byte(request.Payload), 0644)
	if fileErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to update file"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Successfully updated file"))
	return
}

func (endpoint *Endpoints) HandleCreateFile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	request := CreateFileRequest{}

	jsonErr := json.NewDecoder(r.Body).Decode(&request)
	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to read request body"))
		return
	}

	fileName := request.FileName + ".txt"

	if len(request.FileName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("file name is required"))
		return
	}

	_, fileErr := os.Create(fileName)
	if fileErr != nil {
		w.WriteHeader(http.StatusBadRequest)

		_, responseErr := w.Write([]byte(fileErr.Error()))
		if responseErr != nil {
			fmt.Printf("filed create file in CreateFile endpoint %s", responseErr.Error())
			return
		}

		return
	}

	w.WriteHeader(http.StatusCreated)
}
