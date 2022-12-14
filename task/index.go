package task

import (
	"QluTakeLesson/app"
	"QluTakeLesson/utils/RuiJIeNet"
	"QluTakeLesson/utils/config"
	"QluTakeLesson/utils/log"
	"time"
)

// 定时任务

// DailyTask 每日5.45时任务
func DailyTask() {
	// 加载用户配置
	app.ReloadUserConfig()
	// 初始化区域信息
	app.InitArea()
	// 预约时间信息更新
	app.UpdateSegmentList()
	// 预约时间更新
	app.ReserveTime = time.Now().AddDate(0, 0, 1)
	// 登录
	go app.Login()

}

// NetworkCheck 网络检测
func NetworkCheck() {
	if config.Config.RuiJie.Enable {
		// 检测网络
		isNetConn := app.CheckNetwork()
		if !isNetConn {
			log.Warning("网络连接失败，尝试认证校园网")
			// 重新登录
			RuiJIeNet.ExecuteLogin()
			isSuccess := RuiJIeNet.QueryLoginResult()
			if isSuccess {
				log.Info("认证成功")
			} else {
				log.Warning("认证失败")
			}
		}
	}
}

func LoginCheck() {
	isExpire := app.CheckLoginExpire()
	if isExpire {
		log.Warning("登录已过期 尝试重新登陆...")
		go app.Login()
	}
}

func CheckReserveTime() {
	if app.ReserveTime.Day() != time.Now().Add(24*time.Hour).Day() {
		// 预约时间更新
		app.ReserveTime = time.Now().AddDate(0, 0, 1)
	}
	segmentList := app.GetSegmentList()
	// 可预约时间段
	if len(segmentList) == 0 {
		app.UpdateSegmentList()
	} else {
		// 检测预约时间是否过期
		segment := segmentList[0]
		if segment.Day != time.Now().Add(24*time.Hour).Format("2006-01-02") {
			app.UpdateSegmentList()
		}
	}

}

func BootStrap() {
	// 每日5.45时任务
	go func() {
		for {
			now := time.Now()
			// 如果在5.45之前
			if now.Hour() < 5 || (now.Hour() == 5 && now.Minute() < 45) {
				// 等待到5.45
				time.Sleep(time.Duration(5-now.Hour())*time.Hour + time.Duration(45-now.Minute())*time.Minute)
			} else {
				// 等待到第二天5.45
				time.Sleep(time.Duration(24-now.Hour()+5)*time.Hour + time.Duration(45-now.Minute())*time.Minute)
			}
			// 执行任务
			DailyTask()
		}
	}()

	// 每隔5分钟检测一次网络，以及用户登录状态
	go func() {
		for {
			log.Info("执行定时任务")
			// 网络检测
			NetworkCheck()
			time.Sleep(5 * time.Second)
			// 登录检测
			LoginCheck()
			// 预约时间检测
			CheckReserveTime()
			time.Sleep(5 * time.Minute)
		}
	}()

}
