package m74push

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/appleboy/gorush/gorush"
	"github.com/mitchellh/mapstructure"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"
)

// Sound sets the aps sound on the payload.
type Sound struct {
	Critical int     `json:"critical,omitempty"`
	Name     string  `json:"name,omitempty"`
	Volume   float32 `json:"volume,omitempty"`
}

// InitAPNSClient use for initialize APNs Client.
func initAPNSClient() error {
	//if gorush.PushConf.Ios.Enabled {
	var err error
	//var authKey *ecdsa.PrivateKey
	var certificateKey tls.Certificate
	//var ext string

	if gorush.PushConf.Ios.KeyPath != "" {
		//ext = filepath.Ext(gorush.PushConf.Ios.KeyPath)
		certificateKey, err = certificate.FromPemFile(gorush.PushConf.Ios.KeyPath, gorush.PushConf.Ios.Password)
		if err != nil {
			gorush.LogError.Error("Cert Error:", err.Error())

			return err
		}
	} else if gorush.PushConf.Ios.KeyBase64 != "" {
		//ext = "." + gorush.PushConf.Ios.KeyType
		key, err := base64.StdEncoding.DecodeString(gorush.PushConf.Ios.KeyBase64)
		if err != nil {
			gorush.LogError.Error("base64 decode error:", err.Error())

			return err
		}

		certificateKey, err = certificate.FromPemBytes(key, gorush.PushConf.Ios.Password)
		if err != nil {
			gorush.LogError.Error("Cert Error:", err.Error())

			return err
		}
	}

	if gorush.PushConf.Ios.Production {
		ApnsClient = apns2.NewClient(certificateKey).Production()
	} else {
		ApnsClient = apns2.NewClient(certificateKey).Development()
	}
	//	}
	//}

	return nil
}

func iosAlertDictionary(payload *payload.Payload, req pushNotification) *payload.Payload {
	// Alert dictionary

	if len(req.Title) > 0 {
		payload.AlertTitle(req.Title)
	}

	if len(req.Alert.Title) > 0 {
		payload.AlertTitle(req.Alert.Title)
	}

	// Apple Watch & Safari display this string as part of the notification interface.
	if len(req.Alert.Subtitle) > 0 {
		payload.AlertSubtitle(req.Alert.Subtitle)
	}

	if len(req.Alert.TitleLocKey) > 0 {
		payload.AlertTitleLocKey(req.Alert.TitleLocKey)
	}

	if len(req.Alert.LocArgs) > 0 {
		payload.AlertLocArgs(req.Alert.LocArgs)
	}

	if len(req.Alert.TitleLocArgs) > 0 {
		payload.AlertTitleLocArgs(req.Alert.TitleLocArgs)
	}

	if len(req.Alert.Body) > 0 {
		payload.AlertBody(req.Alert.Body)
	}

	if len(req.Alert.LaunchImage) > 0 {
		payload.AlertLaunchImage(req.Alert.LaunchImage)
	}

	if len(req.Alert.LocKey) > 0 {
		payload.AlertLocKey(req.Alert.LocKey)
	}

	if len(req.Alert.Action) > 0 {
		payload.AlertAction(req.Alert.Action)
	}

	if len(req.Alert.ActionLocKey) > 0 {
		payload.AlertActionLocKey(req.Alert.ActionLocKey)
	}

	// General
	if len(req.Category) > 0 {
		payload.Category(req.Category)
	}

	if len(req.Alert.SummaryArg) > 0 {
		payload.AlertSummaryArg(req.Alert.SummaryArg)
	}

	if req.Alert.SummaryArgCount > 0 {
		payload.AlertSummaryArgCount(req.Alert.SummaryArgCount)
	}

	return payload
}

// GetIOSNotification use for define iOS notification.
// The iOS Notification Payload
// ref: https://developer.apple.com/library/content/documentation/NetworkingInternet/Conceptual/RemoteNotificationsPG/PayloadKeyReference.html#//apple_ref/doc/uid/TP40008194-CH17-SW1
func getIOSNotification(req pushNotification) *apns2.Notification {
	notification := &apns2.Notification{
		ApnsID:     req.ApnsID,
		Topic:      req.Topic,
		CollapseID: req.CollapseID,
	}

	if req.Expiration > 0 {
		notification.Expiration = time.Unix(req.Expiration, 0)
	}

	if len(req.Priority) > 0 && req.Priority == "normal" {
		notification.Priority = apns2.PriorityLow
	}

	payload := payload.NewPayload()

	// add alert object if message length > 0
	if len(req.Message) > 0 {
		payload.Alert(req.Message)
	}

	/*
		// zero value for clear the badge on the app icon.
		if req.Badge != nil && *req.Badge >= 0 {
			payload.Badge(*req.Badge)
		}*/

	if req.MutableContent {
		payload.MutableContent()
	}

	switch req.Sound.(type) {
	// from http request binding
	case map[string]interface{}:
		result := &Sound{}
		_ = mapstructure.Decode(req.Sound, &result)
		payload.Sound(result)
	// from http request binding for non critical alerts
	case string:
		payload.Sound(&req.Sound)
	case Sound:
		payload.Sound(&req.Sound)
	}

	if len(req.SoundName) > 0 {
		payload.SoundName(req.SoundName)
	}

	if req.SoundVolume > 0 {
		payload.SoundVolume(req.SoundVolume)
	}

	if req.ContentAvailable {
		payload.ContentAvailable()
	}

	if len(req.URLArgs) > 0 {
		payload.URLArgs(req.URLArgs)
	}

	if len(req.ThreadID) > 0 {
		payload.ThreadID(req.ThreadID)
	}

	for k, v := range req.Data {
		payload.Custom(k, v)
	}

	payload = iosAlertDictionary(payload, req)

	notification.Payload = payload

	return notification
}

func getApnsClient(req pushNotification) (client *apns2.Client) {
	/*if req.Production {
		client = ApnsClient.Production()
	} else if req.Development {
		client = ApnsClient.Development()
	} else {
		if gorush.PushConf.Ios.Production {
			client = ApnsClient.Production()
		} else {
			client = ApnsClient.Development()
		}
	}*/
	client = ApnsClient.Development()
	if req.Production {
		client = ApnsClient.Production()
	}

	log.Println("getApnsClient product ", req.Production, "ios ", gorush.PushConf.Ios.Production, client.Host)
	//log.Printf("getApnsClient %v %v", req.Production, gorush.PushConf.Ios.Production)
	return
}

// PushToIOS provide send notification to APNs server.
func pushToIOS(req pushNotification) bool {
	//gorush.LogAccess.Debug("Start push notification for iOS")
	log.Println("Start push notification for iOS")
	if gorush.PushConf.Core.Sync {
		defer req.WaitDone()
	}

	var (
		retryCount = 0
		maxRetry   = gorush.PushConf.Ios.MaxRetry
	)

	if req.Retry > 0 && req.Retry < maxRetry {
		maxRetry = req.Retry
	}

Retry:
	var (
		isError   = false
		newTokens []string
	)

	notification := getIOSNotification(req)
	client := getApnsClient(req)

	for _, token := range req.Tokens {
		notification.DeviceToken = token

		// send ios notification
		res, err := client.Push(notification)

		if err != nil {
			// apns server error
			//gorush.LogPush(gorush.FailedPush, token, req, err)
			log.Println(gorush.FailedPush)
			/*MARCO if gorush.PushConf.Core.Sync {
				req.AddLog(getLogPushEntry(gorush.FailedPush, token, req, err))
			}*/
			gorush.StatStorage.AddIosError(1)
			newTokens = append(newTokens, token)
			isError = true
			continue
		}

		if res.StatusCode != 200 {
			// error message:
			// ref: https://github.com/sideshow/apns2/blob/master/response.go#L14-L65
			//gorush.LogPush(gorush.FailedPush, token, req, errors.New(res.Reason))
			log.Println(errors.New(res.Reason))
			/*marco if gorush.PushConf.Core.Sync {
				req.AddLog(getLogPushEntry(FailedPush, token, req, errors.New(res.Reason)))
			}*/
			gorush.StatStorage.AddIosError(1)
			newTokens = append(newTokens, token)
			isError = true
			continue
		}

		if res.Sent() {
			//gorush.LogPush(gorush.SucceededPush, token, req, nil)
			log.Println(gorush.SucceededPush)
			gorush.StatStorage.AddIosSuccess(1)
		}
	}

	if isError && retryCount < maxRetry {
		retryCount++

		// resend fail token
		req.Tokens = newTokens
		goto Retry
	}

	return isError
}
