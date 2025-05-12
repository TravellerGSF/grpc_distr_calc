package storage

import (
	"context"
	"os"
	"testing"
)

func setupTestDB(t *testing.T) *Storage {
	path := "./test_user.db"
	_ = os.Remove(path)
	db, err := New(path)
	if err != nil {
		t.Fatalf("failed to create test DB: %v", err)
	}
	return db
}

func TestRegisterAndLoginUser(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.TODO()

	err := db.RegisterUser(ctx, "user1", "pass123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	token, err := db.LoginUser(ctx, "user1", "pass123")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	if token == "" {
		t.Error("expected JWT token, got empty string")
	}
}

func TestLoginWithWrongPassword(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.TODO()

	_ = db.RegisterUser(ctx, "user2", "pass456")
	_, err := db.LoginUser(ctx, "user2", "wrongpass")
	if err == nil {
		t.Error("expected login to fail with wrong password")
	}
}
