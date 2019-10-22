/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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
	"gophercises/task/db"
	"os"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your incomplete tasks",
	Args:  cobra.NoArgs,

	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.Tasks()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		if len(tasks) == 0 {
			fmt.Println("You have no tasks to do ! Take a vacation 🏖")
		} else {
			fmt.Println("You have the following tasks:")
			for i, t := range tasks {
				fmt.Printf("%v. %v\n", i+1, t.Value)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
