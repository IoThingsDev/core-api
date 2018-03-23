package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"context"
	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
	"gopkg.in/mgo.v2/bson"
)

type SigfoxMessage struct {
	Id          string  `json:"id" bson:"_id,omitempty" valid:"-"`
	SigfoxId    string  `json:"sigfoxId" bson:"sigfoxId" valid:"-"`
	FrameNumber uint    `json:"frameNumber" bson:"frameNumber" valid:"-"` //Device : (daily frames under 140)
	Timestamp   int64   `json:"timestamp" bson:"timestamp" valid:"-"`     //Sigfox : time
	Station     string  `json:"station" bson:"station" valid:"-"`         //Sigfox : station
	Snr         float32 `json:"snr" bson:"snr" valid:"-"`                 //Sigfox : snr
	AvgSnr      float32 `json:"avgSnr" bson:"avgSnr" valid:"-"`           //Sigfox : avgSnr
	Rssi        float32 `json:"rssi" bson:"rssi" valid:"-"`               //Sigfox : rssi
	MesType     uint8   `json:"mesType" bson:"mesType" valid:"-"`         //Sigfox : mesType
	Data        string  `json:"data" bson:"data" valid:"-"`               //Sigfox : data
	EventType   string  `json:"eventType" bson:"eventType" valid:"-"`     //Device : eventType
	SwRev       string  `json:"swRev" bson:"swRev" valid:"-"`             //Device : swRev
	Mode        string  `json:"mode" bson:"mode" valid:"-"`               //Device : mode
	Timeframe   string  `json:"timeframe" bson:"timeframe" valid:"-"`     //Device : timeframe
	Data1       float32 `json:"data1" bson:"data1" valid:"-"`             //Device : battery
	Data2       float32 `json:"data2" bson:"data2" valid:"-"`             //Device : temperature
	Data3       float32 `json:"data3" bson:"data3" valid:"-"`             //Device : humidity
	Data4       float32 `json:"data4" bson:"data4" valid:"-"`             //Device : light
	Data5       float32 `json:"data5" bson:"data5" valid:"-"`             //Device : custom
	Data6       float32 `json:"data6" bson:"data6" valid:"-"`             //Device : custom
	Alerts      int64   `json:"alerts" bson:"alerts" valid:"-"`           //Device : alerts
}

func getWifiPosition(ssids string) {
	fmt.Println("WiFi frame")
	ssid1 := string(ssids[0:12])
	ssid2 := string(ssids[12:24])
	fmt.Println("SSID1: ", ssid1, "\t SSID2:", ssid2)
	// TODO: Put Google API Key in config file, like: config.GetString(c, "google_api_key")
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyCN0Z78M1sIT6c2H8PL0KaaFmjkBUE4avQ"))
	if err != nil {
		log.Fatalf("API connection fatal error: %s", err)
	}
	r := &maps.GeolocationRequest{
		ConsiderIP: true,
		WiFiAccessPoints: []maps.WiFiAccessPoint{maps.WiFiAccessPoint{
			MACAddress:         ssid1,
		}, maps.WiFiAccessPoint{
			MACAddress:         ssid2,
		}},
	}
	resp, err := c.Geolocate(context.Background(), r)
	if err != nil {
		log.Fatalf("Fatal Geolocation Request error: %s", err)
	}

	pretty.Println(resp)
}

func decodeGPSFrame(frame string) {
	fmt.Println("GPS frame")
	var latitude, longitude float64
	var latDeg, latMin, latSec float64
	var lngDeg, lngMin, lngSec float64

	isNorth, isEast := false, false
	if string(frame[0:2]) == "4e" {
		isNorth = true
	}
	if string(frame[10:12]) == "45" {
		isEast = true
	}

	if isNorth {
		fmt.Print("N:")
	} else {
		fmt.Print("S:")
	}

	valLatDeg, _ := strconv.ParseInt(frame[2:4], 16, 8)
	latDeg = float64(valLatDeg)
	valLatMin, _ := strconv.ParseInt(frame[4:6], 16, 8)
	latMin = float64(valLatMin)
	valLatSec, _ := strconv.ParseInt(frame[6:8], 16, 8)
	latSec = float64(valLatSec)
	fmt.Println(latDeg, "°\t", latMin, "m\t", latSec, "s")

	latitude = float64(latDeg) + float64(latMin/60) + float64(latSec/3600)

	if isEast {
		fmt.Print("E:")
	} else {
		fmt.Print("W:")
	}
	valLngDeg, _ := strconv.ParseInt(frame[10:12], 16, 8)
	lngDeg = float64(valLngDeg)
	valLngMin, _ := strconv.ParseInt(frame[12:14], 16, 8)
	lngMin = float64(valLngMin)
	valLngSec, _ := strconv.ParseInt(frame[14:16], 16, 8)
	lngSec = float64(valLngSec)
	fmt.Println(lngDeg, "°\t", lngMin, "m\t", lngSec, "s")

	longitude = float64(lngDeg) + float64(lngMin/60) + float64(lngSec/3600)

	fmt.Println("Lat: ", latitude, "\t Lng:", longitude)
}

//MesType, 1=Sensit, 2=Arduino, 3= Wisol EVK
func (mes *SigfoxMessage) BeforeCreate() {
	//*l = decodeSensitFrame(*l)
	mes.Id = bson.NewObjectId().Hex()

	// TODO: Handle modes 2, 3, 4 & 5
	// TODO: Nice to have : Round to 2 digits precision

	data := ""
	if mes.MesType == 1 {
		if len(mes.Data) <= 12 { //8 exactly, 4 bytes
			fmt.Println("Sensit Uplink Message")

			parsed, err := strconv.ParseUint(mes.Data, 16, 32)
			if err != nil {
				log.Fatal(err)
			}
			data = fmt.Sprintf("%08b", parsed)
			/*byte1 := data[0:8]
			byte2 := data[8:16]
			byte3 := data[16:24]
			byte4 := data[24:32]*/

			if len(data) == 25 { //Low battery MSB
				fmt.Println("Sensit Low battery")
				return
			}

			//Byte 1
			mode, _ := strconv.ParseInt(data[5:8], 2, 8)
			timeframe, _ := strconv.ParseInt(data[3:5], 2, 8)
			eventType, _ := strconv.ParseInt(data[1:3], 2, 8)
			batteryMsb := data[0:1]

			//Byte 2
			temperatureMsb := data[8:12]
			batteryLsb := data[12:16]
			battData := []string{batteryMsb, batteryLsb}
			battery, _ := strconv.ParseInt(strings.Join(battData, ""), 2, 8)
			batVal := (float32(battery) * 0.05) + 2.7

			mes.Data1 = batVal

			//Byte 3
			temperature := int64(0)
			tempVal := float32(0)

			reedSwitch := false
			if mode == 0 || mode == 1 {
				temperatureLsb := data[18:24]
				tempData := []string{temperatureMsb, temperatureLsb}
				temperature, _ := strconv.ParseInt(strings.Join(tempData, ""), 2, 16)
				tempVal = (float32(temperature) - 200) / 8
				if data[17] == 1 {
					reedSwitch = true
				}
			} else {
				temperature, _ = strconv.ParseInt(temperatureMsb, 2, 16)
				tempVal = (float32(temperature) - 200) / 8
			}

			mes.Data2 = tempVal

			modeStr := ""
			swRev := ""
			humidity := float32(0.0)
			light := float32(0.0)

			switch mode {
			case 0:
				modeStr = "Button"
				majorSwRev, _ := strconv.ParseInt(data[24:28], 2, 8)
				minorSwRev, _ := strconv.ParseInt(data[28:32], 2, 8)
				swRev = fmt.Sprintf("%d.%d", majorSwRev, minorSwRev)
			case 1:
				modeStr = "Temperature + Humidity"
				humi, _ := strconv.ParseInt(data[24:32], 2, 16)
				humidity = float32(humi) * 0.5
				mes.Data3 = humidity
			case 2:
				modeStr = "Light"
				lightVal, _ := strconv.ParseInt(data[18:24], 2, 8)
				lightMulti, _ := strconv.ParseInt(data[17:18], 2, 8)
				light = float32(lightVal) * 0.01
				if lightMulti == 1 {
					light = light * 8
				}
				mes.Data4 = light
			case 3:
				modeStr = "Door"
			case 4:
				modeStr = "Move"
			case 5:
				modeStr = "Reed switch"
			default:
				modeStr = ""
			}

			timeStr := ""
			switch timeframe {
			case 0:
				timeStr = "10 mins"
			case 1:
				timeStr = "1 hour"
			case 2:
				timeStr = "6 hours"
			case 3:
				timeStr = "24 hours"
			default:
				timeStr = ""
			}

			typeStr := ""
			switch eventType {
			case 0:
				typeStr = "Regular, no alert"
			case 1:
				typeStr = "Button call"
			case 2:
				typeStr = "Alert"
			case 3:
				typeStr = "New mode"
			default:
				timeStr = ""
			}

			switch mode {
			case 0:
				//fmt.Println("v" + swRev)
			case 1:
				//fmt.Println(humidity, "% RH")
			case 2:
				//fmt.Println(light, "lux")
				alerts, _ := strconv.ParseInt(data[24:32], 2, 16)
				mes.Alerts = alerts
			case 3, 4, 5:
				alerts, _ := strconv.ParseInt(data[24:32], 2, 16)
				mes.Alerts = alerts
			}
			if reedSwitch {
				//fmt.Println("Reed switch on")
			}

			mes.SwRev = "v " + swRev
			mes.EventType = typeStr
			mes.Mode = modeStr
			mes.Timeframe = timeStr
		} else { //len: 24 exactly, 12 bytes
			fmt.Println("Sensit Daily Downlink Message")
		}
	} else if mes.MesType == 2 {
		/*
			msg.Seqnbr
			msg.time
		*/
		fmt.Println("Arduino Message")
		mes.Data1 = convertInt16toFloat(mes.Data1, -30, 50)          //Temp
		mes.Data2 = convertUInt16toFloat(mes.Data2, 0, 100)          //Humi
		mes.Data3 = convertUInt16toFloat(mes.Data3, 900, 1100) + 900 //Pres: 900 shift to avoid overflow for numbers above 200
		mes.Data4 = convertUInt16toFloat(mes.Data4, 0, 200)          //Gas
		return

	} else if mes.MesType == 3 {
		fmt.Println("Wisol EVK Message")
		if (string(mes.Data[0:2]) == "4e") || (string(mes.Data[0:2]) == "53") {
			decodeGPSFrame(mes.Data)
		} else {
			getWifiPosition(mes.Data)
		}
	} else {
		return
	}
}

/*
 *	b2 1f		dc 9b		 8b 0c		 7e a4		 01 00		 00 00
 *	19.81, 60.88, 980.14, 64.26,
 * 	1FB2, 9BDC, 0C8B, A47E
 */

func convertInt16toFloat(value float32, min float32, max float32) float32 {
	return (value * (max - min)) / 32768
}

func convertUInt16toFloat(value float32, min float32, max float32) float32 {
	return (value * (max - min)) / 65536
}

const SigfoxMessagesCollection = "sigfox_messages"
