package globals

import (
	"monitoring/config"
	"monitoring/pkg/postgres"
)

var GlobalPG postgres.IPostgres
var GlobalConfig config.Config
