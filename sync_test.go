package partysync

import "testing"

func TestFreshSyncIsEmpty(t *testing.T) {
	s := New()
	if len(s.Version()) != 0 {
		t.Fatal("空台账的版本向量应该是空的")
	}
	if len(s.State("nobody")) != 0 {
		t.Fatal("未知成员的状态应该是空的")
	}
	if s.Pending("nobody") != 0 {
		t.Fatal("未知成员的暂存数应该是 0")
	}
}

func TestDiffWithEqualVersionIsEmpty(t *testing.T) {
	s := New()
	s.Apply("a", 1, map[string]string{"v": "m1"})
	if got := s.Diff(Vector{"a": 1}); len(got) != 0 {
		t.Fatalf("版本相同时不该有差异，得到 %v", got)
	}
}

func TestMissingWithEqualVersionIsEmpty(t *testing.T) {
	s := New()
	s.Apply("a", 1, map[string]string{"v": "m1"})
	if got := s.Missing(Vector{"a": 1}); len(got) != 0 {
		t.Fatalf("版本相同时不该缺东西，得到 %v", got)
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New().Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}
