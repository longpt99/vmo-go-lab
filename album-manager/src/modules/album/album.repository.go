package album

import (
	"album-manager/src/common/repository"
	"album-manager/src/configs/database"
	"album-manager/src/models"
	"context"
	"log"

	"gorm.io/gorm"
)

var (
	TableName = "albums"
)

type Repository interface {
	repository.Repository[models.Album]

	// InsertOne(params *models.Album) (*string, error)
	// List() (*[]models.Album, error)
	// DetailByID(id string) (*models.Album, error)
	// Delete(id string) error
	// DetailByConditions(dest interface{}, params *QueryParams) error
	// UpdateByConditions(dest interface{}, params *UpdateParams) error
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

// In this example, the repo struct embeds the repository.Repository interface directly.
// By doing so, the repo struct inherits all the methods from repository.Repository,
// including the InsertOne method.

// Now, when you initialize a repo instance in the InitRepository function,
// you can assign the repository.InitRepository(store) to the Repository field directly.

// By embedding the repository.Repository interface in the repo struct,
// you avoid the need to redefine the InsertOne method in the repo struct
// while still satisfying the Repository interface.
type repo struct {
	db  *gorm.DB
	ctx context.Context
	repository.Repository[models.Album]
}

func InitRepository(store *database.PostgresConfig) Repository {
	err := store.DB.AutoMigrate(&models.Album{})
	if err != nil {
		log.Panicf(`Migrate table "album" failed: %v\n`, err)
	}

	return &repo{
		db:         store.DB,
		ctx:        store.Ctx,
		Repository: repository.InitRepository[models.Album](store),
	}
}

// func (r *repo) List() (*[]models.Album, error) {
// 	var data []models.Album
// 	// rows, err := r.db.Query(context.Background(), `
// 	// 	SELECT *
// 	// 	FROM Albums
// 	// 	WHERE deleted_at IS NULL
// 	// `)

// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// defer rows.Close()
// 	// data, err := pgx.CollectRows(rows, pgx.RowToStructByName[Album])

// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	return &data, nil
// }

// func (r *repo) DetailByID(id string) (*models.Album, error) {
// 	// rows, err := r.db.Query(context.Background(), `
// 	// 	SELECT * FROM Albums WHERE id = $1
// 	// `, id)

// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// defer rows.Close()
// 	// data, err := pgx.CollectRows(rows, pgx.RowToStructByName[Album])

// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// fmt.Println(len(data))

// 	// if len(data) == 0 {
// 	// 	return nil, nil
// 	// }
// 	var data models.Album
// 	return &data, nil
// }

// func (r *repo) Delete(id string) error {
// 	result := r.db.Delete(&models.Album{ID: id})

// 	if result.Error != nil {
// 		return result.Error
// 	}

// 	return nil
// }

// func (r *repo) InsertOne(params *models.Album) (*string, error) {
// 	var id string

// 	result := r.db.Create(params).Select("id").Scan(&id)

// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return &id, nil
// }

// // CreateAlbum creates a new Album in the db..
// // func (r *repo) CreateAlbum(Album model.Album) (model.Album, error) {
// // 	// TODO handle the potential error below.
// // 	hashedPass, _ := hashPassword(Album.Password)
// // 	Album.Password = hashedPass

// // 	result := repo.db.Create(&Album)
// // 	fmt.Println(result)
// // 	// result := h.db.Create(&Album)
// // 	// if result.Error
// // 	fmt.Println("Inserted a Album with ID:", Album.ID)
// // 	return Album, nil
// // }

// // func hashPassword(password string) (string, error) {
// // 	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
// // 	return string(bytes), err
// // }

// func (r *repo) DetailByConditions(dest interface{}, params *QueryParams) error {
// 	// query, args := r.generateQuery(params)
// 	// err := r.db.QueryRow(context.Background(), query, args...).Scan(r.strutForScan(dest)...)

// 	// if err != nil {
// 	// 	if err == pgx.ErrNoRows {
// 	// 		return nil
// 	// 	}

// 	// 	return err
// 	// }

// 	return nil
// }

// func (r *repo) UpdateByConditions(dest interface{}, params *UpdateParams) error {
// 	// query, args := r.generateQuery(params)
// 	// err := r.db.QueryRow(context.Background(), query, args...).Scan(r.strutForScan(dest)...)

// 	// if err != nil {
// 	// 	if err == pgx.ErrNoRows {
// 	// 		return nil
// 	// 	}

// 	// 	return err
// 	// }

// 	// return nil
// 	return nil
// }

// func (r *repo) generateQuery(params *QueryParams) (string, []interface{}) {
// 	var sb strings.Builder

// 	sb.WriteString("SELECT ")

// 	if len(params.Columns) == 0 {
// 		sb.WriteString("*")
// 	} else {
// 		sb.WriteString(strings.Join(params.Columns, ", "))
// 	}

// 	sb.WriteString(fmt.Sprintf(" FROM %s", params.TableName))

// 	if params.Where != "" {
// 		sb.WriteString(fmt.Sprintf(" WHERE %s", params.Where))
// 	}

// 	if params.OrderBy != "" {
// 		sb.WriteString(fmt.Sprintf(" ORDER BY %s", params.OrderBy))
// 	}

// 	if params.Limit > 0 {
// 		sb.WriteString(fmt.Sprintf(" LIMIT %d", params.Limit))
// 	}

// 	if params.Offset > 0 {
// 		sb.WriteString(fmt.Sprintf(" OFFSET %d", params.Offset))
// 	}

// 	query := sb.String()

// 	return query, params.Args
// }

// func (r *repo) strutForScan(u interface{}) []interface{} {
// 	val := reflect.ValueOf(u).Elem()
// 	v := make([]interface{}, val.NumField())

// 	for i := 0; i < val.NumField(); i++ {
// 		valueField := val.Field(i)
// 		v[i] = valueField.Addr().Interface()
// 	}

// 	return v
// }

// func (r *repo) getStructFields(data interface{}) []FieldData {
// 	var fields []FieldData

// 	value := reflect.ValueOf(data)
// 	if value.Kind() != reflect.Struct {
// 		fmt.Println("Not a struct.")
// 		return nil
// 	}

// 	typeOf := value.Type()

// 	for i := 0; i < value.NumField(); i++ {
// 		field := value.Field(i)
// 		fieldName := typeOf.Field(i).Name

// 		if field.Kind() == reflect.Struct {
// 			embeddedFields := r.getStructFields(field.Interface())
// 			fields = append(fields, embeddedFields...)
// 		} else {
// 			fields = append(fields, FieldData{Key: strings.ToLower(fieldName), Value: field.Interface()})
// 		}
// 	}

// 	return fields
// }
