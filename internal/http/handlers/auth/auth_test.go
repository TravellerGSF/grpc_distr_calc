package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockUserInteractor struct{}

func (m *mockUserInteractor) RegisterUser(ctx context.Context, uname, pswrd string) error {
	if uname == "exists" {
		return ErrAlreadyExists
	}
	return nil
}

func (m *mockUserInteractor) LoginUser(ctx context.Context, uname, pswrd string) (string, error) {
	if uname == "fail" {
		return "", ErrInvalidCredentials
	}
	return "token123", nil
}

var (
	ErrAlreadyExists      = &mockError{"user exists"}
	ErrInvalidCredentials = &mockError{"invalid credentials"}
)

type mockError struct {
	msg string
}

func (e *mockError) Error() string { return e.msg }

func TestRegisterUserHandler_Success(t *testing.T) {
	body := map[string]string{"username": "newuser", "password": "1234"}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/signup/", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()

	handler := RegisterUserHandler(context.TODO(), &mockUserInteractor{})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", rr.Code)
	}
}

func TestLoginUserHandler_Success(t *testing.T) {
	body := map[string]string{"username": "admin", "password": "1234"}
	payload, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/auth/login/", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()

	handler := LoginUserHandler(context.TODO(), &mockUserInteractor{})
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("expected 303 SeeOther, got %d", rr.Code)
	}

	cookie := rr.Result().Cookies()
	if len(cookie) == 0 || cookie[0].Name != "auth_token" {
		t.Error("expected auth_token cookie")
	}
}
