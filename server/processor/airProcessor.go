package processor

import (
	"encoding/json"
	"fmt"
	"github.com/wxmsummer/AirConditioner/common/message"
	"github.com/wxmsummer/AirConditioner/common/utils"
	"github.com/wxmsummer/AirConditioner/server/model"
	"github.com/wxmsummer/AirConditioner/server/repository"
	"net"
)

type AirProcessor struct {
	Conn net.Conn
	Orm  *repository.AirConditionerOrm
}

// 查询空调状态，按房间号查询
func (ap *AirProcessor) FindByRoom(msg *message.Message) (err error) {
	var findByRoom message.AirConditionerFindByRoom
	err = json.Unmarshal([]byte(msg.Data), &findByRoom)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var findByRoomRes message.AirConditionerFindByRoomRes

	roomNum := findByRoom.RoomNum
	airConditioner, err := ap.Orm.FindByRoom(roomNum)
	if err != nil {
		fmt.Println("ap.Orm.FindByRoom err=", err)
		return err
	}

	findByRoomRes.Code = 200
	findByRoomRes.Msg = "根据房间号查询空调成功！"
	findByRoomRes.AirConditioner = airConditioner

	data, err := json.Marshal(findByRoomRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeAirConditionerFindByRoomRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// FindAll 查询所有空调状态
func (ap *AirProcessor) FindAll(msg *message.Message) (err error) {
	var findAll message.AirConditionerFindAll
	err = json.Unmarshal([]byte(msg.Data), &findAll)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var findAllRes message.AirConditionerFindAllRes

	airConditioners, err := ap.Orm.FindAll()
	if err != nil {
		fmt.Println("ap.Orm.FindAll() err=", err)
		return err
	}

	findAllRes.Code = 200
	findAllRes.Msg = "查询所有空调状态成功！"
	findAllRes.AirConditioners = airConditioners

	data, err := json.Marshal(findAllRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeAirConditionerFindAllRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// Create 新增一条空调状态记录
func (ap *AirProcessor) Create(msg *message.Message) (err error) {
	var createMsg message.AirConditionerCreate
	err = json.Unmarshal([]byte(msg.Data), &createMsg)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var CreateRes message.NormalRes

	airConditioner := createMsg.AirConditioner

	existAir, _ := ap.Orm.FindByRoom(airConditioner.RoomNum)
	// 如果已存在，则返回重新创建提示
	if existAir.RoomNum != 0 {
		CreateRes.Code = 501
		CreateRes.Msg = "空调Room已存在，请重新创建！"
	} else {
		err = ap.Orm.Create(&airConditioner)
		if err != nil {
			fmt.Println("ap.Orm.Create(&airConditioner) err=", err)
			return
		}
		CreateRes.Code = 200
		CreateRes.Msg = "创建成功！"
	}

	data, err := json.Marshal(CreateRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeAirConditionerCreateRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

func (ap *AirProcessor) Update(msg *message.Message) (err error) {
	var updateMsg message.AirConditionerUpdate
	err = json.Unmarshal([]byte(msg.Data), &updateMsg)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var updateRes message.NormalRes

	airConditioner := updateMsg.AirConditioner
	err = ap.Orm.Update(airConditioner)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	updateRes.Code = 200
	updateRes.Msg = "更新成功！"

	data, err := json.Marshal(updateRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// AirConditionerOn
func (ap *AirProcessor) PowerOn(msg *message.Message) (err error) {
	var powerOn message.AirConditionerOn
	err = json.Unmarshal([]byte(msg.Data), &powerOn)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var normalRes message.NormalRes

	// 先取出该空调状态数据
	air, err := ap.Orm.FindByRoom(powerOn.RoomNum)
	if err != nil {
		return err
	}
	// 这里需要对 OpenTime 进行处理
	if air.OpenTime == nil { // 如果是第一次开机，就初始化 OpenTime
		air.OpenTime = []int64{powerOn.OpenTime}
	} else { // 否则，就append OpenTime 列表
		air.OpenTime = append(air.OpenTime, powerOn.OpenTime)
	}
	air.Power = model.PowerOn
	air.Mode = powerOn.Mode
	air.WindLevel = powerOn.WindLevel
	air.Temperature = powerOn.Temperature

	err = ap.Orm.Update(air)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	normalRes.Code = 200
	normalRes.Msg = "开机成功！"

	data, err := json.Marshal(normalRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// AirConditionerOff
func (ap *AirProcessor) PowerOff(msg *message.Message) (err error) {
	var powerOff message.AirConditionerOff
	err = json.Unmarshal([]byte(msg.Data), &powerOff)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var normalRes message.NormalRes

	// 先取出该空调状态数据
	air, err := ap.Orm.FindByRoom(powerOff.RoomNum)
	if err != nil {
		return err
	}
	// 这里需要对 CloseTime 进行处理
	if air.CloseTime == nil { // 如果是第一次关机，就初始化 CloseTime
		air.CloseTime = []int64{powerOff.CloseTime}
	} else { // 否则，就append CloseTime 列表
		air.CloseTime = append(air.CloseTime, powerOff.CloseTime)
	}
	air.Power = model.PowerOff
	err = ap.Orm.Update(air)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	normalRes.Code = 200
	normalRes.Msg = "关机成功！"

	data, err := json.Marshal(normalRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// SetParam
func (ap *AirProcessor) SetParam(msg *message.Message) (err error) {
	var setParam message.AirConditionerSetParam
	err = json.Unmarshal([]byte(msg.Data), &setParam)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var normalRes message.NormalRes

	// 先取出该空调状态数据
	air, err := ap.Orm.FindByRoom(setParam.RoomNum)
	if err != nil {
		return err
	}

	// 设置相应参数
	air.Mode = setParam.Mode
	air.WindLevel = setParam.WindLevel
	air.Temperature = setParam.Temperature
	// 调整次数加一
	air.SetParamNum += 1

	err = ap.Orm.Update(air)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	normalRes.Code = 200
	normalRes.Msg = "调整空调参数成功！"

	data, err := json.Marshal(normalRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// StopWind
func (ap *AirProcessor) StopWind(msg *message.Message) (err error) {
	var stopWind message.AirConditionerStopWind
	err = json.Unmarshal([]byte(msg.Data), &stopWind)
	if err != nil {
		fmt.Println("json.Unmarshal fail, err =", err)
		return
	}

	var resMsg message.Message
	var normalRes message.NormalRes

	// 先取出该空调状态数据
	air, err := ap.Orm.FindByRoom(stopWind.RoomNum)
	if err != nil {
		return err
	}

	// 这里需要对 stopWind 进行处理
	if air.StopWind == nil { // 如果是第一次停止送风，就初始化 stopWind
		air.StopWind = []int64{stopWind.StopWindTime}
	} else { // 否则，就append StopWind 列表
		air.StopWind = append(air.StopWind, stopWind.StopWindTime)
	}

	err = ap.Orm.Update(air)
	if err != nil {
		fmt.Println("ap.Orm.Update(airConditioner) err=", err)
		return
	}

	normalRes.Code = 200
	normalRes.Msg = "空调停止送风成功！"

	data, err := json.Marshal(normalRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeNormalRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}

// GetReport
func (ap *AirProcessor) GetReport(msg *message.Message) (err error) {

	var resMsg message.Message
	var getReportRes message.GetReportRes

	airs, err := ap.Orm.FindAll()
	if err != nil {
		fmt.Println("ap.Orm.FindAll() err=", err)
		return err
	}

	for _, air := range airs {
		var report model.Report
		report.RoomNum = air.RoomNum
		report.TotalPower = air.TotalPower
		report.TotalFee = air.TotalPower * 1
		report.CloseNum = len(air.CloseTime)
		report.SetParamNum = air.SetParamNum

		for i := 0; i < len(air.CloseTime); i++ {
			// 计算空调的总开机时长：用关机数组的值逐个减去开机数组的值
			report.UsedTime += int(air.CloseTime[i] - air.OpenTime[i])
		}

		// 将空调报表逐个添加到 getReportRes.Reports 中
		getReportRes.Reports = append(getReportRes.Reports, report)
	}

	getReportRes.Code = 200
	getReportRes.Msg = "获取所有空调报表成功！"

	data, err := json.Marshal(getReportRes)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	resMsg.Type = message.TypeGetReportRes
	resMsg.Data = string(data)
	data, err = json.Marshal(resMsg)
	if err != nil {
		fmt.Println("json.Marshal fail, err=", err)
		return
	}

	tf := &utils.Transfer{Conn: ap.Conn}
	err = tf.WritePkg(data)
	return
}