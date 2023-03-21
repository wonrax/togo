package togo

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"go.uber.org/zap"
)

func DbInsert(_data any, table string) (result sql.Result, err error) {
	if _data == nil {
		err = errors.New("data is nil")
		return
	}

	data, ok := _data.(map[string]any)
	if !ok {
		Log.Debug("Data is not a map, trying to unmarshal it", zap.Any("data", _data))
		inrec, err_ := json.Marshal(_data)
		if err_ != nil {
			return nil, err_
		}
		err_ = json.Unmarshal(inrec, &data)
		if err != nil {
			return nil, err_
		}
	}

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

// TODO handle wildcards in select statement
func DbFind(data map[string]interface{}, cols []string, table string) (results []map[string]interface{}, err error) {
	args := make([]interface{}, 0, len(data))
	keys := make([]string, 0, len(data))
	whereStmt := ""
	selectedColsStmt := ""

	// make sure that the order of the keys
	// is the same as the order of the values
	for k := range data {
		keys = append(keys, k)
	}

	for i, k := range keys {
		args = append(args, data[k])
		if i == 0 {
			whereStmt += fmt.Sprintf("%s = ?", k)
		} else {
			whereStmt += fmt.Sprintf(" AND %s = ?", k)
		}
	}

	for i := range cols {
		if i == 0 {
			selectedColsStmt += cols[i]
		} else {
			selectedColsStmt += ", " + cols[i]
		}
	}

	sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s;", selectedColsStmt, table, whereStmt)

	Log.Debug("Executing SQL", zap.String("sql", sql), zap.Any("args", args))

	rows, err := Db.Query(sql, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		values := make([]interface{}, len(cols))
		for i := range values {
			values[i] = new(interface{})
		}
		err = rows.Scan(values...)
		if err != nil {
			return
		}
		row := make(map[string]interface{})
		for i, col := range cols {
			row[col] = *(values[i]).(*interface{})
		}
		results = append(results, row)
	}

	return
}
