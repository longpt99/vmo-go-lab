package photo

import (
	"album-manager/src/configs/database"
	"context"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type Repository interface {
	InsertOne(params interface{}) (*string, error)
	List() (*[]Photo, error)
	DetailByID(id string) (*Photo, error)
	Delete(id string) error
	DetailByConditions(dest interface{}, params *QueryParams) error
	UpdateByConditions(dest interface{}, params *UpdateParams) error
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
	store.DB.AutoMigrate(&Photo{})
	return &repo{
		db:  store.DB,
		ctx: store.Ctx,
	}
}

func (r *repo) List() (*[]Photo, error) {
	var data []Photo
	// rows, err := r.db.Query(context.Background(), `
	// 	SELECT *
	// 	FROM Photos
	// 	WHERE deleted_at IS NULL
	// `)

	// if err != nil {
	// 	return nil, err
	// }

	// defer rows.Close()
	// data, err := pgx.CollectRows(rows, pgx.RowToStructByName[Photo])

	// if err != nil {
	// 	return nil, err
	// }

	return &data, nil
}

func (r *repo) DetailByID(id string) (*Photo, error) {
	// rows, err := r.db.Query(context.Background(), `
	// 	SELECT * FROM Photos WHERE id = $1
	// `, id)

	// if err != nil {
	// 	return nil, err
	// }

	// defer rows.Close()
	// data, err := pgx.CollectRows(rows, pgx.RowToStructByName[Photo])

	// if err != nil {
	// 	return nil, err
	// }

	// fmt.Println(len(data))

	// if len(data) == 0 {
	// 	return nil, nil
	// }
	var data Photo
	return &data, nil
}

func (r *repo) Delete(id string) error {
	// _, err := r.db.Exec(context.Background(), `
	// 	UPDATE Photos SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1
	// `, id)

	// if err != nil {
	// 	return err
	// }

	return nil
}

func (r *repo) InsertOne(params interface{}) (*string, error) {
	var id string

	// data := r.getStructFields(params)

	// // Solution 1 pass args
	// // query := `
	// // 	INSERT INTO admin(name, description)
	// // 	VALUES (@name, @description)
	// // 	RETURNING id;
	// // `
	// // args := pgx.NamedArgs{
	// // 	"name":        name,
	// // 	"description": description,
	// // }
	// // err := r.db.QueryRow(context.Background(), query, args).Scan(&id)

	// // Solution 2 pass args
	// // query :=

	// var args []interface{}
	// var valueStrings []string
	// var keys []string

	// for i, d := range data {
	// 	valueStrings = append(valueStrings, fmt.Sprintf("$%d", i+1))
	// 	keys = append(keys, d.Key)
	// 	args = append(args, d.Value)
	// }

	// query := fmt.Sprintf(`
	// INSERT INTO Photos(%s)
	// VALUES (%s)
	// RETURNING id;
	// `, strings.Join(keys, ","), strings.Join(valueStrings, ","))

	// err := r.db.QueryRow(context.Background(), query, args...).Scan(&id)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return nil, err
	// }

	return &id, nil
}

// CreatePhoto creates a new Photo in the db..
// func (r *repo) CreatePhoto(Photo model.Photo) (model.Photo, error) {
// 	// TODO handle the potential error below.
// 	hashedPass, _ := hashPassword(Photo.Password)
// 	Photo.Password = hashedPass

// 	result := repo.db.Create(&Photo)
// 	fmt.Println(result)
// 	// result := h.db.Create(&Photo)
// 	// if result.Error
// 	fmt.Println("Inserted a Photo with ID:", Photo.ID)
// 	return Photo, nil
// }

// func hashPassword(password string) (string, error) {
// 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	return string(bytes), err
// }

func (r *repo) DetailByConditions(dest interface{}, params *QueryParams) error {
	// query, args := r.generateQuery(params)
	// err := r.db.QueryRow(context.Background(), query, args...).Scan(r.strutForScan(dest)...)

	// if err != nil {
	// 	if err == pgx.ErrNoRows {
	// 		return nil
	// 	}

	// 	return err
	// }

	return nil
}

func (r *repo) UpdateByConditions(dest interface{}, params *UpdateParams) error {
	// query, args := r.generateQuery(params)
	// err := r.db.QueryRow(context.Background(), query, args...).Scan(r.strutForScan(dest)...)

	// if err != nil {
	// 	if err == pgx.ErrNoRows {
	// 		return nil
	// 	}

	// 	return err
	// }

	// return nil
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
