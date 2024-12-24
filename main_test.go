package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleCreateFile_Success(t *testing.T) {
	endpoint := &Endpoints{}

	//Podgotovka zaprosa
	requestBody := CreateFileRequest{
		FileName: "hello.txt",
		Payload:  "HeyHey fella",
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	//vipolnenie obrabotchika
	endpoint.HandleCreateFile(w, req)

	// check answer
	resp := w.Result()

	//тестовый запрос
	payload := []byte("{\"file_name\":\"testfile\"}")
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()

	//вызываем обработчик
	handler := http.HandlerFunc(server.HandleCreateFile)
	handler.ServeHTTP(rr, req)

	//проверяем ответ
	assert.Equal(t, http.StatusCreated, rr.Code, "Expected 201 Created")

	//проверяем что файл создан
	_, err := os.Stat("testfile.txt")
	assert.NoError(t, err, "File should have been created")

	//Ydalyaem test file
	_ = os.Remove("testfile.txt")

}

func TestHandleCreateFile_InvalidJSON(t *testing.T) {
	server := Endpoint{}

	//nekorrektnyi JSON
	payload := []byte("{invalid}")
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()

	//vizivaem obrabotchik
	handler := http.HandlerFunc(server.HandleCreateFile)
	handler.ServeHTTP(rr, req)

	//Proveryaem otvet
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request")
	assert.Contains(t, rr.Body.String(), "Invalid JSON format", "Expected error message about JSON")

}

func TestHandleCreateFile_FileCreationError(t *testing.T) {
	server := Endpoint{}

	//nekorrektnoe imya file'a

	payload := []byte("{\"file_name\":\"/invalid/path/testfile\"}")
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()

	//vyzyvaem obrabotchik
	handler := http.HandlerFunc(server.HandleCreateFile)
	handler.ServeHTTP(rr, req)

	//check answer
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request")

}

func TestHandleCreateFile_EmptyFileName(t *testing.T) {
	server := Endpoint{}

	//empty file name
	payload := []byte("{\"file_name\":\"\"}")
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()

	//call obrabotchik
	handler := http.HandlerFunc(server.HandleCreateFile)
	handler.ServeHTTP(rr, req)

	//Check answer
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request")
}
