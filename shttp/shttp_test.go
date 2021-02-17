package shttp

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIClient_Do(t *testing.T) {
	apiClient := new(APIClient)
	resp, err := apiClient.Do(http.DefaultClient, "GET", "https://www.google.com", nil, nil)
	assert.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.NoError(t, err)
	t.Log(string(body))
}

func TestAPIClient_Do1(t *testing.T) {
	apiClient := new(APIClient)
	buffer, err := apiClient.DoBuffer(http.DefaultClient, "GET", "https://www.google.com", nil, nil)
	assert.NoError(t, err)
	t.Log(buffer.String())
	apiClient.RecycleBuffer(buffer)
}
