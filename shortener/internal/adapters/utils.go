package adapters

import (
	"database/sql"
	"github.com/wb-go/wbf/zlog"
)

func ClosePostgresRows(rows *sql.Rows) {
	func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			zlog.Logger.Error().Err(err).Msg("error closing rows")
		}
	}(rows)
}
