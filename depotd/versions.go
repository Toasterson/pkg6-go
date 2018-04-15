package depotd

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

type versionType map[string][]int

var versions versionType

func init() {
	versions = versionType{
		"info":      []int{0},
		"index":     []int{0},
		"search":    []int{0, 1},
		"p5i":       []int{0},
		"publisher": []int{0, 1},
		"admin":     []int{0},
		"versions":  []int{0},
		"catalog":   []int{0, 1},
		"filelist":  []int{0},
		"manifest":  []int{0, 1},
		"add":       []int{0},
		"status":    []int{0},
		"file":      []int{0, 1},
		"abandon":   []int{0},
		"close":     []int{0},
		"open":      []int{0},
		"append":    []int{0},
	}
}

func (d *DepotServer) handleVersionsV0(c echo.Context) error {
	output := "pkg-server 6.0+git"
	for k, v := range versions {
		output += fmt.Sprintf("\n%s %v\n", k, v)
	}
	return c.String(http.StatusOK, output)
}
