package data

type AddressEntity struct {
	//Country or region
	Country string
	//State or province
	State   string
	City    string
	ZipCode string
	Line1   string
	Line2   string
	Line3   string
	//TODO database geo?
	Longitude string
	//TODO database geo?
	Latitude string
}

func NewAddressEntityFromPb(s *Address) *AddressEntity {
	return &AddressEntity{
		Country:   s.Country,
		State:     s.State,
		City:      s.City,
		ZipCode:   s.ZipCode,
		Line1:     s.Line1,
		Line2:     s.Line2,
		Line3:     s.Line3,
		Longitude: s.Longitude,
		Latitude:  s.Latitude,
	}
}

func (s *AddressEntity) ToPb() *Address {
	return &Address{
		Country:   s.Country,
		State:     s.State,
		City:      s.City,
		ZipCode:   s.ZipCode,
		Line1:     s.Line1,
		Line2:     s.Line2,
		Line3:     s.Line3,
		Longitude: s.Longitude,
		Latitude:  s.Latitude,
	}
}
