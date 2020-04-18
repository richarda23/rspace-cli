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
	"rspace"
	"os"
	"path/filepath"
	"io/ioutil"
	"github.com/spf13/cobra"
	"regexp"
	"strconv"
)
var recursiveFlag bool = false
// uploadCmd represents the upload command
 var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload one or 'more files",
	Long: ` Uepload files. Add files and folders to the command line. 
	By default, folder contents aren't uploaded recursively.
	
	Use the --recursive flag to upload all folder contents.
	
	The folder structure is flattened in RSpace, files are uploaded to the target folder.
	
	If not set, files will be uploaded to the appropriate 'Api Inbox' Gallery folders,
	depending on the file type. 

	Files or folder names starting with '.' are ignored. But you can use '.' as an argument
	to upload the current folder, e.g.

	    rspace eln upload . --recursive
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := initialiseContext()
		uploadArgs(ctx, args)
	},
}

func validateArguments( args []string ) {
	for _, filePath := range args {
		var err error
		filePath, err = filepath.Abs(filePath)
		_, err = os.Stat(filePath)
		if err != nil {
			exitWithErr(err)
		}
	}
}

func uploadArgs (ctx *Context, args[]string ) {
	// fail fast if files can't be read
	validateArguments(args)
	
	var uploadedFiles = make ([]*rspace.FileInfo,0) 
	for _, filePath := range args {
		filePath, _ = filepath.Abs(filePath)
		fileInfo, _ := os.Stat(filePath)
		if fileInfo.IsDir() {
			messageStdErr("Scanning for files in " +  fileInfo.Name())
			var filesInDir []string
			if recursiveFlag {
				filepath.Walk(filePath, visit(&filesInDir) )
			} else {
				readSingleDir(filePath, &filesInDir)
			}
			messageStdErr(fmt.Sprintf("Found %d files to upload in %s",
				len(filesInDir), fileInfo.Name()))
			for _, fileInDir := range filesInDir {
				fileInfo :=	postFile(ctx, fileInDir);
				if fileInfo != nil {
					uploadedFiles = append(uploadedFiles, fileInfo)
				}
			}
		} else {
			fileInfo :=	postFile(ctx, filePath);
			if fileInfo != nil {
				uploadedFiles = append(uploadedFiles, fileInfo)
			}
		}
	}
	report(ctx, uploadedFiles)
 }

 func report(ctx *Context, uploaded []*rspace.FileInfo) {
	messageStdErr(fmt.Sprintf("Reporting %d results:", len(uploaded)))
	if ctx.Format.isJson() {
		ctx.write(prettyMarshal(uploaded))
	} else if ctx.Format.isQuiet() {
		printIds(ctx, toIdentifiableFile(uploaded))
	} else {
		listToFileTable(ctx, uploaded)
	}
 }

 func listToFileTable (ctx *Context, results []*rspace.FileInfo) {
	headers := []columnDef {columnDef{"Id",8}, columnDef{"GlobalId",10},  columnDef{"Name", 25}, 
	 columnDef{"Created",24},columnDef{"Size",12},columnDef{"ContentType", 25}}

	rows := make([][]string, 0)
	for _, res := range results {
		data := []string {strconv.Itoa(res.Id),res.GlobalId, res.Name,
			   res.Created,strconv.Itoa(res.Size),res.ContentType}
		rows = append(rows, data)
	}
	if ctx.Format.isCsv() {
		printCsv(ctx, headers, rows)
	} else {
		printTable(ctx, headers, rows)
	}
 }
 func toIdentifiableFile (results []*rspace.FileInfo) []identifiable {
	rows := make([]identifiable, 0)
	
	for _, res := range results {
		rows = append(rows, identifiable{strconv.Itoa(res.Id)})
	}
	return rows
}
 // reads non . files from a single folder
 func readSingleDir(filePath string, files *[]string) {

	fileInfos,_:= ioutil.ReadDir(filePath)
	for _,inf:=range fileInfos {
		if !inf.IsDir() && !isDot(inf) {
				*files = append(*files, filePath + string(os.PathSeparator) +inf.Name())
		}
	}
 }
func visit (files *[]string) filepath.WalkFunc {
	return func  (path string, info os.FileInfo, err error) error {
		// always   ignore '.' folders, don't descend
		messageStdErr("processing " + path)
		if info.IsDir() && isDot(info) {
			messageStdErr("Skipping .folder " + path)
			return filepath.SkipDir
		}
		// always add non . files
		if !info.IsDir() && !isDot(info) {
			*files = append(*files, path)
			return nil
		}
		return nil
	}
}
func isDot(info os.FileInfo) bool {
	//return filepath.Base(info.Name())[0] == '.'
	match,_ :=  regexp.MatchString("^\\.[A-Za-z0-9\\-_]+", info.Name())
	return match
}
func postFile (ctx *Context, filePath string) *rspace.FileInfo {
	messageStdErr("Uploading: " + filePath)
	file, err := ctx.WebClient.UploadFile(filePath)
	if err != nil {
		// other files might upload OK, so don't exit here
		messageStdErr(err.Error())
	}
	return file
}
func init() {
	elnCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	uploadCmd.PersistentFlags().BoolVar(&recursiveFlag, "recursive", false,
	 "If uploading a folder, uploads contents recursively.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
