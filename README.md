# ssh-repl

[![License](https://img.shields.io/badge/License-MIT%20License-blue.svg)](https://opensource.org/licenses/MIT)

A simple repl program for ssh

## Use

First, set config
```
opts:
 network: tcp
 user: yourname
 addr: youraddr
 port: yourport
 kpath: yourfilepath
 key: yourkey # if set key, setting kpath isn't required. On the contrary, kpath is also so.
```

Second, launch program
``` shell
:> ac fast -path config // fastest way to create ssh session

xxx@VM-xxx:~$ :quit // quit current session, back to repl

:> quit // save the last enforced commands in record.txt
```
