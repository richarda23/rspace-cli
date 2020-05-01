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
	"github.com/richarda23/rspace-client-go/rspace"
	"strconv"

)
// listDocumentsCmd represents the listDocuments command
var listFormsCmd = &cobra.Command{
	Use:   "listForms",
	Short: "Lists forms",
	Long:`List forms, sorted or paginated, e.g.

		  rspace eln listForms --orderBy name --maxResults 100
	`,

	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()  
		cfg := configurePagination()
		doListForms(context, cfg)
	},
}

func doListForms (ctx *Context, cfg rspace.RecordListingConfig) {
	var formsList *rspace.FormList
	var err error
	formsList, err = ctx.WebClient.Forms(cfg)
	if err != nil {
		exitWithErr(err)
	}
	formatter := &FormListFormatter{formsList}
	ctx.writeResult(formatter)
}

type FormListFormatter struct {
	*rspace.FormList
}

func (fs *FormListFormatter) ToJson () string{
	return prettyMarshal(fs.FormList)
}

func (ds *FormListFormatter) ToQuiet () []identifiable{
	rows := make([]identifiable, 0)
	for _, res := range ds.FormList.Forms {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}

func (ds *FormListFormatter) ToTable () *TableResult {
	results := ds.FormList.Forms

	headers := []columnDef {columnDef{"Id",8}, columnDef{"GlobalId",10}, columnDef{"Name", 25}, 
	 columnDef{"StableId", 25}}

	rows := make([][]string, 0)
	for _, res := range results {
		data := []string {strconv.Itoa(res.Id),res.GlobalId, res.Name,
			   res.StableId}
		rows = append(rows, data)
	}
	return &TableResult{headers, rows}
	
 }
 func toIdentifiableForm (results []*rspace.FormInfo) []identifiable {
	rows := make([]identifiable, 0)
	
	for _, res := range results {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}
func init() {
	elnCmd.AddCommand(listFormsCmd)
	initPaginationFromArgs(listFormsCmd)
}