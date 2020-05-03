package cmd

import (
	"testing"
	"bytes"
	"io/ioutil"
	"errors"
	"strings"
	//"github.com/spf13/cobra"
	"github.com/richarda23/rspace-client-go/rspace"

)

type ErrorStatus struct {}
func (e ErrorStatus) Status ()(*rspace.Status, error) {
	return nil, errors.New("error from status")
}

func TestStatusHandlesErr (t *testing.T) {
	cmd := NewStatusCmd()
	b := bytes.NewBufferString("")
	cmd.SetErr(b)
	//ctx := initialiseContext()
	err := doRun(cmd, []string{}, ErrorStatus{})
	
	if err == nil {
		t.Fatalf("error NOT handled")
	}
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "error") {
		t.Fatalf("expected \"%s\" got \"%s\"", "a string containing 'error'", string(out))
	}	
}