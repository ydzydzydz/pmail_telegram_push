// 定义 API 响应体结构
export interface ApiResponse<T> {
  code: number
  message: string
  data?: T
}

// 定义 Telegram 机器人信息结构
export interface BotInfo {
  username: string
  bot_link: string
}

// 定义插件设置结构
export interface Setting {
  chat_id: string
  show_content: boolean
  spoiler_content: boolean
  send_attachments: boolean
  disable_link_preview: boolean
}
