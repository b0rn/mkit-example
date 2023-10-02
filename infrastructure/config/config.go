package config

import "github.com/b0rn/mkit/pkg/mlog"

type Config struct {
	LogConfig        mlog.Config      `yaml:"logConfig"`
	ApisConfig       ApisConfig       `yaml:"apisConfig"`
	DatastoresConfig DatastoresConfig `yaml:"dataStoresConfig"`
	UsecasesConfig   UsecasesConfig   `yaml:"usecasesConfig"`
}

// API global configuration
type ApisConfig struct {
	Address    string      `yaml:"address"`
	Port       uint16      `yaml:"port"`
	Root       string      `yaml:"root"`
	RESTConfig RESTConfig  `yaml:"restConfig"`
	MQConfig   MQApiConfig `yaml:"mqConfig"`
}

// Configuration for a specific API
type RESTConfig struct {
	Code    string `yaml:"code"`
	Enabled bool   `yaml:"enabled"`
	Path    string `yaml:"path"`
}

type MQApiConfig struct {
	Code              string            `yaml:"code"`
	Stream            string            `yaml:"stream"`
	CreateUserSubject string            `yaml:"createUserSubject"`
	Config            MQDataStoreConfig `yaml:"config"`
}

type DatastoresConfig struct {
	DbConfig DbDataStoreConfig `yaml:"dbConfig"`
	MqConfig MQDataStoreConfig `yaml:"mqConfig"`
}

type DbDataStoreConfig struct {
	Code       string `yaml:"code"`
	DriverName string `yaml:"driverName"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	DbName     string `yaml:"dbName"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	Tx         bool   `yaml:"tx"`
}

type MQDataStoreConfig struct {
	Code string `yaml:"code"`
	URL  string `yaml:"url"`
}

type DbDataServiceConfig struct {
	Code            string            `yaml:"code"`
	DataStoreConfig DbDataStoreConfig `yaml:"dataStoreConfig"`
}

type MQDataServiceConfig struct {
	Code               string            `yaml:"code"`
	StreamName         string            `yaml:"streamName"`
	UserCreatedSubject string            `yaml:"userCreatedSubject"`
	DataStoreConfig    MQDataStoreConfig `yaml:"dataStoreConfig"`
}

type ManageUsersUsecaseConfig struct {
	Code                string              `yaml:"code"`
	DbDataServiceConfig DbDataServiceConfig `yaml:"dbDataServiceConfig"`
	MqDataServiceConfig MQDataServiceConfig `yaml:"mqDataServiceConfig"`
}

type UsecasesConfig struct {
	ManageUsersUsecaseConfig ManageUsersUsecaseConfig `yaml:"manageUsersUsecaseConfig"`
}
