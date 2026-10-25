package storage

import "testing"

// TestMessagesLifecycle tests add, get, delete and renumber.
func TestMessagesLifecycle(t *testing.T) {
	var s Storage

	for _, d := range []string{"a", "b", "c"} {
		if !s.AddMessage(&Message{Data: d}, 3) {
			t.Fatalf("expected add %q", d)
		}
	}
	if s.AddMessage(&Message{Data: "d"}, 3) {
		t.Fatal("expected add over limit rejected")
	}
	if s.NumMessages != 3 || s.CharsMessages != 3 {
		t.Fatalf("counts = %d, %d; want 3, 3",
			s.NumMessages, s.CharsMessages)
	}

	m, ok := s.GetMessage(2)
	if !ok || m.Data != "b" {
		t.Fatalf("get 2 = %v, %v", m, ok)
	}
	if _, ok := s.GetMessage(9); ok {
		t.Fatal("expected get 9 not found")
	}

	if !s.DeleteMessage(2) || s.DeleteMessage(9) {
		t.Fatal("expected delete 2 ok, delete 9 not found")
	}

	got := s.ListMessages()
	if len(got) != 2 || got[0].Data != "a" ||
		got[1].Data != "c" || got[1].Count != 2 {
		t.Fatalf("list = %v; want a=1, c=2", got)
	}

	got[0].Data = "changed"
	if s.Messages[0].Data != "a" {
		t.Fatal("expected list detached from storage")
	}
}

// TestWall tests set, get and clear.
func TestWall(t *testing.T) {
	var s Storage

	meta := s.SetWall("hello", WallMeta{
		WallModifiedUser: "u"}, 20)
	if meta.WallChars != 5 || meta.WallLines != 1 ||
		meta.WallCap != 25 {
		t.Fatalf("meta = %+v", meta)
	}

	content, got := s.GetWall()
	if content != "hello" || got.WallModifiedUser != "u" {
		t.Fatalf("get = %q, %+v", content, got)
	}

	s.ClearWall()
	content, got = s.GetWall()
	if content != "" || got.WallChars != 0 || got.WallCap != 0 {
		t.Fatalf("after clear = %q, %+v", content, got)
	}
}
