package controllers

import (
	"backend/configs"
	"backend/models"
	"backend/responses"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"time"
)

// GetPlaces Obtiene todos los lugares
// @Summary Obtiene todos los lugares
// @Description Obtiene todos los lugares registrados en la base de datos
// @Tags Lugares
// @Accept json
// @Produce json
// @Success 200 {object} responses.PlacesResponse "Lugares obtenidos exitosamente"
// @Failure 500 {object} responses.ErrorResponse "Error al procesar la solicitud"
// @Router /places/ [get]
func GetPlaces(c *gin.Context) {
	session := configs.DB.NewSession(c, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al cerrar la sesión",
				Error:   err.Error(),
			})
		}
	}(session, c)

	// Consulta para obtener todos los lugares
	r, err := session.Run(
		c,
		"MATCH (l:Lugar) RETURN l",
		nil,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error al procesar la solicitud",
			Error:   err.Error(),
		})
		return
	}

	var places []models.Lugar
	for r.Next(c) {
		vals := r.Record().Values[0].(neo4j.Node).Props

		place := models.Lugar{
			Nombre:       vals["Nombre"].(string),
			Departamento: vals["Departamento"].(string),
			Tipo:         vals["Tipo"].(string),
			Direccion:    vals["Dirección"].(string),
			Foto:         vals["Foto"].(string),
		}

		places = append(places, place)
	}

	c.JSON(http.StatusOK, responses.PlacesResponse{
		Status:  http.StatusOK,
		Message: "Lugares obtenidos exitosamente",
		Places:  places,
	})
}

// NewPlace Crea un nuevo lugar
// @Summary Crea un nuevo lugar
// @Description Crea un nuevo lugar en la base de datos
// @Tags Lugares
// @Accept json
// @Produce json
// @Param place body models.Lugar true "Datos del lugar a crear"
// @Success 201 {object} responses.StandardResponse "Lugar creado exitosamente"
// @Failure 400 {object} responses.ErrorResponse "Error al procesar la solicitud"
// @Failure 500 {object} responses.ErrorResponse "Error al procesar la solicitud"
// @Router /places/ [post]
func NewPlace(c *gin.Context) {
	session := configs.DB.NewSession(c, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al cerrar la sesión",
				Error:   err.Error(),
			})
		}
	}(session, c)

	// Bind de la estructura de lugar
	var place models.Lugar
	if err := c.BindJSON(&place); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Error al procesar la solicitud",
			Error:   err.Error(),
		})
		return
	}

	// Consulta para crear un nuevo lugar
	_, err := session.Run(
		c,
		"CREATE (l:Lugar {Nombre: $Nombre, Departamento: $Departamento, Tipo: $Tipo, Dirección: $Dirección, Foto: $Foto}) RETURN l",
		map[string]interface{}{
			"Nombre":       place.Nombre,
			"Departamento": place.Departamento,
			"Tipo":         place.Tipo,
			"Dirección":    place.Direccion,
			"Foto":         place.Foto,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error al procesar la solicitud",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, responses.StandardResponse{
		Status:  http.StatusCreated,
		Message: "Lugar creado exitosamente",
		Data:    map[string]interface{}{"place": place},
	})
}

// CreateRelationVisited Crea una relación de visita de un usuario a un lugar
// @Summary Crea una relación de visita de un usuario a un lugar
// @Description Crea una relación de (Persona)-[VISITA]->(Lugar)
// @Tags Lugares
// @Accept json
// @Produce json
// @Param relation body models.RelationVisitaLugar true "Relación de visita de un lugar"
// @Success 200 {object} responses.StandardResponse "Relación creada exitosamente"
// @Failure 400 {object} responses.ErrorResponse "Error al procesar la solicitud"
// @Failure 500 {object} responses.ErrorResponse "Error al procesar la solicitud"
// @Router /places/visited [post]
func CreateRelationVisited(c *gin.Context) {
	session := configs.DB.NewSession(c, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al cerrar la sesión",
				Error:   err.Error(),
			})
		}
	}(session, c)

	var relation models.RelationVisitaLugar
	err := c.BindJSON(&relation)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Error al procesar la solicitud",
			Error:   err.Error(),
		})
		return
	}

	f, err := time.Parse(time.DateOnly, relation.Cuando)

	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Error al procesar la fecha",
			Error:   err.Error(),
		})
		return
	}

	_, err = session.Run(
		c,
		"MATCH (p:Persona {Usuario: $usuario}), (l:Lugar {Nombre: $lugar}) CREATE (p)-[r:VISITA {Cuando: $cuando, Rating: $rating, Categoria: $categoria}]->(l) RETURN r",
		map[string]interface{}{
			"usuario":   relation.Usuario,
			"lugar":     relation.Lugar,
			"cuando":    f,
			"rating":    relation.Rating,
			"categoria": relation.Categoria,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error al crear la relación",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, responses.StandardResponse{
		Status:  http.StatusOK,
		Message: "Relación creada exitosamente",
		Data:    nil,
	})
}

// CreateRelationDislikesPlace Crea una relación de no le gusta de un usuario a un lugar
// @Summary Crea una relación de no le gusta de un usuario a un lugar
// @Description Crea una relación de (Persona)-[NO_LE_GUSTA]->(Lugar)
// @Tags Lugares
// @Accept json
// @Produce json
// @Param relation body models.RelationNoLeGustaLugar true "Relación de no le gusta de un lugar"
// @Success 200 {object} responses.StandardResponse "Relación creada exitosamente"
// @Failure 400 {object} responses.ErrorResponse "Error al procesar la solicitud"
// @Failure 500 {object} responses.ErrorResponse "Error al procesar la solicitud"
// @Router /places/dislikes [post]
func CreateRelationDislikesPlace(c *gin.Context) {
	session := configs.DB.NewSession(c, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})

	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
				Status:  http.StatusInternalServerError,
				Message: "Error al cerrar la sesión",
				Error:   err.Error(),
			})
		}
	}(session, c)

	var relation models.RelationNoLeGustaLugar
	err := c.BindJSON(&relation)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Error al procesar la solicitud",
			Error:   err.Error(),
		})
		return
	}

	f, err := time.Parse(time.DateOnly, relation.Cuando)

	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Error al procesar la fecha",
			Error:   err.Error(),
		})
		return
	}

	_, err = session.Run(
		c,
		"MATCH (p:Persona {Usuario: $usuario}), (l:Lugar {Nombre: $lugar}) CREATE (p)-[r:NO_LE_GUSTA {Cuando: $cuando, Rating: $rating, Categoria: $categoria}]->(l) RETURN r",
		map[string]interface{}{
			"usuario":   relation.Usuario,
			"lugar":     relation.Lugar,
			"cuando":    f,
			"rating":    relation.Rating,
			"categoria": relation.Categoria,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "Error al crear la relación",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, responses.StandardResponse{
		Status:  http.StatusOK,
		Message: "Relación creada exitosamente",
		Data:    nil,
	})
}
