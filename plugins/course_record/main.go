package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"

	"github.com/MixinNetwork/supergroup.mixin.one/durable"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/plugin"
)

var (
	mixinClientID string
	hostDB        *durable.Database
)

func PluginInit(ctx *plugin.PluginContext) {
	var err error
	db, err = gorm.Open("postgres", ctx.ConfigMustGet("database_url").(string))
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&CourseRecord{}, &CourseMessage{})
	CourseRecordRepo.CheckConsistency()

	mixinClientID = ctx.MixinClientID()
	hostDB = ctx.HostDB()

	ctx.On(plugin.EventTypeProhibitedStatusChanged, func(s interface{}) {
		if s.(bool) {
			CourseRecordRepo.StartRecording()
		} else {
			if currentRecording := CourseRecordRepo.CurrentRecording(); currentRecording != nil {
				currentRecording.FinishRecording()
			}
		}
	})

	ctx.On(plugin.EventTypeMessageCreated, func(m interface{}) {
		if currentRecording := CourseRecordRepo.CurrentRecording(); currentRecording != nil {
			mixinMessage := m.(models.Message)
			message := CourseMessage{
				MessageID:      mixinMessage.MessageId,
				UserID:         mixinMessage.UserId,
				QuoteMessageID: mixinMessage.QuoteMessageId,
				Data:           mixinMessage.Data,
				Category:       mixinMessage.Category,
				CreatedAt:      mixinMessage.CreatedAt,
			}
			currentRecording.RecordMessage(message)
		}
	})

	ctx.RegisterHTTPHandler("course_record", router)
}
