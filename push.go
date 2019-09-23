package m74push

import (
	"log"

	"github.com/appleboy/gorush/config"
	"github.com/appleboy/gorush/gorush"
	"github.com/appleboy/gorush/storage/memory"
)

// InitAppStatus for initialize app status
func initAppStatus() error {
	//LogAccess.Debug("Init App Status Engine as ", PushConf.Stat.Engine)

	gorush.StatStorage = memory.New()

	if err := gorush.StatStorage.Init(); err != nil {
		//gorush.LogError.Error("storage error: " + err.Error())
		log.Println("storage error: " + err.Error())
		return err
	}

	return nil
}

//SendPushNotification ..
func SendPushNotification() {
	gorush.PushConf, _ = config.LoadConf("")

	gorush.PushConf.Ios.Enabled = true
	gorush.PushConf.Ios.Production = true

	gorush.PushConf.Ios.KeyPath = "../cert/mobiliti2in1.pem"
	//err := gorush.InitAPNSClient()
	err := initAPNSClient()
	if err != nil {
		log.Fatal(err)
	}

	//assert.Nil(t, err)
	//err = gorush.InitAppStatus()
	err = initAppStatus()
	//assert.Nil(t, err)
	if err != nil {
		log.Fatal(err)
	}

	req := pushNotification{
		//42c4741429884981d5e949cfb44e571e9fbfecece823565fca9ffb410f558406
		Tokens: []string{"11aa01229f15f0f0c52029d8cf8cd0aeaf2365fe4cebc4af26cd6d76b7919ef7"},
		//Tokens:   []string{"42c4741429884981d5e949cfb44e571e9fbfecece823565fca9ffb410f558406"},
		Platform:   1,
		Production: false,
		//Development: true,
		Message: "Desconsidere Mensagem",
	}

	log.Println(req)
	// send fail
	//isError := gorush.PushToIOS(req)
	isError := pushToIOS(req)

	//assert.True(t, isError)
	if isError {
		log.Printf("Can Send Push Notifications %v: ", isError)
	}

	//FIM */
}
