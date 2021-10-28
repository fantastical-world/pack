package pack

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

type mockThing struct {
	Name   string `json:"name"`
	Points int    `json:"points"`
}

func (m mockThing) Pack() (string, []byte) {
	b, _ := json.Marshal(m)
	return m.Name, b
}
func (m *mockThing) Unpack(b []byte) {
	json.Unmarshal(b, m)
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
	t.Run("validate that save correctly persists...", func(t *testing.T) {
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

	t.Run("validate that save returns an error if invalid...", func(t *testing.T) {
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
	t.Run("validate that save correctly persists...", func(t *testing.T) {
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
	t.Run("validate that save correctly persists...", func(t *testing.T) {
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

	t.Run("validate that delete table returns an error if table does not exist...", func(t *testing.T) {
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

func TestKVPack_ListTables(t *testing.T) {
	db := New("./test.db")

	t.Run("validate table listing...", func(t *testing.T) {
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
			t.Errorf("unexpected error listing tables, %s", err)
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
	_, err = db.List("__TEST__")
	if err == nil {
		t.Error("expected an error, but none encountered\n")
	}
	err = db.Delete("__TEST__", "test")
	if err == nil {
		t.Error("expected an error, but none encountered\n")
	}
}
