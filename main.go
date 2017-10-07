package main

import (
	"github.com/toasterson/pkg6-go/pkg"
	"fmt"
)

var repoPath = "/home/toast/workspace/go/src/github.com/toasterson/pkg6-go/sample_data/repo"

func main() {
	packageInfo := pkg.PackageInfo{}
	packageInfo.SetFmri("pkg://userland/library/desktop/mate/libmatemixer@1.16.0,5.11-2016.1.1.0:20161224T161749Z")
	packageInfo.ReadManifest(repoPath)
	if err := packageInfo.Save(repoPath); err != nil {
		fmt.Printf("Cannot Save Package: %s", err)
	}
}
