package password

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == password {
		t.Error("Hash should not be the same as original password")
	}

	if len(hash) == 0 {
		t.Error("Hash should not be empty")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "testpassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// 测试正确密码
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for correct password")
	}

	// 测试错误密码
	if CheckPassword("wrongpassword", hash) {
		t.Error("CheckPassword should return false for wrong password")
	}
}

func TestHashPasswordWithCost(t *testing.T) {
	password := "testpassword123"

	// 测试不同的 cost 值
	costs := []int{bcrypt.MinCost, bcrypt.DefaultCost, bcrypt.MaxCost}

	for _, cost := range costs {
		hash, err := HashPasswordWithCost(password, cost)
		if err != nil {
			t.Fatalf("HashPasswordWithCost failed with cost %d: %v", cost, err)
		}

		if !CheckPassword(password, hash) {
			t.Errorf("CheckPassword failed for cost %d", cost)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid password", "testpassword123", false},
		{"short password", "123", true},
		{"empty password", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPasswordError(t *testing.T) {
	err := ErrPasswordTooShort
	if err.Error() != "密码长度至少6位" {
		t.Errorf("PasswordError.Error() = %v, want %v", err.Error(), "密码长度至少6位")
	}
}
