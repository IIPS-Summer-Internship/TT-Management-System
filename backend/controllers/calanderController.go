package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"tms-server/config"
	"tms-server/models"

	"github.com/gin-gonic/gin"
)

// GetCalendarSummaryByMonth provides a daily summary of session statuses for a given month.
func GetCalendarSummaryByMonth(c *gin.Context) {
	// --- Parameter Extraction and Validation ---
	monthStr := c.Query("month")
	yearStr := c.Query("year")
	semesterStr := c.Query("semester")
	facultyIDStr := c.Query("faculty_id")
	courseIDStr := c.Query("course_id")

	if monthStr == "" || yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both 'month' and 'year' query parameters are required."})
		return
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'month' parameter. Must be a number between 1 and 12."})
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'year' parameter. Must be a number."})
		return
	}

	// --- Database Query Construction ---
	query := config.DB.Model(&models.Session{}).
		Where("EXTRACT(MONTH FROM date) = ?", month).
		Where("EXTRACT(YEAR FROM date) = ?", year)

	// Apply optional filters by joining related tables
	if semesterStr != "" || courseIDStr != "" {
		query = query.Joins("JOIN timetables ON timetables.id = sessions.timetable_id")
		if semesterStr != "" {
			query = query.Where("timetables.semester = ?", semesterStr)
		}
		if courseIDStr != "" {
			query = query.Where("timetables.course_id = ?", courseIDStr)
		}
	}

	if facultyIDStr != "" {
		// Join lectures based on the composite key (timetable_id, timeslot_id)
		query = query.Joins("JOIN lectures ON lectures.timetable_id = sessions.timetable_id AND lectures.timeslot_id = sessions.timeslot_id").
			Where("lectures.faculty_id = ?", facultyIDStr)
	}

	var sessions []models.Session
	if err := query.Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch session summary."})
		return
	}

	if len(sessions) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No sessions found for the given criteria.", "data": []gin.H{}})
		return
	}

	// --- Data Aggregation ---
	type DayStat struct {
		Held      int `json:"total_held"`
		Cancelled int `json:"total_cancelled"`
		NoData    int `json:"no_data"`
	}
	summary := make(map[string]*DayStat)

	for _, s := range sessions {
		key := s.Date.Format("2006-01-02")
		if _, exists := summary[key]; !exists {
			summary[key] = &DayStat{}
		}
		switch s.Status {
		case "held":
			summary[key].Held++
		case "cancelled":
			summary[key].Cancelled++
		default:
			summary[key].NoData++
		}
	}

	// --- Response Formatting ---
	result := []gin.H{}
	for dateStr, stat := range summary {
		result = append(result, gin.H{
			"date":            dateStr,
			"total_held":      stat.Held,
			"total_cancelled": stat.Cancelled,
			"no_data":         stat.NoData,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetLectureDetailsByDate retrieves detailed information for all lectures scheduled on a specific date.
func GetLectureDetailsByDate(c *gin.Context) {
	// --- Parameter Extraction and Validation ---
	dateStr := c.Query("date")
	semesterStr := c.Query("semester")
	facultyIDStr := c.Query("faculty_id")
	courseIDStr := c.Query("course_id")

	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Date is required in YYYY-MM-DD format."})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, use YYYY-MM-DD."})
		return
	}

	// --- Query Scheduled Lectures for the Day of the Week ---
	// The model uses 1 for Sunday, 2 for Monday, etc. Go's time.Weekday is Sunday=0.
	dayOfWeek := int(date.Weekday()) + 1

	lectureQuery := config.DB.Model(&models.Lecture{}).
		Joins("JOIN timeslots ON timeslots.id = lectures.timeslot_id").
		Joins("JOIN timetables ON timetables.id = lectures.timetable_id").
		Where("timeslots.day_of_week = ?", dayOfWeek)

	// Apply optional filters
	if facultyIDStr != "" {
		lectureQuery = lectureQuery.Where("lectures.faculty_id = ?", facultyIDStr)
	}
	if courseIDStr != "" {
		lectureQuery = lectureQuery.Where("timetables.course_id = ?", courseIDStr)
	}
	if semesterStr != "" {
		lectureQuery = lectureQuery.Where("timetables.semester = ?", semesterStr)
	}

	// Preload all necessary related data for the response
	lectureQuery = lectureQuery.
		Preload("Subject").
		Preload("Faculty").
		Preload("Timeslot").
		Preload("Timetable.Course").
		Preload("Timetable.Batch").
		Preload("Timetable.Section")

	var lectures []models.Lecture
	if err := lectureQuery.Find(&lectures).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch scheduled lectures."})
		return
	}

	if len(lectures) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No lectures scheduled for this day.", "data": []gin.H{}})
		return
	}

	// --- Fetch Session Status for the Specific Date ---
	var sessionsOnDate []models.Session
	config.DB.Where("date = ?", date).Find(&sessionsOnDate)

	// Create a map for quick lookup of session status and ID
	sessionMap := make(map[string]models.Session) // Key: "timetableID-timeslotID"
	for _, s := range sessionsOnDate {
		key := fmt.Sprintf("%d-%d", s.TimetableID, s.TimeslotID)
		sessionMap[key] = s
	}

	// --- Combine Scheduled Lectures with Session Data ---
	result := []gin.H{}
	for _, lecture := range lectures {
		key := fmt.Sprintf("%d-%d", lecture.TimetableID, lecture.TimeslotID)
		session, exists := sessionMap[key]

		status := ""
		var sessionID *uint
		if exists {
			status = session.Status
			sid := session.ID // Create a new variable to take its address
			sessionID = &sid
		}

		sectionName := ""
		if lecture.Timetable.Section != nil {
			sectionName = lecture.Timetable.Section.Name
		}

		result = append(result, gin.H{
			"timetable_id":  lecture.TimetableID,
			"timeslot_id":   lecture.TimeslotID,
			"session_id":    sessionID, // Can be null if no session entry exists
			"subject":       lecture.Subject.Name,
			"faculty":       fmt.Sprintf("%s %s", lecture.Faculty.FirstName, lecture.Faculty.LastName),
			"start_time":    lecture.Timeslot.StartTime.Format("15:04"),
			"end_time":      lecture.Timeslot.EndTime.Format("15:04"),
			"status":        status,
			"semester":      lecture.Timetable.Semester,
			"room":          lecture.Room,
			"batch_year":    lecture.Timetable.Batch.EntryYear,
			"batch_section": sectionName,
			"course_name":   lecture.Timetable.Course.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"date": dateStr,
		"data": result,
	})
}
