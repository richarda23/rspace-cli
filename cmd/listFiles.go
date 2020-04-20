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
var mediaTypeArg = ""
// listDocumentsCmd represents the listDocuments command
var listFilesCmd = &cobra.Command{
	Use:   "listFiles",
	Short: "Lists attachment files.",
	Long:`List files, with optional 'mediaType' argument to restrict the type of files retrieved 

		  rspace eln listFiles --mediaType document
		  rspace eln listFiles --mediaType image
		  rspace eln listFiles --mediaType av
	`,

	Run: func(cmd *cobra.Command, args []string) {
		messageStdErr("listFiles called:")
		context := initialiseContext()  
		cfg := configurePagination()
		doListFiles(context, cfg)
	},
}

func doListFiles (ctx *Context, cfg rspace.RecordListingConfig) {
	var filesList *rspace.FileList
	var err error
	filesList, err = ctx.WebClient.Files(cfg, mediaTypeArg)
	if err != nil {
		exitWithErr(err)
	}
	formatter := &FileListFormatter{FileArrayList{processResults(filesList)}}
	ctx.writeResult(formatter)
}

type FileArrayList struct {
	fList []*rspace.FileInfo
}
type FileListFormatter struct {
	FileArrayList
}

func (fs *FileListFormatter) ToJson () string{
	return prettyMarshal(fs.FileArrayList.fList)
}

func (ds *FileListFormatter) ToQuiet () []identifiable{
	return toIdentifiableFile(ds.FileArrayList.fList)
}

func (ds *FileListFormatter) ToTable () *TableResult {
	results := ds.FileArrayList.fList

	baseInfos := fileListToBaseInfoList(results)
	nameColWidth := getMaxNameLength(baseInfos)
	headers := []columnDef {columnDef{"Id",8}, columnDef{"GlobalId",10},  columnDef{"Name", nameColWidth}, 
	 columnDef{"Created",DISPLAY_TIMESTAMP_WIDTH},columnDef{"Size",12},columnDef{"ContentType", 25}}

	rows := make([][]string, 0)
	for _, res := range results {
		data := []string {strconv.Itoa(res.Id),res.GlobalId, res.Name,
			   res.Created[0:DISPLAY_TIMESTAMP_WIDTH],strconv.Itoa(res.Size),res.ContentType}
		rows = append(rows, data)
	}
	return &TableResult{headers, rows}
	
 }
 func toIdentifiableFile (results []*rspace.FileInfo) []identifiable {
	rows := make([]identifiable, 0)
	
	for _, res := range results {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}
// convert results  so can re-used file-display methods
func processResults (files *rspace.FileList) []*rspace.FileInfo {
	rc := make ([]*rspace.FileInfo, len(files.Files))
	for i, v := range files.Files {
		f := v
		rc[i] = &f
	}
	return rc
}

func init() {
	elnCmd.AddCommand(listFilesCmd)

	initPaginationFromArgs(listFilesCmd)
	// and all subcommands, e.g.:
	listFilesCmd.PersistentFlags().StringVar(&mediaTypeArg, "mediaType", "", "Optional media type, 1 of 'image', 'document' or 'av'")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listDocumentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
