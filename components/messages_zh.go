package components

import "time"

var messagesZH = messages{
	Close: "关闭", Cancel: "取消", Confirm: "确定",

	CarouselRole: "轮播图", SlideRole: "幻灯片", CarouselPrevious: "上一张", CarouselNext: "下一张",
	Pagination: "分页", PagePrevious: "上一页", PageNext: "下一页",
	PrevMonth: "上个月", NextMonth: "下个月",

	Breadcrumb: "面包屑", Search: "搜索", NoResults: "无匹配结果", ShowOptions: "展开选项",
	CommandPalette: "命令面板", Commands: "命令",
	Reactions: "表情回应", Attachments: "附件", SendFailed: "上传失败",
	ContextMenuHint: "右键点击或按 Shift+F10 打开操作菜单。",
	Loading:         "加载中", JumpToEnd: "回到底部",
	OpenMenu: "打开菜单", SideMenu: "侧边栏", Menu: "菜单", Rating: "评分",

	QuestionnaireProgress: "问卷进度", QuestionnaireNext: "下一题", QuestionnaireBack: "上一题",
	QuestionnaireFinish: "完成", QuestionnaireSkip: "跳过",

	StatusOnline: "在线", StatusOffline: "离线", StatusBusy: "忙碌",

	Results:    func(n int) string { return "共 " + itoa(n) + " 条结果" },
	Page:       func(n int) string { return "第 " + itoa(n) + " 页" },
	QuestionOf: func(n, total int) string { return "第 " + itoa(n) + " 题，共 " + itoa(total) + " 题" },
	SlideOf:    func(n, total int) string { return "第 " + itoa(n) + " 张，共 " + itoa(total) + " 张" },
	RatingOf:   func(n, total int) string { return itoa(n) + " 星，共 " + itoa(total) + " 星" },
	Digit:      func(n, total int) string { return "第 " + itoa(n) + " 位，共 " + itoa(total) + " 位" },

	UploadOf:     func(name string) string { return "正在上传 " + name },
	ProcessingOf: func(name string) string { return "正在处理 " + name },

	UploadingStatus:  func(percent int) string { return "正在上传 " + itoa(percent) + "%" },
	ProcessingStatus: "处理中…",

	WeekdayShort: [7]string{"日", "一", "二", "三", "四", "五", "六"},
	WeekdayLong:  zhWeekdays,
	MonthLabel:   func(t time.Time) string { return itoa(t.Year()) + "年" + itoa(int(t.Month())) + "月" },
	DayLabel: func(t time.Time) string {
		return itoa(t.Year()) + "年" + itoa(int(t.Month())) + "月" + itoa(t.Day()) + "日 " + zhWeekdays[t.Weekday()]
	},

	Unread: map[string]string{"other": "（%d 条新消息）"},
}

var zhWeekdays = [7]string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
