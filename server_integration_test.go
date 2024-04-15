//go:build integration

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	url := "http://localhost:8080/api/users"
	contentType := "application/json"
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", "Basic am9lOnNlY3JldA==")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	var target []User

	err = json.NewDecoder(resp.Body).Decode(&target)

	if err != nil {
		fmt.Println(err)

	}
	// fmt.Println(string(body))
	assert.EqualValues(t, http.StatusOK, resp.StatusCode)
	assert.Greater(t, len(target), 0)

}
