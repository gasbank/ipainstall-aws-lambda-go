package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"strings"
)

func Handler(request events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {

	plistUrl := request.QueryStringParameters["plistUrl"]

	path := request.RawPath
	if strings.HasSuffix(path, ".plist") {
		splitPath := strings.Split(path, "/")
		ipaPath := "https://" + strings.Join(splitPath[3:len(splitPath)-4], "/")
		bundleIdentifier := splitPath[len(splitPath) - 4]
		version := splitPath[len(splitPath) - 3]
		name := splitPath[len(splitPath) - 2]
		bodyFmt := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>items</key>
        <array>
            <dict>
                <key>assets</key>
                <array>
                    <dict>
                        <key>kind</key>
                        <string>software-package</string>
                        <key>url</key>
                        <string>%s</string>
                    </dict>
                </array>
                <key>metadata</key>
                <dict>
                    <key>bundle-identifier</key>
                    <string>%s</string>
                    <key>bundle-version</key>
                    <string>%s</string>
                    <key>kind</key>
                    <string>software</string>
                    <key>title</key>
                    <string>%s</string>
                </dict>
            </dict>
        </array>
    </dict>
</plist>
`
		body := fmt.Sprintf(bodyFmt, ipaPath, bundleIdentifier, version, name)

		resp := &events.APIGatewayProxyResponse{Headers: map[string]string{"Content-Type": "application/xml"}}
		resp.StatusCode = 200
		resp.Body = body
		return resp, nil
	} else if strings.HasSuffix(path, "/ipaInstall") {
		redirectUrl := "itms-services://?action=download-manifest&url=" + plistUrl
		resp := &events.APIGatewayProxyResponse{Headers: map[string]string{"Location": redirectUrl}}
		resp.StatusCode = 307
		return resp, nil
	} else {
		resp := &events.APIGatewayProxyResponse{}
		resp.Body = fmt.Sprintf("Unknown - plistUrl:%s / path:%s", plistUrl, path)
		resp.StatusCode = 200
		return resp, nil
	}
}

func main() {
	lambda.Start(Handler)
}