package components

var messagesZH = messages{
	Close: "关闭", Cancel: "取消", Confirm: "确认",

	CarouselRole: "轮播", CarouselPrevious: "上一张", CarouselNext: "下一张",
	Pagination: "分页", PagePrevious: "上一页", PageNext: "下一页",
	PrevMonth: "上个月", NextMonth: "下个月",

	Breadcrumb: "面包屑导航", Search: "搜索", NoResults: "无结果", ShowOptions: "显示选项",
	CommandPalette: "命令面板", Commands: "命令",
	Reactions: "表情回应", Attachments: "附件", SendFailed: "上传失败。",
	ContextMenuHint: "右键单击或按 Shift+F10 查看操作。",
	Loading:         "加载中", JumpToEnd: "跳到末尾",
	OpenMenu: "打开菜单", SideMenu: "侧边栏", Menu: "菜单", Rating: "评分",

	QuestionnaireProgress: "问卷进度", QuestionnaireNext: "下一题", QuestionnaireBack: "上一题",

	StatusOnline: "在线", StatusOffline: "离线", StatusBusy: "忙碌",

	Results:    func(n int) string { return itoa(n) + " 个结果" },
	Page:       func(n int) string { return "第 " + itoa(n) + " 页" },
	QuestionOf: func(n, total int) string { return "第 " + itoa(n) + " 题，共 " + itoa(total) + " 题" },
	OfTotal:    func(n, total int) string { return itoa(n) + " / " + itoa(total) },
	Digit:      func(n, total int) string { return "第 " + itoa(n) + " 位，共 " + itoa(total) + " 位" },

	UploadOf:     func(name string) string { return "正在上传 " + name },
	ProcessingOf: func(name string) string { return "正在处理 " + name },

	Unread: map[string]string{"other": "(%d 条新消息)"},
}
