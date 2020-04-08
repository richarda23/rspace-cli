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
	"github.com/spf13/cobra"
)
var Quiet bool
var PageSize int

// listDocumentsCmd represents the listDocuments command
var listDocumentsCmd = &cobra.Command{
	Use:   "listDocuments",
	Short: "Lists the documents",
	Long:` Lists documents `,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listDocuments called")
		client := setup()
		fmt.Println(client.DocumentS.GetStatus())
	},
}

func init() {
	listDocumentsCmd.PersistentFlags().BoolVarP(&Quiet, "quiet", "q", false, "quiet output")
	listDocumentsCmd.Flags().IntVarP(&PageSize, "page-size", "p", 20, "page size")
	rootCmd.AddCommand(listDocumentsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listDocumentsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listDocumentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
