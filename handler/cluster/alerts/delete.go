/*
Copyright © 2019 Portworx

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

package alerts

import (
	"github.com/portworx/pxc/pkg/cliops"
	"github.com/portworx/pxc/pkg/commander"
	"github.com/portworx/pxc/pkg/portworx"
	"github.com/portworx/pxc/pkg/util"
	"github.com/spf13/cobra"
)

var deleteAlertsCmd *cobra.Command

var _ = commander.RegisterCommandVar(func() {
	deleteAlertsCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete Portworx alerts",
		Example: `
  # Delete portworx related alerts
  pxc alert delete

  # Delete alerts based on particular alert type. Delete all alerts related to "volume"
  pxc alert delete -t "volume"`,
		RunE: deleteAlertsExec,
	}
})

var _ = commander.RegisterCommandInit(func() {
	AlertAddCommand(deleteAlertsCmd)
	deleteAlertsCmd.Flags().StringP("type", "t", "all", "alert type (Valid Values: [volume node cluster drive all])")
	//TODO: Need to support more flags
})

func deleteAlertsExec(cmd *cobra.Command, args []string) error {
	ctx, conn, err := portworx.PxConnectDefault()
	_ = ctx
	if err != nil {
		return err
	}
	defer conn.Close()
	// Parse out all of the common cli volume flags
	cai := cliops.GetCliAlertInputs(cmd, args)

	// Create a cliVolumeOps object
	alertOps := cliops.NewCliAlertOps(cai)

	// initialize alertOP interface
	alertOps.PxAlertOps = portworx.NewPxAlertOps()

	err = alertOps.PxAlertOps.DeletePxAlerts(alertOps.CliAlertInputs.AlertType)

	if err == nil {
		util.Printf("Alerts deleted successfully\n")
	}

	return err
}
