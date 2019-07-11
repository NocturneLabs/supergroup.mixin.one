package main

import (
	"context"
	"sync"
	"time"

	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/MixinNetwork/supergroup.mixin.one/session"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

var db *gorm.DB

type User struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
}

type CourseRecord struct {
	ID         uint64         `gorm:"primary_key" json:"id"`
	Title      string         `gorm:"default:'Untitled'" json:"title" form:"title"`
	SpeakerIDs pq.StringArray `gorm:"type:text[];column:speaker_ids;default:'{}'" json:"-"`
	IsPublic   bool           `sql:"index" gorm:"default:'false'" json:"is_public" form:"is_public"`
	StartedAt  time.Time      `json:"started_at"`
	FinishedAt *time.Time     `sql:"index" json:"finished_at"`
	DeletedAt  *time.Time     `sql:"index" json:"-"`

	CourseMessages []CourseMessage `json:"-"`
	Speakers       []User          `json:"speakers" gorm:"-"`
	SpeakersMap    map[string]User `json:"-" gorm:"-"`
}

type courseRecordRepo struct {
	currentRecording *CourseRecord
	mutex            sync.Mutex
}

var CourseRecordRepo courseRecordRepo

type CourseMessage struct {
	CourseRecordID uint64 `gorm:"primary_key,index:idx_course_record_id_created_at" json:"-"`

	// fields from original mixin supergroup Message struct
	MessageID      string    `gorm:"primary_key" json:"message_id"`
	UserID         string    `json:"-"`
	QuoteMessageID string    `json:"quote_message_id"`
	Data           string    `json:"data"`
	Category       string    `json:"category"`
	CreatedAt      time.Time `gorm:"index:idx_course_record_id_created_at" json:"created_at"`

	CourseRecord CourseRecord `json:"-"`

	Speaker User `json:"speaker"`
}

func (r *CourseRecord) AfterFind() (err error) {
	r.loadSpeakers()
	return nil
}

func (r *CourseRecord) loadSpeakers() {
	r.SpeakersMap = make(map[string]User)
	for _, uid := range r.SpeakerIDs {
		user, err := models.FindUser(session.WithDatabase(context.TODO(), hostDB), uid)
		if err != nil {
			continue
		}
		u := User{
			ID:        uid,
			FullName:  user.FullName,
			AvatarURL: user.AvatarURL,
		}
		r.SpeakersMap[uid] = u
		r.Speakers = append(r.Speakers, u)
	}
}

func (r *CourseRecord) RecordMessage(m CourseMessage) {
	db.Model(r).Association("CourseMessages").Append(m)
}

func (r *CourseRecord) FinishRecording() {
	if db.Model(r).Association("CourseMessages").Count() == 0 {
		log.Infoln("no messages, abort current recording")
		db.Unscoped().Delete(r)
		return
	}

	now := time.Now()
	r.FinishedAt = &now

	speakerIDs := make(map[string]struct{}) // id => fullName
	for _, m := range r.Messages() {
		if m.UserID == mixinClientID {
			continue
		}

		speakerIDs[m.UserID] = struct{}{}
	}

	r.SpeakerIDs = pq.StringArray{}
	for uid, _ := range speakerIDs {
		r.SpeakerIDs = append(r.SpeakerIDs, uid)
	}

	r.loadSpeakers()
	db.Save(r)
}

func (r *CourseRecord) Messages() (messages []CourseMessage) {
	db.Model(r).Association("CourseMessages").Find(&messages)

	for i, m := range messages {
		if user, found := r.SpeakersMap[m.UserID]; found {
			messages[i].Speaker = user
		}
	}
	return
}

func (courseRecordRepo) PublicRecords() (records []CourseRecord) {
	db.Where("is_public is true").Find(&records)
	return
}

func (courseRecordRepo) AllRecords() (records []CourseRecord) {
	db.Where("finished_at is not null").Find(&records)
	return
}

func (courseRecordRepo) CheckConsistency() {
	unfinishedCount := 0
	db.Model(&CourseRecord{}).Where("is_public is false and finished_at is null").Count(&unfinishedCount)

	if unfinishedCount > 1 {
		log.Panicln("course_record database is broken, too many unfinished recording")
	}
}

func (repo *courseRecordRepo) CurrentRecording() *CourseRecord {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	if repo.currentRecording != nil {
		if repo.currentRecording.FinishedAt == nil {
			return repo.currentRecording
		} else {
			repo.currentRecording = nil
		}
	}

	record := &CourseRecord{}
	db.Where("is_public is false and finished_at is null").First(record)
	if record.ID != 0 {
		repo.currentRecording = record
		return record
	} else {
		return nil
	}
}

func (repo *courseRecordRepo) StartRecording() {
	if currentRecording := CourseRecordRepo.CurrentRecording(); currentRecording != nil {
		log.Errorln("unfinished recording found")
		currentRecording.FinishRecording()
	}

	record := &CourseRecord{StartedAt: time.Now()}
	db.Create(record)

	repo.mutex.Lock()
	repo.currentRecording = record
	repo.mutex.Unlock()
}
