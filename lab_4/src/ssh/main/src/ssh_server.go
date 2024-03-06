package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"github.com/gliderlabs/ssh"
	"time"
)

func main(){
	ssh.Handle(func(s ssh.Session) {
		command := s.Command()
		if len(command) == 0 {
			fmt.Fprintln(s, "No command provided.")
			return
		}

		switch command[0] {
		case "mkdir":
			dirName := command[1]
			err := os.Mkdir(dirName, 0755)
			if err != nil {
				fmt.Fprintln(s, "Failed to create directory:", err)
				break
			}
			fmt.Fprintln(s, "Directory", dirName, "created successfully.")
		case "rmdir":
			dirName := command[1]
			err := os.RemoveAll(dirName)
			if err != nil {
				fmt.Fprintln(s, "Failed to remove directory:", err)
				break
			}
			fmt.Fprintln(s, "Directory", dirName, "removed successfully.")
		case "ls":
			files, err := ioutil.ReadDir(".")
			if err != nil {
				fmt.Fprintln(s, "Failed to list directory:", err)
				break
			}
			for _, file := range files {
				fmt.Fprintln(s, file.Name())
			}
		case "mv":
			src := command[1]
			dest := command[2]
			err := os.Rename(src, dest)
			if err != nil {
				fmt.Fprintln(s, "Failed to move file:", err)
				break
			}
			fmt.Fprintln(s, "File", src, "moved to", dest)
		case "rm":
			fileName := command[1]
			err := os.Remove(fileName)
			if err != nil {
				fmt.Fprintln(s, "Failed to remove file:", err)
				break
			}
			fmt.Fprintln(s, "File", fileName, "removed successfully.")
		case "ping":
			if len(command) < 2 {
			  fmt.Println("Site name not provided.\n")
			  break
			}
	
	
			cmd := exec.Command("ping", command[1], command[2], command[3]) // Ограничиваем выполнение 5 пакетами
			stdout, err := cmd.StdoutPipe()
			if err != nil {
			  fmt.Sprintf("Failed to create StdoutPipe for ping command: %s\n", err)
			}
	
			if err := cmd.Start(); err != nil {
			  fmt.Sprintf("Failed to start ping command: %s\n", err)
			}
	
			output, err := ioutil.ReadAll(stdout)
			if err != nil {
			  fmt.Sprintf("Failed to read ping output: %s\n", err)
			}
	
			done := make(chan error, 1)
			go func() { done <- cmd.Wait() }()
			select {
			case <-time.After(10 * time.Second):
			  if err := cmd.Process.Kill(); err != nil {
				fmt.Sprintf("Failed to kill ping process: %s\n", err)
			  }
			  fmt.Println("Ping command timed out.\n")
			case err := <-done:
			  if err != nil {
				fmt.Sprintf("Ping command finished with error: %s\n", err)
			  }
			  fmt.Fprintln(s, string(output))
			}
	
		default:
			fmt.Fprintln(s, "Unknown command.")
		}
	})

	fmt.Print("IP: ")
	ip := ""
	fmt.Scan(&ip)
	fmt.Print("Port: ")
	port := ""
	fmt.Scan(&port)
	err := ssh.ListenAndServe(ip + ":" + port , nil)
	if err != nil {
		fmt.Println("Failed to start SSH server:", err)
		return
	}
}
