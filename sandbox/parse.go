package sandbox

import "encoding/json"

// docker sandbox ls --json の出力構造
type lsOutput struct {
	VMs []rawVM `json:"vms"`
}

type rawVM struct {
	Name       string   `json:"name"`
	Agent      string   `json:"agent"`
	Status     string   `json:"status"`
	SocketPath string   `json:"socket_path"`
	Workspaces []string `json:"workspaces"`
}

// ParseSandboxList は docker sandbox ls --json の出力をパースする
func ParseSandboxList(data []byte) ([]Sandbox, error) {
	var out lsOutput
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	sandboxes := make([]Sandbox, len(out.VMs))
	for i, vm := range out.VMs {
		workspace := ""
		if len(vm.Workspaces) > 0 {
			workspace = vm.Workspaces[0]
		}
		sandboxes[i] = Sandbox{
			ID:        vm.Name, // IDとして Name を使う
			Name:      vm.Name,
			Agent:     vm.Agent,
			State:     normalizeState(vm.Status),
			Workspace: workspace,
			Risk:      assessRisk(vm),
		}
	}
	return sandboxes, nil
}

func normalizeState(s string) State {
	switch s {
	case "running":
		return StateRunning
	case "stopped", "exited":
		return StateStopped
	default:
		return StateCrashed
	}
}

// assessRisk はCLIから取れる情報でRiskを暫定判定する
// MVPでは情報不足なので RiskDev をデフォルトにする
func assessRisk(vm rawVM) Risk {
	return RiskDev
}
