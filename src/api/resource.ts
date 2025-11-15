// 插件设置资源
const settingResource = '/api/plugin/settings'
// 插件名称
const pluginName = 'pmail_telegram_push'
// 插件配置资源
export const pluginSettingsResource = `${settingResource}/${pluginName}`

// 获取插件配置资源
export const getSettingResource = 'getSetting'
// 更新插件配置资源
export const updateSettingResource = 'updateSetting'
// 获取机器人信息资源
export const getBotInfoResource = 'getBotInfo'
