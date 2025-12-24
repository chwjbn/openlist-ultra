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

func (j *JellyfinClient) IsHosted() bool {

	log.Infof("ApiHost=[%v]", j.getApiHost())

	return len(j.getApiHost()) > 0 && len(j.getApiKey()) > 0 && len(j.getApiRootFolder()) > 0

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

	dstUrl = fmt.Sprintf("%v/Videos/%v/stream.flv?static=true", j.getApiHost(), dstId)

	return dstUrl

}

func (j *JellyfinClient) GetItemId(path string) string {

	var dstId string

	srcPath := filepath.Clean(path)
	if len(srcPath) < 1 {
		return dstId
	}

	var srcPathList []string
	for {

		basePath := filepath.Base(srcPath)
		dirPath := filepath.Dir(srcPath)

		if len(basePath) < 1 {
			break
		}

		if strings.EqualFold(basePath, dirPath) {
			break
		}

		srcPathList = append(srcPathList, basePath)

		srcPath = filepath.Clean(dirPath)
		if len(srcPath) < 1 {
			break
		}
	}

	if len(srcPathList) < 1 {
		return dstId
	}

	srcPathList = append(srcPathList, j.getApiRootFolder())

	log.Infof("GetItemId srcPathList=[%v]", strings.Join(srcPathList, "|"))

	var parentId string

	var srcPathListSize int
	srcPathListSize = len(srcPathList)

	for i := 0; i < srcPathListSize; i++ {

		xApiPathUrl := "/Items?parentId=" + parentId
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

		var curPathName string
		curPathName = srcPathList[srcPathListSize-1-i]

		//log.Infof("curPathName=[%v]", curPathName)

		var curItemId string

		for _, item := range xApiResp.Items {

			itemName := curPathName

			if !item.IsFolder {
				ext := filepath.Ext(curPathName)
				if len(ext) > 0 {
					itemName = curPathName[:len(curPathName)-len(ext)]
				}
			}

			if strings.EqualFold(item.Name, itemName) {

				if !item.IsFolder {
					dstId = item.Id
				}

				curItemId = item.Id
				break
			}
		}

		//log.Infof("curItemId=[%v]", curItemId)

		if len(curItemId) < 1 {
			break
		}

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

	xApiUrl = fmt.Sprintf("%v%v", j.getApiHost(), pathUrl)

	xClient := j.getHttpClient()

	xReq, xErr := http.NewRequest("GET", xApiUrl, nil)
	if xErr != nil {
		return dstData
	}

	xReq.Header.Set("User-Agent", "Mozilla/5.5 Openlist HttpClient")
	xReq.Header.Set("Authorization", fmt.Sprintf("MediaBrowser Token=\"%v\"", j.getApiKey()))

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

func (j *JellyfinClient) getApiHost() string {

	var dstUrl string

	dstUrl = os.Getenv("JELLYFIN_API_HOST")

	return dstUrl

}

func (j *JellyfinClient) getApiKey() string {

	var dstKey string

	dstKey = os.Getenv("JELLYFIN_API_KEY")

	return dstKey

}

func (j *JellyfinClient) getApiRootFolder() string {

	var dstData string

	dstData = os.Getenv("JELLYFIN_API_ROOT_FOLDER")

	return dstData
}
