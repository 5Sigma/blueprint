package blueprint

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"

	"github.com/spf13/viper"
)

func setupConfig(configStr string) (*TemplateData, error) {
	var err error
	fs := afero.NewMemMapFs()
	configStr = strings.ReplaceAll(configStr, "\t", "  ")
	println(configStr)
	afero.WriteFile(fs, "/codex.yaml", []byte(configStr), os.ModePerm)
	v := viper.New()
	v.SetFs(fs)
	v.SetConfigName("codex")
	v.AddConfigPath("/")
	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &TemplateData{viper: v}, nil
}

func TestGet(t *testing.T) {
	d, err := setupConfig(`
	test: 123
	root:
		nested:
			item: 421
	`)
	if err != nil {
		t.Fatal(err)
	}
	if v := d.Get("test"); v != "123" {
		t.Errorf("value: %s", v)
	}
	if v := d.Get("root.nested.item"); v != "421" {
		t.Errorf("value: %s", v)
	}
}
