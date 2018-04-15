package main

import (
	"github.com/toasterson/pkg6-go/repo"
	"github.com/toasterson/pkg6-go/depotd"
)

var repoPath = "file://./sample_data/repo"

func main() {
	depotd.NewDepotServer(repoPath).Start(":8080")
}
