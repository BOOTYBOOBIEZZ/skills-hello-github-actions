package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	http.HandleFunc("/read", endpoint.HandleReadFile)
	http.HandleFunc("/delete", endpoint.HandleDeleteFile)

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

func (endpoint *Endpoints) HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	request := CreateFileRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Failed to delete file"))
		return
	}

	err = os.Remove(request.FileName)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to delete file"))
		return
	}

	w.Write([]byte("Successfully deleted file"))
	w.WriteHeader(http.StatusOK)
	return

}

func (endpoint *Endpoints) HandleReadFile(w http.ResponseWriter, r *http.Request) {
	request := CreateFileRequest{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
		w.Write([]byte("Failed to read file"))
		return
	}

	data, err := os.ReadFile(request.FileName)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to read file"))
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
	return
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

	if r.Body == nil {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
	}

	request := CreateFileRequest{}
	jsonErr := json.NewDecoder(r.Body).Decode(&request)
	if jsonErr != nil {
		log.Println("Error decoding JSON:", jsonErr)
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	fileName := request.FileName + ".txt"

	if len(request.FileName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("file name is required"))
		return
	}

	file, fileErr := os.Create(fileName)
	if fileErr != nil {
		log.Println("Error creating file:", fileErr)
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}

	_, writeErr := file.WriteString(request.Payload)
	if writeErr != nil {
		log.Println("Error writing to file:", writeErr)
		http.Error(w, "Failed to write to file", http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File created successfully"))

}
