package commons

type Config struct {
	Servers    []string `mapstructure:"servers"`
	RateLimit  int      `mapstructure:"rate_limit"`
	LogEnabled bool     `mapstructure:"log_enabled"`
	LogFile    string   `mapstructure:"log_file"`
	Port       int      `mapstructure:"port"`
	Database   struct {
		Redis struct {
			Db       int    `mapstructure:"db"`
			Host     string `mapstructure:"host"`
			Port     int    `mapstructure:"port"`
			Password string `mapstructure:"password"`
		} `mapstructure:"redis"`
	} `mapstructure:"database"`
}
