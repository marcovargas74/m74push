
# m74push
How to send A simple push notification to IOS app, using the Apple Push Notifications Server (APNs)
In folder "howToUse" have a example how to send a push notification. 

### Prerequisites

But the token and certification, is necessary to get from apple developer account.
More information [Developer Apple](https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server)

- Make sure you have [Go](https://golang.org/doc/install) installed and have set your [GOPATH](https://golang.org/doc/code.html#GOPATH).

- Install apns2:

```sh
go get -u github.com/sideshow/apns2
```

## Install

- Install m74push:

```sh
go get -u github.com/marcovargas74/m74push
```

## Example

```go 
import (
	"fmt"

	"github.com/marcovargas74/m74push"
)

func main() {
	m74push.SendPushNotification()
	//assert.True(t, isError)
	fmt.Println("SendPush Notification OK")
}
``` 



