package utils

import (
	"errors"
)

var (
	ErrParameterWrong           = errors.New("参数错误")
	ErrInternalServer           = errors.New("网络错误")
	ErrOperationFailed          = errors.New("操作未成功完成")
	ErrCreateTeamFailed         = errors.New("创建团队失败")
	ErrAddMemberOperationFailed = errors.New("添加成员失败")
	ErrFormatWrong              = errors.New("格式错误")
	ErrUserHaveExist            = errors.New("账户已经存在")
	ErrCopyFail                 = errors.New("数据迁移失败")
	ErrHaveInTeam               = errors.New("用户已经加入其他团队")
	ErrUserNotFound             = errors.New("用户不存在")
	ErrTeamNotFound             = errors.New("团队不存在")
	ErrCreateMsgFailed          = errors.New("创建信息失败")
	ErrDelTeamFailed            = errors.New("解散团队失败")
	ErrMsgNotFound              = errors.New("信息不存在")
	ErrHavenoPower              = errors.New("没有权限")
	ErrSubmitFailed             = errors.New("报名人数不符合标准")
)

const (
	SuccessCode             = 200
	ParameterErrorCode      = 400 // BadRequest
	UnAuthorizedCode        = 401
	NotFoundCode            = 404
	OperationFailedCode     = 406
	HaveExistCode           = 407
	InternalServerErrorCode = 500
)

const (
	SuccessMsg             = "Accepted"
	ParameterErrorMsg      = "参数错误"
	InternalServerErrorMsg = "Internal server error"
	UserHaveExistMsg       = "邮箱已被注册"
	OperationFailedMsg     = "操作未成功完成"
	UserNotFoundMsg        = "账户不存在或密码错误"
	TeamNotFoundMsg        = "团队不存在或密码错误"
	FormatWrongMsg         = "格式错误"
	HaveInTeamMsg          = "用户已经加入过团队"
	NotFoundMsg            = "找不到对象"
	CopyFailMsg            = "数据迁移失败"
	AddMemberFailMsg       = "添加成员失败"
	CreateMsgFailedMsg     = "创建信息失败"
	DelTeamFailedMsg       = "团队解散失败"
	FileSaveFailedMsg      = "文件保存失败"
	UnAuthorizedMsg        = "未认证"
	HavenoPowerMsg         = "没有权限"
	SubmitFailedMsg        = "提交失败,团队人数必须在4-6人之间!"
)
