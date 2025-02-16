package cmd

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start holesail proxy",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("Starting proxy")

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if len(strings.Split(r.Host, ".")) == 1 {
				w.Write([]byte("No holesail connector provided"))
				return
			}

			connector := strings.Split(r.Host, ".")[0]

			findProcess := exec.Command("pgrep", "-f", "/holesail/index.js "+connector)

			out, _ := findProcess.Output()
			splitted := strings.Split(string(out), "\n")

			var isRunning bool
			var freePort int

			for _, pid := range splitted {
				getDirectory := exec.Command("readlink", fmt.Sprintf("/proc/%s/cwd", string(pid)))
				out, _ = getDirectory.Output()

				wd, _ := os.Getwd()

				if strings.TrimSpace(string(out)) == wd {
					isRunning = true

					getCommand := exec.Command("cat", fmt.Sprintf("/proc/%s/cmdline", string(pid)))
					out, _ = getCommand.Output()

					outSplitted := strings.Split(string(out), "/")

					outWithoutBloat := outSplitted[len(outSplitted)-1]
					outWithoutBloat = strings.Split(outWithoutBloat, "--port")[1]
					outWithoutBloat = strings.Trim(outWithoutBloat, "\x00")

					freePort, _ = strconv.Atoi(outWithoutBloat)
				}
			}

			if !isRunning {
				freePort, _ = getFreePort()

				cmd := exec.Command("holesail", connector, "--port", strconv.Itoa(freePort))
				cmd.Start()

				// Wait for the connection to be made
				for {
					_, err := http.Get("http://localhost:" + strconv.Itoa(freePort))
					if err == nil {
						break
					}
					time.Sleep(500 * time.Millisecond)
				}
			}

			appUrl, _ := url.Parse("http://localhost:" + strconv.Itoa(freePort))
			httputil.NewSingleHostReverseProxy(appUrl).ServeHTTP(w, r)
		})

		err := http.ListenAndServe(":80", nil)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}
