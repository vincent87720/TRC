//go:generate goversioninfo -manifest=../../tools/goversioninfo/goversioninfo.exe.manifest

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/fatih/color"
)

var (
	INITPATH string
)

type college struct {
	collegeID   string //學院代號
	collegeName string //學院名稱
}

type department struct {
	college
	departmentID   string //學系序號
	departmentName string //學系名稱
}

type teacher struct {
	department
	teacherID         string //教師編號
	teacherName       string //教師姓名
	teacherPhone      string //分機
	teacherMail       string //電子郵件
	teacherSpace      string //研究室編號
	teacherState      string //任職狀態
	teacherLevel      string //職稱
	teacherLastUpdate string //最後更新日期
}

type course struct {
	department
	courseID      string   //課程編號
	courseName    string   //課程名稱
	courseSubName string   //課程副標題
	year          int      //學年度
	semester      int      //學期
	system        string   //學制
	group         string   //組別
	grade         string   //年級
	class         string   //班級
	credit        string   //學分數
	chooseSelect  string   //必選別
	interval      string   //上課時數
	time          []string //上課時間
	classRoom     []string //教室
	numOfPeople   string   //選課人數
	annex         string   //合班註記
	annexID       string   //合班序號
	remark        string   //備註
}

type student struct {
	department
	studentID   string //學生學號
	studentName string //學生姓名
}

type syllabusVideo struct {
	course
	teacher
	timeNPlace      string
	videoID         string
	videoTitle      string
	videoURL        string //數位課綱連結
	videoDuration   string //數位課綱影片長度(string)
	videoSeconds    int    //數位課綱影片長度(second)
	problemOfCourse string //影片問題
}

type scoreAllert struct {
	course
	student
	allertReason string //預警原因
	tutorMethod  string //輔導方式
}

type rapidPrint struct {
	course
	timeNClassRoom string
}

type facing struct {
	facingName string    //面向名稱
	score      []float64 //各評委分數
	difference []float64 //差分
}

//file 檔案
type file struct {
	filePath    string     //檔案路徑
	fileName    string     //檔案名稱
	sheetName   string     //工作表名稱
	firstRow    []string   //第一行
	dataRows    [][]string //資料
	newDataRows [][]string //新資料
	xlsx        *excelize.File
}

//svFile 數位課綱檔案
type svFile struct {
	file
	cdpCol       int                         //開課單位編號欄位(departmentID)
	cidCol       int                         //課程編號欄位(courseID)
	csnCol       int                         //課程名稱欄位(courseName)
	pocCol       int                         //影片問題欄位(problemOfCourse)
	cthCol       int                         //任課老師欄位(courseTeacher)
	gbtd         map[teacher][]syllabusVideo //group by teacher data
	mergedXi     [][]string                  //合併後的陣列
	maxCourseNum int                         //所有老師中，擁有最多科目的數量
}

//saFile 成績預警檔案
type saFile struct {
	file
	csdCol int                       //開課單位(courseDepartment)
	cidCol int                       //科目序號(courseID)
	csnCol int                       //課程名稱(courseName)
	cthCol int                       //任課老師(courseTeacher)
	sidCol int                       //學生學號(studentID)
	stnCol int                       //學生姓名(studentName)
	alrCol int                       //預警原由(alertReason)
	gbtd   map[teacher][]scoreAllert //group by teacher data
}

//thFile 教師資料檔案
type thFile struct {
	file
	didCol int //學院編號
	tidCol int //教師編號
	trnCol int //教師姓名
	tdpCol int //教師系所

	teacherMap map[string][]teacher //教師資料
}

//rpFile 快速印刷檔案
type rpFile struct {
	file
	trnCol int                      //教師姓名欄位
	tstCol int                      //專兼任別欄位
	sysCol int                      //開課學制欄位
	csdCol int                      //開課系所欄位
	cidCol int                      //科目序號欄位
	ifoCol int                      //系-組-年-班欄位
	csnCol int                      //科目名稱欄位
	ccsCol int                      //選修別欄位
	cdtCol int                      //學分欄位
	ctmCol int                      //時數欄位
	wtrCol int                      //星期-時間-教室欄位
	nopCol int                      //選課人數欄位
	annCol int                      //合班註記欄位
	aidCol int                      //合班序號欄位
	rmkCol int                      //備註欄位
	gbtd   map[teacher][]rapidPrint //group by teacher data
}

//dcFile 差分計算檔案
type dcFile struct {
	file
	judgeNum int                  //評審數量
	judge    []string             //評審
	gbsd     map[student][]facing //group by student data
}

//downloadFile 下載的檔案
type downloadTeacherFile struct {
	file
}

//downloadSyllabusVideoFile 下載的檔案
type downloadSVFile struct {
	file
	cidCol int //科目序號欄位
}

type command struct {
	commandString string
	commandAction string
	flagSet       *flag.FlagSet
}

type commandSet struct {
	lyr1commandSet map[string][]command //第一層指令集
	lyr2commandSet map[string][]command //第二層指令集
}

type flags struct {
	help bool

	//general parameter
	inputFilePath   string //-inPath
	inputFileName   string //-inName
	inputSheetName  string //-inSheet
	outputFilePath  string //-outPath
	outputFileName  string //-outName
	outputSheetName string //-outSheet

	//download video parameter
	appendVideoInfo      bool
	academicYear         string //-year
	semester             string //-semester
	youtubeAPIKey        string
	downloadSVInputFile  []string
	downloadSVOutputFile []string

	//merge syllabus video parameter
	tfile         bool
	svInputFile   []string
	svOutputFile  []string
	svTeacherFile []string //-teacher

	//merge course data parameter
	cdInputFile  []string
	cdOutputFile []string

	//split scoreAlert parameter
	scoreAlertFile     []string //-masterlist
	teacherInfoFile    []string //-teacher
	exportTemplateFile []string //-exptemplate

	//calculate difference parameter
	readAllFilesInDir bool

	startGuiFlagSet        *flag.FlagSet
	downloadVideoFlagSet   *flag.FlagSet
	splitScoreAlertFlagSet *flag.FlagSet
	mergeVideoFlagSet      *flag.FlagSet
	mergeCourseFlagSet     *flag.FlagSet
	calcDifferentFlagSet   *flag.FlagSet
}

//選擇使用模式
func selectMode(trcCmd *commandSet, f *flags) (err error) {
	if len(os.Args) > 1 {
		autoMode(trcCmd, f)
	} else {
		err = manualMode(trcCmd, f)
		if err != nil {
			return err
		}
	}
	return nil
}

func autoMode(trcCmd *commandSet, f *flags) {
	flag.Parse()
	analyseCommand(trcCmd, f)
}

func manualMode(trcCmd *commandSet, f *flags) (err error) {

	usr, err := user.Current()
	if err != nil {
		return err
	}

	path, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		color.Set(color.FgHiRed)
		fmt.Print("\n", usr.Username)
		color.Set(color.FgHiCyan)
		fmt.Print(" TRC ")
		color.Set(color.FgWhite)
		fmt.Print(path, "\r\n")
		fmt.Print(">")

		reader := bufio.NewReader(os.Stdin)
		data, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		dataXi := strings.Fields(data)
		os.Args = os.Args[:0]
		for _, val := range dataXi {
			os.Args = append(os.Args, val)
		}
		err = f.initFlag()
		if err != nil {
			return err
		}
		flag.Parse()
		analyseCommand(trcCmd, f)
	}
}

func setInitialPath() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Error.Fatalf("%+v\n", err)
	}
	INITPATH = filepath.ToSlash(dir)
}

func init() {
	setInitialPath()
	exportAssets()
	initLogging()
}

func main() {
	var f flags
	err := f.initFlag()
	if err != nil {
		Error.Printf("%+v\n", err)
	}

	var trcCmd commandSet
	trcCmd.commandSet()
	trcCmd.setLyr1Command("start", "啟動服務(啟動圖形化介面)")
	trcCmd.setLyr2Command("start", "gui", "啟動圖形化介面", f.startGuiFlagSet)
	trcCmd.setLyr1Command("download", "下載檔案(下載教師資料、數位課綱連結)")
	trcCmd.setLyr2Command("download", "video", "下載數位課綱影片連結", f.downloadVideoFlagSet)
	trcCmd.setLyr2Command("download", "teacher", "下載教師資料", nil)
	trcCmd.setLyr1Command("split", "分割檔案(分割預警總表)")
	trcCmd.setLyr2Command("split", "scoreAlert", "分割預警總表", f.splitScoreAlertFlagSet)
	trcCmd.setLyr1Command("merge", "合併檔案(合併數位課綱資料、開課及製版數登記表)")
	trcCmd.setLyr2Command("merge", "video", "依教師合併數位課綱影片問題", f.mergeVideoFlagSet)
	trcCmd.setLyr2Command("merge", "course", "依教師合併開課總表，製作為製版數登記表", f.mergeCourseFlagSet)
	trcCmd.setLyr1Command("calculate", "計算檔案(計算成績差分)")
	trcCmd.setLyr2Command("calculate", "difference", "計算成績差分", f.calcDifferentFlagSet)

	err = selectMode(&trcCmd, &f)
	if err != nil {
		Error.Printf("%+v\n", err)
	}
}

func findSpecificExtentionFiles(path string, extention string) (files []string, err error) {
	fileXi := make([]string, 0)

	allFiles, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, fi := range allFiles {
		match, _ := regexp.MatchString(extention+"$", fi.Name())
		if match {
			fileXi = append(fileXi, fi.Name())
		}
	}

	return fileXi, nil
}
