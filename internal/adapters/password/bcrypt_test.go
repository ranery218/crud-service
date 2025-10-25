package password

import (
	"context"
	"testing"
)

func TestBcryptHasher(t *testing.T) {
	hasher := NewBcryptHasher(0)
	ctx := context.Background()

	hash, err := hasher.Hash(ctx, "mysecret123")
	if err != nil {
		t.Fatalf("Hashing failed: %v", err)
	}

	if hash == "" || hash == "mysecret123" {
		t.Fatalf("Invalid hash generated")
	}

	if err = hasher.Compare(ctx, hash, "mysecret123"); err != nil {
		t.Fatalf("compare err: %v", err)
	}
}

func TestBcryptHasher_WrongPassword(t *testing.T) {
	hasher := NewBcryptHasher(0)
	ctx := context.Background()

	hash, _ := hasher.Hash(ctx, "mysecret123")
	if err := hasher.Compare(ctx, hash, "bad"); err == nil {
		t.Fatal("expected error for wrong password")
	}
}
