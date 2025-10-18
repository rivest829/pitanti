package conf

type ServerType string

const (
	ServerGame    ServerType = "game"
	ServerGateway ServerType = "gateway"
)

const (
	CtxServerId = "serverId"
	CtxRoleId   = "roleId"
)

type Env string

const (
	EnvDev Env = "dev"
	EnvPro Env = "pro"
)
