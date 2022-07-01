package trace

import (
	"fmt"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	eventMaxSize = 250
	maxChunks    = 4
)

func (op *OtelPlugin) before(operation string) func(tx *gorm.DB) {
	return func(tx *gorm.DB) {
		tx.Statement.Context, _ = otel.Tracer("yggdrasil").
			Start(tx.Statement.Context, op.spanName(tx, operation), oteltrace.WithSpanKind(oteltrace.SpanKindClient))
	}
}

func (op *OtelPlugin) after(operation string) gormHookFunc {
	return func(tx *gorm.DB) {
		span := oteltrace.SpanFromContext(tx.Statement.Context)
		if !span.IsRecording() {
			// skip the reporting if not recording
			return
		}
		defer span.End()

		span.SetName(op.spanName(tx, operation))
		// Error
		if tx.Error != nil {
			span.SetStatus(codes.Error, tx.Error.Error())
		}

		// extract the db operation
		query := strings.ToValidUTF8(extractQuery(tx), "")

		// If query is longer then max size log it as chunked event, otherwise log it in attribute
		if len(query) > eventMaxSize {
			chunkBy(query, eventMaxSize, span.AddEvent)
		} else {
			span.SetAttributes(dbStatement(query))
		}

		operation = operationForQuery(query, operation)
		if tx.Statement.Table != "" {
			span.SetAttributes(dbTable(tx.Statement.Table))
		}

		span.SetAttributes(
			dbOperation(operation),
			dbCount(tx.Statement.RowsAffected),
		)
	}
}

func chunkBy(val string, size int, callback func(string, ...oteltrace.EventOption)) {
	if len(val) > maxChunks*size {
		return
	}

	for i := 0; i < maxChunks*size; i += size {
		end := len(val)
		if end > size {
			end = size
		}
		callback(val[0:end])
		if end > len(val)-1 {
			break
		}
		val = val[end:]
	}
}

func (op *OtelPlugin) spanName(tx *gorm.DB, operation string) string {
	operation = operationForQuery(extractQuery(tx), operation)
	target := op.dbName
	if target == "" {
		target = tx.Dialector.Name()
	}

	if tx.Statement != nil && tx.Statement.Table != "" {
		target += "." + tx.Statement.Table
	}

	return fmt.Sprintf("%s %s", operation, target)
}
