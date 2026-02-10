package enum

type ProfileType int8 // 资料类型

const (
	// 基本资料类型
	SingleLineText ProfileType = iota + 1 // 单行文本
	MultiLineText                         // 多行文本
	SingleChoice                          // 单选
	MultipleChoice                        // 多选
	Date                                  // 日期
	Tag                                   // 标签
	Image                                 // 图片
	Video                                 // 视频
	Audio                                 // 音频
	Number                                // 数字（例如年龄、身高）
	Boolean                               // 布尔值（例如是否已婚）
	File                                  // 文件（例如文档、PDF）
	SecretCard                            // 神秘卡片（隐藏的兴趣或小秘密）
	Achievements                          // 成就徽章
	Mood                                  // 心情动态
	Preference                            // 情感偏好（约会方式、恋爱观等）
	VirtualGift                           // 虚拟礼物展示
	Questionnaire                         // 问卷、调查或测验
)

func (ft ProfileType) String() string {
	switch ft {
	case SingleLineText:
		return "SingleLineText"
	case MultiLineText:
		return "MultiLineText"
	case SingleChoice:
		return "SingleChoice"
	case MultipleChoice:
		return "MultipleChoice"
	case Date:
		return "Date"
	case Tag:
		return "Tag"
	case Image:
		return "Image"
	case Video:
		return "Video"
	case Audio:
		return "Audio"
	case Number:
		return "Number"
	case Boolean:
		return "Boolean"
	case File:
		return "File"
	case SecretCard:
		return "SecretCard"
	case Achievements:
		return "Achievements"
	case Mood:
		return "Mood"
	case Preference:
		return "Preference"
	case VirtualGift:
		return "VirtualGift"
	case Questionnaire:
		return "Questionnaire"
	default:
		return "Unknown"
	}
}

type AccessMethod int8 // 解锁方式

const (
	// 基于互动积分解锁
	InteractionPoints AccessMethod = iota // 0
	// 基于关系等级解锁
	RelationshipLevel // 1
	// 基于行为、任务解锁
	BehaviorAndTasks // 2
	// 基于时间解锁
	TimeBased // 3
	// 基于双向同意解锁
	MutualConsent // 4
)
