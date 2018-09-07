package metadata

import "path/filepath"

type Publisher struct {
	Catalog PackageCatalog
}

func NewPublisher(repoPath string) *Publisher {
	pub := &Publisher{}
	pub.Catalog = NewV1Catalog(filepath.Join(repoPath, "publisher", "catalog"))
	return pub
}
