package builder

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os/exec"
	"strings"
)

func NewDoctorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "detects required dependencies and tools",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkRequiredBinary(); len(err) > 0 {
				for _, e := range err {
					fmt.Printf("- %v\n", e)
				}
				return fmt.Errorf("missing required binaries")
			}
			if err := checkExecutable(); err != nil {
				return err
			}

			fmt.Println("all good")
			return nil
		},
	}

	return cmd
}

func checkRequiredBinary() []error {
	var required = []string{
		"rootlesskit",
		"buildkitd",
		"buildctl",
		"newuidmap", // sudo apt install uidmap
		"newgidmap",
		// TODO(iyear): add more
	}

	errs := make([]error, 0)
	for _, r := range required {
		_, err := exec.LookPath(r)
		// go1.19
		if errors.Is(err, exec.ErrDot) {
			err = nil
		}
		if err != nil {
			errs = append(errs, fmt.Errorf("lookup %s failed: %w", r, err))
		}
	}

	return errs
}

// checkExecutable checks if buildkitd can run with rootlesskit
func checkExecutable() error {
	const running = "running server on"
	cmd := exec.Command("rootlesskit", "buildkitd")

	stdout, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	reader := bufio.NewReader(stdout)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		// if contains 'running server on', then it's running
		if strings.Contains(s, running) {
			return cmd.Process.Kill()
		}
	}

	b, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return fmt.Errorf("buildkitd can't run, please check the logs")
}
