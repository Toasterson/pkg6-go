package metadata

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

var catalogAttrsSource = `{"created": "20161224T160808.632593Z","last-modified": "20161224T161749.202734Z","package-count": 2,"package-version-count": 2,"parts": {"catalog.base.C": {"last-modified": "20161224T161749.202734Z","signature-sha-1": "9534c505d72085797cf142d9c69ea0adc5b2449c"},"catalog.dependency.C": {"last-modified": "20161224T161749.202734Z","signature-sha-1": "f98917b705e74ae20e39a06a7c306cae68c7bd72"},"catalog.summary.C": {"last-modified": "20161224T161749.202734Z","signature-sha-1": "3823dece5982820dd3128b5f9eeeb9448ffea1b7"}},"updates": {"update.20161224T16Z.C": {"last-modified": "20161224T161749.202734Z","signature-sha-1": "183259ff19ecb76bb9766776e09e8ab544828419"}},"version": 1,"_SIGNATURE": {"sha-1": "03bf585ff2dfd495fa752154baf85d2add2c968b"}}`

var catalogBaseSource = `{"userland": {"library/desktop/mate/libmatemixer": [{"signature-sha-1": "0a53b57199ea1fc417cb39d9c28ea16c89f8be2c","version": "1.16.0,5.11-2016.1.1.0:20161224T161749Z"}],"library/desktop/mate/mate-common": [{"signature-sha-1": "55955b2ad9bd5b84ff24e8ad6173d521e99c0eda","version": "1.16.0,5.11-2016.1.1.0:20161224T160903Z"}]},"_SIGNATURE": {"sha-1": "9534c505d72085797cf142d9c69ea0adc5b2449c"}}`

var catalogSummSource = `{"userland": {"library/desktop/mate/libmatemixer": [{"actions": ["set name=userland.info.git-remote value=https://github.com/toasterson/oi-userland","set name=com.oracle.info.name value=libmatemixer","set name=userland.info.git-branch value=oi/hipster","set name=userland.info.git-rev value=084011e1dd452a12a218ecd447b7a1ccd48c897a","set name=pkg.summary value=\"Mixer library for MATE desktop\"","set name=info.classification value=org.opensolaris.category.2008:System/Libraries","set name=info.upstream-url value=http://www.mate-desktop.org","set name=info.source-url value=http://pub.mate-desktop.org/releases/1.16/libmatemixer-1.16.0.tar.xz","set name=org.opensolaris.consolidation value=userland","set name=com.oracle.info.version value=1.16.0"],"version": "1.16.0,5.11-2016.1.1.0:20161224T161749Z"}],"library/desktop/mate/mate-common": [{"actions": ["set name=userland.info.git-remote value=https://github.com/toasterson/oi-userland","set name=com.oracle.info.name value=mate-common","set name=userland.info.git-branch value=oi/hipster","set name=userland.info.git-rev value=084011e1dd452a12a218ecd447b7a1ccd48c897a","set name=pkg.summary value=\"Common automake macros for MATE\"","set name=info.classification value=org.opensolaris.category.2008:System/Libraries","set name=info.upstream-url value=http://www.mate-desktop.org","set name=info.source-url value=http://pub.mate-desktop.org/releases/1.16/mate-common-1.16.0.tar.xz","set name=org.opensolaris.consolidation value=userland","set name=com.oracle.info.version value=1.16.0"],"version": "1.16.0,5.11-2016.1.1.0:20161224T160903Z"}]},"_SIGNATURE": {"sha-1": "3823dece5982820dd3128b5f9eeeb9448ffea1b7"}}`

var catalogDepSource = `{"userland": {"library/desktop/mate/libmatemixer": [{"actions": ["depend fmri=pkg:/library/audio/pulseaudio@6.0-2016.0.0.0 type=require","depend fmri=pkg:/library/glib2@2.46.2-2016.1.0.0 type=require","depend fmri=pkg:/system/library@0.5.11-2016.1.1.16068 type=require","set name=variant.arch value=i386"],"version": "1.16.0,5.11-2016.1.1.0:20161224T161749Z"}],"library/desktop/mate/mate-common": [{"actions": ["depend fmri=pkg:/SUNWcs@0.5.11-2016.1.1.16068 type=require","depend fmri=pkg:/shell/bash@4.3.46-2016.0.0.0 type=require","set name=variant.arch value=i386"],"version": "1.16.0,5.11-2016.1.1.0:20161224T160903Z"}]},"_SIGNATURE": {"sha-1": "f98917b705e74ae20e39a06a7c306cae68c7bd72"}}`

var parts = []string{"catalog.base.C", "catalog.dependency.C", "catalog.summary.C"}

func NewTestingCatalog() (*V1Catalog, error) {
	cat := NewV1Catalog("none")
	if err := json.Unmarshal([]byte(catalogAttrsSource), &cat); err != nil {
		return nil, fmt.Errorf("could not Unmarschal JSON: %s", err)
	}
	for _, part := range parts {
		var source []byte
		switch part {
		case "catalog.base.C":
			source = []byte(catalogBaseSource)
		case "catalog.dependency.C":
			source = []byte(catalogDepSource)
		case "catalog.summary.C":
			source = []byte(catalogSummSource)
		}
		if partFile, err := DeSerializeV1Part(source); err != nil {
			return nil, fmt.Errorf("could not Unmarschal JSON: %s", err)
		} else {
			cat.V1PartContent[part] = partFile
		}
	}

	return cat, nil
}

func TestV1Catalog_SerializeToV1(t *testing.T) {
	var cat *V1Catalog
	var err error
	if cat, err = NewTestingCatalog(); err != nil {
		t.Fatal(err)
	}
	var orig []byte
	var blob []byte
	orig = []byte(catalogAttrsSource)
	if blob, err = json.Marshal(&cat); err != nil {
		t.Fatal(err)
	}
	var valueIndentBlob bytes.Buffer
	var resultIndentBlob bytes.Buffer
	json.Indent(&valueIndentBlob, blob, "", "\t")
	json.Indent(&resultIndentBlob, orig, "", "\t")
	if valueIndentBlob.String() != resultIndentBlob.String() {
		fatalBuff := bytes.NewBufferString("Something went wrong. See the two Results below:\n")
		fatalBuff.WriteString(fmt.Sprintf("Gotten:\n%s\n", valueIndentBlob.String()))
		fatalBuff.WriteString(fmt.Sprintf("Expected:\n%s\n", resultIndentBlob.String()))
		t.Fatal(fatalBuff.String())
	} else {
		t.Logf("successfully processed %s", "catalog.attrs")
	}

	for part := range cat.Parts {
		var partOrig []byte
		var partBlob []byte
		switch part {
		case "catalog.attrs":
			partOrig = []byte(catalogAttrsSource)
		case "catalog.base.C":
			partOrig = []byte(catalogBaseSource)
		case "catalog.dependency.C":
			partOrig = []byte(catalogDepSource)
		case "catalog.summary.C":
			partOrig = []byte(catalogSummSource)
		}
		if partBlob, err = cat.SerializeV1Part(part); err != nil {
			t.Fatal(err)
		}
		var valueIndentBlob bytes.Buffer
		var resultIndentBlob bytes.Buffer
		json.Indent(&valueIndentBlob, partBlob, "", "\t")
		json.Indent(&resultIndentBlob, partOrig, "", "\t")
		if valueIndentBlob.String() != resultIndentBlob.String() {
			fatalBuff := bytes.NewBufferString("Something went wrong. See the two Results below:\n")
			fatalBuff.WriteString(fmt.Sprintf("Gotten:\n%s\n", valueIndentBlob.String()))
			fatalBuff.WriteString(fmt.Sprintf("Expected:\n%s\n", resultIndentBlob.String()))
			t.Fatal(fatalBuff.String())
		} else {
			t.Logf("successfully processed %s", part)
		}

	}
}
