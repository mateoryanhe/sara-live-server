<template>
  <div class="login-page">
    <div class="login-visual">
      <img alt="Sara Live CMS" class="login-illustration" src="@/assets/login-illustration.svg"/>
      <div class="visual-copy">
        <h1>{{ t('login.title') }}</h1>
        <p>{{ t('login.subtitle') }}</p>
        <ul class="feature-list">
          <li>{{ t('login.featureLive') }}</li>
          <li>{{ t('login.featureData') }}</li>
          <li>{{ t('login.featureSecure') }}</li>
        </ul>
      </div>
    </div>

    <div class="login-panel">
      <div class="panel-top">
        <LanguageSwitcher compact/>
      </div>

      <div class="panel-body">
        <div class="brand-badge">
          <el-icon :size="28">
            <Platform/>
          </el-icon>
        </div>
        <h2 class="panel-title">{{ t('login.welcome') }}</h2>
        <p class="panel-subtitle">{{ t('login.hint') }}</p>

        <el-form
            ref="loginFormRef"
            :model="loginForm"
            :rules="loginRules"
            class="login-form"
            size="large"
            @keyup.enter="handleLogin"
        >
          <el-form-item prop="userName">
            <el-input
                v-model="loginForm.userName"
                :placeholder="t('login.userNamePlaceholder')"
                :prefix-icon="UserIcon"
            />
          </el-form-item>
          <el-form-item prop="pwd">
            <el-input
                v-model="loginForm.pwd"
                :placeholder="t('login.passwordPlaceholder')"
                :prefix-icon="LockIcon"
                show-password
                type="password"
            />
          </el-form-item>
          <el-form-item>
            <el-button
                :loading="loading"
                class="login-button"
                type="primary"
                @click="handleLogin"
            >
              {{ t('login.submit') }}
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {Platform, Lock as LockIcon, User as UserIcon} from '@element-plus/icons-vue'
import {useI18n} from 'vue-i18n'
import type {FormInstance, FormRules} from 'element-plus'
import {ElMessage} from 'element-plus'
import {authApi} from '@/api'
import type {LoginRes} from '@/types/api'
import LanguageSwitcher from '@/components/LanguageSwitcher.vue'
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
const {t} = useI18n()
const loading = ref(false)
const loginForm = reactive<LoginForm>({
  userName: '',
  pwd: '',
})

const loginFormRef = ref<FormInstance>()

const loginRules = computed<FormRules>(() => ({
  userName: [{required: true, message: t('login.userNameRequired'), trigger: 'blur'}],
  pwd: [
    {required: true, message: t('login.passwordRequired'), trigger: 'blur'},
    {min: 6, message: t('login.passwordMin'), trigger: 'blur'},
  ],
}))

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
  if (!loginFormRef.value) {
    return
  }

  await loginFormRef.value.validate((valid) => {
    if (!valid) {
      return
    }
    loading.value = true
    authApi.cmsLogin({
      userName: loginForm.userName,
      pwd: loginForm.pwd,
    })
        .then((res: LoginRes) => {
          setAuthSession({
            token: res.token,
            authId: res.authId.toString(),
            admin: res.admin,
            superAdmin: res.superAdmin,
            modules: res.modules || [],
          })
          saveLoginCredentials(loginForm.userName, loginForm.pwd)
          ElMessage.success(t('login.success'))
          router.replace(resolveRedirectPath())
        })
        .catch(err => {
          console.error('Login error:', err)
          clearPermissions()
        })
        .finally(() => {
          loading.value = false
        })
  })
}
</script>

<style scoped>
.login-page {
  display: grid;
  grid-template-columns: minmax(0, 1.15fr) minmax(420px, 480px);
  min-height: 100vh;
  background: #f4f6fb;
}

.login-visual {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 48px 56px;
  overflow: hidden;
  background: linear-gradient(145deg, #312e81 0%, #4338ca 45%, #7c3aed 100%);
  color: #fff;
}

.login-illustration {
  width: min(100%, 520px);
  margin-bottom: 32px;
  filter: drop-shadow(0 20px 40px rgba(15, 23, 42, 0.25));
}

.visual-copy h1 {
  margin: 0 0 12px;
  font-size: 36px;
  font-weight: 700;
  letter-spacing: -0.02em;
}

.visual-copy p {
  margin: 0 0 24px;
  max-width: 420px;
  font-size: 16px;
  line-height: 1.6;
  color: rgba(255, 255, 255, 0.88);
}

.feature-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.feature-list li {
  padding: 8px 14px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(6px);
  font-size: 13px;
}

.login-panel {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 32px 40px;
  background: #fff;
  box-shadow: -8px 0 32px rgba(15, 23, 42, 0.06);
}

.panel-top {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 24px;
}

.panel-body {
  width: 100%;
  max-width: 380px;
  margin: 0 auto;
}

.brand-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  margin-bottom: 20px;
  border-radius: 16px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #fff;
  box-shadow: 0 10px 24px rgba(99, 102, 241, 0.35);
}

.panel-title {
  margin: 0 0 8px;
  font-size: 28px;
  font-weight: 700;
  color: #111827;
}

.panel-subtitle {
  margin: 0 0 28px;
  color: #6b7280;
  font-size: 14px;
}

.login-form {
  margin-top: 8px;
}

.login-button {
  width: 100%;
  height: 44px;
  margin-top: 4px;
  border: none;
  background: linear-gradient(90deg, #6366f1, #8b5cf6);
}

@media (max-width: 960px) {
  .login-page {
    grid-template-columns: 1fr;
  }

  .login-visual {
    display: none;
  }

  .login-panel {
    min-height: 100vh;
    box-shadow: none;
  }
}
</style>
