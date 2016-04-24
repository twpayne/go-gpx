// Package gpx provides convenience types for reading and writing GPX
// documents.
// See http://www.topografix.com/gpx.asp.
package gpx

// TODO DGPSID
// TODO Extensions

import (
	"encoding/xml"
	"strconv"
	"time"
)

const timeLayout = "2006-01-02T15:04:05.999999999Z"

type BoundsType struct {
	MinLat float64 `xml:"minlat,attr"`
	MinLon float64 `xml:"minlon,attr"`
	MaxLat float64 `xml:"maxlat,attr"`
	MaxLon float64 `xml:"maxlon,attr"`
}

type CopyrightType struct {
	Author  string `xml:"author,attr"`
	Year    int    `xml:"year,omitempty"`
	License string `xml:"license,omitempty"`
}

type ExtensionsType []byte

type GPXType struct {
	XMLName  string        `xml:"gpx"`
	Version  string        `xml:"version,attr"`
	Creator  string        `xml:"creator,attr"`
	Metadata *MetadataType `xml:"metadata,omitempty"`
	Wpt      []*WptType    `xml:"wpt,omitempty"`
	Rte      []*RteType    `xml:"rte,omitempty"`
	Trk      []*TrkType    `xml:"trk,omitempty"`
	// Extensions
}

type LinkType struct {
	HREF string `xml:"href,attr"`
	Text string `xml:"text,omitempty"`
	Type string `xml:"type,omitempty"`
}

type PersonType struct {
	Name  string    `xml:"name,omitempty"`
	Email string    `xml:"email,omitempty"`
	Link  *LinkType `xml:"link,omitempty"`
}

type MetadataType struct {
	Name      string         `xml:"name,omitempty"`
	Desc      string         `xml:"desc,omitempty"`
	Author    *PersonType    `xml:"author,omitempty"`
	Copyright *CopyrightType `xml:"copyright,omitempty"`
	Link      []*LinkType    `xml:"link,omitempty"`
	Time      time.Time      `xml:"time,omitempty"`
	Keywords  string         `xml:"keywords,omitempty"`
	Bounds    *BoundsType    `xml:"bounds,omitempty"`
	// Extensions
}

type RteType struct {
	Name   string      `xml:"name,omitempty"`
	Cmt    string      `xml:"cmt,omitempty"`
	Desc   string      `xml:"desc,omitempty"`
	Src    string      `xml:"src,omitempty"`
	Link   []*LinkType `xml:"link,omitempty"`
	Number int         `xml:"number,omitempty"`
	Type   string      `xml:"type,omitempty"`
	// Extensions
	RtePt []*WptType `xml:"rtept,omitempty"`
}

type TrkSegType struct {
	TrkPt []*WptType `xml:"trkpt,omitempty"`
	// Extensions
}

type TrkType struct {
	Name   string      `xml:"name,omitempty"`
	Cmt    string      `xml:"cmt,omitempty"`
	Desc   string      `xml:"desc,omitempty"`
	Src    string      `xml:"src,omitempty"`
	Link   []*LinkType `xml:"link,omitempty"`
	Number int         `xml:"number,omitempty"`
	Type   string      `xml:"type,omitempty"`
	// Extensions
	TrkSeg []*TrkSegType `xml:"trkseg,omitempty"`
}

type WptType struct {
	Lat          float64
	Lon          float64
	Ele          float64
	Time         time.Time
	MagVar       float64
	GeoidHeight  float64
	Name         string
	Cmt          string
	Desc         string
	Src          string
	Link         []*LinkType
	Sym          string
	Type         string
	Fix          string
	Sat          int
	HDOP         float64
	VDOP         float64
	PDOP         float64
	AgeOfGPSData float64
	// DGPSID
	// Extensions
}

func emitIntElement(e *xml.Encoder, localName string, value int) error {
	return emitStringElement(e, localName, strconv.FormatInt(int64(value), 10))
}

func emitStringElement(e *xml.Encoder, localName, value string) error {
	return e.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: localName}})
}

func maybeEmitFloatElement(e *xml.Encoder, localName string, value float64) error {
	if value == 0.0 {
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
	if !w.Time.IsZero() {
		if err := maybeEmitStringElement(e, "time", w.Time.UTC().Format(timeLayout)); err != nil {
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
	if err := maybeEmitFloatElement(e, "ageofgpsdata", w.AgeOfGPSData); err != nil {
		return err
	}
	// DGPSID
	// Extensions
	return e.EncodeToken(start.End())
}

func (g *GPXType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start = xml.StartElement{
		Name: xml.Name{Local: "gpx"},
		Attr: []xml.Attr{
			xml.Attr{
				Name:  xml.Name{Local: "version"},
				Value: g.Version,
			},
			xml.Attr{
				Name:  xml.Name{Local: "creator"},
				Value: g.Creator,
			},
			xml.Attr{
				Name:  xml.Name{Local: "xmlns:xsi"},
				Value: "http://www.w3.org/2001/XMLSchema-instance",
			},
			xml.Attr{
				Name:  xml.Name{Local: "xmlns"},
				Value: "http://www.topografix.com/GPX/1/0",
			},
			xml.Attr{
				Name:  xml.Name{Local: "xsi:schemaLocation"},
				Value: "http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd",
			},
		},
	}
	if err := e.EncodeToken(start); err != nil {
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

func (w *WptType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var e struct {
		Lat          float64     `xml:"lat,attr"`
		Lon          float64     `xml:"lon,attr"`
		Ele          float64     `xml:"ele"`
		Time         string      `xml:"time"`
		MagVar       float64     `xml:"magvar"`
		GeoidHeight  float64     `xml:"geoidheight"`
		Name         string      `xml:"name"`
		Cmt          string      `xml:"cmt"`
		Desc         string      `xml:"desc"`
		Src          string      `xml:"src"`
		Link         []*LinkType `xml:"link"`
		Sym          string      `xml:"sym"`
		Type         string      `xml:"type"`
		Fix          string      `xml:"fix"`
		Sat          int         `xml:"sat"`
		HDOP         float64     `xml:"hdop"`
		VDOP         float64     `xml:"vdop"`
		PDOP         float64     `xml:"pdop"`
		AgeOfGPSData float64     `xml:"ageofgpsdata"`
		// DGPSID
		// Extensions
	}
	if err := d.DecodeElement(&e, &start); err != nil {
		return err
	}
	*w = WptType{
		Lat:          e.Lat,
		Lon:          e.Lon,
		Ele:          e.Ele,
		MagVar:       e.MagVar,
		GeoidHeight:  e.GeoidHeight,
		Name:         e.Name,
		Cmt:          e.Cmt,
		Desc:         e.Desc,
		Src:          e.Src,
		Link:         e.Link,
		Sym:          e.Sym,
		Type:         e.Type,
		Fix:          e.Fix,
		Sat:          e.Sat,
		HDOP:         e.HDOP,
		VDOP:         e.VDOP,
		PDOP:         e.PDOP,
		AgeOfGPSData: e.AgeOfGPSData,
		// DGPSID
		// Extensions
	}
	if e.Time != "" {
		t, err := time.ParseInLocation(timeLayout, e.Time, time.UTC)
		if err != nil {
			return err
		}
		w.Time = t
	}
	return nil
}

var (
	_ xml.Marshaler   = &WptType{}
	_ xml.Unmarshaler = &WptType{}
)
