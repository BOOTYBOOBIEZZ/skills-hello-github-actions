package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Тест создания файла
func TestHandleCreateFile(t *testing.T) {
	endpoint := &Endpoints{}

	// Подготовка запроса
	requestBody := CreateFileRequest{
		FileName: "testfile",
		Payload:  "Hello, World!",
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Error marshaling request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Выполнение обработчика
	endpoint.HandleCreateFile(w, req)

	// Проверка ответа
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	// Проверка содержимого созданного файла
	content, err := os.ReadFile("testfile.txt")
	if err != nil {
		t.Fatalf("Error reading created file: %v", err)
	}
	if string(content) != requestBody.Payload {
		t.Errorf("Expected file content %q, got %q", requestBody.Payload, string(content))
	}

	// Удаление файла
	defer os.Remove("testfile.txt")
}

// Тест чтения файла
func TestHandleReadFile(t *testing.T) {
	endpoint := &Endpoints{}

	// Подготовка тестового файла
	fileName := "testread.txt"
	content := "Hello, World!"
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Error writing test file: %v", err)
	}
	defer os.Remove(fileName)

	// Подготовка запроса
	requestBody := CreateFileRequest{
		FileName: fileName,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Error marshaling request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/read", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Выполнение обработчика
	endpoint.HandleReadFile(w, req)

	// Проверка ответа
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	// Проверка содержимого ответа
	respBody := w.Body.String()
	if respBody != content {
		t.Errorf("Expected response body %q, got %q", content, respBody)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println(content)
}

// Тест обновления файла
func TestHandleUpdateFile(t *testing.T) {
	endpoint := &Endpoints{}

	// Подготовка тестового файла
	fileName := "testupdate.txt"
	initialContent := "Initial Content"
	updatedContent := "Updated Content"
	err := os.WriteFile(fileName, []byte(initialContent), 0644)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	defer os.Remove(fileName)

	// Подготовка запроса
	requestBody := CreateFileRequest{
		FileName: fileName,
		Payload:  updatedContent,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Error marshaling request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Выполнение обработчика
	endpoint.HandleUpdateFile(w, req)

	// Проверка ответа
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Проверка содержимого файла
	content, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("Error reading updated file: %v", err)
	}
	if string(content) != updatedContent {
		t.Errorf("Expected file content %q, got %q", updatedContent, string(content))
	}
}

// Тест удаления файла
func TestHandleDeleteFile(t *testing.T) {
	endpoint := &Endpoints{}

	// Подготовка тестового файла
	fileName := "testdelete.txt"
	err := os.WriteFile(fileName, []byte("Temporary Content"), 0644)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	defer os.Remove(fileName)

	// Подготовка запроса
	requestBody := CreateFileRequest{
		FileName: fileName,
	}
	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Error marshaling request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodDelete, "/delete", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// Выполнение обработчика
	endpoint.HandleDeleteFile(w, req)

	// Проверка ответа
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Проверка удаления файла
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		t.Error("Expected file to be deleted")
	}
}
