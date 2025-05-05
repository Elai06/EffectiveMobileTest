package api

import (
	"effectiveMobile/env"
	"effectiveMobile/internal"
	"effectiveMobile/internal/db"
	"effectiveMobile/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Handlers struct {
	repos    *db.Repository
	external *External
}

func NewHandlers(repos *db.Repository, external *External) *Handlers {
	return &Handlers{repos: repos, external: external}
}

func (h *Handlers) Start(cfg env.Config) error {
	r := gin.Default()

	r.POST("/create", h.create)
	r.GET("/get", h.get)
	r.DELETE("/delete", h.delete)
	r.PUT("/update", h.update)

	srv := http.Server{
		Addr:         cfg.HttpPort,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout * time.Second,
		WriteTimeout: cfg.WriteTimeout * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		return fmt.Errorf("http server start error: %s", err.Error())
	}

	utils.InfoLog("http server start success")

	return nil
}

// @Summary Create a new person
// @Description Creates a new person and enriches data via external APIs
// @Tags people
// @Accept json
// @Produce json
// @Param name query string true "Name of the person"
// @Param surname query string true "Surname of the person"
// @Success 200 {object} internal.Person
// @Failure 500 {object} map[string]string
// @Router /create [post]
func (h *Handlers) create(c *gin.Context) {
	name := c.Query("name")
	surname := c.Query("surname")

	utils.InfoLog("the beginning of human creation")

	age, err := h.external.GetAge(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	gender, err := h.external.GetGender(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	nationality, err := h.external.GetNationality(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	person := internal.Person{
		Name:        name,
		Surname:     surname,
		Age:         age,
		Gender:      gender,
		Nationality: nationality}
	err = h.repos.CreatePerson(person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	utils.DebugLog("the creation request was completed successfully", person)

	c.JSON(http.StatusOK, person)
}

// @Summary Update an existing person
// @Description Updates a person by ID using enriched data from external APIs
// @Tags people
// @Accept json
// @Produce json
// @Param id query integer true "ID of the person to update"
// @Param name query string true "New name of the person"
// @Param surname query string true "New surname of the person"
// @Success 200 {object} internal.Person
// @Failure 500 {object} map[string]string
// @Router /update [put]
func (h *Handlers) update(c *gin.Context) {
	name := c.Query("name")
	surname := c.Query("surname")
	id := c.Query("id")

	utils.InfoLog("the beginning of human update")

	age, err := h.external.GetAge(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	gender, err := h.external.GetGender(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	nationality, err := h.external.GetNationality(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	person := internal.Person{
		Name:        name,
		Surname:     surname,
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
		ID:          id}
	err = h.repos.UpdatePerson(person)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("error updating person: %s", err).Error())
		return
	}

	utils.DebugLog("the update was completed successfully", person)

	c.JSON(http.StatusOK, person)
}

// @Summary Delete a person
// @Description Deletes a person by ID
// @Tags people
// @Accept json
// @Produce json
// @Param id query integer true "ID of the person to delete"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /delete [delete]
func (h *Handlers) delete(c *gin.Context) {
	idStr := c.Query("id")
	utils.InfoLog("the beginning of human delete")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("invalid id: %s", idStr).Error())
		return
	}

	err = h.repos.DeletePerson(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("error deleting person: %s", err).Error())
		return
	}

	utils.DebugLog("the delete was completed successfully", id)

	c.JSON(http.StatusOK, "the delete was completed successfully")
}

// @Summary Get all people with filters
// @Description Retrieves list of people with optional filtering and pagination
// @Tags people
// @Accept json
// @Produce json
// @Param name query string false "Filter by name"
// @Param surname query string false "Filter by surname"
// @Param age query int false "Exact age filter"
// @Param age_min query int false "Minimum age"
// @Param age_max query int false "Maximum age"
// @Param gender query string false "Filter by gender"
// @Param nationality query string false "Filter by nationality"
// @Param limit query int false "Limit per page" default(10)
// @Param offset query int false "Page offset" default(0)
// @Success 200 {object} map[string][]internal.Person
// @Failure 500 {object} map[string]string
// @Router /get [get]
func (h *Handlers) get(c *gin.Context) {

	utils.InfoLog("the beginning of human get")

	name := c.DefaultQuery("name", "")
	surname := c.DefaultQuery("surname", "")
	age := parseInt(c.DefaultQuery("age", "0"))
	ageMin := parseInt(c.DefaultQuery("age_min", "0"))
	ageMax := parseInt(c.DefaultQuery("age_max", "0"))
	gender := c.DefaultQuery("gender", "")
	nationality := c.DefaultQuery("nationality", "")
	limit := parseInt(c.DefaultQuery("limit", "10"))
	offset := parseInt(c.DefaultQuery("offset", "0"))
	filter := internal.PersonFilter{
		Person: internal.Person{Name: name, Surname: surname, Age: age, Gender: gender, Nationality: nationality},
		Limit:  limit,
		Offset: offset,
		AgeMax: ageMax,
		AgeMin: ageMin,
	}

	people, err := h.repos.GetPeople(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Errorf("get people error: %s", err).Error())
		return
	}

	utils.DebugLog("the get was completed successfully", people)

	c.JSON(http.StatusOK, gin.H{"people": people})
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
