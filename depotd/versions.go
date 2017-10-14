package main

import (
	"net/http"
	"fmt"
)

type version_t map[string][]int

var versions version_t

func init() {
	versions = version_t{
		"info":      []int{0},
		"index":     []int{0},
		"search":    []int{0, 1},
		"p5i":       []int{0},
		"publisher": []int{0, 1},
		"admin":     []int{0},
		"versions":  []int{0},
		"catalog":   []int{0, 1},
		"filelist":  []int{0},
		"manifest":  []int{0},
		"add":       []int{0},
		"status":    []int{0},
		"file":      []int{0, 1},
		"abandon":   []int{0},
		"close":     []int{0},
		"open":      []int{0},
		"append":    []int{0},
	}
}

func handleVersionsV0(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pkg-server 6.0+git"))
	for k,v := range versions{
		w.Write([]byte(fmt.Sprintf("%s %v\n", k, v)))
	}
}
