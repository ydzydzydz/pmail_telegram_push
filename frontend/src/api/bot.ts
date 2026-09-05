import { get, post } from '@/api'
import type { BotInfo, Setting, ApiResponse } from '@/types'
import { getBotInfoResource, testMessageResource } from '@/api/resource'

// 获取机器人信息
export const getBotInfo = () => {
  return get<BotInfo>(getBotInfoResource)
}

export function testMessage(setting: Setting): Promise<ApiResponse<void>> {
  return post(testMessageResource, setting)
}
