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
	"rspace"
	"strings"
	"time"
	"fmt"
	"strconv"
)

 type acArgsS struct {
	actionsArg string
	usersArg string
	afterDateArg string
	beforeDateArg string
	globalId string
}

func (acargs *acArgsS) actions() []string {
	return _splitAndTrim(acargs.actionsArg)
}

func (acargs *acArgsS) users() []string {
	return _splitAndTrim(acargs.usersArg)
}

func (acargs *acArgsS) after()  (time.Time){
	return parseDateArg(acargs.afterDateArg)
}

func (acargs *acArgsS) before()  (time.Time){
	return parseDateArg(acargs.beforeDateArg)
}

func (acargs *acArgsS) oid()  rspace.GlobalId{
	return rspace.GlobalId(acargs.globalId)
}
 

func parseDateArg (dateArg string) (time.Time) {
	if len(dateArg) == 0 {
		return time.Time{}
	}
	t,err := 	time.Parse("2006-01-02", dateArg)
	if err != nil {
		exitWithErr(err)
	}
	return t
}

func _splitAndTrim (commaSeparatedArg string) []string {
	return strings.Split(commaSeparatedArg, ",")
}
var acArgs = acArgsS{}


// listDocumentsCmd represents the listDocuments command
var listActivityCmd = &cobra.Command{
	Use:   "listActivity",
	Short: "Lists events and actions",
	Long:`Lists activity 

		  rspace eln listFiles 
	`,

	Run: func(cmd *cobra.Command, args []string) {
		messageStdErr("listActivity called:")
		context := initialiseContext()  
		cfg,_ := configureActivityList()
		doListActivity(context, cfg)
	},
}

func configureActivityList () (*rspace.ActivityQuery, error) {
	builder := rspace.ActivityQueryBuilder{}
	for _,v := range acArgs.actions() {
		builder.Action(strings.TrimSpace(v))
	}
	builder.Domain("RECORD")
	for _,v := range acArgs.users() {
		builder.User(strings.TrimSpace(v))
	}
	
	builder.DateFrom(acArgs.after())
	builder.DateTo(acArgs.before())
	builder.Oid(acArgs.oid())
	return builder.Build()
}

func doListActivity (ctx *Context, cfg *rspace.ActivityQuery) {
	var docList *rspace.ActivityList
	var err error
	docList, err = ctx.WebClient.Activities(cfg)
	if err != nil {
		exitWithErr(err)
	}
	if ctx.Format.isJson() {
		ctx.write(prettyMarshal(docList))
	} else if ctx.Format.isQuiet() {
		printIds(ctx, toIdentifiableEvents(docList))
	} else {
		//processedResults := processResults(docList)
		activityListToTable(ctx, docList)
	}
}

func activityListToTable (ctx *Context, results *rspace.ActivityList) {
	var headers = []columnDef{columnDef{"Action",10},columnDef{"Timestamp", DISPLAY_TIMESTAMP_WIDTH}, 
		columnDef{"Id",10}, columnDef {"Name",25}, columnDef{"User",10}}

	rows := make([][]string, 0)
	for _, res := range results.Activities {
		info:=basicInfoFromPayload(res.Payload)
		data := []string {res.Action, res.Timestamp[:DISPLAY_TIMESTAMP_WIDTH],
			 info.GetGlobalId(),  info.GetName(), res.Username}
		rows = append(rows, data) 
	}
	if ctx.Format.isCsv() {
		printCsv(ctx, headers, rows)
	} else {
		printTable(ctx, headers, rows)
	}
}
func toIdentifiableEvents (result *rspace.ActivityList) []identifiable {
	rc := make([]identifiable, 0)
	for _,v := range result.Activities {
		payload := v.Payload
		// payload can be arbitrary data, usually but not always has an id value		
		if id:=basicInfoFromPayload(payload).GetId(); id > 0 {
			rc = append(rc, identifiable{strconv.Itoa(id)})
		}
	}
	return rc
}
//returns 0 if not exists
func basicInfoFromPayload (payload interface{}) rspace.BasicInfo {
	m, ok := payload.(map[string]interface{})
		if !ok {
    		exitWithStdErrMsg(fmt.Sprintf("want type map[string]interface{};  got %T", payload))
		}
		rc := rspace.IdentifiableNamable{}
		if id,ok:=m["id"]; ok {
			rc.GlobalId = id.(string)
			rc.Id,_=strconv.Atoi(id.(string)[2:])
		}
		if name,ok:=m["name"]; ok {
			rc.Name = name.(string)
		}
	return rc
}

func init() {
	elnCmd.AddCommand(listActivityCmd)

	initPaginationFromArgs(listActivityCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	listActivityCmd.PersistentFlags().StringVar(&acArgs.actionsArg, "actions", "", `Comma separated list of Actions 
	   (e.g. READ, WRITE.CREATE etc`)
	   listActivityCmd.PersistentFlags().StringVar(&acArgs.usersArg, "users", "", `Comma separated list of usernames`)
	   listActivityCmd.PersistentFlags().StringVar(&acArgs.afterDateArg, "afterDate", "", `Find events after this date, in format YYYY:MM:DD
		e.g. 2020-01-31`)
		listActivityCmd.PersistentFlags().StringVar(&acArgs.beforeDateArg, "beforeDate", "", `Find events before this date, in format YYYY:MM:DD
		e.g. 2020-01-31`)
		listActivityCmd.PersistentFlags().StringVar(&acArgs.globalId, "id", "", `Find events for a single document`)
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listDocumentsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
