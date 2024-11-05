package controller

import (
	"fmt"
	"log"
	"math"
	"net/url"
	"strconv"
	"strings"

	"github.com/QubelyLabs/bedrock/pkg/contract"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	BeforeCreate = "beforeCreate"
	AfterCreate  = "AfterCreate"
	BeforeUpdate = "beforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "beforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeRead   = "beforeRead"
	AfterRead    = "AfterRead"
)

type Controller[E any] struct {
	*BaseController
	repository contract.Repository[E]
	name       string
	plural     string
	searchable []string
	unique     func(*E) (any, []any)
	morphs     map[string]func(*E, *gin.Context)
	hooks      map[string]func(*E, *gin.Context) error
}

func NewController[E any](
	repository contract.Repository[E],
	name,
	plural string,
	searchable []string,
	unique func(*E) (any, []any),
	morphs map[string]func(*E, *gin.Context),
	hooks map[string]func(*E, *gin.Context) error,
) *Controller[E] {
	return &Controller[E]{&BaseController{}, repository, name, plural, searchable, unique, morphs, hooks}
}

func (ctrl *Controller[E]) UpsertOne(c *gin.Context) {
	entity := new(E)
	if data, ok := ctrl.Validate(c, entity); !ok {
		ctrl.ErrorWithData(c, "Invalid request, check and try again", data)
		return
	}

	if morph, ok := ctrl.morphs[BeforeCreate]; ok {
		morph(entity, c)
	}

	if hook, ok := ctrl.hooks[BeforeCreate]; ok {
		err := hook(entity, c)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, err.Error())
			return
		}
	}

	err := ctrl.repository.UpsertOne(c, entity)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to save %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterCreate]; ok {
		err := hook(entity, c)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, err.Error())
			return
		}
	}

	if morph, ok := ctrl.morphs[AfterCreate]; ok {
		morph(entity, c)
	}

	ctrl.Success(c, fmt.Sprintf("%v record saved successfully", ctrl.name), entity)
}

func (ctrl *Controller[E]) UpsertMany(c *gin.Context) {
	entities := []E{}
	if data, ok := ctrl.Validate(c, &entities); !ok {
		ctrl.ErrorWithData(c, "Invalid request, check and try again", data)
		return
	}

	if morph, ok := ctrl.morphs[BeforeCreate]; ok {
		for i := range entities {
			morph(&entities[i], c)
		}
	}

	if hook, ok := ctrl.hooks[BeforeCreate]; ok {
		for _, entity := range entities {
			err := hook(&entity, c)
			if err != nil {
				log.Println(err)
				ctrl.Error(c, err.Error())
				return
			}
		}
	}

	err := ctrl.repository.UpsertMany(c, entities...)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to save %v records, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterCreate]; ok {
		for _, entity := range entities {
			err := hook(&entity, c)
			if err != nil {
				log.Println(err)
				ctrl.Error(c, err.Error())
				return
			}
		}
	}

	if morph, ok := ctrl.morphs[AfterCreate]; ok {
		for i := range entities {
			morph(&entities[i], c)
		}
	}

	ctrl.Success(c, fmt.Sprintf("%v records saved successfully", ctrl.name), entities)
}

func (ctrl *Controller[E]) CreateOne(c *gin.Context) {
	entity := new(E)
	if data, ok := ctrl.Validate(c, entity); !ok {
		ctrl.ErrorWithData(c, "Invalid request, check and try again", data)
		return
	}

	if morph, ok := ctrl.morphs[BeforeCreate]; ok {
		morph(entity, c)
	}

	if ctrl.unique != nil {
		query, args := ctrl.unique(entity)
		existing, err := ctrl.repository.Count(c, query, args...)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, "Something went wrong, check and try again")
			return
		}

		if existing > 0 {
			log.Println(err)
			ctrl.ErrorWithData(c, fmt.Sprintf("A similar %v record exist, check and try again", ctrl.name), entity)
			return
		}
	}

	if hook, ok := ctrl.hooks[BeforeCreate]; ok {
		err := hook(entity, c)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, err.Error())
			return
		}
	}

	err := ctrl.repository.CreateOne(c, entity)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to save %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterCreate]; ok {
		err := hook(entity, c)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, err.Error())
			return
		}
	}

	if morph, ok := ctrl.morphs[AfterCreate]; ok {
		morph(entity, c)
	}

	ctrl.Success(c, fmt.Sprintf("%v record saved successfully", ctrl.name), entity)
}

func (ctrl *Controller[E]) CreateMany(c *gin.Context) {
	entities := []E{}
	if data, ok := ctrl.Validate(c, &entities); !ok {
		ctrl.ErrorWithData(c, "Invalid request, check and try again", data)
		return
	}

	if morph, ok := ctrl.morphs[BeforeCreate]; ok {
		for i := range entities {
			morph(&entities[i], c)
		}
	}

	if ctrl.unique != nil {
		for _, entity := range entities {
			query, args := ctrl.unique(&entity)
			existing, err := ctrl.repository.Count(c, query, args...)
			if err != nil {
				log.Println(err)
				ctrl.Error(c, "Something went wrong, check and try again")
				return
			}

			if existing > 0 {
				log.Println(err)
				ctrl.ErrorWithData(c, fmt.Sprintf("A similar %v record exist, check and try again", ctrl.name), entity)
				return
			}
		}
	}

	if hook, ok := ctrl.hooks[BeforeCreate]; ok {
		for _, entity := range entities {
			err := hook(&entity, c)
			if err != nil {
				log.Println(err)
				ctrl.Error(c, err.Error())
				return
			}
		}
	}

	err := ctrl.repository.CreateMany(c, entities...)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to save %v records, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterCreate]; ok {
		for _, entity := range entities {
			err := hook(&entity, c)
			if err != nil {
				log.Println(err)
				ctrl.Error(c, err.Error())
				return
			}
		}
	}

	if morph, ok := ctrl.morphs[AfterCreate]; ok {
		for i := range entities {
			morph(&entities[i], c)
		}
	}

	ctrl.Success(c, fmt.Sprintf("%v records saved successfully", ctrl.name), entities)
}

func (ctrl *Controller[E]) UpdateOne(c *gin.Context) {
	entity := new(E)
	id := c.Param("id")
	if data, ok := ctrl.Validate(c, entity); !ok {
		ctrl.ErrorWithData(c, "Invalid request, check and try again", data)
		return
	}

	_, err := ctrl.repository.FindOne(c, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println(err)
			ctrl.ErrorWithCode(c, "Invalid request, record not found", 404)
			return
		}

		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to retrieve %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if morph, ok := ctrl.morphs[BeforeUpdate]; ok {
		morph(entity, c)
	}

	if ctrl.unique != nil {
		query, args := ctrl.unique(entity)
		var existing int64
		err := ctrl.repository.SQL(c).WithContext(c).Where(query, args...).Where("id != ?", id).Model(entity).Count(&existing).Error
		if err != nil {
			log.Println(err)
			ctrl.Error(c, "Something went wrong, check and try again")
			return
		}

		if existing > 0 {
			log.Println(err)
			ctrl.ErrorWithData(c, fmt.Sprintf("A similar %v record exist, check and try again", ctrl.name), entity)
			return
		}
	}

	if hook, ok := ctrl.hooks[BeforeUpdate]; ok {
		err := hook(entity, c)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, err.Error())
			return
		}
	}

	err = ctrl.repository.UpdateOne(c, id, entity)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to update %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterUpdate]; ok {
		err := hook(entity, c)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, err.Error())
			return
		}
	}

	if morph, ok := ctrl.morphs[AfterUpdate]; ok {
		morph(entity, c)
	}

	ctrl.Success(c, fmt.Sprintf("%v record updated successfully", ctrl.name), entity)
}

func (ctrl *Controller[E]) UpdateMany(c *gin.Context) {
	entity := new(E)
	id := c.Param("id")
	if data, ok := ctrl.Validate(c, entity); !ok {
		ctrl.ErrorWithData(c, "Invalid request, check and try again", data)
		return
	}

	ids := strings.Split(id, "|")
	var entities []E
	for range ids {
		existingEntity, err := ctrl.repository.FindOne(c, id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Println(err)
				ctrl.ErrorWithDataAndCode(c, "Invalid request, record not found", gin.H{"id": id}, 404)
				return
			}

			log.Println(err)
			ctrl.ErrorWithDataAndCode(c, fmt.Sprintf("Unable to retrieve %v record, try again in a bit", ctrl.name), gin.H{"id": id}, 500)
			return
		}

		entities = append(entities, existingEntity)
	}

	if morph, ok := ctrl.morphs[BeforeUpdate]; ok {
		for i := range entities {
			morph(&entities[i], c)
		}
	}

	if ctrl.unique != nil {
		for _, id := range ids {
			query, args := ctrl.unique(entity)
			var existing int64
			err := ctrl.repository.SQL(c).WithContext(c).Where(query, args...).Where("id != ?", id).Model(entity).Count(&existing).Error
			if err != nil {
				log.Println(err)
				ctrl.Error(c, "Something went wrong, check and try again")
				return
			}

			if existing > 0 {
				log.Println(err)
				ctrl.ErrorWithData(c, fmt.Sprintf("A similar %v record exist, check and try again", ctrl.name), entity)
				return
			}
		}
	}

	if hook, ok := ctrl.hooks[BeforeUpdate]; ok {
		for _, entity := range entities {
			err := hook(&entity, c)
			if err != nil {
				log.Println(err)
				ctrl.Error(c, err.Error())
				return
			}
		}
	}

	err := ctrl.repository.UpdateMany(c, entity, "id IN ?", ids)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to update %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterUpdate]; ok {
		for _, entity := range entities {
			err := hook(&entity, c)
			if err != nil {
				log.Println(err)
				ctrl.Error(c, err.Error())
				return
			}
		}
	}

	if morph, ok := ctrl.morphs[AfterUpdate]; ok {
		for i := range entities {
			morph(&entities[i], c)
		}
	}

	ctrl.Success(c, fmt.Sprintf("%v record updated successfully", ctrl.name), entity)
}

func (ctrl *Controller[E]) FindOne(c *gin.Context) {
	id := c.Param("id")

	if morph, ok := ctrl.morphs[BeforeRead]; ok {
		morph(new(E), c)
	}

	if hook, ok := ctrl.hooks[BeforeRead]; ok {
		err := hook(new(E), c)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, err.Error())
			return
		}
	}

	entity, err := ctrl.repository.FindOne(c, id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println(err)
			ctrl.ErrorWithCode(c, "Invalid request, record not found", 404)
			return
		}

		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to retrieve %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterRead]; ok {
		err := hook(&entity, c)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, err.Error())
			return
		}
	}

	if morph, ok := ctrl.morphs[AfterRead]; ok {
		morph(&entity, c)
	}

	ctrl.Success(c, fmt.Sprintf("%v record retrieved successfully", ctrl.name), entity)
}

func (ctrl *Controller[E]) FindMany(c *gin.Context) {
	log.Println(c.Query("firstname"), c.Request.Context(), c.Request.URL.Query(), c.Request.URL.Query(), "per_page")
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	perPageStr := c.Query("per_page")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage <= 0 {
		perPage = 12
	}

	offset := (page - 1) * perPage

	if morph, ok := ctrl.morphs[BeforeRead]; ok {
		morph(new(E), c)
	}

	if hook, ok := ctrl.hooks[BeforeRead]; ok {
		err := hook(new(E), c)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, err.Error())
			return
		}
	}

	query, args := ctrl.buildQuery(c.Request.URL.Query())

	log.Println(query, args)

	total, err := ctrl.repository.Count(c, query, args...)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to retrieve %v record, try again in a bit", ctrl.name), 500)
		return
	}

	lastPage := math.Ceil(float64(total) / float64(perPage))
	nextPage := page + 1
	if nextPage > int(lastPage) {
		nextPage = 0
	}

	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 0
	}

	entities, err := ctrl.repository.FindManyWithLimit(c, perPage, offset, query, args...)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to retrieve %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterRead]; ok {
		for _, entity := range entities {
			err := hook(&entity, c)
			if err != nil {
				log.Println(err)
				ctrl.Error(c, err.Error())
				return
			}
		}
	}

	if morph, ok := ctrl.morphs[AfterRead]; ok {
		for i := range entities {
			morph(&entities[i], c)
		}
	}

	meta := map[string]any{
		"page":     page,
		"per_page": perPage,
		"total":    total,
		"prev":     prevPage,
		"next":     nextPage,
	}

	ctrl.SuccessWithMeta(c, fmt.Sprintf("%v records retrieved successfully", ctrl.name), entities, meta)
}

func (ctrl *Controller[E]) DeleteOne(c *gin.Context) {
	id := c.Param("id")

	entity, err := ctrl.repository.FindOne(c, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println(err)
			ctrl.ErrorWithCode(c, "Invalid request, record not found", 404)
			return
		}

		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to retrieve %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if morph, ok := ctrl.morphs[BeforeDelete]; ok {
		morph(&entity, c)
	}

	if hook, ok := ctrl.hooks[BeforeDelete]; ok {
		err := hook(&entity, c)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, err.Error())
			return
		}
	}

	err = ctrl.repository.DeleteOne(c, id)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to remove %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterDelete]; ok {
		err := hook(&entity, c)
		if err != nil {
			log.Println(err)
			ctrl.Error(c, err.Error())
			return
		}
	}

	if morph, ok := ctrl.morphs[AfterDelete]; ok {
		morph(&entity, c)
	}

	ctrl.Success(c, fmt.Sprintf("%v record removed successfully", ctrl.name), nil)
}

func (ctrl *Controller[E]) DeleteMany(c *gin.Context) {
	ids := new([]string)
	if data, ok := ctrl.Validate(c, ids); !ok {
		ctrl.ErrorWithData(c, "Invalid request, check and try again", data)
		return
	}

	var entities []E
	for _, id := range *ids {
		entity, err := ctrl.repository.FindOne(c, id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Println(err)
				ctrl.ErrorWithDataAndCode(c, "Invalid request, record not found", gin.H{"id": id}, 404)
				return
			}

			log.Println(err)
			ctrl.ErrorWithDataAndCode(c, fmt.Sprintf("Unable to retrieve %v record, try again in a bit", ctrl.name), gin.H{"id": id}, 500)
			return
		}

		entities = append(entities, entity)
	}

	if morph, ok := ctrl.morphs[BeforeDelete]; ok {
		for i := range entities {
			morph(&entities[i], c)
		}
	}

	if hook, ok := ctrl.hooks[BeforeDelete]; ok {
		for _, entity := range entities {
			err := hook(&entity, c)
			if err != nil {
				log.Println(err)
				ctrl.Error(c, err.Error())
				return
			}
		}
	}

	err := ctrl.repository.DeleteMany(c, "id IN ?", *ids)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to remove %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterDelete]; ok {
		for _, entity := range entities {
			err := hook(&entity, c)
			if err != nil {
				log.Println(err)
				ctrl.Error(c, err.Error())
				return
			}
		}
	}

	if morph, ok := ctrl.morphs[AfterDelete]; ok {
		for i := range entities {
			morph(&entities[i], c)
		}
	}

	ctrl.Success(c, fmt.Sprintf("%v records removed successfully", ctrl.name), nil)
}

func (ctrl *Controller[E]) buildQuery(queryParams url.Values) (string, []interface{}) {
	var (
		queryParts []string
		args       []interface{}
	)

	// Keys to exclude from the query
	excludeKeys := []string{"page", "per_page"}

	// Determine joiner ("and" or "or"), default to "and"
	joiner := strings.ToUpper(queryParams.Get("joiner"))
	if joiner != "OR" {
		joiner = "AND"
	}

	// Build the query, excluding specified keys
	for key, values := range queryParams {
		// Skip excluded keys and "joiner" key
		if ctrl.contains(excludeKeys, key) || key == "joiner" {
			continue
		}

		if len(values) > 0 {
			// Add condition for each parameter
			queryParts = append(queryParts, fmt.Sprintf("%s = ?", key))
			args = append(args, values[0]) // Add only the first value for simplicity
		}
	}

	query := strings.Join(queryParts, " "+joiner+" ")
	return query, args
}

func (ctrl *Controller[E]) contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
