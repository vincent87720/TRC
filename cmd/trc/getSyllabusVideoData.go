package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

// GetSyllabusVideo 下載數位課綱資料(建立新檔案)
// Goroutine interface for GUI
// For example:
// 	errChan := make(chan error, 2)
// 	exitChan := make(chan string, 2)
// 	defer close(errChan)
// 	defer close(exitChan)
//
// 	var outputFile file
//
// 	outputFile.setFile(filepath.ToSlash(INITPATH+"/output/"), fi.Year+fi.Semester+"數位課綱.xlsx", "工作表")
//
// 	go GetSyllabusVideo(errChan, exitChan, fi.Year, fi.Semester, fi.Key, outputFile)
//
// Loop:
// 	for {
// 		select {
// 		case err := <-errChan:
// 			Error.Printf("%+v\n", err)
// 		case <-exitChan:
// 			break Loop
// 		}
// 	}
func GetSyllabusVideo(progChan chan int, academicYear string, semester string, youtubeAPIKey string, outputFile file) {

	var dsvf downloadSVFile
	var svlreq getSVLRequest
	svlreq.svXi = make(map[string][]syllabusVideo)
	svlreq.newRequest()
	svlreq.setURL("http://syl.dyu.edu.tw/sl_cour_time.php?itimestamp=" + strconv.FormatInt(time.Now().Unix(), 10))
	// err = svlreq.setURLValues(f.academicYear, f.semester, "'1'", "'1'")
	err := svlreq.setURLValues(academicYear, semester, "'1','2','3','4','5','6','7'", "'1','2','3','4','N','5','6','7','8','9','A','B','C','D','E'")
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = svlreq.sendPostRequest()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = svlreq.parseHTML()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	svlreq.findNode(svlreq.htmlNode)

	err = svlreq.getVideoID()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = svlreq.getVideoInfo(youtubeAPIKey)
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dsvf.transportToSlice(&svlreq)
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dsvf.exportDataToExcel(outputFile)
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1

	return
}

// AppendSyllabusVideo 下載數位課綱資料(合併到舊檔案)
// Goroutine interface for GUI
// For example:
// 	errChan := make(chan error, 2)
// 	exitChan := make(chan string, 2)
// 	defer close(errChan)
// 	defer close(exitChan)
//
// 	var inputFile file
// 	var outputFile file
//
// 	inputFile.setFile(inputPathXi[1], inputPathXi[2], fi.InputSheet)
// 	outputFile.setFile(filepath.ToSlash(INITPATH+"/output/"), fi.Year+fi.Semester+"數位課綱.xlsx", "工作表")
//
// 	inputPathXi := r.FindStringSubmatch(fi.InputPath)
//
// 	go AppendSyllabusVideo(errChan, exitChan, fi.Year, fi.Semester, fi.Key, inputFile, outputFile)
//
// Loop:
// 	for {
// 		select {
// 		case err := <-errChan:
// 			Error.Printf("%+v\n", err)
// 		case <-exitChan:
// 			break Loop
// 		}
// 	}
func AppendSyllabusVideo(progChan chan int, academicYear string, semester string, youtubeAPIKey string, inputFile file, outputFile file) {

	dsvf := downloadSVFile{
		file: inputFile,
	}

	err := dsvf.readRawData()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dsvf.findCol("科目序號", &dsvf.cidCol)
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dsvf.fillSliceLength(len(dsvf.firstRow))
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1

	var svlreq getSVLRequest
	svlreq.svXi = make(map[string][]syllabusVideo)
	svlreq.newRequest()
	svlreq.setURL("http://syl.dyu.edu.tw/sl_cour_time.php?itimestamp=" + strconv.FormatInt(time.Now().Unix(), 10))
	// err = svlreq.setURLValues(f.academicYear, f.semester, "'1'", "'1'")
	err = svlreq.setURLValues(academicYear, semester, "'1','2','3','4','5','6','7'", "'1','2','3','4','N','5','6','7','8','9','A','B','C','D','E'")
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = svlreq.sendPostRequest()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = svlreq.parseHTML()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	svlreq.findNode(svlreq.htmlNode)

	err = svlreq.getVideoID()
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = svlreq.getVideoInfo(youtubeAPIKey)
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dsvf.appendVideoInfo(&svlreq)
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = dsvf.exportDataToExcel(outputFile)
	if err != nil {
		Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1

	return
}

//setURLValues 設定發送(數位課綱影片連結)請求的參數
func (svlreq *getSVLRequest) setURLValues(academicYear string, semester string, day string, lesson string) (err error) {
	yearInt, err := strconv.Atoi(academicYear)
	if err != nil || yearInt < 0 {
		fmt.Printf("無法解析\"%s\"，請輸入合法的年分", academicYear)
		err = errors.WithStack(fmt.Errorf("Incorrect year value"))
		return err
	}

	svlreq.values.Add("smye", strconv.Itoa(yearInt))
	svlreq.academicYear = academicYear

	semesterInt, err := strconv.Atoi(semester)
	if err != nil || semesterInt < 0 {
		fmt.Printf("無法解析\"%s\"，請輸入合法的學期", semester)
		err = errors.WithStack(fmt.Errorf("Incorrect year value"))
		return err
	}

	svlreq.values.Add("smty", strconv.Itoa(semesterInt))
	svlreq.semester = semester

	svlreq.values.Add("str_time", day+"sec"+lesson)
	return nil
}

//findNode 尋找class為row的標籤的子標籤，依照欄位放入svlreq.svXi
func (svlreq *getSVLRequest) findNode(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "div" {
		for _, a := range n.Attr {
			if a.Key == "class" && a.Val == "row" {
				var cs syllabusVideo
				var cID string
				for child := n.FirstChild; child != nil; child = child.NextSibling {
					for _, a1 := range child.Attr {
						if a1.Key == "class" {
						}
						if child.FirstChild != nil {
							if a1.Val == "td1" {
								if child.FirstChild != nil {
									cs.grade = child.FirstChild.Data[:1]
									cs.class = child.FirstChild.Data[2:]
									break
								}
							}
							if a1.Val == "td2" {
								if child.FirstChild != nil {
									cs.credit = child.FirstChild.Data[:1]
									cs.chooseSelect = child.FirstChild.Data[2:]
									break
								}
							}
							if a1.Val == "td3" {
								if child.FirstChild != nil {
									cs.courseID = child.FirstChild.Data
									cID = child.FirstChild.Data
									break
								}
							}
							if a1.Val == "td4" {
								if child.FirstChild != nil {
									cs.courseName = child.FirstChild.Data
									break
								}
							}
							if a1.Val == "td5" {
								if child.FirstChild != nil {
									cs.teacher.teacherName = child.FirstChild.Data
									break
								}
							}
							if a1.Val == "td7" {
								if child.FirstChild != nil {
									cs.timeNPlace = child.FirstChild.Data
									break
								}
							}
							if a1.Val == "td8" {
								if child.FirstChild != nil {
									cs.remark = child.FirstChild.Data
									break
								}
							}
							if a1.Val == "td9" {
								if child.LastChild != nil {
									for _, a2 := range child.LastChild.Attr {
										if a2.Key == "href" {
											cs.videoURL = a2.Val
											break
										}
									}
								}

							}
						}

					}
				}
				svlreq.svXi[cID] = append(svlreq.svXi[cID], cs)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		svlreq.findNode(c)
	}
}

//getVideoID 從YoutubeURL中提取VideoID
func (svlreq *getSVLRequest) getVideoID() (err error) {
	r1, err := regexp.Compile(`(\/watch\?v=|youtu.be\/)...........`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	for id, svXi := range svlreq.svXi {
		for index, sv := range svXi {
			if sv.videoURL != "javascript:alert(\"尚未提供連結\")" {
				vid := r1.FindString(sv.videoURL)
				if vid != "" {
					svlreq.svXi[id][index].videoID = vid[9:]
				}
			}
		}
	}
	return nil
}

//getVideoInfo 發送請求給YoutubeAPI取得影片資訊
func (svlreq *getSVLRequest) getVideoInfo(youtubeAPIKey string) (err error) {
	var req getYTVDRequest
	err = req.setYoutubeAPIKey(youtubeAPIKey)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	for id, svXi := range svlreq.svXi {
		for index, sv := range svXi {
			if sv.videoID != "" {
				err = req.getYoutubeVideoInfo(sv.videoID)
				if err != nil {
					err = errors.WithStack(err)
					return err
				}
				err = req.parseData()
				if err != nil {
					err = errors.WithStack(err)
					return err
				}
				svlreq.svXi[id][index].videoTitle = req.title
				svlreq.svXi[id][index].videoSeconds = req.seconds
				svlreq.svXi[id][index].videoDuration = req.duration
			}
		}
	}
	return nil
}

//transportToSlice 將教師資料放入dtFile的newDataRows中，以便使用exportDataToExcel方法輸出
func (dsvFile *downloadSVFile) transportToSlice(svlreq *getSVLRequest) (err error) {
	dsvFile.newDataRows = make([][]string, 0)
	dsvFile.newDataRows = append(dsvFile.newDataRows, []string{"學年度", "學期", "年", "班", "學分數", "必選別", "科目序號", "科目名稱", "授課教師", "上課時間/地點", "備註", "數位課綱URL", "影片標題", "影片長度(second)", "影片長度(string)", "<3min||>5min"})

	if len(svlreq.svXi) <= 0 {
		err = errors.WithStack(fmt.Errorf("svXi has no data"))
		return err
	}
	for _, xi := range svlreq.svXi {
		for _, sv := range xi {
			tempXi := make([]string, 0)
			tempXi = append(tempXi, svlreq.academicYear, svlreq.semester, sv.grade, sv.class, sv.credit, sv.chooseSelect, sv.courseID, sv.courseName, sv.teacherName, sv.timeNPlace, sv.remark, sv.videoURL, sv.videoTitle, strconv.Itoa(sv.videoSeconds), sv.videoDuration)
			if sv.videoSeconds > 300 || sv.videoSeconds < 180 {
				tempXi = append(tempXi, "OutOfRange")
			}
			dsvFile.newDataRows = append(dsvFile.newDataRows, tempXi)
		}
	}
	return nil
}

//appendVideoInfo 將影片資料放入newDataRows中
func (dsvFile *downloadSVFile) appendVideoInfo(svlreq *getSVLRequest) (err error) {
	if len(dsvFile.dataRows) <= 0 {
		err = errors.WithStack(fmt.Errorf("dataRows has no data"))
		return err
	}
	for index, xi := range dsvFile.dataRows {
		if index == 0 {
			dsvFile.dataRows[index] = append(dsvFile.dataRows[index], "數位課綱URL", "影片標題", "影片長度(second)", "影片長度(string)", "<3min||>5min")
			continue
		}
		if len(svlreq.svXi[xi[dsvFile.cidCol]]) > 0 {
			dsvFile.dataRows[index] = append(dsvFile.dataRows[index], svlreq.svXi[xi[dsvFile.cidCol]][0].videoURL, svlreq.svXi[xi[dsvFile.cidCol]][0].videoTitle, strconv.Itoa(svlreq.svXi[xi[dsvFile.cidCol]][0].videoSeconds), svlreq.svXi[xi[dsvFile.cidCol]][0].videoDuration)
			if svlreq.svXi[xi[dsvFile.cidCol]][0].videoSeconds > 300 || svlreq.svXi[xi[dsvFile.cidCol]][0].videoSeconds < 180 {
				dsvFile.dataRows[index] = append(dsvFile.dataRows[index], "OutOfRange")
			}
		}
	}
	dsvFile.newDataRows = dsvFile.dataRows
	return nil
}
