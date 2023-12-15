package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

func main() {
	ctx, shutdown := context.WithCancel(context.Background())
	defer shutdown()

	var (
		wg sync.WaitGroup
		restartCh = make(chan struct{}, 1)
	)

	var srvErr error

	s := http.Server{
		Addr: ":9000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/restart" {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			restartCh <- struct{}{}

			w.Write([]byte("restarting web server"))
		}),
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer shutdown()

		srvErr = s.ListenAndServe()
	}()

	var cmdErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer shutdown()
		defer s.Shutdown(context.Background())

		for {
			// Run the website server and listen for restart signals 

			var (
				ctx, cancel = context.WithCancel(ctx)
				cmdWG sync.WaitGroup
				exiting bool
				restarting bool
			)

			cmd := exec.CommandContext(ctx, "go", "run", ".")
			cmd.Dir = "../../"
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmdWG.Add(1)
			go func() {
				defer cmdWG.Done()
				defer cancel()
				defer fmt.Println("web server shutdown")

				fmt.Println("starting web server")

				cmdErr = cmd.Run()
			}()

			cmdWG.Add(1)
			go func() {
				defer cmdWG.Done()
				defer cancel()

				select {
				case <-restartCh:
					fmt.Println("restart signal captured")
					restarting = true
				case <-ctx.Done():
					exiting = true
				}
			}()

			cmdWG.Wait()
				
			// double check the web server was killed.
			// TODO: fix this! make sure the .Run() call results in the web server being shutdown.

			lsofCmd := exec.Command("lsof", "-i", "tcp:8080")
			buf := bytes.NewBuffer(nil)
			lsofCmd.Stdout = buf
			lsofCmd.Run()

			b := buf.Bytes()
			if webServerStillRunning := len(b) > 0; webServerStillRunning {
				firstMatch := bytes.Split(b, []byte("\n"))[1]
				matchParts := bytes.Split(firstMatch, []byte(" "))

				var pid string
				for _, part := range matchParts[1:] { // skip COMMAND column
					if pid = string(part); pid != "" {
						break
					}
				}

				fmt.Println("web server is still running on process " + pid + "!")

				exec.Command("kill", pid).Run()

				fmt.Println("web server killed")
			}

			// Exit, restart, or capture error incurred

			if exiting {
				return
			}
			if cmdErr != nil && !restarting {
				fmt.Println("web server error incurred. exiting")
				return
			}

			fmt.Println("restarting...")
		}
	}()

	wg.Wait()

	var failed bool
	if srvErr != nil {
		fmt.Println("hot reloading server error", srvErr)
		failed = true
	}
	if cmdErr != nil {
		fmt.Println("web server error:", cmdErr)
		failed = true
	}

	if failed {
		os.Exit(1)
	}
}
