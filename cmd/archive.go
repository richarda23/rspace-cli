package cmd

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Utility for inspecting RSpace XML archives",
	Long: `Inspect and summarise XML archives and their manifest without opening or importing into RSpace.

--summary argument calculates for one or more archives:

- the total number of documents
- the document creators
- the date range of documents created in the archive

Results are printed 1 row per archive file.

--xsummary works with a *single* archive only and lists information about each document in the archive, including name, tags, modification/creation dates and owner.

 Results are printed 1 row per *document*.

`,
	Args: cobra.MinimumNArgs(1),
	Example: `
// show manifest
rspace archive myArchive.zip --manifest

//summarise the content in one or more archive files
rspace archive myArchive.zip otherArchive.zip --summary

//extended summary of a single archive in csv format
rspace archive myArchive.zip --xsummary --outputFormat csv
	`,

	Run: func(cmd *cobra.Command, args []string) {
		ctx := initialiseContext()
		if archiveArgsA.summaryArg {
			summaries, err := inspectArchives(args, &archiveArgsA)
			if err != nil {
				exitWithErr(err)
			}
			ctx.writeResult(&zipSummaryFormatter{&zipSummaryList{summaries}})
		} else if archiveArgsA.summaryXArg {
			summaries, err := xSummary(args, &archiveArgsA)
			if err != nil {
				exitWithErr(err)
			}
			ctx.writeResult(&xSummaryFormatter{&xSummaryList{summaries}})
		}
	},
}

type xSummaryList struct {
	XSummaryList []*xmlDoc
}

type xSummaryFormatter struct {
	xSummaries *xSummaryList
}

func (xs *xSummaryFormatter) ToJson() string {
	return prettyMarshal(xs.xSummaries)
}

func (ds *xSummaryFormatter) ToQuiet() []identifiable {
	rows := make([]identifiable, 0)
	for _, res := range ds.xSummaries.XSummaryList {
		rows = append(rows, identifiable{res.Name})
	}
	return rows
}

func (ds *xSummaryFormatter) ToTable() *TableResult {
	headers := []columnDef{columnDef{"Name", 25}, columnDef{"Tags", 15},
		columnDef{"created", 22}, columnDef{"lastModified", 22}, columnDef{"Owner", 12}}

	rows := make([][]string, 0)
	for _, res := range ds.xSummaries.XSummaryList {
		data := []string{res.Name, res.Tags,
			res.CreationDate.Format(time.RFC3339), res.LastModifiedDate.Format(time.RFC3339),
			res.CreatedBy}
		rows = append(rows, data)
	}
	return &TableResult{headers, rows}
}

type zipSummaryList struct {
	SummaryList []*zipSummary
}

type zipSummaryFormatter struct {
	results *zipSummaryList
}

func (zs *zipSummaryFormatter) ToJson() string {
	return prettyMarshal(zs.results)
}

func (ds *zipSummaryFormatter) ToQuiet() []identifiable {
	rows := make([]identifiable, 0)
	for _, res := range ds.results.SummaryList {
		rows = append(rows, identifiable{res.FileName})
	}
	return rows
}

func (ds *zipSummaryFormatter) ToTable() *TableResult {
	headers := []columnDef{columnDef{"file", 10}, columnDef{"Total Docs", 10},
		columnDef{"minDate", 22}, columnDef{"maxDate", 22}, columnDef{"Authors", 50}}

	rows := make([][]string, 0)
	for _, res := range ds.results.SummaryList {
		data := []string{res.FileName, strconv.Itoa(res.DocCount),
			res.MinDate.Format(time.RFC3339), res.MaxDate.Format(time.RFC3339), strings.Join(res.Authors, ";")}
		rows = append(rows, data)
	}
	return &TableResult{headers, rows}
}

type archiveArgs struct {
	summaryArg  bool
	summaryXArg bool
	manifestArg bool
}

var archiveArgsA archiveArgs

func xSummary(args []string, config *archiveArgs) ([]*xmlDoc, error) {
	if len(args) > 1 {
		return nil, errors.New("Extended summary only available on a single archive file")
	}
	file := args[0]
	if filepath.Ext(file) != ".zip" {
		return nil, errors.New(fmt.Sprintf("%s is not a zip file", file))
	}
	reader, err := zip.OpenReader(file)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Couldn't open the zip file %s", file))
	}
	parsedDocs := parseArchiveFiles(reader)
	return parsedDocs, nil
}

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
			messageStdErr(fmt.Sprintf("Manifest for %s:", file))
			bytes, err := showManifest(reader)
			if err != nil {
				messageStdErr(err.Error())
			} else {
				messageStdErr(string(bytes))
			}
		}
		parsedDocs := parseArchiveFiles(reader)
		if archiveArgsA.summaryArg {
			//messageStdErr(fmt.Sprintf("Summary for %s:", file))
			summary, _ := summarise(parsedDocs)
			summary.FileName = filepath.Base(file)
			zipSummaries = append(zipSummaries, summary)
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
	Tags             string    `xml:"tag"`
}

func parseTimestamp(timestamp string) (time.Time, error) {
	return time.Parse(time.RFC3339, timestamp)
}

func parseArchiveFiles(reader *zip.ReadCloser) []*xmlDoc {
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
	return parsedDocs
}

type zipSummary struct {
	DocCount int
	MinDate  time.Time
	MaxDate  time.Time
	Authors  []string
	FileName string
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
	return &zipSummary{len(docs), minDate, maxDate, uniqueAuthors, ""}, nil
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
	archiveCmd.Flags().BoolVar(&archiveArgsA.summaryXArg, "xsummary", false, "Show Extended summary of content")
	archiveCmd.Flags().BoolVar(&archiveArgsA.manifestArg, "manifest", true, "Shows manifest of the archive")
	archiveCmd.Flags().StringVar(&outputFormatArg, "outputFormat", "table", "Output format: one of 'json','table', 'csv' or 'quiet' ")
	archiveCmd.Flags().StringVar(&outFileArg, "outFile", "", "Output file for program output")
}
