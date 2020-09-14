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
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
	//"fmt"
)

type addUserArgs struct {
	UsernameArg, FNameArg, LNameArg, EmailArg, RoleArg, AffiliationArg, PasswordFileArg, UserCsvFileArg string
}

var userArgs = addUserArgs{}

// addUserCmd represents the createNotebook command
var addUserCmd = &cobra.Command{
	Use:   "addUser",
	Short: "Adds a new user account",
	Long: `Requires sysadmin permission. Add a new user account rm, equires username and email
	at a minimum. Supply an initial password in a file. Alternatively, supply a CSV file of new users, 1 per row.
	`,
	Example: ` 
	rspace eln addUser --username newusername --email someone@somwhere.com --role user|pi| --pwdfile passwordfile

	// from a CSV file
	rspace eln addUser --userfile users.csv 
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(userArgs.UserCsvFileArg) > 0 {
			createUsersFromFile()
		} else {
			createSingleUser()
		}

	},
}

func readUserCsvFile() [][]string {
	files := []string{userArgs.UserCsvFileArg}
	validateInputFilePaths(files)
	file, _ := os.Open(files[0])
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		exitWithErr(err)
	}
	return records
}

func createUsersFromFile() {
	records := readUserCsvFile()
	// TODO test with file, generate file
	ctx := initialiseContext()
	results := make(chan *UserResult, 5)
	tasks := make(chan *rspace.UserPost, 5)
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go doSubmit(&wg, ctx, tasks, results)
	}
	for i, v := range records {
		if i == 0 {
			continue
		}
		messageStdErr(fmt.Sprintf("Creating user from row %d", i))
		fields := v
		builder := &rspace.UserPostBuilder{}
		builder.Password(string(fields[5])).Username(fields[4])
		builder.Email(rspace.Email(fields[2])).FirstName(fields[0])
		builder.LastName(fields[1]).Affiliation(fields[6])
		builder.ApiKey(fields[7])
		userPost, _ := builder.Role(getRoleForArg(fields[3])).Build()
		tasks <- userPost
	}
	close(tasks)
	for i := 0; i < len(records)-1; i++ {
		res := <-results
		if res.success != nil {
			ctx.write(prettyMarshal(res.success))
		}
	}
	messageStdErr(fmt.Sprintf("waiting, tasks = %d, results = %d", len(tasks), len(results)))

	wg.Wait()
	messageStdErr("done")

}

func doSubmit(wg *sync.WaitGroup, ctx *Context, tasks <-chan *rspace.UserPost,
	results chan<- *UserResult) {
	defer func() {
		wg.Done()
		if err := recover(); err != nil {
			log.Println("user creation failed..", err)
		}
	}()

	for userPost := range tasks {
		user, err := ctx.WebClient.UserNew(userPost)
		result := &UserResult{user, err}
		results <- result
	}

}

// encapsulates result of attempt to create a user
type UserResult struct {
	success *rspace.UserInfo
	failure error
}

func createSingleUser() {
	userPost := validateFlags()
	ctx := initialiseContext()
	user, err := ctx.WebClient.UserNew(userPost)
	if err != nil {
		exitWithErr(err)
	} else {
		ctx.write(prettyMarshal(user))
	}
}

func validateFlags() *rspace.UserPost {
	pwd, err := ioutil.ReadFile(userArgs.PasswordFileArg)
	if err != nil {
		exitWithStdErrMsg("No password file supplied. Please put user password in a file and use the 'pwdfile' argument")
	}
	builder := &rspace.UserPostBuilder{}
	builder.Password(string(pwd)).Username(userArgs.UsernameArg)
	builder.Email(rspace.Email(userArgs.EmailArg))
	builder.FirstName(userArgs.FNameArg).LastName(userArgs.LNameArg)
	builder.Affiliation(userArgs.AffiliationArg)
	post, e := builder.Role(getRoleForArg(userArgs.RoleArg)).Build()
	if e != nil {
		exitWithErr(e)
	}
	return post
}

func getRoleForArg(arg string) rspace.UserRoleType {
	if arg == "pi" || arg == "ROLE_PI" {
		return rspace.Pi
	} else {
		return rspace.User
	}
}

func init() {
	elnCmd.AddCommand(addUserCmd)
	addUserCmd.Flags().StringVar(&userArgs.UsernameArg, "username", "", "username")
	addUserCmd.Flags().StringVar(&userArgs.FNameArg, "first", "Unknown", "First name")
	addUserCmd.Flags().StringVar(&userArgs.LNameArg, "last", "Unknown", "Last name")
	addUserCmd.Flags().StringVar(&userArgs.EmailArg, "email", "", "Valid email address")
	addUserCmd.Flags().StringVar(&userArgs.RoleArg, "role", "user", "Role, either 'pi' or 'user'")
	addUserCmd.Flags().StringVar(&userArgs.AffiliationArg, "affiliation", "unknown", "Affiliation (Community only)")

	addUserCmd.Flags().StringVar(&userArgs.PasswordFileArg, "pwdfile", "", "a file containing the password")
	addUserCmd.Flags().StringVar(&userArgs.UserCsvFileArg, "userfile", "", "a CSV file of new users")

}
