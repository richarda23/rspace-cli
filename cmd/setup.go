package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/url"
	"os"
	"rspace"
	"strings"
)

const (
	APIKEY_ENV_NAME   = "RSPACE_API_KEY"
	BASE_URL_ENV_NAME = "RSPACE_URL"
)

var (
	validOutputFormats = []string{"json", "csv", "quiet", "table"}
	validTreeFilters   = []string{"document", "notebook", "folder"}
	validSortOrders    = []string{"asc", "desc"}
	validRecordOrders  = []string{"name", "created", "lastModified"}
	outputFormat       outputFmt
)

type outputFmt string

func (ft outputFmt) isJson() bool {
	return string(ft) == "json"
}
func (ft outputFmt) isCsv() bool {
	return ft == "csv"
}
func (ft outputFmt) isTab() bool {
	return ft == "table"
}
func (ft outputFmt) isQuiet() bool {
	return ft == "quiet"
}

type Context struct {
	WebClient *rspace.RsWebClient
	Writer    io.Writer
	Format    outputFmt
}

func (ctx *Context) writeResult(formatter ResultListFormatter) {
	if ctx.Format.isJson() {
		ctx.write(formatter.ToJson())
	} else if ctx.Format.isQuiet() {
		printIds(ctx, formatter.ToQuiet())
	} else if (ctx.Format.isCsv()) {
		printCsv(ctx, formatter.ToTable())
	} else {
		printTable(ctx, formatter.ToTable())
	}
}

func (ctx *Context) write(toWrite string) {
	fmt.Fprintln(ctx.Writer, toWrite)
}

func initialiseContext() *Context {
	_validateFlagArgs()
	outputFormat = outputFmt(outputFormatArg)
	rc := Context{}
	rc.WebClient = initWebClient()
	rc.Writer = initOutputWriter(outFileArg)
	rc.Format = outputFormat
	return &rc
}

// exits with error if validation fails
func _validateFlagArgs() {
	outputFormat = outputFmt(outputFormatArg)
	validateOutputFormatExit(outputFormat)
	validateTreeFilterExit(treeFilterArg)
}
func validateTreeFilterExit(treeFilterArg string) []string {
	if len(treeFilterArg) == 0 {
		return make([]string, 0)
	}
	rc := strings.Split(treeFilterArg, ",")

	if !validateArrayContains(validTreeFilters, rc) {
		exitWithStdErrMsg("Invalid tree filter, must be comma-separated list of 1 more terms: " + strings.Join(validTreeFilters, ","))
		return nil
	}
	return rc
}
func validateOutputFormatExit(toTest outputFmt) {
	if !validateOutputFormat(toTest) {
		exitWithStdErrMsg("Invalid outputFormat argument: must be one of: " + strings.Join(validOutputFormats, ","))
	}
}
func validateArrayContains(validTerms []string, toTest []string) bool {
	for _, term := range toTest {
		seen := false
		for _, v := range validTerms {
			if v == term {
				seen = true
			}
		}
		if !seen {
			return false
		}
	}
	return true
}
func validateOutputFormat(toTest outputFmt) bool {
	return toTest.isJson() || toTest.isCsv() || toTest.isQuiet() || toTest.isTab()
}

// attempts to open outfile, if set. If is not set, returns std.out writer
func initOutputWriter(outfile string) io.Writer {
	if len(outfile) == 0 {
		return os.Stdout
	} else {
		file, err := os.Create(outfile)
		if err != nil {
			exitWithErr(err)
		}
		return file
	}
}

func initWebClient() *rspace.RsWebClient {
	if len(getenv(BASE_URL_ENV_NAME)) == 0 {
		exitWithStdErrMsg("No URL for RSpace  detected")
	}
	url, _ := url.Parse(getenv(BASE_URL_ENV_NAME))
	apikey := getenv(APIKEY_ENV_NAME)
	if len(apikey) == 0 {
		exitWithStdErrMsg("No API key detected")
	}
	webClient := rspace.NewWebClient(url, apikey)
	return webClient
}
func getenv(envname string) string {
	return os.Getenv(envname)
}

// common setup for a paginating command
func initPaginationFromArgs(cmd *cobra.Command) {
	cmd.Flags().StringVar(&sortOrderArg, "sortOrder", "desc", "'asc' or 'desc'")
	cmd.Flags().StringVar(&orderByArg, "orderBy", "lastModified", "orders results by 'name', 'created' or 'lastModified'")
	cmd.Flags().IntVar(&pageSizeArg, "maxResults", 20, "Maximum number of results to retrieve")
}
