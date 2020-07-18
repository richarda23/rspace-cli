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
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

type uploadCmdArgs struct {
	RecursiveFlag      bool
	GenerateSummaryDoc bool
	DryrunFlag         bool
	LogfileArg         string
	Caption            string
	TemplateFile       string
}

func setupInterrupt(ctx *Context, toUpload *[]*scannedFileInfo) chan bool {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		messageStdErr(sig.String())
		report(ctx, uploadedFiles)
		notUploadedYet := 0
		for _, v := range *toUpload {
			if !v.Uploaded {
				notUploadedYet++
			}
		}
		if notUploadedYet > 0 {
			logWriter := initLogWriter(uploadArgsArg.LogfileArg, os.Stderr)
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

var uploadedFiles = make([]*rspace.FileInfo, 0)
var uploadArgsArg uploadCmdArgs

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload one or more files or folders",
	Long: ` Upload files. Add files and folders to the command line. 
By default, folder contents aren't uploaded recursively.

Use the --recursive flag to upload all folder tree contents.

Any folder structure in the input is flattened in RSpace. Files are uploaded to the target folder.

If an explicit folder is not set, files will be uploaded to the appropriate 'Api Inbox' Gallery folders,
depending on the file type. 

Files or folder names starting with '.' are ignored. But you can use '.' as an argument
to upload the current folder

The --add-summary flag creates a document in your Home Folder with links to all 
the uploaded files, as a reference to the uploaded files.

If you are uploading many files, and cancel the operation while it is still running by a Ctrl-C
or other interrupt signal, the files *not* uploaded will be listed in stderr or in a file
specified by the --logfile argument.
	`,
	Example: `

// upload a single file
rspace eln upload myimage.png

// do a dry run to see what would be uploaded
rspace eln upload file.doc imageFolder --recursive --dry-run

//use a logfile to record what was uploaded, in the event of cancellation or error
rspace eln upload folderWithManyFiles --recursive --logfile progress.txt

// upload a file and a folder, recursively, and generate a summary document
// A caption will be added to all uploaded files - useful for tagging collections of files.
rspace eln upload file.doc imageFolder --recursive --add-summary --caption anti-CDC2-immunofluorescence
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := initialiseContext()
		uploadArgs(ctx, args)
	},
}

func uploadArgs(ctx *Context, args []string) {
	// fail fast if files can't be read
	validateInputFilePaths(args)
	filesToUpload := scanFiles(args, uploadArgsArg.RecursiveFlag, acceptAll())

	messageStdErr(fmt.Sprintf("Found %d files to upload - total amount to upload is %s", len(filesToUpload),
		sumFileSizeHuman(filesToUpload)))
	setupInterrupt(ctx, &filesToUpload)
	for _, fileToUpload := range filesToUpload {
		fileInfo := postFile(ctx, fileToUpload)
		if fileInfo != nil {
			uploadedFiles = append(uploadedFiles, fileInfo)
		}
	}
	report(ctx, uploadedFiles)
}

func report(ctx *Context, uploaded []*rspace.FileInfo) {
	if uploadArgsArg.DryrunFlag {
		messageStdErr(fmt.Sprintf("File upload would upload %d files", len(uploaded)))
		return
	}
	if uploadArgsArg.GenerateSummaryDoc {
		addSummaryDoc(ctx, uploaded)
	}
	messageStdErr(fmt.Sprintf("Reporting %d results:", len(uploaded)))

	var fal FileArrayList = FileArrayList{uploaded}
	var formatter FileListFormatter = FileListFormatter{fal}
	ctx.writeResult(&formatter)
}

func addSummaryDoc(ctx *Context, uploaded []*rspace.FileInfo) {
	contentStr, _ := generateSummaryContent(uploaded)
	messageStdErr(contentStr)
	summaryDocInfo, err := ctx.WebClient.NewBasicDocumentWithContent("fileupload-summary", "", contentStr)
	if err != nil {
		messageStdErr(err.Error())
	} else {
		messageStdErr("Created summary with id " + summaryDocInfo.GlobalId)
	}
}

type FileInfoSummary struct {
	*rspace.FileInfo
}

func (summary *FileInfoSummary) FileIdLink() template.HTML {
	return template.HTML(fmt.Sprintf("<fileId=%d>", summary.Id))
}

func (summary *FileInfoSummary) GlobalIdLink() template.HTML {
	return template.HTML(fmt.Sprintf(`<a href="/globalId/%s">`, summary.GlobalId))
}

// populates an HTML  table template with links to uploaded files
func generateSummaryContent(results2 []*rspace.FileInfo) (string, error) {

	results := make([]*FileInfoSummary, 0)
	for _, v := range results2 {
		f := v
		results = append(results, &FileInfoSummary{f})
	}
	const tmpl = `
	<table>
	 <tr> <th>Name</th><th>Id</th><th>Link</th></tr>
		{{range $val := .}}
		 <tr>
		 <td>{{$val.Name}}</td>
		 <td>{{$val.GlobalIdLink}}{{$val.GlobalId}}</a></td>
		 <td>
		 {{$val.FileIdLink}}
		 </td>
		 </tr>
		{{end}}
	</table>
	`
	var templToUse string = tmpl

	if len(uploadArgsArg.TemplateFile) > 0 {
		bytes, err := ioutil.ReadFile(uploadArgsArg.TemplateFile)
		if err != nil {
			messageStdErr(err.Error())
			return "", err
		}
		templToUse = string(bytes)
	}
	t := template.Must(template.New("tmpl").Parse(templToUse))
	var buf bytes.Buffer
	t.Execute(&buf, results)
	return buf.String(), nil
}

func fileListToBaseInfoList(results []*rspace.FileInfo) []rspace.BasicInfo {
	var baseResults = make([]rspace.BasicInfo, len(results))
	for i, v := range results {
		var x rspace.BasicInfo = v
		baseResults[i] = x
	}
	return baseResults
}

func postFile(ctx *Context, fileInfo *scannedFileInfo) *rspace.FileInfo {
	filePath := fileInfo.Path
	if uploadArgsArg.DryrunFlag {
		return &rspace.FileInfo{}
	}
	messageStdErr("Uploading: " + filePath)
	cfg := rspace.FileUploadConfig{}
	cfg.Caption = uploadArgsArg.Caption
	cfg.FilePath = filePath
	file, err := ctx.WebClient.UploadFile(cfg)
	if err != nil {
		// other files might upload OK, so don't exit here
		messageStdErr(err.Error())
	}
	fileInfo.Uploaded = true
	return file
}
func init() {
	elnCmd.AddCommand(uploadCmd)
	uploadCmd.PersistentFlags().BoolVar(&uploadArgsArg.RecursiveFlag, "recursive", false, "If uploading a folder, uploads contents recursively.")
	uploadCmd.PersistentFlags().BoolVar(&uploadArgsArg.DryrunFlag, "dry-run", false, "Performs a dry-run, reports on what would be uploaded")
	uploadCmd.PersistentFlags().StringVar(&uploadArgsArg.LogfileArg, "logfile", "", "A log file to record upload progress, if not set will log to standard error")
	uploadCmd.PersistentFlags().StringVar(&uploadArgsArg.Caption, "caption", "", "A caption to be added to all uploaded files")
	uploadCmd.PersistentFlags().BoolVar(&uploadArgsArg.GenerateSummaryDoc,
		"add-summary", false, "Generate a summary document containing links to uploaded files")
	uploadCmd.PersistentFlags().StringVar(&uploadArgsArg.TemplateFile, "summary-template", "", "Template for summary document")

}
