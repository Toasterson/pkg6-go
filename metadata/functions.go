package metadata

import (
	"fmt"
	"strings"
)

func FromFMRI(fmri string) *PackageInfo {
	pkg := &PackageInfo{}
	pkg.SetFmri(fmri)
	return pkg
}

func throwError(action string, pkg string, err string) string {
	return fmt.Sprintf("Error %s %s: %s", action, pkg, err)
}

func FMRI2Unicode(fmri string) string {
	fmri = strings.Replace(fmri, "/", "%2F", -1)
	fmri = strings.Replace(fmri, ",", "%2C", -1)
	fmri = strings.Replace(fmri, ":", "%3A", -1)
	fmri = strings.Replace(fmri, "@", "/", -1)
	return fmri
}

func Unicode2FMRI(unicode string) string {
	//TODO Figure Out why we do this translation again
	unicode = strings.Replace(unicode, "/", "@", -1)
	unicode = strings.Replace(unicode, "%2F", "/", -1)
	unicode = strings.Replace(unicode, "%2C", ",", -1)
	unicode = strings.Replace(unicode, "%3A", ":", -1)
	return unicode
}
