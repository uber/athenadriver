package queryfx

import (
	"database/sql"
	"github.com/uber/athenadriver/athenareader/configfx"
	drv "github.com/uber/athenadriver/go"
	"go.uber.org/fx"
)

var Module = fx.Provide(new)

// Params defines the dependencies or inputs
type Params struct {
	fx.In

	MC configfx.MyConfig
}

// Result defines output
type Result struct {
	fx.Out

	QAD QueryAndDB
}

type QueryAndDB struct {
	DB    *sql.DB
	Query string
}

func new(p Params) (Result, error) {
	// 2. Open Connection.
	dsn := p.MC.DrvConfig.Stringify()
	db, _ := sql.Open(drv.DriverName, dsn)
	// 3. Query and print results
	qad := QueryAndDB{
		DB:    db,
		Query: p.MC.Qy,
	}
	return Result{
		QAD: qad,
	}, nil
}
