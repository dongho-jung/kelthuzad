package main

import (
	"bufio"
	"github.com/hpcloud/tail"
	"github.com/jessevdk/go-flags"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"syscall"
	"time"
)

// Kelthuzad monitors a log or stdout, kills a sick one and respawns a normal one.
type Kelthuzad struct {
	cmd        *exec.Cmd
	opt        *opts
	pattern    *regexp.Regexp
	stdout     io.ReadCloser
	isSpawning bool
}

// opts have several options for argument parsing.
type opts struct {
	LogPath    string `short:"l" long:"logPath" description:"The path of the log instead of stdout"`
	CmdPath    string `short:"c" long:"commandPath" description:"The path of a file containing command string to respawn the process"`
	RawCommand string `short:"r" long:"rawCommand" description:"The command string to spawn the process"`
	Pattern    string `short:"p" long:"pattern" description:"The regex pattern to detect a failure" required:"true"`
	Quiet      bool   `short:"q" long:"quiet" description:"Suppress the ouputs of process which is monitored"`
	Delay      int    `short:"d" long:"delay" description:"The seconds for waiting after respawning" default:"5"`
}

// New returns initialized Kelthuzad pointer
func New(opt *opts) *Kelthuzad {
	kel := &Kelthuzad{}
	kel.opt = opt
	kel.pattern = regexp.MustCompile(kel.opt.Pattern)
	kel.spawn()

	return kel
}

// spawn executes the command from k.opt.CmdPath and assigns it into k's cmd field.
func (k *Kelthuzad) spawn() {
	k.isSpawning = false

	var cmd *exec.Cmd
	if k.opt.CmdPath != "" {
		cmd = exec.Command(k.opt.CmdPath)
	} else {
		cmd = exec.Command("bash", "-lc", k.opt.RawCommand+" 2>&1")
	}

	if k.opt.LogPath == "" {
		// get the stdout pipe before it starts and assign it into k.stdout to monitor stdout
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatalln("[FATAL] k.spawn stdout", err)
		}

		k.stdout = stdout
	}

	// this block is necessary when killing a subprocess properly
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	go func() {
		err := cmd.Start()
		if err != nil {
			log.Fatalln("[FATAL] k.spawn Start", err)
		}
		log.Printf("[SYSTEM] %v is spawned\n", cmd.Process.Pid)
		cmd.Wait()
		log.Printf("[SYSTEM] %v is done!\n", cmd.Process.Pid)

		time.Sleep(time.Duration(k.opt.Delay+5) * time.Second)
		if k.isSpawning == false {
			k.spawn()
		}
	}()

	// return the created Cmd struct
	k.cmd = cmd
}

// kill kills current k.cmd.
func (k *Kelthuzad) kill() {
	pgid, err := syscall.Getpgid(k.cmd.Process.Pid)
	if err == nil {
		syscall.Kill(-pgid, 15)
	} else {
		log.Println("[SYSTEM] the proecss was alreday terminated", err)
	}
}

// check checks whether the line matches with the k.pattern.
func (k *Kelthuzad) check(line string) {
	// if the line contains the k.pattern
	if k.pattern.MatchString(line) {
		// notify it
		log.Printf("[FAIL] %v -> %v\n", line, k.opt.Pattern)

		// kill the sick one
		k.kill()

		k.isSpawning = true

		// wait to avoid being with flooded with respawning
		log.Printf("[SYSTEM] Waiting %v seconds...\n", k.opt.Delay)
		time.Sleep(time.Second * time.Duration(k.opt.Delay))

		// respawn the normal one
		k.spawn()

		// if the Quiet flag isn't set, also print normal lines
	} else if k.opt.Quiet == false {
		log.Println(line)
	}
}

// monitorLog monitors the specific log with tail and checks any changes whenever log populated.
func (k *Kelthuzad) monitorLog() {
	// get the Tail struct for monitoring the last part of the log
	t, err := tail.TailFile(k.opt.LogPath, tail.Config{Follow: true, Location: &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}})
	if err != nil {
		log.Fatalln("[FATAL] k.monitorLog tail", err)
	}

	// monitor the log
	for line := range t.Lines {
		k.check(line.Text)
	}
}

// monitorStdout monitors the stdout of the process and checks it.
func (k *Kelthuzad) monitorStdout() {
	for {
		// monitor the stdout
		scanner := bufio.NewScanner(k.stdout)
		for scanner.Scan() {
			k.check(scanner.Text())
		}
	}
}

// Monitor monitors appropriate one depending on LogPath option.
func (k *Kelthuzad) Monitor() {
	if k.opt.LogPath != "" {
		log.Println("[SYSTEM] monitoring log...")
		k.monitorLog()
	} else {
		log.Println("[SYSTEM] monitoring stdout...")
		k.monitorStdout()
	}
}

func main() {
	// initialize empty options
	opt := &opts{}

	// set the log flags
	log.SetFlags(log.Ltime | log.LstdFlags)

	// parse the arguments
	_, err := flags.Parse(opt)
	if err != nil {
		os.Exit(1)
	}

	// make sure that one of these options to be specified
	if (opt.CmdPath == "") == (opt.RawCommand == "") {
		log.Fatalln("[FATAL] You must specify one of CmdPath, RawCommand!")
	}

	// get a kelthuzad object
	kel := New(opt)

	// handle an interrupt for terminate children process and itself gracefully
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		log.Print("[SYSTEM] recieved an interrupt, stopping...\n\n")
		kel.kill()
		os.Exit(0)
	}()

	// start monitoring
	kel.Monitor()
}
