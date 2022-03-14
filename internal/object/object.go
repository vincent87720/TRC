package object

type College struct {
	CollegeID   string //學院代號
	CollegeName string //學院名稱
}

type Department struct {
	College
	DepartmentID   string //學系序號
	DepartmentName string //學系名稱
}

type Teacher struct {
	Department
	TeacherID         string //教師編號
	TeacherName       string //教師姓名
	TeacherPhone      string //分機
	TeacherMail       string //電子郵件
	TeacherSpace      string //研究室編號
	TeacherState      string //任職狀態
	TeacherLevel      string //職稱
	TeacherLastUpdate string //最後更新日期
}

type Course struct {
	Department
	CourseID      string   //課程編號
	CourseName    string   //課程名稱
	CourseSubName string   //課程副標題
	Year          int      //學年度
	Semester      int      //學期
	System        string   //學制
	Group         string   //組別
	Grade         string   //年級
	Class         string   //班級
	Credit        string   //學分數
	ChooseSelect  string   //必選別
	Interval      string   //上課時數
	Time          []string //上課時間
	ClassRoom     []string //教室
	NumOfPeople   string   //選課人數
	Annex         string   //合班註記
	AnnexID       string   //合班序號
	Remark        string   //備註
}

type Student struct {
	Department
	StudentID   string //學生學號
	StudentName string //學生姓名
}
