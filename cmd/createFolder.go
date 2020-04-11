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

	"github.com/spf13/cobra"
	"rspace"
	"strconv"
	"encoding/json"
	"os"
)
var foldername string
var parentFolder string
// createNotebookCmd represents the createNotebook command
var createFolderCmd = &cobra.Command{
	Use:   "createFolder",
	Short: "Creates a new Folder",
	Long: `Create a new Folder, with an optional name and parent folder
	  create-folder --name foldername --infolder FL1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		post := rspace.FolderPost{IsNotebook:false,}
		doCreateFolder(foldername, parentFolder, post)	
	},
}
func marshal(anything interface{}) string {
        bytes, _ := json.MarshalIndent(anything, "", "\t")
        return string(bytes)
}

func doCreateFolder (foldername string, parentFolder string, post rspace.FolderPost) {
		if len(foldername) > 0 {
			post.Name=foldername
		}
		if len(parentFolder) > 0 {
			id,err :=strconv.Atoi(parentFolder)
			if err != nil {
				fmt.Println("Please supply a numeric folder id for the parent folder")
				os.Exit(1)
			}
			post.ParentFolderId=id
		}
		webClient:=setup()
		got, err := webClient.FolderNew(&post)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if isQuiet {
			fmt.Println(got.Id)
		} else if isVerbose{
			fmt.Println(marshal(got))
		} else {
			fmt.Println(got.GlobalId)
		}
}

func init() {
	elnCmd.AddCommand(createFolderCmd)
	 createFolderCmd.Flags().StringVarP(&foldername, "name", "n", "", "A name for the folder")
	 createFolderCmd.Flags().StringVarP(&parentFolder, "folder", "f", "", "An id for the folder that will contain the new folder")
}
