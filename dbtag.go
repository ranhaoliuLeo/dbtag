package dbtag

import (
	"errors"
	"reflect"
	"strings"
)

// Proxy interface of dbtag ability
type Proxy interface {
	// GetColsWithOmit input selectors you wanna omit
	GetColsWithOmit(selectors ...string) []string
	// GetColsWithSelect input selectors you wanna select
	GetColsWithSelect(selectors ...string) []string
}

type instance struct {
	sample           interface{}
	allCols          []string
	colsMapSelectors map[string][]string
}

var ErrStructFieldEmpty = errors.New("empty field in the sample")

func New(sample interface{}) (Proxy, error) {
	elems := reflect.TypeOf(sample).Elem()
	fieldNum := elems.NumField()
	if fieldNum <= 0 {
		return nil, ErrStructFieldEmpty
	}
	// store all db cols
	allCols := make([]string, 0, fieldNum)
	// store the mapping relation between cols and selectors
	colsWithSelectors := make(map[string][]string, fieldNum)
	for i := 0; i < fieldNum; i++ {
		dbCol := elems.Field(i).Tag.Get("db")
		selectors := strings.Split(elems.Field(i).Tag.Get("selector"), ",")
		allCols = append(allCols, dbCol)
		colsWithSelectors[dbCol] = selectors
	}
	return &instance{
		sample:           sample,
		allCols:          allCols,
		colsMapSelectors: colsWithSelectors,
	}, nil
}

// GetColsWithOmit input selectors you wanna omit,
// if you omit empty, it will return all cols.
// The order of the return tag is the same as struct orde.
func (i *instance) GetColsWithOmit(selectors ...string) []string {
	if len(selectors) == 0 {
		return i.allCols
	}
	cols := make([]string, 0, len(i.allCols))
	for _, col := range i.allCols {
		s, ok := i.colsMapSelectors[col]
		if !ok {
			panic("get selector fail, may cause db fail")
		}
		if !isIntersect(s, selectors) {
			cols = append(cols, col)
		}
	}
	return cols
}

// GetColsWithSelect input selectors you wanna select,
// if you select empty, it will return nil.
// The order of the return tag is the same as struct order.
func (i *instance) GetColsWithSelect(selectors ...string) []string {
	// 这个应该不在这里做
	if len(i.allCols) == 0 || len(selectors) == 0 {
		return nil
	}
	cols := make([]string, 0, len(i.allCols))
	for _, col := range i.allCols {
		s, ok := i.colsMapSelectors[col]
		if !ok {
			panic("get selector fail, may cause db fail")
		}
		if isIntersect(s, selectors) {
			cols = append(cols, col)
		}
	}
	return cols
}

func isIntersect(s1, s2 []string) bool {
	m := make(map[string]bool)
	for _, s := range s1 {
		m[s] = true
	}
	for _, s := range s2 {
		if _, ok := m[s]; ok {
			return true
		}
	}
	return false
}
