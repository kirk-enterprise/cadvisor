package gpu

import "github.com/google/cadvisor/gpu/cmd"

type GPUMonitor interface {
	GetGPUFBSize() map[string]map[string]string
	GetGPUUtils() map[string]map[string][]string
}

func NewGPuMonitor() GPUMonitor {
	return &cmd.CmdGPUMonitor{}
}
