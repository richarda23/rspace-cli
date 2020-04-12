package cmd
import (
"fmt"
"net/url"
"rspace"
"os"
"io"
"strings"

)
const (
	APIKEY_ENV_NAME   = "RSPACE_API_KEY"
	BASE_URL_ENV_NAME = "RSPACE_URL"
)
var (
	validOutputFormats = []string {"json", "csv", "quiet", "table",}
	outputFormat outputFmt
)
type outputFmt string

func (ft outputFmt) isJson () bool {
	return string(ft) == "json";
}
func (ft outputFmt) isCsv () bool {
	return ft == "csv";
}
func (ft outputFmt) isTab () bool {
	return ft == "table";
}
func (ft outputFmt) isQuiet () bool {
	return ft == "quiet";
}


type Context struct {
 WebClient *rspace.RsWebClient
 Writer io.Writer
 Format outputFmt
}

func initialiseContext () *Context {
	rc := Context{}
	rc.WebClient = setup()
	rc.Writer = initOutputWriter(outFileArg)
	outputFormat = outputFmt(outputFormatArg)
	validateOutputFormatExit (outputFormat)
	rc.Format = outputFormat
	return &rc
}
func validateOutputFormatExit (toTest outputFmt) {
	if !validateOutputFormat(toTest) {
		exitWithStdErrMsg("Invalid outputFormat argument: must be one of: " + strings.Join(validOutputFormats,","))
	}
}

func validateOutputFormat (toTest outputFmt) bool {
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
	return nil
}

func setup () *rspace.RsWebClient {
	if len(getenv(BASE_URL_ENV_NAME)) ==0 {
		fmt.Println("No URL for RSpace  detected")
		os.Exit(1)
	}
	url, _ := url.Parse(getenv(BASE_URL_ENV_NAME))
	apikey := getenv(APIKEY_ENV_NAME)
	if len(apikey) ==0 {
		fmt.Println("No API key detected")
		os.Exit(1)
	}
	webClient := rspace.NewWebClient(url, apikey)
	return webClient
}
func getenv(envname string) string {
	return os.Getenv(envname)
}
