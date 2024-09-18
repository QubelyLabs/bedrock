package controller

import (
	"fmt"
	"log"
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
)

type Controller[E any] struct {
	*BaseController
	repository contract.Repository[E]
	name       string
	plural     string
	searchable []string
	unique     func(*E) (any, []any)
	morph      func(*E)
	hooks      map[string]func(*E, *gin.Context)
}

func NewController[E any](
	repository contract.Repository[E],
	name,
	plural string,
	searchable []string,
	unique func(*E) (any, []any),
	morph func(*E),
	hooks map[string]func(*E, *gin.Context),
) *Controller[E] {
	return &Controller[E]{&BaseController{}, repository, name, plural, searchable, unique, morph, hooks}
}

// TODO
// implement @Query('update') shouldUpdate: string to udpate if record already exist
func (ctrl *Controller[E]) UpsertOne(c *gin.Context) {
	entity := new(E)
	if data, ok := ctrl.Validate(c, entity); !ok {
		ctrl.ErrorWithData(c, "Invalid request, check and try again", data)
		return
	}

	if ctrl.morph != nil {
		ctrl.morph(entity)
	}

	if hook, ok := ctrl.hooks[BeforeCreate]; ok {
		hook(entity, c)
	}

	err := ctrl.repository.UpsertOne(c, entity)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to save %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterCreate]; ok {
		hook(entity, c)
	}

	ctrl.Success(c, fmt.Sprintf("%v record saved successfully", ctrl.name), entity)
}

func (ctrl *Controller[E]) UpsertMany(c *gin.Context) {
	entities := []E{}
	if data, ok := ctrl.Validate(c, &entities); !ok {
		ctrl.ErrorWithData(c, "Invalid request, check and try again", data)
		return
	}

	for _, entity := range entities {
		ctrl.morph(&entity)
	}

	if hook, ok := ctrl.hooks[BeforeCreate]; ok {
		for _, entity := range entities {
			hook(&entity, c)
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
			hook(&entity, c)
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

	fmt.Println(entity)

	if ctrl.morph != nil {
		ctrl.morph(entity)
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
		hook(entity, c)
	}

	err := ctrl.repository.CreateOne(c, entity)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to save %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterCreate]; ok {
		hook(entity, c)
	}

	ctrl.Success(c, fmt.Sprintf("%v record saved successfully", ctrl.name), entity)
}

func (ctrl *Controller[E]) CreateMany(c *gin.Context) {
	entities := []E{}
	if data, ok := ctrl.Validate(c, &entities); !ok {
		ctrl.ErrorWithData(c, "Invalid request, check and try again", data)
		return
	}

	for _, entity := range entities {
		ctrl.morph(&entity)
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
			hook(&entity, c)
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
			hook(&entity, c)
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

	existingEntity, err := ctrl.repository.FindOne(c, id)
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

	if ctrl.morph != nil {
		ctrl.morph(&existingEntity) // TODO merge with entity
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
		hook(entity, c)
	}

	err = ctrl.repository.UpdateOne(c, id, entity)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to update %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterUpdate]; ok {
		hook(entity, c)
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

	ids := strings.Split(id, ",")
	var entities []E
	for _, _ = range ids {
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

		ctrl.morph(&existingEntity) // TODO merge with entity

		entities = append(entities, existingEntity)
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
			hook(&entity, c)
		}
	}

	err := ctrl.repository.UpdateMany(c, entity, "id IN ?", strings.Split(id, ","))
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to update %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterUpdate]; ok {
		for _, entity := range entities {
			hook(&entity, c)
		}
	}

	ctrl.Success(c, fmt.Sprintf("%v record updated successfully", ctrl.name), entity)
}

func (ctrl *Controller[E]) FindOne(c *gin.Context) {
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

	ctrl.Success(c, fmt.Sprintf("%v record retrieved successfully", ctrl.name), entity)
}

// async findAll(@Query() queries: any, @GetContext() context: any) {
// 	const [selector, page, perPage, columns, relations] = makeFilter(
// 	  queries,
// 	  this.searchable,
// 	);
// 	const _page = page < 1 ? 1 : page;
// 	const _nextPage = _page + 1;
// 	const _prevPage = _page - 1;
// 	const _perPage = perPage;
// 	const _filter = {
// 	  select: columns,
// 	  ...(perPage ? { take: perPage } : {}),
// 	  ...(page && perPage ? { skip: (page - 1) * perPage } : {}),
// 	  where: selector.length > 0 ? selector : undefined,
// 	  relations,
// 	};

//		const total = await this.service.withContext(context).count(_filter);
//		const entities = await this.service.withContext(context).find(_filter);
//		return success(
//		  entities,
//		  `${titleCase(this.name)}`,
//		  `${titleCase(this.name)} list`,
//		  {
//			current_page: _page,
//			next_page: _nextPage > total ? total : _nextPage,
//			prev_page: _prevPage < 1 ? null : _prevPage,
//			per_page: _perPage,
//			total,
//		  },
//		);
//	  }
func (ctrl *Controller[E]) FindMany(c *gin.Context) {
	pageStr := c.Query("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	perPageStr := c.Query("perPage")
	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage <= 0 {
		perPage = 12
	}

	offset := (page - 1) * perPage
	entities, err := ctrl.repository.FindManyWithLimit(c, perPage, offset, nil)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to retrieve %v record, try again in a bit", ctrl.name), 500)
		return
	}

	ctrl.Success(c, fmt.Sprintf("%v records retrieved successfully", ctrl.name), entities)
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

	if hook, ok := ctrl.hooks[BeforeDelete]; ok {
		hook(&entity, c)
	}

	err = ctrl.repository.DeleteOne(c, id)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to remove %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterDelete]; ok {
		hook(&entity, c)
	}

	ctrl.Success(c, fmt.Sprintf("%v record removed successfully", ctrl.name), nil)
}

func (ctrl *Controller[E]) DeleteMany(c *gin.Context) {
	id := c.Param("id")

	ids := strings.Split(id, ",")
	var entities []E
	for _, _ = range ids {
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

	if hook, ok := ctrl.hooks[BeforeDelete]; ok {
		for _, entity := range entities {
			hook(&entity, c)
		}
	}

	err := ctrl.repository.DeleteMany(c, "id IN ?", ids)
	if err != nil {
		log.Println(err)
		ctrl.ErrorWithCode(c, fmt.Sprintf("Unable to remove %v record, try again in a bit", ctrl.name), 500)
		return
	}

	if hook, ok := ctrl.hooks[AfterDelete]; ok {
		for _, entity := range entities {
			hook(&entity, c)
		}
	}

	ctrl.Success(c, fmt.Sprintf("%v records removed successfully", ctrl.name), nil)
}
