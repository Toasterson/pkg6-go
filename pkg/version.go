package pkg

import (
	"strings"
	"fmt"
	"time"
	"strconv"
)

type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"patch,omitempty"`
	Beta int `json:"beta,omitempty"`
	ReleaseCandidate int `json:"release_candidate,omitempty"`
	VersionLetter string `json:"version_letter,omitempty"`
}

func (v *Version) FromVersionString(versionstring string){
	if strings.Contains(versionstring, "beta"){
		fmt.Sscanf(versionstring, "%d.%d.%d%s-beta%d", &v.Major, &v.Minor, &v.Patch, &v.VersionLetter, &v.Beta)
	} else if strings.Contains(versionstring, "rc"){
		fmt.Sscanf(versionstring, "%d.%d.%d%s-rc%d", &v.Major, &v.Minor, &v.Patch, &v.VersionLetter, &v.ReleaseCandidate)
	} else {
		fmt.Sscanf(versionstring, "%d.%d.%d%s", &v.Major, &v.Minor, &v.Patch, &v.VersionLetter)
	}
}

func (v *Version) LesserThan(v2 Version) bool {
	if v.Major < v2.Major {
		return true
	}
	if v.Minor < v2.Minor {
		return true
	}
	if v.Patch < v2.Patch {
		return true
	}
	return false
}

func (v *Version) ToVersionString() string{
	returnString := strconv.Itoa(v.Major)+"."+strconv.Itoa(v.Minor)+"."+strconv.Itoa(v.Patch)
	if v.VersionLetter != "" {
		returnString += v.VersionLetter
	}
	if v.Beta != 0 {
		returnString += "-beta"+strconv.Itoa(v.Beta)
	} else if v.ReleaseCandidate != 0 {
		returnString += "-rc"+strconv.Itoa(v.ReleaseCandidate)
	}
	return returnString
}

func splitFMRIVersion(FMRIVersion string) (Version, string, string, time.Time) {
	var branch, build string
	component := Version{}
	//Everything up to , is component
	verseppos := strings.Index(FMRIVersion, ",")
	component.FromVersionString(FMRIVersion[0:verseppos])
	FMRIVersion = FMRIVersion[verseppos+1:]
	//Everything up to - is build
	branchseppos := strings.Index(FMRIVersion, "-")
	build = FMRIVersion[0:branchseppos]
	FMRIVersion = FMRIVersion[branchseppos+1:]
	//Everything up to : is branch
	relseppos := strings.Index(FMRIVersion, ":")
	branch = FMRIVersion[0:relseppos]
	FMRIVersion = FMRIVersion[relseppos+1:]
	//The Rest is Packaging Date
	packDate, _ := time.Parse("20060102T150405Z", FMRIVersion)
	return component, build, branch, packDate
}