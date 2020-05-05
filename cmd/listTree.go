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
	"strings"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

var folderIdArg string
var treeFilterArg string

// listTreeCmd represents the listTree command
var listTreeCmd = &cobra.Command{
	Use:   "listTree",
	Short: "Lists the contents of a folder or notebook",
	Long: `Lists the content of the specified folder, or the home folder if no folder ID is set,
	  optionally filtered by documents or notebooks`,
	Example: `
		// show subfolders of home folder
		rspace eln listTree --filter folder

		// show notebooks in specified folder
		rspace eln listTree --filter notebook --folder 1234


	`,

	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()
		cfg := configurePagination()
		doListTree(context, cfg)
	},
}

func doListTree(ctx *Context, cfg rspace.RecordListingConfig) {
	var filters = make([]string, 0)
	if len(treeFilterArg) > 0 {
		filters = strings.Split(treeFilterArg, ",")
	}
	folderId, _ := idFromGlobalId(folderIdArg)
	folderList, err := ctx.WebClient.FolderTree(cfg, folderId, filters)
	if err != nil {
		exitWithErr(err)
	}
	ctx.writeResult(&FolderListFormatter{folderList})
}

type FolderListFormatter struct {
	*rspace.FolderList
}

func (ds *FolderListFormatter) ToJson() string {
	return prettyMarshal(ds.FolderList)
}

func (ds *FolderListFormatter) ToQuiet() []identifiable {
	return toIdentifiable(ds.FolderList)
}

func (ds *FolderListFormatter) ToTable() *TableResult {
	baseInfos := resultsToBaseInfoList(ds.FolderList)
	maxNameCol := getMaxNameLength(baseInfos)
	headers := []columnDef{columnDef{"Id", 8}, columnDef{"GlobalId", 10}, columnDef{"Name", maxNameCol}, columnDef{"Type", 9},
		columnDef{"Created", DISPLAY_TIMESTAMP_WIDTH}, columnDef{"Last Modified", DISPLAY_TIMESTAMP_WIDTH}}

	rows := make([][]string, 0)
	for _, res := range ds.FolderList.Records {
		data := []string{strconv.Itoa(res.Id), res.GlobalId, res.Name, res.Type, res.Created[0:DISPLAY_TIMESTAMP_WIDTH], res.LastModified[0:DISPLAY_TIMESTAMP_WIDTH]}
		rows = append(rows, data)

	}
	table := &TableResult{headers, rows}
	return table
}

func toIdentifiable(results *rspace.FolderList) []identifiable {
	rows := make([]identifiable, 0)
	for _, res := range results.Records {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}

func resultsToBaseInfoList(results *rspace.FolderList) []rspace.BasicInfo {
	var baseResults = make([]rspace.BasicInfo, len(results.Records))
	for i, v := range results.Records {
		var x rspace.BasicInfo = v
		baseResults[i] = x
	}
	return baseResults
}

func init() {
	elnCmd.AddCommand(listTreeCmd)
	// is called directly, e.g.:
	listTreeCmd.Flags().StringVar(&folderIdArg, "folder", "", "The id or global Id of the folder or notebook to list")
	listTreeCmd.Flags().StringVar(&treeFilterArg, "filter", "", "Restrict results to 1 or more of: "+strings.Join(validTreeFilters, ","))

	initPaginationFromArgs(listTreeCmd)
}
