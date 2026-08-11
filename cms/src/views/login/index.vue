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
              :prefix-icon="UserIcon"
              size="large"
          />
        </el-form-item>
        <el-form-item prop="pwd">
          <el-input
              v-model="loginForm.pwd"
              placeholder="密码"
              :prefix-icon="LockIcon"
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
import {onMounted, reactive, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {Platform, Lock as LockIcon, User as UserIcon} from '@element-plus/icons-vue'
import {authApi} from '@/api'
import {ElMessage} from 'element-plus'
import type {LoginRes} from '@/types/api'
import {clearPermissions} from '@/utils/permission'
import {resolveAccessiblePath} from '@/utils/accessible-route'
import {
    getSavedCredentials,
    isAuthenticated,
    saveLoginCredentials,
    setAuthSession,
} from '@/utils/auth'

interface LoginForm {
  userName: string
  pwd: string
}

const router = useRouter()
const route = useRoute()
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

const resolveRedirectPath = () => {
  const redirect = route.query.redirect
  if (typeof redirect === 'string' && redirect.startsWith('/')) {
    return resolveAccessiblePath(redirect)
  }
  return resolveAccessiblePath('/dashboard')
}

onMounted(() => {
  if (isAuthenticated()) {
    router.replace(resolveRedirectPath())
    return
  }

  const saved = getSavedCredentials()
  loginForm.userName = saved.userName
  loginForm.pwd = saved.pwd
})

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
            setAuthSession({
              token: res.token,
              authId: res.authId.toString(),
              admin: res.admin,
              modules: res.modules || [],
            })
            saveLoginCredentials(loginForm.userName, loginForm.pwd)

            ElMessage.success('登录成功')
            router.replace(resolveRedirectPath())
          })
          .catch(err => {
            console.error('Login error:', err)
            clearPermissions()
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