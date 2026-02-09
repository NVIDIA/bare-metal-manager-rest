package subnet

import "fmt"

// Init Subnet
func (sub *API) Init() {
	// Validate the Subnet config later
	ManagerAccess.Data.EB.Log.Info().Msg("Subnet: Initializing the Subnet")
}

// GetState Subnet
func (sub *API) GetState() []string {
	state := ManagerAccess.Data.EB.Managers.Workflow.SubnetState
	var strs []string
	strs = append(strs, fmt.Sprintln("subnet_workflow_started", state.WflowStarted.Load()))
	strs = append(strs, fmt.Sprintln("subnet_workflow_activity_failed", state.WflowActFail.Load()))
	strs = append(strs, fmt.Sprintln("subnet_worflow_activity_succeeded", state.WflowActSucc.Load()))
	strs = append(strs, fmt.Sprintln("subnet_workflow_publishing_failed", state.WflowPubFail.Load()))
	strs = append(strs, fmt.Sprintln("subnet_worflow_publishing_succeeded", state.WflowPubSucc.Load()))

	return strs
}
