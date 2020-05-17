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

type addFolderArg struct {
	Name, ParentFolder string
}

var addNotebookArgS addFolderArg

// addNotebookCmd represents the addNotebook command
var addNotebookCmd = &cobra.Command{
	Use:   "addNotebook",
	Short: "Creates a new notebook",
	Long: `Create a new notebook, with an optional name and parent folder
	`,
	Aliases: []string{"createNotebook", "notebookAdd", "nbAdd"},
	Example: `
// add a new notebook 'MyNotebook' in folder FL1234
rspace eln addNotebook --name MyNotebook --folder FL1234

//add an unnamed notebook in home folder
rspace eln addNotebook
	`,
	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()
		doAddNotebook(context, addFolderArgS)
	},
}

func doAddNotebook(ctx *Context, args addFolderArg) {
	post := rspace.FolderPost{IsNotebook: true}
	doAddFolder(ctx, args, post)
}

func init() {
	elnCmd.AddCommand(addNotebookCmd)
	addNotebookCmd.Flags().StringVarP(&addNotebookArgS.Name, "name", "n", "", "A name for the notebook")
	addNotebookCmd.Flags().StringVarP(&addNotebookArgS.ParentFolder, "folder", "p", "", "An id for the folder that will contain the new notebook")
}
