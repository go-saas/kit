package menu

import _ "embed"

//go:embed menu.yaml
var menuData []byte

var (
	seedMenus [][]byte
)

func init() {
	LoadFromYaml(menuData)
}
func LoadFromYaml(data []byte) {
	seedMenus = append(seedMenus, data)
}

func WalkMenus(f func(menu []byte) error) error {
	for _, menu := range seedMenus {
		if err := f(menu); err != nil {
			return err
		}
	}
	return nil
}
