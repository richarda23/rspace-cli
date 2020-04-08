package cmd
import (
"fmt"
"net/url"
"rspace"
"os"

)

const (
	APIKEY_ENV_NAME   = "RSPACE_API_KEY"
	BASE_URL_ENV_NAME = "RSPACE_URL"
)
func setup () *rspace.RsWebClient {
	if len(getenv(BASE_URL_ENV_NAME)) ==0 {
		fmt.Println("No URL for RSpace  detected")
		os.Exit(1)
	}
	url, _ := url.Parse(getenv(BASE_URL_ENV_NAME))
	fmt.Println("url is " + url.String())
	apikey := getenv(APIKEY_ENV_NAME)
	fmt.Println("api is " + apikey)
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
