package databaseutil

import (
	"encoding/json"

	"github.com/samber/lo"
	"golang.org/x/exp/slices"
)

func BytesToStrings(b []byte) ([]string, error) {
	var s []string
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}

	slices.Sort(s)
	s = slices.Compact(s)
	s = lo.Filter(s, func(obj string, _ int) bool { return obj != "" })
	return s, nil
}
