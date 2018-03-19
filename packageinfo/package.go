package packageinfo

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/toasterson/pkg6-go/action"
	"github.com/toasterson/pkg6-go/util"
	"os"
	"strings"
	"time"
)

func FromFMRI(fmri string) PackageInfo {
	pkg := PackageInfo{}
	pkg.SetFmri(fmri)
	return pkg
}

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
		return errors.New("Invalid FMRI given")
	}
	mapFMRI := p.SplitFmri(fmri)
	p.Publisher = mapFMRI["publisher"]
	p.Name = mapFMRI["name"]
	p.ComponentVersion.FromVersionString(mapFMRI["version"])
	p.BuildVersion = mapFMRI["build_release"]
	p.BranchVersion = mapFMRI["branch"]
	p.SetPackagingDate(mapFMRI["packaging_date"])
	return nil
}

func (p *PackageInfo) SplitFmri(fmri string) map[string]string {
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
					act_string := loopVal.(string)
					if strings.Contains(act_string, "set") {
						attr := action.AttributeAction{}
						attr.FromActionString(act_string)
						p.Attributes = append(p.Attributes, attr)
					} else if strings.Contains(act_string, "depend") {
						dep := action.DependAction{}
						dep.FromActionString(act_string)
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
	path := location + "/" + FMRI2Unicode(p)
	file, err := os.Open(path)
	defer file.Close()
	util.Error(err, "Opening Manifest")
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
			return errors.New(fmt.Sprintf("Uknown Action in %p: %a", p.Name, text))
		}
	}
	return nil
}

func (p *PackageInfo) getFMRI() string {
	return p.Name + "@" + p.ComponentVersion.ToVersionString() + "," + p.BuildVersion + "-" + p.BranchVersion + ":" + p.PackagingDate.Format("20060102T150405Z")
}

func (p *PackageInfo) DropManifest(location string) error {
	return os.Remove(location + "/" + FMRI2Unicode(p))
}

func (p *PackageInfo) UpgradeFormat(location string) error {
	if err := p.Save(location); err != nil {
		return err
	}
	return p.DropManifest(location)
}

func (p *PackageInfo) Save(location string) error {
	b, err := json.Marshal(p)
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot Marshal %s", p.Name))
	}
	path := location + "/" + FMRI2Unicode(p) + ".json"
	file, ferr := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666)
	if ferr != nil {
		return errors.New(throwError("Saving", p.getFMRI(), ferr.Error()))
	}
	defer file.Close()
	if _, werr := file.Write(b); werr != nil {
		return errors.New(throwError("Saving", p.getFMRI(), werr.Error()))
	}
	return nil
}

func (p *PackageInfo) Load(location string) error {
	path := location + "/" + FMRI2Unicode(p) + ".json"
	file, ferr := os.OpenFile(path, os.O_RDONLY, 0666)
	if ferr != nil {
		return errors.New(throwError("Loading", p.getFMRI(), ferr.Error()))
	}
	defer file.Close()
	var b = []byte{}
	_, rerr := file.Read(b)
	if rerr != nil {
		return errors.New(throwError("Loading", p.getFMRI(), rerr.Error()))
	}
	return json.Unmarshal(b, p)
}

func (p *PackageInfo) WriteManifest() string {

	return ""
}

func throwError(action string, pkg string, err string) string {
	return fmt.Sprintf("Error %s %s: %s", action, pkg, err)
}

func FMRI2Unicode(p *PackageInfo) string {
	fmri := p.getFMRI()
	fmri = strings.Replace(fmri, "/", "%2F", -1)
	fmri = strings.Replace(fmri, ",", "%2C", -1)
	fmri = strings.Replace(fmri, ":", "%3A", -1)
	fmri = strings.Replace(fmri, "@", "/", -1)
	return fmri
}

func Unicode2FMRI(unicode string) string {
	unicode = strings.Replace(unicode, "/", "@", -1)
	unicode = strings.Replace(unicode, "%2F", "/", -1)
	unicode = strings.Replace(unicode, "%2C", ",", -1)
	unicode = strings.Replace(unicode, "%3A", ":", -1)
	return unicode
}
