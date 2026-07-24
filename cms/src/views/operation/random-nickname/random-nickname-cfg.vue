<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>随机昵称库</span>
          <el-tag :type="cfg.useDB ? 'success' : 'info'" size="small">
            {{ cfg.useDB ? '使用数据库昵称库' : '使用内置英文昵称库' }}
          </el-tag>
        </div>
      </template>

      <el-alert
          :closable="false"
          class="tip-alert"
          show-icon
          title="说明"
          type="info"
      >
        <p>支持英文、西班牙语、印地语、葡萄牙语。一行一个昵称，导入后写入数据库并加载到内存；新用户注册时按 Accept-Language 选择对应语言池。</p>
        <p>数据库为空时，服务器启动后使用 Go 内置 1000 条英文昵称；有数据后仅使用数据库内容，不再读内置库。</p>
      </el-alert>

      <el-table :data="cfg.langs" border class="stat-table" empty-text="暂无数据">
        <el-table-column label="语言" min-width="140" prop="langLabel"/>
        <el-table-column label="代码" prop="langCode" width="80"/>
        <el-table-column label="数量" prop="count" width="100"/>
        <el-table-column label="示例" min-width="260">
          <template #default="{ row }">
            <span v-if="row.samples?.length">{{ row.samples.join('、') }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
        <el-table-column fixed="right" label="操作" width="100">
          <template #default="{ row }">
            <el-popconfirm
                :title="`确定清空 ${row.langLabel} 昵称库？`"
                confirm-button-text="清空"
                confirm-button-type="danger"
                @confirm="handleClear(row.lang)"
            >
              <template #reference>
                <el-button :disabled="row.count <= 0" link type="danger">清空</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-divider content-position="left">批量导入</el-divider>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="import-form" label-width="120px">
        <el-form-item label="语言" prop="lang">
          <el-select v-model="formData.lang" placeholder="请选择语言" style="width: 240px">
            <el-option
                v-for="item in langOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="导入模式">
          <el-radio-group v-model="formData.replace">
            <el-radio :value="false">追加（去重）</el-radio>
            <el-radio :value="true">覆盖（先清空该语言再导入）</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="昵称内容" prop="content">
          <el-input
              v-model="formData.content"
              :rows="14"
              placeholder="一行一个昵称，空行会自动忽略"
              type="textarea"
          />
          <span class="form-tip">已解析 {{ parsedCount }} 条有效昵称（去重后）</span>
        </el-form-item>

        <el-form-item>
          <el-button :loading="submitting" type="primary" @click="handleImport">导入</el-button>
          <el-button @click="fetchCfg">刷新统计</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {randomNicknameApi} from '@/api/modules/randomNickname'
import type {GetRandomNicknameCfgRes, RandomNicknameLangItem} from '@/types/api'

const loading = ref(false)
const submitting = ref(false)
const formRef = ref()

const langOptions = [
  {value: 1, label: 'English（英文，默认）'},
  {value: 2, label: 'Español（西班牙语）'},
  {value: 3, label: 'हिन्दी（印地语）'},
  {value: 4, label: 'Português（葡萄牙语）'},
]

const cfg = reactive<GetRandomNicknameCfgRes>({
  useDB: false,
  langs: [],
})

const formData = reactive({
  lang: 1,
  content: '',
  replace: false,
})

const formRules = reactive({
  lang: [{required: true, message: '请选择语言', trigger: 'change'}],
  content: [{required: true, message: '请输入昵称内容', trigger: 'blur'}],
})

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
      cfg.langs = langOptions.map(item => ({
        lang: item.value,
        langCode: ['en', 'es', 'hi', 'pt'][item.value - 1] || 'en',
        langLabel: item.label.split('（')[0],
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
    ElMessage.warning('没有可导入的有效昵称')
    return
  }
  submitting.value = true
  try {
    const res = await randomNicknameApi.importRandomNicknames({
      lang: formData.lang,
      content: formData.content,
      replace: formData.replace,
    })
    ElMessage.success(`导入成功：本次 ${res.imported} 条，该语言共 ${res.total} 条`)
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
      ElMessage.success('已清空')
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
