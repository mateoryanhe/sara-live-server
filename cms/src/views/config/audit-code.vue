<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.AuditCodeCfgManagement') }}</span>
        </div>
      </template>

      <el-alert
          :closable="false"
          :description="t('pages.auditCode.tip')"
          show-icon
          type="info"
          class="notice"
      />

      <el-form class="cfg-form" label-width="120px">
        <el-form-item :label="t('pages.auditCode.codeLength')">
          <el-input-number
              v-model="codeLength"
              :min="4"
              :max="32"
              :step="1"
              controls-position="right"
          />
          <span class="form-tip inline-tip">{{ t('pages.auditCode.codeLengthTip') }}</span>
        </el-form-item>
        <el-form-item :label="t('pages.auditCode.currentCode')">
          <div class="code-row">
            <el-input
                :model-value="displayCode"
                readonly
                class="code-input"
            />
            <el-button :disabled="!currentCode" @click="handleCopy">
              {{ t('pages.auditCode.copy') }}
            </el-button>
            <el-button :loading="generating" type="primary" @click="handleGenerate">
              {{ t('pages.auditCode.generate') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item v-if="updatedAt" :label="t('pages.auditCode.lastUpdated')">
          <span>{{ updatedAt }}</span>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, onMounted, ref} from 'vue'
import {ElMessage, ElMessageBox} from 'element-plus'
import {auditCodeApi} from '@/api/modules/audit-code'

const {t} = useI18n()
const loading = ref(false)
const generating = ref(false)
const currentCode = ref('')
const updatedAt = ref('')
const codeLength = ref(8)

const displayCode = computed(() => currentCode.value || t('pages.auditCode.emptyCode'))

const fetchCfg = async () => {
  loading.value = true
  try {
    const res = await auditCodeApi.getAuditCodeCfg()
    currentCode.value = res.cfg?.code || ''
    updatedAt.value = res.cfg?.updatedAt || ''
    // 位数仅作下次生成参数，固定默认 8，不跟当前码长度走
  } catch (error) {
    console.error('fetch audit code failed:', error)
    ElMessage.error(t('pages.auditCode.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleCopy = async () => {
  if (!currentCode.value) {
    return
  }
  try {
    await navigator.clipboard.writeText(currentCode.value)
    ElMessage.success(t('pages.auditCode.copySuccess'))
  } catch (error) {
    console.error('copy audit code failed:', error)
    ElMessage.error(t('pages.auditCode.copyFailed'))
  }
}

const handleGenerate = async () => {
  const length = Number(codeLength.value) || 8
  if (length < 4 || length > 32) {
    ElMessage.warning(t('pages.auditCode.codeLengthInvalid'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.auditCode.generateConfirm'),
        t('pages.auditCode.generateTitle'),
        {type: 'warning'},
    )
  } catch {
    return
  }
  generating.value = true
  try {
    const res = await auditCodeApi.generateAuditCodeCfg(length)
    if (res?.success && res.cfg?.code) {
      currentCode.value = res.cfg.code
      updatedAt.value = res.cfg.updatedAt || ''
      ElMessage.success(t('pages.auditCode.generateSuccess'))
    } else {
      ElMessage.error(t('pages.auditCode.generateFailed'))
    }
  } catch (error) {
    console.error('generate audit code failed:', error)
    ElMessage.error(t('pages.auditCode.generateFailed'))
  } finally {
    generating.value = false
  }
}

onMounted(() => {
  fetchCfg()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.notice {
  margin-bottom: 20px;
}

.cfg-form {
  max-width: 860px;
}

.code-row {
  display: flex;
  gap: 8px;
  width: 100%;
}

.code-input {
  flex: 1;
}

.code-input :deep(input) {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
}

.inline-tip {
  margin-left: 12px;
  color: #909399;
  font-size: 13px;
}
</style>
