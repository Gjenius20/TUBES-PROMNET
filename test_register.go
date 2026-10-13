package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	body := map[string]string{
		"name":     "Student2",
		"email":    "student2_new@test.com",
		"nim":      "STU003",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post("http://localhost:1111/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.Status)
	fmt.Println("Body:", string(respBody))
}