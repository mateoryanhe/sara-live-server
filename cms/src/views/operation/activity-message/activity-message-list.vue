<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.ActivityMessageManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.activityMessageList.addMessage') }}</el-button>
          <el-button
              v-if="hasButtonPermission('ActivityMessageManagement', 'sync')"
              :disabled="selectedRows.length === 0"
              :loading="syncing"
              type="warning"
              @click="handleSyncData"
          >
            {{ t('common.syncData') }}
          </el-button>
        </div>

        <div v-if="selectedRows.length" class="selection-tip">
          {{ t('common.selectedCount', {count: selectedRows.length}) }}
        </div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('common.title')">
            <el-input v-model="searchForm.title" clearable :placeholder="t('pages.activityMessageList.titleFuzzy')"/>
          </el-form-item>
          <el-form-item :label="t('pages.activityMessageList.publishStatus')">
            <el-select v-model="searchForm.statusFilter" :placeholder="t('common.all')" style="width: 140px">
              <el-option :value="0" :label="t('common.all')"/>
              <el-option :value="2" :label="t('pages.activityMessageList.onlyPublished')"/>
              <el-option :value="1" :label="t('pages.activityMessageList.onlyUnpublished')"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">{{ t('common.search') }}</el-button>
            <el-button @click="resetSearch">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>

        <el-table
            v-loading="loading"
            :data="tableData"
            row-key="id"
            style="width: 100%"
            @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="48"/>
          <el-table-column label="ID" prop="id" width="90"/>
          <el-table-column :label="t('pages.activityMessageList.iconEn')" width="100">
            <template #default="{ row }">
              <el-image
                  v-if="row.iconEn"
                  :preview-src-list="[row.iconEn]"
                  :src="row.iconEn"
                  fit="cover"
                  preview-teleported
                  style="width: 48px; height: 48px"
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.activityMessageList.bgEn')" width="110">
            <template #default="{ row }">
              <el-image
                  v-if="row.bgEn"
                  :preview-src-list="[row.bgEn]"
                  :src="row.bgEn"
                  fit="cover"
                  preview-teleported
                  style="width: 72px; height: 40px"
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.activityMessageList.titleEn')" prop="titleEn" min-width="160" show-overflow-tooltip/>
          <el-table-column :label="t('pages.activityMessageList.urlEn')" prop="urlEn" min-width="200" show-overflow-tooltip/>
          <el-table-column :label="t('common.publishedAt')" prop="publishedAt" width="160">
            <template #default="{ row }">
              {{ row.publishedAt || '-' }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.activityMessageList.publishStatus')" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? t('pages.activityMessageList.published') : t('pages.activityMessageList.unpublished') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="280">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
              <el-button
                  v-if="row.status !== 1"
                  size="small"
                  type="success"
                  @click="handlePublish(row)"
              >
                {{ t('pages.activityMessageList.publish') }}
              </el-button>
              <el-button
                  v-else
                  size="small"
                  type="warning"
                  @click="handleUnpublish(row)"
              >
                {{ t('pages.activityMessageList.unpublish') }}
              </el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">{{ t('common.delete') }}</el-button>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="total"
              layout="total, sizes, prev, pager, next, jumper"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
          />
        </div>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" destroy-on-close width="760px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="110px">
        <el-tabs v-model="activeLangTab">
          <el-tab-pane label="English" name="en">
            <el-form-item :label="t('common.icon')" prop="iconEn">
              <div class="icon-upload-wrap">
                <el-upload
                    :before-upload="beforeImageUpload"
                    :disabled="assetUploading.icon.en"
                    :http-request="(opt) => doAssetUpload('icon', 'en', opt)"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="icon-uploader"
                >
                  <img v-if="assetPreviewUrls.icon.en" :src="assetPreviewUrls.icon.en" alt="icon" class="icon-preview"/>
                  <div v-else class="icon-uploader-placeholder">
                    <el-icon class="icon-uploader-icon"><Plus/></el-icon>
                    <span>{{ t('pages.activityMessageList.clickUploadIcon') }}</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.icon.en || currentRow.iconEn" link type="danger" @click="clearAsset('icon', 'en')">
                  {{ t('pages.activityMessageList.removeIcon') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.bg')">
              <div class="icon-upload-wrap">
                <el-upload
                    :before-upload="beforeImageUpload"
                    :disabled="assetUploading.bg.en"
                    :http-request="(opt) => doAssetUpload('bg', 'en', opt)"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="bg-uploader"
                >
                  <img v-if="assetPreviewUrls.bg.en" :src="assetPreviewUrls.bg.en" alt="bg" class="bg-preview"/>
                  <div v-else class="bg-uploader-placeholder">
                    <el-icon class="icon-uploader-icon"><Plus/></el-icon>
                    <span>{{ t('pages.activityMessageList.clickUploadBg') }}</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.bg.en || currentRow.bgEn" link type="danger" @click="clearAsset('bg', 'en')">
                  {{ t('pages.activityMessageList.removeBg') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('common.title')" prop="titleEn">
              <el-input v-model="currentRow.titleEn" maxlength="128" :placeholder="t('pages.activityMessageList.titleEnPlaceholder')" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.content')" prop="contentEn">
              <el-input v-model="currentRow.contentEn" :rows="5" :placeholder="t('pages.activityMessageList.contentEnPlaceholder')" type="textarea"/>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.link')" prop="urlEn">
              <el-input v-model="currentRow.urlEn" maxlength="512" :placeholder="t('pages.activityMessageList.urlEnPlaceholder')" show-word-limit/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="Español" name="es">
            <el-form-item :label="t('common.icon')">
              <div class="icon-upload-wrap">
                <el-upload
                    :before-upload="beforeImageUpload"
                    :disabled="assetUploading.icon.es"
                    :http-request="(opt) => doAssetUpload('icon', 'es', opt)"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="icon-uploader"
                >
                  <img v-if="assetPreviewUrls.icon.es" :src="assetPreviewUrls.icon.es" alt="icon" class="icon-preview"/>
                  <div v-else class="icon-uploader-placeholder">
                    <el-icon class="icon-uploader-icon"><Plus/></el-icon>
                    <span>{{ t('pages.activityMessageList.clickUploadIcon') }}</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.icon.es || currentRow.iconEs" link type="danger" @click="clearAsset('icon', 'es')">
                  {{ t('pages.activityMessageList.removeIcon') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.bg')">
              <div class="icon-upload-wrap">
                <el-upload
                    :before-upload="beforeImageUpload"
                    :disabled="assetUploading.bg.es"
                    :http-request="(opt) => doAssetUpload('bg', 'es', opt)"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="bg-uploader"
                >
                  <img v-if="assetPreviewUrls.bg.es" :src="assetPreviewUrls.bg.es" alt="bg" class="bg-preview"/>
                  <div v-else class="bg-uploader-placeholder">
                    <el-icon class="icon-uploader-icon"><Plus/></el-icon>
                    <span>{{ t('pages.activityMessageList.clickUploadBg') }}</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.bg.es || currentRow.bgEs" link type="danger" @click="clearAsset('bg', 'es')">
                  {{ t('pages.activityMessageList.removeBg') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('common.title')">
              <el-input v-model="currentRow.titleEs" maxlength="128" :placeholder="t('pages.activityMessageList.titleEsPlaceholder')" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.content')">
              <el-input v-model="currentRow.contentEs" :rows="5" :placeholder="t('pages.activityMessageList.contentEsPlaceholder')" type="textarea"/>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.link')">
              <el-input v-model="currentRow.urlEs" maxlength="512" :placeholder="t('pages.activityMessageList.urlEsPlaceholder')" show-word-limit/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="Português" name="pt">
            <el-form-item :label="t('common.icon')">
              <div class="icon-upload-wrap">
                <el-upload
                    :before-upload="beforeImageUpload"
                    :disabled="assetUploading.icon.pt"
                    :http-request="(opt) => doAssetUpload('icon', 'pt', opt)"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="icon-uploader"
                >
                  <img v-if="assetPreviewUrls.icon.pt" :src="assetPreviewUrls.icon.pt" alt="icon" class="icon-preview"/>
                  <div v-else class="icon-uploader-placeholder">
                    <el-icon class="icon-uploader-icon"><Plus/></el-icon>
                    <span>{{ t('pages.activityMessageList.clickUploadIcon') }}</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.icon.pt || currentRow.iconPt" link type="danger" @click="clearAsset('icon', 'pt')">
                  {{ t('pages.activityMessageList.removeIcon') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.bg')">
              <div class="icon-upload-wrap">
                <el-upload
                    :before-upload="beforeImageUpload"
                    :disabled="assetUploading.bg.pt"
                    :http-request="(opt) => doAssetUpload('bg', 'pt', opt)"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="bg-uploader"
                >
                  <img v-if="assetPreviewUrls.bg.pt" :src="assetPreviewUrls.bg.pt" alt="bg" class="bg-preview"/>
                  <div v-else class="bg-uploader-placeholder">
                    <el-icon class="icon-uploader-icon"><Plus/></el-icon>
                    <span>{{ t('pages.activityMessageList.clickUploadBg') }}</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.bg.pt || currentRow.bgPt" link type="danger" @click="clearAsset('bg', 'pt')">
                  {{ t('pages.activityMessageList.removeBg') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('common.title')">
              <el-input v-model="currentRow.titlePt" maxlength="128" :placeholder="t('pages.activityMessageList.titlePtPlaceholder')" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.content')">
              <el-input v-model="currentRow.contentPt" :rows="5" :placeholder="t('pages.activityMessageList.contentPtPlaceholder')" type="textarea"/>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.link')">
              <el-input v-model="currentRow.urlPt" maxlength="512" :placeholder="t('pages.activityMessageList.urlPtPlaceholder')" show-word-limit/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="हिन्दी" name="hi">
            <el-form-item :label="t('common.icon')">
              <div class="icon-upload-wrap">
                <el-upload
                    :before-upload="beforeImageUpload"
                    :disabled="assetUploading.icon.hi"
                    :http-request="(opt) => doAssetUpload('icon', 'hi', opt)"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="icon-uploader"
                >
                  <img v-if="assetPreviewUrls.icon.hi" :src="assetPreviewUrls.icon.hi" alt="icon" class="icon-preview"/>
                  <div v-else class="icon-uploader-placeholder">
                    <el-icon class="icon-uploader-icon"><Plus/></el-icon>
                    <span>{{ t('pages.activityMessageList.clickUploadIcon') }}</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.icon.hi || currentRow.iconHi" link type="danger" @click="clearAsset('icon', 'hi')">
                  {{ t('pages.activityMessageList.removeIcon') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.bg')">
              <div class="icon-upload-wrap">
                <el-upload
                    :before-upload="beforeImageUpload"
                    :disabled="assetUploading.bg.hi"
                    :http-request="(opt) => doAssetUpload('bg', 'hi', opt)"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="bg-uploader"
                >
                  <img v-if="assetPreviewUrls.bg.hi" :src="assetPreviewUrls.bg.hi" alt="bg" class="bg-preview"/>
                  <div v-else class="bg-uploader-placeholder">
                    <el-icon class="icon-uploader-icon"><Plus/></el-icon>
                    <span>{{ t('pages.activityMessageList.clickUploadBg') }}</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.bg.hi || currentRow.bgHi" link type="danger" @click="clearAsset('bg', 'hi')">
                  {{ t('pages.activityMessageList.removeBg') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('common.title')">
              <el-input v-model="currentRow.titleHi" maxlength="128" :placeholder="t('pages.activityMessageList.titleHiPlaceholder')" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.content')">
              <el-input v-model="currentRow.contentHi" :rows="5" :placeholder="t('pages.activityMessageList.contentHiPlaceholder')" type="textarea"/>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.link')">
              <el-input v-model="currentRow.urlHi" maxlength="512" :placeholder="t('pages.activityMessageList.urlHiPlaceholder')" show-word-limit/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="Bahasa Indonesia" name="id">
            <el-form-item :label="t('common.icon')">
              <div class="icon-upload-wrap">
                <el-upload
                    :before-upload="beforeImageUpload"
                    :disabled="assetUploading.icon.id"
                    :http-request="(opt) => doAssetUpload('icon', 'id', opt)"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="icon-uploader"
                >
                  <img v-if="assetPreviewUrls.icon.id" :src="assetPreviewUrls.icon.id" alt="icon" class="icon-preview"/>
                  <div v-else class="icon-uploader-placeholder">
                    <el-icon class="icon-uploader-icon"><Plus/></el-icon>
                    <span>{{ t('pages.activityMessageList.clickUploadIcon') }}</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.icon.id || currentRow.iconId" link type="danger" @click="clearAsset('icon', 'id')">
                  {{ t('pages.activityMessageList.removeIcon') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.bg')">
              <div class="icon-upload-wrap">
                <el-upload
                    :before-upload="beforeImageUpload"
                    :disabled="assetUploading.bg.id"
                    :http-request="(opt) => doAssetUpload('bg', 'id', opt)"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="bg-uploader"
                >
                  <img v-if="assetPreviewUrls.bg.id" :src="assetPreviewUrls.bg.id" alt="bg" class="bg-preview"/>
                  <div v-else class="bg-uploader-placeholder">
                    <el-icon class="icon-uploader-icon"><Plus/></el-icon>
                    <span>{{ t('pages.activityMessageList.clickUploadBg') }}</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.bg.id || currentRow.bgId" link type="danger" @click="clearAsset('bg', 'id')">
                  {{ t('pages.activityMessageList.removeBg') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('common.title')">
              <el-input v-model="currentRow.titleId" maxlength="128" :placeholder="t('pages.activityMessageList.titleIdPlaceholder')" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.content')">
              <el-input v-model="currentRow.contentId" :rows="5" :placeholder="t('pages.activityMessageList.contentIdPlaceholder')" type="textarea"/>
            </el-form-item>
            <el-form-item :label="t('pages.activityMessageList.link')">
              <el-input v-model="currentRow.urlId" maxlength="512" :placeholder="t('pages.activityMessageList.urlIdPlaceholder')" show-word-limit/>
            </el-form-item>
          </el-tab-pane>
        </el-tabs>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSave">{{ t('common.save') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {useI18n} from 'vue-i18n'
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {activityMessageApi, type ActivityMessageForm} from '@/api/modules/activityMessage'
import {dataSyncApi, uploadApi} from '@/api'
import type {ActivityMessage} from '@/types/api'
import {hasButtonPermission} from '@/utils/permission'

type LangKey = 'en' | 'es' | 'pt' | 'hi' | 'id'
type AssetKind = 'icon' | 'bg'

const {t} = useI18n()
const assetFieldMap: Record<AssetKind, Record<LangKey, keyof ActivityMessageForm>> = {
  icon: {en: 'iconEn', es: 'iconEs', pt: 'iconPt', hi: 'iconHi', id: 'iconId'},
  bg: {en: 'bgEn', es: 'bgEs', pt: 'bgPt', hi: 'bgHi', id: 'bgId'},
}

const emptyLangRecord = () => ({en: '', es: '', pt: '', hi: '', id: ''})
const emptyLangBoolRecord = () => ({en: false, es: false, pt: false, hi: false, id: false})

interface SearchForm {
  title: string
  statusFilter: number
}

const loading = ref(false)
const syncing = ref(false)
const selectedRows = ref<ActivityMessage[]>([])
const tableData = ref<ActivityMessage[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const activeLangTab = ref('en')

const searchForm = reactive<SearchForm>({
  title: '',
  statusFilter: 0,
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): ActivityMessageForm => ({
  id: '',
  iconEn: '',
  iconEs: '',
  iconPt: '',
  iconHi: '',
  iconId: '',
  bgEn: '',
  bgEs: '',
  bgPt: '',
  bgHi: '',
  bgId: '',
  titleEn: '',
  titleEs: '',
  titlePt: '',
  titleHi: '',
  titleId: '',
  contentEn: '',
  contentEs: '',
  contentPt: '',
  contentHi: '',
  contentId: '',
  urlEn: '',
  urlEs: '',
  urlPt: '',
  urlHi: '',
  urlId: '',
})
const currentRow = ref<ActivityMessageForm>(defaultForm())
const formRef = ref<FormInstance>()
const assetPreviewUrls = reactive<Record<AssetKind, Record<LangKey, string>>>({
  icon: emptyLangRecord(),
  bg: emptyLangRecord(),
})
const assetUploading = reactive<Record<AssetKind, Record<LangKey, boolean>>>({
  icon: emptyLangBoolRecord(),
  bg: emptyLangBoolRecord(),
})
const objectPreviewUrls: Record<string, string> = {}

const objectPreviewKey = (kind: AssetKind, lang: LangKey) => `${kind}:${lang}`

const revokeObjectPreviews = () => {
  for (const key of Object.keys(objectPreviewUrls)) {
    URL.revokeObjectURL(objectPreviewUrls[key])
    delete objectPreviewUrls[key]
  }
}

const resetAssetPreviews = () => {
  revokeObjectPreviews()
  Object.assign(assetPreviewUrls.icon, emptyLangRecord())
  Object.assign(assetPreviewUrls.bg, emptyLangRecord())
}

const setAssetPreview = (kind: AssetKind, lang: LangKey, url: string, fromObject = false) => {
  const key = objectPreviewKey(kind, lang)
  const prevObject = objectPreviewUrls[key]
  if (prevObject) {
    URL.revokeObjectURL(prevObject)
    delete objectPreviewUrls[key]
  }
  assetPreviewUrls[kind][lang] = url
  if (fromObject && url) {
    objectPreviewUrls[key] = url
  }
}

const clearAsset = (kind: AssetKind, lang: LangKey) => {
  currentRow.value[assetFieldMap[kind][lang]] = ''
  setAssetPreview(kind, lang, '')
}

const beforeImageUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.activityMessageList.imageOnly'))
    return false
  }
  return true
}

const doAssetUpload = async (kind: AssetKind, lang: LangKey, options: UploadRequestOptions) => {
  const file = options.file as File
  assetUploading[kind][lang] = true
  try {
    const res = await uploadApi.uploadFile(file)
    currentRow.value[assetFieldMap[kind][lang]] = res.fileName
    setAssetPreview(kind, lang, URL.createObjectURL(file), true)
    ElMessage.success(t('pages.activityMessageList.uploadSuccess'))
  } catch (error) {
    console.error('upload image failed:', error)
    ElMessage.error(t('pages.activityMessageList.uploadImageFailed'))
  } finally {
    assetUploading[kind][lang] = false
  }
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    resetAssetPreviews()
  }
})

const formRules = computed<FormRules>(() => ({
  titleEn: [
    {required: true, message: t('pages.activityMessageList.titleEnRequired'), trigger: 'blur'},
    {max: 128, message: t('pages.activityMessageList.titleMaxLength'), trigger: 'blur'},
  ],
  contentEn: [{required: true, message: t('pages.activityMessageList.contentEnRequired'), trigger: 'blur'}],
  urlEn: [{max: 512, message: t('pages.activityMessageList.urlMaxLength'), trigger: 'blur'}],
}))

const syncAssetPreviewsFromRow = (row: ActivityMessage) => {
  resetAssetPreviews()
  setAssetPreview('icon', 'en', row.iconEn || '')
  setAssetPreview('icon', 'es', row.iconEs || '')
  setAssetPreview('icon', 'pt', row.iconPt || '')
  setAssetPreview('icon', 'hi', row.iconHi || '')
  setAssetPreview('icon', 'id', row.iconId || '')
  setAssetPreview('bg', 'en', row.bgEn || '')
  setAssetPreview('bg', 'es', row.bgEs || '')
  setAssetPreview('bg', 'pt', row.bgPt || '')
  setAssetPreview('bg', 'hi', row.bgHi || '')
  setAssetPreview('bg', 'id', row.bgId || '')
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await activityMessageApi.getActivityMessageList({
      title: searchForm.title,
      statusFilter: searchForm.statusFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value,
    })
    tableData.value = response.data || []
    total.value = response.total || 0
  } catch (error) {
    console.error('fetch activity message list failed:', error)
    ElMessage.error(t('pages.activityMessageList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const resetSearch = () => {
  searchForm.title = ''
  searchForm.statusFilter = 0
  currentPage.value = 1
  fetchList()
}

const handleSizeChange = () => {
  currentPage.value = 1
  fetchList()
}

const handleCurrentChange = () => {
  fetchList()
}

const handleAdd = () => {
  dialogTitle.value = t('pages.activityMessageList.addMessage')
  currentRow.value = defaultForm()
  resetAssetPreviews()
  activeLangTab.value = 'en'
  dialogVisible.value = true
}

const handleEdit = (row: ActivityMessage) => {
  dialogTitle.value = t('pages.activityMessageList.editMessage')
  currentRow.value = {
    id: row.id,
    iconEn: row.iconEnName || '',
    iconEs: row.iconEsName || '',
    iconPt: row.iconPtName || '',
    iconHi: row.iconHiName || '',
    iconId: row.iconIdName || '',
    bgEn: row.bgEnName || '',
    bgEs: row.bgEsName || '',
    bgPt: row.bgPtName || '',
    bgHi: row.bgHiName || '',
    bgId: row.bgIdName || '',
    titleEn: row.titleEn || '',
    titleEs: row.titleEs || '',
    titlePt: row.titlePt || '',
    titleHi: row.titleHi || '',
    titleId: row.titleId || '',
    contentEn: row.contentEn || '',
    contentEs: row.contentEs || '',
    contentPt: row.contentPt || '',
    contentHi: row.contentHi || '',
    contentId: row.contentId || '',
    urlEn: row.urlEn || '',
    urlEs: row.urlEs || '',
    urlPt: row.urlPt || '',
    urlHi: row.urlHi || '',
    urlId: row.urlId || '',
  }
  syncAssetPreviewsFromRow(row)
  activeLangTab.value = 'en'
  dialogVisible.value = true
}

const handleSave = async () => {
  if (!formRef.value) {
    return
  }
  try {
    await formRef.value.validate()
  } catch {
    activeLangTab.value = 'en'
    return
  }
  try {
    const payload = {...currentRow.value}
    if (payload.id) {
      await activityMessageApi.updateActivityMessage(payload)
      ElMessage.success(t('common.updateSuccess'))
    } else {
      const {id, ...createData} = payload
      void id
      await activityMessageApi.createActivityMessage(createData)
      ElMessage.success(t('common.createSuccess'))
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('save activity message failed:', error)
    ElMessage.error(t('pages.activityMessageList.saveFailed'))
  }
}

const handleDelete = async (row: ActivityMessage) => {
  try {
    await ElMessageBox.confirm(
        t('pages.activityMessageList.deleteConfirm', {title: row.titleEn || row.id}),
        t('pages.activityMessageList.promptTitle'),
        {type: 'warning'}
    )
    await activityMessageApi.deleteActivityMessage(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('delete activity message failed:', error)
      ElMessage.error(t('pages.activityMessageList.deleteFailed'))
    }
  }
}

const handlePublish = async (row: ActivityMessage) => {
  try {
    await activityMessageApi.publishActivityMessage(row.id)
    ElMessage.success(t('pages.activityMessageList.publishSuccess'))
    fetchList()
  } catch (error) {
    console.error('publish failed:', error)
    ElMessage.error(t('pages.activityMessageList.publishFailed'))
  }
}

const handleUnpublish = async (row: ActivityMessage) => {
  try {
    await activityMessageApi.unpublishActivityMessage(row.id)
    ElMessage.success(t('pages.activityMessageList.unpublishSuccess'))
    fetchList()
  } catch (error) {
    console.error('unpublish failed:', error)
    ElMessage.error(t('pages.activityMessageList.unpublishFailed'))
  }
}

const handleSelectionChange = (rows: ActivityMessage[]) => {
  selectedRows.value = rows
}

const handleSyncData = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning(t('pages.activityMessageList.selectSyncFirst'))
    return
  }
  const ids = selectedRows.value.map((row) => Number(row.id)).filter((id) => id > 0)
  if (ids.length === 0) {
    ElMessage.warning(t('pages.activityMessageList.invalidSelection'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.activityMessageList.syncConfirm', {count: ids.length}),
        t('common.syncData'),
        {
          confirmButtonText: t('common.confirmSync'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        }
    )
    syncing.value = true
    const response = await dataSyncApi.syncActivityMessage({ids})
    if (response?.success) {
      ElMessage.success(
          response.message || t('pages.activityMessageList.syncSuccessDetail', {
            rows: response.rowCount,
            files: response.fileCount,
          })
      )
    } else {
      ElMessage.error(t('pages.activityMessageList.syncFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('sync failed:', error)
      ElMessage.error(t('pages.activityMessageList.syncFailedCheckConfig'))
    }
  } finally {
    syncing.value = false
  }
}

onMounted(() => {
  fetchList()
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

.table-header {
  margin-bottom: 16px;
}

.selection-tip {
  margin-bottom: 12px;
  font-size: 13px;
  color: var(--el-color-primary);
}

.search-form {
  margin-bottom: 16px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.icon-upload-wrap {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.icon-uploader {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  overflow: hidden;
}

.icon-uploader-placeholder {
  width: 120px;
  height: 120px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-secondary);
}

.icon-uploader-icon {
  font-size: 24px;
  margin-bottom: 8px;
}

.icon-preview {
  width: 120px;
  height: 120px;
  object-fit: cover;
  display: block;
}

.bg-uploader {
  border: 1px dashed var(--el-border-color);
  border-radius: 6px;
  cursor: pointer;
  overflow: hidden;
}

.bg-uploader-placeholder {
  width: 180px;
  height: 100px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: var(--el-text-color-secondary);
}

.bg-preview {
  width: 180px;
  height: 100px;
  object-fit: cover;
  display: block;
}
</style>
