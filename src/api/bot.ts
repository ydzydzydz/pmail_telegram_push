import { get } from '@/api'
import type { BotInfo } from '@/types'
import { getBotInfoResource } from '@/api/resource'

// 获取机器人信息
export const getBotInfo = () => {
  return get<BotInfo>(getBotInfoResource)
}
