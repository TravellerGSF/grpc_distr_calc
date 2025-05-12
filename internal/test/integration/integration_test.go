package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

const (
	baseURL    = "http://localhost:8080"
	testUser   = "integration_test_user"
	testPass   = "integration_test_pass"
	testExpr   = "2 + 3 * 4"
	cookieName = "auth_token"
)

func TestIntegrationWorkflow(t *testing.T) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	signupPayload := map[string]string{"username": testUser, "password": testPass}
	_ = postJSON(client, baseURL+"/auth/signup/", signupPayload, nil)

	var loginResp *http.Response
	loginPayload := map[string]string{"username": testUser, "password": testPass}
	err := postJSON(client, baseURL+"/auth/login/", loginPayload, &loginResp)
	if err != nil {
		t.Fatalf("Login error: %v", err)
	}
	defer loginResp.Body.Close()

	var token string
	for _, cookie := range loginResp.Cookies() {
		if cookie.Name == cookieName {
			token = cookie.Value
			break
		}
	}
	if token == "" {
		t.Fatalf("auth_token not found in cookies")
	}

	exprPayload := map[string]string{"expression": testExpr}
	reqBody, _ := json.Marshal(exprPayload)
	req, _ := http.NewRequest("POST", baseURL+"/expression/", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: cookieName, Value: token})

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to post expression: %v, status: %v", err, resp.Status)
	}

	var result map[string]interface{}
	var found bool
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)

		reqGet, _ := http.NewRequest("GET", baseURL+"/expression/", nil)
		reqGet.AddCookie(&http.Cookie{Name: cookieName, Value: token})
		resp, err := client.Do(reqGet)
		if err != nil {
			t.Fatalf("Error retrieving expressions: %v", err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var results []map[string]interface{}
		if err := json.Unmarshal(body, &results); err != nil {
			t.Fatalf("JSON parse error: %v", err)
		}

		if len(results) == 0 {
			continue
		}

		last := results[len(results)-1]
		if last["expression"] == testExpr {
			result = last
			if result["status"] == "done" {
				found = true
				break
			}
		}
	}

	if !found {
		t.Fatal("Expression was not processed within timeout")
	}

	if result["answer"] != "14" {
		t.Errorf("Expected answer '14', got: %v", result["answer"])
	}

	if result["status"] != "done" {
		t.Errorf("Expected status 'done', got: %v", result["status"])
	}
}

func postJSON(client *http.Client, url string, data any, respOut **http.Response) error {
	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if respOut != nil {
		*respOut = resp
	} else {
		resp.Body.Close()
	}
	return nil
}
