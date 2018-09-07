package metadata

import (
	"encoding/json"
	"fmt"
	"time"
)

const V1CatalogTimeStampFormat = "20060102T150405.000000Z"

type V1CatalogTimeStamp time.Time

func (t *V1CatalogTimeStamp) UnmarshalJSON(blob []byte) error {
	var stringTime string
	if err := json.Unmarshal(blob, &stringTime); err != nil {
		return err
	}
	if parserdTime, err := time.Parse(V1CatalogTimeStampFormat, stringTime); err != nil {
		return err
	} else {
		*t = V1CatalogTimeStamp(parserdTime)
		return nil
	}
}

func (t V1CatalogTimeStamp) MarshalJSON() (blob []byte, err error) {
	tim := time.Time(t)
	stringTime := tim.Format(V1CatalogTimeStampFormat)
	if stringTime == "" {
		return nil, fmt.Errorf("could not Format %s as string", t)
	}
	if blob, err = json.Marshal(stringTime); err != nil {
		return nil, err
	}
	return blob, nil
}
