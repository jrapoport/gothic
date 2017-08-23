package dial

import (
	"github.com/netlify/gotrue/conf"
	"github.com/netlify/gotrue/models"
	"github.com/netlify/gotrue/storage"
	"github.com/netlify/gotrue/storage/mongo"
	"github.com/netlify/gotrue/storage/sql"
)

// Dial will connect to that storage engine
func Dial(config *conf.GlobalConfiguration) (storage.Connection, error) {
	if config.DB.Namespace != "" {
		models.Namespace = config.DB.Namespace
	}

	var conn storage.Connection
	var err error
	if config.DB.Driver == "mongo" {
		conn, err = mongo.Dial(config)
	} else {
		conn, err = sql.Dial(config)
	}

	if err != nil {
		return nil, err
	}

	if config.DB.Automigrate {
		if err := conn.Automigrate(); err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}
