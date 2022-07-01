package trace

import (
	"strings"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"gorm.io/gorm"
)

const (
	dbTableKey        = attribute.Key("db.sql.table")
	dbRowsAffectedKey = attribute.Key("db.rows_affected")
	dbOperationKey    = semconv.DBOperationKey
	dbStatementKey    = semconv.DBStatementKey
)

func dbTable(name string) attribute.KeyValue {
	return dbTableKey.String(name)
}

func dbStatement(stmt string) attribute.KeyValue {
	return dbStatementKey.String(stmt)
}

func dbCount(n int64) attribute.KeyValue {
	return dbRowsAffectedKey.Int64(n)
}

func dbOperation(op string) attribute.KeyValue {
	return dbOperationKey.String(op)
}

func operationForQuery(query, op string) string {
	if op != "" {
		return op
	}
	return strings.ToUpper(strings.Split(query, " ")[0])
}

func extractQuery(tx *gorm.DB) string {
	return tx.Statement.SQL.String()
	// if shouldOmit, _ := tx.Statement.Context.Value(omitVarsKey).(bool); shouldOmit {
	// 	return tx.Statement.SQL.String()
	// }
	// return tx.Dialector.Explain(tx.Statement.SQL.String(), tx.Statement.Vars...)
}
