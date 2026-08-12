<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.RandomNicknameManagement') }}</span>
          <el-tag :type="cfg.useDB ? 'success' : 'info'" size="small">
            {{ cfg.useDB ? t('pages.randomNicknameCfg.useDbPool') : t('pages.randomNicknameCfg.useBuiltinPool') }}
          </el-tag>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          :title="t('pages.randomNicknameCfg.tipTitle')"
          type="info"
      >
        <p>{{ t('pages.randomNicknameCfg.tipLine1') }}</p>
        <p>{{ t('pages.randomNicknameCfg.tipLine2') }}</p>
      </el-alert>

      <el-table :data="cfg.langs" border class="stat-table" :empty-text="t('pages.randomNicknameCfg.noData')">
        <el-table-column :label="t('pages.randomNicknameCfg.language')" min-width="140" prop="langLabel"/>
        <el-table-column :label="t('pages.randomNicknameCfg.code')" prop="langCode" width="80"/>
        <el-table-column :label="t('pages.randomNicknameCfg.count')" prop="count" width="100"/>
        <el-table-column :label="t('pages.randomNicknameCfg.samples')" min-width="260">
          <template #default="{ row }">
            <span v-if="row.samples?.length">{{ row.samples.join('、') }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column fixed="right" :label="t('common.actions')" width="100">
          <template #default="{ row }">
            <el-popconfirm
                :title="t('pages.randomNicknameCfg.clearConfirm', {lang: row.langLabel})"
                :confirm-button-text="t('pages.randomNicknameCfg.clear')"
                confirm-button-type="danger"
                @confirm="handleClear(row.lang)"
            >
              <template #reference>
                <el-button :disabled="row.count <= 0" link type="danger">{{ t('pages.randomNicknameCfg.clear') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-divider content-position="left">{{ t('pages.randomNicknameCfg.batchImport') }}</el-divider>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="import-form" label-width="120px">
        <el-form-item :label="t('pages.randomNicknameCfg.language')" prop="lang">
          <el-select v-model="formData.lang" :placeholder="t('pages.randomNicknameCfg.selectLanguage')" style="width: 240px">
            <el-option
                v-for="item in langOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item :label="t('pages.randomNicknameCfg.importMode')">
          <el-radio-group v-model="formData.replace">
            <el-radio :value="false">{{ t('pages.randomNicknameCfg.appendMode') }}</el-radio>
            <el-radio :value="true">{{ t('pages.randomNicknameCfg.replaceMode') }}</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item :label="t('pages.randomNicknameCfg.nicknameContent')" prop="content">
          <el-input
              v-model="formData.content"
              :rows="14"
              :placeholder="t('pages.randomNicknameCfg.nicknamePlaceholder')"
              type="textarea"
          />
          <span class="form-tip">{{ t('pages.randomNicknameCfg.parsedCount', {count: parsedCount}) }}</span>
        </el-form-item>

        <el-form-item>
          <el-button :loading="submitting" type="primary" @click="handleImport">{{ t('pages.randomNicknameCfg.import') }}</el-button>
          <el-button @click="fetchCfg">{{ t('pages.randomNicknameCfg.refreshStats') }}</el-button>
          <el-button @click="resetForm">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage} from 'element-plus'
import {randomNicknameApi} from '@/api/modules/randomNickname'
import type {GetRandomNicknameCfgRes, RandomNicknameLangItem} from '@/types/api'

const {t} = useI18n()
const loading = ref(false)
const submitting = ref(false)
const formRef = ref()

const langOptions = computed(() => [
  {value: 1, label: t('pages.randomNicknameCfg.langEn')},
  {value: 2, label: t('pages.randomNicknameCfg.langEs')},
  {value: 3, label: t('pages.randomNicknameCfg.langHi')},
  {value: 4, label: t('pages.randomNicknameCfg.langPt')},
  {value: 5, label: t('pages.randomNicknameCfg.langId')},
])

const cfg = reactive<GetRandomNicknameCfgRes>({
  useDB: false,
  langs: [],
})

const formData = reactive({
  lang: 1,
  content: '',
  replace: false,
})

const formRules = computed(() => ({
  lang: [{required: true, message: t('pages.randomNicknameCfg.langRequired'), trigger: 'change'}],
  content: [{required: true, message: t('pages.randomNicknameCfg.contentRequired'), trigger: 'blur'}],
}))

const parsedCount = computed(() => parseNicknames(formData.content).length)

function parseNicknames(text: string): string[] {
  const seen = new Set<string>()
  const out: string[] = []
  for (const line of text.split('\n')) {
    const name = line.trim()
    if (!name || seen.has(name)) {
      continue
    }
    seen.add(name)
    out.push(name)
  }
  return out
}

async function fetchCfg() {
  loading.value = true
  try {
    const res = await randomNicknameApi.getRandomNicknameCfg()
    cfg.useDB = res.useDB
    cfg.langs = res.langs || []
    if (!cfg.langs.length) {
      cfg.langs = langOptions.value.map(item => ({
        lang: item.value,
        langCode: ['en', 'es', 'hi', 'pt', 'id'][item.value - 1] || 'en',
        langLabel: item.label.split('（')[0].split('(')[0].trim(),
        count: 0,
        samples: [],
      })) as RandomNicknameLangItem[]
    }
  } finally {
    loading.value = false
  }
}

async function handleImport() {
  await formRef.value?.validate()
  if (parsedCount.value <= 0) {
    ElMessage.warning(t('pages.randomNicknameCfg.noValidNicknames'))
    return
  }
  submitting.value = true
  try {
    const res = await randomNicknameApi.importRandomNicknames({
      lang: formData.lang,
      content: formData.content,
      replace: formData.replace,
    })
    ElMessage.success(t('pages.randomNicknameCfg.importSuccess', {imported: res.imported, total: res.total}))
    formData.content = ''
    await fetchCfg()
  } finally {
    submitting.value = false
  }
}

async function handleClear(lang: number) {
  loading.value = true
  try {
    const res = await randomNicknameApi.clearRandomNicknames({lang})
    if (res.success) {
      ElMessage.success(t('pages.randomNicknameCfg.cleared'))
      await fetchCfg()
    }
  } finally {
    loading.value = false
  }
}

function resetForm() {
  formData.lang = 1
  formData.content = ''
  formData.replace = false
  formRef.value?.clearValidate()
}

onMounted(fetchCfg)
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.tip-alert {
  margin-bottom: 16px;
}

.tip-alert p {
  margin: 4px 0;
}

.stat-table {
  margin-bottom: 8px;
}

.import-form {
  max-width: 960px;
}

.form-tip {
  display: block;
  margin-top: 6px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.text-muted {
  color: var(--el-text-color-secondary);
}
</style>
