package server

import (
	"githib.com/s4bb4t/leadgen/internal/lib/models"
	"githib.com/s4bb4t/leadgen/internal/storage"

	_ "githib.com/s4bb4t/leadgen/docs"
	"github.com/gin-gonic/gin"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"strconv"
)

type API struct {
	Log  *slog.Logger
	Repo storage.RepositoryI
}

func Run(logger *slog.Logger, storage storage.RepositoryI) {
	api := API{
		Log:  logger,
		Repo: storage,
	}

	r := gin.Default()

	r.GET("/swagger/*any", gin.WrapH(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	)))

	v1 := r.Group("/api/v1")
	{
		v1.POST("/buildings", api.CreateBuilding)

		v1.GET("/buildings/:title", api.GetBuildingByTitle)

		v1.GET("/buildings", api.GetAllBuildings)
	}

	api.Log.Info("Starting Server")
	if err := r.Run("localhost:8080"); err != nil {
		logger.Error("Unable to start server", err)
	}
}

// CreateBuilding godoc
// @Summary Create a new building
// @Description Creates a new building with the provided data
// @Tags buildings
// @Accept json
// @Produce json
// @Param building body models.Building true "Building data"
// @Success 201 {object} models.Building
// @Failure 400 {object} string
// @Router /api/v1/buildings [post]
func (api *API) CreateBuilding(c *gin.Context) {
	var building models.Building
	if err := c.ShouldBindJSON(&building); err != nil {
		api.Log.Info(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBuilding, err := api.Repo.Save(c, building)
	if err != nil {
		api.Log.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	api.Log.Debug("success", slog.Any("building", createdBuilding))
	c.JSON(http.StatusCreated, createdBuilding)
}

// GetBuildingByTitle godoc
// @Summary Get a building by title
// @Description Returns information about a building by its title
// @Tags buildings
// @Accept json
// @Produce json
// @Param title path string true "Building title"
// @Success 200 {object} models.Building
// @Failure 404 {object} string
// @Router /api/v1/buildings/{title} [get]
func (api *API) GetBuildingByTitle(c *gin.Context) {
	title := c.Param("title")

	building, err := api.Repo.Building(c, title)
	if err != nil {
		api.Log.Info(err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Building not found"})
		return
	}

	api.Log.Debug("success", slog.Any("building", building))
	c.JSON(http.StatusOK, building)
}

// GetAllBuildings godoc
// @Summary Get all buildings with filtering
// @Description Returns a list of all buildings with the ability to filter by city, year, and number of floors
// @Tags buildings
// @Accept json
// @Produce json
// @Param city query string false "City"
// @Param year query int false "Year"
// @Param floors query int false "Number of floors"
// @Param limit query int false "Limit of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {object} models.Buildings
// @Failure 500 {object} string
// @Router /api/v1/buildings [get]
func (api *API) GetAllBuildings(c *gin.Context) {
	query := models.Query{
		City: c.DefaultQuery("city", ""),
	}

	if year := c.DefaultQuery("year", ""); year != "" {
		yearInt, err := strconv.Atoi(year)
		if err != nil {
			api.Log.Info(err.Error())
		}
		query.Year = yearInt
	}

	if floors := c.DefaultQuery("floors", ""); floors != "" {
		floorsInt, err := strconv.Atoi(floors)
		if err != nil {
			api.Log.Info(err.Error())
		}
		query.Floors = floorsInt
	}

	if limit := c.DefaultQuery("limit", "10"); limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			api.Log.Info(err.Error())
		}
		query.Limit = limitInt
	}

	if offset := c.DefaultQuery("offset", "0"); offset != "" {
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			api.Log.Info(err.Error())
		}
		query.Offset = offsetInt
	}

	buildings, err := api.Repo.Buildings(c, query)
	if err != nil {
		api.Log.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	api.Log.Debug("success", slog.Any("meta", buildings.Meta))
	c.JSON(http.StatusOK, buildings)
}
