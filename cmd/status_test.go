package cmd

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"
	"testing"

	//"github.com/spf13/cobra"
	"github.com/richarda23/rspace-client-go/rspace"
)

type ErrorStatus struct{}

func (e ErrorStatus) Status() (*rspace.Status, error) {
	return nil, errors.New("error from status")
}

type OKStatus struct{}

func (e OKStatus) Status() (*rspace.Status, error) {
	return &rspace.Status{"message", "version"}, nil
}

func TestStatus(t *testing.T) {

	// test-spy for context error-writer
	errWriter := bytes.NewBufferString("")
	ctx := &Context{}
	ctx.ErrWriter = errWriter
	err := doRun(ErrorStatus{}, ctx)
	if err == nil {
		t.Fatal("error NOT handled")
	}
	errString, err := ioutil.ReadAll(errWriter)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(errString), "error") {
		t.Fatalf("expected '%s' got '%s'", "a string containing 'error'", string(errString))
	}

	// now assert OK status is written to Context OutWriter
	outWriter := bytes.NewBufferString("")
	ctx.Writer = outWriter
	err = doRun(OKStatus{}, ctx)
	if err != nil {
		t.Fatal(err)
	}
	outString, _ := ioutil.ReadAll(outWriter)
	if !strings.Contains(string(outString), "message") {
		t.Fatalf("expected '%s' got '%s'", "a string containing 'error'", string(outString))
	}

}
