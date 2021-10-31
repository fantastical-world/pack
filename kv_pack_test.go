package pack

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

type mockThing struct {
	Name   string   `json:"name"`
	Points int      `json:"points"`
	Meta   mockMeta `json:"meta"`
}

type mockMeta struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func (m mockThing) Pack() (string, []byte) {
	b, _ := json.Marshal(m)
	return m.Name, b
}
func (m *mockThing) Unpack(b []byte) {
	json.Unmarshal(b, m)
}

type mockThingNoMeta struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
}

func (m mockThingNoMeta) Pack() (string, []byte) {
	b, _ := json.Marshal(m)
	return m.Name, b
}
func (m *mockThingNoMeta) Unpack(b []byte) {
	json.Unmarshal(b, m)
}

type notStruct string

func (n notStruct) Pack() (string, []byte) {
	return string(n), []byte(n)
}
func (n *notStruct) Unpack(b []byte) {
}

//this does nothing more than validate Pack interface compliance
var _ Pack = (*KVPack)(nil)

func Test_New(t *testing.T) {
	_ = New("./test.db")
}

func TestKVPack_Save(t *testing.T) {
	db := New("./test.db")
	mock := &mockThing{Name: "heyo", Points: 100}
	name, want := mock.Pack()
	t.Run("validate that save correctly persists", func(t *testing.T) {
		err := db.Save("__TEST__", mock)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}

		got, err := db.Get("__TEST__", name)
		if err != nil {
			t.Errorf("unexpected error getting, %s", err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}

		mock2 := &mockThing{}
		mock2.Unpack(got)
		if !reflect.DeepEqual(mock, mock2) {
			t.Errorf("want %v, got %v", mock, mock2)
		}
	})

	t.Run("validate that save returns an error if invalid", func(t *testing.T) {
		emptyThing := &mockThing{}
		err := db.Save("__TEST__", emptyThing)
		if err == nil {
			t.Error("expected an error, but none encountered")
		}
	})

	err := os.Remove("./test.db")
	if err != nil {
		t.Errorf("unexpected err encountered, %s", err)
	}
}

func TestKVPack_Get(t *testing.T) {
	db := New("./test.db")
	mock := &mockThing{Name: "heyo", Points: 100}
	name, want := mock.Pack()
	t.Run("validate that get correctly returns thing", func(t *testing.T) {
		err := db.Save("__TEST__", mock)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}

		got, err := db.Get("__TEST__", name)
		if err != nil {
			t.Errorf("unexpected error getting, %s", err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	err := os.Remove("./test.db")
	if err != nil {
		t.Errorf("unexpected err encountered, %s", err)
	}
}

func TestKVPack_Delete(t *testing.T) {
	db := New("./test.db")
	mock := &mockThing{Name: "heyo", Points: 100}
	name, want := mock.Pack()
	t.Run("validate that delete correctly deletes", func(t *testing.T) {
		err := db.Save("__TEST__", mock)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}

		got, err := db.Get("__TEST__", name)
		if err != nil {
			t.Errorf("unexpected error getting, %s", err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}

		err = db.Delete("__TEST__", name)
		if err != nil {
			t.Errorf("unexpected error getting, %s", err)
		}

		got, err = db.Get("__TEST__", name)
		if err != ErrThingDoesNotExist {
			t.Errorf("want %s, got %s", ErrThingDoesNotExist, err)
		}

		if got != nil {
			t.Errorf("expected nil, got %s", got)
		}
	})

	t.Run("validate that delete returns an error if thing does not exist", func(t *testing.T) {
		err := db.Delete("IDONTEXIST", "IDONT")
		if err != ErrThingDoesNotExist {
			t.Errorf("want %s, got %s", ErrThingDoesNotExist, err)
		}
	})

	err := os.Remove("./test.db")
	if err != nil {
		t.Errorf("unexpected err encountered, %s", err)
	}
}

func TestKVPack_List(t *testing.T) {
	db := New("./test.db")

	t.Run("validate listing", func(t *testing.T) {
		mock := &mockThing{Name: "heyo", Points: 100}
		mock2 := &mockThing{Name: "heyo2", Points: 200}
		mock3 := &mockThing{Name: "heyo3", Points: 300}

		err := db.Save("__TEST__", mock)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}
		err = db.Save("__TEST__", mock2)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}
		err = db.Save("__TEST__", mock3)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}

		got, err := db.List("__TEST__")
		if err != nil {
			t.Errorf("unexpected error listing things, %s", err)
		}

		want := []string{"heyo", "heyo2", "heyo3"}
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})

	err := os.Remove("./test.db")
	if err != nil {
		t.Errorf("unexpected err encountered, %s", err)
	}
}

func TestKVPack_ListMeta(t *testing.T) {
	db := New("./test.db")

	t.Run("validate meta listing (single)", func(t *testing.T) {
		mock := &mockThing{Name: "heyo", Points: 100, Meta: mockMeta{Name: "meta-name1", Count: 1}}
		want := mockMeta{Name: "meta-name1", Count: 1}
		err := db.Save("__TEST__", mock)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}

		got, err := db.ListMeta("__TEST__")
		if err != nil {
			t.Errorf("unexpected error listing meta, %s", err)
		}

		for _, meta := range got {
			temp, e := json.Marshal(meta)
			if e != nil {
				t.Errorf("unexpected error marshaling meta, %s", e)
			}
			var gotMeta mockMeta
			e = json.Unmarshal(temp, &gotMeta)
			if e != nil {
				t.Errorf("unexpected error unmarshaling meta, %s", e)
			}

			if !reflect.DeepEqual(want, gotMeta) {
				t.Errorf("want %v, got %v", want, gotMeta)
			}
		}
	})

	t.Run("validate meta listing", func(t *testing.T) {
		mock := &mockThing{Name: "heyo", Points: 100, Meta: mockMeta{Name: "meta-name1", Count: 1}}
		mock2 := &mockThing{Name: "heyo2", Points: 200, Meta: mockMeta{Name: "meta-name2", Count: 2}}
		mock3 := &mockThing{Name: "heyo3", Points: 300, Meta: mockMeta{Name: "meta-name3", Count: 3}}

		err := db.Save("__TEST__", mock)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}
		err = db.Save("__TEST__", mock2)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}
		err = db.Save("__TEST__", mock3)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}

		got, err := db.ListMeta("__TEST__")
		if err != nil {
			t.Errorf("unexpected error listing meta, %s", err)
		}

		if len(got) != 3 {
			t.Errorf("want 3, got %d", len(got))
		}

		for _, meta := range got {
			temp, e := json.Marshal(meta)
			if e != nil {
				t.Errorf("unexpected error marshaling meta, %s", e)
			}
			var gotMeta mockMeta
			e = json.Unmarshal(temp, &gotMeta)
			if e != nil {
				t.Errorf("unexpected error unmarshaling meta, %s", e)
			}

			if !strings.HasPrefix(gotMeta.Name, "meta-name") {
				t.Errorf("want meta-name#, got %s", gotMeta.Name)
			}
			if gotMeta.Count == 0 {
				t.Errorf("want #, got %d", gotMeta.Count)
			}
		}
	})

	t.Run("validate empty meta", func(t *testing.T) {
		mock := &mockThing{Name: "heyo", Points: 100}
		mock2 := &mockThing{Name: "heyo2", Points: 200}
		mock3 := &mockThing{Name: "heyo3", Points: 300}

		err := db.Save("__TEST__", mock)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}
		err = db.Save("__TEST__", mock2)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}
		err = db.Save("__TEST__", mock3)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}

		got, err := db.ListMeta("__TEST__")
		if err != nil {
			t.Errorf("unexpected error listing meta, %s", err)
		}

		if len(got) != 3 {
			t.Errorf("want 3, got %d", len(got))
		}

		want := mockMeta{}
		for _, meta := range got {
			temp, e := json.Marshal(meta)
			if e != nil {
				t.Errorf("unexpected error marshaling meta, %s", e)
			}
			var gotMeta mockMeta
			e = json.Unmarshal(temp, &gotMeta)
			if e != nil {
				t.Errorf("unexpected error unmarshaling meta, %s", e)
			}

			if !reflect.DeepEqual(want, gotMeta) {
				t.Errorf("want %v, got %v", want, gotMeta)
			}
		}
	})

	t.Run("validate no meta", func(t *testing.T) {
		mock := &mockThingNoMeta{Name: "heyo", Points: 100}
		mock2 := &mockThingNoMeta{Name: "heyo2", Points: 200}
		mock3 := &mockThingNoMeta{Name: "heyo3", Points: 300}

		err := db.Save("__TEST__", mock)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}
		err = db.Save("__TEST__", mock2)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}
		err = db.Save("__TEST__", mock3)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}

		got, err := db.ListMeta("__TEST__")
		if err != nil {
			t.Errorf("unexpected error listing meta, %s", err)
		}

		if len(got) != 0 {
			t.Errorf("want 0, got %d", len(got))
		}
	})

	t.Run("validate error not struct", func(t *testing.T) {
		mock := notStruct("heyo")

		err := db.Save("__TEST__", &mock)
		if err != nil {
			t.Errorf("unexpected error saving, %s", err)
		}

		got, err := db.ListMeta("__TEST__")
		if err == nil {
			t.Error("expected error, got nil")
		}

		if len(got) != 0 {
			t.Errorf("want 0, got %d", len(got))
		}
	})

	err := os.Remove("./test.db")
	if err != nil {
		t.Errorf("unexpected err encountered, %s", err)
	}
}

func Test_BadDatabaseErrors(t *testing.T) {
	db := KVPack{}
	mock := &mockThing{}

	err := db.Save("__TEST__", mock)
	if err == nil {
		t.Error("expected an error, but none encountered\n")
	}
	_, err = db.Get("__TEST__", "test")
	if err == nil {
		t.Error("expected an error, but none encountered\n")
	}
	err = db.Delete("__TEST__", "test")
	if err == nil {
		t.Error("expected an error, but none encountered\n")
	}
	_, err = db.List("__TEST__")
	if err == nil {
		t.Error("expected an error, but none encountered\n")
	}
	_, err = db.ListMeta("__TEST__")
	if err == nil {
		t.Error("expected an error, but none encountered\n")
	}
}
