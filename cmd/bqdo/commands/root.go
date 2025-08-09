// Copyright 2022-2025, Matthew Winter
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wintermi/bqdo/internal/buildinfo"
)

var (
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "bqdo",
	Short: "bqdo is a CLI for executing BigQuery SQL as part of a pipeline",
	Long:  "bqdo is a CLI for executing BigQuery SQL as part of a pipeline.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("bqdo %s (commit %s)\n", buildinfo.Version, buildinfo.Commit)
		},
	})

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	}

	cobra.OnInitialize(func() {
		_ = os.Setenv("BQDO_CLI", "1")
	})
}
