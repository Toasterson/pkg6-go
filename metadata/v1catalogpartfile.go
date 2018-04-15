package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type V1CatalogPartFile struct {
	Publishers map[string]V1Publisher
	Signature  Signature `json:"_SIGNATURE"`
}

func (f V1CatalogPartFile) MarshalJSON() (blob []byte, err error) {
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

func (f *V1CatalogPartFile) UnmarshalJSON(blob []byte) error {
	raw := map[string]*json.RawMessage{}
	if err := json.Unmarshal(blob, &raw); err != nil {
		return err
	}
	if len(f.Publishers) == 0 {
		f.Publishers = make(map[string]V1Publisher)
	}
	for key, val := range raw {
		if key == "_SIGNATURE" {
			sig := Signature{}
			if err := json.Unmarshal(*val, &sig); err != nil {
				return err
			}
			f.Signature = sig
		} else {
			pub := V1Publisher{}
			if err := json.Unmarshal(*val, &pub); err != nil {
				return err
			}
			f.Publishers[key] = pub
		}
	}
	return nil
}

type V1Publisher struct {
	Packages map[string]V1Packages
}

func (pub V1Publisher) MarshalJSON() (blob []byte, err error) {
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

func (pub *V1Publisher) UnmarshalJSON(blob []byte) error {
	pub.Packages = make(map[string]V1Packages)
	raw := make(map[string]*json.RawMessage)
	if err := json.Unmarshal(blob, &raw); err != nil {
		return err
	}
	for key, value := range raw {
		pak := V1Packages{}
		if err := json.Unmarshal(*value, &pak); err != nil {
			return err
		}
		pub.Packages[key] = pak
	}
	return nil
}

type V1Packages []V1PackageVersion

type V1PackageVersion struct {
	Actions   V1Actions `json:"actions,omitempty"`
	Signature string    `json:"signature-sha-1,omitempty"`
	Version   string    `json:"version"`
}

type V1Actions []string
