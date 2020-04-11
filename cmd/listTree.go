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
	"rspace"
	"github.com/spf13/cobra"
)

var folderId int 
// listTreeCmd represents the listTree command
var listTreeCmd = &cobra.Command{
	Use:   "listTree",
	Short: "Lists the contents of a folder or notebook",
	Long: `Lists the content of the specified folder, or the home folder if no folder ID is set`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listTree called")
		outputFormat = outputFmt(outputFormatArg)
		validateOutputFormatExit (outputFormat)
		webClient:=setup()
		cfg := rspace.NewRecordListingConfig()
		folderList, err := webClient.FolderTree(cfg, folderId,make ([]string,0)) 
		if err != nil {
			exitWithErr(err)
		}
		fmt.Println(prettyMarshal(folderList))
	},
}

func init() {
	elnCmd.AddCommand(listTreeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listTreeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	 listTreeCmd.Flags().IntVar(&folderId, "folder",  0, "The id of the folder or notebook to list")
}
