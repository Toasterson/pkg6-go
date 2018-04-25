package repo

import (
	"crypto/sha1"
	"encoding/json"
	"git.wegmueller.it/toasterson/glog"
	"github.com/cavaliercoder/grab"
	"github.com/spf13/viper"
	"github.com/toasterson/pkg6-go/metadata"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type publisherJson struct {
	Packages   []interface{}      `json:"packages"`
	Version    int                `json:"version"`
	Publishers []publisherJsonPub `json:"publishers"`
}

type publisherJsonPub struct {
	Alias             interface{}   `json:"alias"`
	IntermediateCerts []interface{} `json:"intermediate_certs"`
	Name              string        `json:"name"`
	Packages          []interface{} `json:"packages"`
	Repositories      []interface{} `json:"repositories"`
	SigningCACerts    []interface{} `json:"signing_ca_certs"`
}

type HttpRepo struct {
	Name                       string       `json:"-"`
	BaseUrl                    string       `json:"-"`
	CacheDir                   string       `json:"-"`
	Client                     *grab.Client `json:"-"`
	Version                    int
	TrustAnchorDirectory       string
	CheckCertificateRevocation bool
	SignatureRequiredNames     []string
	Publishers                 []string                       `json:"-"`
	Catalogs                   map[string]*metadata.V1Catalog `json:"-"`
}

func (r *HttpRepo) loadPublishersFromRemote() (err error) {
	dstDir := r.CacheDir
	pubFileName := filepath.Join(dstDir, "publisher.json")
	_, err = grab.Get(pubFileName, r.BaseUrl+"/publisher/0")
	if err != nil {
		return err
	}
	tmp := publisherJson{}
	var f *os.File
	if f, err = os.Open(pubFileName); err != nil {
		return err
	}
	defer func() {
		f.Close()
		os.Remove(pubFileName)
	}()
	if err = json.NewDecoder(f).Decode(&tmp); err != nil {
		return err
	}
	for _, pub := range tmp.Publishers {
		r.Publishers = append(r.Publishers, pub.Name)
	}
	return nil
}

func (r *HttpRepo) GetFile(publisher string, hash string) (*os.File, error) {
	panic("implement me")
}

func (r *HttpRepo) GetPackage(fmri string) (metadata.PackageInfo, error) {
	panic("implement me")
}

func (r *HttpRepo) GetPath() string {
	panic("implement me")
}

func (r *HttpRepo) GetPublishers() []string {
	panic("implement me")
}

func (r *HttpRepo) GetPackageFMRIs(publisher string, partial bool) []string {
	panic("implement me")
}

func (r *HttpRepo) Create() error {
	panic("implement me")
}

func (r *HttpRepo) Load() error {
	r.Client = grab.NewClient()
	if r.CacheDir == "" {
		r.CacheDir = filepath.Join(viper.GetString("cache_dir"), r.Name)
	}
	if stat, err := os.Lstat(r.CacheDir); os.IsNotExist(err) {
		os.MkdirAll(r.CacheDir, 0755)
	} else if stat.Mode().IsRegular() {
		panic("cache directory is a file aborting")
	}
	return r.loadPublishersFromRemote()
}

func (r *HttpRepo) Save() (err error) {
	panic("implement me")
}

func (r *HttpRepo) Destroy() error {
	panic("implement me")
}

func (r *HttpRepo) GetVersion() int {
	return 4
}

func (r *HttpRepo) DownloadAllCatalogs() error {
	for _, pub := range r.Publishers {
		if err := r.DownloadCatalog(pub); err != nil {
			return err
		}
	}
	return nil
}

// This function gets the Catalog from the Remote Repository and Writes is to it's cache directory
func (r *HttpRepo) DownloadCatalog(publisher string) error {
	dstDir := filepath.Join(r.CacheDir, "catalog", publisher)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}
	// Get the Index File
	if _, err := grab.Get(filepath.Join(dstDir, "catalog.attrs"), r.BaseUrl+"/"+publisher+"/catalog/1/catalog.attrs"); err != nil {
		return err
	}
	// Load the Indexfile into a struct for processing
	attrs := metadata.V1Catalog{}
	var content []byte
	var err error
	if content, err = ioutil.ReadFile(filepath.Join(dstDir, "catalog.attrs")); err != nil {
		return err
	}
	if err = json.Unmarshal(content, &attrs); err != nil {
		return err
	}

	// Now get all parts
	for part, partMeta := range attrs.Parts {
		req, err := grab.NewRequest(filepath.Join(dstDir, part), r.BaseUrl+"/"+publisher+"/catalog/1/"+part)
		if err != nil {
			return err
		}
		//TODO make replaceable with whatever hasching method is available.
		req.SetChecksum(sha1.New(), []byte(partMeta.Signature), true)
		glog.Infof("Downloading catalog part %s", part)
		timer := time.NewTimer(1 * time.Second)
		resp := r.Client.Do(req)

	PartDlLoop:
		for {
			select {
			case <-timer.C:
				glog.Infof("Loaded %s%", resp.Progress())
			case <-resp.Done:
				break PartDlLoop
			}
		}
		timer.Stop()
		if err := resp.Err(); err != nil {
			glog.Emergf("Error while Downloading %s: %s", part, err)
		} else {
			glog.Infof("Finished downloading %s", part)
		}
	}

	// Now grab all Update Parts
	for update, meta := range attrs.Updates {
		req, err := grab.NewRequest(filepath.Join(dstDir, update), r.BaseUrl+"/"+publisher+"/catalog/1/"+update)
		if err != nil {
			return err
		}
		req.SetChecksum(sha1.New(), []byte(meta.Signature), true)
		glog.Infof("Downloading catalog update %s", update)
		timer := time.NewTimer(1 * time.Second)
		resp := r.Client.Do(req)

	DlLoop:
		for {
			select {
			case <-timer.C:
				glog.Infof("Loaded %s%", resp.Progress())
			case <-resp.Done:
				break DlLoop
			}
		}

		timer.Stop()
		if err := resp.Err(); err != nil {
			glog.Emergf("Error while Downloading %s: %s", update, err)
		} else {
			glog.Infof("Finished downloading %s", update)
		}
	}
	return nil
}

func (r *HttpRepo) GetCatalogFile(publisher, part string) (*os.File, error) {
	//Todo Auto download
	s := filepath.Join(r.CacheDir, "catalog", publisher, part)
	return os.Open(s)
}

func (r *HttpRepo) GetCatalog(publisher string) *metadata.V1Catalog {
	panic("implement me")
}

func (r *HttpRepo) Search(params map[string]string, query string) string {
	panic("implement me")
}

func (r *HttpRepo) AddPackage(info metadata.PackageInfo) error {
	panic("implement me")
}
