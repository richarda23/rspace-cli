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
)
var outputFormatArg string
var outFileArg string
var sortOrderArg string
var orderByArg string
var pageSizeArg int
// elnCmd represents the eln command
var elnCmd = &cobra.Command{
	Use:   "eln",
	Short: "Top-level command to work with RSpace ELN",
	Long: ` Run rspace eln --help to see all the possible commands.
	`,
	Args:   cobra.MinimumNArgs(1) ,
	Run: func(cmd *cobra.Command, args []string) {
		messageStdErr("Requires a subcommand")
	},
}

func init() {
	rootCmd.AddCommand(elnCmd)
	 elnCmd.PersistentFlags().StringVarP(&outputFormatArg, "outputFormat", "f", "table", "Output format: one of 'json','table', 'csv' or 'quiet' ")
	 elnCmd.PersistentFlags().StringVarP(&outFileArg, "outFile", "o", "", "Output file for program output")
}
