/*
    Listings table 
*/
CREATE TABLE IF NOT EXISTS listings
(
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    active_package text,
    county text NOT NULL,
    observed date NOT NULL,
    removed date,
    asking_price text,
    broker_agency_name text,
    description text,
    fee text,
    housingform text,
    listing_id text UNIQUE,
    land_area text,
    record_type text,
    square_meter_price text,
    street_address text,
    living_and_supplemental_areas text,
    location_description text,
    new_construction boolean,
    project_id text,
    published_at date,
    rooms text,
    slug text,
    price numeric,
    immediate_price numeric,
    PRIMARY KEY (id)
);

/*
    Geo tiles table 
*/
CREATE TABLE IF NOT EXISTS geo_tiles
(
    county text NOT NULL,
    geo_tile text NOT NULL,
    id bigint NOT NULL GENERATED ALWAYS AS IDENTITY,
    UNIQUE (county,geo_tile),
    PRIMARY KEY (id)
);