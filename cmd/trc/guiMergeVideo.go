package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func runMergeVideoDialog(owner walk.Form, fi *MergeVideoFileInfo, iconFilePath string) (int, error) {
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
		Title:      "合併數位課綱資料",
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
		MinSize:       Size{Width: 300, Height: 300},
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
						MinSize:  Size{Width: 250, Height: 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							onOpenFileButtonClicked(owner, inputFilePath, inputSheetSelector)
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
						MinSize:  Size{Width: 250, Height: 0},
						ReadOnly: true,
						Visible:  Bind("tfile.Checked"),
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							onOpenFileButtonClicked(owner, teacherFilePath, teacherSheetSelector)
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
						Text:     "確定",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								Error.Printf("%+v\n", err)
								return
							}
							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "取消",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)
}
