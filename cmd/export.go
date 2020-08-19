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
	"strconv"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

type exportCmdArgs struct {
	Scope  string
	Format string
	// block for export to complete
	Wait         bool
	MaxLinkLevel int
}

var exportCmdArgsArg exportCmdArgs

// im
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports RSpace documents to XML or HTML archive",
	Long: `
Make an export of your work to zipped archive in  HTML (default) or XML format. If you are a PI or admin
user you can export other users' work or a group's work by supplying a group or user ID as a 
single argument.

Exports of selections are exported as well, by supplying IDs as command line arguments

You can opt to wait for the export to complete using --wait. This will cause the command to 
block until the export process has completed.

Launching an export returns a job Id that you can use to download the results using 'job' command.
`,
	Example: `
// export your own work to HTML, waiting for the archive process to complete
rspace eln export --format html --scope user --wait

// submit an export but don't wait for completion
rspace eln export --format xml --scope user

// export your group's work (you need to be a PI or labAdmin to do this)
rspace eln export  12345 --format xml --scope group

// export a selection of records, notebooks or folders (3 in this case)
rspace eln export SD123 NB456 FL789 --format html --scope selection

// export - don't include linked documents. You can use numeric IDs if you prefer.
rspace eln export 123 456 --linkDepth 0
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
	var itemIds []int
	var userOrGroupId int
	if len(args) > 0 {
		if scope == rspace.SELECTION_EXPORT_SCOPE {
			itemIds = globalIdListToIntList(args)
		} else {
			if len(args) > 1 {
				exitWithStdErrMsg("Only a single user or group id is supported")
			}
			var err error
			userOrGroupId, err = strconv.Atoi(args[0])
			if err != nil {
				exitWithErr(err)
			}
		}
	}
	post := rspace.ExportPost{format, scope, userOrGroupId, itemIds, exportCmdArgsArg.MaxLinkLevel}
	messageStdErr("Waiting for export to start...")
	if exportCmdArgsArg.Wait {

		result, err := ctx.WebClient.Export(post, true)
		if err != nil {
			exitWithErr(err)
		}
		if result.IsCompleted() {
			ctx.writeResult(&JobFormatter{result})
		}
	} else {
		result, err := ctx.WebClient.Export(post, false)
		if err != nil {
			exitWithErr(err)
		}
		ctx.writeResult(&JobFormatter{result})
	}
}

type JobFormatter struct {
	*rspace.Job
}

func (fs *JobFormatter) ToQuiet() []identifiable {
	rows := make([]identifiable, 0)
	rows = append(rows, identifiable{strconv.Itoa(fs.Job.Id)})
	return rows
}

func (fs *JobFormatter) ToTable() *TableResult {
	headers := []columnDef{columnDef{"Id", 8}, columnDef{"Status", 10},
		columnDef{"Percent Complete", 18}, columnDef{"Download size", 14}}

	rows := make([][]string, 0)
	var sizeStr = "unknown"
	if fs.Job.IsCompleted() {
		sizeStr = humanizeBytes(uint64(fs.Job.Result.Size))
	}

	data := []string{strconv.Itoa(fs.Job.Id), fs.Job.Status,
		fmt.Sprintf("%3.2f", fs.Job.PercentComplete), sizeStr}
	rows = append(rows, data)
	return &TableResult{headers, rows}
}
func (fs *JobFormatter) ToJson() string {
	return prettyMarshal(fs.Job)
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
	case "selection":
		return rspace.SELECTION_EXPORT_SCOPE
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
	exportCmd.PersistentFlags().BoolVar(&exportCmdArgsArg.Wait,
		"wait", false, "Wait for export to complete")
	exportCmd.PersistentFlags().IntVar(&exportCmdArgsArg.MaxLinkLevel,
		"linkDepth", 1, "Maximum number of links to follow to include in export")
}
