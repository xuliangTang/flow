package conf

var Config *Configuration

type DbConfig struct {
	Host        string
	Port        int32
	Username    string
	Password    string
	Name        string
	MaxIdleConn int
	MaxOpenConn int
}

type RedisConfig struct {
	Host        string
	Port        int32
	Password    string
	Db          int
	PoolSize    int
	MinIdleConn int
}

type ServerConfig struct {
	AppPath string
}

type MiniProgramConfig struct {
	AppId  string
	Secret string
}

type JwtConfig struct {
	Exp int
	Key string
}

type Configuration struct {
	Database    *DbConfig
	Redis       *RedisConfig
	Server      *ServerConfig
	MiniProgram *MiniProgramConfig
	Jwt         *JwtConfig
}
