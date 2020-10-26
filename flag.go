package gotrix

import "strings"

type PathFlags []string

func (i *PathFlags) String() string {
	return strings.Join(*i, ";")
}

func (i *PathFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}
