package blueprint

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestMain(m *testing.M) {
	fs = afero.NewMemMapFs()
	fs.Mkdir("/project", os.ModePerm)
	os.Exit(m.Run())
}

func processFile(config, template string) (string, error) {
	afero.WriteFile(fs, "/project/bp.yaml",
		[]byte(strings.ReplaceAll(config, "\t", "  ")), os.ModePerm)
	afero.WriteFile(fs, "/project/template.txt", []byte(template), os.ModePerm)
	provider := NewFileSystemProvider("/project")
	err := Process(provider)
	if err != nil {
		return "", err
	}

	if v, err := afero.Exists(fs, "build/template.txt"); err != nil || !v {
		return "", errors.New("Template not written")
	}
	b, err := afero.ReadFile(fs, "build/template.txt")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func TestBasicTemplate(t *testing.T) {
	_, err := processFile("basic_variable: 123", `value is {{ .Get "basic_variable"  }}`)
	if err != nil {
		t.Fatal(err)
	}
	if v, err := afero.Exists(fs, "build/template.txt"); err != nil || !v {
		t.Fatalf("template not written: %s", err)
	}
}

func TestEnvironmentReading(t *testing.T) {
	config := ``
	os.Setenv("codex-test", "123")
	template := `variable is {{ .Env "codex-test" }}`
	str, err := processFile(config, template)
	if err != nil {
		t.Fatal(err)
	}
	if str != "variable is 123" {
		t.Errorf("Template: %s", str)
	}
}

func TestIterateVariable(t *testing.T) {
	config := `
		root:
		- item1
		- item2
	`
	template := `{{range $i, $v := .GetStrings "root"}} {{$i}} is {{$v}}{{end}}`
	str, err := processFile(config, template)
	if err != nil {
		t.Fatal(err)
	}
	if str != ` 0 is item1 1 is item2` {
		t.Errorf("Template: %s", str)
	}
}

func TestNestedFiles(t *testing.T) {
	templateStr := "template"
	configStr := ""
	afero.WriteFile(fs, "/project/bp.yaml", []byte(strings.ReplaceAll(configStr, "\t", "  ")), os.ModePerm)
	afero.WriteFile(fs, "/project/sub/a/template.txt", []byte(templateStr), os.ModePerm)
	provider := NewFileSystemProvider("/project")
	err := Process(provider)
	if err != nil {
		t.Fatal(err)
	}
	if v, err := afero.Exists(fs, "build/sub/a/template.txt"); err != nil || !v {
		t.Fatal(errors.New("Template not written"))
	}
	b, err := afero.ReadFile(fs, "build/sub/a/template.txt")
	if err != nil {
		t.Fatal(err)
	}
	str := string(b)
	if str != "template" {
		t.Errorf("template: %s", str)
	}
}
