package conf

import "fmt"

func (x *Endpoints)GetEventOrDefault(name string) *Event {
	var res *Event
	var ok bool
	if name!=""{
		res,ok =x.Events[name]
	}
	if !ok{
		res,ok = x.Events["default"]
		if !ok{
			panic(fmt.Sprintf("cannot resolve event %s",name))
		}
	}
	return res
}