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

// listDocumentsCmd represents the listDocuments command
var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "Lists the groups you are a member of.",
	Long:  `List all groups that you are a member of`,
	Example: `
rspace eln groups
	`,

	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()
		doListGroups(context)
	},
}

func doListGroups(ctx *Context) {
	var groupList *rspace.GroupList
	var err error
	groupList, err = ctx.WebClient.Groups()
	if err != nil {
		exitWithErr(err)
	}
	formatter := &GroupListFormatter{groupList}
	ctx.writeResult(formatter)
}

type GroupListFormatter struct {
	*rspace.GroupList
}

func (fs *GroupListFormatter) ToJson() string {
	return prettyMarshal(fs.GroupList.Groups)
}

func (ds *GroupListFormatter) ToQuiet() []identifiable {
	rows := make([]identifiable, 0)
	for _, res := range ds.GroupList.Groups {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}

func (ds *GroupListFormatter) ToTable() *TableResult {
	results := ds.GroupList.Groups

	headers := []columnDef{columnDef{"Id", 8}, columnDef{"Name", 25},
		columnDef{"Type", 10}, columnDef{"SharedFolderId", 16}}

	rows := make([][]string, 0)
	for _, res := range results {
		data := []string{strconv.Itoa(res.Id), res.Name, res.Type, strconv.Itoa(res.SharedFolderId)}
		rows = append(rows, data)
	}
	return &TableResult{headers, rows}

}
func init() {
	elnCmd.AddCommand(listGroupsCmd)
	initPaginationFromArgs(listGroupsCmd)
}
