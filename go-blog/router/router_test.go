package router

import (
	"testing"
)

func TestAdd(t *testing.T) {
	if got := add(3, 4); got != 7 {
		t.Errorf("add(3,4) = %d, 期望 7", got)
	}
	if got := add(-1, 1); got != 0 {
		t.Errorf("add(-1,1) = %d, 期望 0", got)
	}
}

func TestSub(t *testing.T) {
	if got := sub(10, 3); got != 7 {
		t.Errorf("sub(10,3) = %d, 期望 7", got)
	}
	if got := sub(3, 10); got != -7 {
		t.Errorf("sub(3,10) = %d, 期望 -7", got)
	}
}

func TestMul(t *testing.T) {
	if got := mul(6, 7); got != 42 {
		t.Errorf("mul(6,7) = %d, 期望 42", got)
	}
	if got := mul(0, 100); got != 0 {
		t.Errorf("mul(0,100) = %d, 期望 0", got)
	}
}

func TestDiv(t *testing.T) {
	if got := div(10, 3); got != 3 {
		t.Errorf("div(10,3) = %d, 期望 3", got)
	}
	if got := div(10, 0); got != 0 {
		t.Errorf("div(10,0) = %d, 期望 0（零除返回0）", got)
	}
}

func TestIterate(t *testing.T) {
	result := iterate(5)
	if len(result) != 5 {
		t.Fatalf("iterate(5) 长度=%d, 期望 5", len(result))
	}
	for i, v := range result {
		if v != i {
			t.Errorf("iterate(5)[%d] = %d, 期望 %d", i, v, i)
		}
	}

	if len(iterate(0)) != 0 {
		t.Error("iterate(0) 应返回空切片")
	}
}
