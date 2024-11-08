package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"githib.com/s4bb4t/leadgen/internal/config"
	"githib.com/s4bb4t/leadgen/internal/lib/models"
	"strings"

	_ "github.com/lib/pq"
)

var SaveStmt *sql.Stmt
var BuildingStmt *sql.Stmt

type Repo struct {
	Db *sql.DB
}

type RepositoryI interface {
	Save(context.Context, models.Building) (models.Building, error)
	Building(context.Context, string) (models.Building, error)
	Buildings(context.Context, models.Query) (models.Buildings, error)
	Close() error
}

func Connect(cfg *config.Config) (RepositoryI, error) {
	const op = "storage.postgres.Connect"

	dbStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Db)

	conn, err := sql.Open("postgres", dbStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	storage := &Repo{Db: conn}

	if err := storage.Db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//err = goose.Up(storage.Db, "./internal/migrations")
	//if err != nil {
	//	return nil, fmt.Errorf("%s: failed to apply migrations: %w", op, err)
	//}

	if err := storage.LoadStmts(); err != nil {
		return nil, fmt.Errorf("%s: failed to load statements: %w", op, err)
	}

	return storage, nil
}

func (repo *Repo) LoadStmts() error {
	var err error
	SaveStmt, err = repo.Db.Prepare(`
	INSERT INTO public.buildings (title, city, floors, year) 
	VALUES ($1, $2, $3, $4) 
	RETURNING id
`)
	if err != nil {
		return err
	}

	BuildingStmt, err = repo.Db.Prepare(`
	SELECT title, city, year, floors 
	FROM public.buildings 
	WHERE title = $1
`)
	if err != nil {
		return err
	}

	return err
}

func (repo *Repo) Save(ctx context.Context, building models.Building) (models.Building, error) {
	const op = "repo.pgsql.Save"

	_, err := SaveStmt.ExecContext(ctx, building.Title, building.City, building.Floors, building.Year)
	if err != nil {
		return building, fmt.Errorf("%s: %w", op, err)
	}

	return building, nil
}

func (repo *Repo) Building(ctx context.Context, title string) (models.Building, error) {
	const op = "repo.pgsql.Building"

	var row models.Building

	err := BuildingStmt.QueryRowContext(ctx, title).Scan(&row.Title, &row.City, &row.Floors, &row.Year)
	if err != nil {
		return models.Building{}, fmt.Errorf("%s: %w", op, err)
	}

	return row, nil
}

func (repo *Repo) Buildings(ctx context.Context, query models.Query) (models.Buildings, error) {
	const op = "repo.pgsql.Buildings"

	var all models.Buildings
	var row models.Building

	// Базовый запрос и запрос для подсчета
	baseQuery := "SELECT Title, city, floors, year FROM public.buildings"
	countQuery := "SELECT COUNT(*) FROM public.buildings"
	var whereClauses []string
	var args []any
	argIdx := 0

	// Применение условий по городам, этажам и году, если они заданы
	if query.City != "" {
		argIdx++
		whereClauses = append(whereClauses, fmt.Sprintf("city = $%d", argIdx))
		args = append(args, query.City)
		all.Meta.Query.City = query.City
	}
	if query.Floors > 0 {
		argIdx++
		whereClauses = append(whereClauses, fmt.Sprintf("floors = $%d", argIdx))
		args = append(args, query.Floors)
		all.Meta.Query.Floors = query.Floors
	}
	if query.Year > 0 {
		argIdx++
		whereClauses = append(whereClauses, fmt.Sprintf("year = $%d", argIdx))
		args = append(args, query.Year)
		all.Meta.Query.Year = query.Year
	}

	if len(whereClauses) > 0 {
		whereSQL := " WHERE " + strings.Join(whereClauses, " AND ")
		baseQuery += whereSQL
		countQuery += whereSQL
	}

	if query.Limit == 0 {
		query.Limit = 10
	}

	if query.Limit > 0 {
		argIdx++
		all.Meta.Query.Limit = query.Limit
		baseQuery += fmt.Sprintf(" LIMIT $%d", argIdx)
		args = append(args, query.Limit)
	}

	argIdx++
	all.Meta.Query.Offset = query.Offset
	baseQuery += fmt.Sprintf(" OFFSET $%d", argIdx)
	args = append(args, query.Offset*query.Limit)

	rows, err := repo.Db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return models.Buildings{}, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&row.Title, &row.City, &row.Floors, &row.Year)
		if err != nil {
			return models.Buildings{}, fmt.Errorf("%s: %w", op, err)
		}
		all.Data = append(all.Data, row)
	}

	var countArgs []any
	if len(whereClauses) > 0 {
		countArgs = args[:len(args)-2]
	}

	err = repo.Db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&all.Meta.TotalAmount)
	if err != nil {
		return models.Buildings{}, fmt.Errorf("%s: %w", op, err)
	}

	if all.Data == nil {
		all.Data = []models.Building{}
	}
	return all, nil
}

func (repo *Repo) Close() error {
	var err error

	err = BuildingStmt.Close()
	if err != nil {
		return err
	}

	err = SaveStmt.Close()
	if err != nil {
		return err
	}

	err = repo.Db.Close()
	if err != nil {
		return err
	}

	return nil
}
