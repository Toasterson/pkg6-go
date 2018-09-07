package depotd

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/toasterson/pkg6-go/metadata"
	"net/http"
	"strings"
)

func (d *DepotServer) handleCatalogV0(c echo.Context) error {
	return nil
}

func (d *DepotServer) handleCatalogV1(c echo.Context) error {
	catalogPart := c.Param("catalog")
	publisher := c.Param("publisher")
	var cat *metadata.V1Catalog
	var err error
	if cat, err = d.Repos[0].GetCatalog(publisher); err != nil {
		c.NoContent(http.StatusInternalServerError)
		return err
	}
	var content interface{}
	var ok bool
	if catalogPart == "catalog.attrs" {
		return c.JSON(200, cat)
	} else if strings.Contains(catalogPart, "update") {
		if content, ok = cat.V1UpdateContent[catalogPart]; !ok {
			return c.JSON(404, fmt.Errorf("update file %s does not exist", catalogPart))
		}
	} else {
		if content, ok = cat.V1PartContent[catalogPart]; !ok {
			return c.JSON(404, fmt.Errorf("part file %s does not exist", catalogPart))
		}
	}
	return c.JSON(200, content)
}
