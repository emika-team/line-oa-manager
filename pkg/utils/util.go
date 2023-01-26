package utils

import "github.com/emika-team/line-oa-manager/pkg/firebase/models"

func BuildMessage(event interface{}, destination string) models.Message {
	m := event.(map[string]interface{})["message"].(map[string]interface{})
	userId, ok := event.(map[string]interface{})["source"].(map[string]interface{})["userId"].(string)
	if !ok {
		userId = ""
	}
	mtype, ok := m["type"].(string)
	if !ok {
		mtype = ""
	}
	mtext, ok := m["text"].(string)
	if !ok {
		mtext = ""
	}
	emojis := []models.Emoji{}
	e, ok := m["emojis"].([]interface{})
	if !ok {
		e = []interface{}{}
	}
	for _, e := range e {
		emojis = append(emojis, models.Emoji{
			Index:     int(e.(map[string]interface{})["index"].(float64)),
			ProductID: e.(map[string]interface{})["productId"].(string),
			EmojiID:   e.(map[string]interface{})["emojiId"].(string),
		})
	}
	packageId, ok := m["packageId"].(string)
	if !ok {
		packageId = ""
	}
	stickerId, ok := m["stickerId"].(string)
	if !ok {
		stickerId = ""
	}
	originalContentUrl, ok := m["originalContentUrl"].(string)
	if !ok {
		originalContentUrl = ""
	}
	previewImageUrl, ok := m["previewImageUrl"].(string)
	if !ok {
		previewImageUrl = ""
	}
	trackingId, ok := m["trackingId"].(string)
	if !ok {
		trackingId = ""
	}
	duration, ok := m["duration"].(float64)
	if !ok {
		duration = 0
	}
	sender, ok := event.(map[string]interface{})["sender"].(string)
	if !ok {
		sender = userId
	}
	return models.Message{
		ID:                 m["id"].(string),
		Destination:        destination,
		UID:                userId,
		Type:               mtype,
		Text:               mtext,
		Emojis:             emojis,
		PackageID:          packageId,
		StickerID:          stickerId,
		OriginalContentUrl: originalContentUrl,
		PreviewImageUrl:    previewImageUrl,
		TrackingId:         trackingId,
		Duration:           int(duration),
		Sender:             sender,
	}
}
