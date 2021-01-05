build-fill:
	GOARCH=amd64 GOOS=linux go build -o builds/fill/main src/FillEvent.go
	$$USERPROFILE/go/bin/build-lambda-zip.exe -output builds/fill/fill.zip builds/fill/main

deploy-fill:
	aws lambda update-function-code --function-name event --region us-west-2  --zip-file fileb://builds/fill/fill.zip