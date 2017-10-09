package main

import (
	"github.com/toasterson/pkg6-go/repo"
)

var repoPath = "/home/toast/workspace/go/src/github.com/toasterson/pkg6-go/sample_data/repo"

func main() {
	repo := repository.FileRepo{}
	repo.Path = repoPath
	repo.Load()
	repo.Upgrade()
}
