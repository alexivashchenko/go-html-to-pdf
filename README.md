
# An app to convert Badges HTML to PDF

## Build

### Linux

`$Env:GOOS = "linux"; $Env:GOARCH = "386"; go build`

### Windows

`$Env:GOOS = "windows"; $Env:GOARCH = "386"; go build`


## Usage

### Linux

`cd /home/www-data/laravel/go-html-to-pdf && ./go-html-to-pdf -input_file="/home/www-data/laravel/storage/tmp.html" -output_file="/home/www-data/laravel/storage/tmp.pdf" -width=410 -height=610 -chrome="/home/www-data/.cache/puppeteer/chrome/linux-128.0.6613.119/chrome-linux64/chrome"`

### Windows

`cd W:\SERVER\files\accreditation-app\go-html-to-pdf\; Start-Process -FilePath "go-html-to-pdf.exe" -ArgumentList "-input_file","W:\SERVER\files\accreditation-app\storage\tmp.html","-output_file","W:\SERVER\files\accreditation-app\storage\tmp.pdf","-width",400,"-height",600,"-chrome","C:\Users\lex\.cache\puppeteer\chrome\win64-126.0.6478.182\chrome-win64\chrome.exe"`

`go run main.go -input_file="W:\SERVER\files\accreditation-app\storage\tmp.html" -output_file="W:\SERVER\files\accreditation-app\storage\tmp.pdf" -width=410 -height=610 -chrome="C:\Users\lex\.cache\puppeteer\chrome\win64-126.0.6478.182\chrome-win64\chrome.exe"`
