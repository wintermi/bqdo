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

package defaults

// BqdoTOML is the default scaffold for bqdo.toml written by the init command.
// The format matches the configuration structure defined in the config package.
var BqdoTOML = []byte(`# bqdo default config (TOML)
# This file is written by bqdo init. Adjust fields to suit your project.

directory = "sql/"
project_id = "your-project-id"
dataset = "your_dataset"
location = "US"
impersonate_service_account = ""

[vars]
env = "dev"
start_date = "2025-01-01"
`)
