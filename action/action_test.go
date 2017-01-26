package action

import (
	"testing"
)

const action_string1 = "set name=info.source-url value=http://www.pgpool.net/download.php?f=pgpool-II-3.3.1.tar.gz"
const action_string2 = "set name=variant.opensolaris.zone value=global value=nonglobal"
const action_string3 = "set name=variant.arch value=i386"
const action_string4 = "set name=pkg.summary value=\"Gujarati language support\""
const action_string5 = "set name=pkg.summary value=\\\"'XZ Utils - loss-less file compression application and library.'\\\""
const action_string6 = "set name=pkg.summary value=\\\"provided mouse accessibility enhancements\\\""
const action_string7 = "set name=info.upstream value=X.Org Foundation"
const action_string8 = "set name=pkg.description value=Latvian language support's extra files"
const action_string9 = "set name=illumos-gate.info.git-remote value=git://github.com/illumos/illumos-gate.git"

func TestAttributeAction_FromActionString1(t *testing.T) {
	action := AttributeAction{}
	action.FromActionString(action_string1)
	if action.Name != "info.source-url" {
		t.Errorf("Expected %s got %s", "info.source-url", action.Name)
	}
	if action.Values[0] != "http://www.pgpool.net/download.php?f=pgpool-II-3.3.1.tar.gz" {
		t.Errorf("Expected %s got %s", "http://www.pgpool.net/download.php?f=pgpool-II-3.3.1.tar.gz", action.Values[0])
	}
}

func TestAttributeAction_FromActionString2(t *testing.T) {
	action := AttributeAction{}
	action.FromActionString(action_string2)
	if action.Name != "variant.opensolaris.zone" {
		t.Errorf("Expected %s got %s", "variant.opensolaris.zone", action.Name)
	}
	if action.Values[0] != "global" {
		t.Errorf("Expected %s got %s", "global", action.Values[0])
	}
	if action.Values[1] != "nonglobal" {
		t.Errorf("Expected %s got %s", "nonglobal", action.Values[1])
	}
}

func TestAttributeAction_FromActionString3(t *testing.T) {
	action := AttributeAction{}
	action.FromActionString(action_string3)
	if action.Name != "variant.arch" {
		t.Errorf("Expected %s got %s", "variant.arch", action.Name)
	}
	if action.Values[0] != "i386" {
		t.Errorf("Expected %s got %s", "i386", action.Values[0])
	}
}

func TestAttributeAction_FromActionString4(t *testing.T) {
	action := AttributeAction{}
	action.FromActionString(action_string4)
	if action.Name != "pkg.summary" {
		t.Errorf("Expected %s got %s", "pkg.summary", action.Name)
	}
	if action.Values[0] != "Gujarati language support" {
		t.Errorf("Expected %s got %s", "Gujarati language support", action.Values[0])
	}
}

func TestAttributeAction_FromActionString5(t *testing.T) {
	action := AttributeAction{}
	action.FromActionString(action_string5)
	if action.Name != "pkg.summary" {
		t.Errorf("Expected %s got %s", "pkg.summary", action.Name)
	}
	if action.Values[0] != "XZ Utils - loss-less file compression application and library." {
		t.Errorf("Expected %s got %s", "XZ Utils - loss-less file compression application and library.", action.Values[0])
	}
}

func TestAttributeAction_FromActionString6(t *testing.T) {
	action := AttributeAction{}
	action.FromActionString(action_string6)
	if action.Name != "pkg.summary" {
		t.Errorf("Expected %s got %s", "pkg.summary", action.Name)
	}
	if action.Values[0] != "provided mouse accessibility enhancements" {
		t.Errorf("Expected %s got %s", "provided mouse accessibility enhancements", action.Values[0])
	}
}

func TestAttributeAction_FromActionString7(t *testing.T) {
	action := AttributeAction{}
	action.FromActionString(action_string7)
	if action.Name != "info.upstream" {
		t.Errorf("Expected %s got %s", "info.upstream", action.Name)
	}
	if action.Values[0] != "X.Org Foundation" {
		t.Errorf("Expected %s got %s", "X.Org Foundation", action.Values[0])
	}
}

func TestAttributeAction_FromActionString8(t *testing.T) {
	action := AttributeAction{}
	action.FromActionString(action_string8)
	if action.Name != "pkg.description" {
		t.Errorf("Expected %s got %s", "pkg.description", action.Name)
	}
	if action.Values[0] != "Latvian language support's extra files" {
		t.Errorf("Expected %s got %s", "Latvian language support's extra files", action.Values[0])
	}
}

func TestAttributeAction_FromActionString9(t *testing.T) {
	action := AttributeAction{}
	action.FromActionString(action_string9)
	if action.Name != "illumos-gate.info.git-remote" {
		t.Errorf("Expected %s got %s", "illumos-gate.info.git-remote", action.Name)
	}
	if action.Values[0] != "git://github.com/illumos/illumos-gate.git" {
		t.Errorf("Expected %s got %s", "git://github.com/illumos/illumos-gate.git", action.Values[0])
	}
}