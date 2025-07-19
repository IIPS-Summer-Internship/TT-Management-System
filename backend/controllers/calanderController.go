package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
	"tms-server/config"
	"tms-server/models"
)

func GetCalendarSummaryByMonth(c *gin.Context) {
	month := c.Query("month")
	year := c.Query("year")
	semester := c.Query("semester")
	facultyID := c.Query("faculty_id")
	courseID := c.Query("course_id")

	if month == "" || year == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Both 'month' and 'year' query parameters are required.",
		})
		return
	}

	if _, err := strconv.Atoi(month); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'month' parameter. Must be a number."})
		return
	}
	if _, err := strconv.Atoi(year); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'year' parameter. Must be a number."})
		return
	}

	// 👇 New: Join Session → Timetable → Lecture
	query := config.DB.Model(&models.Session{}).
		Joins("JOIN lectures ON lectures.timetable_id = sessions.timetable_id AND lectures.timeslot_id = sessions.timeslot_id").
		Where("EXTRACT(MONTH FROM sessions.date) = ?", month).
		Where("EXTRACT(YEAR FROM sessions.date) = ?", year)

	if semester != "" {
		query = query.Joins("JOIN timetables ON timetables.timetable_id = sessions.timetable_id").
			Where("timetables.semester = ?", semester)
	}

	if facultyID != "" {
		query = query.Where("lectures.faculty_id = ?", facultyID)
	}

	if courseID != "" {
		query = query.Joins("JOIN timetables ON timetables.timetable_id = sessions.timetable_id").
			Where("timetables.course_id = ?", courseID)
	}

	var sessions []models.Session
	err := query.Find(&sessions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sessions"})
		return
	}

	if len(sessions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no sessions found", "data": []gin.H{}})
		return
	}

	type DayStat struct {
		Held      int
		Cancelled int
		Nil       int
	}
	summary := make(map[string]*DayStat)
	for _, s := range sessions {
		key := s.Date.Format("2006-01-02")
		if summary[key] == nil {
			summary[key] = &DayStat{}
		}
		switch s.Status {
		case "held":
			summary[key].Held++
		case "cancelled":
			summary[key].Cancelled++
		default:
			summary[key].Nil++
		}
	}

	result := []gin.H{}
	for dateStr, stat := range summary {
		result = append(result, gin.H{
			"date":            dateStr,
			"total_held":      stat.Held,
			"total_cancelled": stat.Cancelled,
			"no_data":         stat.Nil,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func GetLectureDetailsByDate(c *gin.Context) {
	dateStr := c.Query("date")
	semester := c.Query("semester")
	facultyID := c.Query("faculty_id")
	courseID := c.Query("course_id")

	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required in YYYY-MM-DD format"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	var sessions []models.Session
	if err := config.DB.Where("date = ?", date).Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch sessions"})
		return
	}

	if len(sessions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no sessions found", "data": []gin.H{}})
		return
	}

	// Now: Fetch lectures by composite key pairs
	lecturePairs := make([][2]uint, 0, len(sessions))
	for _, s := range sessions {
		lecturePairs = append(lecturePairs, [2]uint{s.TimetableID, s.TimeslotID})
	}

	var lectures []models.Lecture
	lectureQuery := config.DB.
		Preload("Subject").
		Preload("Faculty").
		Preload("Room").
		Preload("Batch.Course")

	// No WHERE IN composite, so use explicit filter
	for _, pair := range lecturePairs {
		lectureQuery = lectureQuery.Or("timetable_id = ? AND timeslot_id = ?", pair[0], pair[1])
	}

	if semester != "" {
		lectureQuery = lectureQuery.Joins("JOIN timetables ON timetables.timetable_id = lectures.timetable_id").
			Where("timetables.semester = ?", semester)
	}
	if facultyID != "" {
		lectureQuery = lectureQuery.Where("lectures.faculty_id = ?", facultyID)
	}
	if courseID != "" {
		lectureQuery = lectureQuery.Joins("JOIN timetables ON timetables.timetable_id = lectures.timetable_id").
			Where("timetables.course_id = ?", courseID)
	}

	if err := lectureQuery.Find(&lectures).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch lectures"})
		return
	}

	if len(lectures) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no lectures found", "data": []gin.H{}})
		return
	}

	// Build lookup
	type LectureKey struct {
		TimetableID uint
		TimeslotID  uint
	}
	lectureMap := make(map[LectureKey]models.Lecture)
	for _, l := range lectures {
		lectureMap[LectureKey{l.TimetableID, l.TimeslotID}] = l
	}

	result := []gin.H{}
	for _, s := range sessions {
		lk := LectureKey{s.TimetableID, s.TimeslotID}
		lecture, ok := lectureMap[lk]
		if !ok {
			continue
		}
		result = append(result, gin.H{
			"timetable_id":  s.TimetableID,
			"timeslot_id":   s.TimeslotID,
			"subject":       lecture.Subject.Name,
			"faculty":       lecture.Faculty.Name,
			"day_of_week":   lecture.Timeslot.DayOfWeek,
			"start_time":    lecture.Timeslot.StartTime,
			"end_time":      lecture.Timeslot.EndTime,
			"status":        s.Status,
			"semester":      lecture.Timetable.Semester,
			"room":          lecture.Room.Room_Name,
			"course_name":   lecture.Subject.Course.Name,
			"session_id":    s.SessionID,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"date": dateStr,
		"data": result,
	})
}
