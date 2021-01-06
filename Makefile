build-fill:
	rm -f lambdas/builds/fill/main
	rm -f lambdas/builds/fill/fill.zip
	GOARCH=amd64 GOOS=linux go build -o lambdas/builds/fill/main lambdas/src/FillEvent.go
	$$USERPROFILE/go/bin/build-lambda-zip.exe -output lambdas/builds/fill/fill.zip lambdas/builds/fill/main

deploy-fill:
	aws lambda update-function-code --function-name FillEvent --region us-west-2  --zip-file fileb://lambdas/builds/fill/fill.zip

test:
	go test ./lambdas/src -v