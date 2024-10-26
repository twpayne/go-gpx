package gpx_test

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"time"

	gpx "github.com/twpayne/go-gpx"
)

func ExampleRead() {
	r := bytes.NewBufferString(`
		<?xml version="1.0" encoding="UTF-8"?>
		<gpx version="1.1" creator="ExpertGPS 1.1 - http://www.topografix.com" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://www.topografix.com/GPX/1/1" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 https://www.topografix.com/GPX/1/1/gpx.xsd">
		  <wpt lat="42.438878" lon="-71.119277">
			<ele>44.586548</ele>
			<speed>9.16</speed>
			<time>2001-11-28T21:05:28Z</time>
			<name>5066</name>
			<desc>5066</desc>
			<sym>Crossing</sym>
			<type>Crossing</type>
		  </wpt>
		</gpx>
	`)
	t, err := gpx.Read(r)
	if err != nil {
		fmt.Printf("err == %v", err)
		return
	}
	fmt.Printf("t.Wpt[0] == %+v", t.Wpt[0])
	// Output:
	// t.Wpt[0] == &{Lat:42.438878 Lon:-71.119277 Ele:44.586548 Speed:9.16 Course:0 Time:2001-11-28 21:05:28 +0000 UTC MagVar:0 GeoidHeight:0 Name:5066 Cmt: Desc:5066 Src: Link:[] Sym:Crossing Type:Crossing Fix: Sat:0 HDOP:0 VDOP:0 PDOP:0 AgeOfDGPSData:0 DGPSID:[] Extensions:<nil>}
}

func ExampleGPX_WriteIndent() {
	g := &gpx.GPX{
		Version: "1.1",
		Creator: "ExpertGPS 1.1 - http://www.topografix.com",
		Wpt: []*gpx.WptType{
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
	if _, err := os.Stdout.WriteString(xml.Header); err != nil {
		fmt.Printf("err == %v", err)
	}
	if err := g.WriteIndent(os.Stdout, "", "  "); err != nil {
		fmt.Printf("err == %v", err)
	}
	// Output:
	// <?xml version="1.0" encoding="UTF-8"?>
	// <gpx version="1.1" creator="ExpertGPS 1.1 - http://www.topografix.com" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://www.topografix.com/GPX/1/1" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 https://www.topografix.com/GPX/1/1/gpx.xsd">
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
