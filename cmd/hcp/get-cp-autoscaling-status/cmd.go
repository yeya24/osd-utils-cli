package getcpautoscalingstatus

import (
	"github.com/spf13/cobra"
)

// NewCmdGetCPAutoscalingStatus creates and returns the get-cp-autoscaling-status command
func NewCmdGetCPAutoscalingStatus() *cobra.Command {
	return newCmdAutoscalingAudit()
}
