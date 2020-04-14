package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	//"errors"
)

func exitWithStdErrMsg(message string) {
	messageStdErr(message)
	os.Exit(1)
}
func exitWithErr(err error) {
	exitWithStdErrMsg(err.Error())
}
func messageStdErr(message string) {
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

func printTable(ctx *Context, headers []columnDef, content [][]string) {
	printTableHeaders(ctx, headers)
	printContent(ctx, headers, content)
}
func printCsv(ctx *Context, headers []columnDef, content [][]string) {
	writer := csv.NewWriter(ctx.Writer)
	writer.Write(columnDefsToString(headers))
	if err := writer.WriteAll(content); err != nil {
		exitWithErr(err)
	}
}

func columnDefsToString(headers []columnDef) []string {
	rowToPrint := make([]string, 0)
	for _, header := range headers {
		rowToPrint = append(rowToPrint, header.Title)
	}
	return rowToPrint
}

func printContent(ctx *Context, headers []columnDef, content [][]string) {
	for _, row := range content {
		rowToPrint := make([]string, 0)
		for i, cell := range row {
			toPrint := abbreviate(cell, headers[i].Width)
			toPrint = fmt.Sprintf("%-[1]*s", headers[i].Width, toPrint)
			rowToPrint = append(rowToPrint, toPrint)
		}
		ctx.write(strings.Join(rowToPrint, "\t"))
	}
}
func printTableHeaders(ctx *Context, headers []columnDef) {
	headersToPrint := make([]string, 0)
	for _, header := range headers {
		toPrint := abbreviate(header.Title, header.Width)
		toPrint = fmt.Sprintf("%-[1]*s", header.Width, toPrint)
		headersToPrint = append(headersToPrint, toPrint)
	}
	ctx.write(strings.Join(headersToPrint, "\t"))
}
func abbreviate(toAbbreviate string, maxLen int) string {
	if len(toAbbreviate) > maxLen {
		toAbbreviate = toAbbreviate[0:(maxLen-3)] + ".."
	}
	return toAbbreviate
}

type identifiable struct {
	Id string
}

//
func printIds(ctx *Context, source []identifiable) {
	for _, item := range source {
		ctx.write(item.Id)
	}
}
