package action

import (
	"encoding/json"
	"testing"
)

const (
	dependJSONString = `{"fmri":"pkg:/system/library@0.5.11-2016.0.0.15685","type":"require"}`
	dependActString1 = "depend fmri=pkg:/system/library@0.5.11-2016.0.0.15685 type=require"
)

func TestDepenActionLoad(t *testing.T) {
	var dep DependAction
	if err := json.Unmarshal([]byte(dependJSONString), &dep); err != nil {
		t.Fatalf("could not unmarshal depend Action: %s", err)
	}
	if dep.Type != "require" {
		t.Fatal("type not correct")
	}
	if dep.FMRI != "pkg:/system/library@0.5.11-2016.0.0.15685" {
		t.Fatal("fmri not correct")
	}
}

func TestDependActionFromActionString(t *testing.T) {
	action := DependAction{}
	action.FromActionString(dependActString1)
	if action.FMRI == "fmri=pkg:/system/library@0.5.11-2016.0.0.15685" {
		t.Errorf(errorformat, "fmri=pkg:/system/library@0.5.11-2016.0.0.15685", action.FMRI)
	}

}
