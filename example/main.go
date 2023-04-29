package main

import (
	"fmt"
	"os"

	"github.com/fvahid/gor_middleware/gormw"
)

func OnRequest(gor *gormw.Gor, msg *gormw.GorMessage, kwargs ...interface{}) *gormw.GorMessage {
	gor.On("response", OnResponse, msg.ID, msg)
	return msg
}

func OnResponse(gor *gormw.Gor, msg *gormw.GorMessage, kwargs ...interface{}) *gormw.GorMessage {
	req, _ := kwargs[0].(*gormw.GorMessage)
	gor.On("replay", OnReplay, req.ID, req, msg)
	return msg
}

func OnReplay(gor *gormw.Gor, msg *gormw.GorMessage, kwargs ...interface{}) *gormw.GorMessage {
	req, _ := kwargs[0].(*gormw.GorMessage)
	resp, _ := kwargs[1].(*gormw.GorMessage)
	fmt.Fprintf(os.Stderr, "request raw http: %s\n", req.HTTP)
	fmt.Fprintf(os.Stderr, "response raw http: %s\n", resp.HTTP)
	fmt.Fprintf(os.Stderr, "replay raw http: %s\n", msg.HTTP)
	respStatus, _ := gormw.HTTPStatus(string(resp.HTTP))
	replayStatus, _ := gormw.HTTPStatus(string(msg.HTTP))
	if respStatus != replayStatus {
		fmt.Fprintf(os.Stderr, "replay status [%s] diffs from response status [%s]\n", replayStatus, respStatus)
	} else {
		fmt.Fprintf(os.Stderr, "replay status is same as response status\n")
	}
	return msg
}

func main() {
	gor := gormw.CreateGor()
	gor.On("request", OnRequest, "")
	gor.Run()
}
