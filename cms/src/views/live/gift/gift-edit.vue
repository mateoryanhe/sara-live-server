<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-button @click="goBack">{{ t('pages.giftList.backToList') }}</el-button>
          <span>{{ pageTitle }}</span>
        </div>
      </template>

      <el-skeleton v-if="pageLoading" :rows="12" animated/>

      <template v-else>
        <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="110px">
          <el-tabs v-model="activeTab" class="gift-edit-tabs" type="border-card">
            <el-tab-pane :label="t('pages.giftList.tabBasicInfo')" name="basic">
              <div class="tab-pane-body">
                <el-form-item :label="t('pages.giftList.giftName')" prop="name">
                  <el-input v-model="currentRow.name" :placeholder="t('pages.giftList.giftNamePlaceholder')"/>
                </el-form-item>
                <el-form-item :label="t('common.nameEn')" prop="nameEn">
                  <el-input v-model="currentRow.nameEn" :placeholder="t('pages.giftList.nameEnPlaceholder')"/>
                </el-form-item>
                <el-form-item :label="t('common.nameEs')" prop="nameEs">
                  <el-input v-model="currentRow.nameEs" :placeholder="t('pages.giftList.nameEsPlaceholder')"/>
                </el-form-item>
                <el-form-item :label="t('common.namePt')" prop="namePt">
                  <el-input v-model="currentRow.namePt" :placeholder="t('pages.giftList.namePtPlaceholder')"/>
                </el-form-item>
                <el-form-item :label="t('common.nameHi')" prop="nameHi">
                  <el-input v-model="currentRow.nameHi" :placeholder="t('pages.giftList.nameHiPlaceholder')"/>
                </el-form-item>
                <el-form-item :label="t('common.nameId')" prop="nameId">
                  <el-input v-model="currentRow.nameId" :placeholder="t('pages.giftList.nameIdPlaceholder')"/>
                </el-form-item>
                <el-form-item :label="t('common.icon')" prop="icon">
                  <div class="asset-upload-wrap">
                    <el-upload
                        :before-upload="beforeIconUpload"
                        :disabled="iconUploading"
                        :http-request="(opt: UploadRequestOptions) => doUpload(opt, 'icon')"
                        :show-file-list="false"
                        accept="image/*"
                        action="#"
                        class="icon-uploader"
                    >
                      <img v-if="iconPreviewUrl" :src="iconPreviewUrl" alt="icon" class="icon-preview"/>
                      <div v-else class="asset-uploader-placeholder icon-placeholder">
                        <el-icon class="asset-uploader-icon">
                          <Plus/>
                        </el-icon>
                        <span>{{ t('pages.giftList.clickUploadIcon') }}</span>
                      </div>
                    </el-upload>
                    <el-button
                        v-if="iconPreviewUrl || currentRow.icon"
                        link
                        type="danger"
                        @click="clearAsset('icon')"
                    >
                      {{ t('pages.giftList.removeIcon') }}
                    </el-button>
                  </div>
                </el-form-item>
                <el-form-item :label="t('common.animationResource')" prop="animation">
                  <div class="asset-upload-wrap">
                    <el-upload
                        :before-upload="beforeAnimationUpload"
                        :disabled="animationUploading"
                        :http-request="(opt: UploadRequestOptions) => doUpload(opt, 'animation')"
                        :show-file-list="false"
                        accept=".svga,.pag,.json,.lottie,.mp4,.webm,.zip,.gif,.apng,.png,.webp,.jpg,.jpeg,.bmp"
                        action="#"
                        class="animation-uploader"
                    >
                      <img
                          v-if="animationPreviewType === 'image'"
                          :src="animationPreviewUrl"
                          alt="animation"
                          class="animation-preview"
                      />
                      <video
                          v-else-if="animationPreviewType === 'video'"
                          :src="animationPreviewUrl"
                          class="animation-preview"
                          controls
                          preload="metadata"
                      />
                      <div v-else-if="animationPreviewType === 'file' && animationFileLabel" class="asset-file-label">
                        <el-icon class="asset-uploader-icon">
                          <Document/>
                        </el-icon>
                        <span class="file-name">{{ animationFileLabel }}</span>
                      </div>
                      <div v-else class="asset-uploader-placeholder animation-placeholder">
                        <el-icon class="asset-uploader-icon">
                          <Plus/>
                        </el-icon>
                        <span>{{ t('pages.giftList.clickUploadAnimation') }}</span>
                      </div>
                    </el-upload>
                    <el-button
                        v-if="animationPreviewType !== 'none' || currentRow.animation"
                        link
                        type="danger"
                        @click="clearAsset('animation')"
                    >
                      {{ t('pages.giftList.removeAnimation') }}
                    </el-button>
                  </div>
                </el-form-item>
                <el-form-item :label="t('pages.giftList.priceDiamond')" prop="price">
                  <el-input-number v-model="currentRow.price" :min="0" :precision="NUMBER_INPUT_DECIMALS" controls-position="right"/>
                </el-form-item>
                <el-form-item :label="t('common.category')" prop="categoryType">
                  <el-select v-model="currentRow.categoryType" style="width: 100%">
                    <el-option :label="t('pages.giftList.categoryNormal')" :value="0"/>
                    <el-option :label="t('pages.giftList.categoryLucky')" :value="1"/>
                  </el-select>
                </el-form-item>
                <el-form-item :label="t('common.sort')" prop="sort">
                  <el-input-number v-model="currentRow.sort" controls-position="right"/>
                </el-form-item>
                <el-form-item :label="t('common.description')" prop="description">
                  <el-input v-model="currentRow.description" :placeholder="t('pages.giftList.descriptionPlaceholder')" type="textarea"/>
                </el-form-item>
              </div>
            </el-tab-pane>

            <el-tab-pane :label="t('pages.giftList.tabRtp')" name="rtp" :disabled="!isLuckyGift">
              <div class="tab-pane-body">
                <el-alert
                    v-if="!isLuckyGift"
                    :closable="false"
                    :title="t('pages.giftList.luckyTabsDisabledTip')"
                    show-icon
                    type="info"
                />
                <template v-else>
                  <el-form-item :label="t('pages.giftList.luckyGiftRtp')" prop="luckyGiftRtp">
                    <el-input-number
                        v-model="currentRow.luckyGiftRtp"
                        :max="100"
                        :min="0"
                        :precision="NUMBER_INPUT_DECIMALS"
                        :step="0.0001"
                        controls-position="right"
                        style="width: 220px"
                    />
                    <span class="form-tip-inline">{{ t('pages.giftList.luckyGiftRtpTip') }}</span>
                  </el-form-item>
                </template>
              </div>
            </el-tab-pane>

            <el-tab-pane :label="t('pages.giftList.tabMultiplier')" name="multiplier" :disabled="!isLuckyGift">
              <div class="tab-pane-body">
                <el-alert
                    v-if="!isLuckyGift"
                    :closable="false"
                    :title="t('pages.giftList.luckyTabsDisabledTip')"
                    show-icon
                    type="info"
                />
                <el-form-item v-else :label="t('pages.giftList.luckyMultiplierTiers')">
                  <LuckyGiftTierEditor v-model="currentRow.luckyGiftMultiplierTiers"/>
                </el-form-item>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-form>

        <div class="footer-actions">
          <el-button @click="goBack">{{ t('common.cancel') }}</el-button>
          <el-button :loading="saving" type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
        </div>
      </template>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {computed, onBeforeUnmount, onMounted, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {useRoute, useRouter} from 'vue-router'
import {Document, Plus} from '@element-plus/icons-vue'
import type {UploadRequestOptions} from 'element-plus'
import LuckyGiftTierEditor from '@/components/gift/LuckyGiftTierEditor.vue'
import {useGiftFormEditor} from '@/composables/useGiftFormEditor'
import type {Gift} from '@/types/api.ts'
import {NUMBER_INPUT_DECIMALS} from '@/utils/number-format'

const {t} = useI18n()
const route = useRoute()
const router = useRouter()

const giftIdParam = computed(() => String(route.params.id || ''))
const isCreate = computed(() => giftIdParam.value === 'new' || giftIdParam.value === '')

const pageTitle = computed(() =>
  isCreate.value ? t('pages.giftList.addGift') : t('pages.giftList.editGift'),
)

const {
  pageLoading,
  saving,
  activeTab,
  formRef,
  currentRow,
  isLuckyGift,
  iconUploading,
  animationUploading,
  iconPreviewUrl,
  animationPreviewUrl,
  animationPreviewType,
  animationFileLabel,
  formRules,
  beforeIconUpload,
  beforeAnimationUpload,
  doUpload,
  clearAsset,
  initCreateForm,
  applyGiftRowToForm,
  loadGiftById,
  saveGift,
  dispose,
} = useGiftFormEditor()

watch(isLuckyGift, (lucky) => {
  if (!lucky && (activeTab.value === 'rtp' || activeTab.value === 'multiplier')) {
    activeTab.value = 'basic'
  }
})

const goBack = () => {
  router.push({name: 'GiftManagement'})
}

const handleSave = async () => {
  const ok = await saveGift()
  if (ok) {
    goBack()
  }
}

const bootstrap = async () => {
  if (isCreate.value) {
    initCreateForm()
    return
  }
  const stateRow = history.state?.giftRow as Gift | undefined
  if (stateRow && String(stateRow.id) === giftIdParam.value) {
    pageLoading.value = true
    try {
      await applyGiftRowToForm(stateRow)
    } finally {
      pageLoading.value = false
    }
    return
  }
  const loaded = await loadGiftById(giftIdParam.value)
  if (!loaded) {
    goBack()
  }
}

onMounted(() => {
  void bootstrap()
})

onBeforeUnmount(() => {
  dispose()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 16px;
  font-weight: bold;
}

.gift-edit-tabs {
  margin-top: 8px;
}

.tab-pane-body {
  max-width: 920px;
  padding: 16px 8px 8px;
}

.footer-actions {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.form-tip-inline {
  margin-left: 12px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.asset-upload-wrap {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.icon-uploader :deep(.el-upload),
.animation-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  transition: border-color 0.2s;
}

.icon-uploader :deep(.el-upload:hover),
.animation-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.asset-uploader-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-secondary);
  font-size: 13px;
  gap: 8px;
}

.icon-placeholder {
  width: 96px;
  height: 96px;
}

.animation-placeholder {
  width: 240px;
  height: 120px;
}

.asset-uploader-icon {
  font-size: 28px;
}

.icon-preview {
  width: 96px;
  height: 96px;
  display: block;
  object-fit: cover;
}

.animation-preview {
  width: 240px;
  height: 120px;
  display: block;
  object-fit: cover;
}

.asset-file-label {
  width: 240px;
  height: 120px;
  padding: 12px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  text-align: center;
}

.asset-file-label .file-name {
  font-size: 13px;
  color: var(--el-text-color-regular);
  word-break: break-all;
  line-height: 1.4;
}
</style>
