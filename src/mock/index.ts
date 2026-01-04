import Mock from 'mockjs'
import type { ApiResponse, BotInfo, Setting } from '@/types'
import {
  pluginSettingsResource,
  getBotInfoResource,
  getSettingResource,
  updateSettingResource,
  testMessageResource,
} from '@/api/resource'

// 模拟 Telegram 机器人信息
const Random = Mock.Random

// 模拟超时时长
Mock.setup({
  timeout: '500-1000',
})

// 模拟获取 Telegram 机器人信息
export const getBotInfoMock = Mock.mock(
  `${pluginSettingsResource}/${getBotInfoResource}`,
  'get',
  () => {
    const botInfo: ApiResponse<BotInfo> = {
      code: 0,
      message: '获取 Telegram 机器人信息成功',
      data: {
        username: Random.string('lower', 10),
        bot_link: Random.string('lower', 10),
      },
    }
    // console.log('getBotInfo mock response', botInfo)
    return botInfo
  },
)

// 模拟获取设置信息
export const getSettingMock = Mock.mock(
  `${pluginSettingsResource}/${getSettingResource}`,
  'get',
  () => {
    const setting: ApiResponse<Setting> = {
      code: 0,
      message: '获取设置信息成功',
      data: {
        chat_id: Random.string('number', 9),
        show_content: Random.boolean(),
        spoiler_content: Random.boolean(),
        send_attachments: Random.boolean(),
        disable_link_preview: Random.boolean(),
      },
    }
    // console.log('getSetting mock response', setting)
    return setting
  },
)

// 模拟更新设置信息
export const updateSettingMock = Mock.mock(
  `${pluginSettingsResource}/${updateSettingResource}`,
  'post',
  (_) => {
    // console.log('updateSetting options', options)
    const code = Random.integer(0, -1) // 0 表示成功，其他值表示失败
    const message = code === 0 ? '更新设置信息成功' : '更新设置信息失败'
    const setting: ApiResponse<Setting> = {
      code,
      message,
    }
    // console.log('updateSetting mock response', setting)
    return setting
  },
)

// 模拟测试消息
export const testMessageMock = Mock.mock(
  `${pluginSettingsResource}/${testMessageResource}`,
  'post',
  (_) => {
    // console.log('testMessage options', options)
    const code = Random.integer(0, -1) // 0 表示成功，其他值表示失败
    const message = code === 0 ? '测试消息成功' : '测试消息失败'
    const setting: ApiResponse<void> = {
      code,
      message,
    }
    // console.log('testMessage mock response', setting)
    return setting
  },
)

export default { getBotInfoMock, getSettingMock, updateSettingMock, testMessageMock }
