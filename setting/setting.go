package setting

import (
	"github.com/spf13/viper"
)

type Common struct {
	AppInstallPath string `mapstructure:"app_install_path"`
	AppBackupPath  string `mapstructure:"app_backup_path"`
}

type Hosts struct {
	User string   `mapstructure:"user"`
	Port int      `mapstructure:"port"`
	IPS  []string `mapstructure:"ips"`
}
type Redis struct {
	IP   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
	Pass string `mapstructure:"pass"`
}

type MySQL struct {
	IP   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
}

type App struct {
	Name  string   `mapstructure:"name"`
	Nodes []string `mapstructure:"nodes"`
}

type Apps []*App

type setting struct {
	vp *viper.Viper
}

func New() *setting {
	vp := viper.New()
	return &setting{vp: vp}
}

func (s setting) Init(file string) error {
	s.vp.SetConfigFile(file)
	s.vp.AddConfigPath(".")
	if err := s.vp.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func (s setting) SetSection(k string, v interface{}) error {
	if err := s.vp.UnmarshalKey(k, v); err != nil {
		return err
	}
	return nil
}
