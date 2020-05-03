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
	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
)

func NewStatusCmd() *cobra.Command {
	// statusCmd represents the status command
	return &cobra.Command{
		Use:     "status",
		Short:   "Checks status of RSpace",
		Long:    "Gets version and current status of RSpace",
		Example: "rspace status",
		RunE:    runFunction,
	}
}

type StatusCli interface {
	Status() (*rspace.Status, error)
}

//fixed signature for cobra framework
func runFunction(cmd *cobra.Command, args []string) error {
	context := initialiseContext()
	return doRun(context.WebClient, context)
}

func doRun(cli StatusCli, ctx *Context) error {
	got, err := cli.Status()
	if err != nil {
		ctx.messageStdErr(err.Error())
		return err
	}
	ctx.write(got.RSpaceVersion + ", " + got.Message + "\n")
	return nil
}

func init() {
	rootCmd.AddCommand(NewStatusCmd())
}
