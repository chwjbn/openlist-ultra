package jellyfin

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type JellyfinClient struct {
}

func NewClient() *JellyfinClient {
	return &JellyfinClient{}
}

func (j *JellyfinClient) GetApiHost() string {

	var dstUrl string
	dstUrl = os.Getenv("JELLYFIN_API_HOST")
	return dstUrl

}

func (j *JellyfinClient) GetApiKey() string {

	var dstKey string
	dstKey = os.Getenv("JELLYFIN_API_KEY")
	return dstKey

}

func (j *JellyfinClient) GetApiRootFolder() string {

	var dstData string
	dstData = os.Getenv("JELLYFIN_API_ROOT_FOLDER")
	return dstData
}

func (j *JellyfinClient) IsHosted() bool {

	log.Infof("ApiHost=[%v]", j.GetApiHost())

	return len(j.GetApiHost()) > 0 && len(j.GetApiKey()) > 0 && len(j.GetApiRootFolder()) > 0

}

func (j *JellyfinClient) GetMediaStreamUrl(path string) string {

	log.Infof("GetMediaStreamUrl path=[%v]", path)

	var dstUrl string

	if !j.IsHosted() {
		return dstUrl
	}

	dstId := j.GetItemId(path)

	log.Infof("GetItemId path=[%v] dstId=[%v]", path, dstId)

	if len(dstId) < 1 {
		return dstUrl
	}

	dstUrl = fmt.Sprintf("/px/jellyfin/Videos/%v/main.m3u8?VideoCodec=h264&AudioCodec=aac&AudioStreamIndex=1&AudioSampleRate=48000&MaxFramerate=60&SegmentContainer=mp4&MinSegments=1&BreakOnNonKeyFrames=False", dstId)

	return dstUrl

}

func (j *JellyfinClient) GetItemId(path string) string {

	var dstId string

	srcPath := filepath.Clean(path)
	if len(srcPath) < 1 {
		return dstId
	}

	srcFullPath := filepath.Clean(filepath.Join(j.GetApiRootFolder(), srcPath))

	log.Infof("GetItemId srcFullPath=[%v]", srcFullPath)

	var parentId string

	for {

		xApiPathUrl := "/Items?locationTypes=FileSystem&fields=Path&parentId=" + parentId
		xApiData := j.doGet(xApiPathUrl)

		if len(xApiData) < 1 {
			break
		}

		var xApiResp JellyfinRespItems
		xApiErr := json.Unmarshal([]byte(xApiData), &xApiResp)
		if xApiErr != nil {
			break
		}

		if len(xApiResp.Items) < 1 {
			break
		}

		var curItemId string

		for _, item := range xApiResp.Items {

			if strings.HasPrefix(srcFullPath, item.Path) {

				if !item.IsFolder {
					dstId = item.Id
				}

				curItemId = item.Id
				break
			}
		}

		//没有匹配到
		if len(curItemId) < 1 {
			break
		}

		//匹配到了最终的
		if len(dstId) > 0 {
			break
		}

		parentId = curItemId
	}

	return dstId

}

func (j *JellyfinClient) doGet(pathUrl string) string {

	var dstData string

	var xErr error
	var xApiUrl string

	xApiUrl = fmt.Sprintf("%v%v", j.GetApiHost(), pathUrl)

	xClient := j.getHttpClient()

	xReq, xErr := http.NewRequest("GET", xApiUrl, nil)
	if xErr != nil {
		return dstData
	}

	xReq.Header.Set("User-Agent", "Mozilla/5.5 Openlist HttpClient")
	xReq.Header.Set("Authorization", fmt.Sprintf("MediaBrowser Token=\"%v\"", j.GetApiKey()))

	xResp, xErr := xClient.Do(xReq)
	if xErr != nil {
		return dstData
	}

	xRespBody, xErr := io.ReadAll(xResp.Body)
	if xErr != nil {
		return dstData
	}

	dstData = string(xRespBody)

	return dstData

}

func (j *JellyfinClient) getHttpClient() *http.Client {

	xClient := &http.Client{}
	xClient.Timeout = 60 * time.Second

	return xClient
}
