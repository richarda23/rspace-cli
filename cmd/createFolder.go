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
	"strconv"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

var createFolderArgS = createFolderArg{}

var createFolderCmd = &cobra.Command{
	Use:   "createFolder",
	Short: "Creates a new Folder",
	Long: `
Create a new Folder, with an optional name and parent folder ID
	`,
	Example: `
// make a new folder in folder with id FL1234
rspace eln createFolder --name MyFolder --folder FL1234
	`,
	Run: func(cmd *cobra.Command, args []string) {
		post := rspace.FolderPost{IsNotebook: false}
		ctx := initialiseContext()
		doCreateFolder(ctx, createFolderArgS, post)
	},
}

func doCreateFolder(ctx *Context, args createFolderArg, post rspace.FolderPost) {
	if len(args.Name) > 0 {
		post.Name = args.Name
	}
	if len(args.ParentFolder) > 0 {
		id, err := idFromGlobalId(args.ParentFolder)
		if err != nil {
			exitWithStdErrMsg("Please supply a folder id for the parent folder")
		}
		post.ParentFolderId = id
	}

	got, err := ctx.WebClient.FolderNew(&post)
	if err != nil {
		exitWithErr(err)
	}
	if ctx.Format.isJson() {
		ctx.write(prettyMarshal(got))
	} else if ctx.Format.isTab() || ctx.Format.isCsv() {
		folderToTable(ctx, got)
	} else if ctx.Format.isQuiet() {
		ctx.write(strconv.Itoa(got.Id))
	} else {
		ctx.write("unknown format")
	}
}
func folderToTable(ctx *Context, folder *rspace.Folder) {
	headers := []columnDef{columnDef{"Id", 8}, columnDef{"GlobalId", 10}, columnDef{"Name", 25}, columnDef{"ParentFolderId", 15}, columnDef{"Created", 24}}
	data := []string{strconv.Itoa(folder.Id), folder.GlobalId, folder.Name, strconv.Itoa(folder.ParentFolderId), folder.Created}
	rows := make([][]string, 0)
	rows = append(rows, data)
	if ctx.Format.isCsv() {
		printCsv(ctx, &TableResult{headers, rows})
	} else {
		printTable(ctx, &TableResult{headers, rows})
	}
}
func init() {
	elnCmd.AddCommand(createFolderCmd)
	createFolderCmd.Flags().StringVarP(&createFolderArgS.Name, "name", "n", "", "A name for the folder")
	createFolderCmd.Flags().StringVarP(&createFolderArgS.ParentFolder, "folder", "p", "", "An id for the folder that will contain the new folder")
}
