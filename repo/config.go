package repo

import (
	"bytes"
	"fmt"
)

const (
	Version3 Version = 3
	Version4         = 4
)

const DefaultTrustAnchorDirectory = "/etc/certs/CA"

type Version int

type Config struct {
	Version                    Version
	PublisherPrefix            string
	SignatureRequiredNames     SignatureRequiredNames
	CheckCertificateRevocation StringeAbleBool
	TrustAnchorDirectory       string
}

type StringeAbleBool bool

func (s StringeAbleBool) String() string {
	if s {
		return "True"
	}
	return "False"
}

type SignatureRequiredNames []string

func (l SignatureRequiredNames) String() string {
	if l == nil {
		return "[]"
	}
	buff := bytes.NewBuffer(nil)
	buff.WriteString("[")
	for i, n := range l {
		if i == 0 {
			buff.WriteString(fmt.Sprintf("\"%s\"", n))
		} else {
			buff.WriteString(fmt.Sprintf(",\"%s\"", n))
		}
	}
	buff.WriteString("]")
	return buff.String()
}
