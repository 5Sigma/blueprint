package codex

import (
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/afero"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

// FileHandler is the function handlerpassed to ProcessFiles.
// It is executed once for each file within the project
type FileHandler func(fpath string, fname string, data string) error

// ProjectProvider provides an interface to the project files.
type ProjectProvider interface {
	ReadConfig() (string, error)
	ProcessFiles(h FileHandler) error
}

// FileSystemProvider provides the project files via the file system
type FileSystemProvider struct {
	ProjectFolder string
}

// NewFileSystemProvider creates a new provider using the underlying OS filesystem as a source
func NewFileSystemProvider(fpath string) *FileSystemProvider {
	return &FileSystemProvider{
		ProjectFolder: fpath,
	}
}

// ReadConfig reads the configuration file from the project root
func (fsp *FileSystemProvider) ReadConfig() (string, error) {
	cfgPath := path.Join(fsp.ProjectFolder, "bp.yaml")
	if v, err := afero.Exists(fs, cfgPath); !v || err != nil {
		return "", errors.New("Configuration file not found")
	}
	b, err := afero.ReadFile(fs, cfgPath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ProcessFiles calls the handler for every file in the project
// folder
func (fsp *FileSystemProvider) ProcessFiles(h FileHandler) error {
	return afero.Walk(fs, fsp.ProjectFolder, func(fpath string, info os.FileInfo, err error) error {
		fname := path.Base(fpath)
		relPath, _ := filepath.Rel(fsp.ProjectFolder, fpath)
		if fname == "bp.yaml" {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		bytes, err := afero.ReadFile(fs, fpath)
		return h(path.Dir(relPath), fname, string(bytes))
	})
}

// GitProvider provides the project via a remote git repository.
// The repository is cloned and managed in memory.
type GitProvider struct {
	URL  string
	repo *git.Repository
}

// NewGitProvider initializes a new provider that uses a
// remote git repository as its source
func NewGitProvider(url string) *GitProvider {
	return &GitProvider{
		URL: url,
	}
}

func (gp *GitProvider) clone() error {
	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: gp.URL,
	})
	if err != nil {
		return err
	}
	gp.repo = repo
	return nil
}

// ReadConfig reads the bp.yaml file out of the root of the repository
func (gp *GitProvider) ReadConfig() (string, error) {
	w, err := gp.repo.Worktree()
	if err != nil {
		return "", err
	}

	f, err := w.Filesystem.Open("bp.yaml")

	if err != nil {
		return "", err
	}
	var b []byte
	i, err := f.Read(b)
	if i == 0 {
		return "", errors.New("Nothing read from config")
	}
	if err != nil {
		return "", err
	}
	return string(b), nil
}
