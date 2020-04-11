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
	"rspace"
)
var notebookName string
var parentfolder string
// createNotebookCmd represents the createNotebook command
var createNotebookCmd = &cobra.Command{
	Use:   "createNotebook",
	Short: "Creates a new notebooks",
	Long: `Create a new notebook, with an optional name and parent folder
	  create-notebook --name nbname --infolder FL1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		outputFormat = outputFmt(outputFormatArg)
		validateOutputFormatExit (outputFormat)
		post := rspace.FolderPost{IsNotebook:true,}
		doCreateFolder(notebookName, parentfolder, post)
	},
}

func init() {
	elnCmd.AddCommand(createNotebookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createNotebookCmd.PersistentFlags().String("name", "n","",  "A name for the notebook")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	 createNotebookCmd.Flags().StringVarP(&notebookName, "name", "n", "", "A name for the notebook")
	 createNotebookCmd.Flags().StringVarP(&parentfolder, "folder", "p", "", "An id for the folder that will contain the new notebook")
}
