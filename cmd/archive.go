package cmd

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Utility for inspecting RSpace XML archives",
	Long:  `Inspect and summarise XML archives without opening or importing into RSpace`,
	Args:  cobra.MinimumNArgs(1),
	Example: `
// show manifest
rspace archive myArchive.zip --manifest

//summarise the content in the archive
rspace archive myArchive.zip --summary
	`,

	Run: func(cmd *cobra.Command, args []string) {
		summaries, err := inspectArchives(args, &archiveArgsA)
		if err != nil {
			exitWithErr(err)
		}
		fmt.Println(summaries)
	},
}

type archiveArgs struct {
	summaryArg  bool
	manifestArg bool
}

var archiveArgsA archiveArgs

func inspectArchives(args []string, config *archiveArgs) ([]*zipSummary, error) {
	zipSummaries := make([]*zipSummary, 0)
	for _, file := range args {
		if filepath.Ext(file) != ".zip" {
			messageStdErr(fmt.Sprintf("%s is not a zip file, skipping", file))
			continue
		}

		reader, err := zip.OpenReader(file)
		if err != nil {
			//messageStdErr(fmt.Sprintf("Could not open %s, skipping", file))
		}
		if archiveArgsA.manifestArg {
			//messageStdErr(fmt.Sprintf("Manifest for %s:", file))
			bytes, err := showManifest(reader)
			if err != nil {
				messageStdErr(err.Error())
			} else {
				messageStdErr(string(bytes))
			}
		}
		if archiveArgsA.summaryArg {
			//messageStdErr(fmt.Sprintf("Summary for %s:", file))
			files := parseFiles(reader)
			zipSummaries = append(zipSummaries, files)
		}
	}

	return zipSummaries, nil
}

type xmlDoc struct {
	XMLName          xml.Name
	Name             string    `xml:"name"`
	CreatedBy        string    `xml:"createdBy"`
	CreationDate     time.Time `xml:"creationDate"`
	LastModifiedDate time.Time `xml:"lastModifiedDate"`
}

func parseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}

func parseFiles(reader *zip.ReadCloser) *zipSummary {
	parsedDocs := make([]*xmlDoc, 0)
	for _, f := range reader.File {
		fname := filename(f)
		if strings.HasSuffix(fname, "xml") && strings.HasPrefix(fname, "doc") && !strings.HasSuffix(fname, "_form.xml") {
			fc, _ := f.Open()
			bytes, _ := ioutil.ReadAll(fc)
			mydoc := xmlDoc{}
			xml.Unmarshal(bytes, &mydoc)
			parsedDocs = append(parsedDocs, &mydoc)
		}
	}
	summary, _ := summarise(parsedDocs)
	return summary
}

type zipSummary struct {
	docCount int
	minDate  time.Time
	maxDate  time.Time
	authors  []string
}

func summarise(docs []*xmlDoc) (*zipSummary, error) {
	if len(docs) == 0 {
		errors.New("No documents to summarise")
	}

	maxDate := time.Time{}
	minDate := docs[0].CreationDate
	authors := make(map[string]bool, 0)
	for _, v := range docs {
		if v.CreationDate.After(maxDate) {
			maxDate = v.CreationDate
		}
		if v.CreationDate.Before(minDate) {
			minDate = v.CreationDate
		}
		authors[v.CreatedBy] = true
	}
	var uniqueAuthors = make([]string, 0)
	for k, _ := range authors {
		uniqueAuthors = append(uniqueAuthors, k)
	}
	return &zipSummary{len(docs), minDate, maxDate, uniqueAuthors}, nil
}

func filename(file *zip.File) string {
	return filepath.Base(file.Name)
}

func showManifest(reader *zip.ReadCloser) ([]byte, error) {
	for _, f := range reader.File {
		if filename(f) == "manifest.txt" {
			fc, _ := f.Open()
			bytes, _ := ioutil.ReadAll(fc)
			return bytes, nil
		}
	}
	return nil, errors.New("No manifest.txt file found")
}

func init() {
	rootCmd.AddCommand(archiveCmd)
	// is called directly, e.g.:
	archiveCmd.Flags().BoolVar(&archiveArgsA.summaryArg, "summary", false, "Show summary of content")
	archiveCmd.Flags().BoolVar(&archiveArgsA.manifestArg, "manifest", true, "Shows manifest of the archive")
}
