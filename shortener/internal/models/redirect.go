package models

import "github.com/chempik1234/super-danis-library-golang/pkg/types"

// Redirect - entity representing single click on short link
//
// Aggregate by user-agent, click_at, short_url
type Redirect struct {
	ClickAt   types.DateTime
	UserAgent types.AnyText
	ShortURL  ShortURL
}

// RedirectDataList - grouped list for analytics.
//
// Representation of analytics data snapshot
type RedirectDataList struct {
	Link             *Link
	UniqueUserAgents int
	Data             []*RedirectDataListItem
}

// RedirectDataListItem - item for RedirectDataList.Data
//
// You're not supposed to create values of this type
type RedirectDataListItem struct {
	Minute         types.DateTime
	ClicksInMinute int64
	Data           []*RedirectDataListMinuteItem
}

// RedirectDataListMinuteItem - item for RedirectDataListItem.Data
//
// You're not supposed to create values of this type
type RedirectDataListMinuteItem struct {
	UserAgent types.NotEmptyText
	Clicks    int64
}
