package togo

import (
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

func DbInsert(data map[string]interface{}, table string) (result sql.Result, err error) {
	args := make([]interface{}, 0, len(data))
	keys := make([]string, 0, len(data))
	valPlaceholder := ""
	columns := ""

	// make sure that the order of the keys
	// is the same as the order of the values
	for k := range data {
		keys = append(keys, k)
	}

	for i, k := range keys {
		args = append(args, data[k])
		if i == 0 {
			valPlaceholder += "?"
			columns += k
		} else {
			valPlaceholder += ", ?"
			columns += ", " + k
		}
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", table, columns, valPlaceholder)

	Log.Debug("Executing SQL", zap.String("sql", sql), zap.Any("args", args))
	result, err = Db.Exec(sql, args...)

	return
}
