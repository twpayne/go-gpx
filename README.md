# go-gpx

[![Build Status](https://travis-ci.org/twpayne/go-gpx.svg?branch=master)](https://travis-ci.org/twpayne/go-gpx)
[![GoDoc](https://godoc.org/github.com/twpayne/go-gpx?status.svg)](https://godoc.org/github.com/twpayne/go-gpx)
[![Report Card](https://goreportcard.com/badge/github.com/twpayne/go-gpx)](https://goreportcard.com/report/github.com/twpayne/go-gpx)

Package go-gpx provides convenince methods for reading and writing GPX documents.

Read example:

```go
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
```

Write example:

```go
	g := &GPX{
		Version: "1.0",
		Creator: "ExpertGPS 1.1 - http://www.topografix.com",
		Wpt: []*WptType{
			{
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
	if err := g.WriteIndent(os.Stdout, "", "  "); err != nil {
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
```

[License](LICENSE)
