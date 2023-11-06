package utils

func String(a string) *string {
	if len(a) == 0 {
		return nil
	}
	return &a
}
