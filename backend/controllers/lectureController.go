package controllers

import (
	"net/http"
	"strconv"
	"tms-server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// QueryLectures returns a Gin handler function that searches for lectures.
// It allows filtering by course, batch year, section, semester, faculty, and room.
func QueryLectures(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// --- Parameter Extraction ---
		courseIDStr := c.Query("course_id")
		yearStr := c.Query("year") // Maps to Batch's EntryYear
		sectionName := c.Query("section")
		semesterStr := c.Query("semester")
		facultyIDStr := c.Query("faculty_id")
		roomIDStr := c.Query("room_id")

		// --- Base Query Construction ---
		// Start with the Lecture model and preload all necessary associations for a rich response.
		query := db.Model(&models.Lecture{}).
			Preload("Subject").
			Preload("Faculty.User"). // Also fetch user details for the faculty
			Preload("Timeslot").
			Preload("Timetable.Course").
			Preload("Timetable.Batch").
			Preload("Timetable.Section").
			Preload("Timetable.Room")

		// Join with Timetables, as it's central to most filters.
		query = query.Joins("JOIN timetables ON timetables.id = lectures.timetable_id")

		// --- Apply Optional Filters ---

		if semesterStr != "" {
			if semester, err := strconv.Atoi(semesterStr); err == nil {
				query = query.Where("timetables.semester = ?", semester)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'semester' parameter"})
				return
			}
		}

		if facultyIDStr != "" {
			if facultyID, err := strconv.Atoi(facultyIDStr); err == nil {
				query = query.Where("lectures.faculty_id = ?", facultyID)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'faculty_id' parameter"})
				return
			}
		}

		if roomIDStr != "" {
			if roomID, err := strconv.Atoi(roomIDStr); err == nil {
				query = query.Where("timetables.room_id = ?", roomID)
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'room_id' parameter"})
				return
			}
		}

		// For filters related to Batch (year, course), we need to join the batches table.
		if yearStr != "" || courseIDStr != "" {
			query = query.Joins("JOIN batches ON batches.id = timetables.batch_id")

			if yearStr != "" {
				if year, err := strconv.Atoi(yearStr); err == nil {
					query = query.Where("batches.entry_year = ?", year)
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'year' parameter"})
					return
				}
			}

			if courseIDStr != "" {
				if courseID, err := strconv.Atoi(courseIDStr); err == nil {
					query = query.Where("batches.course_id = ?", courseID)
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'course_id' parameter"})
					return
				}
			}
		}

		// For the section filter, we join the sections table.
		if sectionName != "" {
			// This filter will only work for timetables assigned to a specific section.
			query = query.Joins("JOIN sections ON sections.id = timetables.section_id").
				Where("sections.name = ?", sectionName)
		}

		// --- Execute Query ---
		var lectures []models.Lecture
		if err := query.Find(&lectures).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query lectures: " + err.Error()})
			return
		}

		if len(lectures) == 0 {
			c.JSON(http.StatusOK, gin.H{"message": "No lectures found matching the criteria.", "data": []models.Lecture{}})
			return
		}

		// --- Send Response ---
		c.JSON(http.StatusOK, gin.H{"data": lectures})
	}
}
