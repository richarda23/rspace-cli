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
	var docList *rspace.FileList
	var err error
	docList, err = ctx.WebClient.Files(cfg, mediaTypeArg)
	if err != nil {
		exitWithErr(err)
	}
	if ctx.Format.isJson() {
		ctx.write(prettyMarshal(docList))
	} else if ctx.Format.isQuiet() {
		processedResults := processResults(docList)
		printIds(ctx, toIdentifiableFile(processedResults))
	} else {
		processedResults := processResults(docList)
		listToFileTable(ctx, processedResults)
	}
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
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	listFilesCmd.PersistentFlags().StringVar(&mediaTypeArg, "mediaType", "", "Optional media type, 1 of 'image', 'document' or 'av'")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listDocumentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
