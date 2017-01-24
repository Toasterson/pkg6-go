package main

import (
	"github.com/toasterson/pkg6-go/catalog"
	"fmt"
)

var catalogBasePath = "./sample_data/repo/publisher/userland/catalog"

func main() {
	cat := catalog.Catalog{}
	cat.LoadFromV1(catalogBasePath)
	fmt.Println(cat)
}
