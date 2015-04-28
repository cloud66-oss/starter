package common

import (
	"strings"
)

type Lister struct {
	Items []string
}

func NewLister(seed ...string) *Lister {
	l := &Lister{}
	l.Items = seed

	return l
}

func (l *Lister) ToList(sep string) string {
	return strings.Join(l.Items, sep)
}
