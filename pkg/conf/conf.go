package conf

import "fmt"

func (x *Endpoints) GetEventOrDefault(name string) *Event {
	var res *Event
	var ok bool
	if name != "" {
		res, ok = x.Events[name]
	}
	if !ok {
		res, ok = x.Events["default"]
		if !ok {
			panic(fmt.Sprintf("cannot resolve event %s", name))
		}
	}
	return res
}

func (x *Endpoints) GetDatabaseOrDefault(name string) *Database {
	var res *Database
	var ok bool
	if name != "" {
		res, ok = x.Databases[name]
	}
	if !ok {
		res, ok = x.Databases["default"]
		if !ok {
			panic(fmt.Sprintf("cannot resolve event %s", name))
		}
	}
	return res
}
