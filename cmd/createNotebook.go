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
	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

type createFolderArg struct {
	Name, ParentFolder string
}

var createNotebookArgS createFolderArg

// createNotebookCmd represents the createNotebook command
var createNotebookCmd = &cobra.Command{
	Use:   "createNotebook",
	Short: "Creates a new notebook",
	Long: `Create a new notebook, with an optional name and parent folder
	`,
	Example: `
		// create a new notebook 'MyNotebook' in folder FL1234
		rspace eln createNotebook --name MyNotebook --folder FL1234

		//create an unnamed notebook in home folder
		rspace eln createNotebook
	`,
	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()
		doCreateNotebook(context, createFolderArgS)
	},
}

func doCreateNotebook(ctx *Context, args createFolderArg) {
	post := rspace.FolderPost{IsNotebook: true}
	doCreateFolder(ctx, args, post)
}

func init() {
	elnCmd.AddCommand(createNotebookCmd)
	createNotebookCmd.Flags().StringVarP(&createNotebookArgS.Name, "name", "n", "", "A name for the notebook")
	createNotebookCmd.Flags().StringVarP(&createNotebookArgS.ParentFolder, "folder", "p", "", "An id for the folder that will contain the new notebook")
}
