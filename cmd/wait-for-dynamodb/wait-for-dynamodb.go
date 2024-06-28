package main

import (
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "time"
)

func main() {
    url := "http://dynamodb-local:8000"
    timeout := time.Minute * 2
    ticker := time.NewTicker(time.Second * 2)
    defer ticker.Stop()

    start := time.Now()
    for time.Since(start) < timeout {
        _, err := http.Get(url)
        if err == nil {
            fmt.Println("DynamoDB is available - executing command")
            startMainApp()
            return
        }

        fmt.Println("Waiting for DynamoDB to be available...")
        <-ticker.C
    }

    fmt.Println("Timed out waiting for DynamoDB")
    os.Exit(1)
}

func startMainApp() {
    cmd := exec.Command("/app/voting-service")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    err := cmd.Run()
    if err != nil {
        fmt.Printf("Failed to start main app: %v\n", err)
        os.Exit(1)
    }
    // Wait for the command to finish
    err = cmd.Wait()
    if err != nil {
        fmt.Printf("Main app exited with error: %v\n", err)
        os.Exit(1)
    }
}
