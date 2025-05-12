package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TravellerGSF/grpc_distr_calc/internal/storage"
	"github.com/TravellerGSF/grpc_distr_calc/internal/utils/orchestrator/jwts"
)

type mockUserStorage struct{}

func (m *mockUserStorage) GetUserByID(id int64) (*storage.User, error) {
	return &storage.User{ID: id, Name: "mockUser"}, nil
}

func (m *mockUserStorage) RegisterUser(ctx context.Context, uname, pswrd string) error {
	return nil
}
func (m *mockUserStorage) LoginUser(ctx context.Context, uname, pswrd string) (string, error) {
	return "mockToken", nil
}

func TestGetUsernameHandler_ValidToken(t *testing.T) {
	mock := &mockUserStorage{}
	token, _ := jwts.GenerateJWTToken(42)

	req := httptest.NewRequest("GET", "/auth/username/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler := GetUsernameHandler(mock)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rr.Code)
	}

	if rr.Body.String() == "" {
		t.Error("expected response body with username")
	}
}

func TestGetUsernameHandler_InvalidToken(t *testing.T) {
	mock := &mockUserStorage{}

	req := httptest.NewRequest("GET", "/auth/username/", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rr := httptest.NewRecorder()

	handler := GetUsernameHandler(mock)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", rr.Code)
	}
}
