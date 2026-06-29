package util

import (
	"testing"
)

func TestUintPtr(t *testing.T) {
	v := uint(42)
	ptr := UintPtr(v)
	if ptr == nil {
		t.Fatal("UintPtr 返回 nil")
	}
	if *ptr != 42 {
		t.Errorf("*ptr = %d, 期望 42", *ptr)
	}
}

func TestGenerateAndParseToken(t *testing.T) {
	token, err := GenerateToken(1, "admin")
	if err != nil {
		t.Fatalf("GenerateToken 失败: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateToken 返回空令牌")
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken 失败: %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("UserID = %d, 期望 1", claims.UserID)
	}
	if claims.Username != "admin" {
		t.Errorf("Username = %s, 期望 'admin'", claims.Username)
	}
}

func TestGenerateTokenForDifferentUsers(t *testing.T) {
	token1, _ := GenerateToken(1, "user1")
	token2, _ := GenerateToken(2, "user2")

	claims1, _ := ParseToken(token1)
	claims2, _ := ParseToken(token2)

	if claims1.UserID == claims2.UserID {
		t.Error("不同用户的令牌应包含不同 UserID")
	}
	if claims1.Username == claims2.Username {
		t.Error("不同用户的令牌应包含不同 Username")
	}
}

func TestParseInvalidToken(t *testing.T) {
	_, err := ParseToken("invalid-token-string")
	if err == nil {
		t.Error("解析无效令牌应返回错误")
	}

	_, err = ParseToken("")
	if err == nil {
		t.Error("解析空令牌应返回错误")
	}
}

func TestParseTamperedToken(t *testing.T) {
	token, _ := GenerateToken(1, "admin")
	// 篡改令牌
	tampered := token + "tampered"
	_, err := ParseToken(tampered)
	if err == nil {
		t.Error("篡改的令牌应返回错误")
	}
}
