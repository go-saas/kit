package email

func (x *Config) Normalize() {
	if len(x.Provider) == 0 {
		if x.Smtp != nil {
			x.Provider = "smtp"
		}
	}
}
