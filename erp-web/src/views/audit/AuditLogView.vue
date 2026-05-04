<template>
  <div class="audit-log-view">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="操作人">
          <el-input
            v-model="searchForm.userName"
            placeholder="用户名"
            clearable
            style="width: 140px"
          />
        </el-form-item>
        <el-form-item label="操作描述">
          <el-input
            v-model="searchForm.action"
            placeholder="如: 创建用户"
            clearable
            style="width: 160px"
          />
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="timeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="x"
            style="width: 360px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
          <el-button :icon="RefreshRight" @click="fetchList">刷新</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" style="margin-top: 16px">
      <el-table :data="tableData" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="userName" label="操作人" width="110" />
        <el-table-column prop="action" label="操作描述" width="130">
          <template #default="{ row }">
            <el-tag size="small" type="primary">{{ row.action }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="请求" min-width="200">
          <template #default="{ row }">
            <el-tag
              :type="methodTagType(row.method)"
              size="small"
              style="margin-right: 6px"
            >{{ row.method }}</el-tag>
            <span class="path-text">{{ row.path }}</span>
          </template>
        </el-table-column>
        <el-table-column label="响应码" width="90" align="center">
          <template #default="{ row }">
            <el-tag
              :type="row.statusCode === 200 ? 'success' : 'danger'"
              size="small"
            >{{ row.statusCode }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ip" label="IP" width="130" />
        <el-table-column prop="duration" label="耗时" width="90" align="right">
          <template #default="{ row }">
            <span :class="durationClass(row.duration)">{{ row.duration }}ms</span>
          </template>
        </el-table-column>
        <el-table-column label="操作时间" width="175">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="详情" width="80" align="center" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleDetail(row)">
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 16px; justify-content: flex-end"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </el-card>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="操作详情" width="640px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="操作人">{{ detail.userName }}</el-descriptions-item>
        <el-descriptions-item label="操作描述">{{ detail.action }}</el-descriptions-item>
        <el-descriptions-item label="请求方法">{{ detail.method }}</el-descriptions-item>
        <el-descriptions-item label="请求路径">{{ detail.path }}</el-descriptions-item>
        <el-descriptions-item label="响应码">{{ detail.statusCode }}</el-descriptions-item>
        <el-descriptions-item label="耗时">{{ detail.duration }}ms</el-descriptions-item>
        <el-descriptions-item label="来源IP">{{ detail.ip }}</el-descriptions-item>
        <el-descriptions-item label="操作时间">{{ formatTime(detail.createdAt) }}</el-descriptions-item>
        <el-descriptions-item v-if="detail.query" label="查询参数" :span="2">
          <code class="detail-code">{{ detail.query }}</code>
        </el-descriptions-item>
        <el-descriptions-item v-if="detail.body" label="请求体" :span="2">
          <pre class="detail-pre">{{ formatJson(detail.body) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { onMounted } from 'vue'
import { getAuditLogList } from '@/api/audit'
import { ElMessage } from 'element-plus'
import { Search, Refresh, RefreshRight } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/time'

const loading = ref(false)
const tableData = ref<AuditLog[]>([])
const detailVisible = ref(false)
const detail = ref<AuditLog>({} as AuditLog)

// 时间范围: [startMs, endMs]
const timeRange = ref<[number, number] | null>(null)

const searchForm = reactive({
  userName: '',
  action: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

async function fetchList() {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      pageSize: pagination.pageSize,
      userName: searchForm.userName || undefined,
      action: searchForm.action || undefined
    }
    if (timeRange.value) {
      params.startTime = timeRange.value[0]
      params.endTime = timeRange.value[1]
    }
    const res = await getAuditLogList(params)
    tableData.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (err: any) {
    ElMessage.error(err.message || '获取日志失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  fetchList()
}

function handleReset() {
  searchForm.userName = ''
  searchForm.action = ''
  timeRange.value = null
  pagination.page = 1
  fetchList()
}

function handleDetail(row: AuditLog) {
  detail.value = row
  detailVisible.value = true
}

// HTTP 方法颜色
function methodTagType(method: string): 'success' | 'warning' | 'danger' | 'info' | 'primary' {
  const map: Record<string, 'success' | 'warning' | 'danger' | 'info' | 'primary'> = {
    GET: 'info',
    POST: 'success',
    PUT: 'warning',
    DELETE: 'danger',
    PATCH: 'primary'
  }
  return map[method] || 'info'
}

// 耗时颜色
function durationClass(ms: number): string {
  if (ms < 200) return 'duration-fast'
  if (ms < 1000) return 'duration-normal'
  return 'duration-slow'
}

// 格式化 JSON 字符串
function formatJson(str: string): string {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.audit-log-view {
  padding: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
}
.audit-log-view > :last-child { flex: 1; }
.search-card :deep(.el-form-item) {
  margin-bottom: 0;
}
:deep(.el-table) {
  --el-table-header-bg-color: #f8f9fc;
}
:deep(.el-table .el-table__header th) {
  font-size: 13px;
  color: #606266;
  font-weight: 600;
}
:deep(.el-table .el-table__body td) {
  font-size: 13px;
}
.path-text {
  font-size: 12px;
  color: #606266;
  font-family: 'JetBrains Mono', Consolas, monospace;
}
.duration-fast {
  color: #67c23a;
  font-weight: 500;
}
.duration-normal {
  color: #e6a23c;
  font-weight: 500;
}
.duration-slow {
  color: #f56c6c;
  font-weight: 600;
}
.detail-code {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #606266;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
}
.detail-pre {
  margin: 0;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #303133;
  background: #f5f7fa;
  padding: 10px;
  border-radius: 6px;
  white-space: pre-wrap;
  word-break: break-all;
  max-height: 300px;
  overflow-y: auto;
}
</style>
