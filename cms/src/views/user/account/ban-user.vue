<template>
  <div class="ban-user-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isAnchorBan ? t('pages.banUser.banAnchorTitle') : t('pages.banUser.banUserTitle') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item :label="t('common.userId')" prop="userId">
          <el-input v-model="form.userId" disabled :placeholder="t('common.pleaseEnter')"/>
        </el-form-item>

        <el-form-item :label="t('pages.banUser.openId')" prop="openId">
          <el-input v-model="form.openId" disabled :placeholder="t('pages.banUser.openId')"/>
        </el-form-item>

        <el-form-item :label="t('pages.banUser.ipAddress')" prop="ip">
          <el-input v-model="form.ip" disabled :placeholder="t('pages.banUser.ipAddress')"/>
        </el-form-item>

        <el-form-item :label="t('pages.banUser.channel')" prop="channel">
          <el-input v-model="form.channel" disabled :placeholder="t('pages.banUser.channel')"/>
        </el-form-item>

        <el-form-item :label="t('pages.banUser.banUntil')" prop="banApplyTime">
          <el-date-picker
              v-model="form.banApplyTime"
              :disabled-date="disabledDate"
              format="YYYY-MM-DD HH:mm:ss"
              :placeholder="t('pages.banUser.selectBanUntil')"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submitForm">
            {{ isAnchorBan ? t('pages.banUser.confirmBanAnchor') : t('pages.banUser.confirmBanUser') }}
          </el-button>
          <el-button @click="goBack">{{ t('pages.banUser.back') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {ElForm, ElMessage, type FormRules} from 'element-plus'
import {accountApi} from '@/api'
import type {BanAnchorReq, BanReq} from '@/types/api.ts'

const {t} = useI18n()
const router = useRouter()
const route = useRoute()

const isAnchorBan = computed(() => route.query.type === 'anchor')
const returnPath = computed(() => {
  const path = route.query.returnPath
  if (typeof path === 'string' && path.startsWith('/')) {
    return path
  }
  return '/user/account/user-list'
})

const formRef = ref<InstanceType<typeof ElForm>>()

const form = reactive({
  userId: '',
  openId: '',
  ip: '',
  channel: 0,
  banApplyTime: ''
})

const rules = computed<FormRules>(() => ({
  userId: [
    {required: true, message: t('pages.banUser.userIdRequired'), trigger: 'blur'}
  ],
  banApplyTime: [
    {required: true, message: t('pages.banUser.banUntilRequired'), trigger: 'change'}
  ]
}))

const disabledDate = (time: Date) => {
  return time.getTime() < Date.now()
}

onMounted(() => {
  const userData = route.query
  if (userData && userData.id) {
    form.userId = userData.id as string
    form.openId = userData.openId as string || ''
    form.ip = userData.ip as string || ''
    form.channel = Number(userData.channel) || 0
    form.banApplyTime = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString()
  } else {
    router.push(returnPath.value)
  }
})

const submitForm = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      try {
        if (isAnchorBan.value) {
          const banData: BanAnchorReq = {
            accountId: form.userId,
            banApplyTime: form.banApplyTime,
          }
          const response = await accountApi.banAnchor(banData)
          if (response) {
            ElMessage.success(t('pages.banUser.banAnchorSuccess'))
            router.push(returnPath.value + '?refresh=' + Date.now())
          } else {
            ElMessage.error(t('pages.banUser.banAnchorFailed'))
          }
        } else {
          const banData: BanReq = {
            accountId: form.userId,
            openId: form.openId,
            channel: form.channel,
            banApplyTime: form.banApplyTime,
          }
          const response = await accountApi.ban(banData)
          if (response) {
            ElMessage.success(t('pages.banUser.banUserSuccess'))
            router.push(returnPath.value + '?refresh=' + Date.now())
          } else {
            ElMessage.error(t('pages.banUser.banUserFailed'))
          }
        }
      } catch (error) {
        console.error('Ban request failed:', error)
        ElMessage.error(t('pages.banUser.banRequestFailed'))
      }
    } else {
      ElMessage.error(t('pages.banUser.invalidForm'))
    }
  })
}

const goBack = () => {
  router.push(returnPath.value)
}
</script>

<style scoped>
.ban-user-container {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.el-form {
  max-width: 600px;
  margin-top: 20px;
}
</style>
