package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type DropDownItem struct { // Used in the ComboBox dropdown
	Key  int
	Name string
}

type NormalFileInfo struct {
	InputPath  string
	InputSheet string
}

type DownloadVideoInfo struct {
	InputPath  string
	InputSheet string
	Year       string
	Semester   string
	Key        string
	Append     bool
}

type SplitScoreAlertFileInfo struct {
	MasterPath    string
	MasterSheet   string
	TeacherPath   string
	TeacherSheet  string
	TemplatePath  string
	TemplateSheet string
}

type MergeVideoFileInfo struct {
	InputPath    string
	InputSheet   string
	TeacherPath  string
	TeacherSheet string
	TFile        bool
}

type CalculateDifferenceInfo struct {
	InputPath  string
	InputSheet string
	CalcAll    bool
}

func RunMainWindow() {

	iconMain := filepath.FromSlash("assets/guiImage/blockchain-blueblue.png")
	iconDownload := filepath.FromSlash("assets/guiImage/Those_Icons-download-32.png")
	iconSplit := filepath.FromSlash("assets/guiImage/Those_Icons-split-32.png")
	iconCalculate := filepath.FromSlash("assets/guiImage/Pixel_Perfect-calculate-32.png")
	iconMerge := filepath.FromSlash("assets/guiImage/Those_Icons-merge-32.png")

	font := Font{Family: "Microsoft JhengHei", PointSize: 12}
	titleFont := Font{Family: "Microsoft JhengHei", PointSize: 18}

	var mw *walk.MainWindow

	r, err := regexp.Compile(`(.*\\)([^\\]*.xlsx)`)
	if err != nil {
		log.Print(err)
	}

	if _, err := (MainWindow{
		AssignTo:   &mw,
		Title:      "TRC",
		Icon:       iconMain,
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Size:       Size{1000, 200},
		MinSize:    Size{1000, 200},
		Layout:     VBox{},
		Children: []Widget{
			VSpacer{
				MaxSize: Size{1, 20},
			},
			Composite{
				Layout: Grid{Columns: 4},
				Children: []Widget{
					Label{
						Text: "TRC教學資源中心",
						Font: titleFont,
					},
				},
			},
			VSpacer{
				MaxSize: Size{1, 20},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Composite{
						Layout: VBox{},
						Children: []Widget{
							//download
							PushButton{
								Text:    "Download Teacher",
								Image:   iconDownload,
								Font:    font,
								MinSize: Size{200, 50},
								OnClicked: func() {
									if cmd, err := RunDownloadTeacherDialog(mw, iconDownload); err != nil {
										log.Print(err)
									} else if cmd == walk.DlgCmdOK {
										checkOutputDir()

									}
								},
							},

							//download
							PushButton{
								Text:    "Download Video   ",
								Image:   iconDownload,
								Font:    font,
								MinSize: Size{200, 50},
								OnClicked: func() {
									fi := new(DownloadVideoInfo)
									if cmd, err := RunDownloadVideoDialog(mw, fi, iconDownload); err != nil {
										log.Print(err)
									} else if cmd == walk.DlgCmdOK {
										checkOutputDir()

									}
								},
							},
						},
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							//split
							PushButton{
								Text:    "Split ScoreAlert",
								Image:   iconSplit,
								Font:    font,
								MinSize: Size{200, 50},
								OnClicked: func() {
									fi := new(SplitScoreAlertFileInfo)
									if cmd, err := RunSplitScoreAlertDialog(mw, fi, iconSplit); err != nil {
										log.Print(err)
									} else if cmd == walk.DlgCmdOK {
										checkOutputDir()
										var masterFile file
										var templateFile file
										var teacherFile file
										var outputFile file

										masterPathXi := r.FindStringSubmatch(fi.MasterPath)
										teacherPathXi := r.FindStringSubmatch(fi.TeacherPath)
										templatePathXi := r.FindStringSubmatch(fi.TemplatePath)

										outputFile.setFile(filepath.FromSlash(INITPATH+"/output/"), "", "")
										masterFile.setFile(masterPathXi[1], masterPathXi[2], fi.MasterSheet)
										templateFile.setFile(templatePathXi[1], templatePathXi[2], fi.TemplateSheet)
										teacherFile.setFile(teacherPathXi[1], teacherPathXi[2], fi.TeacherSheet)

										errChan := make(chan error, 2)
										exitChan := make(chan string, 2)
										defer close(errChan)
										defer close(exitChan)

										go SplitScoreAlertData(errChan, exitChan, masterFile, templateFile, teacherFile, outputFile)

									Loop:
										for {
											select {
											case err := <-errChan:
												Error.Printf("%+v\n", err)
											case <-exitChan:
												break Loop
											}
										}
									}
								},
							},
						},
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							//calculate
							PushButton{
								Text:    "Calculate Difference",
								Image:   iconCalculate,
								Font:    font,
								MinSize: Size{200, 50},
								OnClicked: func() {
									fi := new(CalculateDifferenceInfo)
									if cmd, err := RunCalculateDifferenceDialog(mw, fi, iconCalculate); err != nil {
										log.Print(err)
									} else if cmd == walk.DlgCmdOK {
										checkOutputDir()
									}
								},
							},
						},
					},
					Composite{
						Layout: VBox{},
						Children: []Widget{
							//merge
							PushButton{
								Text:    "Merge Course",
								Image:   iconMerge,
								Font:    font,
								MinSize: Size{200, 50},
								OnClicked: func() {
									fi := new(NormalFileInfo)
									if cmd, err := RunMergeCourseDialog(mw, fi, iconMerge); err != nil {
										log.Print(err)
									} else if cmd == walk.DlgCmdOK {
										checkOutputDir()

										var inputFile file
										var outputFile file

										inputPathXi := r.FindStringSubmatch(fi.InputPath)
										inputFile.setFile(inputPathXi[1], inputPathXi[2], fi.InputSheet)
										outputFile.setFile(filepath.FromSlash(INITPATH+"/output/"), "[MERGENCE]開課總表.xlsx", "工作表")

										errChan := make(chan error, 2)
										exitChan := make(chan string, 2)
										defer close(errChan)
										defer close(exitChan)

										go MergeRapidPrintData(errChan, exitChan, inputFile, outputFile)

									Loop:
										for {
											select {
											case err := <-errChan:
												Error.Printf("%+v\n", err)
											case <-exitChan:
												break Loop
											}
										}
									}
								},
							},

							//merge
							PushButton{
								Text:    "Merge Video  ",
								Image:   iconMerge,
								Font:    font,
								MinSize: Size{200, 50},
								OnClicked: func() {
									fi := new(MergeVideoFileInfo)
									if cmd, err := RunMergeVideoDialog(mw, fi, iconMerge); err != nil {
										log.Print(err)
									} else if cmd == walk.DlgCmdOK {
										checkOutputDir()
										var inputFile file
										var outputFile file
										inputPathXi := r.FindStringSubmatch(fi.InputPath)
										inputFile.setFile(inputPathXi[1], inputPathXi[2], fi.InputSheet)
										outputFile.setFile(filepath.FromSlash(INITPATH+"/output/"), "[MERGENCE]數位課綱.xlsx", "工作表")

										errChan := make(chan error, 2)
										exitChan := make(chan string, 2)
										defer close(errChan)
										defer close(exitChan)

										if fi.TFile {
											var teacherFile file
											teacherPathXi := r.FindStringSubmatch(fi.TeacherPath)
											teacherFile.setFile(teacherPathXi[1], teacherPathXi[2], fi.TeacherSheet)

											go MergeSyllabusVideoDataByList(errChan, exitChan, inputFile, outputFile, teacherFile)

										} else {
											go MergeSyllabusVideoData(errChan, exitChan, inputFile, outputFile)
										}

									Loop:
										for {
											select {
											case err := <-errChan:
												Error.Printf("%+v\n", err)
											case <-exitChan:
												break Loop
											}
										}
									}
								},
							},
						},
					},
				},
			},
			HSpacer{},
		},
	}.Run()); err != nil {
		log.Fatal(err)
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

func OnOpenFileButtonClicked(owner walk.Form, filePath *walk.LineEdit, selector *walk.ComboBox) {
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
