package TeacherData

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vincent87720/TRC/internal/file"
	"github.com/vincent87720/TRC/internal/object"
)

//thFile 教師資料檔案
type ThFile struct {
	file.WorksheetFile
	DidCol int //學院編號
	TidCol int //教師編號
	TrnCol int //教師姓名
	TdpCol int //教師系所

	TeacherMap map[string][]object.Teacher //教師資料
}

//groupByTeacher 依照教師名稱將教師資料分群
func (thr *ThFile) GroupByTeacher() (err error) {
	err = thr.FindCol("學院編號", &thr.DidCol)
	if err != nil {
		return err
	}
	thr.FindCol("教師編號", &thr.TidCol)
	thr.FindCol("教師姓名", &thr.TrnCol)
	thr.FindCol("所屬單位名稱", &thr.TdpCol)

	thr.TeacherMap = make(map[string][]object.Teacher)
	if len(thr.DataRows[0]) <= 0 {
		err = errors.WithStack(fmt.Errorf("object.Teacher list dataRows has no data"))
		return err
	}
	for index, value := range thr.DataRows {

		//跳過第零行標題列
		if index == 0 {
			continue
		}

		if value[thr.TidCol] != "" {
			t := object.Teacher{
				Department: object.Department{
					College: object.College{
						CollegeID: value[thr.DidCol],
					},
					DepartmentName: value[thr.TdpCol],
				},
				TeacherID:   value[thr.TidCol],
				TeacherName: value[thr.TrnCol],
			}
			thr.TeacherMap[t.TeacherName] = append(thr.TeacherMap[t.TeacherName], t)
		}
	}
	return nil
}
