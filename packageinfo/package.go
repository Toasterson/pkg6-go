package packageinfo

import (
	"time"
	"github.com/toasterson/pkg6-go/action"
	"strings"
)

type PackageInfo struct {
	Name string `json:"name"`
	ComponentVersion Version `json:"version"`
	BuildVersion Version `json:"build"`
	BranchVersion string `json:"branch"`
	PackagingDate time.Time `json:"packaging_date"`
	SignatureSHA1 string `json:"signature-sha-1"`
	Summary string `json:"summary"`
	Description string `json:"description"`
	Classification []string `json:"classification"`
	Attributes []action.AttributeAction `json:"attributes"`
	Dependencies []action.DependAction `json:"dependencies"`
}

func (p *PackageInfo) FromMap(packMap map[string]interface{}){
	for key,value := range packMap{
		switch key {
		case "version" : {
			p.ComponentVersion, p.BuildVersion, p.BranchVersion, p.PackagingDate = splitFMRIVersion(value.(string))
		}
		case "signature-sha-1": {
			p.SignatureSHA1 = value.(string)
		}
		case "actions" : {
			for _, loopVal := range value.([]interface{}){
				act_string := loopVal.(string)
				if strings.Contains(act_string, "set"){
					attr := action.AttributeAction{}
					attr.FromActionString(act_string)
					p.Attributes = append(p.Attributes, attr)
				} else if strings.Contains(act_string, "depend"){
					dep := action.DependAction{}
					dep.FromActionString(act_string)
					p.Dependencies = append(p.Dependencies, dep)
				}
			}
		}
		default:{

		}
		}
	}
}

func (p *PackageInfo) CompareVersion(p2 PackageInfo) string{
	if p.ComponentVersion.LesserThan(p2.ComponentVersion) {
		if p.PackagingDate.Before(p2.PackagingDate){
			return "lesser"
		}
		return "bigger"
	} else if p2.ComponentVersion.LesserThan(p.ComponentVersion){
		if p2.PackagingDate.Before(p.PackagingDate){
			return "bigger"
		}
		return "lesser"
	} else {
		return "equals"
	}
}

func (p *PackageInfo) Merge(p2 *PackageInfo){
	p.Summary = p2.Summary
	p.Description = p2.Description
	for _, val := range p2.Attributes{
		p.Attributes = append(p.Attributes, val)
	}
	for _, val := range p2.Dependencies{
		p.Dependencies = append(p.Dependencies, val)
	}
}

