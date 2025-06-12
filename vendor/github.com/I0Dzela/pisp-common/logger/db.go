package logger

import (
	"context"
	"github.com/go-pg/pg/v10"
	pg9 "github.com/go-pg/pg/v9"
	log "github.com/sirupsen/logrus"
	"time"
)

type DbLogger struct{}

func (d DbLogger) BeforeQuery(c context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d DbLogger) AfterQuery(_ context.Context, q *pg.QueryEvent) error {
	fq, _ := q.FormattedQuery()
	log.Info("query: ", string(fq))
	log.Info("query - duration: ", time.Since(q.StartTime))
	return nil
}

type DbLogger9 struct{}

func (d DbLogger9) BeforeQuery(c context.Context, _ *pg9.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d DbLogger9) AfterQuery(_ context.Context, q *pg9.QueryEvent) error {
	fq, _ := q.FormattedQuery()
	log.Info("AfterQuery: ", fq)
	return nil
}
