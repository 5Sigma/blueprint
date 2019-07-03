package blueprint

import (
	"bytes"
	"os"
	"text/template"

	"github.com/spf13/afero"

	"path"

	"github.com/spf13/viper"
)

var fs = afero.NewMemMapFs()

// Process runs the engine agianst the templates and
// configuration given by the passed provider.
func Process(provider ProjectProvider) error {
	var (
		err          error
		templateData = NewData(viper.New())
	)

	templateData.viper.SetDefault("build_path", "build")

	cfgStr, err := provider.ReadConfig()
	if err != nil {
		return err
	}
	templateData.viper.SetConfigType("yaml")
	err = templateData.viper.ReadConfig(bytes.NewBufferString(cfgStr))
	if err != nil {
		return err
	}

	provider.ProcessFiles(func(fpath, fname string, tplData string) error {
		buildRoot := templateData.Get("build_path")
		tpl, err := template.New(fname).Parse(string(tplData))
		if err != nil {
			return err
		}
		fs.MkdirAll(path.Join(buildRoot, fpath), os.ModePerm)
		dstPath := path.Join(buildRoot, fpath, fname)
		f, err := fs.Create(dstPath)
		if err != nil {
			return err
		}
		defer f.Close()
		return tpl.Execute(f, templateData)
	})

	return nil
}
