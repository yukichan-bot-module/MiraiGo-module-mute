package mute

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/config"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

var instance *mute
var logger = utils.GetModuleLogger("com.aimerneige.mute")
var defaultMuteTime = 2

type mute struct {
}

func init() {
	instance = &mute{}
	bot.RegisterModule(instance)
}

func (m *mute) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       "com.aimerneige.mute",
		Instance: instance,
	}
}

// Init 初始化过程
// 在此处可以进行 Module 的初始化配置
// 如配置读取
func (m *mute) Init() {
	defaultMuteTime = config.GlobalConfig.GetInt("aimerneige.mute.default")
	if defaultMuteTime == 0 {
		defaultMuteTime = 2
	}
}

// PostInit 第二次初始化
// 再次过程中可以进行跨 Module 的动作
// 如通用数据库等等
func (m *mute) PostInit() {
}

// Serve 注册服务函数部分
func (m *mute) Serve(b *bot.Bot) {
	b.GroupMessageEvent.Subscribe(func(c *client.QQClient, msg *message.GroupMessage) {
		groupCode := msg.GroupCode
		// 检查发送者管理员权限
		senderMemberInfo, err := c.GetMemberInfo(groupCode, msg.Sender.Uin)
		if err != nil {
			errMsg := fmt.Sprintf("在群「%d」获取成员「%d」的用户数据时发成错误，详情请查阅后台日志。", groupCode, msg.Sender.Uin)
			logger.WithError(err).Errorf(errMsg)
			c.SendGroupMessage(groupCode, simpleText(errMsg))
			return
		}
		// 发送者没有管理员权限，忽略消息
		if senderMemberInfo.Permission != client.Administrator && senderMemberInfo.Permission != client.Owner {
			return
		}
		// 检查 bot 管理员权限
		botPermission := true
		botMemberInfo, err := c.GetMemberInfo(groupCode, c.Uin)
		if err != nil {
			errMsg := fmt.Sprintf("在群「%d」获取成员「%d」的用户数据时发成错误，详情请查阅后台日志。", groupCode, msg.Sender.Uin)
			logger.WithError(err).Errorf(errMsg)
			c.SendGroupMessage(groupCode, simpleText(errMsg))
			return
		}
		if botMemberInfo.Permission != client.Administrator && botMemberInfo.Permission != client.Owner {
			botPermission = false
		}
		// 处理全体禁言
		if msg.ToString() == "开启全体禁言" || msg.ToString() == "开启全员禁言" {
			if botPermission == false {
				c.SendGroupMessage(groupCode, simpleText("请先授予机器人管理员权限。"))
				return
			}
			muteAll(c, msg, true)
			return
		}
		if msg.ToString() == "关闭全体禁言" || msg.ToString() == "关闭全员禁言" {
			if botPermission == false {
				c.SendGroupMessage(groupCode, simpleText("请先授予机器人管理员权限。"))
				return
			}
			muteAll(c, msg, false)
			return
		}
		// 处理禁言成员指令
		if len(msg.Elements) < 2 {
			return
		}
		// 解析指令
		isAt := false
		isMute := false
		muteCommand := ""
		var target int64
		for _, ele := range msg.Elements {
			switch e := ele.(type) {
			case *message.AtElement:
				isAt = true
				target = e.Target
			case *message.TextElement:
				if isMute == false {
					muteCommand = e.Content
					muteCommand = strings.TrimSpace(muteCommand)
					if strings.HasPrefix(muteCommand, "禁言") {
						isMute = true
					}
				}
			}
		}
		if isAt && isMute && target != 0 {
			if botPermission == false {
				c.SendGroupMessage(groupCode, simpleText("请先授予机器人管理员权限。"))
				return
			}
			if target == c.Uin {
				c.SendGroupMessage(groupCode, simpleText("不要禁言我啊~"))
				return
			}
			targetMemberInfo, err := c.GetMemberInfo(groupCode, target)
			if err != nil {
				errMsg := fmt.Sprintf("在群「%d」获取成员「%d」的用户数据时发成错误，详情请查阅后台日志。", groupCode, target)
				logger.WithError(err).Errorf(errMsg)
				c.SendGroupMessage(groupCode, simpleText(errMsg))
				return
			}
			if targetMemberInfo.Permission == client.Owner {
				c.SendGroupMessage(groupCode, simpleText("你居然想禁言群主？真是危险的想法呢~"))
				return
			}
			if targetMemberInfo.Permission == client.Administrator {
				c.SendGroupMessage(groupCode, simpleText("管理员是无法禁言的呢~"))
				return
			}
			muteTime := defaultMuteTime
			if muteCommand != "禁言" {
				timeCommand := muteCommand[6:]
				timeCommand = strings.TrimSpace(timeCommand)
				timeI64, err := strconv.ParseInt(timeCommand, 10, 64)
				if err != nil {
					errMsg := fmt.Sprintf("指令解析失败，「%s」不是正确的数字。", timeCommand)
					c.SendGroupMessage(groupCode, simpleText(errMsg))
				}
				if timeI64 > int64(4320) {
					timeI64 = 4320
				}
				if timeI64 < int64(0) {
					timeI64 = 0
				}
				muteTime = int(timeI64)
			}
			if err := targetMemberInfo.Mute(uint32(muteTime * 60)); err != nil {
				errMsg := fmt.Sprintf("在群「%d」尝试禁言成员「%d」的过程中发生错误，详情请查阅后台日志。", groupCode, target)
				logger.WithError(err).Errorf(errMsg)
				c.SendGroupMessage(groupCode, simpleText(errMsg))
				return
			}
			if muteTime == 0 {
				c.SendGroupMessage(groupCode, simpleText(fmt.Sprintf("群成员「%d」已被管理员「%d」解除禁言。", target, msg.Sender.Uin)))
				return
			}
			c.SendGroupMessage(groupCode, simpleText(fmt.Sprintf("群成员「%d」已被管理员「%d」禁言「%d」分钟。", target, msg.Sender.Uin, muteTime)))
		}
	})
}

// Start 此函数会新开携程进行调用
// ```go
//
//	go exampleModule.Start()
//
// ```
// 可以利用此部分进行后台操作
// 如 http 服务器等等
func (m *mute) Start(b *bot.Bot) {
}

// Stop 结束部分
// 一般调用此函数时，程序接收到 os.Interrupt 信号
// 即将退出
// 在此处应该释放相应的资源或者对状态进行保存
func (m *mute) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	// 别忘了解锁
	defer wg.Done()
}

func simpleText(s string) *message.SendingMessage {
	return message.NewSendingMessage().Append(message.NewText(s))
}

func muteAll(c *client.QQClient, msg *message.GroupMessage, mute bool) {
	if group := c.FindGroup(msg.GroupCode); group != nil {
		group.MuteAll(mute)
	}
}
