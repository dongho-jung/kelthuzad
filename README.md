# Overview
kelthuzad monitors a log. If any errors are detected, then he replaces it with normal one.

# Origin of his name
the name kelthuzad comes from `Kel'Thuzad` one of the cards in the TCG game Hearthstone.
![image](https://user-images.githubusercontent.com/19762154/56653541-d08e7480-66c8-11e9-9241-dd67a480309f.png)

As you could already get it, he restores the sick process.

# When I need him?
## TL; DR. You could need him when you want to monitor some outputs of a process and respawn it when the pattern matches one of the outputs.

In my case, I had a pod that ochestrated by k8s. And it communicated with the outer network by proxy. However, the container used to go wrong at times for a reason or another. Whenever the container goes wrong, k8s revived it but not the proxy. So what happens? the actual service was restored and running properly behind scenes, but the proxy just referenced the wrong one. I had two options, one is making a new stateful set for only the proxy and another one is monitoring it and replacing automatically. The former was already implemented by k8s, but it was quite hard to use when I saw. And I also found the latter such as 'Immortal', 'Forever', 'Supervisor', etc... but they were also hard to use or had some dependencies.

So I made up my mind to make my own thing. All you need to do is as follows...

# How can I use him?
1. **Set the log** which is populated with the output of the target process. If need be, you can make use of redirection to logging.
2. **Set the recipe** for executing the target process. That recipe could be anything executable like .sh, .exe, etc...
3. **DONE**. ex) `./kelthuzad -c ./countdown.sh -r bye -v` Give it a shot!

the content of *countdown.sh* could be like as follows.
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

# Usage
```
Usage:
  kelthuzad [OPTIONS]

Application Options:
  -p, --path=    The path of the log
  -c, --command= The path of a command string to respawn the process
  -r, --regex=   The regex pattern to detect a failure
  -v, --verbose  Print a verbose message to stdout
  -d, --delay=   The seconds for waiting after respawning (default: 5)

Help Options:
  -h, --help     Show this help message
```

# Demo
[![asciicast](https://asciinema.org/a/242769.svg)](https://asciinema.org/a/242769)

# History
## 1.1
### Overview
- make LogPath optional
- change default Delay to 5 from 60
- make the usage utilize object-oriented-programming more

### Changed
- New struct, Kelthuzad
- Kelthuzad has only one exported method, Monitor()
- Just use `New` Function returns initialized kelthuzad pointer
- All you need to do is just getting by New() and monitoring it by .Monitor()
