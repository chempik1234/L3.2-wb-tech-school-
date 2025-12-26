package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/adapters"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"github.com/chempik1234/super-danis-library-golang/pkg/types"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"golang.org/x/sync/errgroup"
	"strings"
	"sync"
	"time"
)

// StoragePostgresRepo - adapter for ports.StoragePostgresRepo
//
// PostgresSQL
type StoragePostgresRepo struct {
	db       *dbpg.DB
	strategy retry.Strategy
}

// NewStoragePostgresRepo creates a new StoragePostgresRepo
func NewStoragePostgresRepo(db *dbpg.DB, retryStrategy retry.Strategy) *StoragePostgresRepo {
	return &StoragePostgresRepo{db: db, strategy: retryStrategy}
}

// SaveRedirectsBatch - save a bunch of redirects to DB
//
// batching is fast, batching is everything!
func (s *StoragePostgresRepo) SaveRedirectsBatch(ctx context.Context, redirectsToSave []*models.Redirect) error {
	// step 1. make values slice - [`('1','2','3')`, `('1','2','3')`, ...]
	values := make([]string, len(redirectsToSave))
	for i, r := range redirectsToSave {
		values[i] = fmt.Sprintf("('%s','%s','%s')",
			r.ShortURL.String(),
			r.ClickAt.Value().Format(time.RFC3339),
			r.UserAgent.String())
	}
	// step 2. join the slice [(), ()] -> "(),()"
	valuesStr := strings.Join(values, ",")

	// step 3. query with sprintf -> I <3 SQL injections
	// user can use UserAgent to make sql injections
	query := fmt.Sprintf(`INSERT INTO redirects (short_url, click_at, user_agent) VALUES %s`, valuesStr)

	result, err := s.db.ExecWithRetry(ctx, s.strategy, query)
	if err != nil {
		return fmt.Errorf("error saving batch (%d elements): %w", len(redirectsToSave), err)
	}

	// region check rowsAffected int64 == len redirectsToSave
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if int64(len(redirectsToSave)) != rowsAffected {
		return fmt.Errorf("not enough rows inserted in batch: %d / %d", rowsAffected, len(redirectsToSave))
	}
	//endregion

	return nil
}

// GetAnalytics - get aggregated analytics from inside the DB
//
// # Group By is better than a local Golang function
//
// RESULT IS HALF EMPTY because it can't query LINK model
//
// LINK FIELD IS EMPTY QUERY IT YOURSELF with ShortenerStorageRepository
func (s *StoragePostgresRepo) GetAnalytics(ctx context.Context, shortLink models.ShortURL) (*models.RedirectDataList, error) {
	eg := &errgroup.Group{}

	var uniqueAgentsCount int
	var data []*models.RedirectDataListItem

	tx, err := s.db.BeginTxWithRetry(ctx, s.strategy, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
		ReadOnly:  true,
	})
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}

	eg.Go(func() error {
		var err error
		uniqueAgentsCount, err = s.getUniqueUserAgentCount(ctx, tx, shortLink)
		return err
	})
	eg.Go(func() error {
		var err error
		data, err = s.getRedirectDataList(ctx, tx, shortLink)
		return err
	})

	err = eg.Wait()
	if err != nil {
		return nil, fmt.Errorf("error getting analytics: %w", err)
	}

	return &models.RedirectDataList{
		Link:             nil,
		UniqueUserAgents: uniqueAgentsCount,
		Data:             data,
	}, nil
}

func (s *StoragePostgresRepo) getUniqueUserAgentCount(ctx context.Context, tx *sql.Tx, link models.ShortURL) (int, error) {
	query := `SELECT COUNT(DISTINCT user_agent) FROM redirects WHERE short_url = $1`

	result := 0

	// since we use Tx, we've got to use custom retries
	err := retry.Do(func() error {
		row := tx.QueryRowContext(ctx, query, link.String())
		err := row.Scan(&result)
		if err != nil {
			return fmt.Errorf("error during scan: %w", err)
		}
		return nil
	}, s.strategy)

	if err != nil {
		return result, fmt.Errorf("error querying count user_agent: %w", err)
	}

	return result, nil
}

func (s *StoragePostgresRepo) getRedirectDataList(ctx context.Context, tx *sql.Tx, link models.ShortURL) ([]*models.RedirectDataListItem, error) {
	// return {
	// "unique_user_agent": ...
	// "data": [
	//    {
	//        "minute": ...
	//        "clicks_in_minute": ...
	//        "data": [
	//            "user_agent": ...
	//            "clicks": ...
	//        ]
	//    }
	// ]

	// instead of nested tables, we make 1 sorted denormalized table

	// | minute | user_agent | clicks_sum
	// |--------|            |
	// | 20:00  | useragent1 | 10
	// | 20:00  | useragent2 | 10
	// |--------|            |
	// | 20:03  | useragent1 | 10      <-- blocks by minutes
	// | 20:03  | useragent2 | 10           there could be hours, days, months, but here we've got minutes
	// | 20:03  | useragent3 | 10
	// | 20:03  | useragent4 | 10
	// |--------|            |
	// | 20:10  | useragent3 | 10
	// | 20:10  | useragent4 | 10

	// then we need:
	//
	// TABLE (minute user_agent clicks)
	// order by minute, user_agent
	//
	// we'll count clicks_in_minute in progress
	//
	// check further comments

	query := `SELECT date_trunc('minute', click_at) as minute, user_agent, count(user_agent) as clicks
              FROM redirects
              WHERE short_url = $1
              GROUP BY date_trunc('minute', click_at), user_agent
              ORDER BY minute DESC`

	var rows *sql.Rows

	// since we use Tx, we've got to use custom retries
	err := retry.Do(func() error {
		var err error
		rows, err = tx.QueryContext(ctx, query, link.String())
		if err != nil {
			return fmt.Errorf("error querying rows: %w", err)
		}
		return nil
	}, s.strategy)
	if err != nil {
		return nil, fmt.Errorf("error querying rows: %w", err)
	}

	defer adapters.ClosePostgresRows(rows)

	dataList := make([]*models.RedirectDataListItem, 0)

	// workflow:
	//
	// listForCurrentMinute=[item]
	// listForCurrentMinute=[item, item]
	//
	// new minute! {
	//   dataList = append(dataList, listForCurrentMinute)
	//   listForCurrentMinute=[]
	// }
	//
	// listForCurrentMinute=[item2]
	// listForCurrentMinute=[item2, item2]
	// listForCurrentMinute=[item2, item2, item2]
	//
	// end! { dataList = append(dataList, listForCurrentMinute) }  <- final append after "for" loop
	var listForCurrentMinute []*models.RedirectDataListMinuteItem

	var currentMinute time.Time

	var rowMinute time.Time
	var rowUserAgent string
	var rowClicks int64
	var clicksInCurrentMinute int64

	// we should print error 'some user agents are empty'
	// only once, so we use sync.Once
	userAgentErrorOnce := &sync.Once{}

	for rows.Next() {
		// row: minute, user_agent, clicks
		if err = rows.Scan(&rowMinute, &rowUserAgent, &rowClicks); err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		// exec on row 0 or after every new minute block
		if listForCurrentMinute == nil || rowMinute != currentMinute {
			if listForCurrentMinute != nil {
				dataList = append(dataList, &models.RedirectDataListItem{
					Minute:         types.NewDateTime(currentMinute),
					ClicksInMinute: clicksInCurrentMinute,
					Data:           listForCurrentMinute,
				})
			}
			clicksInCurrentMinute = 0
			listForCurrentMinute = make([]*models.RedirectDataListMinuteItem, 0)
			currentMinute = rowMinute
		}

		// we do that before "continue" because we might lose some clicks that are made but their user agents are incorrect
		clicksInCurrentMinute += rowClicks

		userAgentValid, err := types.NewNotEmptyText(rowUserAgent)
		if err != nil {
			userAgentErrorOnce.Do(func() { zlog.Logger.Error().Msg("database validation error: some user agents are empty") })
			continue
		}

		listForCurrentMinute = append(listForCurrentMinute, &models.RedirectDataListMinuteItem{
			UserAgent: userAgentValid,
			Clicks:    rowClicks,
		})
	}

	// repeat what's in cycle
	if listForCurrentMinute != nil {
		dataList = append(dataList, &models.RedirectDataListItem{
			Minute:         types.NewDateTime(currentMinute),
			ClicksInMinute: clicksInCurrentMinute,
			Data:           listForCurrentMinute,
		})
	}

	return dataList, nil
}
