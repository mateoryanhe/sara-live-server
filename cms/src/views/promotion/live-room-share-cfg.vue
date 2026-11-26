<template>
  <div class="page-container">
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.PromotionLiveRoomShareCfgManagement') }}</span>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="formRules" class="cfg-form" label-width="140px">
        <el-form-item :label="t('pages.liveRoomShare.afLink')" prop="afLink">
          <el-input
              v-model="formData.afLink"
              clearable
              :placeholder="t('pages.liveRoomShare.afLinkPlaceholder')"
          />
          <div class="form-tip">{{ t('pages.liveRoomShare.afLinkTip') }}</div>
        </el-form-item>

        <el-divider>{{ t('pages.liveRoomShare.shareTextTab') }}</el-divider>
        <p class="form-tip builtin-tip">{{ t('pages.liveRoomShare.shareTextBuiltinTip') }}</p>

        <el-tabs v-model="activeLocale" type="border-card">
          <el-tab-pane
              v-for="loc in LIVE_ROOM_SHARE_LOCALES"
              :key="loc"
              :label="LOCALE_LABELS[loc]"
              :name="loc"
          >
            <el-input
                v-model="formData[localeField(loc)]"
                type="textarea"
                :rows="8"
                maxlength="4000"
                show-word-limit
                :placeholder="t('pages.liveRoomShare.shareTextPlaceholder')"
            />
          </el-tab-pane>
        </el-tabs>

        <el-form-item class="preview-block" :label="t('pages.liveRoomShare.previewTitle')">
          <el-input
              :model-value="previewText || t('pages.liveRoomShare.previewEmpty')"
              type="textarea"
              :rows="4"
              readonly
          />
        </el-form-item>

        <el-form-item v-if="metaInfo.updatedAt" :label="t('pages.shareAppDownload.lastUpdated')">
          <span>{{ metaInfo.updatedAt }}</span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSave">{{ t('common.saveConfig') }}</el-button>
          <el-button @click="fetchCfg">{{ t('common.refresh') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref} from 'vue'
import type {FormInstance, FormRules} from 'element-plus'
import {ElMessage} from 'element-plus'
import {liveRoomShareApi} from '@/api/modules/live-room-share'
import type {LiveRoomShareCfgItem} from '@/types/api'
import {LOCALE_LABELS, type CmsLocale} from '@/i18n/locales/types'

/** 房间分享文案：不维护简繁体，默认英文 Tab */
const LIVE_ROOM_SHARE_LOCALES = ['en', 'es', 'pt', 'hi', 'id'] as const satisfies readonly CmsLocale[]
type LiveRoomShareLocale = (typeof LIVE_ROOM_SHARE_LOCALES)[number]

const {t} = useI18n()
const loading = ref(false)
const formRef = ref<FormInstance>()
const activeLocale = ref<LiveRoomShareLocale>('en')

type TextField = 'shareTextEn' | 'shareTextEs' | 'shareTextPt' | 'shareTextHi' | 'shareTextId'

const localeFieldMap: Record<LiveRoomShareLocale, TextField> = {
  en: 'shareTextEn',
  es: 'shareTextEs',
  pt: 'shareTextPt',
  hi: 'shareTextHi',
  id: 'shareTextId',
}

const localeField = (loc: LiveRoomShareLocale) => localeFieldMap[loc]

const formData = reactive({
  id: '0',
  afLink: '',
  shareTextZhCN: '',
  shareTextZhTW: '',
  shareTextEn: '',
  shareTextEs: '',
  shareTextPt: '',
  shareTextHi: '',
  shareTextId: '',
})

const metaInfo = reactive({createdAt: '', updatedAt: ''})

const builtinShareText = reactive<Record<LiveRoomShareLocale, string>>({
  en: '',
  es: '',
  pt: '',
  hi: '',
  id: '',
})

function effectiveShareLine(loc: LiveRoomShareLocale): string {
  const custom = formData[localeField(loc)].trim()
  if (custom) return custom
  return (builtinShareText[loc] || builtinShareText.en || '').trim()
}

const formRules: FormRules = {
  afLink: [{required: true, message: () => t('pages.liveRoomShare.needAfLink'), trigger: 'blur'}],
}

const previewText = computed(() => {
  const line1 = effectiveShareLine('en')
  const base = formData.afLink.trim()
  if (!line1 && !base) return ''
  const line2 = appendPreviewAf(base)
  if (!line1) return line2
  if (!line2) return line1
  return `${line1}\n${line2}`
})

function appendPreviewAf(base: string): string {
  if (!base) return ''
  const sep = base.includes('?') ? '&' : '?'
  return `${base}${sep}inviteCode=YOUR_INVITE_CODE&roomId=123456`
}

const applyBuiltinShareText = (cfg: LiveRoomShareCfgItem | null | undefined) => {
  builtinShareText.en = cfg?.builtinShareTextEn?.trim() || ''
  builtinShareText.es = cfg?.builtinShareTextEs?.trim() || builtinShareText.en
  builtinShareText.pt = cfg?.builtinShareTextPt?.trim() || builtinShareText.en
  builtinShareText.hi = cfg?.builtinShareTextHi?.trim() || builtinShareText.en
  builtinShareText.id = cfg?.builtinShareTextId?.trim() || builtinShareText.en
}

const applyCfg = (cfg: LiveRoomShareCfgItem | null | undefined) => {
  applyBuiltinShareText(cfg)
  if (!cfg) {
    formData.id = '0'
    formData.afLink = ''
    Object.values(localeFieldMap).forEach((key) => {
      formData[key] = ''
    })
    formData.shareTextZhCN = ''
    formData.shareTextZhTW = ''
    metaInfo.createdAt = ''
    metaInfo.updatedAt = ''
    return
  }
  formData.id = cfg.id || '0'
  formData.afLink = cfg.afLink || ''
  formData.shareTextZhCN = ''
  formData.shareTextZhTW = ''
  formData.shareTextEn = cfg.shareTextEn || ''
  formData.shareTextEs = cfg.shareTextEs || ''
  formData.shareTextPt = cfg.shareTextPt || ''
  formData.shareTextHi = cfg.shareTextHi || ''
  formData.shareTextId = cfg.shareTextId || ''
  metaInfo.createdAt = cfg.createdAt || ''
  metaInfo.updatedAt = cfg.updatedAt || ''
}

const fetchCfg = async () => {
  loading.value = true
  try {
    const res = await liveRoomShareApi.getLiveRoomShareCfg()
    applyCfg(res.cfg)
  } catch {
    ElMessage.error(t('common.loadFailed'))
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  const ok = await formRef.value?.validate().catch(() => false)
  if (!ok) return
  loading.value = true
  try {
    const res = await liveRoomShareApi.saveLiveRoomShareCfg({
      id: formData.id !== '0' ? Number(formData.id) : undefined,
      afLink: formData.afLink.trim(),
      shareTextZhCN: '',
      shareTextZhTW: '',
      shareTextEn: formData.shareTextEn.trim(),
      shareTextEs: formData.shareTextEs.trim(),
      shareTextPt: formData.shareTextPt.trim(),
      shareTextHi: formData.shareTextHi.trim(),
      shareTextId: formData.shareTextId.trim(),
    })
    if (res.success) {
      ElMessage.success(t('pages.liveRoomShare.saveSuccess'))
      if (res.id) formData.id = res.id
      await fetchCfg()
    }
  } catch {
    ElMessage.error(t('common.saveFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  void fetchCfg()
})
</script>

<style scoped>
.form-tip {
  margin-top: 6px;
  color: #64748b;
  font-size: 13px;
  line-height: 1.5;
}

.preview-block {
  margin-top: 20px;
}

.builtin-tip {
  margin: 0 0 12px;
}
</style>