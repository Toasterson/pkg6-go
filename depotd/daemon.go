package depotd

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
	"github.com/toasterson/pkg6-go/repo"
	"os"
	"path/filepath"
)

func init() {
	viper.SetDefault("depot_home", "./")
	viper.SetDefault("socket_path", "/var/run/depotd.sock")
	viper.SetDefault("config_path", "depotd.yml")
	viper.SetDefault("cache_dir", "cache")
}

type DepotServer struct {
	*echo.Echo
	HomePath string
	State    DepotState
	Repos    []repo.Repository
	Mirrors  map[string]*MirrorConfig
}

type DepotState struct {
	Repos   []string                `json:"repos"`
	Mirrors map[string]MirrorConfig `json:"mirrors"`
}

func NewDepotServer() (depot *DepotServer) {
	//Load Environment Variables
	viper.AutomaticEnv()
	//Set config Path
	viper.SetConfigFile(viper.GetString("config_path"))
	//Now read the Config. If it exists it will just be loaded from disk
	viper.ReadInConfig()
	depot = &DepotServer{
		Echo: echo.New(),
	}
	depot.HomePath = os.ExpandEnv(viper.GetString("depot_home"))
	depot.loadRepositoriesStartup()
	depot.setupMiddleware()
	depot.mountPublicEndpoints()
	return
}

func (d *DepotServer) loadRepositoriesStartup() {
	d.loadStateFile()
	d.Repos = make([]repo.Repository, 0)
	d.Mirrors = make(map[string]*MirrorConfig)
	// Load the Specific Repository Configurations
	for _, repository := range d.State.Repos {
		r, err := repo.NewRepo(repository)
		if err != nil {
			panic(fmt.Errorf("could not initialize repositories aborting startup: %s", err))
		}
		r.Load()
		d.Repos = append(d.Repos, r)
	}

	for name, mirror := range d.State.Mirrors {
		d.Mirrors[name] = &mirror
		d.Mirrors[name].InitRepos()
	}
}

func (d *DepotServer) loadStateFile() {
	if f, err := os.Open(filepath.Join(d.HomePath, "state.json")); err != nil {
		d.State = DepotState{
			Repos:   make([]string, 0),
			Mirrors: make(map[string]MirrorConfig),
		}
		d.saveStateFile()
	} else {
		defer f.Close()
		if err := json.NewDecoder(f).Decode(&d.State); err != nil {
			panic(err)
		}
	}
}

func (d *DepotServer) saveStateFile() {
	if f, err := os.Create(filepath.Join(d.HomePath, "state.json")); err != nil {
		panic(err)
	} else {
		defer f.Close()
		if err := json.NewEncoder(f).Encode(d.State); err != nil {
			panic(err)
		}
	}
}

func (d *DepotServer) mountPublicEndpoints() {
	d.GET("/versions/0/", d.handleVersionsV0)
	d.GET("/:publisher/search/0/:query", d.handleSearchV0)
	d.GET("/:publisher/search/1/:casesensitive_:returntype_:maxreturn_:startreturn_:query", d.handleSearchV1)
	d.GET("/:publisher/catalog/0", d.handleCatalogV0)
	d.GET("/:publisher/catalog/1/:catalog", d.handleCatalogV1)
	d.GET("/:publisher/manifest/0/:manifest", d.handleManifestV0)
	d.GET("/:publisher/file/1/:fileSHA1", d.handleFileV1)
	d.GET("/:publisher/p5i/0/:fmri", d.handleP5IV0)
}

func (d *DepotServer) setupMiddleware() {
	d.Use(middleware.Logger())
	d.Use(middleware.Recover())
}
