package audit

import "testing"

func TestStoreAddListClear(t *testing.T) {
	s := NewStore(3)
	s.Add(Entry{Method: "POST", Path: "/a"})
	s.Add(Entry{Method: "DELETE", Path: "/b"})
	s.Add(Entry{Method: "PUT", Path: "/c"})
	s.Add(Entry{Method: "POST", Path: "/d"})
	items, total := s.List(10, 0, "")
	if total != 3 || len(items) != 3 {
		t.Fatalf("expected 3 entries, got total=%d len=%d", total, len(items))
	}
	if items[0].Path != "/d" {
		t.Fatalf("newest first expected /d, got %s", items[0].Path)
	}
	posts, ptotal := s.List(10, 0, "POST")
	if ptotal != 1 || len(posts) != 1 || posts[0].Path != "/d" {
		t.Fatalf("method filter failed: total=%d items=%+v", ptotal, posts)
	}
	s.Clear()
	if _, tt := s.List(10, 0, ""); tt != 0 {
		t.Fatalf("clear failed, total=%d", tt)
	}
}
