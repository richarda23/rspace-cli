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
	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type importWordCmdArgs struct {
	RecursiveFlag      bool
	GenerateSummaryDoc bool
	DryrunFlag         bool
	LogfileArg         string
	TargetFolder       int
}

func setUpImportInterrupt(ctx *Context, toUpload *[]*scannedFileInfo) chan bool {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		messageStdErr(sig.String())
		reportImport(ctx, importedDocs)
		notUploadedYet := 0
		for _, v := range *toUpload {
			if !v.Uploaded {
				notUploadedYet++
			}
		}
		if notUploadedYet > 0 {
			logWriter := initLogWriter(importArgsArg.LogfileArg, os.Stderr)
			summary := fmt.Sprintf("%d files weren't uploaded:", notUploadedYet)
			fmt.Fprintln(logWriter, summary)
			for _, v := range *toUpload {
				if !v.Uploaded {
					fmt.Fprintln(logWriter, v.Path)
				}
			}
		}
		os.Exit(1)
	}()
	return done
}

var importedDocs = make([]*rspace.DocumentInfo, 0)
var importArgsArg importWordCmdArgs

// importWordCmd represents the upload command
var importWordCmd = &cobra.Command{
	Use:   "importWord",
	Short: "Import MSOffice Word files (doc, docx or odt)",
	Long: ` Import Word files as RSpace document. Add files and folders to the command line. 
	Files and folders are scanned for Word documents and converter
	
	Use the --recursive flag to scan all folder tree contents.
	
	Any folder structure in the input is flattened in RSpace. 
	Documents are generated in 'folder' or HomeFolder if 'targetFolder' is not set.
	
	Files or folder names starting with '.' are ignored. But you can use '.' as an argument
	to scan the current folder, e.g.

	rspace eln importWord . --recursive

	If you are importing many files, and cancel the operation while it is still running by a Ctrl-C
	or other interrupt signal, the files *not* imported will be listed in stderr or in a file
	specified by the --logfile argument.
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := initialiseContext()
		importArgs(ctx, args)
	},
}

func importArgs(ctx *Context, args []string) {
	// fail fast if files can't be read
	validateInputFilePaths(args)
	filesToUpload := scanFiles(args, importArgsArg.RecursiveFlag, acceptMsDoc())

	messageStdErr(fmt.Sprintf("Found %d files to import - total amount to import is %s", len(filesToUpload),
		sumFileSizeHuman(filesToUpload)))
	setUpImportInterrupt(ctx, &filesToUpload)
	for _, fileToUpload := range filesToUpload {
		docInfo := importFile(ctx, fileToUpload)
		if docInfo != nil {
			importedDocs = append(importedDocs, docInfo)
		}
	}
	reportImport(ctx, importedDocs)
}

type DocArrayList struct {
	docList []*rspace.DocumentInfo
}

type DocArrayListFormatter struct {
	DocArrayList
}

func (fmt *DocArrayListFormatter) ToJson() string {
	return prettyMarshal(fmt.DocArrayList.docList)
}

func (fmt *DocArrayListFormatter) ToQuiet() []identifiable {
	return toIdentifiableDoc(fmt.DocArrayList.docList)
}

func (fmt *DocArrayListFormatter) ToTable() *TableResult {
	results := fmt.DocArrayList.docList

	//baseInfos := fileListToBaseInfoList(results)
	///nameColWidth := getMaxNameLength(baseInfos)
	headers := []columnDef{columnDef{"Id", 8}, columnDef{"GlobalId", 10},
		columnDef{"Name", 25},
		columnDef{"Created", DISPLAY_TIMESTAMP_WIDTH}}

	rows := make([][]string, 0)
	for _, res := range results {
		data := []string{strconv.Itoa(res.Id), res.GlobalId, res.Name,
			res.Created[0:DISPLAY_TIMESTAMP_WIDTH]}
		rows = append(rows, data)
	}
	return &TableResult{headers, rows}
}

func toIdentifiableDoc(results []*rspace.DocumentInfo) []identifiable {
	rows := make([]identifiable, 0)
	for _, res := range results {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}

func reportImport(ctx *Context, uploaded []*rspace.DocumentInfo) {
	if importArgsArg.DryrunFlag {
		messageStdErr(fmt.Sprintf("File upload would upload %d files", len(uploaded)))
		return
	}
	messageStdErr(fmt.Sprintf("reportImporting %d results:", len(uploaded)))

	var dal DocArrayList = DocArrayList{uploaded}
	//TODO FIX THIS, implement DocListFormatter
	var formatter DocArrayListFormatter = DocArrayListFormatter{dal}
	ctx.writeResult(&formatter)
}

func importDocListToBaseInfoList(results []*rspace.FileInfo) []rspace.BasicInfo {
	var baseResults = make([]rspace.BasicInfo, len(results))
	for i, v := range results {
		var x rspace.BasicInfo = v
		baseResults[i] = x
	}
	return baseResults
}

func importFile(ctx *Context, fileInfo *scannedFileInfo) *rspace.DocumentInfo {
	filePath := fileInfo.Path
	if importArgsArg.DryrunFlag {
		return &rspace.DocumentInfo{}
	}
	messageStdErr("Uploading: " + filePath)
	doc, err := ctx.WebClient.ImportWord(filePath, importArgsArg.TargetFolder, 0)
	if err != nil {
		// other files might upload OK, so don't exit here
		messageStdErr(err.Error())
	} else {
		fileInfo.Uploaded = true
	}
	return doc
}
func init() {
	elnCmd.AddCommand(importWordCmd)
	importWordCmd.PersistentFlags().BoolVar(&importArgsArg.RecursiveFlag, "recursive", false, "If uploading a folder, uploads contents recursively.")
	importWordCmd.PersistentFlags().BoolVar(&importArgsArg.DryrunFlag, "dry-run", false, "Performs a dry-run, reportImports on what would be uploaded")
	importWordCmd.PersistentFlags().StringVar(&importArgsArg.LogfileArg, "logfile", "", "A log file to record upload progress, if not set will log to standard error")
	importWordCmd.PersistentFlags().IntVar(&importArgsArg.TargetFolder,
		"folder", 0, "ID of Target folder for imported Word files")
}
