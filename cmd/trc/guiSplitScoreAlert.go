package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func runSplitScoreAlertDialog(owner walk.Form, fi *SplitScoreAlertFileInfo, iconFilePath string) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton

	var db *walk.DataBinder

	var masterFilePath *walk.LineEdit
	var teacherFilePath *walk.LineEdit
	var templateFilePath *walk.LineEdit

	var masterSheetSelector *walk.ComboBox
	var teacherSheetSelector *walk.ComboBox
	var templateSheetSelector *walk.ComboBox

	labelFont := Font{Family: "Microsoft JhengHei", PointSize: 11}

	return Dialog{
		AssignTo:   &dlg,
		Title:      "分割預警總表",
		Icon:       iconFilePath,
		Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
		Font:       Font{Family: "Microsoft JhengHei", PointSize: 9},
		DataBinder: DataBinder{
			AssignTo:       &db,
			Name:           "splitFileInfo",
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
						Text: "預警總表",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &masterFilePath,
						Text:     Bind("MasterPath"),
						MinSize:  Size{Width: 250, Height: 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							onOpenFileButtonClicked(owner, masterFilePath, masterSheetSelector)
						},
					},
					ComboBox{
						AssignTo: &masterSheetSelector,
						Name:     "masterCombobox",
						Editable: false,
						//BindingMember:用來設定綁定的欄位
						//以DropDownItem struct為例
						//type DropDownItem struct {
						//	Key  int
						//	Name string
						//}
						//若想要設定綁定數值，必須使用BindingMember: "Name"以確保選取得數值會正確綁定
						//
						//DisplayMember用來顯示下拉式選單的內容，同樣必須使用DisplayMember: "Name"設定以確保畫面上可看到所有下拉式選單的選項
						BindingMember: "Name",
						DisplayMember: "Name",
						Value:         Bind("MasterSheet"),
					},

					Label{
						Text: "教師名單",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &teacherFilePath,
						Text:     Bind("TeacherPath"),
						MinSize:  Size{Width: 250, Height: 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							onOpenFileButtonClicked(owner, teacherFilePath, teacherSheetSelector)
						},
					},
					ComboBox{
						AssignTo:      &teacherSheetSelector,
						Editable:      false,
						BindingMember: "Name",
						DisplayMember: "Name",
						Value:         Bind("TeacherSheet"),
					},

					Label{
						Text: "空白分表",
						Font: labelFont,
					},
					LineEdit{
						AssignTo: &templateFilePath,
						Text:     Bind("TemplatePath"),
						MinSize:  Size{Width: 250, Height: 0},
						ReadOnly: true,
					},
					PushButton{
						Text: "選擇檔案",
						OnClicked: func() {
							onOpenFileButtonClicked(owner, templateFilePath, templateSheetSelector)
						},
					},
					ComboBox{
						AssignTo:      &templateSheetSelector,
						Editable:      false,
						BindingMember: "Name",
						DisplayMember: "Name",
						Value:         Bind("TemplateSheet"),
					},
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
