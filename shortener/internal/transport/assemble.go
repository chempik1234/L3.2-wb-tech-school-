package transport

import "github.com/wb-go/wbf/ginext"

// AssembleRouter is the function you'd call in `main.go` to get THE app router
func AssembleRouter(shortenerHandler *ShortenerHandler) *ginext.Engine {
	router := ginext.New("release")

	router.POST("/shorten", shortenerHandler.CreateLink)
	// TODO: router.GET("/s/:short_url", shortenerHandler.RedirectLink)
	// TODO: router.DELETE("/analytics/:short_url", shortenerHandler.AnalyticsLink)

	return router
}
