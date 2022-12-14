package setting

import (
	"log"
	"time"

	"gopkg.in/ini.v1"
)

type App struct {
	PageSize    int
	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

type SEQUENCE struct {
	Sequence int64
}

func (s *SEQUENCE) Init() {
	s.Sequence = 0
}

func (s *SEQUENCE) Get() int64 {
	s.Sequence++
	return s.Sequence
}

var DatabaseSetting = &Database{}

var cfg *ini.File
var Sequence SEQUENCE

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("web/conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}
	Sequence.Init()

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
