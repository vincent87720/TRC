package main

import (
	"log"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func RunMergeVideoDialog(owner walk.Form, fi *MergeVideoFileInfo, iconFilePath string) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	var db *walk.DataBinder

	var inputFilePath *walk.LineEdit
	var teacherFilePath *walk.LineEdit

	var inputSheetSelector *walk.ComboBox
	var teacherSheetSelector *walk.ComboBox

	labelFont := Font{Family: "Microsoft JhengHei", PointSize: 11}

	return Dialog{
		AssignTo:   &dlg,
		Title:      "MergeSyllabusVideoData",
		Icon:       iconFilePath,
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Font:       Font{Family: "Microsoft JhengHei", PointSize: 9},
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "mergeVideoInfo",
			DataSource:     fi,
			ErrorPresenter: ToolTipErrorPresenter{},
		},
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{300, 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 4},
				Children: []Widget{
					Label{
						Text: "數位課綱",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &inputFilePath,
						Text:     Bind("InputPath"),
						MinSize:  Size{250, 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							OnOpenFileButtonClicked(owner, inputFilePath, inputSheetSelector)
						},
					},
					ComboBox{
						AssignTo:      &inputSheetSelector,
						Editable:      false,
						BindingMember: "Name",
						DisplayMember: "Name",
						Value:         Bind("InputSheet"),
					},

					Label{
						Text:    "教師名單",
						Font:    labelFont,
						Visible: Bind("tfile.Checked"),
					},
					LineEdit{
						AssignTo: &teacherFilePath,
						Text:     Bind("TeacherPath"),
						MinSize:  Size{250, 0},
						ReadOnly: true,
						Visible:  Bind("tfile.Checked"),
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							OnOpenFileButtonClicked(owner, teacherFilePath, teacherSheetSelector)
						},
						Visible: Bind("tfile.Checked"),
					},
					ComboBox{
						AssignTo:      &teacherSheetSelector,
						Editable:      false,
						BindingMember: "Name",
						DisplayMember: "Name",
						Value:         Bind("TeacherSheet"),
						Visible:       Bind("tfile.Checked"),
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					Label{
						Text: "使用教師名單檔案匯入所屬單位",
					},
					CheckBox{
						Name:    "tfile",
						Checked: Bind("TFile"),
					},
					HSpacer{},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								log.Print(err)
								return
							}
							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
