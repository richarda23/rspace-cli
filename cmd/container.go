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
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

type containerArgs struct {
	yamlFile string
}

var containerArgsa containerArgs

// listDocumentsCmd represents the listDocuments command
var containerCmd = &cobra.Command{
	Use:   "container",
	Short: "Creates a container hierarchy",
	Long: `Create container hierarchy from yaml definition files
	`,
	Example: `
rspace eln container --definitionFile myContainers.yaml
	`,

	Run: func(cmd *cobra.Command, args []string) {
		context := initialiseContext()
		createTree(context, args, &containerArgsa)
	},
}

type AllContainers struct {
	AllContainers []*rspace.ContainerPost `yaml:"allcontainers"`
}

func createTree(ctx *Context, args []string, flags *containerArgs) {
	filePath := containerArgsa.yamlFile
	conList := &AllContainers{}

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		exitWithErr(err)
	}
	err = yaml.Unmarshal([]byte(data), conList)
	if err != nil {
		return
	}
	for _, v := range conList.AllContainers {
		con := v
		con.ParentContainer = nil

		results, _ := ctx.WebClient.CreateContainerTree(con)
		for _, v := range results.Containers {
			fmt.Printf("Created container %s with id %d \n", v.Name, v.Id)
		}
	}
}

func init() {
	elnCmd.AddCommand(containerCmd)
	containerCmd.Flags().StringVar(&containerArgsa.yamlFile, "definitionFile", "",
		"Yaml File defining container tree")
}
