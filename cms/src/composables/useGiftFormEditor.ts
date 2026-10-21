import {computed, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {giftApi, uploadApi} from '@/api'
import type {Gift} from '@/types/api.ts'
import {
  getExt,
  resolveFilePreviewType,
  resolveMediaPreviewType,
  type MediaPreviewType,
} from '@/utils/media-preview'
import {truncateNumber} from '@/utils/number-format'
import {
  GIFT_CATEGORY_LUCKY,
  GIFT_CATEGORY_NORMAL,
  isKnownGiftCategoryType,
} from '@/utils/gift-category'
import {
  createDefaultLuckyGiftMultiplierTiers,
  type LuckyGiftMultiplierTier,
} from '@/utils/lucky-gift-multiplier'

export interface GiftForm {
  id: string
  name: string
  nameEn: string
  nameEs: string
  namePt: string
  nameHi: string
  nameId: string
  icon: string
  animation: string
  price: number
  categoryType: number
  luckyGiftRtp: number
  luckyGiftMultiplierTiers: LuckyGiftMultiplierTier[]
  sort: number
  description: string
}

export function defaultGiftForm(): GiftForm {
  return {
    id: '',
    name: '',
    nameEn: '',
    nameEs: '',
    namePt: '',
    nameHi: '',
    nameId: '',
    icon: '',
    animation: '',
    price: 0,
    categoryType: GIFT_CATEGORY_NORMAL,
    luckyGiftRtp: 0,
    luckyGiftMultiplierTiers: createDefaultLuckyGiftMultiplierTiers(),
    sort: 0,
    description: '',
  }
}

const mapApiTiersToForm = (tiers: {multiplier: number; weight: number; type: string}[]) => {
  if (!tiers?.length) {
    return createDefaultLuckyGiftMultiplierTiers()
  }
  return tiers.map((tier, index) => ({
    key: `tier-${index}`,
    multiplier: tier.multiplier,
    weight: tier.weight,
    type: tier.type as LuckyGiftMultiplierTier['type'],
  }))
}

export function useGiftFormEditor() {
  const {t} = useI18n()

  const pageLoading = ref(false)
  const saving = ref(false)
  const activeTab = ref('basic')
  const formRef = ref<FormInstance>()
  const currentRow = ref<GiftForm>(defaultGiftForm())

  const isLuckyGift = computed(() => currentRow.value.categoryType === GIFT_CATEGORY_LUCKY)

  const iconUploading = ref(false)
  const animationUploading = ref(false)
  const iconPreviewUrl = ref('')
  const animationPreviewUrl = ref('')
  const animationPreviewType = ref<MediaPreviewType>('none')
  const animationFileLabel = ref('')
  const objectPreviewUrls: Partial<Record<'icon' | 'animation', string>> = {}

  const revokeObjectPreview = (field: 'icon' | 'animation') => {
    const url = objectPreviewUrls[field]
    if (url) {
      URL.revokeObjectURL(url)
      delete objectPreviewUrls[field]
    }
  }

  const revokeAllObjectPreviews = () => {
    revokeObjectPreview('icon')
    revokeObjectPreview('animation')
  }

  const resetAssetPreview = () => {
    revokeAllObjectPreviews()
    iconPreviewUrl.value = ''
    animationPreviewUrl.value = ''
    animationPreviewType.value = 'none'
    animationFileLabel.value = ''
  }

  const setIconPreview = (url: string, fromObject = false) => {
    revokeObjectPreview('icon')
    iconPreviewUrl.value = url
    if (fromObject && url) {
      objectPreviewUrls.icon = url
    }
  }

  const setAnimationPreview = (
    url: string,
    fileLabel: string,
    type?: MediaPreviewType,
    fromObject = false,
  ) => {
    revokeObjectPreview('animation')
    const previewType = type ?? resolveMediaPreviewType(url, fileLabel)
    animationPreviewType.value = previewType
    animationPreviewUrl.value = previewType === 'image' || previewType === 'video' ? url : ''
    animationFileLabel.value = previewType === 'file' ? fileLabel : ''
    if (fromObject && url) {
      objectPreviewUrls.animation = url
    }
  }

  const clearAsset = (field: 'icon' | 'animation') => {
    currentRow.value[field] = ''
    if (field === 'icon') {
      setIconPreview('')
    } else {
      setAnimationPreview('', '', 'none')
    }
    formRef.value?.validateField(field).catch(() => undefined)
  }

  watch(
    () => currentRow.value.categoryType,
    (categoryType) => {
      if (categoryType === GIFT_CATEGORY_LUCKY && currentRow.value.luckyGiftMultiplierTiers.length === 0) {
        currentRow.value.luckyGiftMultiplierTiers = createDefaultLuckyGiftMultiplierTiers()
      }
    },
  )

  const allowedAnimationExt = [
    '.jpg', '.jpeg', '.png', '.gif', '.webp', '.bmp', '.apng',
    '.svga', '.pag', '.json', '.lottie', '.mp4', '.webm', '.zip',
  ]

  const beforeIconUpload = (file: File): boolean => {
    if (!file.type.startsWith('image/')) {
      ElMessage.error(t('pages.giftList.iconImageOnly'))
      return false
    }
    return true
  }

  const beforeAnimationUpload = (file: File): boolean => {
    const ext = getExt(file.name)
    if (!allowedAnimationExt.includes(ext)) {
      ElMessage.error(t('pages.giftList.unsupportedFileType', {ext: ext || '-'}))
      return false
    }
    return true
  }

  const doUpload = async (options: UploadRequestOptions, field: 'icon' | 'animation') => {
    const file = options.file as File
    const flag = field === 'icon' ? iconUploading : animationUploading
    flag.value = true
    try {
      const res = await uploadApi.uploadFile(file)
      currentRow.value[field] = res.fileName
      if (field === 'icon') {
        setIconPreview(URL.createObjectURL(file), true)
      } else {
        const previewType = resolveFilePreviewType(file.name)
        if (previewType === 'image' || previewType === 'video') {
          setAnimationPreview(URL.createObjectURL(file), '', previewType, true)
        } else {
          setAnimationPreview('', res.fileName, 'file')
        }
      }
      formRef.value?.validateField(field).catch(() => undefined)
      ElMessage.success(t('pages.giftList.uploadSuccess'))
    } catch (error) {
      console.error('upload failed:', error)
      ElMessage.error(t('pages.giftList.uploadFailed'))
    } finally {
      flag.value = false
    }
  }

  const formRules = computed<FormRules>(() => ({
    name: [
      {required: true, message: t('pages.giftList.nameRequired'), trigger: 'blur'},
      {min: 1, max: 64, message: t('pages.giftList.nameLength'), trigger: 'blur'},
    ],
    nameEn: [{max: 64, message: t('pages.giftList.nameEnMaxLength'), trigger: 'blur'}],
    nameEs: [{max: 64, message: t('pages.giftList.nameEsMaxLength'), trigger: 'blur'}],
    namePt: [{max: 64, message: t('pages.giftList.namePtMaxLength'), trigger: 'blur'}],
    nameHi: [{max: 64, message: t('pages.giftList.nameHiMaxLength'), trigger: 'blur'}],
    nameId: [{max: 64, message: t('pages.giftList.nameIdMaxLength'), trigger: 'blur'}],
    categoryType: [
      {required: true, message: t('pages.giftList.categoryRequired'), trigger: 'change'},
      {
        validator: (_rule, value, callback) => {
          if (isKnownGiftCategoryType(Number(value))) {
            callback()
            return
          }
          callback(new Error(t('pages.giftList.categoryRequired')))
        },
        trigger: 'change',
      },
    ],
    description: [{max: 255, message: t('pages.giftList.descriptionMaxLength'), trigger: 'blur'}],
  }))

  const applyGiftRowToForm = async (row: Gift) => {
    const iconName = row.iconName || ''
    const animationName = row.animationName || ''
    const categoryType = isKnownGiftCategoryType(row.categoryType) ? row.categoryType : GIFT_CATEGORY_NORMAL
    let luckyGiftRtp = truncateNumber(row.luckyGiftRtp ?? 0)
    let luckyGiftMultiplierTiers = createDefaultLuckyGiftMultiplierTiers()
    if (categoryType === GIFT_CATEGORY_LUCKY) {
      try {
        const luckyCfg = await giftApi.getGiftLuckyCfg(row.id)
        luckyGiftRtp = truncateNumber(luckyCfg.luckyGiftRtp ?? 0)
        luckyGiftMultiplierTiers = mapApiTiersToForm(luckyCfg.tiers)
      } catch (error) {
        console.error('load gift lucky cfg failed:', error)
        ElMessage.warning(t('pages.giftList.luckyCfgLoadFailed'))
      }
    }
    currentRow.value = {
      id: row.id,
      name: row.name,
      nameEn: row.nameEn || '',
      nameEs: row.nameEs || '',
      namePt: row.namePt || '',
      nameHi: row.nameHi || '',
      nameId: row.nameId || '',
      icon: iconName,
      animation: animationName,
      price: truncateNumber(row.price),
      categoryType,
      luckyGiftRtp,
      luckyGiftMultiplierTiers,
      sort: Number(row.sort) || 0,
      description: row.description,
    }
    setIconPreview(row.icon || '')
    if (animationName) {
      const previewType = resolveMediaPreviewType(row.animation || '', animationName)
      if (previewType === 'image' || previewType === 'video') {
        setAnimationPreview(row.animation || '', '', previewType)
      } else {
        setAnimationPreview('', animationName, 'file')
      }
    } else {
      setAnimationPreview('', '', 'none')
    }
  }

  const loadGiftById = async (id: string) => {
    pageLoading.value = true
    try {
      const response = await giftApi.getGiftList({
        pageIndex: 1,
        pageSize: 500,
        statusFilter: 0,
      })
      const row = response.data.find((item) => String(item.id) === String(id))
      if (!row) {
        ElMessage.error(t('pages.giftList.giftNotFound'))
        return false
      }
      await applyGiftRowToForm(row)
      return true
    } catch (error) {
      console.error('load gift failed:', error)
      ElMessage.error(t('pages.giftList.fetchFailed'))
      return false
    } finally {
      pageLoading.value = false
    }
  }

  const initCreateForm = () => {
    currentRow.value = defaultGiftForm()
    resetAssetPreview()
    activeTab.value = 'basic'
  }

  const buildLuckyPayload = () => {
    if (currentRow.value.categoryType !== GIFT_CATEGORY_LUCKY) {
      return {}
    }
    for (const tier of currentRow.value.luckyGiftMultiplierTiers) {
      if (!tier.multiplier || tier.multiplier <= 0 || !tier.weight || tier.weight <= 0) {
        ElMessage.warning(t('pages.giftList.luckyTierInvalid'))
        return null
      }
    }
    return {
      luckyGiftRtp: currentRow.value.luckyGiftRtp,
      luckyGiftMultiplierTiers: currentRow.value.luckyGiftMultiplierTiers.map((tier) => ({
        multiplier: tier.multiplier,
        weight: tier.weight,
        type: tier.type,
      })),
    }
  }

  const saveGift = async (): Promise<boolean> => {
    if (!formRef.value) {
      return false
    }
    try {
      await formRef.value.validate()
    } catch {
      ElMessage.warning(t('pages.giftList.formInvalid'))
      return false
    }
    const luckyPayload = buildLuckyPayload()
    if (luckyPayload === null) {
      activeTab.value = 'multiplier'
      return false
    }
    saving.value = true
    try {
      const {id, name, nameEn, nameEs, namePt, nameHi, nameId, icon, animation, price, categoryType, sort, description} =
        currentRow.value
      if (id) {
        await giftApi.updateGift({id, name, nameEn, nameEs, namePt, nameHi, nameId, icon, animation, price, categoryType, sort, description, ...luckyPayload})
      } else {
        await giftApi.createGift({name, nameEn, nameEs, namePt, nameHi, nameId, icon, animation, price, categoryType, sort, description, ...luckyPayload})
      }
      ElMessage.success(id ? t('common.updateSuccess') : t('common.createSuccess'))
      return true
    } catch (error) {
      console.error('save gift failed:', error)
      ElMessage.error(currentRow.value.id ? t('pages.giftList.updateFailed') : t('pages.giftList.createFailed'))
      return false
    } finally {
      saving.value = false
    }
  }

  const dispose = () => {
    resetAssetPreview()
  }

  return {
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
  }
}
