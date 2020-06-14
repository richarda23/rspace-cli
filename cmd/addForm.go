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
	"strings"

	"github.com/pkg/errors"
	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
	//"fmt"
)

type addFormArgs struct {
	Publish bool
}

var formArgs = addFormArgs{}

// addUserCmd represents the createNotebook command
var addFormCmd = &cobra.Command{
	Use:   "addForm",
	Short: "Adds a new Form",
	Long: `
Create 1 or forms from .yaml or .json definition files.
The synatx for the form definitions should be that described in the forms/
POST API documentation, or its YAML equivalent.

Multiple files can be specified as arguments.

--publish flag, if set, will publish the Forms so they can be used 
immediately to create documents. If not set, forms remain in the 'NEW' state.
`,
	Example: ` 
// Create 2 forms; one specified in yaml format and one specified in json format.
// Publish both forms so they are available to use to create documents.
rspace eln addForm myFormDef.yaml mySecondForm.json --publish
`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := initialiseContext()
		doCreateForms(args, ctx)
	},
}

func doCreateForms(args []string, ctx *Context) {
	createdForms := make([]rspace.Form, 0)
	// iterate over arguments
	for _, v := range args {
		form, err := createForm(v, ctx)
		// file not exist, not a valid file, or create failed
		if err != nil {
			messageStdErr(fmt.Sprintf("%s -skipping", err.Error()))
			continue
		}
		createdForms = append(createdForms, *form)
		if formArgs.Publish {
			form, err = ctx.WebClient.PublishForm(form.Id)
			if err != nil {
				messageStdErr(err.Error())
				continue
			} else {
				// update with new publishing status
				createdForms[len(createdForms)-1] = *form
			}
		}
	}
	// now print results
	if len(createdForms) > 0 {
		fList := rspace.FormList{}
		fList.Forms = createdForms
		formatter := &FormListFormatter{&fList}
		ctx.writeResult(formatter)
	} else {
		messageStdErr("No forms created")
	}
}

func createForm(path string, ctx *Context) (*rspace.Form, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var formInfo *rspace.Form
	if strings.ToLower(filepath.Ext(path)) == ".yaml" {
		formInfo, err = ctx.WebClient.CreateFormYaml(file)
	} else if strings.ToLower(filepath.Ext(path)) == ".json" {
		formInfo, err = ctx.WebClient.CreateFormJson(file)
	} else {
		return nil, errors.Errorf(
			"'%s' is not a yaml or json file", path)
	}
	if err != nil {
		return nil, err
	} else {
		return formInfo, nil
	}
}

func init() {
	elnCmd.AddCommand(addFormCmd)
	addFormCmd.Flags().BoolVar(&formArgs.Publish, "publish", false,
		"Publishes the forms")

}
