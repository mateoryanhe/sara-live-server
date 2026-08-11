<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>活动消息管理</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增活动消息</el-button>
          <el-button
              v-if="hasButtonPermission('ActivityMessageManagement', 'sync')"
              :disabled="selectedRows.length === 0"
              :loading="syncing"
              type="warning"
              @click="handleSyncData"
          >
            同步数据
          </el-button>
        </div>

        <div v-if="selectedRows.length" class="selection-tip">已选 {{ selectedRows.length }} 项</div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item label="标题">
            <el-input v-model="searchForm.title" clearable placeholder="标题(模糊匹配)"/>
          </el-form-item>
          <el-form-item label="发布状态">
            <el-select v-model="searchForm.statusFilter" placeholder="全部" style="width: 140px">
              <el-option :value="0" label="全部"/>
              <el-option :value="2" label="只看已发布"/>
              <el-option :value="1" label="只看未发布"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleSearch">搜索</el-button>
            <el-button @click="resetSearch">重置</el-button>
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
          <el-table-column label="图标(英文)" width="100">
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
          <el-table-column label="背景图(英文)" width="110">
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
          <el-table-column label="标题(英文)" prop="titleEn" min-width="160" show-overflow-tooltip/>
          <el-table-column label="跳转链接(英文)" prop="urlEn" min-width="200" show-overflow-tooltip/>
          <el-table-column label="发布时间" prop="publishedAt" width="160">
            <template #default="{ row }">
              {{ row.publishedAt || '-' }}
            </template>
          </el-table-column>
          <el-table-column label="发布状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'">
                {{ row.status === 1 ? '已发布' : '未发布' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="createdAt" width="160"/>
          <el-table-column label="更新时间" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" label="操作" width="280">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">编辑</el-button>
              <el-button
                  v-if="row.status !== 1"
                  size="small"
                  type="success"
                  @click="handlePublish(row)"
              >
                发布
              </el-button>
              <el-button
                  v-else
                  size="small"
                  type="warning"
                  @click="handleUnpublish(row)"
              >
                取消发布
              </el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
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
            <el-form-item label="图标" prop="iconEn">
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
                    <span>点击上传图标</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.icon.en || currentRow.iconEn" link type="danger" @click="clearAsset('icon', 'en')">
                  移除图标
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="背景图">
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
                    <span>点击上传背景图</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.bg.en || currentRow.bgEn" link type="danger" @click="clearAsset('bg', 'en')">
                  移除背景图
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="标题" prop="titleEn">
              <el-input v-model="currentRow.titleEn" maxlength="128" placeholder="英文标题" show-word-limit/>
            </el-form-item>
            <el-form-item label="内容" prop="contentEn">
              <el-input v-model="currentRow.contentEn" :rows="5" placeholder="英文内容" type="textarea"/>
            </el-form-item>
            <el-form-item label="跳转链接" prop="urlEn">
              <el-input v-model="currentRow.urlEn" maxlength="512" placeholder="英文跳转链接" show-word-limit/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="Español" name="es">
            <el-form-item label="图标">
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
                    <span>点击上传图标</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.icon.es || currentRow.iconEs" link type="danger" @click="clearAsset('icon', 'es')">
                  移除图标
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="背景图">
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
                    <span>点击上传背景图</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.bg.es || currentRow.bgEs" link type="danger" @click="clearAsset('bg', 'es')">
                  移除背景图
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="标题">
              <el-input v-model="currentRow.titleEs" maxlength="128" placeholder="西班牙语标题" show-word-limit/>
            </el-form-item>
            <el-form-item label="内容">
              <el-input v-model="currentRow.contentEs" :rows="5" placeholder="西班牙语内容" type="textarea"/>
            </el-form-item>
            <el-form-item label="跳转链接">
              <el-input v-model="currentRow.urlEs" maxlength="512" placeholder="西班牙语跳转链接" show-word-limit/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="Português" name="pt">
            <el-form-item label="图标">
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
                    <span>点击上传图标</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.icon.pt || currentRow.iconPt" link type="danger" @click="clearAsset('icon', 'pt')">
                  移除图标
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="背景图">
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
                    <span>点击上传背景图</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.bg.pt || currentRow.bgPt" link type="danger" @click="clearAsset('bg', 'pt')">
                  移除背景图
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="标题">
              <el-input v-model="currentRow.titlePt" maxlength="128" placeholder="葡萄牙语标题" show-word-limit/>
            </el-form-item>
            <el-form-item label="内容">
              <el-input v-model="currentRow.contentPt" :rows="5" placeholder="葡萄牙语内容" type="textarea"/>
            </el-form-item>
            <el-form-item label="跳转链接">
              <el-input v-model="currentRow.urlPt" maxlength="512" placeholder="葡萄牙语跳转链接" show-word-limit/>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="हिन्दी" name="hi">
            <el-form-item label="图标">
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
                    <span>点击上传图标</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.icon.hi || currentRow.iconHi" link type="danger" @click="clearAsset('icon', 'hi')">
                  移除图标
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="背景图">
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
                    <span>点击上传背景图</span>
                  </div>
                </el-upload>
                <el-button v-if="assetPreviewUrls.bg.hi || currentRow.bgHi" link type="danger" @click="clearAsset('bg', 'hi')">
                  移除背景图
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="标题">
              <el-input v-model="currentRow.titleHi" maxlength="128" placeholder="印地语标题" show-word-limit/>
            </el-form-item>
            <el-form-item label="内容">
              <el-input v-model="currentRow.contentHi" :rows="5" placeholder="印地语内容" type="textarea"/>
            </el-form-item>
            <el-form-item label="跳转链接">
              <el-input v-model="currentRow.urlHi" maxlength="512" placeholder="印地语跳转链接" show-word-limit/>
            </el-form-item>
          </el-tab-pane>
        </el-tabs>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref, watch} from 'vue'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {activityMessageApi, type ActivityMessageForm} from '@/api/modules/activityMessage'
import {dataSyncApi, uploadApi} from '@/api'
import type {ActivityMessage} from '@/types/api'
import {hasButtonPermission} from '@/utils/permission'

type LangKey = 'en' | 'es' | 'pt' | 'hi'
type AssetKind = 'icon' | 'bg'

const assetFieldMap: Record<AssetKind, Record<LangKey, keyof ActivityMessageForm>> = {
  icon: {en: 'iconEn', es: 'iconEs', pt: 'iconPt', hi: 'iconHi'},
  bg: {en: 'bgEn', es: 'bgEs', pt: 'bgPt', hi: 'bgHi'},
}

const emptyLangRecord = () => ({en: '', es: '', pt: '', hi: ''})
const emptyLangBoolRecord = () => ({en: false, es: false, pt: false, hi: false})

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
  bgEn: '',
  bgEs: '',
  bgPt: '',
  bgHi: '',
  titleEn: '',
  titleEs: '',
  titlePt: '',
  titleHi: '',
  contentEn: '',
  contentEs: '',
  contentPt: '',
  contentHi: '',
  urlEn: '',
  urlEs: '',
  urlPt: '',
  urlHi: '',
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
    ElMessage.error('只能上传图片文件')
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
    ElMessage.success('上传成功')
  } catch (error) {
    console.error('上传图片失败:', error)
    ElMessage.error('上传图片失败')
  } finally {
    assetUploading[kind][lang] = false
  }
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    resetAssetPreviews()
  }
})

const formRules: FormRules = {
  titleEn: [
    {required: true, message: '请输入英文标题', trigger: 'blur'},
    {max: 128, message: '标题最长128字符', trigger: 'blur'},
  ],
  contentEn: [{required: true, message: '请输入英文内容', trigger: 'blur'}],
  urlEn: [{max: 512, message: '跳转链接最长512字符', trigger: 'blur'}],
}

const syncAssetPreviewsFromRow = (row: ActivityMessage) => {
  resetAssetPreviews()
  setAssetPreview('icon', 'en', row.iconEn || '')
  setAssetPreview('icon', 'es', row.iconEs || '')
  setAssetPreview('icon', 'pt', row.iconPt || '')
  setAssetPreview('icon', 'hi', row.iconHi || '')
  setAssetPreview('bg', 'en', row.bgEn || '')
  setAssetPreview('bg', 'es', row.bgEs || '')
  setAssetPreview('bg', 'pt', row.bgPt || '')
  setAssetPreview('bg', 'hi', row.bgHi || '')
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
    console.error('获取活动消息列表失败:', error)
    ElMessage.error('获取活动消息列表失败')
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
  dialogTitle.value = '新增活动消息'
  currentRow.value = defaultForm()
  resetAssetPreviews()
  activeLangTab.value = 'en'
  dialogVisible.value = true
}

const handleEdit = (row: ActivityMessage) => {
  dialogTitle.value = '编辑活动消息'
  currentRow.value = {
    id: row.id,
    iconEn: row.iconEnName || '',
    iconEs: row.iconEsName || '',
    iconPt: row.iconPtName || '',
    iconHi: row.iconHiName || '',
    bgEn: row.bgEnName || '',
    bgEs: row.bgEsName || '',
    bgPt: row.bgPtName || '',
    bgHi: row.bgHiName || '',
    titleEn: row.titleEn || '',
    titleEs: row.titleEs || '',
    titlePt: row.titlePt || '',
    titleHi: row.titleHi || '',
    contentEn: row.contentEn || '',
    contentEs: row.contentEs || '',
    contentPt: row.contentPt || '',
    contentHi: row.contentHi || '',
    urlEn: row.urlEn || '',
    urlEs: row.urlEs || '',
    urlPt: row.urlPt || '',
    urlHi: row.urlHi || '',
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
      ElMessage.success('更新成功')
    } else {
      const {id, ...createData} = payload
      void id
      await activityMessageApi.createActivityMessage(createData)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('保存活动消息失败:', error)
    ElMessage.error('保存活动消息失败')
  }
}

const handleDelete = async (row: ActivityMessage) => {
  try {
    await ElMessageBox.confirm(`确定删除活动消息「${row.titleEn || row.id}」吗？`, '提示', {type: 'warning'})
    await activityMessageApi.deleteActivityMessage(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除活动消息失败:', error)
      ElMessage.error('删除活动消息失败')
    }
  }
}

const handlePublish = async (row: ActivityMessage) => {
  try {
    await activityMessageApi.publishActivityMessage(row.id)
    ElMessage.success('发布成功')
    fetchList()
  } catch (error) {
    console.error('发布失败:', error)
    ElMessage.error('发布失败')
  }
}

const handleUnpublish = async (row: ActivityMessage) => {
  try {
    await activityMessageApi.unpublishActivityMessage(row.id)
    ElMessage.success('已取消发布')
    fetchList()
  } catch (error) {
    console.error('取消发布失败:', error)
    ElMessage.error('取消发布失败')
  }
}

const handleSelectionChange = (rows: ActivityMessage[]) => {
  selectedRows.value = rows
}

const handleSyncData = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先勾选要同步的活动消息')
    return
  }
  const ids = selectedRows.value.map((row) => Number(row.id)).filter((id) => id > 0)
  if (ids.length === 0) {
    ElMessage.warning('所选活动消息无效')
    return
  }
  try {
    await ElMessageBox.confirm(
        `将把已选 ${ids.length} 条活动消息及关联图片资源同步到目标环境（按 ID 覆盖或新增）。是否继续？`,
        '同步数据',
        {
          confirmButtonText: '确定同步',
          cancelButtonText: '取消',
          type: 'warning',
        }
    )
    syncing.value = true
    const response = await dataSyncApi.syncActivityMessage({ids})
    if (response?.success) {
      ElMessage.success(response.message || `同步成功：${response.rowCount} 条，${response.fileCount} 个文件`)
    } else {
      ElMessage.error('同步失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('同步失败:', error)
      ElMessage.error('同步失败，请检查数据同步配置与目标服务')
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
