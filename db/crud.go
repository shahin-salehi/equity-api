package db

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/shahin-salehi/equity-api/types"
)

type CRUD interface {
	Listing(l types.Listing) error
	Delta(ids []string) (*types.Delta, error)
	ReadCounties() ([]types.County, error)
	TilesByCounty(county string) ([]types.GeoTile, error)
	InsertManyTiles(tiles []types.GeoTile) error
}

type crud struct {
	db *pgx.Conn
}

func NewRepo(conn *pgx.Conn) *crud {
	return &crud{db: conn}
}

func (c *crud) InsertManyTiles(tiles []types.GeoTile) error {
	tableName := "geo_tiles"
	// when do we release the memory allocated by this array?
	columns := []string{"county", "geo_tile"}
	entries := [][]any{}
	// can we avoid this loop?
	for _, geoTile := range tiles {
		entries = append(entries, []any{geoTile.County, geoTile.GeoTile})
	}

	_, err := c.db.CopyFrom(
		context.Background(),
		pgx.Identifier{tableName},
		columns,
		pgx.CopyFromRows(entries),
	)
	if err != nil {
		slog.Error("failed to insert many tiles", slog.Any("error", err))
		return err
	}
	return nil
}

func (c *crud) TilesByCounty(county string) ([]types.GeoTile, error) {
	sql_stmt := `SELECT geo_tile, county FROM geo_tiles WHERE county=@County;`
	args := pgx.NamedArgs{
		"County": county,
	}

	rows, err := c.db.Query(context.Background(), sql_stmt, args)
	if err != nil {
		slog.Error("failed to read geotiles from database.", slog.Any("error", err))
		return nil, err
	}

	structuredRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.GeoTile])
	if err != nil {
		slog.Error("failed to collect geotile rows to structured format.", slog.Any("error", err))
		return nil, err
	}

	return structuredRows, nil
}

func (c *crud) Delta(ids []string) (*types.Delta, error) {
	// stmt
	sql_stmt := `SELECT listings_delta($1::text[]) AS NewIDs;`

	// execute query
	rows, err := c.db.Query(context.Background(), sql_stmt, ids)
	if err != nil {
		slog.Error("failed to execute query", slog.Any("function", "FilterIDs"), slog.Any("error", err))
		return nil, err
	}

	// strucutre rows
	structuredRows, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[types.Delta])
	if err != nil {
		slog.Error("Failed to strucutre counties ", slog.Any("error", err))
		return nil, err
	}

	return &structuredRows, nil
}

func (c *crud) ReadCounties() ([]types.County, error) {
	sql_stmt := `SELECT DISTINCT county FROM geo_tiles;`
	rows, err := c.db.Query(context.Background(), sql_stmt)
	if err != nil {
		slog.Error("failed to get distinct counties ", slog.Any("error", err))
		return nil, err
	}

	// strucutre rows
	structuredRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[types.County])
	if err != nil {
		slog.Error("Failed to strucutre counties ", slog.Any("error", err))
		return nil, err
	}

	return structuredRows, nil
}

func (c *crud) Listing(l types.Listing) error {
	// stmt
	sql_stmt := `INSERT INTO listings(
		active_package, county, observed, asking_price, broker_agency_name, description, fee, housingform, listing_id, land_area, record_type, square_meter_price, street_address, living_and_supplemental_areas, location_description, new_construction, project_id, published_at, rooms, slug, price, immediate_price)
		VALUES (@active_package, @county, @observed, @asking_price, @broker_agency_name, @description, @fee, @housingform, @listing_id, @land_area, @record_type, @square_meter_price, @street_address, @living_and_supplemental_areas, @location_description, @new_construction, @project_id, @published_at, @rooms, @slug, @price, @immediate_price);`

	// args
	// conv time
	publishedAt, err := strconv.ParseFloat(l.PublishedAt, 64)
	if err != nil {
		slog.Error("failed to conv published at to int", slog.Any("error", err))
		return err
	}
	args := pgx.NamedArgs{
		"active_package":                l.ActivePackage,
		"county":                        l.County,
		"observed":                      time.Now(),
		"asking_price":                  l.AskingPrice,
		"broker_agency_name":            l.BrokerAgencyName,
		"description":                   l.Description,
		"fee":                           l.Fee,
		"housingform":                   l.HousingForm.Symbol,
		"l_id":                          l.Id,
		"land_area":                     l.LandArea,
		"record_type":                   l.RecordType,
		"square_meter_price":            l.SquareMeterPrice,
		"street_address":                l.StreetAddress,
		"living_and_supplemental_areas": l.LivingAndSupplementalAreas,
		"location_description":          l.LocationDescription,
		"new_construction":              l.NewConstruction,
		"project_id":                    l.ProjectId,
		"published_at":                  time.Unix(int64(publishedAt), 0),
		"rooms":                         l.Rooms,
		"slug":                          l.Slug,
		"price":                         l.Price,
		"immediate_price":               l.ImmediatePrice,
	}

	// exec
	_, err = c.db.Exec(context.Background(), sql_stmt, args)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				slog.Warn("duplicate listing rejected.", slog.Any("warning", err))
				return nil
			} // duplicate
		}
		slog.Error("failed to insert listing ", slog.Any("error", err))
		return err
	}
	return nil

}
