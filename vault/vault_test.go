package vault

import (
	"testing"
)

func TestUnlock(t *testing.T) {
	test := []struct {
		name        string
		input           string
		expectError bool
	}{
		{"correct password","secret123",false},
		{"wrong password","abc",true},
		{"empty password","",true},
	}
	for _, tc := range test{
		t.Run(tc.name, func(t *testing.T) {
			testVault, makeError := NewVault("$$oohitsAsecret$$")
			if makeError != nil {
				t.Fatal(makeError)
			}
			err := testVault.Unlock(tc.input)
			if tc.expectError && err == nil{
				t.Fatalf("expected an error, but got nil")
			}
			if !tc.expectError && err != nil {
				t.Fatalf("did not expect an error but got:%v", err)
			}
			if !tc.expectError && testVault.isLocked {
				t.Fatalf("expected vault to be unlocked but vault is still locked")
			}

		})
	}
}
