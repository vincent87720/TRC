package RapidPrint

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/vincent87720/TRC/internal/cmdline"
	"github.com/vincent87720/TRC/internal/file"
	logging "github.com/vincent87720/TRC/internal/logging"
	"github.com/vincent87720/TRC/internal/object"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

//rpFile 快速印刷檔案
type rpFile struct {
	file.WorksheetFile
	trnCol int                             //教師姓名欄位
	tstCol int                             //專兼任別欄位
	sysCol int                             //開課學制欄位
	csdCol int                             //開課系所欄位
	cidCol int                             //科目序號欄位
	ifoCol int                             //系-組-年-班欄位
	csnCol int                             //科目名稱欄位
	ccsCol int                             //選修別欄位
	cdtCol int                             //學分欄位
	ctmCol int                             //時數欄位
	wtrCol int                             //星期-時間-教室欄位
	nopCol int                             //選課人數欄位
	annCol int                             //合班註記欄位
	aidCol int                             //合班序號欄位
	rmkCol []int                           //備註欄位
	gbtd   map[object.Teacher][]rapidPrint //group by teacher data
}

type rapidPrint struct {
	object.Course
	timeNClassRoom string
}

type rapidPrintXi []rapidPrint

func (rp rapidPrintXi) Len() int           { return len(rp) }
func (rp rapidPrintXi) Less(i, j int) bool { return rp[i].CourseName < rp[j].CourseName }
func (rp rapidPrintXi) Swap(i, j int)      { rp[i], rp[j] = rp[j], rp[i] }
func (rp rapidPrintXi) Sort()              { sort.Sort(rp) }

type rapidPrintXi_2 []rapidPrint

func (rp rapidPrintXi_2) Len() int { return len(rp) }
func (rp rapidPrintXi_2) Less(i, j int) bool {
	iid, _ := strconv.Atoi(rp[i].CourseID)
	jid, _ := strconv.Atoi(rp[j].CourseID)

	return iid < jid
}
func (rp rapidPrintXi_2) Swap(i, j int) { rp[i], rp[j] = rp[j], rp[i] }
func (rp rapidPrintXi_2) Sort()         { sort.Sort(rp) }

type teacherXi []object.Teacher

func (t teacherXi) Len() int { return len(t) }
func (t teacherXi) Less(i, j int) bool {
	utf8ToBig5 := traditionalchinese.Big5.NewEncoder()
	ibig5, _, _ := transform.String(utf8ToBig5, t[i].TeacherName)
	jbig5, _, _ := transform.String(utf8ToBig5, t[j].TeacherName)

	return ibig5 < jbig5
}
func (t teacherXi) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t teacherXi) Sort()         { sort.Sort(t) }

// MergeRapidPrintData 剖析並合併開課資料
// Goroutine interface for GUI
// For example:
//
//	var inputFile file
// 	var outputFile file
//
// 	inputFile.setFile("Your file path", "Your file name", "Your sheet name")
// 	outputFile.setFile("Your file path", "Your file name", "Your sheet name")
//
// 	errChan := make(chan error, 2)
// 	exitChan := make(chan string, 2)
// 	defer close(errChan)
// 	defer close(exitChan)
//
// 	go MergeRapidPrintData(errChan, exitChan, inputFile, outputFile)
//
// Loop:
//	for {
// 		select {
// 		case err := <-errChan:
// 			fmt.Println(err)
// 		case <-exitChan:
// 			break Loop
// 		}
// 	}
func MergeRapidPrintData(progChan chan int, inputFile file.WorksheetFile, outputFile file.WorksheetFile) {

	rpf := rpFile{
		WorksheetFile: inputFile,
	}

	err := rpf.ReadRawData()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = rpf.FillSliceLength(len(rpf.FirstRow))
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = rpf.findColumn()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = rpf.groupByTeacher()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = rpf.mergeData()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = rpf.transportToSlice()
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1
	err = rpf.ExportDataToExcel(outputFile)
	if err != nil {
		logging.Error.Printf("%+v\n", err)
		return
	}
	progChan <- 1

	return
}

//mergeCourseData 合併開課總表內可合併的課程
func MergeRapidPrintData_Command(inputFileInfo file.FileInfo, outputFileInfo file.FileInfo) (err error) {

	var inputFile rpFile
	var outputFile file.WorksheetFile
	inputFile.SetFile(inputFileInfo.FilePath, inputFileInfo.FileName, inputFileInfo.SheetName)
	outputFile.SetFile(outputFile.FilePath, outputFile.FileName, outputFile.SheetName)

	quit := make(chan int)
	defer close(quit)

	go cmdline.Spinner("Course data file is loading...", 80*time.Millisecond, quit)
	err = inputFile.ReadRawData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Loading file have completed")

	go cmdline.Spinner("Data is preprocessing...", 80*time.Millisecond, quit)
	err = inputFile.FillSliceLength(15)
	if err != nil {
		quit <- 1
		return err
	}
	err = inputFile.findColumn()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Preprocessing have completed")

	go cmdline.Spinner("Data is grouping by teacher...", 80*time.Millisecond, quit)
	err = inputFile.groupByTeacher()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Grouping have completed")

	go cmdline.Spinner("Data is merging...", 80*time.Millisecond, quit)
	err = inputFile.mergeData()
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Merging have completed")

	go cmdline.Spinner("Files are exporting...", 80*time.Millisecond, quit)
	err = inputFile.transportToSlice()
	if err != nil {
		quit <- 1
		return err
	}
	err = inputFile.ExportDataToExcel(outputFile)
	if err != nil {
		quit <- 1
		return err
	}
	quit <- 1
	fmt.Println("\r> Exporting have completed")
	return nil
}

//findColumn 尋找欄位
func (rpf *rpFile) findColumn() (err error) {

	err = rpf.FindCol("教師姓名", &rpf.trnCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("專兼任別", &rpf.tstCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("開課學制", &rpf.sysCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("開課系所", &rpf.csdCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("科目序號", &rpf.cidCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("系-組-年-班", &rpf.ifoCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("科目名稱", &rpf.csnCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("選修別", &rpf.ccsCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("學分", &rpf.cdtCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("時數", &rpf.ctmCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("星期-時間-教室", &rpf.wtrCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("選課人數", &rpf.nopCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("合班註記", &rpf.annCol)
	if err != nil {
		return err
	}
	err = rpf.FindCol("合班序號", &rpf.aidCol)
	if err != nil {
		return err
	}
	err = rpf.FindAllCol("教師姓名", &rpf.rmkCol)
	if err != nil {
		return err
	}
	return nil
}

//groupByTeacher 依照教師名稱將開課資料分群
func (rpf *rpFile) groupByTeacher() (err error) {

	r1, err := regexp.Compile(`\(.*?\s`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	r2, err := regexp.Compile(`\(.*?\-`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}
	r3, err := regexp.Compile(`[A-Z].*`)
	if err != nil {
		err = errors.WithStack(err)
		return err
	}

	rpf.gbtd = make(map[object.Teacher][]rapidPrint)
	if len(rpf.DataRows) <= 0 {
		err = errors.WithStack(fmt.Errorf("dataRows has no data"))
		return err
	}
	for index, value := range rpf.DataRows {

		//跳過第零行標題列
		if index == 0 {
			continue
		}

		//使用正規表達式拆解"星期-時間-教室"欄位
		thisTime := make([]string, 0)
		thisClassRoom := make([]string, 0)

		substrs := r1.FindAllString(value[rpf.wtrCol], -1)
		for _, loc := range substrs {
			loc = strings.Replace(loc, " ", "", -1)
			splitWtrLoc := r2.FindStringIndex(loc)
			thisTime = append(thisTime, loc[splitWtrLoc[0]:splitWtrLoc[1]-1])
			thisClassRoom = append(thisClassRoom, loc[splitWtrLoc[1]:])
		}

		//使用正規表達式拆解"科目名稱"欄位，若課程名稱包含副標題(ex.G3、J1...)則將其分開，並將課程名稱放入thisCourseName，副標題放入thisCourseSubName
		var thisCourseName string
		var thisCourseSubName string
		splitCsnLoc := r3.FindStringIndex(value[rpf.csnCol])
		if len(splitCsnLoc) > 0 {
			thisCourseName = value[rpf.csnCol][:splitCsnLoc[0]]
			thisCourseSubName = value[rpf.csnCol][splitCsnLoc[0]:]
		} else {
			thisCourseName = value[rpf.csnCol]
		}

		//合併除了第一個教師姓名以外的所有教師姓名，並放入remark欄位
		combineTeachers := ""
		checked := false
		for idx, val := range rpf.rmkCol {
			if idx == 0 {
				continue
			}
			if value[val] != "" {
				if checked == false {
					combineTeachers = combineTeachers + value[val]
					checked = true
				} else {
					combineTeachers = combineTeachers + "、" + value[val]
				}
			}
		}
		// fmt.Println(combineTeachers)

		if value[rpf.trnCol] != "" {
			t := object.Teacher{
				TeacherName:  value[rpf.trnCol],
				TeacherState: value[rpf.tstCol],
			}
			c := rapidPrint{
				Course: object.Course{
					System: value[rpf.sysCol],
					Department: object.Department{
						DepartmentName: value[rpf.csdCol],
						DepartmentID:   value[rpf.ifoCol][0:4],
					},
					CourseID:      value[rpf.cidCol],
					CourseName:    thisCourseName,
					CourseSubName: thisCourseSubName,
					Group:         value[rpf.ifoCol][5:8],
					Grade:         value[rpf.ifoCol][9:10],
					Class:         value[rpf.ifoCol][11:12],
					ChooseSelect:  value[rpf.ccsCol],
					Credit:        value[rpf.cdtCol],
					Interval:      value[rpf.ctmCol],
					Time:          thisTime,
					ClassRoom:     thisClassRoom,
					NumOfPeople:   value[rpf.nopCol],
					Annex:         value[rpf.annCol],
					AnnexID:       value[rpf.aidCol],
					Remark:        combineTeachers,
				},
			}
			rpf.gbtd[t] = append(rpf.gbtd[t], c)
		}
	}
	return nil
}

//mergeData 依照規則合併重複的科目名稱
func (rpf *rpFile) mergeData() (err error) {
	//對一位老師的所有科目進行處理
	mergedRpData := make(map[object.Teacher][]rapidPrint)
	if len(rpf.gbtd) <= 0 {
		err = errors.WithStack(fmt.Errorf("gbtd has no data"))
		return err
	}
	for key, value := range rpf.gbtd {

		//排序
		var rpp rapidPrintXi
		rpp = value
		rpp.Sort()

		var nextLoop int
		nextLoop = 0

		//比對目前科目(rpp[i])與下一個科目(rpp[index])是否相同
		for index, _ := range rpp {
			if index < nextLoop {
				continue
			}

			tempSystemMap := make(map[string]int)     //暫存相同課程名稱的學制
			tempRemarkMap := make(map[string]int)     //暫存相同課程名稱的備註
			tempSubNameMap := make(map[string]int)    //暫存相同課程名稱的副標題
			tempDepartmentMap := make(map[string]int) //暫存相同課程名稱的開課系所
			tempCourseIDXi := make([]string, 0)       //暫存相同課程名稱的課程編號
			tempCourseTimeXi := make([]string, 0)     //暫存相同課程名稱的上課時間地點
			tempNumOfPeopleXi := make([]string, 0)    //暫存相同課程名稱的修課人數

			sameCourseNumCount := 0 //紀錄目前科目以外和目前科目可以合併的科目數量

			//將目前科目的資訊放入暫存變數中
			tempSystemMap[rpp[index].System] = 1
			tempRemarkMap[rpp[index].Remark] = 1
			tempSubNameMap[rpp[index].CourseSubName] = 1
			tempDepartmentMap[rpp[index].Course.Department.DepartmentName] = 1
			tempNumOfPeopleXi = append(tempNumOfPeopleXi, rpp[index].NumOfPeople)
			for i := range rpp[index].Time {
				tempCourseTimeXi = append(tempCourseTimeXi, rpp[index].Time[i]+" "+rpp[index].ClassRoom[i])
			}

			//此為下一個科目(rpp[index])
			for nextIndex := index + 1; nextIndex < len(rpp); nextIndex++ {

				//若符合名稱相同及學制限制，則進行合併
				if rpp[index].CourseName == rpp[nextIndex].CourseName &&
					(((rpp[index].System == "大學日間部" || rpp[index].System == "進修學士班" || rpp[index].System == "四技部") && (rpp[nextIndex].System == "大學日間部" || rpp[nextIndex].System == "進修學士班" || rpp[nextIndex].System == "四技部")) ||
						((rpp[index].System == "大學日間部" || rpp[index].System == "研究所碩士班") && (rpp[nextIndex].System == "大學日間部" || rpp[nextIndex].System == "研究所碩士班")) ||
						((rpp[index].System == "研究所碩士班" || rpp[index].System == "碩士在職專班") && (rpp[nextIndex].System == "研究所碩士班" || rpp[nextIndex].System == "碩士在職專班"))) {

					//設定學制的map為1
					tempSystemMap[rpp[nextIndex].System] = 1

					//設定副標題的map為1，若已存在則累加
					v1, found := tempSubNameMap[rpp[nextIndex].CourseSubName]
					if found {
						tempSubNameMap[rpp[nextIndex].CourseSubName] = v1 + 1
					} else {
						tempSubNameMap[rpp[nextIndex].CourseSubName] = 1
					}

					//設定備註的map為1
					tempRemarkMap[rpp[nextIndex].Remark] = 1

					//設定學系的map為1
					tempDepartmentMap[rpp[nextIndex].Course.Department.DepartmentName] = 1

					//暫存人數欄位到tempNumOfPeopleXi
					tempNumOfPeopleXi = append(tempNumOfPeopleXi, rpp[nextIndex].NumOfPeople)

					//暫存科目序號到tempCourseIDXi
					tempCourseIDXi = append(tempCourseIDXi, rpp[nextIndex].CourseID)

					//暫存時間及教室欄位到tempCourseTimeXi
					for i := range rpp[nextIndex].Time {
						tempCourseTimeXi = append(tempCourseTimeXi, rpp[nextIndex].Time[i]+" "+rpp[nextIndex].ClassRoom[i])
					}

					sameCourseNumCount++

				} else {
					break
				}
			}

			//串接學制map到tempSystemMap
			if len(tempSystemMap) > 0 {
				var newSystem string
				var count bool
				count = false
				for key := range tempSystemMap {
					if count == false {
						newSystem = newSystem + key
						count = true
					} else {
						newSystem = newSystem + "," + key
					}

				}
				rpp[index].System = newSystem
			}

			//串接系所
			if len(tempDepartmentMap) > 0 {
				var newDep string
				var count bool
				count = false
				for key := range tempDepartmentMap {
					if count == false {
						newDep = newDep + key
						count = true
					} else {
						newDep = newDep + "," + key
					}

				}
				rpp[index].Course.Department.DepartmentName = newDep
			}

			//在科目名稱欄位加上科目名稱數量及串接副標題
			if len(tempSubNameMap) > 0 {
				var newSubName string
				count := false
				for key, value := range tempSubNameMap {
					if count == false {
						newSubName = newSubName + key + "*" + strconv.Itoa(value)
						count = true
					} else {
						newSubName = newSubName + "," + key + "*" + strconv.Itoa(value)
					}
				}
				rpp[index].CourseName = rpp[index].CourseName + newSubName
			}

			//串接備註
			if len(tempRemarkMap) > 0 {
				var newRemark string
				count := false
				for key, _ := range tempRemarkMap {
					if key != "" {
						if count == false {
							newRemark = newRemark + key
							count = true
						} else {
							newRemark = newRemark + "、" + key
						}
					}
				}
				rpp[index].Remark = newRemark
			}

			//串接科目序號
			if len(tempCourseIDXi) > 0 {
				var newCourseID string

				for index, idXi := range tempCourseIDXi {
					if index == 0 {
						newCourseID = newCourseID + idXi
					} else {
						newCourseID = newCourseID + "," + idXi
					}
				}

				if len(rpp[index].Remark) == 0 {
					newCourseID = "合" + newCourseID

				} else {
					newCourseID = "合" + newCourseID + "、" + rpp[index].Remark
				}
				rpp[index].Remark = newCourseID
			}

			//串接時間
			if len(tempCourseTimeXi) > 0 {
				var newTime string
				for index, timeXi := range tempCourseTimeXi {
					if index == 0 {
						newTime = newTime + timeXi
					} else {
						newTime = newTime + "," + timeXi
					}
				}
				rpp[index].timeNClassRoom = newTime
			}

			//串接人數
			if len(tempNumOfPeopleXi) > 0 {
				var newNOP string
				for index, nopXi := range tempNumOfPeopleXi {
					if index == 0 {
						newNOP = newNOP + nopXi
					} else {
						newNOP = newNOP + "+" + nopXi
					}
				}
				rpp[index].NumOfPeople = newNOP
			}

			//設定下次要執行的圈數
			nextLoop = index + sameCourseNumCount + 1

			//將tempMerged加入到mergedRpData
			mergedRpData[key] = append(mergedRpData[key], rpp[index])
		}

		//依照課程編號排序
		var rp2 rapidPrintXi_2
		rp2 = mergedRpData[key]
		rp2.Sort()
		mergedRpData[key] = rp2
	}

	//將合併過後的檔案加入
	rpf.gbtd = mergedRpData
	return nil
}

func (rpf *rpFile) transportToSlice() (err error) {
	//設定第一列
	rpf.NewDataRows = append(rpf.NewDataRows, []string{"教師姓名", "專兼任別", "開課學制", "開課系所", "科目序號", "系-組-年-班", "科目名稱", "選修別", "學分", "時數", "星期-時間-教室", "選課人數", "合班註記", "合班序號", "備註"})
	if len(rpf.gbtd) <= 0 {
		err = errors.WithStack(fmt.Errorf("gbtd has no data"))
		return err
	}

	//依照教師姓名排序
	keys := make([]object.Teacher, 0, len(rpf.gbtd))
	for t := range rpf.gbtd {
		keys = append(keys, t)
	}
	var tXi teacherXi
	tXi = keys
	tXi.Sort()

	for _, t := range tXi {
		for _, rpValue := range rpf.gbtd[t] {
			tempXi := make([]string, 0)
			tempXi = append(tempXi, t.TeacherName, t.TeacherState, rpValue.System, rpValue.Course.DepartmentName, rpValue.CourseID, rpValue.Course.DepartmentID+"-"+rpValue.Group+"-"+rpValue.Grade+"-"+rpValue.Class, rpValue.CourseName, rpValue.ChooseSelect, rpValue.Credit, rpValue.Interval, rpValue.timeNClassRoom, rpValue.NumOfPeople, rpValue.Annex, rpValue.AnnexID, rpValue.Remark)
			rpf.NewDataRows = append(rpf.NewDataRows, tempXi)
		}
	}
	return nil
}
