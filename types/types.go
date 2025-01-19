package types

type County struct {
	County string `json:"county"`
}

type GeoTile struct {
	County  string `json:"county"`
	GeoTile string `json:"GeoTile"`
}

type ScrapedIDs struct {
	IDs []string
}

type Delta struct {
	NewIDs []string `json:"newIDs"`
}

type HousingForm struct {
	Symbol string `json:"symbol"`
}

type Listing struct {
	ActivePackage              string      `json:"activePackage"`
	AskingPrice                string      `json:"askingPrice"`
	BrokerAgencyName           string      `json:"brokerAgencyName"`
	Description                string      `json:"description"`
	Fee                        string      `json:"fee"`
	Floor                      string      `json:"floor"`
	HousingForm                HousingForm `json:"housingForm"`
	Id                         string      `json:"id"`
	LandArea                   string      `json:"landArea"`
	LivingAndSupplementalAreas string      `json:"livingAndSupplementalAreas"`
	LocationDescription        string      `json:"locationDescription"`
	NewConstruction            bool        `json:"newConstruction"`
	ProjectId                  string      `json:"projectId"`
	PublishedAt                string      `json:"publishedAt"`
	RecordType                 string      `json:"recordType"`
	Rooms                      string      `json:"rooms"`
	Slug                       string      `json:"slug"`
	SquareMeterPrice           string      `json:"squareMeterPrice"`
	StreetAddress              string      `json:"streetAddress"`
	County                     string      `json:"county"`
	Price                      int         `json:"price"`
	ImmediatePrice             int         `json:"immediatePrice"`
}
