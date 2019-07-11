package main

import (
	"testing"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	var err error
	db, err = gorm.Open("postgres", "postgres://qinix:@localhost/course_record_test?sslmode=disable")
	if err != nil {
		log.Panicln(err)
	}
	db.AutoMigrate(&CourseRecord{}, &CourseMessage{})
}

func TestCourseRecordRepo(t *testing.T) {
	db.Unscoped().Delete(CourseRecord{})
	db.Unscoped().Delete(CourseMessage{})
	CourseRecordRepo.CheckConsistency()
	assert.Nil(t, CourseRecordRepo.CurrentRecording())

	CourseRecordRepo.StartRecording()
	CourseRecordRepo.CurrentRecording().RecordMessage(CourseMessage{MessageID: "1", UserID: "1", Data: "1"})
	CourseRecordRepo.CurrentRecording().RecordMessage(CourseMessage{MessageID: "2", UserID: "2", Data: "2"})
	CourseRecordRepo.CurrentRecording().FinishRecording()

	assert.Nil(t, CourseRecordRepo.CurrentRecording())

	records := CourseRecordRepo.AllRecords()
	assert.Len(t, records, 1)
	// t.Log(records)

	messages := records[0].Messages()
	assert.Len(t, messages, 2)
	// t.Log(messages)
}
