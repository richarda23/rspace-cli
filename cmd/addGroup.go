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
	"errors"
	"strings"

	"github.com/richarda23/rspace-client-go/rspace"
	"github.com/spf13/cobra"
	//"fmt"
)

type addGroupArgs struct {
	PiArg, GroupNameArg, MembersArg string
}

var groupArgs = addGroupArgs{}

// addUserCmd represents the createNotebook command
var addGroupCmd = &cobra.Command{
	Use:   "addGroup",
	Short: "Adds a new LabGroup",
	Long: `Requires sysadmin permission. Creates a new lab group with a PI and 
	 optionally other group members. All users must be already existing in RSpace.

	 Note that members must have also initialised their accounts by logging in at least once
	`,
	Example: ` 
## create a group with 2 members and a PI
rspace eln addGroup --name Bob-Group --pi bobsmith --members anabelz,sarahs

# if the group name has spaces, enclose in double-quotes
rspace eln addGroup --name "Prof Smith Lab" --pi bobsmith --members anabelz,sarahs
`,
	Run: func(cmd *cobra.Command, args []string) {
		groupPost, err := validateGrpFlags()
		if err != nil {
			exitWithErr(err)
		}
		ctx := initialiseContext()
		group, err := ctx.WebClient.GroupNew(groupPost)
		if err != nil {
			exitWithErr(err)
		} else {

			list := rspace.GroupList{[]*rspace.GroupInfo{group}}
			formatter := &GroupListFormatter{&list}
			ctx.writeResult(formatter)
		}
	},
}

func validateGrpFlags() (*rspace.GroupPost, error) {
	if len(groupArgs.GroupNameArg) == 0 {
		return nil, errors.New("Group name is required, using --name flag")
	}
	var userGroupPosts []rspace.UserGroupPost = make([]rspace.UserGroupPost, 0, 5)
	if len(groupArgs.PiArg) == 0 {
		return nil, errors.New("PI username is required, using --pi flag")
	} else {
		userGroupPosts = append(userGroupPosts, rspace.UserGroupPost{groupArgs.PiArg, "PI"})
	}
	if len(groupArgs.MembersArg) > 0 {
		members := strings.Split(groupArgs.MembersArg, ",")
		for _, v := range members {
			username := strings.TrimSpace(v)
			if len(username) == 0 {
				messageStdErr("empty username, skipping")
			} else {
				userGroupPosts = append(userGroupPosts, rspace.UserGroupPost{username, "DEFAULT"})
			}
		}
	}

	groupPost, err := rspace.GroupPostNew(groupArgs.GroupNameArg, userGroupPosts)
	if err != nil {
		return nil, err
	}
	return groupPost, nil
}

func init() {
	elnCmd.AddCommand(addGroupCmd)
	addGroupCmd.Flags().StringVar(&groupArgs.GroupNameArg, "name", "", "The name of the group")
	addGroupCmd.Flags().StringVar(&groupArgs.PiArg, "pi", "", "Username of the group PI")
	addGroupCmd.Flags().StringVar(&groupArgs.MembersArg, "members", "", "Comma-separated list of usernames of other group members")

}
