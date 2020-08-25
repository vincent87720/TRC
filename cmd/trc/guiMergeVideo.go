package main

import (
	"fmt"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func RunMergeVideoDialog(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	var masterFilePath *walk.LineEdit
	var teacherFilePath *walk.LineEdit
	var templateFilePath *walk.LineEdit

	var masterSheetSelector *walk.ComboBox
	var teacherSheetSelector *walk.ComboBox
	var templateSheetSelector *walk.ComboBox

	labelFont := Font{Family: "Microsoft JhengHei", PointSize: 11}

	return Dialog{
		AssignTo: &dlg,
		Title:    "Split",
		Icon:          "./assets/guiImage/Those_Icons-split-32.png",
		Background:    SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Font:          Font{Family: "Microsoft JhengHei", PointSize: 9},
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize:       Size{300, 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 4},
				Children: []Widget{
					Label{
						Text: "預警總表",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &masterFilePath,
						MinSize:  Size{250, 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							OnOpenFileButtonClicked(owner, masterFilePath, masterSheetSelector)
						},
					},
					ComboBox{
						AssignTo:      &masterSheetSelector,
						Editable:      false,
						BindingMember: "key",
						DisplayMember: "Name",
					},

					Label{
						Text: "教師名單",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &teacherFilePath,
						MinSize:  Size{250, 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							OnOpenFileButtonClicked(owner, teacherFilePath, teacherSheetSelector)
						},
					},
					ComboBox{
						AssignTo:      &teacherSheetSelector,
						Editable:      false,
						BindingMember: "key",
						DisplayMember: "Name",
					},

					Label{
						Text: "空白分表",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &templateFilePath,
						MinSize:  Size{250, 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							OnOpenFileButtonClicked(owner, templateFilePath, templateSheetSelector)
						},
					},
					ComboBox{
						AssignTo:      &templateSheetSelector,
						Editable:      false,
						BindingMember: "key",
						DisplayMember: "Name",
					},
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
							fmt.Println(masterFilePath.Text())
							fmt.Println(teacherFilePath.Text())
							fmt.Println(templateFilePath.Text())

							fmt.Println(masterSheetSelector.Text())
							fmt.Println(teacherSheetSelector.Text())
							fmt.Println(templateSheetSelector.Text())

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
