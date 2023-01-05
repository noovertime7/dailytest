package scheme

import (
	"errors"
	"fmt"
	"github.com/noovertime7/dailytest/scheme_demo/object"
	"reflect"
)

type Scheme struct {
	Name      string
	Price     int64
	typeToGVK map[reflect.Type][]GroupVersionKind
	product   map[GroupVersionKind]object.Product
}

type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

func (s *Scheme) AddKnownTypes(gvk GroupVersionKind, obj object.Product) {
	t := reflect.TypeOf(obj)
	if len(gvk.Version) == 0 {
		panic(fmt.Sprintf("version is required on all types: %s %v", gvk, t))
	}
	if t.Kind() != reflect.Ptr {
		panic("All types must be pointers to structs.")
	}
	t = t.Elem()
	if t.Kind() != reflect.Struct {
		panic("All types must be pointers to structs.")
	}
	s.typeToGVK[t] = append(s.typeToGVK[t], gvk)
	s.product[gvk] = obj
}

func (s *Scheme) GetName(g GroupVersionKind) (string, error) {
	obj, ok := s.product[g]
	if !ok {
		return "", errors.New("not found ")
	}
	return obj.GetName(), nil
}

func (s *Scheme) GetPrice(g GroupVersionKind) (int64, error) {
	obj, ok := s.product[g]
	if !ok {
		return 0, errors.New("not found ")
	}
	return obj.GetPrice(), nil
}

func NewScheme() *Scheme {
	return &Scheme{
		typeToGVK: map[reflect.Type][]GroupVersionKind{},
		product:   map[GroupVersionKind]object.Product{},
	}
}

type SchemeBuilder []func(s *Scheme) error

func (sb *SchemeBuilder) AddScheme(s *Scheme) error {
	for _, f := range *sb {
		if err := f(s); err != nil {
			return err
		}
	}
	return nil
}

func (sb *SchemeBuilder) Register(funcs ...func(s *Scheme) error) {
	for _, f := range funcs {
		*sb = append(*sb, f)
	}
}

func NewSchemeBuilder(funcs ...func(*Scheme) error) SchemeBuilder {
	sb := SchemeBuilder{}
	sb.Register(funcs...)
	return sb
}
