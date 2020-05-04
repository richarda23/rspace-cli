package cmd

import (
	"testing"
)

func TestIdFromGlobalId(t *testing.T) {
	id, err := idFromGlobalId("not an id")
	if err != nil {
		t.Fatalf("Error should b non-nil")
	}
	id, err = idFromGlobalId("FL12334")
	assertIdEquals(t, id, 12334)

	id, err = idFromGlobalId("12334")
	assertIdEquals(t, id, 12334)

}

func assertIdEquals(t *testing.T, got, expected int) {
	if got != expected {
		t.Fatalf("Expected %d but was %d", expected, got)
	}
}
