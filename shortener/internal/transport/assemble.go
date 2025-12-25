package transport

import (
	"fmt"
	"github.com/wb-go/wbf/ginext"
)

// AssembleRouter is the function you'd call in `main.go` to get THE app router
func AssembleRouter(shortenerHandler *ShortenerHandler) *ginext.Engine {
	router := ginext.New("release")

	router.POST("/shorten", shortenerHandler.CreateLink)
	router.GET(fmt.Sprintf("/s/:%s", shortLinkParam), shortenerHandler.RedirectLink)
	router.GET(fmt.Sprintf("/analytics/:%s", shortLinkParam), shortenerHandler.AnalyticsLink)

	return router
}
