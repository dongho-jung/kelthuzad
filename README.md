[![Codacy Badge](https://api.codacy.com/project/badge/Grade/753a3a93a96e45149d7e19fb1639fcb7)](https://app.codacy.com/app/0xF4D3C0D3/kelthuzad?utm_source=github.com&utm_medium=referral&utm_content=0xF4D3C0D3/kelthuzad&utm_campaign=Badge_Grade_Dashboard)

## Overview
kelthuzad watches a process. Whenever any errors are detected from its output or a log, then he replaces it with a normal one.

## Origin of his name
the name kelthuzad comes from `Kel' Thuzad` one of the cards in the TCG game Hearthstone.
![image](https://user-images.githubusercontent.com/19762154/56653541-d08e7480-66c8-11e9-9241-dd67a480309f.png)

As you could already get it, he restores an erroneous process.

## When I need him
### TL; DR. You would need him when you want to monitor, pattern detect, respawn a process

In my case, I had a pod that orchestrated by k8s. It communicated with the external network by proxy. However, the container used to go wrong at times for a reason or another. Whenever the container goes wrong, k8s revived it but not the proxy. So what happens? The actual service recovered and running with no problems behind the scenes, but the proxy just referenced the previous wrong one. I had two options, one is making a new stateful set for only the proxy, and another one is monitoring it and replacing automatically. K8s already implemented the former, but it was quite hard to use when I saw. I also found the latter such as 'Immortal', 'Forever', 'Supervisor' but they were also hard to use or had some dependencies.

So I made up my mind to make my own thing. All you need to do is as follows.

## How can I use him
### Basic
1.`./kelthuzad -r 'fallibleCommand foo bar' -p 'error|fail'`

### Use the log
1.**Set the log** which is populated with the output of the target process. If need be, you can make use of redirection for logging.
2.`./kelthuzad -r 'fallibleCommand foo bar' -p 'error|fail' -l <logPath>`

### Use the recipe
1.**Set the recipe** for executing the target process. That recipe could be anything executable like .sh, .exe, etc...
2.`./kelthuzad -c <fallibleRecipePath> -p 'error|fail'`

for example, a recipe could be like as follows:
```sh
#!/bin/bash
 while :
do
    for n in {3..1}; do
        echo "$n"
        sleep 1
    done
    if [ $((RANDOM % 3)) -eq 0 ]; then
        echo 'bye...'
        sleep 99999
    else
        echo 'hello!'
    fi
done
```

## Usage
```sh
Usage:
  kelthuzad [OPTIONS]

Application Options:
  -l, --logPath=     The path of the log instead of stdout
  -c, --commandPath= The path of a file containing command string to respawn the process
  -r, --rawCommand=  The command string to spawn the process
  -p, --pattern=     The regex pattern to detect a failure
  -q, --quiet        Suppress the ouputs of process which is monitored
  -d, --delay=       The seconds for waiting after respawning (default: 5)

Help Options:
  -h, --help         Show this help message
```

## Demo
[![asciicast](https://asciinema.org/a/242769.svg)](https://asciinema.org/a/242769)

## How to build him
-   Linux: GOOS=linux GOARCH=amd64 go build -o kelthuzad_linux_amd64 kelthuzad.go
-   Mac: GOOS=darwin GOARCH=amd64 go build -o kelthuzad_darwin_amd64 kelthuzad.go

## History
### 1.2
-   change flag options
    -   LogPath(p) -> LogPath(l)
    -   Regex(r) -> Pattern(p)
    -   use Quiet(q) instead of Verbose(v)
    -   new flag RawCommand(r) so you don't have to write a script with CmdPath to spawn 
-   support raw command string!
    -   don't have to write a script. if the command is simple enough, you can just pass it by -r 'soSimpleCommand arg0 arg1' 
-   improve logging to identify the source

### 1.1
-   make LogPath optional
-   change default Delay to 5 from 60
-   make the usage utilize object-oriented-programming more
-   New struct, **Kelthuzad**
-   Kelthuzad has only one exported method, Monitor()
-   **Just use `New`** Function that returns initialized kelthuzad pointer
-   **All you need to do is just getting by New() and monitoring it by .Monitor()**
