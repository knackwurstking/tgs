package journal

import (
	"fmt"
	"log/slog"
	"os/exec"
)

type Unit struct {
	// Name or the unit file to get
	Name string `json:"name" yaml:"name"`
	// Output will be used for the shell command `journalctl` as `--output ${output}`
	//
	// optional
	Output string `json:"output,omitempty" yaml:"output,omitempty"`
}

type Units struct {
	System []Unit `json:"system,omitempty" yaml:"system,omitempty"`
	User   []Unit `json:"user,omitempty" yaml:"user,omitempty"`
}

func NewUnits() *Units {
	return &Units{
		System: make([]Unit, 0),
		User:   make([]Unit, 0),
	}
}

func (u *Units) GetSystemUnit(name string) (*Unit, error) {
	for i, unit := range u.System {
		if unit.Name == name {
			return &u.System[i], nil
		}
	}

	return nil, fmt.Errorf("user unit %s not found", name)
}

func (u *Units) GetUserUnit(name string) (*Unit, error) {
	for i, unit := range u.User {
		if unit.Name == name {
			return &u.User[i], nil
		}
	}

	return nil, fmt.Errorf("system unit %s not found", name)
}

func (u *Units) GetOutput(name string) (data []byte, err error) {
	isUser := true

	unit, err := u.GetUserUnit(name)
	if err != nil {
		isUser = false

		unit, err = u.GetSystemUnit(name)
		if err != nil {
			return nil, err
		}
	}

	var output string
	if unit.Output == "default" || unit.Output == "" {
		output = "short"
	} else {
		output = unit.Output
	}

	var cmd *exec.Cmd
	if isUser {
		cmd = exec.Command("bash", "-c",
			fmt.Sprintf(
				"journalctl --user -u %s --output %s --reverse --no-pager | sed 's/\x1b\\[[0-9;]*m//g'",
				unit.Name, output,
			),
		)
	} else {
		cmd = exec.Command("journalctl",
			"-u", unit.Name,
			"--output", output,
			"--reverse",
			"--no-pager",
		)
	}

	slog.Debug("Run journalctl command", "args", cmd.Args)
	if data, err = cmd.CombinedOutput(); err != nil {
		return nil, err
	} else {
		return data, nil
	}
}
