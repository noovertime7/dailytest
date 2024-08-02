package process

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"plugin-demo/common"
	"plugin-demo/grpcplugin/plugin/filesystem"
	"regexp"
	"runtime"
)

type PluginIdentifier struct {
	Command string
	Kind    common.PluginKind
	Version string
}

type Registry interface {
	DiscoverPlugins() error
	Get(kind common.PluginKind, name string) (PluginIdentifier, error)
}

// KindAndName is a convenience struct that combines a PluginKind and a name.
type KindAndVersion struct {
	Kind    common.PluginKind
	Version string
}

// registry implements Registry.
type registry struct {
	dir      string
	logger   logrus.FieldLogger
	logLevel logrus.Level

	processFactory Factory
	fs             filesystem.Interface
	pluginsByID    map[KindAndVersion]PluginIdentifier
	pluginsByKind  map[common.PluginKind]PluginIdentifier
}

func (r *registry) Get(kind common.PluginKind, version string) (PluginIdentifier, error) {
	p, found := r.pluginsByID[KindAndVersion{Kind: kind, Version: version}]
	if !found {
		return PluginIdentifier{}, errors.Errorf("plugin not found")
	}
	return p, nil
}

// NewRegistry returns a new registry.
func NewRegistry(dir string, logger logrus.FieldLogger, logLevel logrus.Level) Registry {
	return &registry{
		dir:      dir,
		logger:   logger,
		logLevel: logLevel,

		processFactory: newProcessFactory(),
		fs:             filesystem.NewFileSystem(),
		pluginsByID:    make(map[KindAndVersion]PluginIdentifier),
		pluginsByKind:  make(map[common.PluginKind]PluginIdentifier),
	}
}

func (r *registry) DiscoverPlugins() error {
	plugins, err := r.readPluginsDir(r.dir)
	if err != nil {
		return err
	}

	for _, command := range plugins {
		version, kind, err := parseFilePath(command)
		if err != nil {
			return err
		}
		if err := r.register(PluginIdentifier{
			Command: command,
			Kind:    common.PluginKind(kind),
			Version: version,
		}); err != nil {
			return err
		}
	}

	return nil
}

func parseFilePath(filePath string) (version, kind string, err error) {
	// 提取文件名
	fileName := filepath.Base(filePath)

	// 使用更精确的正则表达式解析文件名
	re := regexp.MustCompile(`^([a-zA-Z]+)-(\d+\.\d+\.\d+)(\.exe)?$`)
	matches := re.FindStringSubmatch(fileName)

	if len(matches) != 4 {
		return "", "", fmt.Errorf("invalid file format")
	}

	kind = matches[1]
	version = matches[2]

	return version, kind, nil
}

func (r *registry) register(id PluginIdentifier) error {
	key := KindAndVersion{Kind: id.Kind, Version: id.Version}
	if _, found := r.pluginsByID[key]; found {
		return errors.Errorf("duplicate plugin name %q", id.Kind)
	}

	r.pluginsByID[key] = id
	r.pluginsByKind[id.Kind] = id

	return nil
}

// readPluginsDir recursively reads dir looking for plugins.
func (r *registry) readPluginsDir(dir string) ([]string, error) {
	if _, err := r.fs.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, errors.WithStack(err)
	}

	files, err := r.fs.ReadDir(dir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	fullPaths := make([]string, 0, len(files))
	for _, file := range files {
		fullPath := filepath.Join(dir, file.Name())

		if file.IsDir() {
			subDirPaths, err := r.readPluginsDir(fullPath)
			if err != nil {
				return nil, err
			}
			fullPaths = append(fullPaths, subDirPaths...)
			continue
		}

		if !executable(file) {
			continue
		}

		fullPaths = append(fullPaths, fullPath)
	}
	return fullPaths, nil
}

func executable(info os.FileInfo) bool {
	/*
		When we AND the mode with 0111:

		- 0100 (user executable)
		- 0010 (group executable)
		- 0001 (other executable)

		the result will be 0 if and only if none of the executable bits is set.
	*/
	if runtime.GOOS == "windows" && info.Mode()&os.ModeType == 0 {
		return true
	}
	return (info.Mode() & 0111) != 0
}

func isExecutableOnWindows(info os.FileInfo) bool {
	// Check if the file is a regular file and has the ModeType bit set,
	// which indicates it's an executable on Windows.
	return info.Mode().IsRegular() && info.Mode()&os.ModeType != 0
}
