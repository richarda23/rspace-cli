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
	"github.com/spf13/cobra"
	"github.com/richarda23/rspace-client-go/rspace"
	"os"
	"fmt"
)

type downloadArgs struct {
	OutfolderArg string
}

var dArgs= downloadArgs{}


// createNotebookCmd represents the createNotebook command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads a file with the given id",
	Long: `Downloads a file with given id. outfile is optional; if not set
	will download to current folder.
	`,
	Args: cobra.MinimumNArgs(1),
	
	Run: func(cmd *cobra.Command, args []string) {
	//	post := rspace.FolderPost{IsNotebook: false}
		ctx := initialiseContext()
		ids:=validateDownloadArgs(args)
		info := doDownload(ctx, ids)
		ctx.writeResult(info)
	},
}
// TODO hande multiple FileIds
func doDownload (ctx *Context, ids []int) *FileListFormatter {
	var results  = make([]*rspace.FileInfo, 0)
	for _,id := range ids {
		info,err :=ctx.WebClient.Download(id, dArgs.OutfolderArg)
		if err != nil {
			messageStdErr(err.Error())
		} else {			
			results = append(results, info)
		}
	}
	var fal = FileArrayList{results}
	return &FileListFormatter{fal}
}

func validateDownloadArgs (args [] string) []int{
	ids:=make([]int, 0)
	for _,idStr:=range args {
		id,err:=idFromGlobalId(idStr)
		if err != nil {
			messageStdErr(idStr + " is not valid id, skipping")
		}
		ids=append(ids, id)
	}
	if len(dArgs.OutfolderArg) > 0 {
		stats,err := os.Stat(dArgs.OutfolderArg)
		if err != nil {
			exitWithErr(err)
		}
		if !stats.IsDir() {
			exitWithStdErrMsg("Not a directory")
		}
	} else {
		dArgs.OutfolderArg="./"
	}
	fmt.Println(ids)
	return ids
}
	
func init() {
	elnCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVar(&dArgs.OutfolderArg, "dir",  "", "Optional directory to download into")
}
 