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
	"errors"
	"fmt"
	"gophercises/task/db"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

// doCmd represents the do command
var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task on your TODO list as complet",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least 1 argument")
		}
		for _, id := range args {
			if _, err := strconv.Atoi(id); err != nil {
				return errors.New("requires integers argument, got " + id)
			}
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		var taskIDs []int
		for _, i := range args {
			ID, err := strconv.Atoi(i)
			taskIDs = append(taskIDs, ID)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		tasks, err := db.Tasks()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, id := range taskIDs {
			if id < 1 || id > len(tasks) {
				fmt.Println("Invalid task number", id)
				continue
			}

			err = db.DeleteTask(tasks[id-1].Key)
			if err != nil {
				fmt.Printf("Failed to mark %v as complete. Error: %s\n", id, err.Error())
			} else {
				fmt.Printf("You have completed the \"%s\" task.\n", tasks[id-1].Value)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
