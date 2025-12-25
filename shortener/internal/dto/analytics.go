package dto

import (
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"time"
)

// AnalyticsBody - DTO for aggregated analytics:
//
//  1. redirect count
//  2. source/short URL
//  3. all redirects by minutes (1 for 20:01, 2 for 20:05, zeros SKIPPED)
//  4. clicks by user agents for each minute
//
// Body example:
//
//	{
//	  "source_url": "https://ya.ru",
//	  "short_url": "<short_url>",
//	  "total_redirects": 2,
//	  "data": [
//	    {
//	      "minute_timestamp": 1766563140,
//	      "clicks_in_minute": 120469101,
//	      "data": [
//	        {
//	          "user_agent": "...",
//	          "clicks": 1
//	        }
//	      ]
//	    }
//	  ]
//	}
type AnalyticsBody struct {
	SourceURL        string              `json:"source_url"`
	ShortURL         string              `json:"short_url"`
	TotalRedirects   int                 `json:"total_redirects"`
	UniqueUserAgents int                 `json:"unique_user_agents"`
	Data             []analyticsDataItem `json:"data"`
}

type analyticsDataItem struct {
	Minute         string          `json:"minute_timestamp"`
	ClicksInMinute int64           `json:"clicks_in_minute"`
	Data           []userAgentItem `json:"data"`
}

type userAgentItem struct {
	UserAgent string `json:"user_agent"`
	Clicks    int64  `json:"clicks"`
}

// AnalyticsBodyFromDataList - serialize models.RedirectDataList into AnalyticsBody
func AnalyticsBodyFromDataList(redirects *models.RedirectDataList) AnalyticsBody {
	dataList := make([]analyticsDataItem, len(redirects.Data))

	for i, minute := range redirects.Data {
		minuteDataList := make([]userAgentItem, len(minute.Data))
		for j, userAgent := range minute.Data {
			minuteDataList[j] = userAgentItem{
				UserAgent: userAgent.UserAgent.String(),
				Clicks:    userAgent.Clicks,
			}
		}

		dataList[i] = analyticsDataItem{
			Minute:         minute.Minute.Value().Format(time.RFC3339),
			ClicksInMinute: minute.ClicksInMinute,
			Data:           minuteDataList,
		}
	}

	return AnalyticsBody{
		SourceURL:        redirects.Link.SourceURL.String(),
		ShortURL:         redirects.Link.ShortURL.String(),
		UniqueUserAgents: redirects.UniqueUserAgents,
		TotalRedirects:   len(redirects.Data),
		Data:             dataList,
	}
}
