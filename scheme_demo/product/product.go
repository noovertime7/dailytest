package product

import "github.com/noovertime7/dailytest/scheme_demo/scheme"

type Food struct {
	Name  string
	Price int64
}

func (f Food) GetName() string {
	return f.Name
}

func (f Food) GetPrice() int64 {
	return f.Price
}

var SchemeGroupVersion = scheme.GroupVersionKind{Group: "food", Version: "v1"}

var (
	schemeBuilder      = scheme.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &schemeBuilder
	AddToScheme        = localSchemeBuilder.AddScheme
)

func addKnownTypes(scheme *scheme.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion, &Food{Name: "food", Price: 10})
	return nil
}
