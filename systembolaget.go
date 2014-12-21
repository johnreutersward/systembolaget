// Package systembolaget provides a client library for accessing the Systembolaget API.
package systembolaget

import (
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

const (
	libraryVersion = "0.2"
	baseURL        = "http://www.systembolaget.se/"
	userAgent      = "systembolaget/" + libraryVersion
)

type Client struct {
	httpClient *http.Client
	BaseURL    *url.URL
	UserAgent  string
}

type Articles struct {
	CreatedAt   string    `xml:"skapad-tid,omitempty"`
	InfoMessage string    `xml:"info>meddelande,omitempty"`
	Articles    []Article `xml:"artikel"`
}

type Article struct {
	Number            string `xml:"nr,omitempty"`
	ArticleId         string `xml:"Artikelid,omitempty"` // Unique number.
	ProductNumber     string `xml:"Varnummer,omitempty"` // One or more articles can share the same ProductNumber.
	Name              string `xml:"Namn,omitempty"`
	SubName           string `xml:"Namn2,omitempty"`
	Price             string `xml:"Prisinklmoms,omitempty"` // Price, VAT included.
	VolumeMl          string `xml:"Volymiml,omitempty"`     // Volume in millilitre.
	PricePerLitre     string `xml:"PrisPerLiter,omitempty"`
	SoldSince         string `xml:"Saljstart,omitempty"`
	DiscontinuedSince string `xml:"Slutlev,omitempty"`
	ProductGroup      string `xml:"Varugrupp,omitempty"`
	ContainerType     string `xml:"Forpackning,omitempty"`
	SealingType       string `xml:"Forslutning,omitempty"`
	Origin            string `xml:"Ursprung,omitempty"`
	OriginCountry     string `xml:"Ursprunglandnamn,omitempty"`
	Producer          string `xml:"Producent,omitempty"`
	Distributer       string `xml:"Leverantor,omitempty"`
	AnnualVolume      string `xml:"Argang,omitempty"`
	TestAnnualValume  string `xml:"Provadargang,omitempty"`
	AlcoholPrecentage string `xml:"Alkoholhalt,omitempty"`
	Assortment        string `xml:"Sortiment,omitempty"`
	Organic           string `xml:"Ekologisk,omitempty"`
	Kosher            string `xml:"Koscher,omitempty"`
	PrimaryProducts   string `xml:"RavarorBeskrivning,omitempty"`
}

type Stores struct {
	InfoMessage string  `xml:"Info>Meddelande"`
	Stores      []Store `xml:"ButikOmbud"`
}

type Store struct {
	Number       string `xml:"Nr,omitempty"`
	Type         string `xml:"Typ,omitempty"`
	Adress1      string `xml:"Address1,omitempty"`
	Adress2      string `xml:"Address2,omitempty"`
	Adress3      string `xml:"Address3,omitempty"`
	Adress4      string `xml:"Address4,omitempty"`
	Adress5      string `xml:"Address5,omitempty"`
	PhoneNumber  string `xml:"Telefon,omitempty"`
	StoreType    string `xml:"ButiksTyp,omitempty"`
	Services     string `xml:"Tjanster,omitempty"`
	SearchWords  string `xml:"SokOrd,omitempty"`
	OpeningHours string `xml:"Oppettider,omitempty"`
	RT90x        string `xml:"RT90x,omitempty"`
	RT90y        string `xml:"RT90y,omitempty"`
}

// Returns a new Systembolaget API client. If a nil
// httpClient is provided, http.DefaultClient will be used.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(baseURL)
	c := &Client{
		httpClient: httpClient,
		BaseURL:    baseURL,
		UserAgent:  userAgent,
	}
	return c
}

// Returns all stores. Note: downloads large XML file.
func (c *Client) Stores() (*Stores, error) {
	ref := "Assortment.aspx?butikerombud=1"

	stores := &Stores{}
	res, err := c.get(ref)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = xml.NewDecoder(res.Body).Decode(stores); err != nil {
		return nil, err
	}

	return stores, err
}

// Returns all articles. Note: downloads large XML file.
func (c *Client) Articles() (*Articles, error) {
	ref := "Assortment.aspx?Format=Xml"

	articles := &Articles{}
	res, err := c.get(ref)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err = xml.NewDecoder(res.Body).Decode(articles); err != nil {
		return nil, err
	}

	return articles, err
}

func (c *Client) get(endp string) (*http.Response, error) {
	ref, err := url.Parse(endp)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(ref)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return nil, errors.New("api error, response code: " + strconv.Itoa(resp.StatusCode))
	}

	return resp, err
}
