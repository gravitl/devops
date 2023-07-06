/*
Copyright © 2023 Matthew R Kasun <mkasun@nusak.ca>

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
	"os"

	"github.com/gravitl/devops/do"
	"github.com/gravitl/devops/netmaker"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

// deleteEgressSnapshotCmd represents the deleteEgressSnapshot command
var deleteEgressSnapshotCmd = &cobra.Command{
	Use:   "deleteEgressSnapshot",
	Short: "deletes EgressSnapshot ",
	Args:  cobra.ExactArgs(0),
	Long:  "deletes EgressSnapshot ",
	Run: func(cmd *cobra.Command, args []string) {
		setupLoging("deleteEgressSnapshot")
		deleteEgressSnapshot(&config)
		fmt.Println("deleteEgressSnapshot called")
	},
}

func init() {
	rootCmd.AddCommand(deleteEgressSnapshotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteEgressSnapshotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteEgressSnapshotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func deleteEgressSnapshot(c *netmaker.Config) {
	request, _, _, _ := do.Default()
	request.Token = c.DigitalOcean_Token
	name := "egresssnapshot" + c.Tag
	slog.Info("getting snapshot", "name", name)
	snapshot, err := request.ListSnapshot(name)
	if err != nil {
		slog.Error("error getting snapshot", "name", name, "error", err)
		os.Exit(1)
	}
	slog.Info("deleting snapshot egresssnapshot", "ID", snapshot.ID)
	if err := request.DeleteSnapshot(snapshot.ID); err != nil {
		slog.Error("error deleting snapshot", "name", name)
		os.Exit(1)
	}
	slog.Info("snapshot deleted")
}
