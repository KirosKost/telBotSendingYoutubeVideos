package youtube

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	
)

const (
	YOUTUBE_SEARCH_URL = "https://www.googleapis.com/youtube/v3/search"
	YOUTUBE_API_TOKEN  = "<YOUTUBE_API_TOKEN>"
	YOUTUBE_VIDEO_URL  = "https://www.youtube.com/watch?v="
)

//GET https://youtube.googleapis.com/youtube/v3/search?part=id&channelId=UCNau7mXLYeaz_LGymncoWNw&maxResults=1&order=date&key=[YOUR_API_KEY] HTTP/1.1

//Authorization: Bearer [YOUR_ACCESS_TOKEN]
//Accept: application/json

func GetLastVideo(channelUrl string) (string, error) {
	items, err := retrieveVideos(channelUrl)
	if err != nil {
		return "", err
	}
	if len(items) < 1 {
		return "", errors.New("Video not found!")
	}
	return YOUTUBE_VIDEO_URL + items[0].Id.VideoId, nil
}

func retrieveVideos(channelUrl string) ([]Item, error) {
	req, err := makeRequest(channelUrl, 1)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Items, nil
}

func makeRequest(channelUrl string, maxResult int) (*http.Request, error) {
	lastSlashIndex := strings.LastIndex(channelUrl, "/")
	channelId := channelUrl[lastSlashIndex+1:]

	req, err := http.NewRequest("GET", YOUTUBE_SEARCH_URL, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Add("part", "id")
	query.Add("channelId", channelId)
	query.Add("maxResult", strconv.Itoa(maxResult))
	query.Add("order", "date")
	query.Add("key", YOUTUBE_API_TOKEN)
	req.URL.RawQuery = query.Encode()
	fmt.Println(req.URL.String())
	return req, nil
}
