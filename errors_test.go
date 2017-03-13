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
	err := New("", "param", 1)
	if err.Error() != "" {
		t.Fatal(err)
	}

	frames := err.GetStack()
	frm, _ := frames.Next()
	_, foo := filepath.Split(frm.Function)
	if foo != "errors.TestNew" {
		t.Fatalf(foo)
	}

	kv := err.GetKeyval()
	if kv[0].(string) != "param" {
		t.Fatal()
	}

	if kv[1].(int) != 1 {
		t.Fatal()
	}

}

func TestForamtMsg(t *testing.T) {
	tests := []struct {
		msg    string
		keyval []interface{}
		want   string
	}{
		{"", nil, ""},
		{"", []interface{}{1, 2, 3, 4}, ""},
		{"%d %d", []interface{}{1, 2, 3, 4}, "2 4"},
		{"%d %% %d", []interface{}{1, 2, 3, 4}, "2 % 4"},
		{"%v %v", []interface{}{1, 2, 3, 4}, "2 4"},
		{"%v %v", []interface{}{1, 2}, "2 %!v(MISSING)"},
		{"%v %v", nil, "%!v(MISSING) %!v(MISSING)"},
	}

	for _, tt := range tests {
		got := foramtMsg(tt.msg, tt.keyval)
		if got != tt.want {
			t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
		}
	}
}
