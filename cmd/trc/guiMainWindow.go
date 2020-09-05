// +build windows

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

// DropDownItem 下拉式選單結構
type DropDownItem struct { // Used in the ComboBox dropdown
	Key  int
	Name string
}

// NormalFileInfo structure for DataBinder
type NormalFileInfo struct {
	InputPath  string
	InputSheet string
}

// DownloadVideoInfo structure for DataBinder
type DownloadVideoInfo struct {
	InputPath  string
	InputSheet string
	Year       string
	Semester   string
	Key        string
	Append     bool
}

// SplitScoreAlertFileInfo structure for DataBinder
type SplitScoreAlertFileInfo struct {
	MasterPath    string
	MasterSheet   string
	TeacherPath   string
	TeacherSheet  string
	TemplatePath  string
	TemplateSheet string
}

// MergeVideoFileInfo structure for DataBinder
type MergeVideoFileInfo struct {
	InputPath    string
	InputSheet   string
	TeacherPath  string
	TeacherSheet string
	TFile        bool
}

// CalculateDifferenceInfo structure for DataBinder
type CalculateDifferenceInfo struct {
	InputPath  string
	InputSheet string
	Folder     string
	CalcAll    bool
}

func runMainWindow() {

	iconMain := filepath.FromSlash("assets/guiImage/blockchain-blueblue.png")
	iconDownload := filepath.FromSlash("assets/guiImage/Those_Icons-download-32.png")
	iconSplit := filepath.FromSlash("assets/guiImage/Those_Icons-split-32.png")
	iconCalculate := filepath.FromSlash("assets/guiImage/Pixel_Perfect-calculate-32.png")
	iconMerge := filepath.FromSlash("assets/guiImage/Those_Icons-merge-32.png")

	font := Font{Family: "Microsoft JhengHei", PointSize: 12}
	subTitleFont := Font{Family: "Microsoft JhengHei", PointSize: 15}
	titleFont := Font{Family: "Microsoft JhengHei", PointSize: 18}

	var mw *walk.MainWindow
	var pb *walk.ProgressBar

	r, err := regexp.Compile(`(.*\\)([^\\]*.xlsx)`)
	if err != nil {
		Error.Printf("%+v\n", err)
	}

	if _, err := (MainWindow{
		AssignTo:   &mw,
		Title:      "TRC",
		Icon:       iconMain,
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Size:       Size{Width: 1000, Height: 200},
		MinSize:    Size{Width: 1000, Height: 200},
		Layout:     VBox{},
		Children: []Widget{
			VSpacer{
				MaxSize: Size{Width: 1, Height: 20},
			},
			Label{
				Text:          "教學資源中心",
				Font:          titleFont,
				TextAlignment: AlignCenter,
			},
			VSpacer{
				MaxSize: Size{Width: 1, Height: 20},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Composite{
						Layout: VBox{},
						Children: []Widget{
							//download
							Label{
								Text:          "下載",
								Font:          subTitleFont,
								TextAlignment: AlignCenter,
							},
							PushButton{
								Text:    "　教師資料　　",
								Image:   iconDownload,
								Font:    font,
								MinSize: Size{Width: 200, Height: 50},
								OnClicked: func() {
									checkOutputDir()

									progChan := make(chan int, 2)

									var outputFile file

									outputFile.setFile(filepath.ToSlash(INITPATH+"/output/"), "教師名單.xlsx", "工作表")

									go setProgressBar(progChan, mw, pb, 0, 8)
									go GetTeacher(progChan, outputFile)
								},
							},

							//download
							PushButton{
								Text:    "　數位課綱資料",
								Image:   iconDownload,
								Font:    font,
								MinSize: Size{Width: 200, Height: 50},
								OnClicked: func() {
									fi := new(DownloadVideoInfo)
									if cmd, err := runDownloadVideoDialog(mw, fi, iconDownload); err != nil {
										Error.Printf("%+v\n", err)
									} else if cmd == walk.DlgCmdOK {
										checkOutputDir()

										progChan := make(chan int, 2)

										var outputFile file

										outputFile.setFile(filepath.ToSlash(INITPATH+"/output/"), fi.Year+fi.Semester+"數位課綱.xlsx", "工作表")

										if fi.Append {
											var inputFile file
											inputPathXi := r.FindStringSubmatch(fi.InputPath)

											inputFile.setFile(inputPathXi[1], inputPathXi[2], fi.InputSheet)

											go setProgressBar(progChan, mw, pb, 0, 10)
											go AppendSyllabusVideo(progChan, fi.Year, fi.Semester, fi.Key, inputFile, outputFile)

										} else {
											go setProgressBar(progChan, mw, pb, 0, 7)
											go GetSyllabusVideo(progChan, fi.Year, fi.Semester, fi.Key, outputFile)
										}
									}
								},
							},
							VSpacer{},
						},
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							//calculate
							Label{
								Text:          "計算",
								Font:          subTitleFont,
								TextAlignment: AlignCenter,
							},
							PushButton{
								Text:    "　成績差分",
								Image:   iconCalculate,
								Font:    font,
								MinSize: Size{Width: 200, Height: 50},
								OnClicked: func() {
									fi := new(CalculateDifferenceInfo)
									if cmd, err := runCalculateDifferenceDialog(mw, fi, iconCalculate); err != nil {
										Error.Printf("%+v\n", err)
									} else if cmd == walk.DlgCmdOK {
										checkOutputDir()
										//Calculate all files in fi.Folder
										if fi.CalcAll {

											files, err := findSpecificExtentionFiles(fi.Folder, ".xlsx")

											if err != nil {
												Error.Printf("%+v\n", err)
											}

											progChan := make(chan int, 5)
											go setProgressBar(progChan, mw, pb, 0, len(files)*7)

											for _, value := range files {

												var inputFile file
												var outputFile file

												inputFile.setFile(filepath.ToSlash(fi.Folder+"/"), value, "學系彙整版")
												outputFile.setFile(filepath.ToSlash(INITPATH+"/output/"), "[DIFFERENCE]"+value, "工作表")

												go DifferenceCalculate(progChan, inputFile, outputFile)
											}

										} else {
											var inputFile file
											var outputFile file

											inputPathXi := r.FindStringSubmatch(fi.InputPath)
											inputFile.setFile(inputPathXi[1], inputPathXi[2], fi.InputSheet)
											outputFile.setFile(filepath.FromSlash(INITPATH+"/output/"), "[DIFFERENCE]"+inputPathXi[2], "工作表")

											progChan := make(chan int, 2)

											go setProgressBar(progChan, mw, pb, 0, 7)
											go DifferenceCalculate(progChan, inputFile, outputFile)

										}
									}
								},
							},
							VSpacer{},
						},
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							//split
							Label{
								Text:          "分割",
								Font:          subTitleFont,
								TextAlignment: AlignCenter,
							},
							PushButton{
								Text:    "　預警總表",
								Image:   iconSplit,
								Font:    font,
								MinSize: Size{Width: 200, Height: 50},
								OnClicked: func() {
									fi := new(SplitScoreAlertFileInfo)
									if cmd, err := runSplitScoreAlertDialog(mw, fi, iconSplit); err != nil {
										Error.Printf("%+v\n", err)
									} else if cmd == walk.DlgCmdOK {
										checkOutputDir()
										var masterFile file
										var templateFile file
										var teacherFile file
										var outputFile file

										masterPathXi := r.FindStringSubmatch(fi.MasterPath)
										teacherPathXi := r.FindStringSubmatch(fi.TeacherPath)
										templatePathXi := r.FindStringSubmatch(fi.TemplatePath)

										outputFile.setFile(filepath.ToSlash(INITPATH+"/output/"), "", "")
										masterFile.setFile(masterPathXi[1], masterPathXi[2], fi.MasterSheet)
										templateFile.setFile(templatePathXi[1], templatePathXi[2], fi.TemplateSheet)
										teacherFile.setFile(teacherPathXi[1], teacherPathXi[2], fi.TeacherSheet)

										progChan := make(chan int, 2)

										go setProgressBar(progChan, mw, pb, 0, 6)
										go SplitScoreAlertData(progChan, masterFile, templateFile, teacherFile, outputFile)
									}
								},
							},
							VSpacer{},
						},
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							//merge
							Label{
								Text:          "合併",
								Font:          subTitleFont,
								TextAlignment: AlignCenter,
							},
							PushButton{
								Text:    "　製版數登記表",
								Image:   iconMerge,
								Font:    font,
								MinSize: Size{Width: 200, Height: 50},
								OnClicked: func() {
									fi := new(NormalFileInfo)
									if cmd, err := runMergeCourseDialog(mw, fi, iconMerge); err != nil {
										Error.Printf("%+v\n", err)
									} else if cmd == walk.DlgCmdOK {
										checkOutputDir()

										var inputFile file
										var outputFile file

										inputPathXi := r.FindStringSubmatch(fi.InputPath)
										inputFile.setFile(inputPathXi[1], inputPathXi[2], fi.InputSheet)
										outputFile.setFile(filepath.FromSlash(INITPATH+"/output/"), "[MERGENCE]開課總表.xlsx", "工作表")

										progChan := make(chan int, 2)

										go setProgressBar(progChan, mw, pb, 0, 7)
										go MergeRapidPrintData(progChan, inputFile, outputFile)
									}
								},
							},

							//merge
							PushButton{
								Text:    "　數位課綱資料",
								Image:   iconMerge,
								Font:    font,
								MinSize: Size{Width: 200, Height: 50},
								OnClicked: func() {
									fi := new(MergeVideoFileInfo)
									if cmd, err := runMergeVideoDialog(mw, fi, iconMerge); err != nil {
										Error.Printf("%+v\n", err)
									} else if cmd == walk.DlgCmdOK {
										checkOutputDir()
										var inputFile file
										var outputFile file
										inputPathXi := r.FindStringSubmatch(fi.InputPath)
										inputFile.setFile(inputPathXi[1], inputPathXi[2], fi.InputSheet)
										outputFile.setFile(filepath.FromSlash(INITPATH+"/output/"), "[MERGENCE]數位課綱.xlsx", "工作表")

										progChan := make(chan int, 2)

										if fi.TFile {
											var teacherFile file
											teacherPathXi := r.FindStringSubmatch(fi.TeacherPath)
											teacherFile.setFile(teacherPathXi[1], teacherPathXi[2], fi.TeacherSheet)

											go setProgressBar(progChan, mw, pb, 0, 7)
											go MergeSyllabusVideoDataByList(progChan, inputFile, outputFile, teacherFile)
										} else {
											go setProgressBar(progChan, mw, pb, 0, 5)
											go MergeSyllabusVideoData(progChan, inputFile, outputFile)
										}
									}
								},
							},
							VSpacer{},
						},
					},
				},
			},
			HSpacer{},
			ProgressBar{
				AssignTo: &pb,
				MaxValue: 100,
				MinValue: 0,
				Visible:  false,
			},
		},
	}.Run()); err != nil {
		Error.Fatalf("%+v\n", err)
	}
}

func exportAssets() {
	if _, err := os.Stat(INITPATH + "/assets/guiImage"); os.IsNotExist(err) {
		err = os.MkdirAll("assets/guiImage", os.ModeDir)
		if err != nil {
			Error.Printf("%+v\n", err)
		}
	}

	hasMainIcon := false
	hasCalcIcon := false
	hasDownloadIcon := false
	hasMergeIcon := false
	hasSplitIcon := false

	allFiles, err := ioutil.ReadDir("assets/guiImage")
	if err != nil {
		Error.Printf("%+v\n", err)
	}

	for _, fi := range allFiles {
		if fi.Name() == "blockchain-blueblue.png" {
			hasMainIcon = true
		} else if fi.Name() == "Pixel_Perfect-calculate-32.png" {
			hasCalcIcon = true
		} else if fi.Name() == "Those_Icons-download-32.png" {
			hasDownloadIcon = true
		} else if fi.Name() == "Those_Icons-merge-32.png" {
			hasMergeIcon = true
		} else if fi.Name() == "Those_Icons-split-32.png" {
			hasSplitIcon = true
		}
	}

	if !hasMainIcon || !hasCalcIcon || !hasDownloadIcon || !hasMergeIcon || !hasSplitIcon {
		fileXi, err := AssetDir("")
		if err != nil {
			Error.Printf("%+v\n", err)
		}

		for _, value := range fileXi {
			xi, err := Asset(value)
			if err != nil {
				Error.Printf("%+v\n", err)
			}
			// convert []byte to image for saving to file
			img, _, _ := image.Decode(bytes.NewReader(xi))

			//save the imgByte to file
			out, err := os.Create("./assets/" + value)

			if err != nil {
				Error.Printf("%+v\n", err)
				os.Exit(1)
			}

			err = png.Encode(out, img)

			if err != nil {
				Error.Printf("%+v\n", err)
				os.Exit(1)
			}
		}
	}
}

func checkOutputDir() {
	err := os.MkdirAll(INITPATH+"\\output", os.ModeDir)
	if err != nil {
		Error.Printf("%+v\n", err)
	}
}

func onOpenFileButtonClicked(owner walk.Form, filePath *walk.LineEdit, selector *walk.ComboBox) {
	dlg := new(walk.FileDialog)
	dlg.Title = "Open File"
	dlg.Filter = "Excel檔案 (*.xlsx)|*.xlsx|所有檔案 (*.*)|*.*"

	if ok, err := dlg.ShowOpen(owner); err != nil {
		fmt.Fprintln(os.Stderr, "Error:Open master file")
		return
	} else if !ok {
		fmt.Fprintln(os.Stderr, "Cancel file selection")
		return
	}

	filePath.SetText(dlg.FilePath)

	keys := []*DropDownItem{}

	xlsx, err := excelize.OpenFile(dlg.FilePath)
	if err != nil {
		Error.Printf("%+v\n", err)
	}
	xlsxSht := xlsx.GetSheetMap()
	for idx, val := range xlsxSht {
		keys = append(keys, &DropDownItem{Key: idx, Name: val})
	}
	selector.SetModel(keys)
}

func onBrowseFolderButtonClicked(owner walk.Form, filePath *walk.LineEdit) {
	dlg := new(walk.FileDialog)
	dlg.Title = "Browse Folder"

	if ok, err := dlg.ShowBrowseFolder(owner); err != nil {
		fmt.Fprintln(os.Stderr, "Error:Open folder")
		return
	} else if !ok {
		fmt.Fprintln(os.Stderr, "Cancel folder browsing")
		return
	}

	filePath.SetText(dlg.FilePath)
}

func setProgressBar(progChan chan int, mw *walk.MainWindow, pb *walk.ProgressBar, min int, max int) {
	currentVal := 0
	mw.Synchronize(func() {
		pb.SetVisible(true)
		pb.SetRange(min, max)
	})
Loop:
	for {
		select {
		case val := <-progChan:
			currentVal += val
			mw.Synchronize(func() {
				pb.SetValue(currentVal)
			})
			if currentVal >= max {
				time.Sleep(time.Duration(1) * time.Second)
				mw.Synchronize(func() {
					pb.SetVisible(false)
					pb.SetValue(0)
				})
				close(progChan)
				break Loop
			}
		}
	}
}
