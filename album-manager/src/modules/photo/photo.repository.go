package photo

import (
	"album-manager/src/configs/database"
	"album-manager/src/models"
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

var (
	TableName = "users"
)

type Repository interface {
	InsertOne(params interface{}) (*string, error)
	List() (*[]models.Photo, error)
	DetailByID(id string) (*models.Photo, error)
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

type repo struct {
	db  *gorm.DB
	ctx context.Context
}

func InitRepository(store *database.PostgresConfig) Repository {
	err := store.DB.AutoMigrate(&models.Photo{})
	if err != nil {
		log.Panicf(`Migrate table "users" failed: %v\n`, err)
	}

	return &repo{
		db:  store.DB,
		ctx: store.Ctx,
	}
}

func (r *repo) List() (*[]models.Photo, error) {
	var data []models.Photo
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

func (r *repo) DetailByID(id string) (*models.Photo, error) {
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
	var data models.Photo
	return &data, nil
}

func (r *repo) Delete(id string) error {
	// _, err := r.db.Exec(context.Background(), `
	// 	UPDATE Users SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1
	// `, id)

	// if err != nil {
	// 	return err
	// }

	return nil
}

func (r *repo) InsertOne(params interface{}) (*string, error) {
	var id string
	result := r.db.Table(TableName).Create(params).Select("id").Scan(&id)

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

func (r *repo) CountByConditions(params *QueryParams) (*int64, error) {
	var count int64
	result := r.db.Table(params.TableName).Where(params.Where, params.Args...).Count(&count)

	if result.Error != nil {
		return nil, result.Error
	}

	return &count, nil
}

func (r *repo) DetailByConditions(dest interface{}, params *QueryParams) error {
	result := r.db.Table(params.TableName).Where(params.Where, params.Args...).First(&dest)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	return nil
}

func (r *repo) UpdateByConditions(params *UpdateParams) error {
	result := r.db.Table(params.TableName).Where(params.Where, params.Args...).Updates(params.Data)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *repo) generateQuery(params *QueryParams) (string, []interface{}) {
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

func (r *repo) strutForScan(u interface{}) []interface{} {
	val := reflect.ValueOf(u).Elem()
	v := make([]interface{}, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		v[i] = valueField.Addr().Interface()
	}

	return v
}

func (r *repo) getStructFields(data interface{}) []FieldData {
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
