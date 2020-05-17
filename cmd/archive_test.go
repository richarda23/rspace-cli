package cmd

import (
	"testing"
)

func TestReadArchives(t *testing.T) {
	archiveArgsA.summaryArg = true
	summaries, _ := inspectArchives([]string{"testData/rs3.zip", "testData/rs2.zip"},
		&archiveArgsA)
	if len(summaries) != 2 {
		t.Fatalf("Expected %d results but got %d", 2, len(summaries))
	}
}

func TestReadArchive(t *testing.T) {
	archiveArgsA.summaryArg = true
	summaries, _ := inspectArchives([]string{"testData/rs3.zip"},
		&archiveArgsA)
	if len(summaries) != 1 {
		t.Fatalf("Expected %d results but got %d", 1, len(summaries))
	}
	summary := summaries[0]
	if summary.docCount != 3 {
		t.Fatalf("Expected %d docs in archive but got %d", 3, summary.docCount)
	}
	if summary.minDate.IsZero() {
		t.Fatalf("Min date is Zero but should be set")
	}
	if summary.maxDate.IsZero() {
		t.Fatalf("Max date is Zero but should be set")
	}

	if summary.minDate.After(summary.maxDate) {
		t.Fatalf("Min date %s must be before %s", summary.minDate.String(), summary.maxDate.String())
	}

	authors := summary.authors
	if authors[0] != "user5e" {
		t.Fatalf("Authors should be user5e")
	}
}
