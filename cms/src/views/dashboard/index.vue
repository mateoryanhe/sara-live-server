<template>
  <div v-loading="loading" class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>系统总数据</span>
          <el-button :loading="loading" @click="fetchSysStat">刷新</el-button>
        </div>
      </template>

      <el-row :gutter="20">
        <el-col :lg="6" :md="12" :sm="24" :xs="24">
          <div class="stat-card stat-card-gold">
            <div class="stat-label">金币总额</div>
            <div class="stat-value">{{ formatAmount(sysStat.totalGold) }}</div>
          </div>
        </el-col>
        <el-col :lg="6" :md="12" :sm="24" :xs="24">
          <div class="stat-card stat-card-recharge">
            <div class="stat-label">总充值金额</div>
            <div class="stat-value">{{ formatAmount(sysStat.totalRecharge) }}</div>
          </div>
        </el-col>
        <el-col :lg="6" :md="12" :sm="24" :xs="24">
          <div class="stat-card stat-card-withdraw">
            <div class="stat-label">总提现金额</div>
            <div class="stat-value">{{ formatAmount(sysStat.totalWithdraw) }}</div>
          </div>
        </el-col>
        <el-col :lg="6" :md="12" :sm="24" :xs="24">
          <div class="stat-card stat-card-register">
            <div class="stat-label">总注册用户数</div>
            <div class="stat-value">{{ formatCount(sysStat.totalRegisterUser) }}</div>
          </div>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, reactive, ref} from 'vue'
import {ElMessage} from 'element-plus'
import {sysStatApi} from '@/api'
import type {SysStat} from '@/types/api'

const loading = ref(false)
const sysStat = reactive<SysStat>({
  totalGold: 0,
  totalRecharge: 0,
  totalWithdraw: 0,
  totalRegisterUser: 0,
})

const fetchSysStat = async () => {
  loading.value = true
  try {
    const data = await sysStatApi.getSysStat()
    sysStat.totalGold = data.totalGold ?? 0
    sysStat.totalRecharge = data.totalRecharge ?? 0
    sysStat.totalWithdraw = data.totalWithdraw ?? 0
    sysStat.totalRegisterUser = data.totalRegisterUser ?? 0
  } catch (error) {
    console.error('获取系统总数据失败:', error)
    ElMessage.error('获取系统总数据失败')
  } finally {
    loading.value = false
  }
}

const formatAmount = (value: number | null | undefined) => {
  if (value === null || value === undefined) {
    return '-'
  }
  return Number(value).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}

const formatCount = (value: string | number | null | undefined) => {
  if (value === null || value === undefined || value === '') {
    return '-'
  }
  return Number(value).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchSysStat()
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

.stat-card {
  border-radius: 12px;
  padding: 24px 20px;
  margin-bottom: 20px;
  color: #fff;
  min-height: 120px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.stat-label {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 12px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
  word-break: break-all;
}

.stat-card-gold {
  background: linear-gradient(135deg, #f6ad55, #ed8936);
}

.stat-card-recharge {
  background: linear-gradient(135deg, #63b3ed, #3182ce);
}

.stat-card-withdraw {
  background: linear-gradient(135deg, #68d391, #38a169);
}

.stat-card-register {
  background: linear-gradient(135deg, #b794f4, #805ad5);
}
</style>
