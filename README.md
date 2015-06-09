# hostsup

A tool for easily editing your hosts file.

## Installation

```sh
$ go get github.com/damonkelley/hostsup
```

## Usage
```console
NAME:
   hostsup - A tool to easily manage your hosts file.

USAGE:
   hostsup [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   add          Add a hosts entry
   rm           Remove a hosts entry
   ls           List hosts entries
   clean        Remove all hosts entries added by hostsup
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version
```

## Example
```console
$ sudo hostsup add 192.168.0.1 myhost.dev
INFO[0000] Added "192.168.0.1 myhost.dev" to /etc/hosts.

$ cat /etc/hosts
##
# Host Database
#
# localhost is used to configure the loopback interface
# when the system is booting.  Do not change this entry.
##
127.0.0.1       localhost
255.255.255.255 broadcasthost
::1             localhost

192.168.0.1     myhost.dev      # HOSTSUP aa60767294065d50d667a88ad2275888

$ hostsup ls
HOSTNAME     IP
myhost.dev   192.168.0.1

$ sudo hostsup rm myhost.dev
INFO[0000] Removed "192.168.0.1 myhost.dev" from /etc/hosts.

## Warning
This tool edits an important file on your machine. Use at your own risk.
