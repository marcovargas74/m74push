package main

import (
	"fmt"

	"github.com/marcovargas74/m74push"
)

func main() {
	m74push.SendPushNotification()
	//assert.True(t, isError)
	fmt.Println("SendPush Notification OK")
}
