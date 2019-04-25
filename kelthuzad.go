package main

import (
	"github.com/hpcloud/tail"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"os/exec"
	"regexp"
	"syscall"
	"time"
)

var logger *log.Logger

// a function to respawn the process
func spawn(cmdPath string) *exec.Cmd {
	cmd := exec.Command(cmdPath)

	// this block is necessary when killing a subprocess properly
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	go func() {
		err := cmd.Start()
		if err != nil {
			logger.Fatalln(err)
		}
		defer cmd.Wait()
		logger.Printf("%v is spawned\n", cmd.Process.Pid)
	}()

	// return the created Cmd struct
	return cmd
}

func main() {

	// the argument options
	var opts struct {
		LogPath string `short:"p" long:"path" description:"The path of the log" required:"true"`
		CmdPath string `short:"c" long:"command" description:"The path of a command string to respawn the process" required:"true"`
		Regex   string `short:"r" long:"regex" description:"The regex pattern to detect a failure" required:"true"`
		Verbose bool   `short:"v" long:"verbose" description:"Print a verbose message to stdout"`
		Delay   int    `short:"d" long:"delay" description:"The seconds for waiting after respawning" default:"60"`
	}

	// set the logger
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Ltime)

	// parse the arguments
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	// get the Tail struct for monitoring the last part of the log
	t, err := tail.TailFile(opts.LogPath, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
	if err != nil {
		logger.Fatalln(err)
	}

	regexFail := regexp.MustCompile(opts.Regex) // compile the Regex pattern
	currentCmd := spawn(opts.CmdPath)           // get the current Cmd struct

	// monitor the log(Tail struct)
	for line := range t.Lines {

		// if the current line contains the pattern that declared by argument
		if regexFail.MatchString(line.Text) {
			// notify it
			logger.Printf("[FAIL] %v -> %v\n", line.Text, opts.Regex)

			// kill the sick process
			pgid, err := syscall.Getpgid(currentCmd.Process.Pid)
			if err == nil {
				syscall.Kill(-pgid, 15)
			}

			// wait to avoid being with flooded with respawning
			logger.Printf("Waiting %v seconds...\n", opts.Delay)
			time.Sleep(time.Second * time.Duration(opts.Delay))

			// respawn normal one
			currentCmd = spawn(opts.CmdPath)

			// if the Verbose flag is set, also print normal lines
		} else if opts.Verbose {
			logger.Println(line.Text)
		}
	}
}
