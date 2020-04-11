package cmd
import (
"os"
"encoding/json"
"fmt"
"strings"
//"errors"
)
func exitWithStdErrMsg (message string) {
	messageStdErr (message)
	os.Exit(1)
}
func exitWithErr (err error) {
	exitWithStdErrMsg(err.Error())
}
func messageStdErr (message string) {
	fmt.Fprintln(os.Stderr, message)
}
func prettyMarshal(anything interface{}) string {
        bytes, _ := json.MarshalIndent(anything, "", "\t")
        return string(bytes)
}
type columnDef struct {
	Title string
	Width int
}
func printTable(headers []columnDef , content [][]string ){
	printTableHeaders (headers)
	printContent(headers, content)
}
func printContent (headers []columnDef, content [][]string) {
	for _, row := range content {
		rowToPrint := make([]string, 0)
		for i, cell := range row {
			toPrint:= abbreviate(cell, headers[i].Width)
			toPrint = fmt.Sprintf("%-[1]*s", headers[i].Width, toPrint)
			rowToPrint = append(rowToPrint, toPrint)
		}
		fmt.Println(strings.Join(rowToPrint, "\t"))
	}
}
func printTableHeaders (headers []columnDef) {
	headersToPrint := make([]string, 0)
	for _, header := range headers {
		toPrint:= abbreviate(header.Title, header.Width)
		toPrint = fmt.Sprintf("%-[1]*s", header.Width, toPrint)
		headersToPrint = append(headersToPrint, toPrint)
	}
	fmt.Println(strings.Join(headersToPrint, "\t"))
}
func abbreviate(toAbbreviate string, maxLen int) string {
        if len(toAbbreviate) > maxLen {
                toAbbreviate = toAbbreviate[0:(maxLen-3)] + ".."
        }
        return toAbbreviate
}


