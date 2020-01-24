package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	help bool

	//download video parameter
	academicYear   string //-year
	semester       string //-semester
	outputFileName string //-filename

	//split scoreAlert parameter
	scoreAlertFilePath     string //-masterlist
	teacherInfoFilePath    string //-teacher
	exportTemplateFilePath string //-exptemplate
)

func getSyllabusVideo() {

	quit := make(chan struct{})

	sr := newCourseList()
	co := &course{}
	sr.setAcademicYear(academicYear)
	sr.setSemester(semester)

	go spinner("Downloading", 80*time.Millisecond, quit)
	sr.getCourseData()
	close(quit)
	fmt.Println("\r> Download completed")

	quit = make(chan struct{})
	go spinner("Parsing", 80*time.Millisecond, quit)
	co.find(sr.htmlNode)
	close(quit)
	fmt.Println("\r> Parsing completed")

	quit = make(chan struct{})
	go spinner("Exporting", 80*time.Millisecond, quit)
	co.exportToExcel(outputFileName)
	close(quit)
	fmt.Println("\r> Export completed")

}

func getTeacher() {
	teacher := &teacher{}

	engineering := &searchTeacherRequest{
		dept:       "'2000','5040','5011','5022','5080','5100','5240','6013','5023','5081','5082','6312'",
		deptItem:   "1",
		searchItem: "4",
	}
	management := &searchTeacherRequest{
		dept:       "'2003','8021','7003','6410','5110','5120','5150','5131','5140','5190','7004'",
		deptItem:   "2",
		searchItem: "4",
	}
	foreignLanguages := &searchTeacherRequest{
		dept:       "'2005','5212','5220','5230','5211','5231','6861'",
		deptItem:   "3",
		searchItem: "4",
	}
	desighAndArts := &searchTeacherRequest{
		dept:       "'2004','5070','5030','5060','5090','6001','5091','5096','6600'",
		deptItem:   "4",
		searchItem: "4",
	}
	biotechnology := &searchTeacherRequest{
		dept:       "'2006','5052','5512','5180','5250'",
		deptItem:   "5",
		searchItem: "4",
	}
	tourism := &searchTeacherRequest{
		dept:       "'2007','5260','5270','7180','5161','5163','5659'",
		deptItem:   "6",
		searchItem: "4",
	}
	studentAffairsAndPhysical := &searchTeacherRequest{
		dept:       "'3200','3300'",
		deptItem:   "7",
		searchItem: "4",
	}
	center := &searchTeacherRequest{
		dept:       "'9010','4210','4020','4007','4025'",
		deptItem:   "8",
		searchItem: "4",
	}
	nursing := &searchTeacherRequest{
		dept:       "'2008','7173','5290','5280','5172','5242','6800'",
		deptItem:   "9",
		searchItem: "4",
	}

	quit := make(chan struct{})
	go spinner("Downloading", 80*time.Millisecond, quit)
	engineering.getTeacherData()
	fmt.Println("\r> 工 學 院OK ")
	management.getTeacherData()
	fmt.Println("\r> 管理學院OK ")
	foreignLanguages.getTeacherData()
	fmt.Println("\r> 外語學院OK ")
	desighAndArts.getTeacherData()
	fmt.Println("\r> 設藝學院OK ")
	biotechnology.getTeacherData()
	fmt.Println("\r> 生資學院OK ")
	tourism.getTeacherData()
	fmt.Println("\r> 觀光學院OK ")
	studentAffairsAndPhysical.getTeacherData()
	fmt.Println("\r> 學程單位OK ")
	center.getTeacherData()
	fmt.Println("\r> 學程中心OK ")
	nursing.getTeacherData()
	close(quit)
	fmt.Println("\r> Download completed")

	quit = make(chan struct{})
	go spinner("Parsing", 80*time.Millisecond, quit)
	teacher.parseTeacherData(engineering.deptItem, engineering.sitemap)
	teacher.parseTeacherData(management.deptItem, management.sitemap)
	teacher.parseTeacherData(foreignLanguages.deptItem, foreignLanguages.sitemap)
	teacher.parseTeacherData(desighAndArts.deptItem, desighAndArts.sitemap)
	teacher.parseTeacherData(biotechnology.deptItem, biotechnology.sitemap)
	teacher.parseTeacherData(tourism.deptItem, tourism.sitemap)
	teacher.parseTeacherData(studentAffairsAndPhysical.deptItem, studentAffairsAndPhysical.sitemap)
	teacher.parseTeacherData(center.deptItem, center.sitemap)
	teacher.parseTeacherData(nursing.deptItem, nursing.sitemap)
	close(quit)
	fmt.Println("\r> Parsing completed")

	quit = make(chan struct{})
	go spinner("Exporting", 80*time.Millisecond, quit)
	teacher.ExportToExcel("教師名單")
	close(quit)
	fmt.Println("\r> Export completed")
}

func splitScoreAlertFile() {
	quit := make(chan struct{})
	sa := &scoreAlert{}

	if teacherInfoFilePath == "" {
		fmt.Println("\rERROR:未設定\"教師名單\"檔案路徑名稱")
		os.Exit(2)
	}
	if scoreAlertFilePath == "" {
		fmt.Println("\rERROR:未設定\"預警總表\"檔案路徑名稱")
		os.Exit(2)
	}
	if exportTemplateFilePath == "" {
		fmt.Println("\rERROR:未設定\"空白分表\"檔案路徑名稱")
		os.Exit(2)
	}

	go spinner("Loading teacher list", 80*time.Millisecond, quit)
	err := sa.loadTeacherInfo()
	if err != nil {
		panic(err)
	}
	close(quit)
	fmt.Println("\r> Loaded teacher list")

	quit = make(chan struct{})
	go spinner("Loading score alert list", 80*time.Millisecond, quit)
	err = sa.loadScoreAlertList()
	if err != nil {
		panic(err)
	}
	close(quit)
	fmt.Println("\r> Loaded score alert list ")

	quit = make(chan struct{})
	go spinner("Splitting", 80*time.Millisecond, quit)
	sa.splitScoreAlertData()
	close(quit)
	fmt.Println("\r> Splitting completed")
}

func spinner(status string, delay time.Duration, ch chan struct{}) {
	for {
		select {
		case <-ch:
			fmt.Printf("\r                              ")
			return
		default:
			for _, r := range `-\|/` {
				fmt.Printf("\r%c %s", r, status)
				time.Sleep(delay)
			}
		}
	}
}

func init() {
	flag.BoolVar(&help, "help", false, "Usage")
	flag.Usage = func() {
		fmt.Println("Usage: trc <command> [<args>]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  download [<args>]\n        下載檔案(教師資料、數位課綱連結)")
		fmt.Println("  split [<args>]\n        分割檔案(分割預警總表)")
	}
}

func setDownloadVideoFlag(downloadVideoCommand *flag.FlagSet) {
	//set default video parameter
	defaultAcademicYear := ""
	defaultSemester := ""
	if time.Now().Month() <= 7 {
		defaultAcademicYear = strconv.Itoa(time.Now().Year() - 1911 - 1)
		defaultSemester = "2"
	} else {
		defaultAcademicYear = strconv.Itoa(time.Now().Year() - 1911)
		defaultSemester = "1"
	}
	downloadVideoCommand.StringVar(&academicYear, "year", defaultAcademicYear, "設定學年度，預設為當前學年度")
	downloadVideoCommand.StringVar(&semester, "semester", defaultSemester, "設定學期，預設為當前學期")
	downloadVideoCommand.StringVar(&outputFileName, "filename", "數位課綱", "設定檔名，預設為查詢目標學年+學期+數位課綱")
}

func setSplitScoreAlertFlag(splitScoreAlertCommand *flag.FlagSet) {
	splitScoreAlertCommand.StringVar(&scoreAlertFilePath, "masterlist", "", "設定預警總表檔案路徑名稱")
	splitScoreAlertCommand.StringVar(&teacherInfoFilePath, "teacher", "", "設定教師名單檔案路徑名稱")
	splitScoreAlertCommand.StringVar(&exportTemplateFilePath, "exptemplate", "", "設定空白分表檔案名稱")
}

func main() {
	flag.Parse()
	downloadVideoCommand := flag.NewFlagSet("download video", flag.ExitOnError)
	setDownloadVideoFlag(downloadVideoCommand)
	splitScoreAlertCommand := flag.NewFlagSet("split scoreAlert", flag.ExitOnError)
	setSplitScoreAlertFlag(splitScoreAlertCommand)

	if len(os.Args) == 1 || help {
		flag.Usage()
		return
	}

	switch os.Args[1] {
	case "download":

		downloadUsage := func() {
			fmt.Println("Usage: trc download <command> [<args>]")
			fmt.Println()
			fmt.Println("Commands:")
			// fmt.Println("  video [<args>]\n        下載數位課綱影片連結")
			fmt.Println("  video [<args>]\t下載數位課綱影片連結")
			for key, value := range downloadVideoCommand.GetFlags() {
				fmt.Println("        -" + key + " " + value.Usage)
			}
			// fmt.Println("  teacher [<args>]\n        下載教師資料")
			fmt.Println("  teacher [<args>]\t下載教師資料")
		}

		if len(os.Args) <= 2 {
			downloadUsage()
			return
		}

		if len(os.Args) >= 2 {
			switch os.Args[2] {
			case "video":
				err := downloadVideoCommand.Parse(os.Args[3:])
				if err != nil {
					log.Println(err)
				}

				fmt.Println(">> GetSyllabusVideoLink")
				getSyllabusVideo()
			case "teacher":
				fmt.Println(">> GetTeacherData")
				getTeacher()
			default:
				downloadUsage()
				return
			}
		}

	case "split":

		splitUsage := func() {
			fmt.Println("Usage: trc split <command> [<args>]")
			fmt.Println()
			fmt.Println("Commands:")
			// fmt.Println("  scoreAlert [<args>]\n        分割預警總表")
			fmt.Println("  scoreAlert [<args>]\t分割預警總表")
			for key, value := range splitScoreAlertCommand.GetFlags() {
				fmt.Println("        -" + key + " " + value.Usage)
			}
		}

		if len(os.Args) <= 2 {
			splitUsage()
			return
		}

		if len(os.Args) >= 2 {
			switch os.Args[2] {
			case "scoreAlert":
				err := splitScoreAlertCommand.Parse(os.Args[3:])
				if err != nil {
					log.Println(err)
				}

				fmt.Println(">> SplitScoreAlertFile")
				splitScoreAlertFile()
			default:
				splitUsage()
				return
			}
		}
	default:
		fmt.Printf("%q is not a valid command.\n", os.Args[1])
		os.Exit(2)
	}

}
