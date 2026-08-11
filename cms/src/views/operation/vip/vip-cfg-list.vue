<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>VIP配置管理</span>
        </div>
      </template>
      <div class="content">
        <div class="table-header">
          <el-button type="primary" @click="handleAdd">新增VIP等级</el-button>
          <el-button
              v-if="hasButtonPermission('VipCfgManagement', 'sync')"
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
          <el-form-item label="等级名称">
            <el-input v-model="searchForm.levelName" clearable placeholder="等级名称(模糊匹配)"/>
          </el-form-item>
          <el-form-item label="提现开关">
            <el-select v-model="searchForm.withdrawSwitchFilter" placeholder="全部" style="width: 140px">
              <el-option :value="0" label="全部"/>
              <el-option :value="2" label="只看开启"/>
              <el-option :value="1" label="只看关闭"/>
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
          <el-table-column label="ID" prop="id" width="100"/>
          <el-table-column label="等级" prop="level" width="80"/>
          <el-table-column label="等级名称" prop="levelName" min-width="120"/>
          <el-table-column label="等级图标" width="100">
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
          <el-table-column label="进场特效" min-width="200">
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
          <el-table-column label="进场特效图标" width="110">
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
          <el-table-column label="提现图标" width="100">
            <template #default="{ row }">
              <el-image
                  v-if="row.withdrawIcon"
                  :preview-src-list="[row.withdrawIcon]"
                  :src="row.withdrawIcon"
                  class="table-icon-preview"
                  fit="cover"
                  preview-teleported
              />
              <span v-else>-</span>
            </template>
          </el-table-column>
          <el-table-column label="提现开关" width="100">
            <template #default="{ row }">
              <el-tag :type="row.withdrawSwitch === 1 ? 'success' : 'info'">
                {{ row.withdrawSwitch === 1 ? '开' : '关' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="进场特效开关" width="120">
            <template #default="{ row }">
              <el-tag :type="row.animationSwitch === 1 ? 'success' : 'info'">
                {{ row.animationSwitch === 1 ? '开' : '关' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="公屏评论特效" min-width="200">
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
          <el-table-column label="公屏评论特效图标" width="130">
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
          <el-table-column label="公屏评论特效开关" width="140">
            <template #default="{ row }">
              <el-tag :type="row.commentEffectSwitch === 1 ? 'success' : 'info'">
                {{ row.commentEffectSwitch === 1 ? '开' : '关' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="客服优先图标" width="120">
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
          <el-table-column label="客服优先开关" width="120">
            <template #default="{ row }">
              <el-tag :type="row.customerServiceSwitch === 1 ? 'success' : 'info'">
                {{ row.customerServiceSwitch === 1 ? '开' : '关' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="升级充值上限" width="130">
            <template #default="{ row }">
              {{ formatAmount(row.upgradeRechargeLimit) }}
            </template>
          </el-table-column>
          <el-table-column label="最低提现金额" width="130">
            <template #default="{ row }">
              {{ formatAmount(row.minWithdrawAmount) }}
            </template>
          </el-table-column>
          <el-table-column label="最高提现金额" width="130">
            <template #default="{ row }">
              {{ formatAmount(row.maxWithdrawAmount) }}
            </template>
          </el-table-column>
          <el-table-column label="手续费" width="100">
            <template #default="{ row }">
              {{ formatAmount(row.fee) }}
            </template>
          </el-table-column>
          <el-table-column label="创建时间" prop="createdAt" width="160"/>
          <el-table-column label="更新时间" prop="updatedAt" width="160"/>
          <el-table-column fixed="right" label="操作" width="180">
            <template #default="{ row }">
              <el-button size="small" @click="handleEdit(row)">编辑</el-button>
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="720px">
      <el-form ref="formRef" :model="currentRow" :rules="formRules" label-width="130px">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="基础信息" name="basic">
            <el-form-item label="等级" prop="level">
              <el-input-number v-model="currentRow.level" :min="1" controls-position="right"/>
            </el-form-item>
            <el-form-item label="等级名称" prop="levelName">
              <el-input v-model="currentRow.levelName" placeholder="请输入等级名称"/>
            </el-form-item>
            <el-form-item label="等级图标" prop="levelIcon">
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
                    <span>点击上传图标</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="levelIconPreviewUrl || currentRow.levelIcon"
                    link
                    type="danger"
                    @click="clearLevelIcon"
                >
                  移除图标
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="升级充值上限" prop="upgradeRechargeLimit">
              <el-input-number
                  v-model="currentRow.upgradeRechargeLimit"
                  :min="0"
                  :precision="NUMBER_INPUT_DECIMALS"
                  :step="0.0001"
                  controls-position="right"
              />
              <div class="form-tip">保留4位小数，例如 100.0000</div>
            </el-form-item>
            <el-form-item label="最低提现金额" prop="minWithdrawAmount">
              <el-input-number
                  v-model="currentRow.minWithdrawAmount"
                  :min="0"
                  :precision="NUMBER_INPUT_DECIMALS"
                  :step="0.0001"
                  controls-position="right"
              />
              <div class="form-tip">保留4位小数</div>
            </el-form-item>
            <el-form-item label="最高提现金额" prop="maxWithdrawAmount">
              <el-input-number
                  v-model="currentRow.maxWithdrawAmount"
                  :min="0"
                  :precision="NUMBER_INPUT_DECIMALS"
                  :step="0.0001"
                  controls-position="right"
              />
              <div class="form-tip">保留4位小数，0 表示不限制</div>
            </el-form-item>
            <el-form-item label="手续费" prop="fee">
              <el-input-number
                  v-model="currentRow.fee"
                  :min="0"
                  :precision="NUMBER_INPUT_DECIMALS"
                  :step="0.0001"
                  controls-position="right"
              />
              <div class="form-tip">保留4位小数</div>
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane label="房间进场特效" name="entryEffect">
            <el-form-item label="进场特效开关" prop="animationSwitch">
              <el-radio-group v-model="currentRow.animationSwitch">
                <el-radio :label="1">开</el-radio>
                <el-radio :label="0">关</el-radio>
              </el-radio-group>
              <div class="form-tip">仅控制 App 端是否展示该等级进场特效，后端不参与业务判断</div>
            </el-form-item>
            <el-form-item label="特效动画" prop="animation">
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
                    <span>点击上传 MP4 动画</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="animationPreviewUrl || currentRow.animation"
                    link
                    type="danger"
                    @click="clearAnimation"
                >
                  移除动画
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="特效图标" prop="animationIcon">
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
                    <span>点击上传图标</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="animationIconPreviewUrl || currentRow.animationIcon"
                    link
                    type="danger"
                    @click="clearAnimationIcon"
                >
                  移除图标
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="特效说明(英文)" prop="animationDescEn">
              <el-input
                  v-model="currentRow.animationDescEn"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Entry effect description in English"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="特效说明(西班牙语)" prop="animationDescEs">
              <el-input
                  v-model="currentRow.animationDescEs"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descripción del efecto de entrada en español"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="特效说明(葡萄牙语)" prop="animationDescPt">
              <el-input
                  v-model="currentRow.animationDescPt"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descrição do efeito de entrada em português"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="特效说明(印地语)" prop="animationDescHi">
              <el-input
                  v-model="currentRow.animationDescHi"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="प्रवेश प्रभाव विवरण हिंदी में"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane label="直播间公屏评论特效" name="commentEffect">
            <el-form-item label="公屏评论特效开关" prop="commentEffectSwitch">
              <el-radio-group v-model="currentRow.commentEffectSwitch">
                <el-radio :label="1">开</el-radio>
                <el-radio :label="0">关</el-radio>
              </el-radio-group>
              <div class="form-tip">仅控制 App 端是否展示该等级公屏评论特效，后端不参与业务判断</div>
            </el-form-item>
            <el-form-item label="特效动画" prop="commentEffect">
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
                    <span>点击上传 MP4 动画</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="commentEffectPreviewUrl || currentRow.commentEffect"
                    link
                    type="danger"
                    @click="clearCommentEffect"
                >
                  移除动画
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="特效图标" prop="commentEffectIcon">
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
                    <span>点击上传图标</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="commentEffectIconPreviewUrl || currentRow.commentEffectIcon"
                    link
                    type="danger"
                    @click="clearCommentEffectIcon"
                >
                  移除图标
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="特效说明(英文)" prop="commentEffectDescEn">
              <el-input
                  v-model="currentRow.commentEffectDescEn"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Public screen comment effect description in English"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="特效说明(西班牙语)" prop="commentEffectDescEs">
              <el-input
                  v-model="currentRow.commentEffectDescEs"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descripción del efecto de comentario en pantalla pública en español"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="特效说明(葡萄牙语)" prop="commentEffectDescPt">
              <el-input
                  v-model="currentRow.commentEffectDescPt"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descrição do efeito de comentário na tela pública em português"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="特效说明(印地语)" prop="commentEffectDescHi">
              <el-input
                  v-model="currentRow.commentEffectDescHi"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="सार्वजनिक स्क्रीन टिप्पणी प्रभाव विवरण हिंदी में"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane label="提现文本" name="withdrawText">
            <el-form-item label="提现开关" prop="withdrawSwitch">
              <el-radio-group v-model="currentRow.withdrawSwitch">
                <el-radio :label="1">开</el-radio>
                <el-radio :label="0">关</el-radio>
              </el-radio-group>
              <div class="form-tip">仅控制 App 端是否展示该等级提现特权，后端不参与业务判断</div>
            </el-form-item>
            <el-form-item label="提现图标" prop="withdrawIcon">
              <div class="asset-upload-wrap">
                <el-upload
                    :before-upload="beforeIconUpload"
                    :disabled="withdrawIconUploading"
                    :http-request="(opt) => doAssetUpload(opt, 'withdrawIcon')"
                    :show-file-list="false"
                    accept="image/*"
                    action="#"
                    class="icon-uploader"
                >
                  <img
                      v-if="withdrawIconPreviewUrl"
                      :src="withdrawIconPreviewUrl"
                      alt="withdraw icon"
                      class="icon-preview"
                  />
                  <div v-else class="asset-uploader-placeholder icon-placeholder">
                    <el-icon class="asset-uploader-icon">
                      <Plus/>
                    </el-icon>
                    <span>点击上传图标</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="withdrawIconPreviewUrl || currentRow.withdrawIcon"
                    link
                    type="danger"
                    @click="clearWithdrawIcon"
                >
                  移除图标
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="提现须知(英文)" prop="withdrawNoticeEn">
              <el-input
                  v-model="currentRow.withdrawNoticeEn"
                  :autosize="{ minRows: 4, maxRows: 8 }"
                  maxlength="2000"
                  placeholder="Withdraw notice in English"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="提现须知(西班牙语)" prop="withdrawNoticeEs">
              <el-input
                  v-model="currentRow.withdrawNoticeEs"
                  :autosize="{ minRows: 4, maxRows: 8 }"
                  maxlength="2000"
                  placeholder="Aviso de retiro en español"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="提现须知(葡萄牙语)" prop="withdrawNoticePt">
              <el-input
                  v-model="currentRow.withdrawNoticePt"
                  :autosize="{ minRows: 4, maxRows: 8 }"
                  maxlength="2000"
                  placeholder="Aviso de saque em português"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="提现须知(印地语)" prop="withdrawNoticeHi">
              <el-input
                  v-model="currentRow.withdrawNoticeHi"
                  :autosize="{ minRows: 4, maxRows: 8 }"
                  maxlength="2000"
                  placeholder="हिंदी में निकासी सूचना"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
          </el-tab-pane>

          <el-tab-pane label="客服优先" name="customerService">
            <el-form-item label="客服优先开关" prop="customerServiceSwitch">
              <el-radio-group v-model="currentRow.customerServiceSwitch">
                <el-radio :label="1">开</el-radio>
                <el-radio :label="0">关</el-radio>
              </el-radio-group>
              <div class="form-tip">仅控制 App 端是否展示该等级客服优先特权，后端不参与业务判断</div>
            </el-form-item>
            <el-form-item label="特权图标" prop="customerServiceIcon">
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
                    <span>点击上传图标</span>
                  </div>
                </el-upload>
                <el-button
                    v-if="customerServiceIconPreviewUrl || currentRow.customerServiceIcon"
                    link
                    type="danger"
                    @click="clearCustomerServiceIcon"
                >
                  移除图标
                </el-button>
              </div>
            </el-form-item>
            <el-form-item label="特权说明(英文)" prop="customerServiceDescEn">
              <el-input
                  v-model="currentRow.customerServiceDescEn"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Customer service priority description in English"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="特权说明(西班牙语)" prop="customerServiceDescEs">
              <el-input
                  v-model="currentRow.customerServiceDescEs"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descripción de prioridad de atención al cliente en español"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="特权说明(葡萄牙语)" prop="customerServiceDescPt">
              <el-input
                  v-model="currentRow.customerServiceDescPt"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="Descrição de prioridade de atendimento ao cliente em português"
                  show-word-limit
                  type="textarea"
              />
            </el-form-item>
            <el-form-item label="特权说明(印地语)" prop="customerServiceDescHi">
              <el-input
                  v-model="currentRow.customerServiceDescHi"
                  :autosize="{ minRows: 3, maxRows: 6 }"
                  maxlength="2000"
                  placeholder="ग्राहक सेवा प्राथमिकता विवरण हिंदी में"
                  show-word-limit
                  type="textarea"
              />
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
import {dataSyncApi, uploadApi, vipCfgApi} from '@/api'
import type {VipCfg} from '@/types/api.ts'
import {hasButtonPermission} from '@/utils/permission'
import {getExt, isVideoUrl, resolveMediaPreviewType} from '@/utils/media-preview'
import {formatAmount, NUMBER_INPUT_DECIMALS, truncateNumber} from '@/utils/number-format'

interface SearchForm {
  levelName: string
  withdrawSwitchFilter: number
}

interface VipCfgForm {
  id: string
  level: number
  levelName: string
  levelIcon: string
  withdrawSwitch: number
  animationSwitch: number
  commentEffectSwitch: number
  customerServiceSwitch: number
  upgradeRechargeLimit: number
  minWithdrawAmount: number
  maxWithdrawAmount: number
  fee: number
  animation: string
  animationIcon: string
  animationDescEn: string
  animationDescEs: string
  animationDescPt: string
  animationDescHi: string
  commentEffect: string
  commentEffectIcon: string
  commentEffectDescEn: string
  commentEffectDescEs: string
  commentEffectDescPt: string
  commentEffectDescHi: string
  withdrawIcon: string
  withdrawNoticeEn: string
  withdrawNoticeEs: string
  withdrawNoticePt: string
  withdrawNoticeHi: string
  customerServiceIcon: string
  customerServiceDescEn: string
  customerServiceDescEs: string
  customerServiceDescPt: string
  customerServiceDescHi: string
}

const loading = ref(false)
const syncing = ref(false)
const selectedRows = ref<VipCfg[]>([])
const activeTab = ref('basic')
const tableData = ref<VipCfg[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const searchForm = reactive<SearchForm>({
  levelName: '',
  withdrawSwitchFilter: 0
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const defaultForm = (): VipCfgForm => ({
  id: '',
  level: 1,
  levelName: '',
  levelIcon: '',
  withdrawSwitch: 1,
  animationSwitch: 1,
  commentEffectSwitch: 1,
  customerServiceSwitch: 1,
  upgradeRechargeLimit: 0,
  minWithdrawAmount: 0,
  maxWithdrawAmount: 0,
  fee: 0,
  animation: '',
  animationIcon: '',
  animationDescEn: '',
  animationDescEs: '',
  animationDescPt: '',
  animationDescHi: '',
  commentEffect: '',
  commentEffectIcon: '',
  commentEffectDescEn: '',
  commentEffectDescEs: '',
  commentEffectDescPt: '',
  commentEffectDescHi: '',
  withdrawIcon: '',
  withdrawNoticeEn: '',
  withdrawNoticeEs: '',
  withdrawNoticePt: '',
  withdrawNoticeHi: '',
  customerServiceIcon: '',
  customerServiceDescEn: '',
  customerServiceDescEs: '',
  customerServiceDescPt: '',
  customerServiceDescHi: ''
})
const currentRow = ref<VipCfgForm>(defaultForm())
const formRef = ref<FormInstance>()

const animationUploading = ref(false)
const animationIconUploading = ref(false)
const levelIconUploading = ref(false)
const commentEffectUploading = ref(false)
const commentEffectIconUploading = ref(false)
const withdrawIconUploading = ref(false)
const customerServiceIconUploading = ref(false)
const animationPreviewUrl = ref('')
const animationIconPreviewUrl = ref('')
const levelIconPreviewUrl = ref('')
const commentEffectPreviewUrl = ref('')
const commentEffectIconPreviewUrl = ref('')
const withdrawIconPreviewUrl = ref('')
const customerServiceIconPreviewUrl = ref('')
let animationObjectPreviewUrl = ''
let animationIconObjectPreviewUrl = ''
let levelIconObjectPreviewUrl = ''
let commentEffectObjectPreviewUrl = ''
let commentEffectIconObjectPreviewUrl = ''
let withdrawIconObjectPreviewUrl = ''
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

const revokeWithdrawIconObjectPreview = () => {
  if (withdrawIconObjectPreviewUrl) {
    URL.revokeObjectURL(withdrawIconObjectPreviewUrl)
    withdrawIconObjectPreviewUrl = ''
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

const resetWithdrawIconPreview = () => {
  revokeWithdrawIconObjectPreview()
  withdrawIconPreviewUrl.value = ''
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

const setWithdrawIconPreview = (url: string, fromObject = false) => {
  revokeWithdrawIconObjectPreview()
  withdrawIconPreviewUrl.value = url
  if (fromObject && url) {
    withdrawIconObjectPreviewUrl = url
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
  currentRow.value.animation = ''
  resetAnimationPreview()
}

const clearAnimationIcon = () => {
  currentRow.value.animationIcon = ''
  resetAnimationIconPreview()
}

const clearWithdrawIcon = () => {
  currentRow.value.withdrawIcon = ''
  resetWithdrawIconPreview()
}

const clearCustomerServiceIcon = () => {
  currentRow.value.customerServiceIcon = ''
  resetCustomerServiceIconPreview()
}

const clearCommentEffect = () => {
  currentRow.value.commentEffect = ''
  resetCommentEffectPreview()
}

const clearCommentEffectIcon = () => {
  currentRow.value.commentEffectIcon = ''
  resetCommentEffectIconPreview()
}

const clearLevelIcon = () => {
  currentRow.value.levelIcon = ''
  resetLevelIconPreview()
}

watch(dialogVisible, (visible) => {
  if (!visible) {
    resetAnimationPreview()
    resetAnimationIconPreview()
    resetLevelIconPreview()
    resetCommentEffectPreview()
    resetCommentEffectIconPreview()
    resetWithdrawIconPreview()
    resetCustomerServiceIconPreview()
    activeTab.value = 'basic'
  }
})

const beforeAnimationUpload = (file: File): boolean => {
  const ext = getExt(file.name)
  if (ext !== '.mp4') {
    ElMessage.error('仅支持 MP4 格式')
    return false
  }
  return true
}

const beforeIconUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error('图标只能上传图片文件')
    return false
  }
  return true
}

const doAssetUpload = async (
    options: UploadRequestOptions,
    field: 'animation' | 'animationIcon' | 'levelIcon' | 'commentEffect' | 'commentEffectIcon' | 'withdrawIcon' | 'customerServiceIcon'
) => {
  const file = options.file as File
  const uploadingMap = {
    animation: animationUploading,
    animationIcon: animationIconUploading,
    levelIcon: levelIconUploading,
    commentEffect: commentEffectUploading,
    commentEffectIcon: commentEffectIconUploading,
    withdrawIcon: withdrawIconUploading,
    customerServiceIcon: customerServiceIconUploading
  }
  const uploading = uploadingMap[field]
  uploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    currentRow.value[field] = res.fileName
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
    } else {
      setWithdrawIconPreview(URL.createObjectURL(file), true)
    }
    ElMessage.success('上传成功')
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

const validateWithdrawRange = (_rule: unknown, _value: unknown, callback: (error?: Error) => void) => {
  const {minWithdrawAmount, maxWithdrawAmount} = currentRow.value
  if (maxWithdrawAmount > 0 && minWithdrawAmount > maxWithdrawAmount) {
    callback(new Error('最低提现金额不能大于最高提现金额'))
    return
  }
  callback()
}

const formRules: FormRules = {
  level: [{required: true, message: '请输入等级', trigger: 'change'}],
  levelName: [
    {required: true, message: '请输入等级名称', trigger: 'blur'},
    {min: 1, max: 64, message: '等级名称长度在1-64个字符', trigger: 'blur'}
  ],
  withdrawSwitch: [{required: true, message: '请选择提现开关', trigger: 'change'}],
  animationSwitch: [{required: true, message: '请选择进场特效开关', trigger: 'change'}],
  commentEffectSwitch: [{required: true, message: '请选择公屏评论特效开关', trigger: 'change'}],
  customerServiceSwitch: [{required: true, message: '请选择客服优先开关', trigger: 'change'}],
  minWithdrawAmount: [{validator: validateWithdrawRange, trigger: 'change'}],
  maxWithdrawAmount: [{validator: validateWithdrawRange, trigger: 'change'}]
}

const fetchList = async () => {
  loading.value = true
  try {
    const response = await vipCfgApi.getVipCfgList({
      levelName: searchForm.levelName,
      withdrawSwitchFilter: searchForm.withdrawSwitchFilter,
      pageIndex: currentPage.value,
      pageSize: pageSize.value
    })
    tableData.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取VIP配置列表失败:', error)
    ElMessage.error('获取VIP配置列表失败')
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
  dialogTitle.value = '新增VIP等级'
  currentRow.value = defaultForm()
  activeTab.value = 'basic'
  resetAnimationPreview()
  resetAnimationIconPreview()
  resetLevelIconPreview()
  resetCommentEffectPreview()
  resetCommentEffectIconPreview()
  resetWithdrawIconPreview()
  resetCustomerServiceIconPreview()
  dialogVisible.value = true
}

const handleSelectionChange = (rows: VipCfg[]) => {
  selectedRows.value = rows
}

const handleSyncData = async () => {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请先勾选要同步的配置')
    return
  }
  const ids = selectedRows.value.map((row) => Number(row.id)).filter((id) => id > 0)
  if (ids.length === 0) {
    ElMessage.warning('所选配置无效')
    return
  }
  try {
    await ElMessageBox.confirm(
        `将把已选 ${ids.length} 条 VIP 配置及关联图片/动画资源同步到目标环境（按 ID 覆盖或新增）。是否继续？`,
        '同步数据',
        {
          confirmButtonText: '确定同步',
          cancelButtonText: '取消',
          type: 'warning',
        }
    )
    syncing.value = true
    const response = await dataSyncApi.syncVipCfg({ids})
    if (response?.success) {
      ElMessage.success(response.message || `同步成功：${response.rowCount} 条配置，${response.fileCount} 个文件`)
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

const handleEdit = (row: VipCfg) => {
  dialogTitle.value = '编辑VIP等级'
  activeTab.value = 'basic'
  const animationName = row.animationName || ''
  const animationIconName = row.animationIconName || ''
  const levelIconName = row.levelIconName || ''
  const commentEffectName = row.commentEffectName || ''
  const commentEffectIconName = row.commentEffectIconName || ''
  const withdrawIconName = row.withdrawIconName || ''
  const customerServiceIconName = row.customerServiceIconName || ''
  currentRow.value = {
    id: row.id,
    level: Number(row.level) || 1,
    levelName: row.levelName,
    levelIcon: levelIconName,
    withdrawSwitch: Number(row.withdrawSwitch) || 0,
    animationSwitch: Number(row.animationSwitch) || 0,
    commentEffectSwitch: Number(row.commentEffectSwitch) || 0,
    customerServiceSwitch: Number(row.customerServiceSwitch) || 0,
    upgradeRechargeLimit: truncateNumber(row.upgradeRechargeLimit),
    minWithdrawAmount: truncateNumber(row.minWithdrawAmount),
    maxWithdrawAmount: truncateNumber(row.maxWithdrawAmount),
    fee: truncateNumber(row.fee),
    animation: animationName,
    animationIcon: animationIconName,
    animationDescEn: row.animationDescEn || '',
    animationDescEs: row.animationDescEs || '',
    animationDescPt: row.animationDescPt || '',
    animationDescHi: row.animationDescHi || '',
    commentEffect: commentEffectName,
    commentEffectIcon: commentEffectIconName,
    commentEffectDescEn: row.commentEffectDescEn || '',
    commentEffectDescEs: row.commentEffectDescEs || '',
    commentEffectDescPt: row.commentEffectDescPt || '',
    commentEffectDescHi: row.commentEffectDescHi || '',
    withdrawIcon: withdrawIconName,
    withdrawNoticeEn: row.withdrawNoticeEn || '',
    withdrawNoticeEs: row.withdrawNoticeEs || '',
    withdrawNoticePt: row.withdrawNoticePt || '',
    withdrawNoticeHi: row.withdrawNoticeHi || '',
    customerServiceIcon: customerServiceIconName,
    customerServiceDescEn: row.customerServiceDescEn || '',
    customerServiceDescEs: row.customerServiceDescEs || '',
    customerServiceDescPt: row.customerServiceDescPt || '',
    customerServiceDescHi: row.customerServiceDescHi || ''
  }
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
  if (withdrawIconName) {
    setWithdrawIconPreview(row.withdrawIcon || '')
  } else {
    resetWithdrawIconPreview()
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
    await ElMessageBox.confirm(`确定要删除 VIP等级 "${row.levelName}" 吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await vipCfgApi.deleteVipCfg(row.id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const handleSave = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const payload = {
        level: currentRow.value.level,
        levelName: currentRow.value.levelName,
        levelIcon: currentRow.value.levelIcon,
        withdrawSwitch: currentRow.value.withdrawSwitch,
        animationSwitch: currentRow.value.animationSwitch,
        commentEffectSwitch: currentRow.value.commentEffectSwitch,
        customerServiceSwitch: currentRow.value.customerServiceSwitch,
        upgradeRechargeLimit: currentRow.value.upgradeRechargeLimit,
        minWithdrawAmount: currentRow.value.minWithdrawAmount,
        maxWithdrawAmount: currentRow.value.maxWithdrawAmount,
        fee: currentRow.value.fee,
        animation: currentRow.value.animation,
        animationIcon: currentRow.value.animationIcon,
        animationDescEn: currentRow.value.animationDescEn,
        animationDescEs: currentRow.value.animationDescEs,
        animationDescPt: currentRow.value.animationDescPt,
        animationDescHi: currentRow.value.animationDescHi,
        commentEffect: currentRow.value.commentEffect,
        commentEffectIcon: currentRow.value.commentEffectIcon,
        commentEffectDescEn: currentRow.value.commentEffectDescEn,
        commentEffectDescEs: currentRow.value.commentEffectDescEs,
        commentEffectDescPt: currentRow.value.commentEffectDescPt,
        commentEffectDescHi: currentRow.value.commentEffectDescHi,
        withdrawIcon: currentRow.value.withdrawIcon,
        withdrawNoticeEn: currentRow.value.withdrawNoticeEn,
        withdrawNoticeEs: currentRow.value.withdrawNoticeEs,
        withdrawNoticePt: currentRow.value.withdrawNoticePt,
        withdrawNoticeHi: currentRow.value.withdrawNoticeHi,
        customerServiceIcon: currentRow.value.customerServiceIcon,
        customerServiceDescEn: currentRow.value.customerServiceDescEn,
        customerServiceDescEs: currentRow.value.customerServiceDescEs,
        customerServiceDescPt: currentRow.value.customerServiceDescPt,
        customerServiceDescHi: currentRow.value.customerServiceDescHi
      }
      if (currentRow.value.id) {
        await vipCfgApi.updateVipCfg({id: currentRow.value.id, ...payload})
      } else {
        await vipCfgApi.createVipCfg(payload)
      }
      ElMessage.success(currentRow.value.id ? '更新成功' : '创建成功')
      dialogVisible.value = false
      fetchList()
    } catch (error) {
      console.error('保存失败:', error)
      ElMessage.error('保存失败')
    }
  })
}

const resetSearch = () => {
  searchForm.levelName = ''
  searchForm.withdrawSwitchFilter = 0
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
