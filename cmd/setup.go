package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net/url"
	"os"
	"github.com/richarda23/rspace-client-go/rspace"
	"strings"
	"github.com/spf13/viper"
)

const (
	APIKEY_ENV_NAME   = "RSPACE_API_KEY"
	BASE_URL_ENV_NAME = "RSPACE_URL"
)

var (
	validOutputFormats = []string{"json", "csv", "quiet", "table"}
	// for listTree
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
// Context maintains references to the webClient and result Writers
type Context struct {
	WebClient *rspace.RsWebClient
	Writer    io.Writer
	Format    outputFmt
}

//Writes the result of a command to one of the supported output formats
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
// writes a string to the output stream (either stdout or a defined file)
func (ctx *Context) write(toWrite string) {
	fmt.Fprintln(ctx.Writer, toWrite)
}

// main initialisation method. Creates an RsWebClient, output writer and format
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
// asserts that an outputFmt argument is valid
func validateOutputFormat(toTest outputFmt) bool {
	return toTest.isJson() || toTest.isCsv() || toTest.isQuiet() || toTest.isTab()
}
// returns an io.Writer for a log file. If logfile is empty, return default writer.
func initLogWriter (logfile string, defaultWriter io.Writer) io.Writer {
	if len(logfile) == 0 {
		return defaultWriter
	} else {
		file, err := os.Create(logfile)
		if err != nil {
			exitWithErr(err)
		}
		return file
	}
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

// reads apikey and url from viper configuration,
// then sets these into an RsWebClient instance
func initWebClient() *rspace.RsWebClient {
	urlCfg,ok:=viper.Get(BASE_URL_ENV_NAME).(string)
	if !ok || len(urlCfg) == 0 {
		exitWithStdErrMsg("No URL for RSpace  detected")
	}
	url, _ := url.Parse(urlCfg)
	messageStdErr("RSpace URL: " + urlCfg)
	apikey,ok := viper.Get(APIKEY_ENV_NAME).(string)
	if !ok || len(apikey) == 0 {
		exitWithStdErrMsg("No API key detected")
	}
	messageStdErr("Api key:" + apikey[0:4] + "...")
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
