package main

import (
	"os"

	togo "github.com/wonrax/togo/internal"
)

func main() {
	environment := os.Getenv("ENVIRONMENT")
	togo.Start(environment)
}
