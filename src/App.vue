<template>
  <div class="pmail-telegram-push-settings">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>Telegram 推送设置</span>
          <div style="float: right">
            <el-badge :value="saved ? 1 : ''" class="item">
              <el-button
                type="primary"
                :disabled="!botInfo?.bot_link"
                @click="contactBot"
                style="float: right"
                id="contact_bot_button"
              >
                <TelegramIcon class="icon" />
                联系机器人
              </el-button>
            </el-badge>
          </div>
        </div>
      </template>

      <el-form
        :model="formData"
        :rules="rules"
        label-width="180px"
        label-position="left"
        v-loading="loading"
      >
        <el-form-item label="Telegram 聊天 ID" prop="chat_id">
          <el-input
            v-model="formData.chat_id"
            placeholder="请输入 Telegram Chat ID，置空则禁用推送"
            id="chat_id"
          />
        </el-form-item>

        <el-form-item label="显示邮件内容">
          <el-switch
            v-model="formData.show_content"
            :disabled="formData.chat_id.trim().length === 0"
          />
        </el-form-item>

        <el-form-item label="邮件内容防剧透">
          <el-switch
            v-model="formData.spoiler_content"
            :disabled="formData.chat_id.trim().length === 0"
          />
        </el-form-item>

        <el-form-item label="发送附件">
          <el-switch
            v-model="formData.send_attachments"
            :disabled="formData.chat_id.trim().length === 0"
          />
        </el-form-item>

        <el-form-item label="禁用链接预览">
          <el-switch
            v-model="formData.disable_link_preview"
            :disabled="formData.chat_id.trim().length === 0"
          />
        </el-form-item>
      </el-form>

      <el-form-item>
        <div class="form-item-buttons">
          <el-button
            type="success"
            :disabled="loading || formData.chat_id.trim().length === 0"
            style="margin: 0 auto"
            @click="postTestMessage"
            id="test_message_button"
          >
            <LinkIcon class="icon" />
            测试消息
          </el-button>
          <el-button
            type="primary"
            @click="confirmSubmit"
            :disabled="loading"
            style="margin: 0 auto"
            id="save_button"
          >
            <SaveIcon class="icon" />
            保存设置
          </el-button>
        </div>
      </el-form-item>
      <el-tour v-model="tourOpen" @finish="finishTour" @close="finishTour">
        <el-tour-step
          target="#chat_id"
          title="第一步"
          description="请输入 Chat ID， 置空则禁用推送"
        />
        <el-tour-step
          target="#contact_bot_button"
          title="第二步"
          description="点击联系机器人按钮，与机器人进行一次交互。"
        />
        <el-tour-step
          target="#test_message_button"
          title="第三步"
          description="点击测试消息按钮，测试是否能成功发送消息。"
        />
        <el-tour-step
          target="#save_button"
          title="第四步"
          description="点击保存设置按钮，保存当前配置。"
        />
        <template #indicators="{ current, total }">
          <span>{{ current + 1 }} / {{ total }}</span>
        </template>
      </el-tour>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import resizeIframeHeight from '@/utils/resize'
import { ElMessage, ElMessageBox } from 'element-plus'
import SaveIcon from '@/assets/icons/save-line.svg'
import TelegramIcon from '@/assets/icons/telegram.svg'
import LinkIcon from '@/assets/icons/link.svg'
import type { BotInfo, Setting } from '@/types'
import { getBotInfo, testMessage } from '@/api/bot'
import { getSetting, updateSetting } from '@/api/setting'

const saved = ref(false)
const loading = ref(false)
const botInfo = ref<BotInfo>({
  username: '',
  bot_link: '',
})

// 获取机器人信息
const fetchBotInfo = () => {
  loading.value = true
  getBotInfo()
    .then((response) => {
      if (response.data) {
        botInfo.value = response.data
      }
    })
    .catch((error) => {
      ElMessage.error('获取机器人信息失败')
      console.error(error)
    })
    .finally(() => {
      loading.value = false
    })
}

// 联系机器人
const contactBot = () => {
  if (!botInfo.value?.bot_link) {
    ElMessage.error('机器人链接为空')
    return
  }
  saved.value = false
  window.open(botInfo.value.bot_link, '_blank')
}

const formData = ref<Setting>({
  chat_id: '',
  show_content: true,
  spoiler_content: true,
  send_attachments: true,
  disable_link_preview: true,
})
const rules = {
  chat_id: [
    {
      required: false,
      message: '请输入 Chat ID， 置空则禁用推送',
      trigger: 'blur',
      whitespace: false,
    },
  ],
}
// 获取插件配置
const fetchSetting = () => {
  loading.value = true
  getSetting()
    .then((response) => {
      if (response.data) {
        formData.value = response.data
      }
    })
    .catch((error) => {
      ElMessage.error('获取设置信息失败')
      console.error(error)
    })
    .finally(() => {
      loading.value = false
    })
}

// 保存插件配置
const saveSetting = () => {
  loading.value = true
  updateSetting(formData.value || {})
    .then(() => {
      saved.value = true
      ElMessage.success('设置已保存')
    })
    .catch((error) => {
      ElMessage.error('保存设置失败')
      console.error(error)
    })
    .finally(() => {
      loading.value = false
    })
}

// 确认保存设置
const confirmSubmit = () => {
  ElMessageBox.confirm('确认保存设置吗？', '保存设置', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    saveSetting()
  })
}

const postTestMessage = () => {
  loading.value = true
  if (!formData.value.chat_id) {
    ElMessage.error('请输入 Chat ID')
    return
  }
  testMessage(formData.value)
    .then((response) => {
      if (response.code === 0) {
        ElMessage.success(response.message)
        saved.value = true
      } else {
        ElMessage.error(response.message)
      }
    })
    .catch((error) => {
      ElMessage.error('测试消息失败')
      console.error(error)
    })
    .finally(() => {
      loading.value = false
    })
}

const tourOpen = ref(false)
const finishTour = () => {
  localStorage.setItem('hasShownTour', 'true')
  console.log('tour finished')
}

onMounted(() => {
  const hasShown = localStorage.getItem('hasShownTour') === 'true'
  if (!hasShown) {
    tourOpen.value = true
  }
  resizeIframeHeight()
  fetchBotInfo()
  fetchSetting()
})
</script>

<style scoped>
.pmail-telegram-push-settings {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.box-card {
  margin-top: 20px;
}

.card-header {
  font-size: 18px;
  font-weight: bold;
}

.icon {
  margin-right: 5px;
  width: 20px;
  height: 20px;
  vertical-align: middle;
  color: white;
}

.form-item-buttons {
  margin-top: 20px;
  margin: 0 auto;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
}
</style>
