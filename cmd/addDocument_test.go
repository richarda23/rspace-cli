package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/richarda23/rspace-client-go/rspace"
)

func TestGetContent(t *testing.T) {

	// plain text content
	args := addDocArgs{}
	args.Content = "abcdefg"
	content := getContent(args)

	if content != "abcdefg" {
		t.Fatalf("unexpected content")
	}

	// text file
	args.Content = ""
	args.ContentFile = "testData/textContent.txt"
	content = getContent(args)
	if !strings.Contains(content, "<pre>") {
		t.Fatalf("expected content should be wrapped in <pre> tag but was %s", content)
	}

	//html file
	args.ContentFile = "testData/textContent.html"
	content = getContent(args)
	if !strings.Contains(content, "<p> some html </p>") {
		t.Fatalf("expected verbatim html but was '%s'", content)
	}
}

func TestAddDocumentOK(t *testing.T) {
	ctx := Context{}
	outWriter := bytes.NewBufferString("")
	ctx.Writer = outWriter
	// output only gets ID field, so we don't have to make a fully populated DocumentInfo
	ctx.Format = outputFmt("quiet")
	// supply all required fields
	doAddDocRun(addDocArgs{Content: "hello", NameArg: "name", Tags: "tags"}, &ctx, &DocumentAddOK{})
	assertOKDoc(outWriter, t)

	// content is optional
	doAddDocRun(addDocArgs{NameArg: "name", Tags: "tags"}, &ctx, &DocumentAddOK{})
	assertOKDoc(outWriter, t)

	// tags is optional
	doAddDocRun(addDocArgs{NameArg: "name"}, &ctx, &DocumentAddOK{})
	assertOKDoc(outWriter, t)

	// name is optional
	doAddDocRun(addDocArgs{}, &ctx, &DocumentAddOK{})
	assertOKDoc(outWriter, t)
}

func assertOKDoc(outWriter io.Reader, t *testing.T) {
	outString, _ := ioutil.ReadAll(outWriter)
	// quiet ouput only outputs the ID
	if string(outString) != "1234\n" {
		t.Fatalf("Expected quiet output of id 1234 but was '%s'", string(outString))
	}
}

// stub for returning a created document
type DocumentAddOK struct{}

// simulates a correct response
func (ds *DocumentAddOK) NewBasicDocumentWithContent(name, tags, content string) (*rspace.DocumentInfo, error) {
	namable := &rspace.IdentifiableNamable{Id: 1234}
	rc := rspace.DocumentInfo{IdentifiableNamable: namable}
	return &rc, nil
}
