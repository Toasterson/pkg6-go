package main

import (
	"github.com/toasterson/pkg6-go/repo"
)

var repoPath = "./sample_data/repo"

func main() {
	repository, _ := repo.NewRepo("file://"+repoPath)
	repository.Load()
	repository.Upgrade()
}
