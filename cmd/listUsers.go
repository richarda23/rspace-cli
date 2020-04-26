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
	"strconv"
	"time"

)
// listDocumentsCmd represents the listDocuments command
var listUsersCmd = &cobra.Command{
	Use:   "listUsers",
	Short: "Lists users",
	Long:`List users, sorted or paginated, e.g.

		  rspace eln listusers
	`,

	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()  
		//cfg := configurePagination()
		doListusers(context)
	},
}

func doListusers (ctx *Context) {
	var usersList *rspace.UserList
	var err error
	usersList, err = ctx.WebClient.Users(time.Time{}, time.Time{})
	if err != nil {
		exitWithErr(err)
	}
	formatter := &UserListFormatter{usersList}
	ctx.writeResult(formatter)
}

type UserListFormatter struct {
	*rspace.UserList
}

func (fs *UserListFormatter) ToJson () string{
	return prettyMarshal(fs.UserList)
}

func (ds *UserListFormatter) ToQuiet () []identifiable{
	rows := make([]identifiable, 0)
	for _, res := range ds.UserList.Users {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}

func (ds *UserListFormatter) ToTable () *TableResult {
	results := ds.UserList.Users

	headers := []columnDef {columnDef{"Id",8}, columnDef{"Username",10}, columnDef{"FirstName", 25}, 
	 columnDef{"LastName", 25}}

	rows := make([][]string, 0)
	for _, res := range results {
		data := []string {strconv.Itoa(res.Id),res.Username, res.FirstName,
			   res.LastName}
		rows = append(rows, data)
	}
	return &TableResult {headers, rows}
	
 }
 func toIdentifiableUser (results []*rspace.UserInfo) []identifiable {
	rows := make([]identifiable, 0)
	
	for _, res := range results {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}
func init() {
	elnCmd.AddCommand(listUsersCmd)
	//initPaginationFromArgs(listUsersCmd)
	// and all subcommands, e.g.:
	//listUsersCmd.PersistentFlags().StringVar(&mediaTypeArg, "mediaType", "", "Optional media type, 1 of '
}