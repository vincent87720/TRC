
# TRC Program

此程式提供給教資中心使用，程式包含教師資料下載、預警資料分割、製版數統計表合併、數位課綱資料合併以及尺規評量差分計算。

# 下載

點擊以下連結下載，或進入bin資料夾下載執行檔  
[Windows版本](https://github.com/vincent87720/TRC/raw/ver2.0/bin/windows/trc.exe)
[MacOS版本](https://github.com/vincent87720/TRC/raw/ver2.0/bin/darwin/trc)
[Linux版本](https://github.com/vincent87720/TRC/raw/ver2.0/bin/linux/trc)

![進入bin目錄](https://github.com/vincent87720/TRC/blob/ver2.0/assets/mdImage/bin.png)


# 開始

[教學影片](https://youtu.be/hKECNRl-Zqs)

# 指令

第一層指令包含五種，分別是start、download、split、merge和calculate，各指令又包含了不同的第二層指令，其架構如下

![指令](https://github.com/vincent87720/TRC/blob/ver2.0/assets/mdImage/cmdEng.png)

![功能](https://github.com/vincent87720/TRC/blob/ver2.0/assets/mdImage/cmdCht.png)

  

依照所需功能由上往下選擇指令即可，例如"下載教師資料"的指令為"trc download teacher"、"合併製版數登記表"的指令為"trc merge course"

  

## 啟動圖形化介面
**指令**
- 輸入指令`trc start gui`可啟動圖形化使用者介面
![圖形化使用者介面](https://github.com/vincent87720/TRC/blob/ver2.0/assets/mdImage/gui.png)


## 下載教師資料

  

**規則**

無

**下載資料**

- 輸入指令`trc download teacher`可下載教師資料

  

**參數**

無

**範例**
![trcDownloadTeacher](https://github.com/vincent87720/TRC/blob/ver2.0/assets/mdImage/trcDownloadTeacher.png)

  

## 下載數位課綱連結

**規則**
 - 若有指定參數會依照參數設定值運行
 - 若無指定參數或依照預設設定值運行
 - 檔案格式必須為.xlsx

**下載新檔案**
- 必須指定-key參數
- 建議指定-year、-semester參數
 - 輸入指令`trc download video`+參數可下載數位課綱影片資料

  **合併影片資料到原有檔案**
- 必須指定-key、-append參數
- 需指定輸入檔案，輸入檔案需包含**科目序號**欄位
- 建議指定-year、-semester、-inName參數
 - 輸入指令`trc download video`+參數可下載數位課綱影片資料

**參數**

 - -**append**在原有檔案內增加影片資訊
 - -**key**設定YoutubeAPIKey
 - -**year**設定學年度  
預設值：當前學年度  

 - -**semester**設定學期  
預設值：當前學期  

 - -**inName** 指定輸入檔案名稱參數  
預設值：數位課綱.xlsx  
指定方式為`-inName 檔案名稱.xlsx` ex: -inName 數位課綱.xlsx

 -  **-inSheet** 指定輸入檔案工作表參數  
預設值：工作表  
指定方式為`-inSheet 工作表名稱` ex: -inSheet 工作表

 -  **-inPath** 指定輸入檔案路徑(點擊右鍵>內容>位置，可得知檔案位置)，若檔案與trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-inPath 路徑` ex: -inPath C:\Users\User\Desktop\TRCProgram\

 - -**outName** 指定輸出檔案名稱  
預設值：數位課綱.xlsx  
指定方式為`-outName 檔案名稱.xlsx` ex: -outName 數位課綱.xlsx

 -  **-outSheet** 指定輸出檔案工作表  
預設值：工作表  
指定方式為`-outSheet 工作表名稱` ex: -outSheet 工作表

 -  **-outPath** 指定輸出檔案路徑(點擊右鍵>內容>位置，可得知檔案位置)，若檔案與trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-outPath 路徑` ex: -outPath C:\Users\User\Desktop\TRCProgram\

**範例**
![trcDownloadVideo](https://github.com/vincent87720/TRC/blob/ver2.0/assets/mdImage/trcDownloadVideo.png)

## 分割預警總表

  

**規則**

- 預警總表須包含**開課學系**、**科目序號**、**預警科目**、**授課教師**、**學號**、**學生姓名**及**預警原由**欄位

- 教師名單須包含**學院編號**、**教師編號**、**教師姓名**及**所屬單位名稱**欄位

- 檔案格式必須為.xlsx

**依照指定名稱及路徑**

- 須將所有要計算的檔案(.xlsx檔)放置於和trc.exe同個資料夾

- 預警總表檔名須設為**預警總表.xlsx**

- 教師名單檔名須設為**教師名單.xlsx**

- 空白模板檔名須設為**空白.xlsx**

- 預警總表、教師名單和空白模板的工作表名稱必須為"工作表"

- 輸入指令`trc split scoreAlert`可計算所有檔案

  

**自訂名稱及路徑**

  

- 需指定-masterName、-masterPath、-masterSheet、-teacherName、-teacherPath、-teacherSheet、-templateName、-templatePath、-templateSheet參數

- 輸入指令`trc split scoreAlert`+參數即可計算所有檔案

  

**參數**

  

- -**masterName** 指定輸入預警總表檔案名稱  
預設值：預警總表.xlsx  
指定方式為`-masterName 檔案名稱.xlsx` ex: -masterName 預警總表.xlsx

- -**masterSheet** 指定預警總表工作表名稱  
預設值：工作表  
指定方式為`-masterSheet 工作表` ex: -masterSheet 工作表

- -**masterPath** 指定輸入預警總表檔案路徑參數(點擊右鍵>內容>位置，可得知檔案位置)，若檔案與trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-masterPath 路徑` ex: -masterPath C:\Users\User\Desktop\TRCProgram\

- -**teacherName** 指定輸入教師名單檔案名稱  
預設值：教師名單.xlsx  
指定方式為`-teacherName 檔案名稱.xlsx` ex: -teacherName 教師名單.xlsx

- -**teacherSheet** 指定教師名單工作表名稱  
預設值：工作表  
指定方式為`-teacherSheet 工作表` ex: -teacherSheet 工作表

- -**teacherPath** 指定輸入教師名單檔案路徑參數(點擊右鍵>內容>位置，可得知檔案位置)，若檔案trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-teacherPath 路徑` ex: -teacherPath C:\Users\User\Desktop\TRCProgram\

- -**templateName** 指定空白分表模板檔案名稱  
預設值：空白.xlsx  
指定方式為`-templateName 檔案名稱.xlsx` ex: -templateName 空白.xlsx

- -**templateSheet** 指定空白分表模板工作表名稱  
預設值：工作表  
指定方式為`-templateSheet 工作表` ex: -templateSheet 工作表

- -**templatePath** 指定空白分表模板檔案路徑參數(點擊右鍵>內容>位置，可得知檔案位置)，若檔案與trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-templatePath 路徑` ex: -templatePath C:\Users\User\Desktop\TRCProgram\

**輸入檔案範例**     
[分割預警總表範例檔案](https://github.com/vincent87720/TRC/tree/ver2.0/assets/exampleInputFile/%E5%88%86%E5%89%B2%E9%A0%90%E8%AD%A6%E7%B8%BD%E8%A1%A8)

## 計算尺規成績差分

  

**規則**

- 面向、學生及評委可任意增加或刪減，其餘須按照特定格式

- 檔案格式必須為.xlsx

**計算大量檔案**

- 須將所有要計算的檔案(.xlsx檔)放置於和trc.exe同個資料夾

- 工作表名稱必須為"學系彙整版"

- 輸入指令`trc calculate difference -A`可計算所有檔案

  

**計算特定檔案**

- 建議指定-inPath、-inName、-inSheet參數

- 輸入指令`trc calculate difference`+參數即可計算所有檔案

  

**參數**

- -**inName** 指定輸入檔案名稱  
預設值：評分表.xlsx  
指定方式為`-inName 檔案名稱.xlsx` ex: -inName 評分表.xlsx

-  **-inSheet** 指定輸入檔案工作表  
預設值：學系彙整版  
指定方式為`-inSheet 工作表名稱` ex: -inSheet 學系彙整版

-  **-inPath** 指定輸入檔案路徑(點擊右鍵>內容>位置，可得知檔案位置)，若檔案與trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-inPath 路徑` ex: -inPath C:\Users\User\Desktop\TRCProgram\
- -**outName** 指定輸出檔案名稱  
預設值：評分表.xlsx  
指定方式為`-outName 檔案名稱.xlsx` ex: -outName 評分表.xlsx

-  **-outSheet** 指定輸出檔案工作表  
預設值：學系彙整版  
指定方式為`-outSheet 工作表名稱` ex: -outSheet 學系彙整版

-  **-outPath** 指定輸出檔案路徑(點擊右鍵>內容>位置，可得知檔案位置)，若檔案與trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-outPath 路徑` ex: -outPath C:\Users\User\Desktop\TRCProgram\

**輸入檔案範例**     
[計算尺規成績差分範例檔案](https://github.com/vincent87720/TRC/tree/ver2.0/assets/exampleInputFile/%E8%A8%88%E7%AE%97%E5%B0%BA%E8%A6%8F%E6%88%90%E7%B8%BE%E5%B7%AE%E5%88%86)

**範例**
![trcCalculateDifference](https://github.com/vincent87720/TRC/blob/ver2.0/assets/mdImage/trcCalculateDifference.png)

## 合併製版數登記表

  

**規則**

- 預警總表須包含**教師姓名、專兼任別、開課學制、開課系所、科目序號、系-組-年-班、科目名稱、選修別、學分、時數、星期-時間-教室、選課人數、合班註記、合班序號**及**備註**欄位

- 可同時有多個欄位的名稱為"教師姓名"(第一個以外的教師姓名會放入備註)

- 檔案格式必須為.xlsx

**依照指定名稱及路徑**

- 須將所有要計算的檔案(.xlsx檔)放置於和trc.exe同個資料夾

- 預警總表檔名須設為**開課總表.xlsx**

- 開課總表的工作表名稱必須為"工作表"

- 輸入指令`trc merge course`可計算所有檔案

  

**自訂名稱及路徑**

- 需指定-inName、-inSheet、-inPath參數

- 輸入指令`trc merge course`+參數即可計算所有檔案

  

**參數**

- -**inName** 指定輸入檔案名稱  
預設值：開課總表.xlsx  
指定方式為`-inName 檔案名稱.xlsx` ex: -inName 開課總表.xlsx

-  **-inSheet** 指定輸入檔案工作表  
預設值：工作表  
指定方式為`-inSheet 工作表名稱` ex: -inSheet 工作表

-  **-inPath** 指定輸入檔案路徑(點擊右鍵>內容>位置，可得知檔案位置)，若檔案與trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-inPath 路徑` ex: -inPath C:\Users\User\Desktop\TRCProgram\

- -**outName** 指定輸出檔案名稱  
預設值：[MERGENCE]開課總表.xlsx  
指定方式為`-outName 檔案名稱.xlsx` ex: -outName [MERGENCE]開課總表.xlsx  

-  **-outSheet** 指定輸出檔案工作表  
預設值：工作表  
指定方式為`-outSheet 工作表名稱` ex: -outSheet 工作表

-  **-outPath** 指定輸出檔案路徑(點擊右鍵>內容>位置，可得知檔案位置)，若檔案與trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-outPath 路徑` ex: -outPath C:\Users\User\Desktop\TRCProgram\

**輸入檔案範例**     
[合併製版數登記表範例檔案](https://github.com/vincent87720/TRC/tree/ver2.0/assets/exampleInputFile/%E5%90%88%E4%BD%B5%E8%A3%BD%E7%89%88%E6%95%B8%E7%99%BB%E8%A8%98%E8%A1%A8)

**範例**
![trcMergeCourse](https://github.com/vincent87720/TRC/blob/ver2.0/assets/mdImage/trcMergeCourse.png)

## 合併數位課綱資料

  

**規則**

- 數位課綱資料表須包含**教師姓名、所屬單位、科目序號、科目名稱**及**影片問題**欄位

- 檔案格式必須為.xlsx

**依照指定名稱及路徑**

- 須將所有要計算的檔案(.xlsx檔)放置於和trc.exe同個資料夾

- 數位課綱資料檔名須設為**數位課綱.xlsx**

- 開課總表的工作表名稱必須為"工作表"

- 輸入指令`trc merge video`可計算所有檔案

  

**自訂名稱及路徑**

- 需指定-inName、-inSheet、-inPath參數

- 輸入指令`trc merge video`+參數即可計算所有檔案


**使用教師名單檔案填入所屬單位**

- 需指定-tfile參數

- 教師資料檔名須設為**教師名單.xlsx**

- 輸入指令`trc merge video -tfile`+參數即可合併檔案

  

**參數**

- -**inName** 指定輸入檔案名稱  
預設值：數位課綱.xlsx  
指定方式為`-inName 檔案名稱.xlsx` ex: -inName 數位課綱.xlsx

-  **-inSheet** 指定輸入檔案工作表  
預設值：工作表  
指定方式為`-inSheet 工作表名稱` ex: -inSheet 工作表

-  **-inPath** 指定輸入檔案路徑(點擊右鍵>內容>位置，可得知檔案位置)，若檔案與trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-inPath 路徑` ex: -inPath C:\Users\User\Desktop\TRCProgram\

- -**outName** 指定輸出檔案名稱  
預設值：[MERGENCE]數位課綱.xlsx  
指定方式為`-outName 檔案名稱.xlsx` ex: -outName [MERGENCE]數位課綱.xlsx

-  **-outSheet** 指定輸出檔案工作表  
預設值：工作表  
指定方式為`-outSheet 工作表名稱` ex: -outSheet 工作表

-  **-outPath** 指定輸出檔案路徑(點擊右鍵>內容>位置，可得知檔案位置)，若檔案與trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-outPath 路徑` ex: -outPath C:\Users\User\Desktop\TRCProgram\

- -**tfile** 使用教師名單檔案填入所屬單位  
預設值：false  
指定方式為輸入-tfile

- -**tfName** 指定教師名單檔案名稱  
預設值：數位課綱.xlsx  
指定方式為`-inName 檔案名稱.xlsx` ex: -inName 教師名單.xlsx

-  **-tfSheet** 指定教師名單檔案工作表  
預設值：工作表  
指定方式為`-inSheet 工作表名稱` ex: -inSheet 工作表

-  **-tfPath** 指定教師名單檔案路徑(點擊右鍵>內容>位置，可得知檔案位置)，若檔案與trc.exe同目錄則無需指定  
預設值：目前所在路徑  
指定方式為`-inPath 路徑` ex: -inPath C:\Users\User\Desktop\TRCProgram\

**輸入檔案範例**     
[合併數位課綱資料範例檔案](https://github.com/vincent87720/TRC/tree/ver2.0/assets/exampleInputFile/%E5%90%88%E4%BD%B5%E6%95%B8%E4%BD%8D%E8%AA%B2%E7%B6%B1%E8%B3%87%E6%96%99)

**範例**
![trcMergeVideo](https://github.com/vincent87720/TRC/blob/ver2.0/assets/mdImage/trcMergeVideo.png)