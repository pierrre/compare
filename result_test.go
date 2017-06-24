package compare

import (
	"fmt"
	"testing"
)

func TestResultFormat(t *testing.T) {
	r := Result{
		Difference{
			Message: "test1",
			V1:      1,
			V2:      2,
		},
		Difference{
			Message: "test2",
			V1:      "a",
			V2:      "b",
		},
	}
	s := fmt.Sprintf("%+v", r)
	expected := ".: test1\n\tv1=1\n\tv2=2\n.: test2\n\tv1=\"a\"\n\tv2=\"b\""
	if s != expected {
		t.Fatalf("unexpected result:\ngot:\n%s\nwant:\n%s", s, expected)
	}
}

func TestResultFormatEmpty(t *testing.T) {
	var r Result
	s := fmt.Sprintf("%+v", r)
	expected := "<none>"
	if s != expected {
		t.Fatalf("unexpected result:\ngot:\n%s\nwant:\n%s", s, expected)
	}
}

func TestResultFormatUnsupportedVerb(t *testing.T) {
	var r Result
	s := fmt.Sprintf("%s", r)
	expected := "%!s(compare.Result)"
	if s != expected {
		t.Fatalf("unexpected result:\ngot:\n%s\nwant:\n%s", s, expected)
	}
}

func TestDifferenceFormatUnsupportedVerb(t *testing.T) {
	var d Difference
	s := fmt.Sprintf("%s", d)
	expected := "%!s(compare.Difference)"
	if s != expected {
		t.Fatalf("unexpected result:\ngot:\n%s\nwant:\n%s", s, expected)
	}
}
