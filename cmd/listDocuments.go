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
)
var Quiet bool
var PageSize int
var SearchQuery string

// listDocumentsCmd represents the listDocuments command
var listDocumentsCmd = &cobra.Command{
	Use:   "listDocuments",
	Short: "Lists the documents",
	Long:` Lists documents `,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listDocuments called")
		context := initialiseContext()  
		cfg := configurePagination()
		doListDocs(context, cfg)
	},
}
func doListDocs (ctx *Context, cfg rspace.RecordListingConfig) {
	var doclist *rspace.DocumentList
	var err error
	if len(SearchQuery) > 0 {
		doclist,err=ctx.WebClient.DocumentS.SearchDocuments(cfg, SearchQuery)
	} else {
		doclist, err = ctx.WebClient.DocumentS.Documents(cfg )
	}
	if err != nil {
		exitWithErr(err)
	}
	if ctx.Format.isJson() {
		ctx.write(prettyMarshal(doclist))
	} else if ctx.Format.isQuiet() {
		printIds(ctx, toIdentifiableDoc(doclist))
	} else {
		listToDocTable(ctx, doclist)
	}
}
func listToDocTable(ctx *Context, results *rspace.DocumentList) {
	headers := []columnDef {columnDef{"Id",8}, columnDef{"GlobalId",10},  columnDef{"Name", 25},  columnDef{"Created",24},columnDef{"Last Modified",24}}

	rows := make([][]string, 0)
	for _, res := range results.Documents {
		data := []string {strconv.Itoa(res.Id),res.GlobalId, res.Name,   res.Created,res.LastModified}
		rows = append(rows, data)
	}
	if ctx.Format.isCsv() {
		printCsv(ctx, headers, rows)
	} else {
		printTable(ctx, headers, rows)
	}
}
func toIdentifiableDoc (results *rspace.DocumentList) []identifiable {
	rows := make([]identifiable, 0)
	for _, res := range results.Documents {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}
func init() {
	elnCmd.AddCommand(listDocumentsCmd)

	initPaginationFromArgs(listDocumentsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	listDocumentsCmd.PersistentFlags().StringVar(&SearchQuery, "query", "", "Search query term")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listDocumentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
