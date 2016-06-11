package gpx

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/d4l3k/messagediff"
	"github.com/twpayne/go-geom"
)

func ExampleRead() {
	r := bytes.NewBufferString("<gpx" +
		" version=\"1.0\"" +
		" creator=\"ExpertGPS 1.1 - http://www.topografix.com\"" +
		" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"" +
		" xmlns=\"http://www.topografix.com/GPX/1/0\"" +
		" xsi:schemaLocation=\"http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd\">\n" +
		"\t<wpt lat=\"42.438878\" lon=\"-71.119277\">\n" +
		"\t\t<ele>44.586548</ele>\n" +
		"\t\t<time>2001-11-28T21:05:28Z</time>\n" +
		"\t\t<name>5066</name>\n" +
		"\t\t<desc>5066</desc>\n" +
		"\t\t<sym>Crossing</sym>\n" +
		"\t\t<type>Crossing</type>\n" +
		"\t</wpt>\n" +
		"</gpx>")
	t, err := Read(r)
	if err != nil {
		fmt.Printf("err == %v", err)
		return
	}
	fmt.Printf("t.Wpt[0] == %+v", t.Wpt[0])
	// Output:
	// t.Wpt[0] == &{Lat:42.438878 Lon:-71.119277 Ele:44.586548 Time:2001-11-28 21:05:28 +0000 UTC MagVar:0 GeoidHeight:0 Name:5066 Cmt: Desc:5066 Src: Link:[] Sym:Crossing Type:Crossing Fix: Sat:0 HDOP:0 VDOP:0 PDOP:0 AgeOfGPSData:0 DGPSID:[] Extensions:<nil>}
}

func ExampleT_WriteIndent() {
	t := &T{
		Version: "1.0",
		Creator: "ExpertGPS 1.1 - http://www.topografix.com",
		Wpt: []*WptType{
			&WptType{
				Lat:  42.438878,
				Lon:  -71.119277,
				Ele:  44.586548,
				Time: time.Date(2001, 11, 28, 21, 5, 28, 0, time.UTC),
				Name: "5066",
				Desc: "5066",
				Sym:  "Crossing",
				Type: "Crossing",
			},
		},
	}
	if err := t.WriteIndent(os.Stdout, "", "  "); err != nil {
		fmt.Printf("err == %v", err)
	}
	// Output:
	// <gpx version="1.0" creator="ExpertGPS 1.1 - http://www.topografix.com" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://www.topografix.com/GPX/1/0" xsi:schemaLocation="http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd">
	//   <wpt lat="42.438878" lon="-71.119277">
	//     <ele>44.586548</ele>
	//     <time>2001-11-28T21:05:28Z</time>
	//     <name>5066</name>
	//     <desc>5066</desc>
	//     <sym>Crossing</sym>
	//     <type>Crossing</type>
	//   </wpt>
	// </gpx>
}

func TestWpt(t *testing.T) {
	for _, tc := range []struct {
		data          string
		wpt           *WptType
		layout        geom.Layout
		g             *geom.Point
		noTestMarshal bool
	}{
		{
			data: "<wpt lat=\"42.438878\" lon=\"-71.119277\"></wpt>",
			wpt: &WptType{
				Lat: 42.438878,
				Lon: -71.119277,
			},
			layout: geom.XY,
			g:      geom.NewPoint(geom.XY).MustSetCoords([]float64{-71.119277, 42.438878}),
		},
		{
			data: "<wpt lat=\"42.438878\" lon=\"-71.119277\">\n" +
				"\t<ele>44.586548</ele>\n" +
				"</wpt>",
			wpt: &WptType{
				Lat: 42.438878,
				Lon: -71.119277,
				Ele: 44.586548,
			},
			layout: geom.XYZ,
			g:      geom.NewPoint(geom.XYZ).MustSetCoords([]float64{-71.119277, 42.438878, 44.586548}),
		},
		{
			data: "<wpt lat=\"42.438878\" lon=\"-71.119277\">\n" +
				"\t<time>2001-11-28T21:05:28Z</time>\n" +
				"</wpt>",
			wpt: &WptType{
				Lat:  42.438878,
				Lon:  -71.119277,
				Time: time.Date(2001, 11, 28, 21, 5, 28, 0, time.UTC),
			},
			layout: geom.XYM,
			g:      geom.NewPoint(geom.XYM).MustSetCoords([]float64{-71.119277, 42.438878, 1006981528}),
		},
		{
			data: "<wpt lat=\"42.438878\" lon=\"-71.119277\">\n" +
				"\t<ele>44.586548</ele>\n" +
				"\t<time>2001-11-28T21:05:28Z</time>\n" +
				"\t<name>5066</name>\n" +
				"\t<desc><![CDATA[5066]]></desc>\n" +
				"\t<sym>Crossing</sym>\n" +
				"\t<type><![CDATA[Crossing]]></type>\n" +
				"</wpt>\n",
			wpt: &WptType{
				Lat:  42.438878,
				Lon:  -71.119277,
				Ele:  44.586548,
				Time: time.Date(2001, 11, 28, 21, 5, 28, 0, time.UTC),
				Name: "5066",
				Desc: "5066",
				Sym:  "Crossing",
				Type: "Crossing",
			},
			layout:        geom.XYZM,
			g:             geom.NewPoint(geom.XYZM).MustSetCoords([]float64{-71.119277, 42.438878, 44.586548, 1006981528}),
			noTestMarshal: true,
		},
		{
			data: "<wpt lat=\"42.438878\" lon=\"-71.119277\">\n" +
				"\t<ele>44.586548</ele>\n" +
				"\t<time>2001-11-28T21:05:28Z</time>\n" +
				"\t<magvar>1.1</magvar>\n" +
				"\t<geoidheight>2.2</geoidheight>\n" +
				"\t<name>5066</name>\n" +
				"\t<cmt>Comment</cmt>\n" +
				"\t<desc>5066</desc>\n" +
				"\t<src>Source</src>\n" +
				"\t<link href=\"http://example.com\">\n" +
				"\t\t<text>Text</text>\n" +
				"\t\t<type>Type</type>\n" +
				"\t</link>\n" +
				"\t<sym>Crossing</sym>\n" +
				"\t<type>Crossing</type>\n" +
				"\t<fix>3d</fix>\n" +
				"\t<sat>3</sat>\n" +
				"\t<hdop>4.4</hdop>\n" +
				"\t<vdop>5.5</vdop>\n" +
				"\t<pdop>6.6</pdop>\n" +
				"\t<ageofgpsdata>7.7</ageofgpsdata>\n" +
				"\t<dgpsid>8</dgpsid>\n" +
				"</wpt>",
			wpt: &WptType{
				Lat:         42.438878,
				Lon:         -71.119277,
				Ele:         44.586548,
				MagVar:      1.1,
				Time:        time.Date(2001, 11, 28, 21, 5, 28, 0, time.UTC),
				GeoidHeight: 2.2,
				Name:        "5066",
				Cmt:         "Comment",
				Desc:        "5066",
				Src:         "Source",
				Link: []*LinkType{
					&LinkType{
						HREF: "http://example.com",
						Text: "Text",
						Type: "Type",
					},
				},
				Sym:          "Crossing",
				Type:         "Crossing",
				Fix:          "3d",
				Sat:          3,
				HDOP:         4.4,
				VDOP:         5.5,
				PDOP:         6.6,
				AgeOfGPSData: 7.7,
				DGPSID:       []int{8},
			},
			layout: geom.XYZM,
			g:      geom.NewPoint(geom.XYZM).MustSetCoords([]float64{-71.119277, 42.438878, 44.586548, 1006981528}),
		},
	} {
		var gotWpt WptType
		if err := xml.Unmarshal([]byte(tc.data), &gotWpt); err != nil {
			t.Errorf("xml.Unmarshal([]byte(%q), &gotWpt) == %v, want nil", tc.data, err)
		}
		if diff, equal := messagediff.PrettyDiff(tc.wpt, &gotWpt); !equal {
			t.Errorf("xml.Unmarshal([]byte(%q), &gotWpt); got == %#v, diff\n%s", tc.data, gotWpt, diff)
		}
		if tc.layout != geom.NoLayout {
			gotG := tc.wpt.Geom(tc.layout)
			if diff, equal := messagediff.PrettyDiff(tc.g, gotG); !equal {
				t.Errorf("%#v.Geom() == %#v, diff\n%s", tc.wpt, gotG, diff)
			}
		}
		if !tc.noTestMarshal {
			var b bytes.Buffer
			e := xml.NewEncoder(&b)
			e.Indent("", "\t")
			start := xml.StartElement{Name: xml.Name{Local: "wpt"}}
			if err := e.EncodeElement(tc.wpt, start); err != nil {
				t.Errorf("e.EncodeElement(%#v, %#v) == _, %v, want _, nil", tc.wpt, start, err)
			}
			if diff, equal := messagediff.PrettyDiff(strings.Split(tc.data, "\n"), strings.Split(b.String(), "\n")); !equal {
				t.Errorf("xml.Marshal(%#v) == %q, nil, want %q, diff\n%s", tc.wpt, b.String(), tc.data, diff)
			}
		}
	}
}

func TestRte(t *testing.T) {
	for _, tc := range []struct {
		data          string
		rte           *RteType
		layout        geom.Layout
		g             *geom.LineString
		noTestMarshal bool
	}{
		{
			data: "<rte>\n" +
				"\t<name>BELLEVUE</name>\n" +
				"\t<desc>Bike Loop Bellevue</desc>\n" +
				"\t<number>1</number>\n" +
				"\t<rtept lat=\"42.43095\" lon=\"-71.107628\">\n" +
				"\t\t<ele>23.4696</ele>\n" +
				"\t\t<time>2001-06-02T00:18:15Z</time>\n" +
				"\t\t<name>BELLEVUE</name>\n" +
				"\t\t<cmt>BELLEVUE</cmt>\n" +
				"\t\t<desc>Bellevue Parking Lot</desc>\n" +
				"\t\t<sym>Parking Area</sym>\n" +
				"\t\t<type>Parking</type>\n" +
				"\t</rtept>\n" +
				"\t<rtept lat=\"42.43124\" lon=\"-71.109236\">\n" +
				"\t\t<ele>26.56189</ele>\n" +
				"\t\t<time>2001-11-07T23:53:41Z</time>\n" +
				"\t\t<name>GATE6</name>\n" +
				"\t\t<desc>Gate 6</desc>\n" +
				"\t\t<sym>Trailhead</sym>\n" +
				"\t\t<type>Trail Head</type>\n" +
				"\t</rtept>\n" +
				"</rte>",
			rte: &RteType{
				Name:   "BELLEVUE",
				Desc:   "Bike Loop Bellevue",
				Number: 1,
				RtePt: []*WptType{
					&WptType{
						Lat:  42.43095,
						Lon:  -71.107628,
						Ele:  23.4696,
						Time: time.Date(2001, 6, 2, 0, 18, 15, 0, time.UTC),
						Name: "BELLEVUE",
						Cmt:  "BELLEVUE",
						Desc: "Bellevue Parking Lot",
						Sym:  "Parking Area",
						Type: "Parking",
					},
					&WptType{
						Lat:  42.43124,
						Lon:  -71.109236,
						Ele:  26.56189,
						Time: time.Date(2001, 11, 7, 23, 53, 41, 0, time.UTC),
						Name: "GATE6",
						Desc: "Gate 6",
						Sym:  "Trailhead",
						Type: "Trail Head",
					},
				},
			},
			layout: geom.XYZM,
			g: geom.NewLineString(geom.XYZM).MustSetCoords(
				[]geom.Coord{
					geom.Coord{-71.107628, 42.43095, 23.4696, 991441095},
					geom.Coord{-71.109236, 42.43124, 26.56189, 1005177221},
				},
			),
		},
	} {
		var gotRte RteType
		if err := xml.Unmarshal([]byte(tc.data), &gotRte); err != nil {
			t.Errorf("xml.Unmarshal([]byte(%q), &gotRte) == %v, want nil", tc.data, err)
		}
		if diff, equal := messagediff.PrettyDiff(tc.rte, &gotRte); !equal {
			t.Errorf("xml.Unmarshal([]byte(%q), &gotRte); got == %#v, diff\n%s", tc.data, gotRte, diff)
		}
		if tc.layout != geom.NoLayout {
			gotG := tc.rte.Geom(tc.layout)
			if diff, equal := messagediff.PrettyDiff(tc.g, gotG); !equal {
				t.Errorf("%#v.Geom() == %#v, diff\n%s", tc.rte, gotG, diff)
			}
		}
		if !tc.noTestMarshal {
			var b bytes.Buffer
			e := xml.NewEncoder(&b)
			e.Indent("", "\t")
			start := xml.StartElement{Name: xml.Name{Local: "rte"}}
			if err := e.EncodeElement(tc.rte, start); err != nil {
				t.Errorf("e.EncodeElement(%#v, %#v) == _, %v, want _, nil", tc.rte, start, err)
			}
			if diff, equal := messagediff.PrettyDiff(strings.Split(tc.data, "\n"), strings.Split(b.String(), "\n")); !equal {
				t.Errorf("xml.Marshal(%#v) == %q, nil, want %q, diff\n%s", tc.rte, b.String(), tc.data, diff)
			}
		}
	}
}

func TestTrk(t *testing.T) {
	for _, tc := range []struct {
		data          string
		trk           *TrkType
		layout        geom.Layout
		g             *geom.MultiLineString
		noTestMarshal bool
	}{
		{
			data: "<trk>\n" +
				"\t<trkseg>\n" +
				"\t\t<trkpt lat=\"47.644548\" lon=\"-122.326897\">\n" +
				"\t\t\t<ele>4.46</ele>\n" +
				"\t\t\t<time>2009-10-17T18:37:26Z</time>\n" +
				"\t\t</trkpt>\n" +
				"\t\t<trkpt lat=\"47.644548\" lon=\"-122.326897\">\n" +
				"\t\t\t<ele>4.94</ele>\n" +
				"\t\t\t<time>2009-10-17T18:37:31Z</time>\n" +
				"\t\t</trkpt>\n" +
				"\t\t<trkpt lat=\"47.644548\" lon=\"-122.326897\">\n" +
				"\t\t\t<ele>6.87</ele>\n" +
				"\t\t\t<time>2009-10-17T18:37:34Z</time>\n" +
				"\t\t</trkpt>\n" +
				"\t</trkseg>\n" +
				"</trk>",
			trk: &TrkType{
				TrkSeg: []*TrkSegType{
					&TrkSegType{
						TrkPt: []*WptType{
							&WptType{
								Lat:  47.644548,
								Lon:  -122.326897,
								Ele:  4.46,
								Time: time.Date(2009, 10, 17, 18, 37, 26, 0, time.UTC),
							},
							&WptType{
								Lat:  47.644548,
								Lon:  -122.326897,
								Ele:  4.94,
								Time: time.Date(2009, 10, 17, 18, 37, 31, 0, time.UTC),
							},
							&WptType{
								Lat:  47.644548,
								Lon:  -122.326897,
								Ele:  6.87,
								Time: time.Date(2009, 10, 17, 18, 37, 34, 0, time.UTC),
							},
						},
					},
				},
			},
			layout: geom.XYZM,
			g: geom.NewMultiLineString(geom.XYZM).MustSetCoords(
				[][]geom.Coord{
					[]geom.Coord{
						geom.Coord{-122.326897, 47.644548, 4.46, 1255804646},
						geom.Coord{-122.326897, 47.644548, 4.94, 1255804651},
						geom.Coord{-122.326897, 47.644548, 6.87, 1255804654},
					},
				},
			),
		},
	} {
		var gotTrk TrkType
		if err := xml.Unmarshal([]byte(tc.data), &gotTrk); err != nil {
			t.Errorf("xml.Unmarshal([]byte(%q), &gotTrk) == %v, want nil", tc.data, err)
		}
		if diff, equal := messagediff.PrettyDiff(tc.trk, &gotTrk); !equal {
			t.Errorf("xml.Unmarshal([]byte(%q), &gotTrk); got == %#v, diff\n%s", tc.data, gotTrk, diff)
		}
		if tc.layout != geom.NoLayout {
			gotG := tc.trk.Geom(tc.layout)
			if diff, equal := messagediff.PrettyDiff(tc.g, gotG); !equal {
				t.Errorf("%#v.Geom() == %#v, diff\n%s", tc.trk, gotG, diff)
			}
		}
		if !tc.noTestMarshal {
			var b bytes.Buffer
			e := xml.NewEncoder(&b)
			e.Indent("", "\t")
			start := xml.StartElement{Name: xml.Name{Local: "trk"}}
			if err := e.EncodeElement(tc.trk, start); err != nil {
				t.Errorf("e.EncodeElement(%#v, %#v) == _, %v, want _, nil", tc.trk, start, err)
			}
			if diff, equal := messagediff.PrettyDiff(strings.Split(tc.data, "\n"), strings.Split(b.String(), "\n")); !equal {
				t.Errorf("xml.Marshal(%#v) == %q, nil, want %q, diff\n%s", tc.trk, b.String(), tc.data, diff)
			}
		}
	}
}

func TestRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		data string
		gpx  *T
	}{
		{
			data: "<gpx" +
				" version=\"1.0\"" +
				" creator=\"ExpertGPS 1.1 - http://www.topografix.com\"" +
				" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"" +
				" xmlns=\"http://www.topografix.com/GPX/1/0\"" +
				" xsi:schemaLocation=\"http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd\">" +
				"</gpx>",
			gpx: &T{
				Version: "1.0",
				Creator: "ExpertGPS 1.1 - http://www.topografix.com",
			},
		},
		{
			data: "<gpx" +
				" version=\"1.0\"" +
				" creator=\"ExpertGPS 1.1 - http://www.topografix.com\"" +
				" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"" +
				" xmlns=\"http://www.topografix.com/GPX/1/0\"" +
				" xsi:schemaLocation=\"http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd\">\n" +
				"\t<wpt lat=\"42.438878\" lon=\"-71.119277\">\n" +
				"\t\t<ele>44.586548</ele>\n" +
				"\t\t<time>2001-11-28T21:05:28Z</time>\n" +
				"\t\t<name>5066</name>\n" +
				"\t\t<desc>5066</desc>\n" +
				"\t\t<sym>Crossing</sym>\n" +
				"\t\t<type>Crossing</type>\n" +
				"\t</wpt>\n" +
				"</gpx>",
			gpx: &T{
				Version: "1.0",
				Creator: "ExpertGPS 1.1 - http://www.topografix.com",
				Wpt: []*WptType{
					&WptType{
						Lat:  42.438878,
						Lon:  -71.119277,
						Ele:  44.586548,
						Time: time.Date(2001, 11, 28, 21, 5, 28, 0, time.UTC),
						Name: "5066",
						Desc: "5066",
						Sym:  "Crossing",
						Type: "Crossing",
					},
				},
			},
		},
		{
			data: "<gpx" +
				" version=\"1.0\"" +
				" creator=\"ExpertGPS 1.1 - http://www.topografix.com\"" +
				" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\"" +
				" xmlns=\"http://www.topografix.com/GPX/1/0\"" +
				" xsi:schemaLocation=\"http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd\">\n" +
				"\t<rte>\n" +
				"\t\t<name>BELLEVUE</name>\n" +
				"\t\t<desc>Bike Loop Bellevue</desc>\n" +
				"\t\t<number>1</number>\n" +
				"\t\t<rtept lat=\"42.43095\" lon=\"-71.107628\">\n" +
				"\t\t\t<ele>23.4696</ele>\n" +
				"\t\t\t<time>2001-06-02T00:18:15Z</time>\n" +
				"\t\t\t<name>BELLEVUE</name>\n" +
				"\t\t\t<cmt>BELLEVUE</cmt>\n" +
				"\t\t\t<desc>Bellevue Parking Lot</desc>\n" +
				"\t\t\t<sym>Parking Area</sym>\n" +
				"\t\t\t<type>Parking</type>\n" +
				"\t\t</rtept>\n" +
				"\t\t<rtept lat=\"42.43124\" lon=\"-71.109236\">\n" +
				"\t\t\t<ele>26.56189</ele>\n" +
				"\t\t\t<time>2001-11-07T23:53:41Z</time>\n" +
				"\t\t\t<name>GATE6</name>\n" +
				"\t\t\t<desc>Gate 6</desc>\n" +
				"\t\t\t<sym>Trailhead</sym>\n" +
				"\t\t\t<type>Trail Head</type>\n" +
				"\t\t</rtept>\n" +
				"\t</rte>\n" +
				"</gpx>",
			gpx: &T{
				Version: "1.0",
				Creator: "ExpertGPS 1.1 - http://www.topografix.com",
				Rte: []*RteType{
					&RteType{
						Name:   "BELLEVUE",
						Desc:   "Bike Loop Bellevue",
						Number: 1,
						RtePt: []*WptType{
							&WptType{
								Lat:  42.43095,
								Lon:  -71.107628,
								Ele:  23.4696,
								Time: time.Date(2001, 6, 2, 0, 18, 15, 0, time.UTC),
								Name: "BELLEVUE",
								Cmt:  "BELLEVUE",
								Desc: "Bellevue Parking Lot",
								Sym:  "Parking Area",
								Type: "Parking",
							},
							&WptType{
								Lat:  42.43124,
								Lon:  -71.109236,
								Ele:  26.56189,
								Time: time.Date(2001, 11, 7, 23, 53, 41, 0, time.UTC),
								Name: "GATE6",
								Desc: "Gate 6",
								Sym:  "Trailhead",
								Type: "Trail Head",
							},
						},
					},
				},
			},
		},
	} {
		got, err := Read(bytes.NewBufferString(tc.data))
		if err != nil {
			t.Errorf("Read(bytes.NewBuffer(%v)) == _, %v, want _, nil", tc.data, err)
		}
		if diff, equal := messagediff.PrettyDiff(tc.gpx, got); !equal {
			t.Errorf("xml.Unmarshal([]byte(%q), &got); got == %#v, diff\n%s", tc.data, got, diff)
		}
		b := &bytes.Buffer{}
		if err := tc.gpx.WriteIndent(b, "", "\t"); err != nil {
			t.Errorf("%#v.WriteIndent(...) == %v, want nil", tc.gpx, err)
		}
		if diff, equal := messagediff.PrettyDiff(strings.Split(tc.data, "\n"), strings.Split(b.String(), "\n")); !equal {
			t.Errorf("xml.Marshal(%#v) ==\n%s\nwant\n%s\ndiff\n%s", tc.gpx, b.String(), tc.data, diff)
		}
	}
}

func TestTime(t *testing.T) {
	for _, tc := range []struct {
		t time.Time
		m float64
	}{
		{
			t: time.Unix(0, 0),
			m: 0,
		},
		{
			t: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			m: 946684800,
		},
		{
			t: time.Date(2006, 1, 2, 15, 4, 5, 500000000, time.UTC),
			m: 1136214245.5,
		},
	} {
		if gotM := timeToM(tc.t); gotM != tc.m {
			t.Errorf("timeToM(%v) == %v, want %v", tc.t, gotM, tc.m)
		}
		if gotT := mToTime(tc.m); gotT != tc.t {
			t.Errorf("mToTime(%v) == %v, want %v", tc.m, gotT, tc.t)
		}
	}
}

func TestParseFellsLoop(t *testing.T) {
	if _, err := Read(bytes.NewBuffer(fellsLoopData)); err != nil {
		t.Errorf("Read(bytes.NewBuffer(fellsLoopData)) == _, %v, want _, nil", err)
	}
}

var fellsLoopData = []byte(`<?xml version="1.0"?>
<gpx
 version="1.0"
 creator="ExpertGPS 1.1 - http://www.topografix.com"
 xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
 xmlns="http://www.topografix.com/GPX/1/0"
 xsi:schemaLocation="http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd">
<time>2002-02-27T17:18:33Z</time>
<bounds minlat="42.401051" minlon="-71.126602" maxlat="42.468655" maxlon="-71.102973"/>
<wpt lat="42.438878" lon="-71.119277">
 <ele>44.586548</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>5066</name>
 <desc><![CDATA[5066]]></desc>
 <sym>Crossing</sym>
 <type><![CDATA[Crossing]]></type>
</wpt>
<wpt lat="42.439227" lon="-71.119689">
 <ele>57.607200</ele>
 <time>2001-06-02T03:26:55Z</time>
 <name>5067</name>
 <desc><![CDATA[5067]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.438917" lon="-71.116146">
 <ele>44.826904</ele>
 <time>2001-11-16T23:03:38Z</time>
 <name>5096</name>
 <desc><![CDATA[5096]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.443904" lon="-71.122044">
 <ele>50.594727</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>5142</name>
 <desc><![CDATA[5142]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.447298" lon="-71.121447">
 <ele>127.711200</ele>
 <time>2001-06-02T03:26:58Z</time>
 <name>5156</name>
 <desc><![CDATA[5156]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.454873" lon="-71.125094">
 <ele>96.926400</ele>
 <time>2001-06-02T03:26:59Z</time>
 <name>5224</name>
 <desc><![CDATA[5224]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.459079" lon="-71.124988">
 <ele>82.600800</ele>
 <time>2001-06-02T03:26:59Z</time>
 <name>5229</name>
 <desc><![CDATA[5229]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.456979" lon="-71.124474">
 <ele>82.905600</ele>
 <time>2001-06-02T03:26:59Z</time>
 <name>5237</name>
 <desc><![CDATA[5237]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.454401" lon="-71.120990">
 <ele>66.696655</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>5254</name>
 <desc><![CDATA[5254]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.451442" lon="-71.121746">
 <ele>74.627442</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>5258</name>
 <desc><![CDATA[5258]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.454404" lon="-71.120660">
 <ele>65.254761</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>5264</name>
 <desc><![CDATA[5264]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.457761" lon="-71.121045">
 <ele>77.419200</ele>
 <time>2001-06-02T03:27:00Z</time>
 <name>526708</name>
 <desc><![CDATA[526708]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.457089" lon="-71.120313">
 <ele>74.676000</ele>
 <time>2001-06-02T03:27:00Z</time>
 <name>526750</name>
 <desc><![CDATA[526750]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.456592" lon="-71.119676">
 <ele>78.713135</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>527614</name>
 <desc><![CDATA[527614]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.456252" lon="-71.119356">
 <ele>78.713135</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>527631</name>
 <desc><![CDATA[527631]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.458148" lon="-71.119135">
 <ele>68.275200</ele>
 <time>2001-06-02T03:27:00Z</time>
 <name>5278</name>
 <desc><![CDATA[5278]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.459377" lon="-71.117693">
 <ele>64.008000</ele>
 <time>2001-06-02T03:27:01Z</time>
 <name>5289</name>
 <desc><![CDATA[5289]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.464183" lon="-71.119828">
 <ele>52.997925</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>5374FIRE</name>
 <desc><![CDATA[5374FIRE]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.465650" lon="-71.119399">
 <ele>56.388000</ele>
 <time>2001-06-02T03:27:02Z</time>
 <name>5376</name>
 <desc><![CDATA[5376]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.439018" lon="-71.114456">
 <ele>56.388000</ele>
 <time>2001-06-02T03:26:55Z</time>
 <name>6006</name>
 <desc><![CDATA[600698]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.438594" lon="-71.114803">
 <ele>46.028564</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>6006BLUE</name>
 <desc><![CDATA[6006BLUE]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.436757" lon="-71.113223">
 <ele>37.616943</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>6014MEADOW</name>
 <desc><![CDATA[6014MEADOW]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.441754" lon="-71.113220">
 <ele>56.388000</ele>
 <time>2001-06-02T03:26:55Z</time>
 <name>6029</name>
 <desc><![CDATA[6029]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.436243" lon="-71.109075">
 <ele>50.292000</ele>
 <time>2001-06-02T03:27:05Z</time>
 <name>6053</name>
 <desc><![CDATA[6053]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.439250" lon="-71.107500">
 <ele>25.603200</ele>
 <time>2001-06-02T03:26:57Z</time>
 <name>6066</name>
 <desc><![CDATA[6066]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.439764" lon="-71.107582">
 <ele>34.442400</ele>
 <time>2001-06-02T03:26:57Z</time>
 <name>6067</name>
 <desc><![CDATA[6067]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.434766" lon="-71.105874">
 <ele>30.480000</ele>
 <time>2001-06-02T03:26:57Z</time>
 <name>6071</name>
 <desc><![CDATA[6071]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.433304" lon="-71.106599">
 <ele>15.240000</ele>
 <time>2001-06-02T03:26:56Z</time>
 <name>6073</name>
 <desc><![CDATA[6073]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.437338" lon="-71.104772">
 <ele>37.795200</ele>
 <time>2001-06-02T03:26:57Z</time>
 <name>6084</name>
 <desc><![CDATA[6084]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.442196" lon="-71.110975">
 <ele>64.008000</ele>
 <time>2001-06-02T03:26:55Z</time>
 <name>6130</name>
 <desc><![CDATA[6130]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.442981" lon="-71.111441">
 <ele>64.008000</ele>
 <time>2001-06-02T03:26:58Z</time>
 <name>6131</name>
 <desc><![CDATA[6131]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.444773" lon="-71.108882">
 <ele>62.788800</ele>
 <time>2001-06-02T03:27:05Z</time>
 <name>6153</name>
 <desc><![CDATA[6153]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.443592" lon="-71.106301">
 <ele>55.473600</ele>
 <time>2001-06-02T03:27:05Z</time>
 <name>6171</name>
 <desc><![CDATA[6171]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.447804" lon="-71.106624">
 <ele>62.484000</ele>
 <time>2001-06-02T03:27:04Z</time>
 <name>6176</name>
 <desc><![CDATA[6176]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.448448" lon="-71.106158">
 <ele>62.179200</ele>
 <time>2001-06-02T03:27:04Z</time>
 <name>6177</name>
 <desc><![CDATA[6177]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.453415" lon="-71.106783">
 <ele>69.799200</ele>
 <time>2001-06-02T03:26:55Z</time>
 <name>6272</name>
 <desc><![CDATA[6272]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.453434" lon="-71.107253">
 <ele>73.152000</ele>
 <time>2001-06-02T03:26:56Z</time>
 <name>6272</name>
 <desc><![CDATA[6272]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.458298" lon="-71.106771">
 <ele>70.104000</ele>
 <time>2001-06-02T03:27:04Z</time>
 <name>6278</name>
 <desc><![CDATA[6278]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.451430" lon="-71.105413">
 <ele>57.564209</ele>
 <time>2001-11-16T23:03:38Z</time>
 <name>6280</name>
 <desc><![CDATA[6280]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.453845" lon="-71.105206">
 <ele>66.696655</ele>
 <time>2001-11-16T23:03:38Z</time>
 <name>6283</name>
 <desc><![CDATA[6283]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.459986" lon="-71.106170">
 <ele>72.945191</ele>
 <time>2001-11-16T23:03:38Z</time>
 <name>6289</name>
 <desc><![CDATA[6289]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.457616" lon="-71.105116">
 <ele>72.847200</ele>
 <time>2001-06-02T03:27:04Z</time>
 <name>6297</name>
 <desc><![CDATA[6297]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.467110" lon="-71.113574">
 <ele>53.644800</ele>
 <time>2001-06-02T03:27:02Z</time>
 <name>6328</name>
 <desc><![CDATA[6328]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.464202" lon="-71.109863">
 <ele>43.891200</ele>
 <time>2001-06-02T03:27:03Z</time>
 <name>6354</name>
 <desc><![CDATA[6354]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.466459" lon="-71.110067">
 <ele>48.768000</ele>
 <time>2001-06-02T03:27:02Z</time>
 <name>635722</name>
 <desc><![CDATA[635722]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.466557" lon="-71.109410">
 <ele>49.072800</ele>
 <time>2001-06-02T03:27:02Z</time>
 <name>635783</name>
 <desc><![CDATA[635783]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.463495" lon="-71.107117">
 <ele>62.484000</ele>
 <time>2001-06-02T03:27:03Z</time>
 <name>6373</name>
 <desc><![CDATA[6373]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.401051" lon="-71.110241">
 <ele>3.962400</ele>
 <time>2001-06-02T03:26:56Z</time>
 <name>6634</name>
 <desc><![CDATA[6634]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.432621" lon="-71.106532">
 <ele>13.411200</ele>
 <time>2001-06-02T03:26:56Z</time>
 <name>6979</name>
 <desc><![CDATA[6979]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.431033" lon="-71.107883">
 <ele>34.012085</ele>
 <time>2001-11-16T23:03:38Z</time>
 <name>6997</name>
 <desc><![CDATA[6997]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</wpt>
<wpt lat="42.465687" lon="-71.107360">
 <ele>87.782400</ele>
 <time>2001-06-02T03:27:03Z</time>
 <name>BEAR HILL</name>
 <cmt>BEAR HILL TOWER</cmt>
 <desc><![CDATA[Bear Hill Tower]]></desc>
 <sym>Tall Tower</sym>
 <type><![CDATA[Tower]]></type>
</wpt>
<wpt lat="42.430950" lon="-71.107628">
 <ele>23.469600</ele>
 <time>2001-06-02T00:18:15Z</time>
 <name>BELLEVUE</name>
 <cmt>BELLEVUE</cmt>
 <desc><![CDATA[Bellevue Parking Lot]]></desc>
 <sym>Parking Area</sym>
 <type><![CDATA[Parking]]></type>
</wpt>
<wpt lat="42.438666" lon="-71.114079">
 <ele>43.384766</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>6016</name>
 <desc><![CDATA[Bike Loop Connector]]></desc>
 <sym>Trailhead</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.456469" lon="-71.124651">
 <ele>89.916000</ele>
 <time>2001-06-02T03:26:59Z</time>
 <name>5236BRIDGE</name>
 <desc><![CDATA[Bridge]]></desc>
 <sym>Bridge</sym>
 <type><![CDATA[Bridge]]></type>
</wpt>
<wpt lat="42.465759" lon="-71.119815">
 <ele>55.473600</ele>
 <time>2001-06-02T03:27:01Z</time>
 <name>5376BRIDGE</name>
 <desc><![CDATA[Bridge]]></desc>
 <sym>Bridge</sym>
 <type><![CDATA[Bridge]]></type>
</wpt>
<wpt lat="42.442993" lon="-71.105878">
 <ele>52.730400</ele>
 <time>2001-06-02T03:27:05Z</time>
 <name>6181CROSS</name>
 <desc><![CDATA[Crossing]]></desc>
 <sym>Crossing</sym>
 <type><![CDATA[Crossing]]></type>
</wpt>
<wpt lat="42.435472" lon="-71.109664">
 <ele>45.110400</ele>
 <time>2001-06-02T03:27:05Z</time>
 <name>6042CROSS</name>
 <desc><![CDATA[Crossing]]></desc>
 <sym>Crossing</sym>
 <type><![CDATA[Crossing]]></type>
</wpt>
<wpt lat="42.458516" lon="-71.103646">
 <name>DARKHOLLPO</name>
 <desc><![CDATA[Dark Hollow Pond]]></desc>
 <sym>Fishing Area</sym>
</wpt>
<wpt lat="42.443109" lon="-71.112675">
 <ele>56.083200</ele>
 <time>2001-06-02T03:26:57Z</time>
 <name>6121DEAD</name>
 <desc><![CDATA[Dead End]]></desc>
 <sym>Danger Area</sym>
 <type><![CDATA[Dead End]]></type>
</wpt>
<wpt lat="42.449866" lon="-71.119298">
 <ele>117.043200</ele>
 <time>2001-06-02T03:26:59Z</time>
 <name>5179DEAD</name>
 <desc><![CDATA[Dead End]]></desc>
 <sym>Danger Area</sym>
 <type><![CDATA[Dead End]]></type>
</wpt>
<wpt lat="42.459629" lon="-71.116524">
 <ele>69.494400</ele>
 <time>2001-06-02T03:27:01Z</time>
 <name>5299DEAD</name>
 <desc><![CDATA[Dead End]]></desc>
 <sym>Danger Area</sym>
 <type><![CDATA[Dead End]]></type>
</wpt>
<wpt lat="42.465485" lon="-71.119148">
 <ele>56.997600</ele>
 <time>2001-06-02T03:27:02Z</time>
 <name>5376DEAD</name>
 <desc><![CDATA[Dead End]]></desc>
 <sym>Danger Area</sym>
 <type><![CDATA[Dead End]]></type>
</wpt>
<wpt lat="42.462776" lon="-71.109986">
 <ele>46.939200</ele>
 <time>2001-06-02T03:27:03Z</time>
 <name>6353DEAD</name>
 <desc><![CDATA[Dead End]]></desc>
 <sym>Danger Area</sym>
 <type><![CDATA[Dead End]]></type>
</wpt>
<wpt lat="42.446793" lon="-71.108784">
 <ele>61.264800</ele>
 <time>2001-06-02T03:27:04Z</time>
 <name>6155DEAD</name>
 <desc><![CDATA[Dead End]]></desc>
 <sym>Danger Area</sym>
 <type><![CDATA[Dead End]]></type>
</wpt>
<wpt lat="42.451204" lon="-71.126602">
 <ele>110.947200</ele>
 <time>2001-06-02T03:26:59Z</time>
 <name>GATE14</name>
 <desc><![CDATA[Gate 14]]></desc>
 <sym>Truck Stop</sym>
 <type><![CDATA[Road]]></type>
</wpt>
<wpt lat="42.458499" lon="-71.122078">
 <ele>77.724000</ele>
 <time>2001-06-02T03:27:00Z</time>
 <name>GATE16</name>
 <desc><![CDATA[Gate 16]]></desc>
 <sym>Truck Stop</sym>
 <type><![CDATA[Road]]></type>
</wpt>
<wpt lat="42.459376" lon="-71.119238">
 <ele>65.836800</ele>
 <time>2001-06-02T03:27:01Z</time>
 <name>GATE17</name>
 <desc><![CDATA[Gate 17]]></desc>
 <sym>Truck Stop</sym>
 <type><![CDATA[Road]]></type>
</wpt>
<wpt lat="42.466353" lon="-71.119240">
 <ele>57.302400</ele>
 <time>2001-06-02T03:27:02Z</time>
 <name>GATE19</name>
 <desc><![CDATA[Gate 19]]></desc>
 <sym>Truck Stop</sym>
 <type><![CDATA[Road]]></type>
</wpt>
<wpt lat="42.468655" lon="-71.107697">
 <ele>49.377600</ele>
 <time>2001-06-02T03:27:03Z</time>
 <name>GATE21</name>
 <desc><![CDATA[Gate 21]]></desc>
 <sym>Truck Stop</sym>
 <type><![CDATA[Road]]></type>
</wpt>
<wpt lat="42.456718" lon="-71.102973">
 <ele>81.076800</ele>
 <time>2001-06-02T03:27:03Z</time>
 <name>GATE24</name>
 <desc><![CDATA[Gate 24]]></desc>
 <sym>Truck Stop</sym>
 <type><![CDATA[Road]]></type>
</wpt>
<wpt lat="42.430847" lon="-71.107690">
 <ele>21.515015</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>GATE5</name>
 <desc><![CDATA[Gate 5]]></desc>
 <sym>Truck Stop</sym>
 <type><![CDATA[Truck Stop]]></type>
</wpt>
<wpt lat="42.431240" lon="-71.109236">
 <ele>26.561890</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>GATE6</name>
 <desc><![CDATA[Gate 6]]></desc>
 <sym>Trailhead</sym>
 <type><![CDATA[Trail Head]]></type>
</wpt>
<wpt lat="42.439502" lon="-71.106556">
 <ele>32.004000</ele>
 <time>2001-06-02T00:18:16Z</time>
 <name>6077LOGS</name>
 <desc><![CDATA[Log Crossing]]></desc>
 <sym>Amusement Park</sym>
 <type><![CDATA[Obstacle]]></type>
</wpt>
<wpt lat="42.449765" lon="-71.122320">
 <ele>119.809082</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>5148NANEPA</name>
 <desc><![CDATA[Nanepashemet Road Crossing]]></desc>
 <sym>Trailhead</sym>
 <type><![CDATA[Trail Head]]></type>
</wpt>
<wpt lat="42.457388" lon="-71.119845">
 <ele>73.761600</ele>
 <time>2001-06-02T03:27:00Z</time>
 <name>5267OBSTAC</name>
 <desc><![CDATA[Obstacle]]></desc>
 <sym>Amusement Park</sym>
 <type><![CDATA[Obstacle]]></type>
</wpt>
<wpt lat="42.434980" lon="-71.109942">
 <ele>45.307495</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>PANTHRCAVE</name>
 <desc><![CDATA[Panther Cave]]></desc>
 <sym>Tunnel</sym>
 <type><![CDATA[Tunnel]]></type>
</wpt>
<wpt lat="42.453256" lon="-71.121211">
 <ele>77.992066</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>5252PURPLE</name>
 <desc><![CDATA[Purple Rock Hill]]></desc>
 <sym>Summit</sym>
 <type><![CDATA[Summit]]></type>
</wpt>
<wpt lat="42.457734" lon="-71.117481">
 <ele>67.970400</ele>
 <time>2001-06-02T03:27:01Z</time>
 <name>5287WATER</name>
 <desc><![CDATA[Reservoir]]></desc>
 <sym>Swimming Area</sym>
 <type><![CDATA[Reservoir]]></type>
</wpt>
<wpt lat="42.459278" lon="-71.124574">
 <ele>81.076800</ele>
 <time>2001-06-02T03:27:00Z</time>
 <name>5239ROAD</name>
 <desc><![CDATA[Road]]></desc>
 <sym>Truck Stop</sym>
 <type><![CDATA[Road]]></type>
</wpt>
<wpt lat="42.458782" lon="-71.118991">
 <ele>67.360800</ele>
 <time>2001-06-02T03:27:01Z</time>
 <name>5278ROAD</name>
 <desc><![CDATA[Road]]></desc>
 <sym>Truck Stop</sym>
 <type><![CDATA[Road]]></type>
</wpt>
<wpt lat="42.439993" lon="-71.120925">
 <ele>53.949600</ele>
 <time>2001-06-02T00:18:14Z</time>
 <name>5058ROAD</name>
 <cmt>ROAD CROSSING</cmt>
 <desc><![CDATA[Road Crossing]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Road Crossing]]></type>
</wpt>
<wpt lat="42.453415" lon="-71.106782">
 <ele>69.799200</ele>
 <time>2001-06-02T00:18:13Z</time>
 <name>SHEEPFOLD</name>
 <desc><![CDATA[Sheepfold Parking Lot]]></desc>
 <sym>Parking Area</sym>
 <type><![CDATA[Parking]]></type>
</wpt>
<wpt lat="42.455956" lon="-71.107483">
 <ele>64.008000</ele>
 <time>2001-06-02T03:27:04Z</time>
 <name>SOAPBOX</name>
 <desc><![CDATA[Soap Box Derby Track]]></desc>
 <sym>Cemetery</sym>
 <type><![CDATA[Intersection]]></type>
</wpt>
<wpt lat="42.465913" lon="-71.119328">
 <ele>64.533692</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>5376STREAM</name>
 <desc><![CDATA[Stream Crossing]]></desc>
 <sym>Bridge</sym>
 <type><![CDATA[Bridge]]></type>
</wpt>
<wpt lat="42.445359" lon="-71.122845">
 <ele>61.649902</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>5144SUMMIT</name>
 <desc><![CDATA[Summit]]></desc>
 <sym>Summit</sym>
 <type><![CDATA[Summit]]></type>
</wpt>
<wpt lat="42.441727" lon="-71.121676">
 <ele>67.360800</ele>
 <time>2001-06-02T00:18:16Z</time>
 <name>5150TANK</name>
 <cmt>WATER TANK</cmt>
 <desc><![CDATA[Water Tank]]></desc>
 <sym>Museum</sym>
 <type><![CDATA[Water Tank]]></type>
</wpt>
<rte>
 <name>BELLEVUE</name>
 <desc><![CDATA[Bike Loop Bellevue]]></desc>
 <number>1</number>
<rtept lat="42.430950" lon="-71.107628">
 <ele>23.469600</ele>
 <time>2001-06-02T00:18:15Z</time>
 <name>BELLEVUE</name>
 <cmt>BELLEVUE</cmt>
 <desc><![CDATA[Bellevue Parking Lot]]></desc>
 <sym>Parking Area</sym>
 <type><![CDATA[Parking]]></type>
</rtept>
<rtept lat="42.431240" lon="-71.109236">
 <ele>26.561890</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>GATE6</name>
 <desc><![CDATA[Gate 6]]></desc>
 <sym>Trailhead</sym>
 <type><![CDATA[Trail Head]]></type>
</rtept>
<rtept lat="42.434980" lon="-71.109942">
 <ele>45.307495</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>PANTHRCAVE</name>
 <desc><![CDATA[Panther Cave]]></desc>
 <sym>Tunnel</sym>
 <type><![CDATA[Tunnel]]></type>
</rtept>
<rtept lat="42.436757" lon="-71.113223">
 <ele>37.616943</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>6014MEADOW</name>
 <desc><![CDATA[6014MEADOW]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.439018" lon="-71.114456">
 <ele>56.388000</ele>
 <time>2001-06-02T03:26:55Z</time>
 <name>6006</name>
 <desc><![CDATA[600698]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.438594" lon="-71.114803">
 <ele>46.028564</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>6006BLUE</name>
 <desc><![CDATA[6006BLUE]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.438917" lon="-71.116146">
 <ele>44.826904</ele>
 <time>2001-11-16T23:03:38Z</time>
 <name>5096</name>
 <desc><![CDATA[5096]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.438878" lon="-71.119277">
 <ele>44.586548</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>5066</name>
 <desc><![CDATA[5066]]></desc>
 <sym>Crossing</sym>
 <type><![CDATA[Crossing]]></type>
</rtept>
<rtept lat="42.439227" lon="-71.119689">
 <ele>57.607200</ele>
 <time>2001-06-02T03:26:55Z</time>
 <name>5067</name>
 <desc><![CDATA[5067]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.439993" lon="-71.120925">
 <ele>53.949600</ele>
 <time>2001-06-02T00:18:14Z</time>
 <name>5058ROAD</name>
 <cmt>ROAD CROSSING</cmt>
 <desc><![CDATA[Road Crossing]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Road Crossing]]></type>
</rtept>
<rtept lat="42.441727" lon="-71.121676">
 <ele>67.360800</ele>
 <time>2001-06-02T00:18:16Z</time>
 <name>5150TANK</name>
 <cmt>WATER TANK</cmt>
 <desc><![CDATA[Water Tank]]></desc>
 <sym>Museum</sym>
 <type><![CDATA[Water Tank]]></type>
</rtept>
<rtept lat="42.443904" lon="-71.122044">
 <ele>50.594727</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>5142</name>
 <desc><![CDATA[5142]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.445359" lon="-71.122845">
 <ele>61.649902</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>5144SUMMIT</name>
 <desc><![CDATA[Summit]]></desc>
 <sym>Summit</sym>
 <type><![CDATA[Summit]]></type>
</rtept>
<rtept lat="42.447298" lon="-71.121447">
 <ele>127.711200</ele>
 <time>2001-06-02T03:26:58Z</time>
 <name>5156</name>
 <desc><![CDATA[5156]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.449765" lon="-71.122320">
 <ele>119.809082</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>5148NANEPA</name>
 <desc><![CDATA[Nanepashemet Road Crossing]]></desc>
 <sym>Trailhead</sym>
 <type><![CDATA[Trail Head]]></type>
</rtept>
<rtept lat="42.451442" lon="-71.121746">
 <ele>74.627442</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>5258</name>
 <desc><![CDATA[5258]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.453256" lon="-71.121211">
 <ele>77.992066</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>5252PURPLE</name>
 <desc><![CDATA[Purple Rock Hill]]></desc>
 <sym>Summit</sym>
 <type><![CDATA[Summit]]></type>
</rtept>
<rtept lat="42.456252" lon="-71.119356">
 <ele>78.713135</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>527631</name>
 <desc><![CDATA[527631]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.456592" lon="-71.119676">
 <ele>78.713135</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>527614</name>
 <desc><![CDATA[527614]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.457388" lon="-71.119845">
 <ele>73.761600</ele>
 <time>2001-06-02T03:27:00Z</time>
 <name>5267OBSTAC</name>
 <desc><![CDATA[Obstacle]]></desc>
 <sym>Amusement Park</sym>
 <type><![CDATA[Obstacle]]></type>
</rtept>
<rtept lat="42.458148" lon="-71.119135">
 <ele>68.275200</ele>
 <time>2001-06-02T03:27:00Z</time>
 <name>5278</name>
 <desc><![CDATA[5278]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.459377" lon="-71.117693">
 <ele>64.008000</ele>
 <time>2001-06-02T03:27:01Z</time>
 <name>5289</name>
 <desc><![CDATA[5289]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.464183" lon="-71.119828">
 <ele>52.997925</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>5374FIRE</name>
 <desc><![CDATA[5374FIRE]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.465650" lon="-71.119399">
 <ele>56.388000</ele>
 <time>2001-06-02T03:27:02Z</time>
 <name>5376</name>
 <desc><![CDATA[5376]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.465913" lon="-71.119328">
 <ele>64.533692</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>5376STREAM</name>
 <desc><![CDATA[Stream Crossing]]></desc>
 <sym>Bridge</sym>
 <type><![CDATA[Bridge]]></type>
</rtept>
<rtept lat="42.467110" lon="-71.113574">
 <ele>53.644800</ele>
 <time>2001-06-02T03:27:02Z</time>
 <name>6328</name>
 <desc><![CDATA[6328]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.466459" lon="-71.110067">
 <ele>48.768000</ele>
 <time>2001-06-02T03:27:02Z</time>
 <name>635722</name>
 <desc><![CDATA[635722]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.466557" lon="-71.109410">
 <ele>49.072800</ele>
 <time>2001-06-02T03:27:02Z</time>
 <name>635783</name>
 <desc><![CDATA[635783]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.463495" lon="-71.107117">
 <ele>62.484000</ele>
 <time>2001-06-02T03:27:03Z</time>
 <name>6373</name>
 <desc><![CDATA[6373]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.465687" lon="-71.107360">
 <ele>87.782400</ele>
 <time>2001-06-02T03:27:03Z</time>
 <name>BEAR HILL</name>
 <cmt>BEAR HILL TOWER</cmt>
 <desc><![CDATA[Bear Hill Tower]]></desc>
 <sym>Tall Tower</sym>
 <type><![CDATA[Tower]]></type>
</rtept>
<rtept lat="42.459986" lon="-71.106170">
 <ele>72.945191</ele>
 <time>2001-11-16T23:03:38Z</time>
 <name>6289</name>
 <desc><![CDATA[6289]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.457616" lon="-71.105116">
 <ele>72.847200</ele>
 <time>2001-06-02T03:27:04Z</time>
 <name>6297</name>
 <desc><![CDATA[6297]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.453845" lon="-71.105206">
 <ele>66.696655</ele>
 <time>2001-11-16T23:03:38Z</time>
 <name>6283</name>
 <desc><![CDATA[6283]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.451430" lon="-71.105413">
 <ele>57.564209</ele>
 <time>2001-11-16T23:03:38Z</time>
 <name>6280</name>
 <desc><![CDATA[6280]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.448448" lon="-71.106158">
 <ele>62.179200</ele>
 <time>2001-06-02T03:27:04Z</time>
 <name>6177</name>
 <desc><![CDATA[6177]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.447804" lon="-71.106624">
 <ele>62.484000</ele>
 <time>2001-06-02T03:27:04Z</time>
 <name>6176</name>
 <desc><![CDATA[6176]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.444773" lon="-71.108882">
 <ele>62.788800</ele>
 <time>2001-06-02T03:27:05Z</time>
 <name>6153</name>
 <desc><![CDATA[6153]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.443592" lon="-71.106301">
 <ele>55.473600</ele>
 <time>2001-06-02T03:27:05Z</time>
 <name>6171</name>
 <desc><![CDATA[6171]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.442981" lon="-71.111441">
 <ele>64.008000</ele>
 <time>2001-06-02T03:26:58Z</time>
 <name>6131</name>
 <desc><![CDATA[6131]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.442196" lon="-71.110975">
 <ele>64.008000</ele>
 <time>2001-06-02T03:26:55Z</time>
 <name>6130</name>
 <desc><![CDATA[6130]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.441754" lon="-71.113220">
 <ele>56.388000</ele>
 <time>2001-06-02T03:26:55Z</time>
 <name>6029</name>
 <desc><![CDATA[6029]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.439018" lon="-71.114456">
 <ele>56.388000</ele>
 <time>2001-06-02T03:26:55Z</time>
 <name>6006</name>
 <desc><![CDATA[600698]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Intersection]]></type>
</rtept>
<rtept lat="42.436757" lon="-71.113223">
 <ele>37.616943</ele>
 <time>2001-11-28T21:05:28Z</time>
 <name>6014MEADOW</name>
 <desc><![CDATA[6014MEADOW]]></desc>
 <sym>Dot</sym>
 <type><![CDATA[Dot]]></type>
</rtept>
<rtept lat="42.434980" lon="-71.109942">
 <ele>45.307495</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>PANTHRCAVE</name>
 <desc><![CDATA[Panther Cave]]></desc>
 <sym>Tunnel</sym>
 <type><![CDATA[Tunnel]]></type>
</rtept>
<rtept lat="42.431240" lon="-71.109236">
 <ele>26.561890</ele>
 <time>2001-11-07T23:53:41Z</time>
 <name>GATE6</name>
 <desc><![CDATA[Gate 6]]></desc>
 <sym>Trailhead</sym>
 <type><![CDATA[Trail Head]]></type>
</rtept>
<rtept lat="42.430950" lon="-71.107628">
 <ele>23.469600</ele>
 <time>2001-06-02T00:18:15Z</time>
 <name>BELLEVUE</name>
 <cmt>BELLEVUE</cmt>
 <desc><![CDATA[Bellevue Parking Lot]]></desc>
 <sym>Parking Area</sym>
 <type><![CDATA[Parking]]></type>
</rtept>
</rte>
</gpx>
`)
