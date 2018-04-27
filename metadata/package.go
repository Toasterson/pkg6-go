package metadata

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/toasterson/pkg6-go/action"
	"os"
	"strings"
	"time"
)

type PackageInfo struct {
	Publisher        string                   `json:"publisher,-"`
	Name             string                   `json:"name"`
	ComponentVersion Version                  `json:"version"`
	BuildVersion     string                   `json:"build"`
	BranchVersion    string                   `json:"branch"`
	PackagingDate    time.Time                `json:"packaging_date"`
	SignatureSHA1    string                   `json:"signature-sha-1"`
	Summary          string                   `json:"summary"`
	Description      string                   `json:"description"`
	Classification   []string                 `json:"classification"`
	Attributes       []action.AttributeAction `json:"attributes"`
	Dependencies     []action.DependAction    `json:"dependencies"`
	Directories      []action.DirectoryAction `json:"directories"`
	Files            []action.FileAction      `json:"files"`
	Links            []action.LinkAction      `json:"links"`
	Licenses         []action.LicenseAction   `json:"licenses"`
}

func (p *PackageInfo) SetFmri(fmri string) error {
	if !strings.HasPrefix(fmri, "pkg://") {
		return errors.New("invalid FMRI given")
	}
	mapFMRI := SplitFmri(fmri)
	p.Publisher = mapFMRI["publisher"]
	p.Name = mapFMRI["name"]
	p.ComponentVersion.FromVersionString(mapFMRI["version"])
	p.BuildVersion = mapFMRI["build_release"]
	p.BranchVersion = mapFMRI["branch"]
	p.SetPackagingDate(mapFMRI["packaging_date"])
	return nil
}

func SplitFmri(fmri string) map[string]string {
	var mapFMRI = map[string]string{}
	tmpFMRI := fmri
	if strings.HasPrefix(tmpFMRI, "pkg://") {
		tmpFMRI = strings.Replace(tmpFMRI, "pkg://", "", 1)
		tmp := strings.SplitN(tmpFMRI, "/", 2)
		mapFMRI["publisher"] = tmp[0]
		tmpFMRI = tmp[1]
	} else {
		tmpFMRI = strings.Replace(tmpFMRI, "pkg:/", "", -1)
	}
	tmpName := strings.SplitN(tmpFMRI, "@", 2)
	mapFMRI["name"] = tmpName[0]
	tmpFMRI = tmpName[1]

	tmpVersion := strings.SplitN(tmpFMRI, ",", 2)
	mapFMRI["version"] = tmpVersion[0]
	tmpFMRI = tmpVersion[1]

	tmpBuild := strings.SplitN(tmpFMRI, "-", 2)
	mapFMRI["build_release"] = tmpBuild[0]
	tmpFMRI = tmpBuild[1]

	tmpBranch := strings.SplitN(tmpFMRI, ":", 2)
	mapFMRI["branch"] = tmpBranch[0]
	tmpFMRI = tmpBranch[1]

	mapFMRI["packaging_date"] = tmpFMRI
	return mapFMRI
}

func (p *PackageInfo) SetPackagingDate(datestring string) {
	t, err := time.Parse("20060102T150405Z", datestring)
	if err == nil {
		p.PackagingDate = t
	}
}

func (p *PackageInfo) FromMap(packMap map[string]interface{}) {
	for key, value := range packMap {
		switch key {
		case "version":
			{
				p.ComponentVersion, p.BuildVersion, p.BranchVersion, p.PackagingDate = splitFMRIVersion(value.(string))
			}
		case "signature-sha-1":
			{
				p.SignatureSHA1 = value.(string)
			}
		case "actions":
			{
				for _, loopVal := range value.([]interface{}) {
					actString := loopVal.(string)
					if strings.Contains(actString, "set") {
						attr := action.AttributeAction{}
						attr.FromActionString(actString)
						p.Attributes = append(p.Attributes, attr)
					} else if strings.Contains(actString, "depend") {
						dep := action.DependAction{}
						dep.FromActionString(actString)
						p.Dependencies = append(p.Dependencies, dep)
					}
				}
			}
		default:
			{

			}
		}
	}
}

func (p *PackageInfo) CompareVersion(p2 PackageInfo) string {
	if p.ComponentVersion.LesserThan(p2.ComponentVersion) {
		if p.PackagingDate.Before(p2.PackagingDate) {
			return "lesser"
		}
		return "bigger"
	} else if p2.ComponentVersion.LesserThan(p.ComponentVersion) {
		if p2.PackagingDate.Before(p.PackagingDate) {
			return "bigger"
		}
		return "lesser"
	} else {
		return "equals"
	}
}

func (p *PackageInfo) Merge(p2 *PackageInfo) {
	p.Summary = p2.Summary
	p.Description = p2.Description
	for _, val := range p2.Attributes {
		p.Attributes = append(p.Attributes, val)
	}
	for _, val := range p2.Dependencies {
		p.Dependencies = append(p.Dependencies, val)
	}
}

func (p *PackageInfo) ReadManifest(location string) error {
	path := location + "/" + FMRI2Unicode(p.GetFMRI())
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "set") {
			attr := action.AttributeAction{}
			attr.FromActionString(text)
			p.Attributes = append(p.Attributes, attr)
		} else if strings.HasPrefix(text, "depend") {
			dep := action.DependAction{}
			dep.FromActionString(text)
			p.Dependencies = append(p.Dependencies, dep)
		} else if strings.HasPrefix(text, "file") {
			fileAction := action.FileAction{}
			fileAction.FromActionString(text)
			p.Files = append(p.Files, fileAction)
		} else if strings.HasPrefix(text, "license") {
			lic := action.LicenseAction{}
			lic.FromActionString(text)
			p.Licenses = append(p.Licenses, lic)
		} else if strings.HasPrefix(text, "link") {
			linkAction := action.LinkAction{}
			linkAction.FromActionString(text)
			p.Links = append(p.Links, linkAction)
		} else {
			return fmt.Errorf("uknown Action in %p: %a", p.Name, text)
		}
	}
	return nil
}

func (p *PackageInfo) GetFMRI() string {
	return p.Name + "@" + p.ComponentVersion.ToVersionString() + "," + p.BuildVersion + "-" + p.BranchVersion + ":" + p.PackagingDate.Format("20060102T150405Z")
}

func (p *PackageInfo) Save(location string) error {
	b, err := json.Marshal(p)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot Marshal %s", p.Name))
	}
	path := location + "/" + FMRI2Unicode(p.GetFMRI()) + ".json"
	file, ferr := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if ferr != nil {
		return errors.New(throwError("Saving", p.GetFMRI(), ferr.Error()))
	}
	defer file.Close()
	if _, werr := file.Write(b); werr != nil {
		return errors.New(throwError("Saving", p.GetFMRI(), werr.Error()))
	}
	return nil
}

func (p *PackageInfo) Load(location string) error {
	path := location + "/" + FMRI2Unicode(p.GetFMRI()) + ".json"
	file, ferr := os.OpenFile(path, os.O_RDONLY, 0666)
	if ferr != nil {
		return errors.New(throwError("Loading", p.GetFMRI(), ferr.Error()))
	}
	defer file.Close()
	var b = []byte{}
	_, rerr := file.Read(b)
	if rerr != nil {
		return errors.New(throwError("Loading", p.GetFMRI(), rerr.Error()))
	}
	return json.Unmarshal(b, p)
}

func (p PackageInfo) WriteManifest() string {
	buff := bytes.NewBufferString("")
	for _, attr := range p.Attributes {
		buff.WriteString(fmt.Sprintf("set name=%s", attr.Name))
		for _, val := range attr.Values {
			buff.WriteString(fmt.Sprintf(" value=\"%s\"", val))
		}
		if len(attr.Optionals) > 0 {
			for optKey, optValue := range attr.Optionals {
				buff.WriteString(fmt.Sprintf(" %s=\"%s\"", optKey, optValue))
			}
		}
		buff.WriteString("\n")
	}
	for _, dir := range p.Directories {
		buff.WriteString(fmt.Sprintf("dir group=%s mode=%s owner=%s path=\"%s\"", dir.Group, dir.Mode, dir.Owner, dir.Path))
		if len(dir.Facets) > 0 {
			for key, val := range dir.Facets {
				buff.WriteString(fmt.Sprintf("%s=%s", key, val))
			}
		}
		buff.WriteString("\n")
	}
	/*
		for _, file := range p.Files {

		}
		for _, link := range p.Links {

		}
		for _, license := range p.Licenses {

		}
		for _, dep := range p.Dependencies {

		}
	*/
	return buff.String()
}
