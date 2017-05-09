//Package cmd uses raw cmd to get the per-pid gpu sm util and mem util
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type CmdGPUMonitor struct {
}

// IsCmdExist checks whether nvidia-smi cmd exist or not
func (c *CmdGPUMonitor) IsCmdExist() bool {
	_, err := exec.LookPath("nvidia-smi")
	if err != nil {
		fmt.Println("nvidia-smi is not in your PATH")
		return false
	}

	return true
}

// GetGPUUtils runs nvidia-smi command and get per pid sm util and mem util
func (c *CmdGPUMonitor) GetGPUUtils() (map[string]map[string][]string, error) {
	// first key is pid, second key is device id, value is two-elemet util slice, first element is sm util
	// second element is mem util
	res := make(map[string]map[string][]string)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)

	defer cancel()

	cmd := exec.CommandContext(ctx, "nvidia-smi", "pmon", "-c", "1")

	env := os.Environ()
	env = append(env, "NVSMI_SHOW_ALL_DEVICES=1")
	cmd.Env = env

	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	// sample output
	/*# gpu     pid  type    sm   mem   enc   dec   command
	# Idx       #   C/G     %     %     %     %   name
		0       -     -     -     -     -     -   -
		1       -     -     -     -     -     -   -
		2       -     -     -     -     -     -   -
		3       -     -     -     -     -     -   -
		4       -     -     -     -     -     -   -
		5       -     -     -     -     -     -   -
		6       -     -     -     -     -     -   -
		7       -     -     -     -     -     -   -
		8       -     -     -     -     -     -   -
		9       -     -     -     -     -     -   -
		10   64756     C     0     0     0     0   pulpf
		11       -     -     -     -     -     -   -
		12       -     -     -     -     -     -   -
		13       -     -     -     -     -     -   -
		14       -     -     -     -     -     -   -
		15 1426541     C    66    26     0     0   python */

	for c, line := range strings.Split(string(out), "\n") {
		vals := strings.Fields(line)
		if c < 2 || len(vals) != 8 {
			continue
		}

		if vals[1] == "-" {
			continue
		}

		res[vals[1]] = make(map[string][]string)
		res[vals[1]][vals[0]] = append(res[vals[1]][vals[0]], vals[3])
		res[vals[1]][vals[0]] = append(res[vals[1]][vals[0]], vals[4])
	}

	return res, nil
}

// GetGPUFBSize runs nvidia-smi command and get fb utilization
func (c *CmdGPUMonitor) GetGPUFBSize() (map[string]map[string]string, error) {
	// first key is pid, second key is device id, value is fb size(unit is MB)
	res := make(map[string]map[string]string)
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)

	defer cancel()

	cmd := exec.CommandContext(ctx, "nvidia-smi", "pmon", "-c", "1", "-s", "m")

	env := os.Environ()
	env = append(env, "NVSMI_SHOW_ALL_DEVICES=1")
	cmd.Env = env

	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	// sample output
	/*# gpu     pid  type    fb   command
	# Idx       #   C/G    MB   name
	    0       -     -     -   -
	    1       -     -     -   -
	    2       -     -     -   -
	    3       -     -     -   -
	    4       -     -     -   -
	    5       -     -     -   -
	    6       -     -     -   -
	    7       -     -     -   -
	    8       -     -     -   -
	    9       -     -     -   -
	   10   38148     C   284   pulpf
	   11       -     -     -   -
	   12       -     -     -   -
	   13       -     -     -   -
	   14       -     -     -   -
	   15       -     -     -   - */

	for c, line := range strings.Split(string(out), "\n") {
		vals := strings.Fields(line)
		if c < 2 || len(vals) != 5 {
			continue
		}

		if vals[1] == "-" {
			continue
		}

		res[vals[1]] = make(map[string]string)
		res[vals[1]][vals[0]] = vals[3]
	}

	return res, nil
}
