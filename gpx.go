// Package gpx provides convenience types for reading and writing GPX
// documents.
// See http://www.topografix.com/gpx.asp.
package gpx

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	geom "github.com/twpayne/go-geom"
	"golang.org/x/net/html/charset"
)

const (
	http  = "http://"
	https = "https://"
)

var timeLayouts = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999",
}

// StartElement is the XML start element for GPX files.
var StartElement = xml.StartElement{
	Name: xml.Name{Local: "gpx"},
}

var copyrightYearLayouts = []string{
	"2006",
	"2006Z",
	"2006-07:00",
}

var errNoTimeLayout = errors.New("no time layout")

// A BoundsType is a boundsType.
type BoundsType struct {
	MinLat float64 `xml:"minlat,attr"`
	MinLon float64 `xml:"minlon,attr"`
	MaxLat float64 `xml:"maxlat,attr"`
	MaxLon float64 `xml:"maxlon,attr"`
}

// A CopyrightType is a copyrightType.
type CopyrightType struct {
	Author  string `xml:"author,attr"`
	Year    int    `xml:"year,omitempty"`
	License string `xml:"license,omitempty"`
}

// An ExtensionsType contains elements from another schema.
type ExtensionsType struct {
	XML []byte `xml:",innerxml"`
}

// A GPX is a gpxType.
type GPX struct {
	XMLName            string            `xml:"gpx"`
	XMLSchemaLocations []string          `xml:"xsi:schemaLocation,attr"`
	XMLAttrs           map[string]string `xml:"-"`
	Version            string            `xml:"version,attr"`
	Creator            string            `xml:"creator,attr"`
	Metadata           *MetadataType     `xml:"metadata,omitempty"`
	Wpt                []*WptType        `xml:"wpt,omitempty"`
	Rte                []*RteType        `xml:"rte,omitempty"`
	Trk                []*TrkType        `xml:"trk,omitempty"`
	Extensions         *ExtensionsType   `xml:"extensions,omitempty"`
}

// A LinkType is a linkType.
type LinkType struct {
	HREF string `xml:"href,attr"`
	Text string `xml:"text,omitempty"`
	Type string `xml:"type,omitempty"`
}

// A PersonType is a personType.
type PersonType struct {
	Name  string     `xml:"name,omitempty"`
	Email *EmailType `xml:"email,omitempty"`
	Link  *LinkType  `xml:"link,omitempty"`
}

// An EmailType is an emailType.
type EmailType struct {
	Name   string `xml:"id,attr"`
	Domain string `xml:"domain,attr"`
}

// A MetadataType is a metadataType.
type MetadataType struct {
	Name       string          `xml:"name,omitempty"`
	Desc       string          `xml:"desc,omitempty"`
	Author     *PersonType     `xml:"author,omitempty"`
	Copyright  *CopyrightType  `xml:"copyright,omitempty"`
	Link       []*LinkType     `xml:"link,omitempty"`
	Time       time.Time       `xml:"time,omitempty"`
	Keywords   string          `xml:"keywords,omitempty"`
	Bounds     *BoundsType     `xml:"bounds,omitempty"`
	Extensions *ExtensionsType `xml:"extensions,omitempty"`
}

// A RteType is a rteType.
type RteType struct {
	Name       string          `xml:"name,omitempty"`
	Cmt        string          `xml:"cmt,omitempty"`
	Desc       string          `xml:"desc,omitempty"`
	Src        string          `xml:"src,omitempty"`
	Link       []*LinkType     `xml:"link,omitempty"`
	Number     int             `xml:"number,omitempty"`
	Type       string          `xml:"type,omitempty"`
	Extensions *ExtensionsType `xml:"extensions,omitempty"`
	RtePt      []*WptType      `xml:"rtept,omitempty"`
}

// A TrkSegType is a trkSegType.
type TrkSegType struct {
	TrkPt      []*WptType      `xml:"trkpt,omitempty"`
	Extensions *ExtensionsType `xml:"extensions,omitempty"`
}

// A TrkType is a trkType.
type TrkType struct {
	Name       string          `xml:"name,omitempty"`
	Cmt        string          `xml:"cmt,omitempty"`
	Desc       string          `xml:"desc,omitempty"`
	Src        string          `xml:"src,omitempty"`
	Link       []*LinkType     `xml:"link,omitempty"`
	Number     int             `xml:"number,omitempty"`
	Type       string          `xml:"type,omitempty"`
	Extensions *ExtensionsType `xml:"extensions,omitempty"`
	TrkSeg     []*TrkSegType   `xml:"trkseg,omitempty"`
}

// A WptType is a wptType.
type WptType struct {
	Lat           float64         `xml:"lat,omitempty"`
	Lon           float64         `xml:"lon,omitempty"`
	Ele           float64         `xml:"ele,omitempty"`
	Speed         float64         `xml:"speed,omitempty"`
	Course        float64         `xml:"course,omitempty"`
	Time          time.Time       `xml:"time,omitempty"`
	MagVar        float64         `xml:"magvar,omitempty"`
	GeoidHeight   float64         `xml:"geoidheight,omitempty"`
	Name          string          `xml:"name,omitempty"`
	Cmt           string          `xml:"cmt,omitempty"`
	Desc          string          `xml:"desc,omitempty"`
	Src           string          `xml:"src,omitempty"`
	Link          []*LinkType     `xml:"link,omitempty"`
	Sym           string          `xml:"sym,omitempty"`
	Type          string          `xml:"type,omitempty"`
	Fix           string          `xml:"fix,omitempty"`
	Sat           int             `xml:"sat,omitempty"`
	HDOP          float64         `xml:"hdop,omitempty"`
	VDOP          float64         `xml:"vdop,omitempty"`
	PDOP          float64         `xml:"pdop,omitempty"`
	AgeOfDGPSData float64         `xml:"ageofdgpsdata,omitempty"`
	DGPSID        []int           `xml:"dgpsid,omitempty"`
	Extensions    *ExtensionsType `xml:"extensions,omitempty"`
}

// UnmarshalXML implements xml.Unmarshaler.UnmarshalXML.
func (c *CopyrightType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	alias := struct {
		Author  string  `xml:"author,attr"`
		Year    *string `xml:"year,omitempty"`
		License string  `xml:"license,omitempty"`
	}{}

	err := d.DecodeElement(&alias, &start)
	if err != nil {
		return err
	}

	c.Author = alias.Author
	c.License = alias.License

	if alias.Year == nil {
		return nil
	}

	for _, layout := range copyrightYearLayouts {
		var date time.Time
		date, err = time.Parse(layout, *alias.Year)
		if err == nil {
			c.Year = date.Year()
			return nil
		}
	}

	return fmt.Errorf("couldn't parse Copyright year: %s", *alias.Year)
}

type ReadOption func()

// Read reads a new GPX from r.
func Read(r io.Reader, options ...ReadOption) (*GPX, error) {
	for _, option := range options {
		option()
	}
	gpx := &GPX{}
	d := xml.NewDecoder(r)
	d.CharsetReader = charset.NewReaderLabel
	return gpx, d.Decode(gpx)
}

// WithTimeLayout applies a custom time layout for the decoding of the GPX source.
func WithTimeLayout(layout string) ReadOption {
	return func() {
		timeLayouts = []string{layout}
	}
}

// WithTimeLayouts applies a custom time layouts for the decoding of the GPX
// source.
func WithTimeLayouts(layouts []string) ReadOption {
	return func() {
		timeLayouts = layouts
	}
}

// MarshalXML implements xml.Marshaler.MarshalXML.
func (g *GPX) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	baseURL := "www.topografix.com/GPX/" + strings.Join(strings.Split(g.Version, "."), "/")
	xmlSchemaLocations := append([]string{
		http + baseURL,
		https + baseURL + "/gpx.xsd",
	}, g.XMLSchemaLocations...)
	attr := []xml.Attr{
		{
			Name:  xml.Name{Local: "version"},
			Value: g.Version,
		},
		{
			Name:  xml.Name{Local: "creator"},
			Value: g.Creator,
		},
		{
			Name:  xml.Name{Local: "xmlns:xsi"},
			Value: http + "www.w3.org/2001/XMLSchema-instance",
		},
		{
			Name:  xml.Name{Local: "xmlns"},
			Value: http + baseURL,
		},
		{
			Name:  xml.Name{Local: "xsi:schemaLocation"},
			Value: strings.Join(xmlSchemaLocations, " "),
		},
	}
	for k, v := range g.XMLAttrs {
		attr = append(attr, xml.Attr{
			Name:  xml.Name{Local: k},
			Value: v,
		})
	}
	start := xml.StartElement{
		Name: xml.Name{Local: "gpx"},
		Attr: attr,
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := e.EncodeElement(g.Metadata, xml.StartElement{Name: xml.Name{Local: "metadata"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(g.Wpt, xml.StartElement{Name: xml.Name{Local: "wpt"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(g.Rte, xml.StartElement{Name: xml.Name{Local: "rte"}}); err != nil {
		return err
	}
	if err := e.EncodeElement(g.Trk, xml.StartElement{Name: xml.Name{Local: "trk"}}); err != nil {
		return err
	}
	return e.EncodeToken(start.End())
}

// Write writes g to w.
func (g *GPX) Write(w io.Writer) error {
	return xml.NewEncoder(w).EncodeElement(g, StartElement)
}

// WriteIndent writes g to w.
func (g *GPX) WriteIndent(w io.Writer, prefix, indent string) error {
	e := xml.NewEncoder(w)
	e.Indent(prefix, indent)
	return e.EncodeElement(g, StartElement)
}

// UnmarshalXML implements xml.Unmarshaler.UnmarshalXML.
func (m *MetadataType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var e struct {
		Name       string          `xml:"name"`
		Desc       string          `xml:"desc"`
		Author     *PersonType     `xml:"author"`
		Copyright  *CopyrightType  `xml:"copyright"`
		Link       []*LinkType     `xml:"link"`
		Time       string          `xml:"time"`
		Keywords   string          `xml:"keywords"`
		Bounds     *BoundsType     `xml:"bounds"`
		Extensions *ExtensionsType `xml:"extensions"`
	}
	if err := d.DecodeElement(&e, &start); err != nil {
		return err
	}
	mt := MetadataType{
		Name:       e.Name,
		Desc:       e.Desc,
		Author:     e.Author,
		Copyright:  e.Copyright,
		Link:       e.Link,
		Keywords:   e.Keywords,
		Bounds:     e.Bounds,
		Extensions: e.Extensions,
	}
	if e.Time != "" {
		t, err := parseTime(e.Time)
		if err != nil {
			return err
		}
		mt.Time = t
	}
	*m = mt
	return nil
}

// NewRteType returns a new RteType with geometry g.
func NewRteType(g *geom.LineString) *RteType {
	return &RteType{
		RtePt: newWptTypes(g),
	}
}

// Geom returns r's geometry.
func (r *RteType) Geom(layout geom.Layout) *geom.LineString {
	flatCoords := make([]float64, 0, layout.Stride()*len(r.RtePt))
	for _, rp := range r.RtePt {
		flatCoords = rp.appendFlatCoords(flatCoords, layout)
	}
	return geom.NewLineStringFlat(layout, flatCoords)
}

// NewTrkType returns a new TrkType with geometry g.
func NewTrkType(g *geom.MultiLineString) *TrkType {
	trkSegs := make([]*TrkSegType, g.NumLineStrings())
	for i := range trkSegs {
		trkSegs[i] = NewTrkSegType(g.LineString(i))
	}
	return &TrkType{
		TrkSeg: trkSegs,
	}
}

// Geom returns t's geometry.
func (t *TrkType) Geom(layout geom.Layout) *geom.MultiLineString {
	ends := make([]int, len(t.TrkSeg))
	end := 0
	for i, ts := range t.TrkSeg {
		end += layout.Stride() * len(ts.TrkPt)
		ends[i] = end
	}
	flatCoords := make([]float64, 0, end)
	for _, ts := range t.TrkSeg {
		for _, tp := range ts.TrkPt {
			flatCoords = tp.appendFlatCoords(flatCoords, layout)
		}
	}
	return geom.NewMultiLineStringFlat(layout, flatCoords, ends)
}

// NewTrkSegType returns a new TrkSegType with geometry g.
func NewTrkSegType(g *geom.LineString) *TrkSegType {
	return &TrkSegType{
		TrkPt: newWptTypes(g),
	}
}

// Geom returns ts's geometry.
func (ts *TrkSegType) Geom(layout geom.Layout) *geom.LineString {
	flatCoords := make([]float64, 0, layout.Stride()*len(ts.TrkPt))
	for _, tp := range ts.TrkPt {
		flatCoords = tp.appendFlatCoords(flatCoords, layout)
	}
	return geom.NewLineStringFlat(layout, flatCoords)
}

// NewWptType returns a new WptType with geometry g.
func NewWptType(g *geom.Point) *WptType {
	flatCoords := g.FlatCoords()
	layout := g.Layout()
	w := &WptType{
		Lat: flatCoords[1],
		Lon: flatCoords[0],
	}
	if zIndex := layout.ZIndex(); zIndex != -1 {
		w.Ele = flatCoords[zIndex]
	}
	if mIndex := layout.MIndex(); mIndex != -1 {
		w.Time = MToTime(flatCoords[mIndex])
	}
	return w
}

// Geom returns w's geometry.
func (w *WptType) Geom(layout geom.Layout) *geom.Point {
	return geom.NewPointFlat(layout, w.appendFlatCoords(make([]float64, 0, layout.Stride()), layout))
}

// MarshalXML implements xml.Marshaler.MarshalXML.
func (w *WptType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	latAttr := xml.Attr{
		Name:  xml.Name{Local: "lat"},
		Value: strconv.FormatFloat(w.Lat, 'f', -1, 64),
	}
	lonAttr := xml.Attr{
		Name:  xml.Name{Local: "lon"},
		Value: strconv.FormatFloat(w.Lon, 'f', -1, 64),
	}
	start.Attr = append(start.Attr, latAttr, lonAttr)
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	if err := maybeEmitFloatElement(e, "ele", w.Ele); err != nil {
		return err
	}
	if err := maybeEmitFloatElement(e, "speed", w.Speed); err != nil {
		return err
	}
	if err := maybeEmitFloatElement(e, "course", w.Course); err != nil {
		return err
	}
	if !w.Time.IsZero() {
		if err := maybeEmitStringElement(e, "time", w.Time.UTC().Format(timeLayouts[0])); err != nil {
			return err
		}
	}
	if err := maybeEmitFloatElement(e, "magvar", w.MagVar); err != nil {
		return err
	}
	if err := maybeEmitFloatElement(e, "geoidheight", w.GeoidHeight); err != nil {
		return err
	}
	if err := maybeEmitStringElement(e, "name", w.Name); err != nil {
		return err
	}
	if err := maybeEmitStringElement(e, "cmt", w.Cmt); err != nil {
		return err
	}
	if err := maybeEmitStringElement(e, "desc", w.Desc); err != nil {
		return err
	}
	if err := maybeEmitStringElement(e, "src", w.Src); err != nil {
		return err
	}
	if w.Link != nil {
		if err := e.EncodeElement(w.Link, xml.StartElement{Name: xml.Name{Local: "link"}}); err != nil {
			return err
		}
	}
	if err := maybeEmitStringElement(e, "sym", w.Sym); err != nil {
		return err
	}
	if err := maybeEmitStringElement(e, "type", w.Type); err != nil {
		return err
	}
	if err := maybeEmitStringElement(e, "fix", w.Fix); err != nil {
		return err
	}
	if err := maybeEmitIntElement(e, "sat", w.Sat); err != nil {
		return err
	}
	if err := maybeEmitFloatElement(e, "hdop", w.HDOP); err != nil {
		return err
	}
	if err := maybeEmitFloatElement(e, "vdop", w.VDOP); err != nil {
		return err
	}
	if err := maybeEmitFloatElement(e, "pdop", w.PDOP); err != nil {
		return err
	}
	if err := maybeEmitFloatElement(e, "ageofdgpsdata", w.AgeOfDGPSData); err != nil {
		return err
	}
	for _, dgpsid := range w.DGPSID {
		if err := emitIntElement(e, "dgpsid", dgpsid); err != nil {
			return err
		}
	}
	if w.Extensions != nil {
		if err := e.EncodeElement(w.Extensions, xml.StartElement{Name: xml.Name{Local: "extensions"}}); err != nil {
			return err
		}
	}
	return e.EncodeToken(start.End())
}

// UnmarshalXML implements xml.Unmarshaler.UnmarshalXML.
func (w *WptType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var e struct {
		Lat           float64         `xml:"lat,attr"`
		Lon           float64         `xml:"lon,attr"`
		Ele           float64         `xml:"ele"`
		Speed         float64         `xml:"speed"`
		Course        float64         `xml:"course"`
		Time          string          `xml:"time"`
		MagVar        float64         `xml:"magvar"`
		GeoidHeight   float64         `xml:"geoidheight"`
		Name          string          `xml:"name"`
		Cmt           string          `xml:"cmt"`
		Desc          string          `xml:"desc"`
		Src           string          `xml:"src"`
		Link          []*LinkType     `xml:"link"`
		Sym           string          `xml:"sym"`
		Type          string          `xml:"type"`
		Fix           string          `xml:"fix"`
		Sat           int             `xml:"sat"`
		HDOP          float64         `xml:"hdop"`
		VDOP          float64         `xml:"vdop"`
		PDOP          float64         `xml:"pdop"`
		AgeOfDGPSData float64         `xml:"ageofdgpsdata"`
		DGPSID        []int           `xml:"dgpsid"`
		Extensions    *ExtensionsType `xml:"extensions"`
	}
	if err := d.DecodeElement(&e, &start); err != nil {
		return err
	}
	wt := WptType{
		Lat:           e.Lat,
		Lon:           e.Lon,
		Ele:           e.Ele,
		Speed:         e.Speed,
		Course:        e.Course,
		MagVar:        e.MagVar,
		GeoidHeight:   e.GeoidHeight,
		Name:          e.Name,
		Cmt:           e.Cmt,
		Desc:          e.Desc,
		Src:           e.Src,
		Link:          e.Link,
		Sym:           e.Sym,
		Type:          e.Type,
		Fix:           e.Fix,
		Sat:           e.Sat,
		HDOP:          e.HDOP,
		VDOP:          e.VDOP,
		PDOP:          e.PDOP,
		AgeOfDGPSData: e.AgeOfDGPSData,
		DGPSID:        e.DGPSID,
		Extensions:    e.Extensions,
	}
	if e.Time != "" {
		t, err := parseTime(e.Time)
		if err != nil {
			return err
		}
		wt.Time = t
	}
	*w = wt
	return nil
}

func (w *WptType) appendFlatCoords(flatCoords []float64, layout geom.Layout) []float64 {
	switch layout {
	case geom.NoLayout:
		return flatCoords
	case geom.XY:
		return append(flatCoords, w.Lon, w.Lat)
	case geom.XYZ:
		return append(flatCoords, w.Lon, w.Lat, w.Ele)
	case geom.XYM:
		return append(flatCoords, w.Lon, w.Lat, TimeToM(w.Time))
	case geom.XYZM:
		return append(flatCoords, w.Lon, w.Lat, w.Ele, TimeToM(w.Time))
	default:
		flatCoords = append(flatCoords, w.Lon, w.Lat, w.Ele, TimeToM(w.Time))
		flatCoords = append(flatCoords, make([]float64, int(layout)-4)...)
		return flatCoords
	}
}

func MToTime(m float64) time.Time {
	if m == 0 {
		return time.Unix(0, 0)
	}
	return time.Unix(int64(m), int64(m*float64(time.Second))%int64(time.Second)).UTC()
}

func TimeToM(t time.Time) float64 {
	if t.IsZero() {
		return 0
	}
	return float64(t.UnixNano()) / float64(time.Second)
}

func emitIntElement(e *xml.Encoder, localName string, value int) error {
	return emitStringElement(e, localName, strconv.Itoa(value))
}

func emitStringElement(e *xml.Encoder, localName, value string) error {
	return e.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: localName}})
}

func maybeEmitFloatElement(e *xml.Encoder, localName string, value float64) error {
	if value == 0 {
		return nil
	}
	return emitStringElement(e, localName, strconv.FormatFloat(value, 'f', -1, 64))
}

func maybeEmitIntElement(e *xml.Encoder, localName string, value int) error {
	if value == 0 {
		return nil
	}
	return emitIntElement(e, localName, value)
}

func maybeEmitStringElement(e *xml.Encoder, localName, value string) error {
	if value == "" {
		return nil
	}
	return emitStringElement(e, localName, value)
}

func newWptTypes(g *geom.LineString) []*WptType {
	flatCoords := g.FlatCoords()
	layout := g.Layout()
	mIndex := layout.MIndex()
	zIndex := layout.ZIndex()
	stride := layout.Stride()
	wpts := make([]*WptType, g.NumCoords())
	start := 0
	for i := range wpts {
		wpt := &WptType{
			Lat: flatCoords[start+1],
			Lon: flatCoords[start],
		}
		if zIndex != -1 {
			wpt.Ele = flatCoords[start+zIndex]
		}
		if mIndex != -1 {
			wpt.Time = MToTime(flatCoords[start+mIndex])
		}
		start += stride
		wpts[i] = wpt
	}
	return wpts
}

func parseTime(value string) (time.Time, error) {
	firstErr := errNoTimeLayout
	for i, timeLayout := range timeLayouts {
		switch t, err := time.Parse(timeLayout, value); {
		case err == nil:
			return t, nil
		case i == 0:
			firstErr = err
		}
	}
	return time.Time{}, firstErr
}
