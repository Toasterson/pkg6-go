package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type V1CatalogUpdateFile struct {
	Publishers map[string]V1PublisherUpdate
	Signature  Signature `json:"_SIGNATURE"`
}

func (f V1CatalogUpdateFile) MarshalJSON() (blob []byte, err error) {
	buff := bytes.NewBufferString("{")
	counter := 0
	lenOfMap := len(f.Publishers)
	for key, value := range f.Publishers {
		if pubBlob, err := json.Marshal(value); err != nil {
			return nil, err
		} else {
			buff.WriteString(fmt.Sprintf("\"%s\": %s", key, pubBlob))
		}
		if counter < (lenOfMap - 1) {
			buff.WriteString(",")
		}
		counter++
	}
	buff.WriteString(",")
	if sigBlob, err := json.Marshal(f.Signature); err != nil {
		return nil, err
	} else {
		buff.WriteString(fmt.Sprintf("\"_SIGNATURE\": %s", sigBlob))
	}
	buff.WriteString("}")
	return buff.Bytes(), nil
}

func (f *V1CatalogUpdateFile) UnmarshalJSON(blob []byte) error {
	raw := map[string]*json.RawMessage{}
	if err := json.Unmarshal(blob, &raw); err != nil {
		return err
	}
	if len(f.Publishers) == 0 {
		f.Publishers = make(map[string]V1PublisherUpdate)
	}
	for key, val := range raw {
		if key == "_SIGNATURE" {
			sig := Signature{}
			if err := json.Unmarshal(*val, &sig); err != nil {
				return err
			}
			f.Signature = sig
		} else {
			pubUp := V1PublisherUpdate{}
			if err := json.Unmarshal(*val, &pubUp); err != nil {
				return err
			}
			f.Publishers[key] = pubUp
		}
	}
	return nil
}

type V1PublisherUpdate struct {
	Packages map[string]V1PackageUpdates
}

func (pub V1PublisherUpdate) MarshalJSON() (blob []byte, err error) {
	buff := bytes.NewBufferString("{")
	counter := 0
	lenOfMap := len(pub.Packages)
	for pkgName, value := range pub.Packages {
		if pkgBlob, err := json.Marshal(value); err != nil {
			return nil, err
		} else {
			buff.WriteString(fmt.Sprintf("\"%s\": %s", pkgName, pkgBlob))
		}
		if counter < (lenOfMap - 1) {
			buff.WriteString(",")
		}
		counter++
	}
	buff.WriteString("}")
	return buff.Bytes(), nil
}

func (pub *V1PublisherUpdate) UnmarshalJSON(blob []byte) error {
	pub.Packages = make(map[string]V1PackageUpdates)
	raw := make(map[string]*json.RawMessage)
	if err := json.Unmarshal(blob, &raw); err != nil {
		return err
	}
	for key, value := range raw {
		pak := make([]V1PackageUpdate, 1)
		if err := json.Unmarshal(*value, &pak); err != nil {
			return err
		}
		pub.Packages[key] = pak
	}
	return nil
}

type V1PackageUpdates []V1PackageUpdate

type V1PackageUpdate struct {
	Base        V1PackageUpdateBase       `json:"catalog.base.C"`
	Dependency  V1PackageUpdateDependency `json:"catalog.dependency.C"`
	Summary     V1PackageUpdateSummary    `json:"catalog.summary.C"`
	Operation   string                    `json:"op-type"`
	OpTimestamp V1CatalogTimeStamp        `json:"op-time"`
	Version     string                    `json:"version"`
}

type V1PackageUpdateBase struct {
	Signature string `json:"signature-sha-1"`
	Version   string `json:"version"`
}

type V1PackageUpdateDependency struct {
	Actions V1Actions `json:"actions"`
	Version string    `json:"version"`
}

type V1PackageUpdateSummary struct {
	Actions V1Actions `json:"actions"`
	Version string    `json:"version"`
}
