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

var Quiet bool
var PageSize int
var searchQuery string

// arguments for advanced search
var orSrchArg bool = false
var nameSearchArg string
var tagSrchArg string
var formSearchArg string
var createdBeforeSrchArg string
var createdAfterSrcArg string
var modifiedBeforeSrchArg string
var modifiedAfterSrchArg string

// listDocumentsCmd represents the listDocuments command
var listDocumentsCmd = &cobra.Command{
	Use:   "listDocuments",
	Short: "Lists the documents",
	Long: ` List or search for documents. Search term  is optional.
	`,
	Example: `
// A global search over names, tags, file and text content
rspace eln listDocuments --query myexperiment

// Get documents whose name starts with 'experiment' OR are created in 2020 or later
rspace eln listDocuments --or --name experiment* --createdAfter 2020-01-01

//  list the first hundred documents created
rspace eln listDocuments --orderBy created --sortOrder asc --maxResults 100

// list documents created from a particular form:
rspace eln listDocuments --form FM12345
	`,

	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContextWithTimeout(30)
		cfg := configurePagination()
		doListDocs(context, cfg)
	},
}

type DocListFormatter struct {
	*rspace.DocumentList
}

func (ds *DocListFormatter) ToJson() string {
	return prettyMarshal(ds.DocumentList)
}

func (ds *DocListFormatter) ToQuiet() []identifiable {
	rows := make([]identifiable, 0)
	for _, res := range ds.DocumentList.Documents {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}

func (ds *DocListFormatter) ToTable() *TableResult {
	results := ds.DocumentList
	baseResults := docListToBaseInfoList(results)
	maxNameColWidth := getMaxNameLength(baseResults)
	headers := []columnDef{columnDef{"GlobalId", 10}, columnDef{"Name", maxNameColWidth},
		columnDef{"Form", 10}, columnDef{"Created", 19}, columnDef{"Last Modified", 19}, columnDef{"Owner", 15}}

	rows := make([][]string, 0)
	for _, res := range results.Documents {

		data := []string{res.GlobalId, res.GetName(), res.Form.GlobalId,
			res.Created[0:DISPLAY_TIMESTAMP_WIDTH], res.LastModified[0:DISPLAY_TIMESTAMP_WIDTH], res.UserInfo.Username} // ignore seconds/millis to save space
		rows = append(rows, data)
	}
	return &TableResult{headers, rows}
}

func doListDocs(ctx *Context, cfg rspace.RecordListingConfig) {
	var docList *rspace.DocumentList
	var err error
	if len(searchQuery) > 0 {
		docList, err = ctx.WebClient.SearchDocuments(cfg, searchQuery)
	} else if advancedSrchArgsAreProvided() {
		docList, err = doAdvancedSearch(ctx, cfg)
	} else {
		docList, err = ctx.WebClient.Documents(cfg)
	}

	if err != nil {
		exitWithErr(err)
	}
	formatter := DocListFormatter{docList}
	ctx.writeResult(&formatter)
}
func advancedSrchArgsAreProvided() bool {
	var advSearchArgs = []string{nameSearchArg, tagSrchArg, createdBeforeSrchArg, createdAfterSrcArg,
		modifiedAfterSrchArg, modifiedBeforeSrchArg, formSearchArg}
	for _, v := range advSearchArgs {
		if len(v) > 0 {
			return true
		}
	}
	return false
}
func listToDocTable(ctx *Context, formatter ResultListFormatter) {
	table := formatter.ToTable()
	if ctx.Format.isCsv() {
		printCsv(ctx, table)
	} else {
		printTable(ctx, table)
	}
}

func docListToBaseInfoList(results *rspace.DocumentList) []rspace.BasicInfo {
	var baseResults = make([]rspace.BasicInfo, len(results.Documents))
	for i, v := range results.Documents {
		var x rspace.BasicInfo = v
		baseResults[i] = x
	}
	return baseResults
}

func doAdvancedSearch(ctx *Context, cfg rspace.RecordListingConfig) (*rspace.DocumentList, error) {
	messageStdErr("in advanced search")
	builder := rspace.SearchQueryBuilder{}
	if orSrchArg {
		builder.Operator(rspace.Or)
	}
	builder.AddTerm(nameSearchArg, rspace.NAME)
	builder.AddTerm(tagSrchArg, rspace.TAG)
	builder.AddTerm(formSearchArg, rspace.FORM)
	if createdTerm := createdAfterSrcArg + ";" + createdBeforeSrchArg; createdTerm != ";" {
		builder.AddTerm(createdTerm, rspace.CREATED)
	}
	if modifiedTerm := modifiedAfterSrchArg + ";" + modifiedBeforeSrchArg; modifiedTerm != ";" {
		builder.AddTerm(modifiedAfterSrchArg+";"+modifiedBeforeSrchArg, rspace.LAST_MODIFIED)
	}

	q := builder.Build()
	return ctx.WebClient.AdvancedSearchDocuments(cfg, q)
}

func init() {
	elnCmd.AddCommand(listDocumentsCmd)

	initPaginationFromArgs(listDocumentsCmd)

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
	listDocumentsCmd.PersistentFlags().StringVar(&formSearchArg, "form", "",
		"Documents created by a form; either name or globalID (e.g. FM5)")

}
