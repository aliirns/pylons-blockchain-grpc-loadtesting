package main

import (
	"fmt"
	"time"

	"github.com/myzhan/boomer"
)

var globalBoomer *boomer.Boomer

func Task1() {
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start)

	_, err := getBalance("pylo128m5z84ydpvmngyv5njyjwrwc0guml0kvp7naf", "upylon")
	if err != nil {
		globalBoomer.RecordFailure("get-balance query", "grpc", elapsed.Nanoseconds()/int64(time.Millisecond), "failure")
	}
	globalBoomer.RecordSuccess("get-balance query", "grpc", elapsed.Nanoseconds()/int64(time.Millisecond), int64(10))
}

func Task2() {
	start := time.Now()
	time.Sleep(100 * time.Millisecond)
	elapsed := time.Since(start)

	_, err := createAccount()
	if err != nil {
		globalBoomer.RecordFailure("create-account", "grpc", elapsed.Nanoseconds()/int64(time.Millisecond), "failure")
	}
	globalBoomer.RecordSuccess("create-account", "grpc", elapsed.Nanoseconds()/int64(time.Millisecond), int64(10))
}

func main() {

	// task1 := &boomer.Task{
	// 	Name: "Task 1",
	// 	// The weight is used to distribute goroutines over multiple tasks.
	// 	Weight: 10,
	// 	Fn:     Task1,
	// }

	// numClients := 1000000
	// spawnRate := float64(100)
	// globalBoomer = boomer.NewStandaloneBoomer(numClients, spawnRate)
	// // globalBoomer.AddOutput(boomer.NewConsoleOutput())
	// globalBoomer.Run(task1)

	// cmd := exec.Command("pylonsd", "tx", "pylons", "create-account", "alii", "", "", "--from alii")
	// stdout, err := cmd.Output()
	// if err != nil {
	// 	fmt.Println(err, stdout)
	// 	return
	// }

	// fmt.Print(string(stdout))
	fmt.Println(createAccount())
}
