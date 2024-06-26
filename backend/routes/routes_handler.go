package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

func Routes(router *gin.Engine) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Service is up and running!",
		})
	})

	admin := router.Group("/admin")
	{
		admin.GET("/users", controllers.GetAllUsers)
		admin.POST("/tag", controllers.TagUsers)
		admin.POST("/tag/remove", controllers.RemoveTag)
		admin.POST("/users/delete", controllers.DeleteUsers)
		admin.GET("/metrics", controllers.Metrics)
	}

	users := router.Group("/users")
	{
		users.POST("/student", controllers.NewStudent)
		users.POST("/teacher", controllers.NewTeacher)
		users.POST("/teacher-student", controllers.NewProfesorStudent)
		users.GET("/details/:username", controllers.GetUserDetails)

		users.POST("/login", controllers.Login)
		users.POST("/post", controllers.NewPublication)
		users.DELETE("/clear/:username", controllers.ClearPublications)

		users.GET("/relations/:username", controllers.GetUserRelations)
		users.POST("/relations/delete", controllers.DeleteSingleRelation)
		users.DELETE("/relations/delete-all/:username", controllers.DeleteAllRelations)

		users.GET("/recommendation/:username", controllers.Recommendation)
	}

	careers := router.Group("/careers")
	{
		careers.GET("/", controllers.GetCareers)
		careers.POST("/studies", controllers.CreateRelationStudiesCareer)
		careers.POST("/interests", controllers.CreateRelationInterestsCareer)
	}

	signs := router.Group("/signs")
	{
		signs.GET("/", controllers.GetZodiacalSigns)
		signs.POST("/is", controllers.CreateRelationIsSign)
	}

	teams := router.Group("/teams")
	{
		teams.GET("/", controllers.GetTeams)
		teams.POST("/", controllers.NewTeam)
		teams.POST("/likes", controllers.CreateRelationSupportsTeam)
		teams.POST("/dislikes", controllers.CreateRelationDislikesTeam)
	}

	places := router.Group("/places")
	{
		places.GET("/", controllers.GetPlaces)
		places.POST("/", controllers.NewPlace)
		places.POST("/visited", controllers.CreateRelationVisited)
		places.POST("/dislikes", controllers.CreateRelationDislikesPlace)
	}

	songs := router.Group("/songs")
	{
		songs.GET("/", controllers.GetSongs)
		songs.POST("/", controllers.NewSong)
		songs.POST("/likes", controllers.CreateRelationLikesSong)
		songs.POST("/dislikes", controllers.CreateRelationDislikesSong)
		songs.POST("/favorite", controllers.CreateRelationFavoriteSong)

		songs.PUT("/remembers", controllers.SetSongNewProperty)
		songs.POST("/remembers/remove", controllers.DeleteSongRememberProperty)

		songs.PUT("/music-player", controllers.SetPreferredMusicPlayer)
		songs.DELETE("/music-player/:username", controllers.DeletePreferredMusicPlayer)
	}

}
