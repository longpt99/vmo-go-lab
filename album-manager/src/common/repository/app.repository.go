package repository

import (
	"album-manager/src/common/models"
	"album-manager/src/configs/database"
	"context"
	"fmt"
	"math"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	InsertOne(params *T) (*string, error)
	List(params *FindParams) (*ResponseData[T], error)
	DetailByID(id string) (*T, error)
	Delete(id string) error
	DetailByConditions(params *QueryParams) (*T, error)
	UpdateByConditions(params *UpdateParams) error
	CountByConditions(params *QueryParams) (*int64, error)
}

type QueryParams struct {
	Args      []interface{}
	TableName string
	Where     string
	OrderBy   string
	Columns   []string
	Joins     []string
	Limit     int
	Offset    int
}

type FindParams struct {
	Args   []interface{}
	Where  string
	Select []string
	Joins  []string
	models.QueryStringParams
}

type UpdateParams struct {
	TableName string
	Where     string
	OrderBy   string
	Columns   []string
	Data      interface{}
	Args      []interface{}
	Limit     int
	Offset    int
}

type FieldData struct {
	Value interface{}
	Key   string
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type ResponseData[T any] struct {
	Data *[]T        `json:"data"`
	Meta *Pagination `json:"meta"`
}

type repo[T any] struct {
	db  *gorm.DB
	ctx context.Context
}

func InitRepository[T any](store *database.PostgresConfig) Repository[T] {
	return &repo[T]{
		db:  store.DB,
		ctx: store.Ctx,
	}
}

func (r *repo[T]) List(params *FindParams) (*ResponseData[T], error) {
	var (
		entity    T
		data      []T
		totalItem int64
	)

	offset := 0
	limit := 10

	if params.Page > 0 {
		offset = params.Page - 1
	}

	if params.PageSize > 0 {
		limit = params.PageSize
	}

	queryBuilder := r.db.Model(&entity)

	if len(params.Joins) > 0 {
		for _, r := range params.Joins {
			queryBuilder.Joins(r)
		}
	}

	if params.Where != "" {
		queryBuilder = queryBuilder.Where(params.Where, params.Args...)
	}

	_ = queryBuilder.Count(&totalItem)

	if params.OrderBy != "" {
		queryBuilder.Order(fmt.Sprintf("%s.%s %s", queryBuilder.Statement.Table, params.OrderBy, params.OrderDirection))
		params.Select = append(params.Select, fmt.Sprintf("%s.%s", queryBuilder.Statement.Table, params.OrderBy))
	}

	if len(params.Select) > 0 {
		queryBuilder.Select(params.Select)
	}

	//TODO: Join table and select field
	// err := queryBuilder.Preload("Users", func(db *gorm.DB) *gorm.DB {
	// 	return db.Select("id, name, email")
	// }).Find(&data).Error

	err := queryBuilder.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return &ResponseData[T]{
		Data: &data,
		Meta: &Pagination{
			Page:       offset + 1,
			PageSize:   limit,
			TotalItems: totalItem,
			TotalPages: int(math.Ceil(float64(totalItem) / float64(limit))),
		},
	}, nil
}

func (r *repo[T]) DetailByID(id string) (*T, error) {
	// rows, err := r.db.Query(context.Background(), `
	// 	SELECT * FROM Users WHERE id = $1
	// `, id)

	// if err != nil {
	// 	return nil, err
	// }

	// defer rows.Close()
	// data, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])

	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Println(len(data))

	// if len(data) == 0 {
	// 	return nil, nil
	// }
	var data T
	return &data, nil
}

func (r *repo[T]) Delete(id string) error {
	var entity T

	err := r.db.Where("id = ?", id).Delete(&entity).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repo[T]) InsertOne(params *T) (*string, error) {
	var id string

	err := r.db.Create(params).Select("id").Scan(&id).Error
	if err != nil {
		return nil, err
	}

	return &id, nil
}

// CreateUser creates a new user in the db..
// func (r *repo) CreateUser(user model.User) (model.User, error) {
// 	// TODO handle the potential error below.
// 	hashedPass, _ := hashPassword(user.Password)
// 	user.Password = hashedPass

// 	result := repo.db.Create(&user)
// 	fmt.Println(result)
// 	// result := h.db.Create(&user)
// 	// if result.Error
// 	fmt.Println("Inserted a user with ID:", user.ID)
// 	return user, nil
// }

// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

func (r *repo[T]) CountByConditions(params *QueryParams) (*int64, error) {
	var (
		count  int64
		entity T
	)

	queryBuilder := r.db.Model(&entity)

	if len(params.Joins) > 0 {
		for _, r := range params.Joins {
			queryBuilder.Joins(r)
		}
	}

	err := queryBuilder.Where(params.Where, params.Args...).Count(&count).Error
	if err != nil {
		return nil, err
	}

	return &count, nil
}

func (r *repo[T]) DetailByConditions(params *QueryParams) (*T, error) {
	var entity T

	queryBuilder := r.db.Model(&entity)

	if len(params.Joins) > 0 {
		for _, r := range params.Joins {
			queryBuilder.Joins(r)
		}
	}

	err := queryBuilder.Where(params.Where, params.Args...).First(&entity).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return &entity, nil
}

func (r *repo[T]) UpdateByConditions(params *UpdateParams) error {
	err := r.db.Table(params.TableName).Where(params.Where, params.Args...).Updates(params.Data).Error
	if err != nil {
		return err
	}

	return nil
}

func (*repo[T]) generateQuery(params *QueryParams) (string, []interface{}) {
	var sb strings.Builder

	sb.WriteString("SELECT ")

	if len(params.Columns) == 0 {
		sb.WriteString("*")
	} else {
		sb.WriteString(strings.Join(params.Columns, ", "))
	}

	sb.WriteString(fmt.Sprintf(" FROM %s", params.TableName))

	if params.Where != "" {
		sb.WriteString(fmt.Sprintf(" WHERE %s", params.Where))
	}

	if params.OrderBy != "" {
		sb.WriteString(fmt.Sprintf(" ORDER BY %s", params.OrderBy))
	}

	if params.Limit > 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", params.Limit))
	}

	if params.Offset > 0 {
		sb.WriteString(fmt.Sprintf(" OFFSET %d", params.Offset))
	}

	query := sb.String()

	return query, params.Args
}

func (r *repo[T]) strutForScan(u interface{}) []interface{} {
	val := reflect.ValueOf(u).Elem()
	v := make([]interface{}, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		v[i] = valueField.Addr().Interface()
	}

	return v
}

func (r *repo[T]) getStructFields(data interface{}) []FieldData {
	var fields []FieldData

	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Struct {
		fmt.Println("Not a struct.")
		return nil
	}

	typeOf := value.Type()

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := typeOf.Field(i).Name

		if field.Kind() == reflect.Struct {
			embeddedFields := r.getStructFields(field.Interface())
			fields = append(fields, embeddedFields...)
		} else {
			fields = append(fields, FieldData{Key: strings.ToLower(fieldName), Value: field.Interface()})
		}
	}

	return fields
}
