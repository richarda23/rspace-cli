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
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
 var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload one or more files",
	Long: ``,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := initialiseContext()
		uploadArgs(ctx, args)
	},
}

func uploadArgs (ctx *Context, args[]string ) {
	for _, filePath := range args {
		
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			exitWithErr(err)
		}

		if fileInfo.IsDir() {
			messageStdErr("Scanning for files in " +  fileInfo.Name())
			var filesInDir []string
			filepath.Walk(filePath, visit(&filesInDir) )
			messageStdErr(fmt.Sprintf("Found %d files to upload in %s",
				 len(filesInDir), fileInfo.Name()))
			for _, fileInDir := range filesInDir {
				postFile(ctx, fileInDir);
			}
		}
	}
 }
func visit (files *[]string) filepath.WalkFunc {
	return func  (path string, info os.FileInfo, err error) error {
		// don't recurse, ignore '.' files
		if !info.IsDir() && info.Name()[0] != '.' {
			fmt.Println(info.Name())
			*files = append(*files, path)	
		}
		return nil
	}
}
func postFile (ctx *Context, filePath string){
	file, err := ctx.WebClient.FileS.UploadFile(filePath)
	if err != nil {
		exitWithErr(err)
	}
	fmt.Println(prettyMarshal(file))
}
func init() {
	elnCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
