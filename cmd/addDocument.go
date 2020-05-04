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
	"io/ioutil"
	"strings"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

type addDocArgs struct {
	ParentfolderArg string
	NameArg         string
	Tags            string
	ContentFile     string
	Content         string
}

var addDocArgV = addDocArgs{}

// addDocumentCmd represents the createNotebook command
var addDocumentCmd = &cobra.Command{
	Use:   "addDocument",
	Short: "Creates a new basic document with optional tags and content",
	Long: `Create a new document, with an optional name and parent folder.

	  If a file is a file of HTML content, it is loaded verbatim, otherwise, plain text files are wrapped in '<pre>'
	  tags to preserve formatting.
	`,
	Example: `
	// create a doc tags and HTML content
	rspace eln  addDocument --name doc1 --tags tag1,tag2 --contentFile textToPutInDoc.html

	// create a doc with tags and  plain-text content, which will be wrapped in '<pre>' tag
	rspace eln  addDocument --name doc1 --tags tag1,tag2 --contentFile textToPutInDoc.txt

	// create a doc using verbatim text
	rspace eln  addDocument --name doc1  --content "some content"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()
		doAddDocRun(addDocArgV, context, context.WebClient)
	},
}

type DocClient interface {
	NewBasicDocumentWithContent(name, tags, content string) (*rspace.DocumentInfo, error)
}

func doAddDocRun(addDocArgV addDocArgs, context *Context, docClient DocClient) {
	content := getContent(addDocArgV)
	newDoc, err := docClient.NewBasicDocumentWithContent(addDocArgV.NameArg,
		addDocArgV.Tags, content)

	if err != nil {
		exitWithErr(err)
	}
	docList := rspace.DocumentList{}
	docList.Documents = []rspace.DocumentInfo{*newDoc}
	var dlf = DocListFormatter{&docList}
	context.writeResult(&dlf)
}

func getContent(addDocArgV addDocArgs) string {
	if len(addDocArgV.Content) > 0 {
		return addDocArgV.Content
	} else if len(addDocArgV.ContentFile) > 0 {
		bytes, err := ioutil.ReadFile(addDocArgV.ContentFile)
		if err != nil {
			exitWithErr(err)
		}
		lowerCaseFile := strings.ToLower(addDocArgV.ContentFile)
		if strings.HasSuffix(lowerCaseFile, "html") ||
			strings.HasSuffix(lowerCaseFile, "htm") {
			return string(bytes)
		} else {
			return "<pre>" + string(bytes) + "</pre>"
		}
	}
	return ""
}

func init() {
	elnCmd.AddCommand(addDocumentCmd)
	addDocumentCmd.Flags().StringVar(&addDocArgV.NameArg, "name", "", "A name for the document")
	addDocumentCmd.Flags().StringVar(&addDocArgV.ParentfolderArg, "folder", "", "An id for the folder that will contain the new document")
	addDocumentCmd.Flags().StringVar(&addDocArgV.Tags, "tags", "", "One or more tags, comma separated")
	addDocumentCmd.Flags().StringVar(&addDocArgV.ContentFile, "file", "", "A file of text or HTML content to put in the document")
	addDocumentCmd.Flags().StringVar(&addDocArgV.Content, "content", "", "Text or HTML content to put in the document")
}
