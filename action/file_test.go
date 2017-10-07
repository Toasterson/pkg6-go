package action

import "testing"

const (
	fileActionString1 = "file 84ecd47ef66bb9be2f7cce732e4ec67af6c473ff path=usr/share/caja/browser.xml pkg.size=3988 pkg.csize=1145 chash=1eb0bf5396f154ae9646be477a2ad990204f4da2 group=bin mode=0444 owner=root pkg.depend.bypass-generate=usr/lib(.*)/libpq.so.*"
)

func TestFileActionFromActionString(t *testing.T) {
	fact := FileAction{}
	fact.FromActionString(fileActionString1)
	if fact.Sha1 != "84ecd47ef66bb9be2f7cce732e4ec67af6c473ff" {
		t.Fatal("sha1 mismatch")
	}
	if fact.Group != "bin" {
		t.Fatal("group mismatch")
	}
	if fact.Path != "usr/share/caja/browser.xml" {
		t.Fatal("path mismatch")
	}
	if fact.Size != 3988 {
		t.Fatal("size mismatch")
	}
	if fact.Csize != 1145 {
		t.Fatal("csize mismatch")
	}
}
