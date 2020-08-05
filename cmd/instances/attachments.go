// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package instances

import (
	"github.com/spf13/cobra"
)

//AttachCmd to manage instances
var AttachCmd = &cobra.Command{
	Use:   "attachments",
	Short: "Manage environments to instances",
	Long:  "Manage environments to instances",
}

var environment string

func init() {

	AttachCmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	AttachCmd.AddCommand(CreateAttachCmd)
	AttachCmd.AddCommand(ListAttachCmd)
	AttachCmd.AddCommand(GetAttachCmd)
	AttachCmd.AddCommand(DeleteAttachCmd)
}
