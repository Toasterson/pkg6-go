package depotd

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/toasterson/pkg6-go/repo"
)

type DepotServer struct {
	*echo.Echo
	repo.Repository
}

func NewDepotServer(repository string) (depot *DepotServer) {
	if r, err := repo.NewRepo(repository); err != nil {
		panic(fmt.Errorf("can not instantiate depotserver: %s", err))
	} else {
		depot = &DepotServer{
			Echo:       echo.New(),
			Repository: r,
		}
	}
	depot.setupMiddleware()
	depot.mountPublicEndpoints()
	return
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
