import { get, post } from '@/api'
import type { Setting, ApiResponse } from '@/types'
import { getSettingResource, updateSettingResource } from '@/api/resource'

// 获取插件配置
export function getSetting(): Promise<ApiResponse<Setting>> {
  return get(getSettingResource)
}

// 更新插件配置
export function updateSetting(setting: Setting): Promise<ApiResponse<void>> {
  return post(updateSettingResource, setting)
}
