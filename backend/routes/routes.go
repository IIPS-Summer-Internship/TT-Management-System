package routes

import (
	"tms-server/config"
	"tms-server/controllers"
	"tms-server/middleware"
	"tms-server/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine) {
	r.Use(middleware.CORSMiddleware())

	db := config.DB
	api := r.Group("/api/v1")

	// Public routes
	api.GET("/ping", controllers.Ping)
	api.POST("/login", controllers.Login)

	// Protected routes for faculty and above
	api.Use(middleware.JWTAuthMiddleware())
	api.POST("/logout", controllers.Logout)
	registerFacultyRoutes(api, db)

	// Admin-only routes
	admin := api.Group("/")
	admin.Use(middleware.RoleAuthMiddleware("admin", "superadmin"))
	registerAdminRoutes(admin, db)

	// Superadmin-only routes
	super := api.Group("/")
	super.Use(middleware.RoleAuthMiddleware("superadmin"))
	registerSuperAdminRoutes(super, db)
}

func registerFacultyRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Courses
	r.GET("/course", controllers.All[models.Course](db))
	r.GET("/course/:id", controllers.Get[models.Course](db))

	// Subjects
	r.GET("/subject", controllers.All[models.Subject](db))
	r.GET("/subject/:id", controllers.Get[models.Subject](db))

	// Faculty
	r.GET("/faculty", controllers.All[models.Faculty](db))
	r.GET("/faculty/:id", controllers.Get[models.Faculty](db))

	// Rooms
	r.GET("/room", controllers.All[models.Room](db))
	r.GET("/room/:id", controllers.Get[models.Room](db))

	// Batch
	r.GET("/batch", controllers.All[models.Batch](db))
	r.GET("/batch/:id", controllers.Get[models.Batch](db))

	//Timeslots
	r.GET("/timeslot", controllers.All[models.Timeslot](db))
	r.GET("/timeslot/:id", controllers.Get[models.Timeslot](db))

	// Timetables
	r.GET("/timetable", controllers.All[models.Timetable](db))
	r.GET("/timetable/:id", controllers.Get[models.Timetable](db))

	// Sections
	r.GET("/section", controllers.All[models.Section](db))
	r.GET("/section/:id", controllers.Get[models.Section](db))

	// Session Notes
	r.GET("/session_note", controllers.All[models.SessionNote](db))
	r.GET("/session_note/:id", controllers.Get[models.SessionNote](db))

	// Lectures (Composite Key!)
	r.GET("/lecture", controllers.QueryLectures(db))
	r.GET("/lecture/query", controllers.QueryLectures(db))
	r.GET("/lecture/:timetable_id/:timeslot_id", controllers.GetLectureByCompositeKey(db))

	//  Sessions
	r.GET("/session", controllers.All[models.Session](db))
	r.GET("/session/:id", controllers.Get[models.Session](db))

	// Calendar endpoints (updated join logic)
	r.GET("/calendar", controllers.GetCalendarSummaryByMonth)
	r.GET("/calendar/day", controllers.GetLectureDetailsByDate)
}

func registerAdminRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// CRUD for each core table
	r.POST("/course", controllers.Create[models.Course](db))
	r.PUT("/course/:id", controllers.Update[models.Course](db))
	r.DELETE("/course/:id", controllers.Delete[models.Course](db))

	r.POST("/subject", controllers.Create[models.Subject](db))
	r.PUT("/subject/:id", controllers.Update[models.Subject](db))
	r.DELETE("/subject/:id", controllers.Delete[models.Subject](db))

	r.POST("/faculty", controllers.Create[models.Faculty](db))
	r.PUT("/faculty/:id", controllers.Update[models.Faculty](db))
	r.DELETE("/faculty/:id", controllers.Delete[models.Faculty](db))

	r.POST("/room", controllers.Create[models.Room](db))
	r.PUT("/room/:id", controllers.Update[models.Room](db))
	r.DELETE("/room/:id", controllers.Delete[models.Room](db))

	r.POST("/batch", controllers.Create[models.Batch](db))
	r.PUT("/batch/:id", controllers.Update[models.Batch](db))
	r.DELETE("/batch/:id", controllers.Delete[models.Batch](db))

	r.POST("/timeslot", controllers.Create[models.Timeslot](db))
	r.PUT("/timeslot/:id", controllers.Update[models.Timeslot](db))
	r.DELETE("/timeslot/:id", controllers.Delete[models.Timeslot](db))

	r.POST("/timetable", controllers.Create[models.Timetable](db))
	r.PUT("/timetable/:id", controllers.Update[models.Timetable](db))
	r.DELETE("/timetable/:id", controllers.Delete[models.Timetable](db))

	r.POST("/section", controllers.Create[models.Section](db))
	r.PUT("/section/:id", controllers.Update[models.Section](db))
	r.DELETE("/section/:id", controllers.Delete[models.Section](db))

	r.POST("/session_note", controllers.Create[models.SessionNote](db))
	r.PUT("/session_note/:id", controllers.Update[models.SessionNote](db))
	r.DELETE("/session_note/:id", controllers.Delete[models.SessionNote](db))

	// Lecture CRUD with composite key handling
	r.POST("/lecture", controllers.CreateLecture(db))
	r.PUT("/lecture/:timetable_id/:timeslot_id", controllers.UpdateLectureByCompositeKey(db))
	r.DELETE("/lecture/:timetable_id/:timeslot_id", controllers.DeleteLectureByCompositeKey(db))

	r.POST("/session", controllers.Create[models.Session](db))
	r.PUT("/session/:id", controllers.Update[models.Session](db))
	r.DELETE("/session/:id", controllers.Delete[models.Session](db))
}

func registerSuperAdminRoutes(r *gin.RouterGroup, db *gorm.DB) {

	r.GET("/user", controllers.All[models.User](db))
	r.POST("/user", controllers.Create[models.User](db))
	r.GET("/user/:id", controllers.Get[models.User](db))
	r.PUT("/user/:id", controllers.Update[models.User](db))
	r.DELETE("/user/:id", controllers.Delete[models.User](db))

	r.GET("/role", controllers.All[models.Role](db))
	r.POST("/role", controllers.Create[models.Role](db))
	r.GET("/role/:id", controllers.Get[models.Role](db))
	r.PUT("/role/:id", controllers.Update[models.Role](db))
	r.DELETE("/role/:id", controllers.Delete[models.Role](db))
}
