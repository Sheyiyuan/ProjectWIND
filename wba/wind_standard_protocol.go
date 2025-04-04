package wba

// WindStandardProtocolAPI Wind标准协议API,提供了onebot11标准中的API接口。
type WindStandardProtocolAPI interface {
	// UnsafelySendMsg [不安全][需要master权限]按照QQ号或群号发送消息到指定的群组或用户。
	// 参数:
	// - msgType: 消息类型，如 "group" 或 "private"。
	// - groupId: 群组 ID，若为私聊则为 0。
	// - userId: 用户 ID，若为群聊则为 0。
	// - message: 要发送的消息内容。
	// - autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码）。
	UnsafelySendMsg(msgType string, groupId int64, userId int64, message string, autoEscape bool)

	// UnsafelySendPrivateMsg [不安全][需要master权限]按照QQ号发送私聊消息给指定用户。
	// 参数:
	// - userId: 用户 ID。
	// - message: 要发送的消息内容。
	// - autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码）。
	UnsafelySendPrivateMsg(userId int64, message string, autoEscape bool)

	// UnsafelySendGroupMsg [不安全][需要master权限]按照群号发送群聊消息到指定群组。
	// 参数:
	// - groupId: 群组 ID。
	// - message: 要发送的消息内容。
	// - autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码）。
	UnsafelySendGroupMsg(groupId int64, message string, autoEscape bool)

	// SendMsg 按照提供的消息事件信息回复指定消息。
	// 参数:
	// - msg: 消息事件信息。
	// - message: 要回复的消息内容。
	// - autoEscape: 是否自动转义消息内容。
	SendMsg(msg MessageEventInfo, message string, autoEscape bool)

	// SendPrivateMsg 按照提供的消息事件信息回复私聊消息。
	// 参数:
	// - msg: 消息事件信息。
	// - message: 要回复的消息内容。
	// - autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码）。
	SendPrivateMsg(msg MessageEventInfo, message string, autoEscape bool)

	// SendGroupMsg 按照提供的消息事件信息回复群聊消息。
	// 参数:
	// - msg: 消息事件信息。
	// - message: 要回复的消息内容。
	// - autoEscape: 消息内容是否作为纯文本发送（即不解析 CQ 码）。
	SendGroupMsg(msg MessageEventInfo, message string, autoEscape bool)

	// UnsafelyDeleteMsg [不安全][需要master权限]撤回指定消息。
	// 参数:
	// - msgId: 消息 ID。
	UnsafelyDeleteMsg(msgId int32)

	// DeleteMsg 按照提供的消息事件信息撤回指定的消息。
	// 参数:
	// - msg: 消息事件信息。
	DeleteMsg(msg MessageEventInfo)

	// SendLike 给指定用户发送点赞，支持指定点赞次数。
	// 参数:
	// - userId: 用户 ID。
	// - times: 赞的次数，每个好友每天最多 10 次。
	SendLike(userId int64, times int)

	// SetGroupKick 将指定用户踢出指定群组，支持拒绝该用户再次加入。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - rejectAddRequest: 是否拒绝该用户再次加入。
	SetGroupKick(groupId int64, userId int64, rejectAddRequest bool)

	// SetGroupBan 禁言指定用户在指定群组的发言，支持指定禁言时长。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - duration: 禁言时长，单位为秒，0 表示取消禁言。
	SetGroupBan(groupId int64, userId int64, duration int32)

	// SetGroupWholeBan 开启或关闭指定群组的全员禁言。
	// 参数:
	// - groupId: 群组 ID。
	// - enable: 是否开启全员禁言。
	SetGroupWholeBan(groupId int64, enable bool)

	// SetGroupAdmin 设置或取消指定用户在指定群组的管理员权限。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - enable: 是否设置为管理员。
	SetGroupAdmin(groupId int64, userId int64, enable bool)

	// SetGroupLeave 退出指定群组。
	// 参数:
	// - groupId: 群组 ID。
	// - isDismiss: 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散。
	SetGroupLeave(groupId int64, isDismiss bool)

	// SetGroupCard 设置指定用户在指定群组的群名片。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - card: 群名片内容。
	SetGroupCard(groupId int64, userId int64, card string)

	// SetGroupName 设置指定群组的名称。
	// 参数:
	// - groupId: 群组 ID。
	// - groupName: 群组名称。
	SetGroupName(groupId int64, groupName string)

	// SetGroupSpecialTitle 设置指定用户在指定群组的专属头衔，支持指定头衔有效期。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - specialTitle: 专属头衔内容。
	SetGroupSpecialTitle(groupId int64, userId int64, specialTitle string)

	// SetFriendAddRequest 处理好友添加请求，支持同意或拒绝并添加备注。
	// 参数:
	// - flag: 请求标识(需从上报的数据中获得)。
	// - approve: 是否同意请求。
	// - remark: 备注信息（仅在同意时有效）。
	SetFriendAddRequest(flag string, approve bool, remark string)

	// SetGroupAddRequest 处理群组添加请求，支持同意或拒绝并添加理由。
	// 参数:
	// - flag: 请求标识。
	// - subType: 请求子类型。
	// - approve: 是否同意请求。
	// - reason: 拒绝理由。
	SetGroupAddRequest(flag string, subType string, approve bool, reason string)

	// GetLoginInfo 获取当前登录的账号信息。
	// 返回: API 响应信息。
	GetLoginInfo() APIResponseInfo

	// GetVersionInfo 获取当前程序的版本信息。
	// 返回: API 响应信息。
	GetVersionInfo() APIResponseInfo

	// GetMsg 获取指定消息 ID 的消息内容。
	// 参数:
	// - msgId: 消息 ID。
	// 返回: API 响应信息。
	GetMsg(msgId int32) APIResponseInfo

	// GetForwardMsg 获取指定转发消息 ID 的消息内容。
	// 参数:
	// - msgId: 转发消息 ID。
	// 返回: API 响应信息。
	GetForwardMsg(msgId string) APIResponseInfo

	// GetGroupList 获取当前登录账号加入的所有群组列表。
	// 返回: API 响应信息。
	GetGroupList() APIResponseInfo

	// GetGroupMemberList 获取指定群组的所有成员列表。
	// 参数:
	// - groupId: 群组 ID。
	// 返回: API 响应信息。
	GetGroupMemberList(groupId int64) APIResponseInfo

	// GetGroupMemberInfo 获取指定群组中指定用户的详细信息，支持缓存控制。
	// 参数:
	// - groupId: 群组 ID。
	// - userId: 用户 ID。
	// - noCache: 是否不使用缓存。
	// 返回: API 响应信息。
	GetGroupMemberInfo(groupId int64, userId int64, noCache bool) APIResponseInfo

	// GetFriendList 获取当前登录账号的所有好友列表。
	// 返回: API 响应信息。
	GetFriendList() APIResponseInfo

	// GetStrangerInfo 获取指定用户的详细信息，支持缓存控制。
	// 参数:
	// - userId: 用户 ID。
	// - noCache: 是否不使用缓存。
	// 返回: API 响应信息。
	GetStrangerInfo(userId int64, noCache bool) APIResponseInfo

	// GetGroupInfo 获取指定群组的详细信息，支持缓存控制。
	// 参数:
	// - groupId: 群组 ID。
	// - noCache: 是否不使用缓存。
	// 返回: API 响应信息。
	GetGroupInfo(groupId int64, noCache bool) APIResponseInfo

	// GetGroupHonorInfo 获取指定群组的荣誉信息。
	// 参数:
	// - groupId: 群组 ID。
	// - Type: 荣誉类型。
	// 返回: API 响应信息。
	GetGroupHonorInfo(groupId int64, Type string) APIResponseInfo

	// GetStatus 获取当前程序的运行状态信息。
	// 返回: API 响应信息。
	GetStatus() APIResponseInfo

	//	// GetCookies 获取指定域名的 Cookies 信息。
	// 参数:
	// - domain: 域名。
	// 返回: API 响应信息。
	GetCookies(domain string) APIResponseInfo

	// GetCSRFToken 获取 CSRF 令牌。
	// 返回: API 响应信息。
	GetCSRFToken() APIResponseInfo

	// GetCredentials 获取指定域名的凭证信息。
	// 参数:
	// - domain: 域名。
	// 返回: API 响应信息。
	GetCredentials(domain string) APIResponseInfo

	// GetImage 获取指定文件的图片信息。
	// 参数:
	// - file: 文件路径或标识符。
	// 返回: API 响应信息。
	GetImage(file string) APIResponseInfo

	// GetRecord 获取指定文件的语音记录信息，支持指定输出格式。
	// 参数:
	// - file: 文件路径或标识符。
	// - outFormat: 输出格式。
	// 返回: API 响应信息。
	GetRecord(file string, outFormat string) APIResponseInfo

	// CanSendImage 检查是否可以发送图片。
	// 返回: API 响应信息。
	CanSendImage() APIResponseInfo

	// CanSendRecord 检查是否可以发送语音记录。
	// 返回: API 响应信息。
	CanSendRecord() APIResponseInfo

	// SetRestart 设置程序在指定延迟后重启。
	// 参数:
	// - delay: 延迟时间，单位为秒。
	SetRestart(delay int32)

	// CleanCache 清理程序的缓存。
	CleanCache()
}

type WindStandardDataBaseAPI interface {
	// SetUserVariable 设置用户变量
	// 参数：
	// - app: 应用信息。
	// - msg: 消息事件信息。
	// - key: 变量名称。
	// - value: 变量值。
	SetUserVariable(app AppInfo, msg MessageEventInfo, key string, value string)

	// SetGroupVariable 设置群组变量
	// 参数：
	// - app: 应用信息。
	// - msg: 消息事件信息。
	// - key: 变量名称。
	// - value: 变量值。
	SetGroupVariable(app AppInfo, msg MessageEventInfo, key string, value string)

	// SetOutUserVariable  [需要master权限]设置其他数据库中的用户变量
	// 参数：
	// - app: 应用信息。
	// - msg: 消息事件信息。
	// - key: 变量名称。
	// - value: 变量值。
	// - datamap: 数据表名称。
	SetOutUserVariable(app AppInfo, datamap string, msg MessageEventInfo, key string, value string)

	// SetOutGroupVariable [需要master权限]设置其他数据库中的群组变量
	// 参数：
	// - app: 应用信息。
	// - msg: 消息事件信息。
	// - key: 变量名称。
	// - value: 变量值。
	// - datamap: 数据表名称。
	SetOutGroupVariable(app AppInfo, datamap string, msg MessageEventInfo, key string, value string)

	// UnsafelySetUserVariable [不安全][需要master权限]设置用户变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// - value: 变量值。
	UnsafelySetUserVariable(app AppInfo, id string, key string, value string)

	// UnsafelySetGroupVariable [不安全][需要master权限]设置群组变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// - value: 变量值。
	UnsafelySetGroupVariable(app AppInfo, id string, key string, value string)

	// UnsafelySetGlobalVariable [不安全][需要master权限]设置全局变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// - value: 变量值。
	UnsafelySetGlobalVariable(app AppInfo, id string, key string, value string)

	// 	UnsafelySetOutUserVariable [不安全][需要master权限]设置其他数据库中的用户变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// - value: 变量值。
	// - datamap: 数据表名称。
	UnsafelySetOutUserVariable(app AppInfo, datamap string, id string, key string, value string)

	// UnsafelySetOutGroupVariable [不安全][需要master权限]设置其他数据库中的群组变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// - value: 变量值。
	// - datamap: 数据表名称。
	UnsafelySetOutGroupVariable(app AppInfo, datamap string, id string, key string, value string)

	// UnsafelySetOutGlobalVariable [不安全][需要master权限]设置其他数据库中的全局变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// - value: 变量值。
	// - datamap: 数据表名称。
	UnsafelySetOutGlobalVariable(app AppInfo, datamap string, id string, key string, value string)

	// GetUserVariable 获取用户变量
	// 参数：
	// - app: 应用信息。
	// - msg: 消息变量名称。
	// - key: 变量名称。
	// 返回: 变量值，是否存在。
	GetUserVariable(app AppInfo, msg MessageEventInfo, key string) (string, bool)

	// GetGroupVariable 获取群组变量
	// 参数：
	// - app: 应用信息。
	// - msg: 消息变量名称。
	// - key: 变量名称。
	// 返回: 变量值，是否存在。
	GetGroupVariable(app AppInfo, msg MessageEventInfo, key string) (string, bool)

	// GetOutUserVariable [需要master权限]获取其他数据库中的用户变量
	// 参数：
	// - app: 应用信息。
	// - msg: 消息变量名称。
	// - key: 变量名称。
	// - datamap：数据表名称。
	// 返回: 变量值，是否存在。
	GetOutUserVariable(app AppInfo, datamap string, msg MessageEventInfo, key string) (string, bool)

	// GetOutGroupVariable [需要master权限]获取其他数据库中的群组变量
	// 参数：
	// - app: 应用信息。
	// - msg: 消息变量名称。
	// - key: 变量名称。
	// - datamap：数据表名称。
	// 返回: 变量值，是否存在。
	GetOutGroupVariable(app AppInfo, datamap string, msg MessageEventInfo, key string) (string, bool)

	// UnsafelyGetUserVariable [不安全][需要master权限]获取用户变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// 返回: 变量值，是否存在。
	UnsafelyGetUserVariable(app AppInfo, id string, key string) (string, bool)

	// UnsafelyGetGroupVariable [不安全][需要master权限]获取群组变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// 返回: 变量值，是否存在。
	UnsafelyGetGroupVariable(app AppInfo, id string, key string) (string, bool)

	// UnsafelyGetGlobalVariable [不安全][需要master权限]获取全局变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// 返回: 变量值，是否存在。
	UnsafelyGetGlobalVariable(app AppInfo, id string, key string) (string, bool)

	// UnsafelyGetOutUserVariable [不安全][需要master权限]获取其他数据库中的用户变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// - datamap：数据表名称。
	// 返回: 变量值，是否存在。
	UnsafelyGetOutUserVariable(app AppInfo, datamap string, id string, key string) (string, bool)

	// UnsafelyGetOutGroupVariable [不安全][需要master权限]获取其他数据库中的群组变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// - datamap：数据表名称。
	// 返回: 变量值，是否存在。
	UnsafelyGetOutGroupVariable(app AppInfo, datamap string, id string, key string) (string, bool)

	// UnsafelyGetOutGlobalVariable [不安全][需要master权限]获取其他数据库中的全局变量
	// 参数：
	// - app: 应用信息。
	// - id: 数据单元 ID。
	// - key: 变量名称。
	// - datamap：数据表名称。
	// 返回: 变量值，是否存在。
	UnsafelyGetOutGlobalVariable(app AppInfo, datamap string, id string, key string) (string, bool)

	// GetIntConfig 获取指定数据单元的整数型配置。
	// 参数:
	// - app: 应用信息。
	// - datamap: 数据单元名称。
	// - key: 配置名称。
	// 返回: 配置值，是否存在。
	GetIntConfig(app AppInfo, datamap string, key string) (int64, bool)

	// GetStringConfig 获取指定数据单元的字符串型配置。
	// 参数:
	// - app: 应用信息。
	// - datamap: 数据单元名称。
	// - key: 配置名称。
	// 返回: 配置值，是否存在。
	GetStringConfig(app AppInfo, datamap string, key string) (string, bool)

	// GetFloatConfig 获取指定数据单元的浮点型配置。
	// 参数:
	// - app: 应用信息。
	// - datamap: 数据单元名称。
	// - key: 配置名称。
	// 返回: 配置值，是否存在。
	GetFloatConfig(app AppInfo, datamap string, key string) (float64, bool)

	// GetIntSliceConfig 获取指定数据单元的整数型切片配置。
	// 参数:
	// - app: 应用信息。
	// - datamap: 数据单元名称。
	// - key: 配置名称。
	// 返回: 配置值，是否存在。
	GetIntSliceConfig(app AppInfo, datamap string, key string) ([]int64, bool)

	// GetStringSliceConfig 获取指定数据单元的字符串型切片配置。
	// 参数:
	// - app: 应用信息。
	// - datamap: 数据单元名称。
	// - key: 配置名称。
	// 返回: 配置值，是否存在。
	GetStringSliceConfig(app AppInfo, datamap string, key string) ([]string, bool)

	// UnsafelyCreatePublicDatamap [不安全][需要master权限]创建公共数据表
	// 参数：
	// - app: 应用信息。
	// - datamapId: 数据表名称。
	UnsafelyCreatePublicDatamap(app AppInfo, datamapId string)
}
