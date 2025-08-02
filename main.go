package main

import (
"fmt"
"os"

"github.com/jdk-manager/cmd"
)

func main() {
fmt.Println("JDK Manager is starting...") // Temporary debug line
if err := cmd.Execute(); err != nil {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}
}
