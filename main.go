package main

import (
	"github.com/toasterson/pkg6-go/depotd"
)

var repoPath = "file://./sample_data/repo"

func main() {
	depot := depotd.NewDepotServer(repoPath)
	if err := depot.Load(); err != nil {
		panic(err)
	}
	depot.Start(":8080")
}
