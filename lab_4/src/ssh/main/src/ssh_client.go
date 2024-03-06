package main

import (
    "bufio"
	"golang.org/x/crypto/ssh"
	"fmt"
    "io"
    "os"
    "log"
)


func main() {
    fmt.Print("IP: ")
    ip := ""
    fmt.Scan(&ip)
    fmt.Print("Port: ")
    port := ""
    fmt.Scan(&port)
    fmt.Print("login: ")
    user := ""
    fmt.Scan(&user)
    fmt.Print("password: ")
    password := ""
    fmt.Scan(&password)
    fmt.Println(user, password, ip + ":" + port)
	sshConfig := &ssh.ClientConfig{
        User: user,
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Auth: []ssh.AuthMethod{
            ssh.Password(password)},
    }

    connection, err := ssh.Dial("tcp",  ip + ":" + port, sshConfig)
    if err != nil {
        log.Println("Failed to dial:", err)
        return
    }
    comand := ""
    for ; comand != "e_exit"; { 

    session, err := connection.NewSession()
    if err != nil {
        log.Println("Failed to create session:", err)
	    return 
    }

    stdin, err := session.StdinPipe()
    if err != nil {
        log.Println("Unable to setup stdin for session:", err)
	    return
    }
    go io.Copy(stdin, os.Stdin)

    stdout, err := session.StdoutPipe()
    if err != nil {
        log.Println("Unable to setup stdout for session:", err)
	    return
    }
    go io.Copy(os.Stdout, stdout)

    stderr, err := session.StderrPipe()
    if err != nil {
        log.Println("Unable to setup stderr for session:", err)
	    return
    }
    go io.Copy(os.Stderr, stderr)
    
    reader := bufio.NewReader(os.Stdin)
    comand, err := reader.ReadString('\n')
    if err != nil {
     log.Fatal(err)
    }
        err = session.Start(comand)
        if err != nil {
            log.Println(err)
            return
        }
    }
}