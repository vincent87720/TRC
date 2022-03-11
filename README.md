# TRC Program
此程式提供給教資中心使用，程式包含教師資料下載、預警資料分割、製版數統計表合併、數位課綱資料合併以及尺規評量差分計算。


# 開始
可觀看[教學影片](https://youtu.be/hKECNRl-Zqs)觀看如何下載及啟動使用者介面  


# 下載
進入[下載頁面](https://github.com/vincent87720/TRC/releases)下載最新版本程式  


# 圖形化介面
## 啟動
- 輸入指令`trc start gui`可啟動圖形化使用者介面  
![圖形化使用者介面](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/gui.png)  

## 下載教師資料
### 規則
無  
### 介面
點選`下載>教師資料`，下載完成的教師資料會儲存在output資料夾中  
![download_teacherData](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/download_teacherData.png)  
### 輸入檔案範例
無  
### 輸出檔案範例
![trcDownloadTeacher](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/trcDownloadTeacher.png)  

## 下載數位課綱連結
### 規則
 - 若有指定參數會依照參數設定值運行
 - 若無指定參數或依照預設設定值運行
 - 檔案格式必須為.xlsx
### 介面
無  
### 輸入檔案範例
無  
### 輸出檔案範例
![trcDownloadVideo](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/trcDownloadVideo.png)  

## 分割預警總表
### 規則
- 預警總表須包含**開課學系**、**科目序號**、**預警科目**、**授課教師**、**學號**、**學生姓名**及**預警原由**欄位
- 教師名單須包含**學院編號**、**教師編號**、**教師姓名**及**所屬單位名稱**欄位
- 檔案格式必須為.xlsx
### 介面
![split_scoreAlert](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/split_scoreAlert.png)  
點選`分割>預警總表`後會跳出檔案選擇視窗  
選擇預警總表、教師名單及空白分表並於後方指定工作表  
點選確定後程式會將分割後的表單匯出到output資料夾中  
![split_scoreAlert_selectfile](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/split_scoreAlert_selectfile.png)  
### 輸入檔案範例
[分割預警總表範例檔案](https://github.com/vincent87720/TRC/blob/master/assets/exampleInputFile/%E5%88%86%E5%89%B2%E9%A0%90%E8%AD%A6%E7%B8%BD%E8%A1%A8)  
### 輸出檔案範例
無  

## 計算尺規成績差分
### 規則
- 面向、學生及評委可任意增加或刪減，其餘須按照特定格式
- 檔案格式必須為.xlsx
### 介面
![calculate_differenceCalculation](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/calculate_differenceCalculation.png)  
點選`計算>成績差分`後會跳出檔案選擇視窗  
可選擇處理單一檔案或多個檔案  
處理單一檔案需指定檔案及工作表  
![calculate_difference_singlefile](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/calculate_difference_singlefile.png)  
處理多個檔案需指定檔案所在資料夾  
![calculate_difference_multiplefile](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/calculate_difference_multiplefile.png)  
點選確定後程式會將計算完成的的表單匯出到output資料夾中  
### 輸入檔案範例
[計算尺規成績差分範例檔案](https://github.com/vincent87720/TRC/blob/master/assets/exampleInputFile/%E8%A8%88%E7%AE%97%E5%B0%BA%E8%A6%8F%E6%88%90%E7%B8%BE%E5%B7%AE%E5%88%86)  
### 輸出檔案範例
![trcCalculateDifference](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/trcCalculateDifference.png)  

## 合併製版數登記表
### 規則
- 預警總表須包含**教師姓名、專兼任別、開課學制、開課系所、科目序號、系-組-年-班、科目名稱、選修別、學分、時數、星期-時間-教室、選課人數、合班註記、合班序號**及**備註**欄位
- 可同時有多個欄位的名稱為"教師姓名"(第一個以外的教師姓名會放入備註)
- 檔案格式必須為.xlsx
- 開課學制需使用以下名稱命名：大學日間部、進修學士班、四技部、研究所碩士班、碩士在職專班、研究所博士班、二年制在職專班
### 介面
![merge_rapidPrint](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/merge_rapidPrint.png)  
點選`合併>製版數登記表`後會跳出檔案選擇視窗  
選擇開課總表檔案並於後方指定工作表  
點選確定後程式會合併檔案並匯出到output資料夾中  
![merge_rapidPrint_selectfile](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/merge_rapidPrint_selectfile.png)  
### 輸入檔案範例
[合併製版數登記表範例檔案](https://github.com/vincent87720/TRC/blob/master/assets/exampleInputFile/%E5%90%88%E4%BD%B5%E8%A3%BD%E7%89%88%E6%95%B8%E7%99%BB%E8%A8%98%E8%A1%A8)  
### 輸出檔案範例
![trcMergeCourse](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/trcMergeCourse.png)  

## 合併數位課綱資料
### 規則
- 數位課綱資料表須包含**教師姓名、所屬單位、科目序號、科目名稱**及**影片問題**欄位
- 檔案格式必須為.xlsx
### 介面
![merge_SyllabusVideo](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/merge_SyllabusVideo.png)  
點選`合併>數位課綱資料`後會跳出檔案選擇視窗  
選擇數位課綱檔案並於後方指定工作表  
另外可指定是否匯入教師名單，提供所屬單位資料  
點選確定後程式會合併檔案並匯出到output資料夾中  
![merge_SyllabusVideo_selectfile](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/merge_SyllabusVideo_selectfile.png)  
![merge_SyllabusVideo_selectfile_withTeacherData](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/merge_SyllabusVideo_selectfile_withTeacherData.png)  
### 輸入檔案範例
[合併數位課綱資料範例檔案](https://github.com/vincent87720/TRC/blob/master/assets/exampleInputFile/%E5%90%88%E4%BD%B5%E6%95%B8%E4%BD%8D%E8%AA%B2%E7%B6%B1%E8%B3%87%E6%96%99)  
### 輸出檔案範例
![trcMergeVideo](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/trcMergeVideo.png)  


# 指令
第一層指令包含五種，分別是start、download、split、merge和calculate，各指令又包含了不同的第二層指令，其架構如下  
![指令](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/cmdEng.png)  
![功能](https://github.com/vincent87720/TRC/blob/master/assets/mdImage/cmdCht.png)  
依照所需功能由上往下選擇指令即可，例如"下載教師資料"的指令為"trc download teacher"、"合併製版數登記表"的指令為"trc merge course"  

## 下載教師資料
### 下載資料
- 輸入指令`trc download teacher`可下載教師資料
### 參數
無  

## 下載數位課綱連結
### 下載新檔案
- 必須指定-key參數
- 建議指定-year、-semester參數
- 輸入指令`trc download video`+參數可下載數位課綱影片資料
### 合併影片資料到原有檔案
- 必須指定-key、-append參數
- 需指定輸入檔案，輸入檔案需包含**科目序號**欄位
- 建議指定-year、-semester、-inName參數
- 輸入指令`trc download video`+參數可下載數位課綱影片資料
### 參數
| 參數      | 備註                                                                                              | 預設值        |
| --------- |:------------------------------------------------------------------------------------------------- |:------------- |
| -append   | 在原有檔案內增加影片資訊                                                                          |               |
| -key      | 設定YoutubeAPIKey                                                                                 |               |
| -year     | 設定學年度                                                                                        | 當前學年度    |
| -semester | 設定學期                                                                                          | 當前學期      |
| -inName   | 指定輸入檔案名稱參數                                                                              | 數位課綱.xlsx |
| -inSheet  | 指定輸入檔案工作表參數                                                                            | 工作表        |
| -inPath   | 指定輸入檔案路徑<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案與trc.exe同目錄則無需指定 | 目前所在路徑  |
| -outName  | 指定輸出檔案名稱                                                                                  | 數位課綱.xlsx |
| -outSheet | 指定輸出檔案工作表                                                                                | 工作表        |
| -outPath  | 指定輸出檔案路徑<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案與trc.exe同目錄則無需指定 | 目前所在路徑  |

## 分割預警總表
### 依照指定名稱及路徑
- 須將所有要計算的檔案(.xlsx檔)放置於和trc.exe同個資料夾
- 預警總表檔名須設為**預警總表.xlsx**
- 教師名單檔名須設為**教師名單.xlsx**
- 空白模板檔名須設為**空白.xlsx**
- 預警總表、教師名單和空白模板的工作表名稱必須為"工作表"
- 輸入指令`trc split scoreAlert`可計算所有檔案
### 自訂名稱及路徑
- 需指定-masterName、-masterPath、-masterSheet、-teacherName、-teacherPath、-teacherSheet、-templateName、-templatePath、-templateSheet參數
- 輸入指令`trc split scoreAlert`+參數即可計算所有檔案
### 參數
| 參數           | 備註                                                                                                      | 預設值        |
| -------------- |:--------------------------------------------------------------------------------------------------------- | ------------- |
| -masterName    | 指定輸入預警總表檔案名稱                                                                                  | 預警總表.xlsx |
| -masterSheet   | 指定預警總表工作表名稱                                                                                    | 工作表        |
| -masterPath    | 指定輸入預警總表檔案路徑參數<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案與trc.exe同目錄則無需指定 | 目前所在路徑  |
| -teacherName   | 指定輸入教師名單檔案名稱                                                                                  | 教師名單.xlsx |
| -teacherSheet  | 指定教師名單工作表名稱                                                                                    | 工作表        |
| -teacherPath   | 指定輸入教師名單檔案路徑參數<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案trc.exe同目錄則無需指定   | 目前所在路徑  |
| -templateName  | 指定空白分表模板檔案名稱                                                                                  | 空白.xlsx     |
| -templateSheet | 指定空白分表模板工作表名稱                                                                                | 工作表        |
| -templatePath  | 指定空白分表模板檔案路徑參數<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案與trc.exe同目錄則無需指定 | 目前所在路徑  |

## 計算尺規成績差分
### 計算大量檔案
- 須將所有要計算的檔案(.xlsx檔)放置於和trc.exe同個資料夾
- 工作表名稱必須為"學系彙整版"
- 輸入指令`trc calculate difference -A`可計算所有檔案

### 計算特定檔案
- 建議指定-inPath、-inName、-inSheet參數
- 輸入指令`trc calculate difference`+參數即可計算所有檔案
### 參數
| 參數      | 備註                                                                                          | 預設值       |
| --------- |:--------------------------------------------------------------------------------------------- | ------------ |
| -inName   | 指定輸入檔案名稱                                                                              | 評分表.xlsx  |
| -inSheet  | 指定輸入檔案工作表                                                                            | 學系彙整版   |
| -inPath   | 指定輸入檔案路徑<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案與trc.exe同目錄則無需指定 | 目前所在路徑 |
| -outName  | 指定輸出檔案名稱                                                                              | 評分表.xlsx  |
| -outSheet | 指定輸出檔案工作表                                                                            | 學系彙整版   |
| -outPath  | 指定輸出檔案路徑<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案與trc.exe同目錄則無需指定 | 目前所在路徑 |

## 合併製版數登記表
### 依照指定名稱及路徑
- 須將所有要計算的檔案(.xlsx檔)放置於和trc.exe同個資料夾
- 預警總表檔名須設為**開課總表.xlsx**
- 開課總表的工作表名稱必須為"工作表"
- 輸入指令`trc merge course`可計算所有檔案
### 自訂名稱及路徑
- 需指定-inName、-inSheet、-inPath參數
- 輸入指令`trc merge course`+參數即可計算所有檔案
### 參數
| 參數 | 備註 | 預設值 |
| ---- | ---- | ------ |
| -inName | 指定輸入檔案名稱 | 開課總表.xlsx |
| -inSheet     | 指定輸入檔案工作表 |工作表|
| -inPath     | 指定輸入檔案路徑<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案與trc.exe同目錄則無需指定 |目前所在路徑|
| -outName     |指定輸出檔案名稱|[MERGENCE]<br>開課總表.xlsx|
| -outSheet     |指定輸出檔案工作表 | 工作表|
| -outPath | 指定輸出檔案路徑<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案與trc.exe同目錄則無需指定 |目前所在路徑 |

## 合併數位課綱資料
### 依照指定名稱及路徑
- 須將所有要計算的檔案(.xlsx檔)放置於和trc.exe同個資料夾
- 數位課綱資料檔名須設為**數位課綱.xlsx**
- 開課總表的工作表名稱必須為"工作表"
- 輸入指令`trc merge video`可計算所有檔案
### 自訂名稱及路徑
- 需指定-inName、-inSheet、-inPath參數
- 輸入指令`trc merge video`+參數即可計算所有檔案
### 使用教師名單檔案填入所屬單位
- 需指定-tfile參數
- 教師資料檔名須設為**教師名單.xlsx**
- 輸入指令`trc merge video -tfile`+參數即可合併檔案
### 參數
| 參數 | 備註 | 預設值 |
| ---- | ---- | ------ |
|-inName | 指定輸入檔案名稱 | 數位課綱.xlsx |
|-inSheet | 指定輸入檔案工作表|工作表|
|-inPath | 指定輸入檔案路徑<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案與trc.exe同目錄則無需指定|目前所在路徑|
|-outName | 指定輸出檔案名稱|[MERGENCE]<br>數位課綱.xlsx|
|-outSheet | 指定輸出檔案工作表|工作表|
|-outPath | 指定輸出檔案路徑<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案與trc.exe同目錄則無需指定|目前所在路徑|
|-tfile | 使用教師名單檔案填入所屬單位|false|
|-tfName | 指定教師名單檔案名稱|數位課綱.xlsx|
|-tfSheet | 指定教師名單檔案工作表|工作表|
|-tfPath | 指定教師名單檔案路徑<br>(點擊右鍵>內容>位置，可得知檔案位置)，<br>若檔案與trc.exe同目錄則無需指定|目前所在路徑|
