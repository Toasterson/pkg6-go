package depotd

import (
	"fmt"
	"git.wegmueller.it/toasterson/glog"
	"github.com/toasterson/pkg6-go/metadata"
	"github.com/toasterson/pkg6-go/repo"
	"os"
	"path/filepath"
)

type MirrorConfig struct {
	Name       string
	BaseUrl    string
	LocalPath  string
	LocalRepo  *repo.FileRepo `yaml:"-" json:"-" toml:"-"`
	RemoteRepo *repo.HttpRepo `yaml:"-" json:"-" tome:"-"`
}

func NewMirrorConfig(name, url string) *MirrorConfig {
	return &MirrorConfig{
		Name:    name,
		BaseUrl: url,
	}
}

// A function designed to be run in the background the grab some files and be able to start mirroring.
func (m *MirrorConfig) InitRepos() {
	m.RemoteRepo = &repo.HttpRepo{
		Name:    m.Name,
		BaseUrl: m.BaseUrl,
	}
	if err := m.RemoteRepo.Load(); err != nil {
		glog.Emergf("could not initialize remote repo for mirror %s: %s", m.Name, err)
	}
	m.LocalRepo = &repo.FileRepo{
		Path:       m.LocalPath,
		IsMirror:   true,
		Publishers: m.RemoteRepo.Publishers,
	}
	if m.LocalRepo.Exists() {
		if err := m.LocalRepo.Load(); err != nil {
			glog.Emergf("could not initialize local repo for mirror %s: %s", m.Name, err)
		}
	} else {
		if err := m.LocalRepo.Create(); err != nil {
			glog.Emergf("could not initialize local repo for mirror %s: %s", m.Name, err)
		}
	}
}

func (d *DepotServer) AddMirror(mirror *MirrorConfig) error {
	mirror.LocalPath = filepath.Join(d.HomePath, "repositories", mirror.Name)
	d.State.Mirrors[mirror.Name] = *mirror
	d.saveStateFile()
	d.Mirrors[mirror.Name] = mirror
	d.Mirrors[mirror.Name].InitRepos()
	return nil
}

//The main Worker Function Grabs files from the Remote Repo and Puts them into the local one.
func (d *DepotServer) Mirror(name string) error {
	if mirror, ok := d.Mirrors[name]; !ok {
		return fmt.Errorf("mirror with name %s does not exist", name)
	} else {
		var cat *metadata.V1Catalog
		var err error
		publishers := mirror.RemoteRepo.GetPublishers()
		for _, pub := range publishers {
			// Grab Catalog from Remote or Cache
			if cat, err = mirror.RemoteRepo.GetCatalog(pub); err != nil {
				return err
			}
			catalogFiles := []string{"catalog.attrs"}
			for part := range cat.Parts {
				catalogFiles = append(catalogFiles, part)
			}
			for update := range cat.Updates {
				catalogFiles = append(catalogFiles, update)
			}
			// Add all catalogs to local Repository
			for _, catF := range catalogFiles {
				var f *os.File
				f, err = mirror.RemoteRepo.GetCatalogFile(pub, catF)
				if err != nil {
					f.Close()
					return err
				}
				err = mirror.LocalRepo.AddCatalogFile(pub, catF, f)
				if err != nil {
					f.Close()
					return err
				}
				f.Close()
			}
			//TODO do the same for all manifests
			//TODO do the same for all files
		}
		return nil
	}
}
