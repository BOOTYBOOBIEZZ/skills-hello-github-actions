package main

import (
	"encoding/json"
	"fmt"
	"io"
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
	server := Endpoint{Address: "127.0.0.1:8081"}

	http.HandleFunc("/create", server.HandleCreateFile)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println(err)
	}

}

type CreateFileRequest struct {
	FileName string `json:"file_name"`
}

type Endpoint struct {
	Address string
}

func (endpoint *Endpoint) HandleCreateFile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	jsonBody, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to read request body"))
		return
	}

	a := string(jsonBody)
	fmt.Println(a)

	request := CreateFileRequest{}
	jsonErr := json.Unmarshal(jsonBody, &request)
	if jsonErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, fileErr := os.Create(request.FileName + ".txt")
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
