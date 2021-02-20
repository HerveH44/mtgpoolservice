package setting

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/go-ini/ini"
)

type Settings struct {
	App
	Database
	Server
}

type App struct {
	MTGJsonEndpoint string
	AdminPassword   string
}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Database struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
	Log      bool
	SslMode  string
}

var cfg *ini.File

var dbConnectionRegexp = regexp.MustCompile(`^(?P<Type>.*)://(?P<Username>.*):(?P<Password>.*)@(?P<Host>.*):(?P<Port>\d*)/(?P<DatabaseName>.*)`)

// GetSettings initialize the configuration instance
func GetSettings() (settings Settings) {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.GetSettings, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", &settings.App)
	mapTo("server", &settings.Server)
	mapTo("database", &settings.Database)

	// Special Heroku settings
	if envPort, ok := os.LookupEnv("PORT"); ok {
		if parseInt, err := strconv.Atoi(envPort); err == nil {
			settings.Server.HttpPort = parseInt
		}
	}

	if dbURL, ok := os.LookupEnv("DATABASE_URL"); ok {
		paramsMap := getParams(dbConnectionRegexp, dbURL)
		if val, ok := paramsMap["Username"]; ok {
			settings.Database.User = val
		}
		if val, ok := paramsMap["Password"]; ok {
			settings.Database.Password = val
		}
		if val, ok := paramsMap["Host"]; ok {
			settings.Database.Host = val
		}
		if val, ok := paramsMap["Port"]; ok {
			settings.Database.Port = val
		}
		if val, ok := paramsMap["DatabaseName"]; ok {
			settings.Database.Name = val
		}
	}

	// basic auth settings
	if password, ok := os.LookupEnv("ADMIN_PWD"); ok {
		settings.App.AdminPassword = password
	}

	// AWS
	if hostname, ok := os.LookupEnv("RDS_HOSTNAME"); ok {
		settings.Database.Host = hostname
	}
	if port, ok := os.LookupEnv("RDS_PORT"); ok {
		settings.Database.Port = port
	}
	if dbName, ok := os.LookupEnv("RDS_DB_NAME"); ok {
		settings.Database.Name = dbName
	}
	if userName, ok := os.LookupEnv("RDS_USERNAME"); ok {
		settings.Database.User = userName
	}
	if password, ok := os.LookupEnv("RDS_PASSWORD"); ok {
		settings.Database.Password = password
	}

	settings.Server.ReadTimeout = settings.Server.ReadTimeout * time.Second
	settings.Server.WriteTimeout = settings.Server.WriteTimeout * time.Second

	return
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
