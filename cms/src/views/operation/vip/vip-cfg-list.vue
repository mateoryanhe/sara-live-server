<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.VipCfgManagement') }}</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">{{ t('pages.vipCfgList.addVipLevel') }}</el-button>
          <el-button
              v-if="hasButtonPermission('VipCfgManagement', 'sync')"
              :disabled="selectedRows.length === 0"
              :loading="syncing"
              type="warning"
              @click="handleSyncData"
          >
            {{ t('common.syncData') }}
          </el-button>
        </div>

        <div v-if="selectedRows.length" class="selection-tip">{{ t('common.selectedCount', {count: selectedRows.length}) }}</div>

        <el-form :model="searchForm" class="search-form" inline>
          <el-form-item :label="t('pages.vipCfgList.levelName')">
            <el-input v-model="searchForm.levelName" clearable :placeholder="t('pages.vipCfgList.levelNameFuzzy')"/>
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
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column :label="t('common.level')" prop="level" width="80"/>
          <el-table-column :label="t('pages.vipCfgList.levelName')" prop="levelName" min-width="120"/>
          <el-table-column :label="t('pages.vipCfgList.levelIcon')" width="100">
            <template #default="{ row }">
              <el-image
                  v-if="row.levelIcon"
                  :preview-src-list="[row.levelIcon]"
                  :src="row.levelIcon"
                  class="table-icon-preview"
                  fit="cover"
                  preview-teleported
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.vipCfgList.entryEffect')" min-width="200">
            <template #default="{ row }">
              <video
                  v-if="isVideoUrl(row.animation)"
                  :src="row.animation"
                  class="table-media-preview"
                  controls
                  preload="metadata"
              />
              <span v-else class="media-url-text">{{ row.animationName || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.vipCfgList.entryEffectIcon')" width="110">
            <template #default="{ row }">
              <el-image
                  v-if="row.animationIcon"
                  :preview-src-list="[row.animationIcon]"
                  :src="row.animationIcon"
                  class="table-icon-preview"
                  fit="cover"
                  preview-teleported
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.vipCfgList.entryEffectSwitch')" width="120">
            <template #default="{ row }">
              <el-tag :type="row.animationSwitch === 1 ? 'success' : 'info'">
                {{ row.animationSwitch === 1 ? t('common.on') : t('common.off') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.vipCfgList.commentEffect')" min-width="200">
            <template #default="{ row }">
              <video
                  v-if="isVideoUrl(row.commentEffect)"
                  :src="row.commentEffect"
                  class="table-media-preview"
                  controls
                  preload="metadata"
              />
              <span v-else class="media-url-text">{{ row.commentEffectName || '-' }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.vipCfgList.commentEffectIcon')" width="130">
            <template #default="{ row }">
              <el-image
                  v-if="row.commentEffectIcon"
                  :preview-src-list="[row.commentEffectIcon]"
                  :src="row.commentEffectIcon"
                  class="table-icon-preview"
                  fit="cover"
                  preview-teleported
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.vipCfgList.commentEffectSwitch')" width="140">
            <template #default="{ row }">
              <el-tag :type="row.commentEffectSwitch === 1 ? 'success' : 'info'">
                {{ row.commentEffectSwitch === 1 ? t('common.on') : t('common.off') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.vipCfgList.customerServiceIcon')" width="120">
            <template #default="{ row }">
              <el-image
                  v-if="row.customerServiceIcon"
                  :preview-src-list="[row.customerServiceIcon]"
                  :src="row.customerServiceIcon"
                  class="table-icon-preview"
                  fit="cover"
                  preview-teleported
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.vipCfgList.customerServiceSwitch')" width="120">
            <template #default="{ row }">
              <el-tag :type="row.customerServiceSwitch === 1 ? 'success' : 'info'">
                {{ row.customerServiceSwitch === 1 ? t('common.on') : t('common.off') }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.vipCfgList.upgradeRechargeLimit')" width="130">
            <template #default="{ row }">
              {{ formatAmount(row.upgradeRechargeLimit) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="160"/>
          <el-table-column :label="t('common.updatedAt')" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" :label="t('common.actions')" width="180">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" destroy-on-close width="720px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="130px">
        <el-tabs v-model="activeTab">
          <el-tab-pane :label="t('pages.vipCfgList.tabBasic')" name="basic">
            <el-form-item :label="t('common.level')" prop="level">
              <el-input-number v-model="currentRow.level" :min="1" controls-position="right"/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.levelName')" prop="levelName">
              <el-input v-model="currentRow.levelName" :placeholder="t('pages.vipCfgList.levelNamePlaceholder')"/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.levelIcon')" prop="levelIcon">
              <div class="asset-upload-wrap">
                <el-upload
                    :before-upload="beforeIconUpload"
                    :disabled="levelIconUploading"
                    :http-request="(opt) => doAssetUpload(opt, 'levelIcon')"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="icon-uploader"
                >
                  <img
                      v-if="levelIconPreviewUrl"
                      :src="levelIconPreviewUrl"
                      alt="level icon"
                      class="icon-preview"
                  />
                  <div v-else class="asset-uploader-placeholder icon-placeholder">
                    <el-icon class="asset-uploader-icon">
                      <Plus/>
                    </el-icon>
                    <span>{{ t('pages.vipCfgList.clickUploadIcon') }}</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="levelIconPreviewUrl || currentRow.levelIcon"
                    link
                    type="danger"
                    @click="clearLevelIcon"
                >
                  {{ t('pages.vipCfgList.removeIcon') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.upgradeRechargeLimit')" prop="upgradeRechargeLimit">
              <el-input-number
                  v-model="currentRow.upgradeRechargeLimit"
                  :min="0"
                  :precision="NUMBER_INPUT_DECIMALS"
                  :step="0.0001"
                  controls-position="right"
              />
              <div class="form-tip">{{ t('pages.vipCfgList.decimalExample') }}</div>
            </el-form-item>
            <el-alert :closable="false" show-icon type="info" class="form-tip-block">
              {{ t('pages.vipCfgList.titleEditTip') }}
            </el-alert>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.vipCfgList.tabEntryEffect')" name="entryEffect">
            <el-form-item :label="t('pages.vipCfgList.entryEffectSwitch')" prop="animationSwitch">
              <el-radio-group v-model="currentRow.animationSwitch">
                <el-radio :label="1">{{ t('common.on') }}</el-radio>
                <el-radio :label="0">{{ t('common.off') }}</el-radio>
              </el-radio-group>
              <div class="form-tip">{{ t('pages.vipCfgList.appDisplayTip') }}</div>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectAnimation')" prop="animation">
              <div class="asset-upload-wrap">
                <el-upload
                    :before-upload="beforeAnimationUpload"
                    :disabled="animationUploading"
                    :http-request="(opt) => doAssetUpload(opt, 'animation')"
                    :show-file-list="false"
                    accept=".mp4"
                    action="#"
                    class="animation-uploader"
                >
                  <video
                      v-if="animationPreviewUrl"
                      :src="animationPreviewUrl"
                      class="animation-preview"
                      controls
                      preload="metadata"
                  />
                  <div v-else class="asset-uploader-placeholder animation-placeholder">
                    <el-icon class="asset-uploader-icon">
                      <Plus/>
                    </el-icon>
                    <span>{{ t('pages.vipCfgList.clickUploadMp4') }}</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="animationPreviewUrl || currentRow.animation"
                    link
                    type="danger"
                    @click="clearAnimation"
                >
                  {{ t('pages.vipCfgList.removeAnimation') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectIcon')" prop="animationIcon">
              <div class="asset-upload-wrap">
                <el-upload
                    :before-upload="beforeIconUpload"
                    :disabled="animationIconUploading"
                    :http-request="(opt) => doAssetUpload(opt, 'animationIcon')"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="icon-uploader"
                >
                  <img
                      v-if="animationIconPreviewUrl"
                      :src="animationIconPreviewUrl"
                      alt="animation icon"
                      class="icon-preview"
                  />
                  <div v-else class="asset-uploader-placeholder icon-placeholder">
                    <el-icon class="asset-uploader-icon">
                      <Plus/>
                    </el-icon>
                    <span>{{ t('pages.vipCfgList.clickUploadIcon') }}</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="animationIconPreviewUrl || currentRow.animationIcon"
                    link
                    type="danger"
                    @click="clearAnimationIcon"
                >
                  {{ t('pages.vipCfgList.removeIcon') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectTitleEn')" prop="animationTitleEn">
              <el-input v-model="currentRow.animationTitleEn" maxlength="128" placeholder="Entry effect title in English" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectTitleEs')" prop="animationTitleEs">
              <el-input v-model="currentRow.animationTitleEs" maxlength="128" placeholder="Título del efecto de entrada en español" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectTitlePt')" prop="animationTitlePt">
              <el-input v-model="currentRow.animationTitlePt" maxlength="128" placeholder="Título do efeito de entrada em português" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectTitleHi')" prop="animationTitleHi">
              <el-input v-model="currentRow.animationTitleHi" maxlength="128" placeholder="प्रवेश प्रभाव शीर्षक हिंदी में" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectTitleId')" prop="animationTitleId">
              <el-input v-model="currentRow.animationTitleId" maxlength="128" placeholder="Judul efek masuk dalam Bahasa Indonesia" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectDescEn')" prop="animationDescEn">
              <el-input
                  v-model="currentRow.animationDescEn"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Entry effect description in English"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectDescEs')" prop="animationDescEs">
              <el-input
                  v-model="currentRow.animationDescEs"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descripción del efecto de entrada en español"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectDescPt')" prop="animationDescPt">
              <el-input
                  v-model="currentRow.animationDescPt"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descrição do efeito de entrada em português"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectDescHi')" prop="animationDescHi">
              <el-input
                  v-model="currentRow.animationDescHi"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="प्रवेश प्रभाव विवरण हिंदी में"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectDescId')" prop="animationDescId">
              <el-input
                  v-model="currentRow.animationDescId"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Deskripsi efek masuk dalam Bahasa Indonesia"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.vipCfgList.tabCommentEffect')" name="commentEffect">
            <el-form-item :label="t('pages.vipCfgList.commentEffectSwitch')" prop="commentEffectSwitch">
              <el-radio-group v-model="currentRow.commentEffectSwitch">
                <el-radio :label="1">{{ t('common.on') }}</el-radio>
                <el-radio :label="0">{{ t('common.off') }}</el-radio>
              </el-radio-group>
              <div class="form-tip">{{ t('pages.vipCfgList.appDisplayCommentTip') }}</div>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectAnimation')" prop="commentEffect">
              <div class="asset-upload-wrap">
                <el-upload
                    :before-upload="beforeAnimationUpload"
                    :disabled="commentEffectUploading"
                    :http-request="(opt) => doAssetUpload(opt, 'commentEffect')"
                    :show-file-list="false"
                    accept=".mp4"
                    action="#"
                    class="animation-uploader"
                >
                  <video
                      v-if="commentEffectPreviewUrl"
                      :src="commentEffectPreviewUrl"
                      class="animation-preview"
                      controls
                      preload="metadata"
                  />
                  <div v-else class="asset-uploader-placeholder animation-placeholder">
                    <el-icon class="asset-uploader-icon">
                      <Plus/>
                    </el-icon>
                    <span>{{ t('pages.vipCfgList.clickUploadMp4') }}</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="commentEffectPreviewUrl || currentRow.commentEffect"
                    link
                    type="danger"
                    @click="clearCommentEffect"
                >
                  {{ t('pages.vipCfgList.removeAnimation') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectIcon')" prop="commentEffectIcon">
              <div class="asset-upload-wrap">
                <el-upload
                    :before-upload="beforeIconUpload"
                    :disabled="commentEffectIconUploading"
                    :http-request="(opt) => doAssetUpload(opt, 'commentEffectIcon')"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="icon-uploader"
                >
                  <img
                      v-if="commentEffectIconPreviewUrl"
                      :src="commentEffectIconPreviewUrl"
                      alt="comment effect icon"
                      class="icon-preview"
                  />
                  <div v-else class="asset-uploader-placeholder icon-placeholder">
                    <el-icon class="asset-uploader-icon">
                      <Plus/>
                    </el-icon>
                    <span>{{ t('pages.vipCfgList.clickUploadIcon') }}</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="commentEffectIconPreviewUrl || currentRow.commentEffectIcon"
                    link
                    type="danger"
                    @click="clearCommentEffectIcon"
                >
                  {{ t('pages.vipCfgList.removeIcon') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectTitleEn')" prop="commentEffectTitleEn">
              <el-input v-model="currentRow.commentEffectTitleEn" maxlength="128" placeholder="Comment effect title in English" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectTitleEs')" prop="commentEffectTitleEs">
              <el-input v-model="currentRow.commentEffectTitleEs" maxlength="128" placeholder="Título del efecto de comentario en español" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectTitlePt')" prop="commentEffectTitlePt">
              <el-input v-model="currentRow.commentEffectTitlePt" maxlength="128" placeholder="Título do efeito de comentário em português" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectTitleHi')" prop="commentEffectTitleHi">
              <el-input v-model="currentRow.commentEffectTitleHi" maxlength="128" placeholder="टिप्पणी प्रभाव शीर्षक हिंदी में" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectTitleId')" prop="commentEffectTitleId">
              <el-input v-model="currentRow.commentEffectTitleId" maxlength="128" placeholder="Judul efek komentar dalam Bahasa Indonesia" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectDescEn')" prop="commentEffectDescEn">
              <el-input
                  v-model="currentRow.commentEffectDescEn"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Public screen comment effect description in English"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectDescEs')" prop="commentEffectDescEs">
              <el-input
                  v-model="currentRow.commentEffectDescEs"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descripción del efecto de comentario en pantalla pública en español"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectDescPt')" prop="commentEffectDescPt">
              <el-input
                  v-model="currentRow.commentEffectDescPt"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descrição do efeito de comentário na tela pública em português"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectDescHi')" prop="commentEffectDescHi">
              <el-input
                  v-model="currentRow.commentEffectDescHi"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="सार्वजनिक स्क्रीन टिप्पणी प्रभाव विवरण हिंदी में"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.effectDescId')" prop="commentEffectDescId">
              <el-input
                  v-model="currentRow.commentEffectDescId"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Deskripsi efek komentar layar publik dalam Bahasa Indonesia"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane :label="t('pages.vipCfgList.tabCustomerService')" name="customerService">
            <el-form-item :label="t('pages.vipCfgList.customerServiceSwitch')" prop="customerServiceSwitch">
              <el-radio-group v-model="currentRow.customerServiceSwitch">
                <el-radio :label="1">{{ t('common.on') }}</el-radio>
                <el-radio :label="0">{{ t('common.off') }}</el-radio>
              </el-radio-group>
              <div class="form-tip">{{ t('pages.vipCfgList.appDisplayCsTip') }}</div>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.privilegeIcon')" prop="customerServiceIcon">
              <div class="asset-upload-wrap">
                <el-upload
                    :before-upload="beforeIconUpload"
                    :disabled="customerServiceIconUploading"
                    :http-request="(opt) => doAssetUpload(opt, 'customerServiceIcon')"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="icon-uploader"
                >
                  <img
                      v-if="customerServiceIconPreviewUrl"
                      :src="customerServiceIconPreviewUrl"
                      alt="customer service icon"
                      class="icon-preview"
                  />
                  <div v-else class="asset-uploader-placeholder icon-placeholder">
                    <el-icon class="asset-uploader-icon">
                      <Plus/>
                    </el-icon>
                    <span>{{ t('pages.vipCfgList.clickUploadIcon') }}</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="customerServiceIconPreviewUrl || currentRow.customerServiceIcon"
                    link
                    type="danger"
                    @click="clearCustomerServiceIcon"
                >
                  {{ t('pages.vipCfgList.removeIcon') }}
                </el-button>
              </div>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.privilegeTitleEn')" prop="customerServiceTitleEn">
              <el-input v-model="currentRow.customerServiceTitleEn" maxlength="128" placeholder="Customer service privilege title in English" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.privilegeTitleEs')" prop="customerServiceTitleEs">
              <el-input v-model="currentRow.customerServiceTitleEs" maxlength="128" placeholder="Título de atención al cliente en español" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.privilegeTitlePt')" prop="customerServiceTitlePt">
              <el-input v-model="currentRow.customerServiceTitlePt" maxlength="128" placeholder="Título de atendimento ao cliente em português" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.privilegeTitleHi')" prop="customerServiceTitleHi">
              <el-input v-model="currentRow.customerServiceTitleHi" maxlength="128" placeholder="ग्राहक सेवा विशेषाधिकार शीर्षक हिंदी में" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.privilegeTitleId')" prop="customerServiceTitleId">
              <el-input v-model="currentRow.customerServiceTitleId" maxlength="128" placeholder="Judul hak layanan pelanggan dalam Bahasa Indonesia" show-word-limit/>
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.privilegeDescEn')" prop="customerServiceDescEn">
              <el-input
                  v-model="currentRow.customerServiceDescEn"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Customer service priority description in English"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.privilegeDescEs')" prop="customerServiceDescEs">
              <el-input
                  v-model="currentRow.customerServiceDescEs"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descripción de prioridad de atención al cliente en español"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.privilegeDescPt')" prop="customerServiceDescPt">
              <el-input
                  v-model="currentRow.customerServiceDescPt"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descrição de prioridade de atendimento ao cliente em português"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.privilegeDescHi')" prop="customerServiceDescHi">
              <el-input
                  v-model="currentRow.customerServiceDescHi"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="ग्राहक सेवा प्राथमिकता विवरण हिंदी में"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item :label="t('pages.vipCfgList.privilegeDescId')" prop="customerServiceDescId">
              <el-input
                  v-model="currentRow.customerServiceDescId"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Deskripsi prioritas layanan pelanggan dalam Bahasa Indonesia"
                  show-word-limit
                  type="textarea"
              />
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
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {Plus} from '@element-plus/icons-vue'
import {dataSyncApi, uploadApi, vipCfgApi} from '@/api'
import type {VipCfg} from '@/types/api.ts'
import {hasButtonPermission} from '@/utils/permission'
import {getExt, isVideoUrl, resolveMediaPreviewType} from '@/utils/media-preview'
import {formatAmount, NUMBER_INPUT_DECIMALS, truncateNumber} from '@/utils/number-format'

interface SearchForm {
  levelName: string
}

interface VipCfgForm {
  id: string
  level: number
  levelName: string
  levelIcon: string
  animationSwitch: number
  commentEffectSwitch: number
  customerServiceSwitch: number
  upgradeRechargeLimit: number
  animation: string
  animationIcon: string
  animationTitleEn: string
  animationTitleEs: string
  animationTitlePt: string
  animationTitleHi: string
  animationTitleId: string
  animationDescEn: string
  animationDescEs: string
  animationDescPt: string
  animationDescHi: string
  animationDescId: string
  commentEffect: string
  commentEffectIcon: string
  commentEffectTitleEn: string
  commentEffectTitleEs: string
  commentEffectTitlePt: string
  commentEffectTitleHi: string
  commentEffectTitleId: string
  commentEffectDescEn: string
  commentEffectDescEs: string
  commentEffectDescPt: string
  commentEffectDescHi: string
  commentEffectDescId: string
  customerServiceIcon: string
  customerServiceTitleEn: string
  customerServiceTitleEs: string
  customerServiceTitlePt: string
  customerServiceTitleHi: string
  customerServiceTitleId: string
  customerServiceDescEn: string
  customerServiceDescEs: string
  customerServiceDescPt: string
  customerServiceDescHi: string
  customerServiceDescId: string
}

const {t} = useI18n()

const loading = ref(false)
const syncing = ref(false)
const selectedRows = ref<VipCfg[]>([])
const activeTab = ref('basic')
const tableData = ref<VipCfg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  levelName: ''
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): VipCfgForm => ({
  id: '',
  level: 1,
  levelName: '',
  levelIcon: '',
  animationSwitch: 1,
  commentEffectSwitch: 1,
  customerServiceSwitch: 1,
  upgradeRechargeLimit: 0,
  animation: '',
  animationIcon: '',
  animationTitleEn: '',
  animationTitleEs: '',
  animationTitlePt: '',
  animationTitleHi: '',
  animationTitleId: '',
  animationDescEn: '',
  animationDescEs: '',
  animationDescPt: '',
  animationDescHi: '',
  animationDescId: '',
  commentEffect: '',
  commentEffectIcon: '',
  commentEffectTitleEn: '',
  commentEffectTitleEs: '',
  commentEffectTitlePt: '',
  commentEffectTitleHi: '',
  commentEffectTitleId: '',
  commentEffectDescEn: '',
  commentEffectDescEs: '',
  commentEffectDescPt: '',
  commentEffectDescHi: '',
  commentEffectDescId: '',
  customerServiceIcon: '',
  customerServiceTitleEn: '',
  customerServiceTitleEs: '',
  customerServiceTitlePt: '',
  customerServiceTitleHi: '',
  customerServiceTitleId: '',
  customerServiceDescEn: '',
  customerServiceDescEs: '',
  customerServiceDescPt: '',
  customerServiceDescHi: '',
  customerServiceDescId: ''
})
const currentRow = reactive(defaultForm())
const formRef = ref<FormInstance>()

const animationUploading = ref(false)
const animationIconUploading = ref(false)
const levelIconUploading = ref(false)
const commentEffectUploading = ref(false)
const commentEffectIconUploading = ref(false)
const customerServiceIconUploading = ref(false)
const animationPreviewUrl = ref('')
const animationIconPreviewUrl = ref('')
const levelIconPreviewUrl = ref('')
const commentEffectPreviewUrl = ref('')
const commentEffectIconPreviewUrl = ref('')
const customerServiceIconPreviewUrl = ref('')
let animationObjectPreviewUrl = ''
let animationIconObjectPreviewUrl = ''
let levelIconObjectPreviewUrl = ''
let commentEffectObjectPreviewUrl = ''
let commentEffectIconObjectPreviewUrl = ''
let customerServiceIconObjectPreviewUrl = ''

const revokeAnimationObjectPreview = () => {
  if (animationObjectPreviewUrl) {
    URL.revokeObjectURL(animationObjectPreviewUrl)
    animationObjectPreviewUrl = ''
  }
}

const revokeAnimationIconObjectPreview = () => {
  if (animationIconObjectPreviewUrl) {
    URL.revokeObjectURL(animationIconObjectPreviewUrl)
    animationIconObjectPreviewUrl = ''
  }
}

const revokeCustomerServiceIconObjectPreview = () => {
  if (customerServiceIconObjectPreviewUrl) {
    URL.revokeObjectURL(customerServiceIconObjectPreviewUrl)
    customerServiceIconObjectPreviewUrl = ''
  }
}

const revokeCommentEffectObjectPreview = () => {
  if (commentEffectObjectPreviewUrl) {
    URL.revokeObjectURL(commentEffectObjectPreviewUrl)
    commentEffectObjectPreviewUrl = ''
  }
}

const revokeCommentEffectIconObjectPreview = () => {
  if (commentEffectIconObjectPreviewUrl) {
    URL.revokeObjectURL(commentEffectIconObjectPreviewUrl)
    commentEffectIconObjectPreviewUrl = ''
  }
}

const revokeLevelIconObjectPreview = () => {
  if (levelIconObjectPreviewUrl) {
    URL.revokeObjectURL(levelIconObjectPreviewUrl)
    levelIconObjectPreviewUrl = ''
  }
}

const resetAnimationPreview = () => {
  revokeAnimationObjectPreview()
  animationPreviewUrl.value = ''
}

const resetAnimationIconPreview = () => {
  revokeAnimationIconObjectPreview()
  animationIconPreviewUrl.value = ''
}

const resetCustomerServiceIconPreview = () => {
  revokeCustomerServiceIconObjectPreview()
  customerServiceIconPreviewUrl.value = ''
}

const resetCommentEffectPreview = () => {
  revokeCommentEffectObjectPreview()
  commentEffectPreviewUrl.value = ''
}

const resetCommentEffectIconPreview = () => {
  revokeCommentEffectIconObjectPreview()
  commentEffectIconPreviewUrl.value = ''
}

const resetLevelIconPreview = () => {
  revokeLevelIconObjectPreview()
  levelIconPreviewUrl.value = ''
}

const setAnimationPreview = (url: string, fromObject = false) => {
  revokeAnimationObjectPreview()
  animationPreviewUrl.value = url
  if (fromObject && url) {
    animationObjectPreviewUrl = url
  }
}

const setAnimationIconPreview = (url: string, fromObject = false) => {
  revokeAnimationIconObjectPreview()
  animationIconPreviewUrl.value = url
  if (fromObject && url) {
    animationIconObjectPreviewUrl = url
  }
}

const setCustomerServiceIconPreview = (url: string, fromObject = false) => {
  revokeCustomerServiceIconObjectPreview()
  customerServiceIconPreviewUrl.value = url
  if (fromObject && url) {
    customerServiceIconObjectPreviewUrl = url
  }
}

const setCommentEffectPreview = (url: string, fromObject = false) => {
  revokeCommentEffectObjectPreview()
  commentEffectPreviewUrl.value = url
  if (fromObject && url) {
    commentEffectObjectPreviewUrl = url
  }
}

const setCommentEffectIconPreview = (url: string, fromObject = false) => {
  revokeCommentEffectIconObjectPreview()
  commentEffectIconPreviewUrl.value = url
  if (fromObject && url) {
    commentEffectIconObjectPreviewUrl = url
  }
}

const setLevelIconPreview = (url: string, fromObject = false) => {
  revokeLevelIconObjectPreview()
  levelIconPreviewUrl.value = url
  if (fromObject && url) {
    levelIconObjectPreviewUrl = url
  }
}

const clearAnimation = () => {
  currentRow.animation = ''
  resetAnimationPreview()
}

const clearAnimationIcon = () => {
  currentRow.animationIcon = ''
  resetAnimationIconPreview()
}

const clearCustomerServiceIcon = () => {
  currentRow.customerServiceIcon = ''
  resetCustomerServiceIconPreview()
}

const clearCommentEffect = () => {
  currentRow.commentEffect = ''
  resetCommentEffectPreview()
}

const clearCommentEffectIcon = () => {
  currentRow.commentEffectIcon = ''
  resetCommentEffectIconPreview()
}

const clearLevelIcon = () => {
  currentRow.levelIcon = ''
  resetLevelIconPreview()
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    resetAnimationPreview()
    resetAnimationIconPreview()
    resetLevelIconPreview()
    resetCommentEffectPreview()
    resetCommentEffectIconPreview()
    resetCustomerServiceIconPreview()
    activeTab.value = 'basic'
  }
})

const beforeAnimationUpload = (file: File): boolean => {
  const ext = getExt(file.name)
  if (ext !== '.mp4') {
    ElMessage.error(t('pages.vipCfgList.mp4Only'))
    return false
  }
  return true
}

const beforeIconUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.vipCfgList.iconImageOnly'))
    return false
  }
  return true
}

const doAssetUpload = async (
    options: UploadRequestOptions,
    field: 'animation' | 'animationIcon' | 'levelIcon' | 'commentEffect' | 'commentEffectIcon' | 'customerServiceIcon'
) => {
  const file = options.file as File
  const uploadingMap = {
    animation: animationUploading,
    animationIcon: animationIconUploading,
    levelIcon: levelIconUploading,
    commentEffect: commentEffectUploading,
    commentEffectIcon: commentEffectIconUploading,
    customerServiceIcon: customerServiceIconUploading
  }
  const uploading = uploadingMap[field]
  uploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    currentRow[field] = res.fileName
    if (field === 'animation') {
      setAnimationPreview(URL.createObjectURL(file), true)
    } else if (field === 'animationIcon') {
      setAnimationIconPreview(URL.createObjectURL(file), true)
    } else if (field === 'levelIcon') {
      setLevelIconPreview(URL.createObjectURL(file), true)
    } else if (field === 'commentEffect') {
      setCommentEffectPreview(URL.createObjectURL(file), true)
    } else if (field === 'commentEffectIcon') {
      setCommentEffectIconPreview(URL.createObjectURL(file), true)
    } else if (field === 'customerServiceIcon') {
      setCustomerServiceIconPreview(URL.createObjectURL(file), true)
    }
    ElMessage.success(t('pages.vipCfgList.uploadSuccess'))
  } catch (error) {
    console.error('upload failed:', error)
    ElMessage.error(t('pages.vipCfgList.uploadFailed'))
  } finally {
    uploading.value = false
  }
}

const formRules = computed<FormRules>(() => ({
  level: [{required: true, message: t('pages.vipCfgList.levelRequired'), trigger: 'change'}],
  levelName: [
    {required: true, message: t('pages.vipCfgList.levelNameRequired'), trigger: 'blur'},
    {min: 1, max: 64, message: t('pages.vipCfgList.levelNameLength'), trigger: 'blur'}
  ],
  animationSwitch: [{required: true, message: t('pages.vipCfgList.animationSwitchRequired'), trigger: 'change'}],
  commentEffectSwitch: [{required: true, message: t('pages.vipCfgList.commentEffectSwitchRequired'), trigger: 'change'}],
  customerServiceSwitch: [{required: true, message: t('pages.vipCfgList.customerServiceSwitchRequired'), trigger: 'change'}]
}))

const fetchList = async () => {
  loading.value = true
  try {
    const response = await vipCfgApi.getVipCfgList({
      levelName: searchForm.levelName,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('fetch vip cfg list failed:', error)
    ElMessage.error(t('pages.vipCfgList.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchList()
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchList()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchList()
}

const handleAdd = () => {
  dialogTitle.value = t('pages.vipCfgList.addVipLevel')
  Object.assign(currentRow, defaultForm())
  activeTab.value = 'basic'
  resetAnimationPreview()
  resetAnimationIconPreview()
  resetLevelIconPreview()
  resetCommentEffectPreview()
  resetCommentEffectIconPreview()
  resetCustomerServiceIconPreview()
  dialogVisible.value = true
}

const handleSelectionChange = (rows: VipCfg[]) => {
  selectedRows.value = rows
}

const handleSyncData = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning(t('pages.vipCfgList.selectSyncFirst'))
    return
  }
  const ids = selectedRows.value.map((row) => Number(row.id)).filter((id) => id > 0)
  if (ids.length === 0) {
    ElMessage.warning(t('pages.vipCfgList.invalidSelection'))
    return
  }
  try {
    await ElMessageBox.confirm(
        t('pages.vipCfgList.syncConfirm', {count: ids.length}),
        t('common.syncData'),
        {
          confirmButtonText: t('common.confirmSync'),
          cancelButtonText: t('common.cancel'),
          type: 'warning',
        }
    )
    syncing.value = true
    const response = await dataSyncApi.syncVipCfg({ids})
    if (response?.success) {
      ElMessage.success(response.message || t('pages.vipCfgList.syncSuccessDetail', {
        rows: response.rowCount,
        files: response.fileCount,
      }))
    } else {
      ElMessage.error(t('pages.vipCfgList.syncFailed'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('sync failed:', error)
      ElMessage.error(t('pages.vipCfgList.syncFailedCheckConfig'))
    }
  } finally {
    syncing.value = false
  }
}

const handleEdit = (row: VipCfg) => {
  dialogTitle.value = t('pages.vipCfgList.editVipLevel')
  activeTab.value = 'basic'
  const animationName = row.animationName || ''
  const animationIconName = row.animationIconName || ''
  const levelIconName = row.levelIconName || ''
  const commentEffectName = row.commentEffectName || ''
  const commentEffectIconName = row.commentEffectIconName || ''
  const customerServiceIconName = row.customerServiceIconName || ''
  Object.assign(currentRow, {
    id: row.id,
    level: Number(row.level) || 1,
    levelName: row.levelName,
    levelIcon: levelIconName,
    animationSwitch: Number(row.animationSwitch) || 0,
    commentEffectSwitch: Number(row.commentEffectSwitch) || 0,
    customerServiceSwitch: Number(row.customerServiceSwitch) || 0,
    upgradeRechargeLimit: truncateNumber(row.upgradeRechargeLimit),
    animation: animationName,
    animationIcon: animationIconName,
    animationTitleEn: row.animationTitleEn || '',
    animationTitleEs: row.animationTitleEs || '',
    animationTitlePt: row.animationTitlePt || '',
    animationTitleHi: row.animationTitleHi || '',
    animationTitleId: row.animationTitleId || '',
    animationDescEn: row.animationDescEn || '',
    animationDescEs: row.animationDescEs || '',
    animationDescPt: row.animationDescPt || '',
    animationDescHi: row.animationDescHi || '',
    animationDescId: row.animationDescId || '',
    commentEffect: commentEffectName,
    commentEffectIcon: commentEffectIconName,
    commentEffectTitleEn: row.commentEffectTitleEn || '',
    commentEffectTitleEs: row.commentEffectTitleEs || '',
    commentEffectTitlePt: row.commentEffectTitlePt || '',
    commentEffectTitleHi: row.commentEffectTitleHi || '',
    commentEffectTitleId: row.commentEffectTitleId || '',
    commentEffectDescEn: row.commentEffectDescEn || '',
    commentEffectDescEs: row.commentEffectDescEs || '',
    commentEffectDescPt: row.commentEffectDescPt || '',
    commentEffectDescHi: row.commentEffectDescHi || '',
    commentEffectDescId: row.commentEffectDescId || '',
    customerServiceIcon: customerServiceIconName,
    customerServiceTitleEn: row.customerServiceTitleEn || '',
    customerServiceTitleEs: row.customerServiceTitleEs || '',
    customerServiceTitlePt: row.customerServiceTitlePt || '',
    customerServiceTitleHi: row.customerServiceTitleHi || '',
    customerServiceTitleId: row.customerServiceTitleId || '',
    customerServiceDescEn: row.customerServiceDescEn || '',
    customerServiceDescEs: row.customerServiceDescEs || '',
    customerServiceDescPt: row.customerServiceDescPt || '',
    customerServiceDescHi: row.customerServiceDescHi || '',
    customerServiceDescId: row.customerServiceDescId || ''
  })
  if (animationName && resolveMediaPreviewType(row.animation || '', animationName) === 'video') {
    setAnimationPreview(row.animation || '')
  } else {
    resetAnimationPreview()
  }
  if (animationIconName) {
    setAnimationIconPreview(row.animationIcon || '')
  } else {
    resetAnimationIconPreview()
  }
  if (levelIconName) {
    setLevelIconPreview(row.levelIcon || '')
  } else {
    resetLevelIconPreview()
  }
  if (commentEffectName && resolveMediaPreviewType(row.commentEffect || '', commentEffectName) === 'video') {
    setCommentEffectPreview(row.commentEffect || '')
  } else {
    resetCommentEffectPreview()
  }
  if (commentEffectIconName) {
    setCommentEffectIconPreview(row.commentEffectIcon || '')
  } else {
    resetCommentEffectIconPreview()
  }
  if (customerServiceIconName) {
    setCustomerServiceIconPreview(row.customerServiceIcon || '')
  } else {
    resetCustomerServiceIconPreview()
  }
  dialogVisible.value = true
}

const handleDelete = async (row: VipCfg) => {
  try {
    await ElMessageBox.confirm(t('pages.vipCfgList.deleteConfirm', {name: row.levelName}), t('common.confirmDelete'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning'
    })
    await vipCfgApi.deleteVipCfg(row.id)
    ElMessage.success(t('common.deleteSuccess'))
    fetchList()
  } catch (error) {
    console.error('delete failed:', error)
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const payload = {
        level: currentRow.level,
        levelName: currentRow.levelName,
        levelIcon: currentRow.levelIcon,
        animationSwitch: currentRow.animationSwitch,
        commentEffectSwitch: currentRow.commentEffectSwitch,
        customerServiceSwitch: currentRow.customerServiceSwitch,
        upgradeRechargeLimit: currentRow.upgradeRechargeLimit,
        animation: currentRow.animation,
        animationIcon: currentRow.animationIcon,
        animationTitleEn: currentRow.animationTitleEn,
        animationTitleEs: currentRow.animationTitleEs,
        animationTitlePt: currentRow.animationTitlePt,
        animationTitleHi: currentRow.animationTitleHi,
        animationTitleId: currentRow.animationTitleId,
        animationDescEn: currentRow.animationDescEn,
        animationDescEs: currentRow.animationDescEs,
        animationDescPt: currentRow.animationDescPt,
        animationDescHi: currentRow.animationDescHi,
        animationDescId: currentRow.animationDescId,
        commentEffect: currentRow.commentEffect,
        commentEffectIcon: currentRow.commentEffectIcon,
        commentEffectTitleEn: currentRow.commentEffectTitleEn,
        commentEffectTitleEs: currentRow.commentEffectTitleEs,
        commentEffectTitlePt: currentRow.commentEffectTitlePt,
        commentEffectTitleHi: currentRow.commentEffectTitleHi,
        commentEffectTitleId: currentRow.commentEffectTitleId,
        commentEffectDescEn: currentRow.commentEffectDescEn,
        commentEffectDescEs: currentRow.commentEffectDescEs,
        commentEffectDescPt: currentRow.commentEffectDescPt,
        commentEffectDescHi: currentRow.commentEffectDescHi,
        commentEffectDescId: currentRow.commentEffectDescId,
        customerServiceIcon: currentRow.customerServiceIcon,
        customerServiceTitleEn: currentRow.customerServiceTitleEn,
        customerServiceTitleEs: currentRow.customerServiceTitleEs,
        customerServiceTitlePt: currentRow.customerServiceTitlePt,
        customerServiceTitleHi: currentRow.customerServiceTitleHi,
        customerServiceTitleId: currentRow.customerServiceTitleId,
        customerServiceDescEn: currentRow.customerServiceDescEn,
        customerServiceDescEs: currentRow.customerServiceDescEs,
        customerServiceDescPt: currentRow.customerServiceDescPt,
        customerServiceDescHi: currentRow.customerServiceDescHi,
        customerServiceDescId: currentRow.customerServiceDescId
      }
      if (currentRow.id) {
        await vipCfgApi.updateVipCfg({id: currentRow.id, ...payload})
      } else {
        await vipCfgApi.createVipCfg(payload)
      }
      ElMessage.success(currentRow.id ? t('common.updateSuccess') : t('common.createSuccess'))
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('save failed:', error)
      ElMessage.error(t('pages.vipCfgList.saveFailed'))
    }
  })
}

const resetSearch = () => {
  searchForm.levelName = ''
  currentPage.value = 1
  fetchList()
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
  font-size: 16px;
  font-weight: bold;
}

.table-header {
  margin-bottom: 20px;
}

.selection-tip {
  margin-bottom: 12px;
  font-size: 13px;
  color: var(--el-color-primary);
}

.search-form {
  margin-bottom: 20px;
}

.search-form .el-form-item {
  margin-bottom: 12px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}

.form-tip-block {
  margin-bottom: 12px;
}

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.asset-upload-wrap {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
}

.animation-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  transition: border-color 0.2s;
}

.animation-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.icon-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 8px;
  cursor: pointer;
  overflow: hidden;
  transition: border-color 0.2s;
}

.icon-uploader :deep(.el-upload:hover) {
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

.animation-placeholder {
  width: 240px;
  height: 120px;
}

.icon-placeholder {
  width: 120px;
  height: 120px;
}

.asset-uploader-icon {
  font-size: 28px;
}

.animation-preview {
  width: 240px;
  height: 120px;
  display: block;
  object-fit: cover;
}

.icon-preview {
  width: 120px;
  height: 120px;
  display: block;
  object-fit: contain;
}

.table-media-preview {
  width: 160px;
  max-height: 90px;
  display: block;
  background: #000;
  border-radius: 4px;
}

.table-icon-preview {
  width: 48px;
  height: 48px;
  display: block;
  border-radius: 4px;
}

.media-url-text {
  word-break: break-all;
  line-height: 1.4;
}
</style>
