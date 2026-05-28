<template>
  <div class="login-container">
    <el-card class="login-card">
      <div class="logo-container">
        <div class="logo-icon">
          <el-icon class="icon" size="60">
            <Platform/>
          </el-icon>
        </div>
      </div>
      <h2 class="title"></h2>
      <p class="subtitle">管理系统</p>
      <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="loginRules"
          class="login-form"
      >
        <el-form-item prop="userName">
          <el-input
              v-model="loginForm.userName"
              placeholder="用户名"
              prefix-icon="User"
              size="large"
          />
        </el-form-item>
        <el-form-item prop="pwd">
          <el-input
              v-model="loginForm.pwd"
              placeholder="密码"
              prefix-icon="Lock"
              size="large"
              type="password"
              @keyup.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item>
          <el-button
              :loading="loading"
              class="login-button"
              size="large"
              type="primary"
              @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {reactive, ref} from 'vue'
import {useRouter} from 'vue-router'
import {Platform} from '@element-plus/icons-vue'
import {authApi} from '@/api'
import {ElMessage} from 'element-plus'
import type {LoginRes} from '@/types/api'
import {clearPermissions, setUserPermissions} from '@/utils/permission'

interface LoginForm {
  userName: string
  pwd: string
}

const router = useRouter()
const loading = ref(false)
const loginForm = reactive<LoginForm>({
  userName: '',
  pwd: ''
})

const loginFormRef = ref()

const loginRules = {
  userName: [
    {required: true, message: '请输入用户名', trigger: 'blur'}
  ],
  pwd: [
    {required: true, message: '请输入密码', trigger: 'blur'},
    {min: 6, message: '密码长度不能少于6位', trigger: 'blur'}
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  await loginFormRef.value.validate((valid: boolean) => {
    if (valid) {
      loading.value = true
      authApi.cmsLogin({
        userName: loginForm.userName,
        pwd: loginForm.pwd
      })
          .then((res: LoginRes) => {
            // 登录成功，保存token和authId
            localStorage.setItem('token', res.token)
            localStorage.setItem('authId', res.authId.toString())

            // 设置用户权限信息
            setUserPermissions(res.modules || [], res.admin)

            ElMessage.success('登录成功')
            // 登录成功后跳转到用户列表页面，而不是仪表盘
            router.push('/dashboard')
          })
          .catch(err => {
            console.error('Login error:', err)
            // 登录失败时清除权限信息
            clearPermissions()
            //ElMessage.error('登录失败: ' + (err.message || '网络错误'))
          })
          .finally(() => {
            loading.value = false
          })
    }
  })
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f0f2f5;
  background-image: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 20px;
}

.login-card {
  width: 400px;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  border: none;
  background-color: white;
}

.logo-container {
  text-align: center;
  margin-bottom: 10px;
}

.logo-icon {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 80px;
  height: 80px;
  margin: 0 auto 10px;
  background: linear-gradient(135deg, #409eff 0%, #4a9eff 100%);
  border-radius: 50%;
}

.icon {
  color: white;
}

.title {
  text-align: center;
  margin: 0 0 5px 0;
  color: #303133;
  font-size: 24px;
  font-weight: 600;
}

.subtitle {
  text-align: center;
  margin: 0 0 30px 0;
  color: #909399;
  font-size: 14px;
}

.login-form {
  margin-top: 20px;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  letter-spacing: 1px;
  background: linear-gradient(90deg, #409eff, #4a9eff);
  border: none;
}
</style>