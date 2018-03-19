package catalog

import (
	"fmt"
)

type Signature struct {
	SHA1 string `json:"sha-1"`
}

func (s *Signature) Check(signature string) error {
	if s.SHA1 == signature {
		return nil
	}
	return fmt.Errorf("signature check failed %s != %s", s.SHA1, signature)
}
