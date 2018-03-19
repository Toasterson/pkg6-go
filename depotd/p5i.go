package depotd

import "github.com/labstack/echo"

/*
sample p5i
/{publisher}/p5i/0/{fmri}
{
  "packages": [],
  "publishers": [
    {
      "alias": null,
      "intermediate_certs": [],
      "name": "openindiana.org",
      "packages": [
        "web/server/nginx"
      ],
      "repositories": [],
      "signing_ca_certs": []
    }
  ],
  "version": 1
}
*/
func (d *DepotServer) handleP5IV0(c echo.Context) error {
	return nil
}
