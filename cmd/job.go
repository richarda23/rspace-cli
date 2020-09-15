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
	"os"
	"path/filepath"
	"strconv"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

type jobCmdArgs struct {
	// Download if complete
	Download bool
	// optional download path
	DownloadPath string
}

var jobCmdArgsArg jobCmdArgs

// im
var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "Query progress of a Job",
	Long: ` Query a job status, or download result. You get a jobId after submitting an export
	 request using the 'export' command.
	`,
	Example: `
// get progress of job in tabular format
rspace eln job  22

// get raw JSON
rspace eln job 22 -f json

// download (if complete) to a file of your choice
rspace eln job  22 --download --outfile /home/me/myexport.zip

// download (if complete) to current directory
rspace eln job  22 --download
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// initial wait for job might take some time
		ctx := initialiseContext()
		doJob(ctx, args)
	},
}

func doJob(ctx *Context, args []string) {
	id, e := strconv.Atoi(args[0])
	if e != nil || id <= 0 {
		exitWithStdErrMsg("Invalid job ID, must be an integer > 0, but was " + args[0])
	}
	download := jobCmdArgsArg.Download
	result, err := ctx.WebClient.GetJob(id)
	if err != nil {
		exitWithErr(err)
	}
	ctx.writeResult(&JobFormatter{result})
	if download {
		if result.IsCompleted() {
			downloadpath := getOutfile(result)
			path, err := os.Create(downloadpath)
			if err != nil {
				exitWithErr(err)
			}
			messageStdErr(fmt.Sprintf("downloading to %s (%s)", downloadpath, humanizeBytes(uint64(result.Result.Size))))
			err = ctx.WebClient.DownloadExport(result.DownloadLink(), path)
			if err != nil {
				exitWithErr(err)
			}
		} else {
			messageStdErr(fmt.Sprintf("Job %d not completed, nothing to download", result.Id))
		}
	}
}

func getOutfile(job *rspace.Job) string {
	downloadpath := jobCmdArgsArg.DownloadPath
	if len(downloadpath) == 0 {
		link := job.DownloadLink()
		path := link.Path
		downloadpath = filepath.Base(path)
	}
	return downloadpath

}
func init() {
	elnCmd.AddCommand(jobCmd)
	jobCmd.PersistentFlags().BoolVar(&jobCmdArgsArg.Download,
		"download", false, "Download result, if complete")
	jobCmd.PersistentFlags().StringVar(&jobCmdArgsArg.DownloadPath,
		"outfile", "", "file path to download to...")
}
