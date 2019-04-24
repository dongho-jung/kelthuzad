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
3. **DONE**. ex) `./kelthuzad -c ./spawn.sh -p ./postgres_proxy.log -r 'timeout|error' -v` Give it a shot!
(the content of *spawn.sh* could be like this. `kubectl port-forward --namespace postgresql svc/postgresql-postgresql 3****:3**** --address 175.*.*.* &> postgres_proxy.log`)
