package errors

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		err  error
		want error
	}{
		{nil, nil},
		{fmt.Errorf("foo"), New("foo")},
		{New("a"), New("a")},
	}

	for _, tt := range tests {
		got := Wrap(tt.err)
		if got != tt.want && got.Error() != tt.want.Error() {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}

func TestNew(t *testing.T) {
	err := Newf("%s %d", "param", 1).With("param", 1)
	if err.Error() != "param 1" {
		t.Fatal(err)
	}

	frames := err.StackTrace()
	frm, _ := frames.Next()
	_, foo := filepath.Split(frm.Function)
	if foo != "errors.TestNew" {
		t.Fatalf(foo)
	}

	kv := err.KeyValues()
	if kv[0].(string) != "param" {
		t.Fatalf("kv: %v", kv)
	}

	if kv[1].(int) != 1 {
		t.Fatal()
	}
}

func TestDefer(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			err := New("err in defer")
			frames := err.StackTrace()
			frm, _ := frames.Next()
			_, foo := filepath.Split(frm.Function)
			if foo != "errors.TestDefer.func1" {
				t.Fatal()
			}
		}
	}()

	panic("paniker")
}
