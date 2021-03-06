package models

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strconv"
	"strings"
	//"github.com/IoThingsDev/api/store"
	"gopkg.in/gin-gonic/gin.v1"
	"time"
)

/*func NewStoreController() UserController {
	return UserController{}
}
*/
type SigfoxMessage struct {
	Id          string  `json:"id" bson:"_id,omitempty" valid:"-"`
	SigfoxId    string  `json:"sigfoxId" bson:"sigfoxId" valid:"-"`
	FrameNumber uint    `json:"frameNumber" bson:"frameNumber" valid:"-"` //Device : (daily frames under 140)
	Timestamp   int64   `json:"timestamp" bson:"timestamp" valid:"-"`     //Sigfox : time
	Station     string  `json:"station" bson:"station" valid:"-"`         //Sigfox : station
	Snr         float64 `json:"snr" bson:"snr" valid:"-"`                 //Sigfox : snr
	AvgSnr      float64 `json:"avgSnr" bson:"avgSnr" valid:"-"`           //Sigfox : avgSnr
	Rssi        float64 `json:"rssi" bson:"rssi" valid:"-"`               //Sigfox : rssi
	MesType     uint8   `json:"mesType" bson:"mesType" valid:"-"`         //Sigfox : mesType
	Data        string  `json:"data" bson:"data" valid:"-"`               //Sigfox : data
	EventType   string  `json:"eventType" bson:"eventType" valid:"-"`     //Device : eventType
	SwRev       string  `json:"swRev" bson:"swRev" valid:"-"`             //Device : swRev
	Mode        string  `json:"mode" bson:"mode" valid:"-"`               //Device : mode
	Timeframe   string  `json:"timeframe" bson:"timeframe" valid:"-"`     //Device : timeframe
	Data1       float64 `json:"data1" bson:"data1" valid:"-"`             //Device : battery
	Data2       float64 `json:"data2" bson:"data2" valid:"-"`             //Device : temperature
	Data3       float64 `json:"data3" bson:"data3" valid:"-"`             //Device : humidity
	Data4       float64 `json:"data4" bson:"data4" valid:"-"`             //Device : light
	Data5       float64 `json:"data5" bson:"data5" valid:"-"`             //Device : custom
	Data6       float64 `json:"data6" bson:"data6" valid:"-"`             //Device : custom
	Alerts      int64   `json:"alerts" bson:"alerts" valid:"-"`           //Device : alerts
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
			batVal := (float64(battery) * 0.05) + 2.7

			mes.Data1 = batVal

			//Byte 3
			temperature := int64(0)
			tempVal := float64(0)

			reedSwitch := false
			if mode == 0 || mode == 1 {
				temperatureLsb := data[18:24]
				tempData := []string{temperatureMsb, temperatureLsb}
				temperature, _ := strconv.ParseInt(strings.Join(tempData, ""), 2, 16)
				tempVal = (float64(temperature) - 200) / 8
				if data[17] == 1 {
					reedSwitch = true
				}
			} else {
				temperature, _ = strconv.ParseInt(temperatureMsb, 2, 16)
				tempVal = (float64(temperature) - 200) / 8
			}

			mes.Data2 = tempVal

			modeStr := ""
			swRev := ""
			humidity := float64(0.0)
			light := float64(0.0)

			switch mode {
			case 0:
				modeStr = "Button"
				majorSwRev, _ := strconv.ParseInt(data[24:28], 2, 8)
				minorSwRev, _ := strconv.ParseInt(data[28:32], 2, 8)
				swRev = fmt.Sprintf("%d.%d", majorSwRev, minorSwRev)
			case 1:
				modeStr = "Temperature + Humidity"
				humi, _ := strconv.ParseInt(data[24:32], 2, 16)
				humidity = float64(humi) * 0.5
				mes.Data3 = humidity
			case 2:
				modeStr = "Light"
				lightVal, _ := strconv.ParseInt(data[18:24], 2, 8)
				lightMulti, _ := strconv.ParseInt(data[17:18], 2, 8)
				light = float64(lightVal) * 0.01
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
	} else {
		return
	}
}

/*
 *	b2 1f		dc 9b		 8b 0c		 7e a4		 01 00		 00 00
 *	19.81, 60.88, 980.14, 64.26,
 * 	1FB2, 9BDC, 0C8B, A47E
 */

func convertInt16toFloat(value float64, min float64, max float64) float64 {
	return (value * (max - min)) / 32768
}

func convertUInt16toFloat(value float64, min float64, max float64) float64 {
	return (value * (max - min)) / 65536
}

func (l *SigfoxMessage) FormatData(valueType string) gin.H {
	receivedTime := time.Unix(l.Timestamp, 0)
	var value interface{}
	//TODO: Handle devices configurations
	switch valueType {
	case "humidity":
		value = l.Data3
	case "temperature":
		value = l.Data2
	}

	return gin.H{"value": value, "latitude": 50.434995, "longitude": 30.823634, "timestamp": receivedTime.String()}
}

const SigfoxMessagesCollection = "sigfox_messages"
