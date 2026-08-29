<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ t('menu.UserList') }}</span>
        </div>
      </template>
      <div class="search-form">
        <el-form :model="searchForm" inline label-width="100px">
          <el-form-item :label="t('common.keyword')">
            <el-input
                v-model="searchForm.key"
                clearable
                :placeholder="t('pages.userList.keywordPlaceholder')"
                style="width: 200px"
            />
          </el-form-item>
          <el-form-item :label="t('common.startTime')">
            <el-date-picker
                v-model="searchForm.startTime"
                clearable
                format="YYYY-MM-DD"
                :placeholder="t('common.selectStartTime')"
                style="width: 160px"
                type="date"
                value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item :label="t('common.endTime')">
            <el-date-picker
                v-model="searchForm.endTime"
                clearable
                format="YYYY-MM-DD"
                :placeholder="t('common.selectEndTime')"
                style="width: 160px"
                type="date"
                value-format="YYYY-MM-DD"
            />
          </el-form-item>
          <el-form-item :label="t('pages.userList.rechargeWhitelist')">
            <el-select v-model="searchForm.rechargeWhitelist" clearable :placeholder="t('common.all')" style="width: 120px">
              <el-option :value="1" :label="t('common.yes')"/>
              <el-option :value="0" :label="t('common.no')"/>
            </el-select>
          </el-form-item>
          <el-form-item :label="t('pages.userList.isAnchor')">
            <el-select v-model="searchForm.isAnchor" clearable :placeholder="t('common.all')" style="width: 120px">
              <el-option :value="1" :label="t('common.yes')"/>
              <el-option :value="0" :label="t('common.no')"/>
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button v-if="can('search')" type="primary" @click="handleSearch">{{ t('common.query') }}</el-button>
            <el-button v-if="can('search')" @click="handleReset">{{ t('common.reset') }}</el-button>
          </el-form-item>
        </el-form>
      </div>
      <div class="table-scroll">
        <el-table v-loading="loading" :data="userList" style="width: 100%">
          <el-table-column fixed label="#" type="index" width="55" :index="formatRowIndex"/>
          <el-table-column label="ID" prop="id" min-width="240">
            <template #default="scope">
              <div class="id-cell">
                <el-button v-if="canViewDetail" link type="primary" @click="openDetail(scope.row)">
                  {{ scope.row.id }}
                </el-button>
                <span v-else>{{ scope.row.id }}</span>
                <el-button link type="primary" @click="copyUserId(scope.row.id)">{{ t('common.copy') }}</el-button>
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="t('common.createdAt')" prop="createdAt" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.createdAt) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.lastLoginTime')" prop="lastLoginTime" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.lastLoginTime) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('common.nickname')" prop="nickname" width="140">
            <template #default="scope">{{ scope.row.nickname || '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.avatar')" prop="avatar" width="80">
            <template #default="scope">
              <el-image
                  v-if="scope.row.avatar"
                  :preview-src-list="[scope.row.avatar]"
                  :src="scope.row.avatar"
                  fit="cover"
                  hide-on-click-modal
                  preview-teleported
                  style="width:40px;height:40px;border-radius:50%"
              />
              <span v-else class="avatar-empty">-</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.goldBalance')" align="right" min-width="130" prop="gold">
            <template #default="scope">
              <span class="currency-gold">{{ formatWalletBalance(scope.row.gold) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.diamondBalance')" align="right" min-width="130" prop="diamond">
            <template #default="scope">
              <span class="currency-diamond">{{ formatWalletBalance(scope.row.diamond) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.vipLevel')" prop="vipLevel" width="100">
            <template #default="scope">{{ formatVipLevel(scope.row.vipLevel) }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.userType')" prop="userType" width="120">
            <template #default="scope">{{ formatUserType(scope.row.userType) }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.anchor')" prop="isAnchor" width="90">
            <template #default="scope">
              <el-button
                  v-if="isAnchorUser(scope.row) && canViewAnchorDetail"
                  link
                  type="primary"
                  @click="openAnchorDetail(scope.row)"
              >
                {{ t('common.yes') }}
              </el-button>
              <el-tag v-else-if="isAnchorUser(scope.row)" type="success">{{ t('common.yes') }}</el-tag>
              <el-tag v-else type="info">{{ t('common.no') }}</el-tag>
            </template>
          </el-table-column>

          <el-table-column label="IP" prop="ip" width="150"/>
          <el-table-column :label="t('pages.userList.loginCountry')" prop="loginCountry" width="120">
            <template #default="scope">{{ scope.row.loginCountry || '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.registerIp')" prop="registerIp" width="150">
            <template #default="scope">{{ scope.row.registerIp || '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.registerCountry')" prop="registerCountry" width="120">
            <template #default="scope">{{ scope.row.registerCountry || '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.deviceType')" prop="deviceType" width="100">
            <template #default="scope">{{ scope.row.deviceType || '-' }}</template>
          </el-table-column>
           <el-table-column :label="t('pages.userList.packageName')" prop="packageName" width="180" show-overflow-tooltip>
            <template #default="scope">{{ scope.row.packageName || '-' }}</template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.appVersion')" prop="appVersion" width="120" show-overflow-tooltip>
            <template #default="scope">{{ scope.row.appVersion || '-' }}</template>
          </el-table-column>
          
          <el-table-column :label="t('pages.userList.banStatus')" prop="ban" width="120">
            <template #default="scope">
              <el-tag v-if="scope.row.ban" type="danger">{{ t('pages.userList.banned') }}</el-tag>
              <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.cancelStatus')" prop="cancel" width="120">
            <template #default="scope">
              <el-tag v-if="scope.row.cancel" type="warning">{{ t('pages.userList.canceled') }}</el-tag>
              <el-tag v-else type="success">{{ t('common.normal') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.userList.rechargeWhitelist')" prop="rechargeWhitelist" width="110">
            <template #default="scope">
              <el-tag v-if="scope.row.rechargeWhitelist" type="success">{{ t('common.yes') }}</el-tag>
              <el-tag v-else type="info">{{ t('common.no') }}</el-tag>
            </template>
          </el-table-column>
          
          <el-table-column :label="t('pages.userList.banTime')" prop="banApplyTime" width="200">
            <template #default="scope">
              {{ formatDate(scope.row.banApplyTime) }}
            </template>
          </el-table-column>
           <el-table-column :label="t('pages.userList.openId')" prop="openId" width="200"/>
          <el-table-column :label="t('pages.userList.channel')" prop="channel" width="100"/>
          <el-table-column :label="t('common.remark')" prop="remark" min-width="160" show-overflow-tooltip>
            <template #default="scope">{{ scope.row.remark || '-' }}</template>
          </el-table-column>
          <el-table-column fixed="right" :label="t('common.actions')" width="100">
            <template #default="scope">
              <el-dropdown v-if="hasRowActions" trigger="click" @command="(cmd: string) => handleRowCommand(scope.row, cmd)">
                <el-button size="small" type="primary">
                  {{ t('common.actions') }}
                  <el-icon class="el-icon--right">
                    <ArrowDown/>
                  </el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item v-if="can('openGame')" command="openGame">{{ t('pages.userList.openGame') }}</el-dropdown-item>
                    <el-dropdown-item v-if="canViewDetail" :divided="can('openGame')" command="viewDetail">{{ t('pages.userList.viewDetail') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('setAnchor') && !scope.row.isAnchor" :divided="canViewDetail || can('openGame')" command="setAnchor">{{ t('pages.userList.setAnchor') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('setSeniorAnchor') && !scope.row.isAnchor" command="setSeniorAnchor">{{ t('pages.userList.setSeniorAnchor') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('goldAdd')" :divided="!scope.row.isAnchor" command="gold-add">
                      {{ t('pages.userList.addGold') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('goldSub')" command="gold-sub">{{ t('pages.userList.subGold') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('diamondAdd')" divided command="diamond-add">{{ t('pages.userList.addDiamond') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('diamondSub')" command="diamond-sub">{{ t('pages.userList.subDiamond') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('ban')" divided command="ban">
                      {{ scope.row.ban ? t('pages.userList.unban') : t('pages.userList.ban') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('rankOff') && scope.row.canRank !== false" command="rank-off">{{ t('pages.userList.rankOff') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('rankOn') && scope.row.canRank === false" command="rank-on">{{ t('pages.userList.rankOn') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('rechargeWhitelistOn') && !scope.row.rechargeWhitelist" command="recharge-whitelist-on">{{ t('pages.userList.rechargeWhitelistOn') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('rechargeWhitelistOff') && scope.row.rechargeWhitelist" command="recharge-whitelist-off">{{ t('pages.userList.rechargeWhitelistOff') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('cancel')" divided command="cancel">
                      {{ scope.row.cancel ? t('pages.userList.uncancel') : t('pages.userList.cancelAccount') }}
                    </el-dropdown-item>
                    <el-dropdown-item v-if="can('setUserType') && !scope.row.isAnchor" command="setUserType">{{ t('pages.userList.setUserType') }}</el-dropdown-item>
                    <el-dropdown-item v-if="can('uploadAvatar')" command="uploadAvatar">{{ t('pages.userList.uploadAvatar') }}</el-dropdown-item>
                    <el-dropdown-item v-if="showSetAnchorType && isAnchorUser(scope.row)" command="setAnchorType">{{ t('pages.guildMembers.setAnchorType') }}</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-table-column>

        </el-table>
        <div class="pagination">
          <el-pagination
              :current-page="pagination.pageIndex"
              :page-size="pagination.pageSize"
              :total="pagination.total"
              layout="total, sizes, prev, pager, next, jumper"
              @current-change="handlePageChange"
              @size-change="handleSizeChange"
          />
        </div>
      </div>
    </el-card>

    <el-dialog
        v-model="currencyDialogVisible"
        :title="currencyDialogTitle"
        width="440px"
        @closed="resetCurrencyForm"
    >
      <el-form ref="currencyFormRef" :model="currencyForm" :rules="currencyFormRules" label-width="100px">
        <el-form-item :label="t('common.userId')">
          <el-input v-model="currencyForm.userId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
          <el-input v-model="currencyForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="currencyBalanceLabel">
          <el-input v-model="currencyForm.currentBalanceText" disabled/>
        </el-form-item>
        <el-form-item :label="currencyAmountLabel" prop="amount">
          <el-input-number
              v-model="currencyForm.amount"
              :min="0.0001"
              :precision="NUMBER_INPUT_DECIMALS"
              :step="0.0001"
              controls-position="right"
              style="width: 100%"
          />
        </el-form-item>
        <el-form-item :label="t('pages.userList.adjustReason')" prop="reason">
          <el-select v-model="currencyForm.reason" :placeholder="t('pages.userList.selectAdjustReason')" style="width: 100%">
            <el-option
                v-for="item in currencyAdjustReasonOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="currencyDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="submitCurrencyChange">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="banDialogVisible"
        :title="t('pages.userList.banDialogTitle')"
        width="480px"
        @closed="resetBanForm"
    >
      <el-form ref="banFormRef" :model="banForm" :rules="banFormRules" label-width="110px">
        <el-form-item :label="t('common.userId')">
          <el-input v-model="banForm.userId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
          <el-input v-model="banForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.userList.banUntil')" prop="banApplyTime">
          <el-date-picker
              v-model="banForm.banApplyTime"
              :disabled-date="disabledBanDate"
              format="YYYY-MM-DD HH:mm:ss"
              :placeholder="t('pages.userList.selectBanUntil')"
              style="width: 100%"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="banDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="banSubmitting" type="primary" @click="submitBan">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="userTypeDialogVisible"
        :title="t('pages.userList.editUserTypeTitle')"
        width="440px"
        @closed="resetUserTypeForm"
    >
      <el-form ref="userTypeFormRef" :model="userTypeForm" :rules="userTypeFormRules" label-width="100px">
        <el-form-item :label="t('common.userId')">
          <el-input v-model="userTypeForm.userId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
          <el-input v-model="userTypeForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.userList.userType')" prop="userType">
          <el-select v-model="userTypeForm.userType" :placeholder="t('pages.userList.selectUserType')" style="width: 100%">
            <el-option
                v-for="item in userTypeOptions"
                :key="item.value"
                :label="item.label"
                :value="item.value"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="userTypeDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="userTypeSubmitting" type="primary" @click="submitUserType">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="anchorTypeDialogVisible"
        :title="anchorTypeDialogTitle"
        width="480px"
        @closed="resetAnchorTypeForm"
    >
      <el-form ref="anchorTypeFormRef" :model="anchorTypeForm" :rules="anchorTypeFormRules" label-width="100px">
        <el-form-item :label="t('common.userId')">
          <el-input v-model="anchorTypeForm.userId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
          <el-input v-model="anchorTypeForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.anchorList.anchorType')" prop="anchorType">
          <el-select v-model="anchorTypeForm.anchorType" style="width: 100%">
            <el-option :label="t('pages.anchorList.anchorTypeNormal')" :value="1"/>
            <el-option :label="t('pages.anchorList.anchorTypeSenior')" :value="7"/>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="anchorTypeDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="anchorTypeSubmitting" type="primary" @click="submitAnchorType">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="avatarDialogVisible"
        :title="t('pages.userList.uploadAvatarTitle')"
        width="440px"
        @closed="resetAvatarForm"
    >
      <el-form label-width="100px">
        <el-form-item :label="t('common.userId')">
          <el-input v-model="avatarForm.userId" disabled/>
        </el-form-item>
        <el-form-item :label="t('common.nickname')">
          <el-input v-model="avatarForm.nickname" disabled/>
        </el-form-item>
        <el-form-item :label="t('pages.userList.avatar')">
          <div class="avatar-upload-wrap">
            <el-upload
                :before-upload="beforeAvatarUpload"
                :disabled="avatarUploading"
                :http-request="doAvatarUpload"
                :show-file-list="false"
                accept="image/*"
                class="avatar-uploader"
            >
              <el-image
                  v-if="avatarPreviewUrl"
                  :src="avatarPreviewUrl"
                  fit="cover"
                  style="width:80px;height:80px;border-radius:50%"
              />
              <div v-else class="avatar-uploader-placeholder">
                <el-icon class="avatar-uploader-icon">
                  <Plus/>
                </el-icon>
              </div>
            </el-upload>
            <el-button v-if="avatarForm.avatar || avatarPreviewUrl" link type="danger" @click="clearAvatar">
              {{ t('pages.userList.clearAvatar') }}
            </el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="avatarDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button :loading="avatarSubmitting" type="primary" @click="submitAvatar">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <el-dialog
        v-model="gamePickerDialogVisible"
        :title="gamePickerDialogTitle"
        destroy-on-close
        width="920px"
        @closed="resetGamePicker"
        @open="handleGamePickerOpen"
    >
      <p class="game-picker-tip">{{ gamePickerTipText }}</p>
      <el-form :model="gamePickerSearchForm" class="search-form" inline>
        <el-form-item :label="t('pages.gameList.gameCode')">
          <el-input v-model="gamePickerSearchForm.gameCode" clearable :placeholder="t('pages.gameList.gameCodePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.gameList.gameName')">
          <el-input v-model="gamePickerSearchForm.name" clearable :placeholder="t('pages.gameList.gameNamePlaceholder')"/>
        </el-form-item>
        <el-form-item :label="t('pages.gameList.platform')">
          <el-input v-model="gamePickerSearchForm.platform" clearable :placeholder="t('pages.gameList.platformPlaceholder')"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleGamePickerSearch">{{ t('common.search') }}</el-button>
          <el-button @click="resetGamePickerSearch">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
      <el-table v-loading="gamePickerLoading" :data="gamePickerTableData" row-key="gameCode" style="width: 100%">
        <el-table-column :label="t('pages.gameList.gameCode')" min-width="140" prop="gameCode"/>
        <el-table-column :label="t('common.name')" min-width="120">
          <template #default="{ row }">{{ row.name || row.nameEn || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameList.nameEn')" min-width="120" prop="nameEn">
          <template #default="{ row }">{{ row.nameEn || '-' }}</template>
        </el-table-column>
        <el-table-column :label="t('pages.gameList.platform')" prop="platform" width="100"/>
        <el-table-column fixed="right" :label="t('common.actions')" width="120">
          <template #default="{ row }">
            <el-button
                :loading="gamePickerStartingCode === row.gameCode"
                link
                type="success"
                @click="handleGamePickerStart(row)"
            >
              {{ t('pages.gameShelfList.startGame') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
            v-model:current-page="gamePickerPagination.pageIndex"
            v-model:page-size="gamePickerPagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="gamePickerPagination.total"
            layout="total, sizes, prev, pager, next"
            @current-change="fetchGamePickerList"
            @size-change="handleGamePickerSizeChange"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts" setup>
import {computed, onMounted, reactive, ref, watch} from 'vue'
import {useI18n} from 'vue-i18n'
import {accountApi, diamondApi, goldApi, guildApi, uploadApi} from '@/api'
import {ArrowDown, Plus} from '@element-plus/icons-vue'
import {ElForm, ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadRequestOptions} from 'element-plus'
import {useRoute, useRouter} from 'vue-router'
import type {BanReq, CancelReq, GameShelfItem, UnBanReq, UnCancelReq, UserInfo} from '@/types/api.ts'
import {gamePlatformApi} from '@/api/modules/gamePlatform'
import {formatWalletBalance, NUMBER_INPUT_DECIMALS} from '@/utils/number-format'
import {usePagePermission} from '@/composables/usePagePermission'

const {t, locale} = useI18n()
const {can} = usePagePermission('UserList')
const canViewDetail = computed(() => can('viewDetail'))
const canViewAnchorDetail = computed(() => can('viewAnchorDetail'))
const showSetAnchorType = computed(() => can('setAnchorType'))

const USER_TYPE_ANCHOR = 1
const USER_TYPE_SENIOR_ANCHOR = 7

const ROW_ACTION_KEYS = [
  'setAnchor',
  'setSeniorAnchor',
  'setAnchorType',
  'goldAdd',
  'goldSub',
  'diamondAdd',
  'diamondSub',
  'ban',
  'rankOff',
  'rankOn',
  'rechargeWhitelistOn',
  'rechargeWhitelistOff',
  'cancel',
  'setUserType',
  'uploadAvatar',
  'openGame',
] as const

const hasRowActions = computed(() => canViewDetail.value || ROW_ACTION_KEYS.some(key => can(key)))

const userList = ref<UserInfo[]>([])
const loading = ref(false)
type CurrencyType = 'gold' | 'diamond'
type CurrencyMode = 'add' | 'sub'

const GM_ADJUST_REASON_TEST = 6
const GM_ADJUST_REASON_COMPENSATION = 7

const currencyDialogVisible = ref(false)
const currencyType = ref<CurrencyType>('gold')
const currencyMode = ref<CurrencyMode>('add')
const currencyFormRef = ref<FormInstance>()
const banDialogVisible = ref(false)
const banSubmitting = ref(false)
const banFormRef = ref<FormInstance>()
const userTypeDialogVisible = ref(false)
const userTypeSubmitting = ref(false)
const userTypeFormRef = ref<FormInstance>()
const anchorTypeDialogVisible = ref(false)
const anchorTypeDialogTitle = ref('')
const anchorTypeSubmitting = ref(false)
const anchorTypeFormRef = ref<InstanceType<typeof ElForm>>()

const avatarDialogVisible = ref(false)
const avatarSubmitting = ref(false)
const avatarUploading = ref(false)
const avatarPreviewUrl = ref('')
const avatarChanged = ref(false)
const avatarForm = reactive({
  userId: '',
  nickname: '',
  avatar: '',
})

const gamePickerDialogVisible = ref(false)
const gamePickerLoading = ref(false)
const gamePickerStartingCode = ref('')
const gamePickerTableData = ref<GameShelfItem[]>([])
const gamePickerUser = ref<UserInfo | null>(null)
const gamePickerSearchForm = reactive({
  gameCode: '',
  name: '',
  platform: '',
})
const gamePickerPagination = reactive({
  pageIndex: 1,
  pageSize: 20,
  total: 0,
})

const gamePickerDialogTitle = computed(() => {
  const user = gamePickerUser.value
  if (!user) {
    return t('pages.userList.openGame')
  }
  if (user.nickname) {
    return t('pages.userList.openGameDialogTitle', {name: user.nickname})
  }
  return t('pages.userList.openGameDialogTitleNoName', {id: user.id})
})

const gamePickerTipText = computed(() => {
  const user = gamePickerUser.value
  if (!user) {
    return ''
  }
  if (user.nickname) {
    return t('pages.gameShelfList.pickUserHint', {name: user.nickname, id: user.id})
  }
  return t('pages.gameShelfList.pickUserHintNoName', {id: user.id})
})

const userTypeLabelMap = computed<Record<number, string>>(() => ({
  0: t('pages.userList.userTypeNormal'),
  1: t('pages.userList.userTypeAnchor'),
  2: t('pages.userList.userTypeBotAnchor'),
  3: t('pages.userList.userTypeBotViewer'),
  4: t('pages.userList.userTypeTester'),
  5: t('pages.userList.userTypeCmsAuthor'),
  7: t('pages.userList.userTypeSeniorAnchor'),
}))

const userTypeOptions = computed(() => [
  {label: t('pages.userList.userTypeNormal'), value: 0},
  {label: t('pages.userList.userTypeTester'), value: 4}
])

interface UserTypeForm {
  userId: string
  nickname: string
  userType: number
}

const userTypeForm = reactive<UserTypeForm>({
  userId: '',
  nickname: '',
  userType: 0
})

const userTypeFormRules = computed<FormRules>(() => ({
  userType: [
    {required: true, message: t('pages.userList.selectUserTypeRequired'), trigger: 'change'}
  ]
}))

interface AnchorTypeForm {
  userId: string
  nickname: string
  guildId: string
  anchorType: 1 | 7
}

const anchorTypeForm = reactive<AnchorTypeForm>({
  userId: '',
  nickname: '',
  guildId: '0',
  anchorType: 1,
})

const anchorTypeFormRules = computed<FormRules>(() => ({
  anchorType: [
    {required: true, message: t('pages.guildMembers.anchorTypeRequired'), trigger: 'change'},
  ],
}))

const isAnchorUser = (row: UserInfo) =>
    !!row.isAnchor || row.userType === USER_TYPE_ANCHOR || row.userType === USER_TYPE_SENIOR_ANCHOR

interface BanForm {
  userId: string
  nickname: string
  openId: string
  channel: number
  banApplyTime: string
}

const banForm = reactive<BanForm>({
  userId: '',
  nickname: '',
  openId: '',
  channel: 0,
  banApplyTime: ''
})

const banFormRules = computed<FormRules>(() => ({
  banApplyTime: [
    {required: true, message: t('pages.userList.selectBanUntilRequired'), trigger: 'change'}
  ]
}))

interface CurrencyForm {
  userId: string
  nickname: string
  currentBalanceText: string
  amount: number
  reason: number
}

const currencyForm = reactive<CurrencyForm>({
  userId: '',
  nickname: '',
  currentBalanceText: '',
  amount: 1,
  reason: GM_ADJUST_REASON_TEST,
})

const currencyAdjustReasonOptions = computed(() => [
  {label: t('pages.userList.adjustReasonTest'), value: GM_ADJUST_REASON_TEST},
  {label: t('pages.userList.adjustReasonCompensation'), value: GM_ADJUST_REASON_COMPENSATION},
])

const currencyName = computed(() =>
    currencyType.value === 'gold' ? t('pages.userList.gold') : t('pages.userList.diamond')
)

const currencyDialogTitle = computed(() =>
    currencyMode.value === 'add'
        ? t('pages.userList.addCurrency', {currency: currencyName.value})
        : t('pages.userList.subCurrency', {currency: currencyName.value})
)

const currencyBalanceLabel = computed(() =>
    t('pages.userList.currentBalance', {balance: currencyName.value})
)

const currencyAmountLabel = computed(() =>
    currencyMode.value === 'add' ? t('pages.userList.addAmount') : t('pages.userList.subAmount')
)

const currencyFormRules = computed<FormRules>(() => ({
  amount: [
    {required: true, message: t('pages.userList.amountRequired'), trigger: 'blur'},
    {
      validator: (_rule, value, callback) => {
        const n = Number(value)
        if (!n || n <= 0) {
          callback(new Error(t('pages.userList.amountPositive')))
          return
        }
        callback()
      },
      trigger: 'blur'
    }
  ],
  reason: [
    {required: true, message: t('pages.userList.selectAdjustReasonRequired'), trigger: 'change'},
  ],
}))

const searchForm = reactive({
  key: '',
  startTime: '',
  endTime: '',
  rechargeWhitelist: undefined as number | undefined,
  isAnchor: undefined as number | undefined,
})

const router = useRouter()
const route = useRoute()

const pagination = reactive({
  pageIndex: 1,
  pageSize: 10,
  total: 0
})

watch(() => route.query.refresh, (newRefresh) => {
  if (newRefresh) {
    fetchUserList()
  }
}, {immediate: false})

const fetchUserList = async (silent = false) => {
  if (!silent) {
    loading.value = true
  }
  try {
    const params: Record<string, unknown> = {
      pageIndex: pagination.pageIndex,
      pageSize: pagination.pageSize,
      key: searchForm.key,
      startTime: searchForm.startTime,
      endTime: searchForm.endTime,
    }
    if (searchForm.rechargeWhitelist !== undefined && searchForm.rechargeWhitelist !== null) {
      params.rechargeWhitelist = searchForm.rechargeWhitelist
    }
    if (searchForm.isAnchor !== undefined && searchForm.isAnchor !== null) {
      params.isAnchor = searchForm.isAnchor
    }

    const response = await accountApi.getUserInfo(params as Parameters<typeof accountApi.getUserInfo>[0])

    const result = response.data
    userList.value = result || []
    pagination.total = response?.total || 0
  } catch (error) {
    console.error('fetchUserList failed:', error)
    ElMessage.error(t('pages.userList.fetchFailed'))
  } finally {
    if (!silent) {
      loading.value = false
    }
  }
}

const handleSearch = () => {
  pagination.pageIndex = 1
  fetchUserList()
}

const handleReset = () => {
  searchForm.key = ''
  searchForm.startTime = ''
  searchForm.endTime = ''
  searchForm.rechargeWhitelist = undefined
  searchForm.isAnchor = undefined

  pagination.pageIndex = 1
  fetchUserList()
}

const handlePageChange = (page: number) => {
  pagination.pageIndex = page
  fetchUserList()
}

const resolveAccountOpenId = (row: UserInfo): string => {
  let openId = row.openId || ''
  const canceledIdx = openId.indexOf('__canceled__')
  if (canceledIdx > 0) {
    openId = openId.slice(0, canceledIdx)
  }
  return openId
}

const accountActionPayload = (row: UserInfo) => ({
  accountId: row.id,
  openId: resolveAccountOpenId(row),
  channel: row.channel ?? 0,
})

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.pageIndex = 1
  fetchUserList()
}

const formatRowIndex = (index: number) =>
    (pagination.pageIndex - 1) * pagination.pageSize + index + 1

const resetCurrencyForm = () => {
  currencyForm.userId = ''
  currencyForm.nickname = ''
  currencyForm.currentBalanceText = ''
  currencyForm.amount = 1
  currencyForm.reason = GM_ADJUST_REASON_TEST
  currencyFormRef.value?.clearValidate()
}

const openCurrencyDialog = (row: UserInfo, type: CurrencyType, mode: CurrencyMode) => {
  currencyType.value = type
  currencyMode.value = mode
  currencyForm.userId = String(row.id)
  currencyForm.nickname = row.nickname || '-'
  currencyForm.currentBalanceText = formatWalletBalance(type === 'gold' ? row.gold : row.diamond)
  currencyForm.amount = 1
  currencyForm.reason = GM_ADJUST_REASON_TEST
  currencyDialogVisible.value = true
}

const openDetail = (row: UserInfo) => {
  router.push({
    name: 'UserDetail',
    query: {id: String(row.id)},
  })
}

const copyUserId = async (userId: string | number | undefined) => {
  const value = String(userId ?? '').trim()
  if (!value) {
    return
  }
  try {
    await navigator.clipboard.writeText(value)
    ElMessage.success(t('common.operationSuccess'))
  } catch (error) {
    console.error('copy user id failed:', error)
    ElMessage.error(t('common.failed'))
  }
}

const openAnchorDetail = (row: UserInfo) => {
  router.push({
    path: '/user/anchor/anchor-detail',
    query: {id: String(row.id)},
  })
}

const resetGamePickerSearch = () => {
  gamePickerSearchForm.gameCode = ''
  gamePickerSearchForm.name = ''
  gamePickerSearchForm.platform = ''
  gamePickerPagination.pageIndex = 1
}

const resetGamePicker = () => {
  gamePickerUser.value = null
  gamePickerTableData.value = []
  gamePickerStartingCode.value = ''
  gamePickerPagination.total = 0
  resetGamePickerSearch()
}

const fetchGamePickerList = async () => {
  gamePickerLoading.value = true
  try {
    const response = await gamePlatformApi.getGameShelfList({
      gameCode: gamePickerSearchForm.gameCode,
      name: gamePickerSearchForm.name,
      platform: gamePickerSearchForm.platform,
      pageIndex: gamePickerPagination.pageIndex,
      pageSize: gamePickerPagination.pageSize,
    })
    gamePickerTableData.value = response.data || []
    gamePickerPagination.total = response.total || 0
  } catch (error) {
    console.error('fetch game shelf list failed:', error)
    ElMessage.error(t('pages.gameShelfList.fetchFailed'))
  } finally {
    gamePickerLoading.value = false
  }
}

const handleGamePickerOpen = () => {
  fetchGamePickerList()
}

const handleGamePickerSearch = () => {
  gamePickerPagination.pageIndex = 1
  fetchGamePickerList()
}

const handleGamePickerSizeChange = (size: number) => {
  gamePickerPagination.pageSize = size
  gamePickerPagination.pageIndex = 1
  fetchGamePickerList()
}

const handleGamePickerStart = async (row: GameShelfItem) => {
  const user = gamePickerUser.value
  if (!user?.id) {
    return
  }
  gamePickerStartingCode.value = row.gameCode
  try {
    const response = await gamePlatformApi.getCMSGameStartLink({
      userId: String(user.id),
      gameCode: row.gameCode,
      platform: row.platform,
    })
    const link = response?.link?.trim()
    if (!link) {
      ElMessage.error(t('pages.gameShelfList.startGameEmpty'))
      return
    }
    window.open(link, '_blank', 'noopener,noreferrer')
  } catch (error) {
    console.error('get cms game start link failed:', error)
    ElMessage.error(t('pages.gameShelfList.startGameFailed'))
  } finally {
    gamePickerStartingCode.value = ''
  }
}

const openGamePicker = (row: UserInfo) => {
  gamePickerUser.value = row
  resetGamePickerSearch()
  gamePickerPagination.pageSize = 20
  gamePickerDialogVisible.value = true
}

const handleRowCommand = (row: UserInfo, command: string) => {
  switch (command) {
    case 'openGame':
      openGamePicker(row)
      break
    case 'viewDetail':
      openDetail(row)
      break
    case 'setUserType':
      openUserTypeDialog(row)
      break
    case 'gold-add':
      openCurrencyDialog(row, 'gold', 'add')
      break
    case 'gold-sub':
      openCurrencyDialog(row, 'gold', 'sub')
      break
    case 'diamond-add':
      openCurrencyDialog(row, 'diamond', 'add')
      break
    case 'diamond-sub':
      openCurrencyDialog(row, 'diamond', 'sub')
      break
    case 'ban':
      handleBanAction(row)
      break
    case 'rank-on':
      handleCanRankAction(row, true)
      break
    case 'rank-off':
      handleCanRankAction(row, false)
      break
    case 'recharge-whitelist-on':
      handleRechargeWhitelistAction(row, true)
      break
    case 'recharge-whitelist-off':
      handleRechargeWhitelistAction(row, false)
      break
    case 'cancel':
      toggleCancelStatus(row)
      break
    case 'setAnchor':
      handleSetAnchor(row)
      break
    case 'setSeniorAnchor':
      handleSetSeniorAnchor(row)
      break
    case 'setAnchorType':
      openAnchorTypeDialog(row)
      break
    case 'uploadAvatar':
      openAvatarDialog(row)
      break
  }
}

const resetAvatarForm = () => {
  avatarForm.userId = ''
  avatarForm.nickname = ''
  avatarForm.avatar = ''
  avatarPreviewUrl.value = ''
  avatarChanged.value = false
}

const openAvatarDialog = (row: UserInfo) => {
  avatarForm.userId = String(row.id)
  avatarForm.nickname = row.nickname || '-'
  avatarForm.avatar = ''
  avatarPreviewUrl.value = row.avatar || ''
  avatarChanged.value = false
  avatarDialogVisible.value = true
}

const beforeAvatarUpload = (file: File): boolean => {
  if (!file.type.startsWith('image/')) {
    ElMessage.error(t('pages.userList.imageOnly'))
    return false
  }
  return true
}

const doAvatarUpload = async (options: UploadRequestOptions) => {
  const file = options.file as File
  avatarUploading.value = true
  try {
    const res = await uploadApi.uploadFile(file)
    avatarForm.avatar = res.fileName
    avatarChanged.value = true
    avatarPreviewUrl.value = URL.createObjectURL(file)
    ElMessage.success(t('pages.userList.uploadSuccess'))
  } catch (error) {
    console.error('upload avatar failed:', error)
    ElMessage.error(t('pages.userList.uploadFailed'))
  } finally {
    avatarUploading.value = false
  }
}

const clearAvatar = () => {
  avatarForm.avatar = ''
  avatarPreviewUrl.value = ''
  avatarChanged.value = true
}

const submitAvatar = async () => {
  if (!avatarChanged.value) {
    avatarDialogVisible.value = false
    return
  }
  avatarSubmitting.value = true
  try {
    const res = await accountApi.setUserAvatar({
      accountId: avatarForm.userId,
      avatar: avatarForm.avatar,
    })
    if (res?.success) {
      avatarDialogVisible.value = false
      ElMessage.success(t('pages.userList.uploadAvatarSuccess'))
      await fetchUserList(true)
    } else {
      ElMessage.error(t('pages.userList.uploadFailed'))
    }
  } catch (error) {
    console.error('setUserAvatar failed:', error)
  } finally {
    avatarSubmitting.value = false
  }
}

const resetUserTypeForm = () => {
  userTypeForm.userId = ''
  userTypeForm.nickname = ''
  userTypeForm.userType = 0
  userTypeFormRef.value?.clearValidate()
}

const openUserTypeDialog = (row: UserInfo) => {
  userTypeForm.userId = String(row.id)
  userTypeForm.nickname = row.nickname || '-'
  const currentType = row.userType ?? 0
  userTypeForm.userType = currentType === 4 ? 4 : 0
  userTypeDialogVisible.value = true
}

const submitUserType = async () => {
  if (!userTypeFormRef.value) return
  await userTypeFormRef.value.validate(async (valid) => {
    if (!valid) return
    userTypeSubmitting.value = true
    try {
      await accountApi.setUserType({
        accountId: userTypeForm.userId,
        userType: userTypeForm.userType
      })
      userTypeDialogVisible.value = false
      ElMessage.success(t('pages.userList.userTypeUpdated'))
      await fetchUserList()
    } catch (error) {
      console.error('setUserType failed:', error)
    } finally {
      userTypeSubmitting.value = false
    }
  })
}

const handleSetAnchor = async (row: UserInfo) => {
  try {
    await ElMessageBox.confirm(
        t('pages.userList.setAnchorConfirm', {id: row.id}),
        t('pages.userList.setAnchorTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await accountApi.setAnchor({accountId: String(row.id)})
    ElMessage.success(t('pages.userList.setAnchorSuccess'))
    await fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('setAnchor failed:', error)
    }
  }
}

const handleSetSeniorAnchor = async (row: UserInfo) => {
  try {
    await ElMessageBox.confirm(
        t('pages.userList.setSeniorAnchorConfirm', {id: row.id}),
        t('pages.userList.setSeniorAnchorTitle'),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    await accountApi.setSeniorAnchor({accountId: String(row.id)})
    ElMessage.success(t('pages.userList.setSeniorAnchorSuccess'))
    await fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('setSeniorAnchor failed:', error)
    }
  }
}

const resetAnchorTypeForm = () => {
  anchorTypeForm.userId = ''
  anchorTypeForm.nickname = ''
  anchorTypeForm.guildId = '0'
  anchorTypeForm.anchorType = 1
  anchorTypeFormRef.value?.clearValidate()
}

const openAnchorTypeDialog = (row: UserInfo) => {
  anchorTypeDialogTitle.value = t('pages.guildMembers.setAnchorTypeTitle', {id: row.id})
  anchorTypeForm.userId = String(row.id)
  anchorTypeForm.nickname = row.nickname || '-'
  anchorTypeForm.guildId = row.guildId != null && String(row.guildId) !== '' ? String(row.guildId) : '0'
  const userType = row.userType ?? USER_TYPE_ANCHOR
  anchorTypeForm.anchorType = userType === USER_TYPE_SENIOR_ANCHOR ? 7 : 1
  anchorTypeDialogVisible.value = true
}

const submitAnchorType = async () => {
  if (!anchorTypeFormRef.value) {
    return
  }
  await anchorTypeFormRef.value.validate(async (valid: boolean) => {
    if (!valid) {
      return
    }
    anchorTypeSubmitting.value = true
    try {
      const payload = {
        userId: anchorTypeForm.userId,
        anchorType: anchorTypeForm.anchorType,
      }
      const response = anchorTypeForm.guildId === '0'
          ? await accountApi.setPlatformAnchorType(payload)
          : await guildApi.setGuildAnchorType({
            guildId: anchorTypeForm.guildId,
            ...payload,
          })
      if (response?.success) {
        ElMessage.success(t('pages.guildMembers.setAnchorTypeSuccess'))
        anchorTypeDialogVisible.value = false
        await fetchUserList()
      } else {
        ElMessage.error(t('pages.guildMembers.setAnchorTypeFailed'))
      }
    } catch (error) {
      console.error('setAnchorType failed:', error)
      ElMessage.error(t('pages.guildMembers.setAnchorTypeRequestFailed'))
    } finally {
      anchorTypeSubmitting.value = false
    }
  })
}

const afterCurrencyChangeSuccess = () => {
  currencyDialogVisible.value = false
  ElMessage.success(t('common.operationSuccess'))
  fetchUserList(true)
}

const submitCurrencyChange = async () => {
  if (!currencyFormRef.value) return
  await currencyFormRef.value.validate(async (valid) => {
    if (!valid) return
    try {
      const payload = {
        userId: currencyForm.userId,
        amount: currencyForm.amount,
        reason: currencyForm.reason,
      }
      const api = currencyType.value === 'gold' ? goldApi : diamondApi
      if (currencyMode.value === 'add') {
        await api.add(payload)
      } else {
        await api.sub(payload)
      }
      afterCurrencyChangeSuccess()
    } catch (error) {
      console.error('currencyChange failed:', error)
    }
  })
}

const defaultBanApplyTime = () => formatServerNowPlusDays(7)

const disabledBanDate = (time: Date) => time.getTime() < Date.now()

const resetBanForm = () => {
  banForm.userId = ''
  banForm.nickname = ''
  banForm.banApplyTime = ''
  banFormRef.value?.clearValidate()
}

const openBanDialog = (row: UserInfo) => {
  banForm.userId = String(row.id)
  banForm.nickname = row.nickname || '-'
  banForm.openId = resolveAccountOpenId(row)
  banForm.channel = row.channel ?? 0
  banForm.banApplyTime = defaultBanApplyTime()
  banDialogVisible.value = true
}

const submitBan = async () => {
  if (!banFormRef.value) return
  await banFormRef.value.validate(async (valid) => {
    if (!valid) return
    banSubmitting.value = true
    try {
      const banData: BanReq = {
        accountId: banForm.userId,
        openId: banForm.openId,
        channel: banForm.channel,
        banApplyTime: banForm.banApplyTime
      }
      const response = await accountApi.ban(banData)
      if (response) {
        banDialogVisible.value = false
        ElMessage.success(t('pages.userList.banSuccess'))
        await fetchUserList()
      } else {
        ElMessage.error(t('pages.userList.banFailed'))
      }
    } catch (error) {
      console.error('ban failed:', error)
    } finally {
      banSubmitting.value = false
    }
  })
}

const handleCanRankAction = async (row: UserInfo, canRank: boolean) => {
  const actionText = canRank ? t('pages.userList.rankOnAction') : t('pages.userList.rankOffAction')
  try {
    await ElMessageBox.confirm(
        t('pages.userList.rankConfirm', {id: row.id, action: actionText}),
        t('pages.userList.rankTitle', {action: actionText}),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    const response = await accountApi.setCanRank({
      accountId: row.id,
      canRank
    })
    if (response) {
      ElMessage.success(t('pages.userList.rankSuccess', {action: actionText}))
      await fetchUserList()
    } else {
      ElMessage.error(t('pages.userList.rankFailed', {action: actionText}))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('canRank action failed:', error)
    }
  }
}

const handleRechargeWhitelistAction = async (row: UserInfo, rechargeWhitelist: boolean) => {
  const actionText = rechargeWhitelist
      ? t('pages.userList.rechargeWhitelistJoin')
      : t('pages.userList.rechargeWhitelistRemove')
  const extra = rechargeWhitelist ? t('pages.userList.rechargeWhitelistExtra') : ''
  try {
    await ElMessageBox.confirm(
        t('pages.userList.rechargeWhitelistConfirm', {id: row.id, action: actionText, extra}),
        t('pages.userList.rankTitle', {action: actionText}),
        {
          confirmButtonText: t('common.confirm'),
          cancelButtonText: t('common.cancel'),
          type: 'warning'
        }
    )
    const response = await accountApi.setRechargeWhitelist({
      accountId: row.id,
      rechargeWhitelist
    })
    if (response) {
      ElMessage.success(t('pages.userList.rankSuccess', {action: actionText}))
      await fetchUserList()
    } else {
      ElMessage.error(t('pages.userList.rankFailed', {action: actionText}))
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('rechargeWhitelist action failed:', error)
    }
  }
}

const handleBanAction = async (row: UserInfo) => {
  if (row.ban) {
    try {
      await ElMessageBox.confirm(
          t('pages.userList.unbanConfirm', {id: row.id}),
          t('pages.userList.unbanTitle'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
          }
      )

      const response = await accountApi.unBan(accountActionPayload(row))

      if (response) {
        ElMessage.success(t('pages.userList.unbanSuccess'))
        await fetchUserList()
      } else {
        ElMessage.error(t('pages.userList.unbanFailed'))
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('unban failed:', error)
      }
    }
    return
  }
  openBanDialog(row)
}

const patchUserRow = (userId: number | string, patch: Partial<UserInfo>) => {
  const id = String(userId)
  const idx = userList.value.findIndex(item => String(item.id) === id)
  if (idx >= 0) {
    userList.value[idx] = {...userList.value[idx], ...patch}
  }
}

const toggleCancelStatus = async (row: UserInfo) => {
  if (row.cancel) {
    try {
      await ElMessageBox.confirm(
          t('pages.userList.uncancelConfirm', {id: row.id}),
          t('pages.userList.uncancelTitle'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
          }
      )

      const unCancelData: UnCancelReq = accountActionPayload(row)
      const response = await accountApi.unCancel(unCancelData)

      if (response) {
        patchUserRow(row.id, {cancel: false})
        ElMessage.success(t('pages.userList.uncancelSuccess'))
        await fetchUserList(true)
      } else {
        ElMessage.error(t('pages.userList.uncancelFailed'))
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('uncancel failed:', error)
      }
    }
  } else {
    try {
      await ElMessageBox.confirm(
          t('pages.userList.cancelConfirm', {id: row.id}),
          t('pages.userList.cancelTitle'),
          {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            type: 'warning'
          }
      )

      const cancelData: CancelReq = accountActionPayload(row)
      const response = await accountApi.cancel(cancelData)

      if (response) {
        patchUserRow(row.id, {cancel: true})
        ElMessage.success(t('pages.userList.cancelSuccess'))
        await fetchUserList(true)
      } else {
        ElMessage.error(t('pages.userList.cancelFailed'))
      }
    } catch (error) {
      if (error !== 'cancel') {
        console.error('cancel failed:', error)
      }
    }
  }
}

const formatUserType = (val: number | null | undefined) => {
  if (val === null || val === undefined) {
    return t('pages.userList.userTypeNormal')
  }
  return userTypeLabelMap.value[val] || t('pages.userList.userTypeNormal')
}

const formatVipLevel = (val: number | null | undefined) => {
  if (val === null || val === undefined || val <= 0) {
    return '-'
  }
  return String(val)
}

const formatDate = (dateString: string | null | undefined) => {
  if (!dateString) {
    return '-'
  }
  try {
    return new Date(dateString).toLocaleString(locale.value)
  } catch {
    return '-'
  }
}

onMounted(() => {
  fetchUserList()
})
</script>

<style scoped>
.page-container {
  padding: 20px;
  max-width: 100%;
  min-width: 0;
}

.page-container :deep(.el-card__body) {
  max-width: 100%;
  overflow-x: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.search-form {
  margin-bottom: 20px;
  width: 100%;
  min-width: 0;
}

.search-form :deep(.el-form-item__label) {
  white-space: nowrap;
}

.search-form :deep(.el-form--inline .el-form-item) {
  margin-right: 12px;
  margin-bottom: 8px;
}

.table-scroll {
  width: 100%;
  max-width: 100%;
  overflow-x: auto;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.game-picker-tip {
  margin: 0 0 16px;
  color: var(--el-text-color-secondary);
  line-height: 1.5;
}

.currency-gold,
.currency-diamond {
  font-variant-numeric: tabular-nums;
  font-weight: 500;
}

.currency-gold {
  color: #d48806;
}

.currency-diamond {
  color: #1677ff;
}

.id-cell {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  flex-wrap: nowrap;
  max-width: 100%;
}

.id-cell .el-button--link,
.id-cell .el-button.is-link {
  flex-shrink: 0;
}

.avatar-upload-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar-uploader :deep(.el-upload) {
  border: 1px dashed var(--el-border-color);
  border-radius: 50%;
  cursor: pointer;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader :deep(.el-upload:hover) {
  border-color: var(--el-color-primary);
}

.avatar-uploader-placeholder {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-uploader-icon {
  font-size: 24px;
  color: #8c939d;
}
</style>
