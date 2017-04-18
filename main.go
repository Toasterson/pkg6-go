package main

import (
	"github.com/toasterson/pkg6-go/packageinfo"
	"github.com/toasterson/pkg6-go/util"
)

var repoPath = "/home/toast/workspace/go/src/github.com/toasterson/pkg6-go/sample_data/repo"

func main() {
	pkg := packageinfo.PackageInfo{}
	pkg.SetFmri("pkg://userland/library/desktop/mate/libmatemixer@1.16.0,5.11-2016.1.1.0:20161224T161749Z")
	pkg.ReadManifest(repoPath)
	util.Error(pkg.Save(repoPath), "Saving Package:")
}
