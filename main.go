package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

//PIDFile ..
var PIDFile = "/tmp/kulwinder.pid"

func savePID(pid int) {

	file, err := os.Create(PIDFile)
	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	defer file.Close()

	_, err = file.WriteString(strconv.Itoa(pid))

	if err != nil {
		log.Printf("Unable to create pid file : %v\n", err)
		os.Exit(1)
	}

	file.Sync() // flush to disk

}

func main() {

	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s [start|stop] \n ", os.Args[0]) // return the program name back to %s
		os.Exit(0)                                            // graceful exit
	}

	var osData string
	if len(os.Args) != 2 {
		osData = ""
	} else {
		osData = os.Args[1]
	}
	if strings.ToLower(osData) == "main" {

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGTERM)

		go func() {
			signalType := <-ch
			signal.Stop(ch)
			fmt.Println("Exit command received. Exiting...")

			// this is a good place to flush everything to disk
			// before terminating.
			fmt.Println("Received signal type : ", signalType)

			// remove PID file
			os.Remove(PIDFile)

			os.Exit(0)

		}()

		//---------
		// Your Projects main function here
		//---------

	}

	if strings.ToLower(os.Args[1]) == "start" {
		if _, err := os.Stat(PIDFile); err == nil {
			fmt.Println("Already running or /tmp/kulwinder.pid file exist.")
			os.Exit(1)
		}

		cmd := exec.Command(os.Args[0], "main")
		cmd.Start()
		fmt.Println("Daemon process ID is : ", cmd.Process.Pid)
		savePID(cmd.Process.Pid)
		os.Exit(0)

	}
	if strings.ToLower(os.Args[1]) == "stop" {
		if _, err := os.Stat(PIDFile); err == nil {
			data, err := ioutil.ReadFile(PIDFile)
			if err != nil {
				fmt.Println("Not running")
				os.Exit(1)
			}
			ProcessID, err := strconv.Atoi(string(data))

			if err != nil {
				fmt.Println("Unable to read and parse process id found in ", PIDFile)
				os.Exit(1)
			}

			process, err := os.FindProcess(ProcessID)

			if err != nil {
				fmt.Printf("Unable to find process ID [%v] with error %v \n", ProcessID, err)
				os.Exit(1)
			}
			// remove PID file
			os.Remove(PIDFile)

			fmt.Printf("Killing process ID [%v] now.\n", ProcessID)
			// kill process and exit immediately
			err = process.Kill()

			if err != nil {
				fmt.Printf("Unable to kill process ID [%v] with error %v \n", ProcessID, err)
				os.Exit(1)
			} else {
				fmt.Printf("Killed process ID [%v]\n", ProcessID)
				os.Exit(0)
			}

		} else {

			fmt.Println("Not running.")
			os.Exit(1)
		}
	} else {
		fmt.Printf("Unknown command : %v\n", os.Args[1])
		fmt.Printf("Usage : %s [start|stop]\n", os.Args[0]) // return the program name back to %s
		os.Exit(1)
	}

}
