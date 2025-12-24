package handles

import (
	"fmt"
	"github.com/OpenListTeam/OpenList/v4/pkg/jellyfin"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ProxyJellyfin(c *gin.Context) {

	dstPath := c.Param("path")
	dstQuery := c.Request.URL.RawQuery

	dstUrl := fmt.Sprintf("%v?%v", dstPath, dstQuery)

	jellyfinClient := jellyfin.NewClient()
	if !jellyfinClient.IsHosted() {
		c.String(406, "Jellyfin not set")
		return
	}

	backendUrl, backendErr := url.Parse(fmt.Sprintf("%v%v", jellyfinClient.GetApiHost(), dstUrl))
	if backendErr != nil {

		log.Errorf("[ProxyJellyfin]ParseUrl=[%v] error:%v", jellyfinClient.GetApiHost(), backendErr.Error())

		c.String(http.StatusGone, "Jellyfin set error")
		return
	}

	log.Infof("[ProxyJellyfin]Proxy Url=[%v]", backendUrl.String())

	backendErrHandler := func(resp http.ResponseWriter, req *http.Request, err error) {

		log.Errorf("[ProxyJellyfin]request=[%v] response error:%v", req.RequestURI, err.Error())

		resp.WriteHeader(http.StatusServiceUnavailable)
		resp.Write([]byte("Jellyfin response error"))
	}

	backendDirector := func(req *http.Request) {
		req.URL = backendUrl
		req.Header.Set("User-Agent", "Mozilla/5.5 Openlist ProxyJellyfin")
		req.Header.Set("Authorization", fmt.Sprintf("MediaBrowser Token=\"%v\"", jellyfinClient.GetApiKey()))
	}

	xProxy := &httputil.ReverseProxy{Director: backendDirector, ErrorHandler: backendErrHandler}
	xProxy.ServeHTTP(c.Writer, c.Request)
}
