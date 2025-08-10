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

	"github.com/spf13/cobra"

	"github.com/wintermi/bqdo/internal/config"
)

var (
	runConfigPath                string
	runProjectID                 string
	runDataset                   string
	runLocation                  string
	runImpersonateServiceAccount string
	runDryRun                    bool
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a pipeline of BigQuery SQL queries",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(runConfigPath)
		if err != nil {
			return err
		}

		// Merge precedence: flags override config.
		projectID := firstNonEmpty(runProjectID, cfg.ProjectID)
		dataset := firstNonEmpty(runDataset, cfg.Dataset)
		location := firstNonEmpty(runLocation, cfg.Location)
		impersonateServiceAccount := firstNonEmpty(runImpersonateServiceAccount, cfg.ImpersonateServiceAccount)

		fmt.Printf("Running bqdo with config: %s\n", runConfigPath)
		fmt.Printf("Directory: %s\n", cfg.Directory)
		fmt.Printf("Project ID: %s\n", projectID)
		fmt.Printf("Dataset: %s\n", dataset)
		fmt.Printf("Location: %s\n", location)
		if impersonateServiceAccount != "" {
			fmt.Printf("Impersonate Service Account: %s\n", impersonateServiceAccount)
		}
		if len(cfg.Vars) > 0 {
			fmt.Printf("Vars: %v\n", cfg.Vars)
		}
		if runDryRun {
			fmt.Println("Dry run enabled: queries will be validated but not executed.")
		}
		return nil
	},
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&runConfigPath, "config", "c", config.DefaultConfigFilename, "Path to the configuration file")
	runCmd.Flags().StringVarP(&runProjectID, "project", "p", "", "Google Cloud Project ID")
	runCmd.Flags().StringVarP(&runDataset, "dataset", "d", "", "BigQuery Dataset")
	runCmd.Flags().StringVarP(&runLocation, "location", "l", "", "BigQuery data processing location (e.g. australia-southeast1)")
	runCmd.Flags().StringVar(&runImpersonateServiceAccount, "impersonate-service-account", "", "Service account email to impersonate for Google Cloud API calls")
	runCmd.Flags().BoolVar(&runDryRun, "dry-run", false, "Dry run: validate and show actions without executing")
}
