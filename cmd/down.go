package cmd

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/process"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downCmd)
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop running holesail processes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		processes, err := process.Processes()
		if err != nil {
			fmt.Println(err)
			return
		}

		var stoppedProcesses int

		for _, process := range processes {
			processName, _ := process.Name()
			cmdLine, _ := process.Cmdline()

			if processName == "node" {
				cmdLineSplitted := strings.Split(cmdLine, "/")
				cmdLineWithoutBloat := strings.Join(cmdLineSplitted[len(cmdLineSplitted)-3:], "/")

				if strings.HasPrefix(cmdLineWithoutBloat, "node_modules/holesail/index.js") {
					process.SendSignal(syscall.SIGINT)
					stoppedProcesses++
				}
			}
		}

		fmt.Printf("Stopped %v holesail processes\n", stoppedProcesses)
	},
}
