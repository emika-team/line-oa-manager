package message

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/firestore"
	pb "github.com/emika-team/grpc-proto/line-oa/go"
	"github.com/emika-team/line-oa-manager/pkg/dto"
	"github.com/emika-team/line-oa-manager/pkg/firebase"
	"github.com/emika-team/line-oa-manager/pkg/firebase/models"
	"github.com/emika-team/line-oa-manager/pkg/http/httpclient"
	"github.com/emika-team/line-oa-manager/pkg/utils"
	"github.com/google/uuid"
)

func GetContent(ctx context.Context, in *pb.GetMessageContentRequest) (*pb.GetMessageContentResponse, error) {
	messageID := in.MessageId
	channelID := in.ChannelId
	documentRef := firebase.FirestoreClient.Collection("channel").Doc(channelID)
	channel := models.Channel{}
	documentSnapshot, err := documentRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	err = documentSnapshot.DataTo(&channel)
	if err != nil {
		return nil, err
	}
	h := map[string]interface{}{
		"Authorization": fmt.Sprintf("Bearer %s", channel.AccessToken),
	}
	url := fmt.Sprintf("%s/v2/bot/message/%s/content", dto.LineAPIDataEndPoint, messageID)
	result, err := httpclient.HttpRequest("GET", url, nil, nil, h)
	if err != nil {
		return nil, err
	}
	return &pb.GetMessageContentResponse{Result: base64.StdEncoding.EncodeToString(result), Success: true}, nil
}

func SendMessage(ctx context.Context, in *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	channelID := in.ChannelId
	documentRef := firebase.FirestoreClient.Collection("channel").Doc(channelID)
	channel := models.Channel{}
	documentSnapshot, err := documentRef.Get(ctx)
	if err != nil {
		return nil, err
	}
	err = documentSnapshot.DataTo(&channel)
	if err != nil {
		return nil, err
	}

	if len(in.GetMessages()) == 0 {
		return nil, fmt.Errorf("messages is empty")
	}
	mpd := in.GetMessages()[0]
	mbyte, _ := mpd.MarshalJSON()
	message := map[string]interface{}{}
	json.Unmarshal(mbyte, &message)
	message["id"] = uuid.New().String()

	event := interface{}(map[string]interface{}{
		"sender": channel.ChannelID,
		"source": map[string]interface{}{
			"userId": in.To,
		},
		"type":    "message",
		"message": message,
	})
	m := utils.BuildMessage(event, "")
	err = firebase.FirestoreClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		err := m.Create(tx)
		if err != nil {
			return err
		}
		h := map[string]interface{}{
			"Authorization":    fmt.Sprintf("Bearer %s", channel.AccessToken),
			"Content-Type":     "application/json",
			"X-Line-Retry-Key": message["id"].(string),
		}
		url := fmt.Sprintf("%s/v2/bot/message/push", dto.LineAPIEndPoint)
		data := map[string]interface{}{
			"to":       in.To,
			"messages": []interface{}{message},
		}
		_, err = httpclient.HttpRequest("POST", url, data, nil, h)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pb.SendMessageResponse{Success: true}, nil
}
