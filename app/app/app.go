package app

import (
	"github.com/satorunooshie/Yatter/app/config"
	"github.com/satorunooshie/Yatter/app/dao"
)

// Dependency manager for whole application
type App struct {
	Dao dao.Dao
}

// Create dependency manager
func NewApp() (*App, error) {
	// panic if lacking something
	daoCfg := config.MySQLConfig()

	dao, err := dao.New(daoCfg)
	if err != nil {
		return nil, err
	}

	return &App{Dao: dao}, nil
}
