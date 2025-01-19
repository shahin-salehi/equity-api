package listing

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

type Store struct {
	db *pgx.Conn
}

func NewStore(db *pgx.Conn) *Store {
	return &Store{db: db}
}

func (s *Store) Listing(l types.Listing) error {
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
	_, err = s.db.Exec(context.Background(), sql_stmt, args)
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
