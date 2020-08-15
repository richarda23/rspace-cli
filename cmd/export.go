/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

type exportCmdArgs struct {
	Scope  string
	Format string
	// user or group id
	Id int
}

var exportCmdArgsArg exportCmdArgs

// im
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports RSpace documents to XML or HTML archive",
	Long: ` Import Word files as RSpace document. Add files and folders to the command line. 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// initial wait for job might take some time
		ctx := initialiseContextWithTimeout(1200)
		exportArgs(ctx, args)
	},
}

func exportArgs(ctx *Context, args []string) {
	scope := getExportScope(exportCmdArgsArg.Scope)
	format := getExportFormat(exportCmdArgsArg.Format)
	id := exportCmdArgsArg.Id
	post := rspace.ExportPost{format, scope, id}
	result, err := ctx.WebClient.Export(post)
	if err != nil {
		exitWithErr(err)
	}
	if result.IsCompleted() {
		ctx.write(fmt.Sprintf("Completed - download link is %s (%s)",
			result.DownloadLink(), humanizeBytes(uint64(result.Result.Size))))
	}

}

func getExportFormat(format string) rspace.ExportFormat {
	switch format {
	case "xml":
		return rspace.XML_FORMAT
	case "html":
		return rspace.HTML_FORMAT
	}
	exitWithStdErrMsg("export format must be 'xml' or 'html'")
	return 0
}

func getExportScope(arg string) rspace.ExportScope {
	switch arg {
	case "user":
		return rspace.USER_EXPORT_SCOPE
	case "group":
		return rspace.GROUP_EXPORT_SCOPE
	}
	exitWithStdErrMsg("export scope must be 'user' or 'group'")
	return 0
}

func init() {
	elnCmd.AddCommand(exportCmd)
	exportCmd.PersistentFlags().StringVar(&exportCmdArgsArg.Scope,
		"scope", "user", "user or group")
	exportCmd.PersistentFlags().StringVar(&exportCmdArgsArg.Format,
		"format", "html", "xml or html")
	exportCmd.PersistentFlags().IntVar(&exportCmdArgsArg.Id,
		"id", 0, "User or group id to export")

}
