package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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
	Alerts      int64   `json:"alerts" bson:"alerts" valid:"-"`           //Device : alerts
}

/*func decodeSensitFrame(rawData SigfoxMessage) (message SigfoxMessage) {
	return decodedMessage
}*/

func (mes *SigfoxMessage) BeforeCreate() {
	//*l = decodeSensitFrame(*l)
	mes.Id = bson.NewObjectId().Hex()

	// TODO : Fix shift when battery MSB=0
	// TODO : Handle modes 2, 3, 4 & 5

	data := ""
	if mes.MesType == 1 {
		fmt.Println("Sensit Message")
		parsed, err := strconv.ParseUint(mes.Data, 16, 32)
		if err != nil {
			log.Fatal(err)
		}
		data = fmt.Sprintf("%08b", parsed)
		/*byte1 := data[0:8]
		byte2 := data[8:16]
		byte3 := data[16:24]
		byte4 := data[24:32]*/

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
		batVal := float32(battery) * 0.05 * 2.7

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

	} else if mes.MesType == 2 {
		fmt.Println("Arduino Message")
		return

	} else {
		return
	}
}

const SigfoxMessagesCollection = "sigfox_messages"
