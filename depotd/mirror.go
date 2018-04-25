package depotd

import (
	"fmt"
	"git.wegmueller.it/toasterson/glog"
	"github.com/toasterson/pkg6-go/repo"
)

type MirrorConfig struct {
	Name       string
	BaseUrl    string
	LocalPath  string
	LocalRepo  *repo.FileRepo `yaml:"-" json:"-" toml:"-"`
	RemoteRepo *repo.HttpRepo `yaml:"-" json:"-" tome:"-"`
}

// A function designed to be run in the background the grab some files and be able to start mirroring.
func (m *MirrorConfig) InitRepos() {
	m.LocalRepo = &repo.FileRepo{
		Path:     m.LocalPath,
		IsMirror: true,
	}
	if err := m.LocalRepo.Load(); err != nil {
		glog.Emergf("could not initialize local repo for mirror %s: %s", m.Name, err)
	}
	m.RemoteRepo = &repo.HttpRepo{
		Name:    m.Name,
		BaseUrl: m.BaseUrl,
	}
	if err := m.RemoteRepo.Load(); err != nil {
		glog.Emergf("could not initialize remote repo for mirror %s: %s", m.Name, err)
	}
}

func (d *DepotServer) AddMirror(mirror MirrorConfig) error {
	d.State.Mirrors[mirror.Name] = mirror
	d.saveStateFile()
	d.Mirrors[mirror.Name] = &mirror
	d.Mirrors[mirror.Name].InitRepos()
	return nil
}

//The main Worker Function Grabs files from the Remote Repo and Puts them into the local one.
func (d *DepotServer) Mirror(name string) error {
	if mirror, ok := d.Mirrors[name]; !ok {
		return fmt.Errorf("mirror with name %s does not exist", name)
	} else {
		return mirror.RemoteRepo.DownloadAllCatalogs()
	}
}
