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
	"io/ioutil"
	"rspace"
	"fmt"
)
type addUserArgs struct {
 UsernameArg string
  EmailArg, RoleArg, PasswordFileArg string
}
var userArgs = addUserArgs{}
// addUserCmd represents the createNotebook command
var addUserCmd = &cobra.Command{
	Use:   "addUser",
	Short: "Adds a new user account",
	Long: `Requires sysadmin permission. Add a new user account
	  addUser --username newusername --email someone@somwhere.com --role user|pi| --pwdfile passwordfile

	  or minimal shorthand as comma-separated triple:

	        addUser username,email,role,passwordfile

	`,
	Run: func(cmd *cobra.Command, args []string) {
		userPost := validateArgs()
		ctx := initialiseContext()
		user, err := ctx.WebClient.UserNew(userPost)
		if err != nil {
			exitWithErr(err)
		} else {
			ctx.write(prettyMarshal(user))
		}
	},
}

func validateArgs () *rspace.UserPost {
	pwd,err := ioutil.ReadFile(userArgs.PasswordFileArg)
	if err != nil {
		exitWithErr(err)
	}
	builder := &rspace.UserPostBuilder{}
	builder.Password(string(pwd)).Username(userArgs.UsernameArg)
	builder.Email(rspace.Email(userArgs.EmailArg))
	builder.Affiliation("Unknown")
	post,e:= builder.Role(getRoleForArg(userArgs.RoleArg)).FirstName("unknown").LastName("unknown").Build()
	if e != nil {
		exitWithErr(e)
	}
	fmt.Println(post)
	return post;
}

func getRoleForArg(arg string) rspace.UserRoleType {
	if arg == "pi" {
		return rspace.Pi
	} else {
		return rspace.User
	}
}

func init() {
	elnCmd.AddCommand(addUserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addUserCmd.PersistentFlags().String("name", "n","",  "A name for the notebook")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	 addUserCmd.Flags().StringVar(&userArgs.UsernameArg, "username", "", "Username, >= 6 chars")
	 addUserCmd.Flags().StringVar(&userArgs.EmailArg, "email", "", "Valid email address")
	 addUserCmd.Flags().StringVar(&userArgs.RoleArg, "role", "user", "Role, either 'pi' or 'user'")
	 addUserCmd.Flags().StringVar(&userArgs.PasswordFileArg, "pwdfile", "", "a file containing the password")
	}
