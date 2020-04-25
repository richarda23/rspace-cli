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
	"github.com/dustin/go-humanize"
	"regexp"
	"text/template"
	"bytes"
    "os/signal"
    "syscall"
)
type uploadCmdArgs struct {
 RecursiveFlag bool
 GenerateSummaryDoc bool
 DryrunFlag bool
 LogfileArg string
}

func setupInterrupt(ctx *Context, toUpload *[]*scannedFileInfo) chan bool {

    sigs := make(chan os.Signal, 1)
    done := make(chan bool, 1)

    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        sig := <-sigs
        fmt.Println()
		fmt.Println(sig)
		report(ctx, uploadedFiles)
		notUploadedYet := 0
		for _,v := range *toUpload {
			if !v.Uploaded {
				notUploadedYet++
			}
		}
		if notUploadedYet > 0 {
			logWriter := initLogWriter(uploadArgsArg.LogfileArg, os.Stderr)
			summary :=fmt.Sprintf("%d files weren't uploaded:", notUploadedYet)
			fmt.Fprintln(logWriter, summary)
			for _,v := range *toUpload {
				if ! v.Uploaded {
					fmt.Fprintln(logWriter, v.Path)
				}
			}
		}
        os.Exit(1)
    }()
    return done
}


type scannedFileInfo struct {
	Path string
	Info os.FileInfo
	Uploaded bool
}
var uploadedFiles = make ([]*rspace.FileInfo,0) 
var uploadArgsArg uploadCmdArgs
// uploadCmd represents the upload command
 var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload one or 'more files",
	Long: ` Upload files. Add files and folders to the command line. 
	By default, folder contents aren't uploaded recursively.
	
	Use the --recursive flag to upload all folder tree contents.
	
	The folder structure is flattened in RSpace. Files are uploaded to the target folder.
	
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
	
	var filesToUpload []*scannedFileInfo = make([]*scannedFileInfo,0)
	for _, filePath := range args {
		filePath, _ = filepath.Abs(filePath)
		fileInfo, _ := os.Stat(filePath)
		if fileInfo.IsDir() {
			messageStdErr("Scanning for files in " +  fileInfo.Name())
			if uploadArgsArg.RecursiveFlag {
				filepath.Walk(filePath, visit(&filesToUpload) )
			} else {
				readSingleDir(filePath, &filesToUpload)
			}
		} else {
			info,_:= os.Stat(filePath)
			filesToUpload = append(filesToUpload, &scannedFileInfo{filePath,info,false})
		}
	}
	messageStdErr(fmt.Sprintf("Found %d files to upload - total amount to upload is %s",len(filesToUpload),
			humanize.Bytes(sumFileSize(filesToUpload))))
	setupInterrupt(ctx, &filesToUpload)
	for _, fileToUpload := range filesToUpload {
		fileInfo :=	postFile(ctx, fileToUpload);
		if fileInfo != nil {
			uploadedFiles = append(uploadedFiles, fileInfo)
		}
	}
	report(ctx, uploadedFiles)
 }

 func sumFileSize(toUpload []*scannedFileInfo) uint64 {
	 var sum int64 = 0
	 for _,v :=range toUpload {
		 sum += v.Info.Size()
	 }
	 return uint64(sum)
 }

 func report(ctx *Context, uploaded []*rspace.FileInfo) {
	if uploadArgsArg.DryrunFlag {
		messageStdErr(fmt.Sprintf("File upload would upload %d files", len(uploaded)))
		return;
	}
	if uploadArgsArg.GenerateSummaryDoc {
		contentStr := generateSummaryContent(uploaded)
		summaryDocInfo:=ctx.WebClient.NewBasicDocumentWithContent("fileupload-summary","", contentStr)
		messageStdErr("Created summary with id " + summaryDocInfo.GlobalId)
	}
	messageStdErr(fmt.Sprintf("Reporting %d results:", len(uploaded)))

	var fal FileArrayList = FileArrayList{uploaded}
	var formatter FileListFormatter = FileListFormatter{fal}
	ctx.writeResult(&formatter)
 }
 // populates an HTML  table template with links to uploaded files
 func generateSummaryContent(results []*rspace.FileInfo) string {
	const tmpl = `
	<table>
	 <tr> <th>Name</th><th>Id</th><th>Link</th></tr>
		{{range $val := .}}
		 <tr><td>{{$val.Name}}</td><td><a href="/globalId/{{$val.GlobalId}}">{{$val.GlobalId}}</a></td><td><fileId={{$val.Id}}></td></tr>
		{{end}}
	</table>
	`
	t := template.Must(template.New("tmpl").Parse(tmpl))
	var buf bytes.Buffer

	t.Execute(&buf, results)
	return buf.String()
 }

 func fileListToBaseInfoList (results []*rspace.FileInfo) []rspace.BasicInfo {
	var baseResults = make([]rspace.BasicInfo, len(results))
	for i,v := range results {
		var x rspace.BasicInfo = v
		baseResults [i] = x
	}
	return baseResults
}
 // reads non . files from a single folder
 func readSingleDir(filePath string, files *[]*scannedFileInfo) {

	fileInfos,_:= ioutil.ReadDir(filePath)
	for _,inf:=range fileInfos {
		if !inf.IsDir() && !isDot(inf) {
				path := filePath + string(os.PathSeparator) +inf.Name()
				*files = append(*files, &scannedFileInfo{path,inf,false})
		}
	}
 }
func visit (files *[]*scannedFileInfo) filepath.WalkFunc {
	return func  (path string, info os.FileInfo, err error) error {
		// always   ignore '.' folders, don't descend
		messageStdErr("processing " + path)
		if info.IsDir() && isDot(info) {
			messageStdErr("Skipping .folder " + path)
			return filepath.SkipDir
		}
		// always add non . files
		if !info.IsDir() && !isDot(info) {
			*files = append(*files, &scannedFileInfo{path,info,false})
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
func postFile (ctx *Context, fileInfo *scannedFileInfo) *rspace.FileInfo {
	filePath := fileInfo.Path
	if uploadArgsArg.DryrunFlag {
		return &rspace.FileInfo{}
	}
	messageStdErr("Uploading: " + filePath)
	file, err := ctx.WebClient.UploadFile(filePath)
	if err != nil {
		// other files might upload OK, so don't exit here
		messageStdErr(err.Error())
	}
	fileInfo.Uploaded = true
	return file
}
func init() {
	elnCmd.AddCommand(uploadCmd)
	uploadCmd.PersistentFlags().BoolVar(&uploadArgsArg.RecursiveFlag, "recursive", false,	"If uploading a folder, uploads contents recursively.")
	uploadCmd.PersistentFlags().BoolVar(&uploadArgsArg.DryrunFlag, "dry-run", false,"Performs a dry-run, reports on what would be uploaded")
	uploadCmd.PersistentFlags().StringVar(&uploadArgsArg.LogfileArg, "logfile", "","A log file to record upload progress, if not set will log to standard error")
	uploadCmd.PersistentFlags().BoolVar(&uploadArgsArg.GenerateSummaryDoc,
		 "add-summary", false, "Generate a summary document containing links to uploaded files")
}
