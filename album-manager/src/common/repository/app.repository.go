package repository

import (
	"album-manager/src/configs/database"
	"album-manager/src/models"
	"context"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Repository[T any] interface {
	InsertOne(params *T) (*string, error)
	List() (*[]T, error)
	DetailByID(id string) (*T, error)
	Delete(id string) error
	DetailByConditions(dest interface{}, params *QueryParams) error
	UpdateByConditions(params *UpdateParams) error
	CountByConditions(params *QueryParams) (*int64, error)
}

type QueryParams struct {
	Columns   []string
	TableName string
	Where     string
	OrderBy   string
	Limit     int
	Offset    int
	Args      []interface{}
}

type UpdateParams struct {
	Columns   []string
	TableName string
	Where     string
	OrderBy   string
	Limit     int
	Offset    int
	Args      []interface{}
	Data      interface{}
}

type FieldData struct {
	Key   string
	Value interface{}
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

func (r *repo[T]) List() (*[]T, error) {
	var data []T
	// rows, err := r.db.Query(context.Background(), `
	// 	SELECT *
	// 	FROM users
	// 	WHERE deleted_at IS NULL
	// `)

	// if err != nil {
	// 	return nil, err
	// }

	// defer rows.Close()
	// data, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])

	// if err != nil {
	// 	return nil, err
	// }

	return &data, nil
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
	result := r.db.Delete(models.User{ID: id})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repo[T]) InsertOne(params *T) (*string, error) {
	var id string
	result := r.db.Create(params).Select("id").Scan(&id)

	if result.Error != nil {
		return nil, result.Error
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
	var count int64
	result := r.db.Table(params.TableName).Where(params.Where, params.Args...).Count(&count)

	if result.Error != nil {
		return nil, result.Error
	}

	return &count, nil
}

func (r *repo[T]) DetailByConditions(dest interface{}, params *QueryParams) error {
	result := r.db.Table(params.TableName).Where(params.Where, params.Args...).First(&dest)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	return nil
}

func (r *repo[T]) UpdateByConditions(params *UpdateParams) error {
	result := r.db.Table(params.TableName).Where(params.Where, params.Args...).Updates(params.Data)

	if result.Error != nil {
		return result.Error
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
