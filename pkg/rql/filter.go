package rql

type FilterWrapper struct {
	v map[string]interface{}
}

func NewFilterWrapper(v map[string]interface{}) *FilterWrapper {
	return &FilterWrapper{
		v: v,
	}
}

func (f *FilterWrapper) GetFilter() map[string]interface{} {
	return f.v
}
