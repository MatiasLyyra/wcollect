package wcollect

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func InsertData[T ClickhouseData](ctx context.Context, conn clickhouse.Conn, data []T) error {
	var tmp T
	inserStmt := fmt.Sprintf("INSERT INTO %v", tmp.table())
	batch, err := conn.PrepareBatch(ctx, inserStmt)
	if err != nil {
		return err
	}
	for _, v := range data {
		if err := batch.AppendStruct(&v); err != nil {
			return err
		}
	}
	if err := batch.Send(); err != nil {
		return err
	}
	return nil
}
