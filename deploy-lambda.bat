set GOOS=linux
go build -o main main.go
build-lambda-zip -output main.zip main
aws lambda update-function-code --function-name ipaInstall --zip-file fileb://main.zip --profile default --region ap-northeast-2