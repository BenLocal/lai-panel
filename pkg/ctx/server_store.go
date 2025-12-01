package ctx

import (
	"fmt"
	"log"
	"sync"

	"github.com/benlocal/lai-panel/pkg/options"
)

var GlobalServerStore *ServerStore
var once sync.Once

var baseIP = "127.0.0.1"

func GetServerStore(opt *options.AgentOptions) *ServerStore {
	once.Do(func() {
		dataPath := opt.DataPath()
		GlobalServerStore = newServerStore(false,
			opt.MasterHost(),
			opt.MasterPort(),
			opt.Name,
			opt.Port,
			opt.Address,
			&dataPath)
		log.Println(GlobalServerStore.str())
	})

	return GlobalServerStore
}

func GetServerStoreForLocal(opt *options.ServeOptions) *ServerStore {
	once.Do(func() {
		dataPath := opt.DataPath()
		GlobalServerStore = newServerStore(true,
			baseIP,
			opt.Port,
			"local",
			opt.Port,
			baseIP,
			&dataPath)
		log.Println(GlobalServerStore.str())
	})
	return GlobalServerStore
}

type ServerStore struct {
	isLocal    bool
	masterHost string
	masterPort int
	name       string
	id         int64
	agentPort  int
	address    string
	dataPath   *string

	mu sync.Mutex
}

func newServerStore(isLocal bool,
	masterHost string,
	masterPort int,
	name string,
	agentPort int,
	address string,
	dataPath *string) *ServerStore {
	return &ServerStore{
		isLocal:    isLocal,
		masterHost: masterHost,
		masterPort: masterPort,
		name:       name,
		agentPort:  agentPort,
		address:    address,
		dataPath:   dataPath,
	}
}

func (s *ServerStore) str() string {
	return fmt.Sprintf(`ServerStore{isLocal: %v, masterHost: %s, masterPort: %d, name: %s, agentPort: %d, address: %s, dataPath: %v}`,
		s.isLocal,
		s.masterHost,
		s.masterPort,
		s.name,
		s.agentPort,
		s.address,
		*s.dataPath)
}

func (s *ServerStore) GetID() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.id
}

func (s *ServerStore) SetID(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.id = id
}

func (s *ServerStore) GetName() string {
	return s.name
}

func (s *ServerStore) GetMasterHost() string {
	return s.masterHost
}

func (s *ServerStore) GetMasterPort() int {
	return s.masterPort
}

func (s *ServerStore) IsLocal() bool {
	return s.isLocal
}

func (s *ServerStore) GetAddress() string {
	return s.address
}

func (s *ServerStore) GetAgentPort() int {
	return s.agentPort
}

func (s *ServerStore) GetDataPath() *string {
	return s.dataPath
}
