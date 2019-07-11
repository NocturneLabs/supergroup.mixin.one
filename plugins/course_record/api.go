package main

import (
	"fmt"
	"strconv"

	"github.com/MixinNetwork/supergroup.mixin.one/middlewares"
	"github.com/MixinNetwork/supergroup.mixin.one/models"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type (
	recordsController  struct{}
	messagesController struct{}
)

var (
	RecordsController  recordsController
	MessagesController messagesController
	router             *gin.Engine
)

func init() {
	router = gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	records := router.Group("/course_record/records")
	{
		// GET /course_record/records
		// list all records, return public records only if user is not admin
		records.GET("", RecordsController.Index)

		record := records.Group("/:record_id")
		{
			// PUT/PATCH /course_record/records/:record_id
			// update record title
			// params:
			//   - title: string
			//   - is_public: boolean
			record.PUT("", RecordsController.Update)
			record.PATCH("", RecordsController.Update)

			messages := record.Group("/messages")
			{
				// GET /course_record/records/:record_id/messages
				// list record messages
				messages.GET("", MessagesController.Index)
			}
		}
	}
}

func (recordsController) Index(c *gin.Context) {
	var records []CourseRecord
	if currentUser(c).GetRole() == "admin" {
		records = CourseRecordRepo.AllRecords()
	} else {
		records = CourseRecordRepo.PublicRecords()
	}
	c.JSON(200, records)
}

func (recordsController) Update(c *gin.Context) {
	if currentUser(c).GetRole() != "admin" {
		c.JSON(403, gin.H{
			"error": "forbidden",
		})
		return
	}

	record, err := getRecord(c)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	var request struct {
		Title    *string `json:"title" form:"title"`
		IsPublic *bool   `json:"is_public" form:"is_public"`
	}

	err = c.Bind(&request)
	if err != nil {
		return
	}

	if request.Title != nil {
		record.Title = *request.Title
	}
	if request.IsPublic != nil {
		record.IsPublic = *request.IsPublic
	}
	db.Save(&record)
	c.JSON(200, record)
}

func (messagesController) Index(c *gin.Context) {
	record, err := getRecord(c)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, record.Messages())
}

func currentUser(c *gin.Context) *models.User {
	return middlewares.CurrentUser(c.Request)
}

func getRecord(c *gin.Context) (record CourseRecord, err error) {
	recordID, err := strconv.ParseUint(c.Param("record_id"), 10, 64)
	if err != nil {
		err = fmt.Errorf("wrong record_id")
		return
	}

	db.First(&record, recordID)
	if record.ID == 0 {
		err = fmt.Errorf("can not find course record with id=%l", recordID)
	}
	return
}
