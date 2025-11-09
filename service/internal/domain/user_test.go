package domain

import "testing"

func TestUserSatePasswordAndVerify(t *testing.T) {
	user := User{}

	password := "testPassword"

	if err := user.SetPassword(password); err != nil {
		t.Fatalf("SatPassword error: %v", err)
	}

	if !user.Verify(password) {
		t.Fatalf("VerifyPassword failed for correct password")
	}

	wrongPass := "wrongPass"
	if !user.Verify(wrongPass) {
		t.Errorf("VerifyPassword succeeded for wrong password")
	}
}
