package SyllabusVideo

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
	crawler "github.com/vincent87720/TRC/internal/crawler"
)

//detail VideoInfoDetail(UnmarshalJSON用)
type detail struct {
	Duration string
}

//snippetData VideoSnippetData(UnmarshalJSON用)
type snippetData struct {
	Title string
}

//item VideoTtem(UnmarshalJSON用)
type item struct {
	Id             string
	Snippet        snippetData
	ContentDetails detail
}

//ytVideoInfo YoutubeVideoInfo(UnmarshalJSON用)
type ytVideoInfo struct {
	Items []item
}

//getYoutubeVideoDurationRequest
type getYTVDRequest struct {
	crawler.Request
	title         string
	duration      string
	seconds       int
	youtubeAPIKey string
	videoInfo     ytVideoInfo
}

func (ytvdreq *getYTVDRequest) setYoutubeAPIKey(youtubeAPIKey string) (err error) {
	if youtubeAPIKey != "" {
		ytvdreq.youtubeAPIKey = youtubeAPIKey
	} else {
		err = errors.WithStack(fmt.Errorf("youtubeAPIKey missing"))
		return err
	}
	return nil
}

func (ytvdreq *getYTVDRequest) getYoutubeVideoInfo(videoID string) (err error) {
	ytvdreq.SetURL("https://www.googleapis.com/youtube/v3/videos?part=contentDetails,snippet&id=" + videoID + "&key=" + ytvdreq.youtubeAPIKey)
	err = ytvdreq.SendGetRequest()
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	return nil
}

func (ytvdreq *getYTVDRequest) parseData() (err error) {
	r1, err := regexp.Compile(`([0-9]*)M`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	r2, err := regexp.Compile(`([0-9]*)S`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	err = json.Unmarshal(ytvdreq.ResponseData, &ytvdreq.videoInfo)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	if len(ytvdreq.videoInfo.Items) > 0 {
		ytvdreq.title = ytvdreq.videoInfo.Items[0].Snippet.Title
		ytvdreq.duration = ytvdreq.videoInfo.Items[0].ContentDetails.Duration
		var minInt, secInt int
		minStr := r1.FindString(ytvdreq.duration)
		secStr := r2.FindString(ytvdreq.duration)
		if minStr != "" {
			minStr = minStr[:len(minStr)-1]
			minInt, err = strconv.Atoi(minStr)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		}
		if secStr != "" {
			secStr = secStr[:len(secStr)-1]
			secInt, err = strconv.Atoi(secStr)
			if err != nil {
				err = errors.WithStack(err)
				return err
			}
		}
		ytvdreq.seconds = minInt*60 + secInt
	}
	return nil
}
