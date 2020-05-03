package cmd

import (
	"strings"
	"testing"
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
