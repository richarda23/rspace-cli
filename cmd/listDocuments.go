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
)
var Quiet bool
var PageSize int
var searchQuery string
// arguments for advanced search
var orSrchArg bool = false
var nameSearchArg string
var tagSrchArg string
var createdBeforeSrchArg string
var createdAfterSrcArg string
var modifiedBeforeSrchArg string
var modifiedAfterSrchArg string

// listDocumentsCmd represents the listDocuments command
var listDocumentsCmd = &cobra.Command{
	Use:   "listDocuments",
	Short: "Lists the documents",
	Long:` Lists documents. Search is optional E.g.
	
		// A global search over names, tags, file and text content
		 rspace eln listDocuments --query myexperiment

		// Get documents whose name starts with 'experiment' OR is created in 2020 or later
		rspace eln listDocuments --or --name experiment* --createdAfter 2020-01-01

		//  list the 1st hundred documents created
		rspace eln listDocuments --orderBy created --sortOrder asc --maxResults 100
	`,

	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()  
		cfg := configurePagination()
		doListDocs(context, cfg)
	},
}

func doListDocs (ctx *Context, cfg rspace.RecordListingConfig) {
	var docList *rspace.DocumentList
	var err error
	if len(searchQuery) > 0 {
		docList,err=ctx.WebClient.SearchDocuments(cfg, searchQuery)
	} else if (advancedSrchArgsAreProvided()){
		docList,err = doAdvancedSearch(ctx, cfg)
	} else {
		docList, err = ctx.WebClient.Documents(cfg )
	}

	
	if err != nil {
		exitWithErr(err)
	}
	if ctx.Format.isJson() {
		ctx.write(prettyMarshal(docList))
	} else if ctx.Format.isQuiet() {
		printIds(ctx, toIdentifiableDoc(docList))
	} else {
		listToDocTable(ctx, docList)
	}
}
func advancedSrchArgsAreProvided() bool {
	var advSearchArgs = []string {nameSearchArg, tagSrchArg, createdBeforeSrchArg, createdAfterSrcArg,
		modifiedAfterSrchArg, modifiedBeforeSrchArg}
	for _, v := range advSearchArgs {
		messageStdErr(v)
		if len(v) > 0 {
			return true
		}
	}
	return false
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

func doAdvancedSearch (ctx *Context, cfg rspace.RecordListingConfig)(*rspace.DocumentList, error){
	messageStdErr("in advanced search")
	builder := rspace.SearchQueryBuilder{}
	if orSrchArg {
		builder.Operator(rspace.Or)
	}
	builder.AddTerm(nameSearchArg, rspace.NAME)
	builder.AddTerm(tagSrchArg, rspace.TAG)
	builder.AddTerm(createdAfterSrcArg + ";" + createdBeforeSrchArg, rspace.CREATED)
	builder.AddTerm(modifiedAfterSrchArg + ";" + modifiedBeforeSrchArg, rspace.LAST_MODIFIED)
	q:=builder.Build()
	return ctx.WebClient.AdvancedSearchDocuments(cfg, q)
}


func init() {
	elnCmd.AddCommand(listDocumentsCmd)

	initPaginationFromArgs(listDocumentsCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	listDocumentsCmd.PersistentFlags().StringVar(&searchQuery, "query", "", "Search query term")
	listDocumentsCmd.PersistentFlags().BoolVar(&orSrchArg, "or", false, "Combines search terms together with boolean 'OR'")
	listDocumentsCmd.PersistentFlags().StringVar(&nameSearchArg, "name", "", "Search by name")
	listDocumentsCmd.PersistentFlags().StringVar(&tagSrchArg, "tag", "", "Search by tag")
	listDocumentsCmd.PersistentFlags().StringVar(&createdAfterSrcArg, "createdAfter", "",
			 "Documents created after date, in format '2019-03-26' ")
	listDocumentsCmd.PersistentFlags().StringVar(&createdBeforeSrchArg, "createdBefore", "",
			 "Documents created before date, in format '2019-03-26' ")
	listDocumentsCmd.PersistentFlags().StringVar(&modifiedAfterSrchArg, "modifiedAfter", "",
			 "Documents modified date, in format '2019-03-26' ")
	listDocumentsCmd.PersistentFlags().StringVar(&modifiedBeforeSrchArg, "modifiedBefore", "",
			 "Documents modified before date, in format '2019-03-26' ")
			 
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listDocumentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
