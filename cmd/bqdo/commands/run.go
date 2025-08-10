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
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"text/template"

	"cloud.google.com/go/bigquery"
	"github.com/wintermi/bqdo/internal/config"
	"google.golang.org/api/impersonate"
	"google.golang.org/api/option"
)

var (
	runConfigPath                string
	runDirectory                 string
	runProjectID                 string
	runDataset                   string
	runLocation                  string
	runImpersonateServiceAccount string
	runDryRun                    bool
)

var runCmd = &cobra.Command{
	Use:          "run",
	Short:        "Run a pipeline of BigQuery SQL queries",
	RunE:         runPipeline,
	SilenceUsage: true,
}

func runPipeline(cmd *cobra.Command, args []string) error {
	// Ensure the config file exists
	if _, err := os.Stat(runConfigPath); err != nil {
		return fmt.Errorf("config file %q not found: %w", runConfigPath, err)
	}

	// Resolve absolute config path and change working directory to its location
	absConfigPath, err := filepath.Abs(runConfigPath)
	if err != nil {
		return fmt.Errorf("resolve config path: %w", err)
	}
	configDir := filepath.Dir(absConfigPath)
	if chdirErr := os.Chdir(configDir); chdirErr != nil {
		return fmt.Errorf("change working directory to %q: %w", configDir, chdirErr)
	}

	// Load configuration
	cfg, err := config.Load(absConfigPath)
	if err != nil {
		return err
	}

	// Merge precedence: flags override config.
	projectID := firstNonEmpty(runProjectID, cfg.ProjectID)
	dataset := firstNonEmpty(runDataset, cfg.Dataset)
	location := firstNonEmpty(runLocation, cfg.Location)
	impersonateServiceAccount := firstNonEmpty(runImpersonateServiceAccount, cfg.ImpersonateServiceAccount)
	directory := firstNonEmpty(runDirectory, cfg.Directory)

	// Basic validations
	if projectID == "" {
		return fmt.Errorf("project ID is required (set in %s or via --project)", runConfigPath)
	}
	if directory == "" {
		return fmt.Errorf("directory is required in %s", runConfigPath)
	}
	if fi, err := os.Stat(directory); err != nil || !fi.IsDir() {
		return fmt.Errorf("directory %q not found or is not a directory", directory)
	}

	fmt.Printf("Using config: %s\n", absConfigPath)
	fmt.Printf("Directory: %s\n", directory)
	fmt.Printf("Project ID: %s\n", projectID)
	if dataset != "" {
		fmt.Printf("Dataset: %s\n", dataset)
	}
	if location != "" {
		fmt.Printf("Location: %s\n", location)
	}
	if impersonateServiceAccount != "" {
		fmt.Printf("Impersonate Service Account: %s\n", impersonateServiceAccount)
	}
	if len(cfg.Vars) > 0 {
		fmt.Printf("Vars: %v\n", cfg.Vars)
	}
	if runDryRun {
		fmt.Println("Dry run enabled: queries will be validated but not executed.")
	}

	// Authenticate BigQuery client
	ctx := context.Background()
	clientOpts := []option.ClientOption{}
	if impersonateServiceAccount != "" {
		ts, tsErr := impersonate.CredentialsTokenSource(ctx, impersonate.CredentialsConfig{
			TargetPrincipal: impersonateServiceAccount,
			Scopes:          []string{"https://www.googleapis.com/auth/cloud-platform"},
		})
		if tsErr != nil {
			return fmt.Errorf("impersonation setup failed: %w", tsErr)
		}
		clientOpts = append(clientOpts, option.WithTokenSource(ts))
	}
	bqClient, err := bigquery.NewClient(ctx, projectID, clientOpts...)
	if err != nil {
		return fmt.Errorf("create BigQuery client: %w", err)
	}
	defer func() { _ = bqClient.Close() }()

	// Collect SQL files recursively in alphabetical order
	var files []string
	walkErr := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		// Only process .sql files
		if strings.EqualFold(filepath.Ext(path), ".sql") {
			files = append(files, path)
		}
		return nil
	})
	if walkErr != nil {
		return fmt.Errorf("walk directory %q: %w", directory, walkErr)
	}
	sort.Strings(files)

	if len(files) == 0 {
		fmt.Printf("No .sql files found under %s\n", directory)
		return nil
	}

	fmt.Printf("Found %d file(s) to process.\n", len(files))

	// Prepare templating
	templateVars := map[string]string{}
	for k, v := range cfg.Vars {
		templateVars[k] = v
	}
	// Add commonly useful automatic variables
	if dataset != "" {
		templateVars["dataset"] = dataset
	}
	if projectID != "" {
		templateVars["project_id"] = projectID
	}

	for _, file := range files {
		start := time.Now()
		fmt.Printf("\n→ Processing %s\n", file)

		// Validate path is within configured directory (satisfies gosec G304)
		rel, rerr := filepath.Rel(directory, file)
		if rerr != nil || strings.HasPrefix(rel, "..") {
			return fmt.Errorf("invalid file path detected: %s", file)
		}
		// #nosec G304: file list is derived from a validated base directory and enumerated via WalkDir
		contents, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read %s: %w", file, err)
		}

		tmpl, err := template.New(filepath.Base(file)).Option("missingkey=error").Parse(string(contents))
		if err != nil {
			return fmt.Errorf("parse template for %s: %w", file, err)
		}
		var rendered bytes.Buffer
		if err := tmpl.Execute(&rendered, templateVars); err != nil {
			return fmt.Errorf("execute template for %s: %w", file, err)
		}

		sqlText := rendered.String()
		if strings.TrimSpace(sqlText) == "" {
			fmt.Printf("Skipping empty SQL in %s\n", file)
			continue
		}

		q := bqClient.Query(sqlText)
		if location != "" {
			q.Location = location
		}
		q.DryRun = runDryRun

		fmt.Printf("Executing%s...\n", func() string {
			if runDryRun {
				return " (dry-run)"
			}
			return ""
		}())
		job, err := q.Run(ctx)
		if err != nil {
			return fmt.Errorf("start job for %s: %w", file, err)
		}
		status, err := job.Wait(ctx)
		if err != nil {
			return fmt.Errorf("wait job for %s: %w", file, err)
		}
		if err := status.Err(); err != nil {
			return fmt.Errorf("job failed for %s: %w", file, err)
		}

		duration := time.Since(start)
		if runDryRun {
			fmt.Printf("Validated %s in %s\n", file, duration.Truncate(time.Millisecond))
		} else {
			fmt.Printf("Completed %s in %s\n", file, duration.Truncate(time.Millisecond))
		}
	}

	fmt.Println("\nAll files processed successfully.")
	return nil
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
	runCmd.Flags().StringVarP(&runDirectory, "directory", "r", "", "Directory containing .sql files to execute (overrides config)")
	runCmd.Flags().StringVarP(&runProjectID, "project", "p", "", "Google Cloud Project ID")
	runCmd.Flags().StringVarP(&runDataset, "dataset", "d", "", "BigQuery Dataset")
	runCmd.Flags().StringVarP(&runLocation, "location", "l", "", "BigQuery data processing location (e.g. australia-southeast1)")
	runCmd.Flags().StringVar(&runImpersonateServiceAccount, "impersonate-service-account", "", "Service account email to impersonate for Google Cloud API calls")
	runCmd.Flags().BoolVar(&runDryRun, "dry-run", false, "Dry run: validate and show actions without executing")
}
