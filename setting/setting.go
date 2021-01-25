package setting

import (
	"github.com/go-ini/ini"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

type App struct {
	PageSize  int
	PrefixUrl string

	RuntimeRootPath string

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
	Type     string
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	Log      bool
	SslMode  string
}

var DatabaseSetting = &Database{}

var cfg *ini.File

var dbConnectionRegexp = regexp.MustCompile(`^(?P<Type>.*)://(?P<Username>.*):(?P<Password>.*)@(?P<Host>.*):(?P<Port>\d*)/(?P<DatabaseName>.*)`)

// Setup initialize the configuration instance
func Setup() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.Setup, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)

	// Special Heroku settings
	if envPort, ok := os.LookupEnv("PORT"); ok {
		if parseInt, err := strconv.Atoi(envPort); err == nil {
			ServerSetting.HttpPort = parseInt
		}
	}

	if dbURL, ok := os.LookupEnv("DATABASE_URL"); ok {
		paramsMap := getParams(dbConnectionRegexp, dbURL)
		if val, ok := paramsMap["Type"]; ok {
			DatabaseSetting.Type = val
		}
		if val, ok := paramsMap["Username"]; ok {
			DatabaseSetting.User = val
		}
		if val, ok := paramsMap["Password"]; ok {
			DatabaseSetting.Password = val
		}
		if val, ok := paramsMap["Host"]; ok {
			DatabaseSetting.Host = val
		}
		if val, ok := paramsMap["Port"]; ok {
			DatabaseSetting.Port = val
		}
		if val, ok := paramsMap["DatabaseName"]; ok {
			DatabaseSetting.Name = val
		}
	}
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
}

/**
 * Parses url with the given regular expression and returns the
 * group values defined in the expression.
 *
 */
func getParams(regEx *regexp.Regexp, url string) (paramsMap map[string]string) {
	match := regEx.FindStringSubmatch(url)
	paramsMap = make(map[string]string)
	for i, name := range regEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
