<template>
  <div class="ban-user-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>封号用户</span>
        </div>
      </template>

      <el-form ref="formRef" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="用户ID" prop="userId">
          <el-input v-model="form.userId" disabled placeholder="请输入用户ID"/>
        </el-form-item>

        <el-form-item label="OpenId" prop="openId">
          <el-input v-model="form.openId" disabled placeholder="OpenId"/>
        </el-form-item>

        <el-form-item label="IP地址" prop="ip">
          <el-input v-model="form.ip" disabled placeholder="IP地址"/>
        </el-form-item>

        <el-form-item label="渠道" prop="channel">
          <el-input v-model="form.channel" disabled placeholder="渠道"/>
        </el-form-item>

        <el-form-item label="封号截止时间" prop="banApplyTime">
          <el-date-picker
              v-model="form.banApplyTime"
              :disabled-date="disabledDate"
              format="YYYY-MM-DD HH:mm:ss"
              placeholder="选择封号截止时间"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>

        <el-form-item>
          <el-button :loading="loading" type="primary" @click="submitForm">确认封号</el-button>
          <el-button @click="goBack">返回</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {ElForm, ElMessage} from 'element-plus'
import {accountApi} from '@/api'
import type {BanReq} from '@/types/api'

const router = useRouter()
const route = useRoute()

// 表单引用
const formRef = ref<InstanceType<typeof ElForm>>()

// 加载状态
const loading = ref(false)

// 表单数据
const form = reactive({
  userId: '',
  openId: '',
  ip: '',
  channel: 0,
  banApplyTime: ''
})

// 表单验证规则
const rules = {
  userId: [
    {required: true, message: '请输入用户ID', trigger: 'blur'}
  ],
  banApplyTime: [
    {required: true, message: '请选择封号截止时间', trigger: 'change'}
  ]
}

// 禁用过去的日期
const disabledDate = (time: Date) => {
  return time.getTime() < Date.now()
}

// 初始化数据
onMounted(() => {
  // 从路由参数获取用户数据
  const userData = route.query
  if (userData && userData.id) {
    form.userId = userData.id as string
    form.openId = userData.openId as string || ''
    form.ip = userData.ip as string || ''
    form.channel = Number(userData.channel) || 0
    // 默认设置为7天后
    form.banApplyTime = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString()
  } else {
    // 如果没有传入用户数据，返回用户列表
    router.push('/account/user-list')
  }
})

// 提交封号表单
const submitForm = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid: boolean) => {
    if (valid) {
      loading.value = true
      try {
        const banData: BanReq = {
          accountId: form.userId,
          banApplyTime: form.banApplyTime
        }

        const response = await accountApi.ban(banData)
        if (response) {
          ElMessage.success('封号成功')
          // 封号成功后跳转回用户列表并强制刷新
          router.push('/account/user-list?refresh=' + Date.now())
        } else {
          ElMessage.error('封号失败')
        }
      } catch (error) {
        console.error('封号请求失败:', error)
        ElMessage.error('封号请求失败')
      } finally {
        loading.value = false
      }
    } else {
      ElMessage.error('请填写正确的表单信息')
    }
  })
}

// 返回用户列表
const goBack = () => {
  router.push('/account/user-list')
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