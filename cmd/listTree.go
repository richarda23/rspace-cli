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
	"rspace"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

var folderId int 
var treeFilterArg string
// listTreeCmd represents the listTree command
var listTreeCmd = &cobra.Command{
	Use:   "listTree",
	Short: "Lists the contents of a folder or notebook",
	Long: `Lists the content of the specified folder, or the home folder if no folder ID is set`,

	Run: func(cmd *cobra.Command, args []string) {
		context:=initialiseContext()  
		cfg:=configurePagination()
		doListTree(context, cfg)
	},
}
func configurePagination () rspace.RecordListingConfig{
	cfg := rspace.NewRecordListingConfig()
	if len(sortOrderArg)  > 0 && validateArrayContains(validSortOrders, []string{sortOrderArg}) {
		cfg.SortOrder = sortOrderArg
	}
	if len(orderByArg)  > 0 && validateArrayContains(validRecordOrders, []string{orderByArg}) {
		cfg.OrderBy = orderByArg
	}
	if pageSizeArg > 0  {
		cfg.PageSize = pageSizeArg
	}
	return cfg
}
func doListTree (ctx *Context, cfg rspace.RecordListingConfig) {
	var filters =  make([]string, 0)
	if len (treeFilterArg) > 0 {
		filters = strings.Split(treeFilterArg,",")
	}
	folderList, err := ctx.WebClient.FolderTree(cfg, folderId, filters) 
	if err != nil {
		exitWithErr(err)
	}
	if ctx.Format.isJson() {
		ctx.write(prettyMarshal(folderList))
	} else if ctx.Format.isQuiet() {
		printIds(ctx, toIdentifiable(folderList))
	} else {
		listToTable(ctx, folderList)
	}
}
func toIdentifiable (results *rspace.FolderList) []identifiable {
	rows := make([]identifiable, 0)
	for _, res := range results.Records {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}

func listToTable(ctx *Context, results *rspace.FolderList) {
	headers := []columnDef {columnDef{"Id",8}, columnDef{"GlobalId",10},  columnDef{"Name", 25},columnDef{"Type", 9},  columnDef{"Created",24},columnDef{"Last Modified",24}}

	rows := make([][]string, 0)
	for _, res := range results.Records {
		data := []string {strconv.Itoa(res.Id),res.GlobalId, res.Name, res.Type,  res.Created,res.LastModified}
		rows = append(rows, data)
	}
	if ctx.Format.isCsv() {
		printCsv(ctx, headers, rows)
	} else {
		printTable(ctx, headers, rows)
	}
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
	 listTreeCmd.Flags().StringVar(&treeFilterArg, "filter",  "", "Restrict results to 1 or more of: " + strings.Join(validTreeFilters, ","))
	 listTreeCmd.Flags().StringVar(&sortOrderArg, "sortOrder",  "desc", "'asc' or 'desc'")
	 listTreeCmd.Flags().StringVar(&orderByArg, "orderBy",  "lastModified", "orders results by 'name', 'created' or 'lastModified'")
	 listTreeCmd.Flags().IntVar(&pageSizeArg, "maxResults",  20, "Maximum number of results to retrieve")
}