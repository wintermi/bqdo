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
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wintermi/bqdo/internal/config"
	"github.com/wintermi/bqdo/internal/defaults"
)

var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "Initialize a bqdo pipeline in the current project",
	RunE:         runInitCmd,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInitCmd(cmd *cobra.Command, args []string) error {
	// Check for existing default config file and ask for overwrite confirmation
	if _, err := os.Stat(config.DefaultConfigFilename); err == nil {
		fmt.Printf("A %s already exists in this directory. Overwrite? [y/N]: ", config.DefaultConfigFilename)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			fmt.Println("Aborted. No files were changed.")
			return nil
		}
	}

	// Write default config from embedded template
	// If file exists, it will be overwritten only if user confirmed above
	if err := os.WriteFile(config.DefaultConfigFilename, defaults.BqdoTOML, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", config.DefaultConfigFilename, err)
	}
	fmt.Printf("Created %s\n", config.DefaultConfigFilename)
	return nil
}
