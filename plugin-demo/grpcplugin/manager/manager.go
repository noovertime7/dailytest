package manager

import (
	"github.com/sirupsen/logrus"
	"plugin-demo/common"
	"plugin-demo/grpcplugin/process"
	"plugin-demo/grpcplugin/types"
	"sync"
)

type Manager interface {
	CleanupClients()
	GetKV(version string) (types.KV, error)
}

// manager implements Manager.
type manager struct {
	logger   logrus.FieldLogger
	logLevel logrus.Level
	registry process.Registry

	restartableProcessFactory process.RestartableProcessFactory

	// lock guards restartableProcesses
	lock                 sync.Mutex
	restartableProcesses map[string]process.RestartableProcess
}

// NewManager constructs a manager for getting plugins.
func NewManager(logger logrus.FieldLogger, level logrus.Level, registry process.Registry) Manager {
	return &manager{
		logger:   logger,
		logLevel: level,
		registry: registry,

		restartableProcessFactory: process.NewRestartableProcessFactory(),

		restartableProcesses: make(map[string]process.RestartableProcess),
	}
}

func (m *manager) CleanupClients() {
	m.lock.Lock()

	for _, restartableProcess := range m.restartableProcesses {
		restartableProcess.Stop()
	}

	m.lock.Unlock()
}

func (m *manager) getRestartableProcess(kind common.PluginKind, version string) (process.RestartableProcess, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	logger := m.logger.WithFields(logrus.Fields{
		"kind":    kind.String(),
		"version": version,
	})
	logger.Debug("looking for plugin in registry")

	info, err := m.registry.Get(kind, version)
	if err != nil {
		return nil, err
	}

	logger = logger.WithField("command", info.Command)

	restartableProcess, found := m.restartableProcesses[info.Command]
	if found {
		logger.Debug("found preexisting restartable plugin process")
		return restartableProcess, nil
	}

	logger.Debug("creating new restartable plugin process")

	restartableProcess, err = m.restartableProcessFactory.NewRestartableProcess(info.Command, m.logger, m.logLevel)
	if err != nil {
		return nil, err
	}

	m.restartableProcesses[info.Command] = restartableProcess

	return restartableProcess, nil
}

func (m *manager) GetKV(version string) (types.KV, error) {
	restartableProcess, err := m.getRestartableProcess(common.PluginKV, version)
	if err != nil {
		return nil, err
	}
	plugin, err := restartableProcess.GetByKindAndVersion(process.KindAndVersion{Kind: common.PluginKV, Version: version})
	if err != nil {
		return nil, err
	}
	return plugin.(types.KV), nil
}
