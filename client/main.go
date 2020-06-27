package main

import (
	"github.com/wxmsummer/AirConditioner/client/process"
	"github.com/wxmsummer/AirConditioner/common/message"
	"time"
)

func main() {

	// 测试AirProcessor.Create()
	up := &process.AirProcessor{}
	// // _ = up.Create(1002)

	// powerOn := message.AirConditionerOn{
	// 	RoomNum:     1001,
	// 	Mode:        "cold",
	// 	WindLevel:   "high",
	// 	Temperature: 26,
	// 	OpenTime:    1591873607,
	// }
	// _ = up.PowerOn(powerOn)

	// powerOff := message.AirConditionerOff{
	// 	RoomNum:     1005,
	// 	CloseTime:    1591873628,
	// }
	// _ = up.PowerOff(powerOff)

	// setParam := message.AirConditionerSetParam {
	// 	RoomNum: 1001,
	// 	Mode: "cold",
	// 	WindLevel: "high",
	// 	Temperature: 26,
	// }

	// _ = up.SetParam(setParam)

	// _ = up.WatchAir()

	// _ = up.GetReport()

	// getDetail := message.GetDetailList{
	// 	RoomNum: 1001,
	// }
	// _ = up.GetDetail(getDetail)



	stopWind := message.AirConditionerStopWind{
		RoomNum:1001,
	}
	_ = up.StopWind(stopWind)

	time.Sleep(time.Second*5)

	sp := &process.ScheduleProcessor{}
	_ = sp.GetServingQueue()
}
