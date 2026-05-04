<template>
  <div class="audit-manage">
    <el-card shadow="never" class="toolbar-card" style="margin-bottom: 16px">
      <el-button :icon="RefreshRight" @click="fetchList">刷新</el-button>
    </el-card>
    <el-card shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="userName" label="用户名" width="120" />
        <el-table-column prop="realName" label="真实姓名" width="120" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="phone" label="手机" width="130" />
        <el-table-column prop="createdAt" label="注册时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-popconfirm title="确定通过该用户的注册申请？" @confirm="handleApprove(row.id)">
              <template #reference>
                <el-button type="success" link size="small">通过</el-button>
              </template>
            </el-popconfirm>
            <el-popconfirm title="确定拒绝该用户的注册申请？" @confirm="handleReject(row.id)">
              <template #reference>
                <el-button type="danger" link size="small">拒绝</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 16px; justify-content: flex-end"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getAuditList, approveRegistration, rejectRegistration } from '@/api/audit'
import { ElMessage } from 'element-plus'
import { RefreshRight } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/time'

const loading = ref(false)
const tableData = ref<any[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

async function fetchList() {
  loading.value = true
  try {
    const res = await getAuditList({
      page: pagination.page,
      pageSize: pagination.pageSize
    })
    tableData.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (err: any) {
    ElMessage.error(err.message || '获取审核列表失败')
  } finally {
    loading.value = false
  }
}

async function handleApprove(id: number) {
  try {
    await approveRegistration(id)
    ElMessage.success('审核通过')
    fetchList()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

async function handleReject(id: number) {
  try {
    await rejectRegistration(id)
    ElMessage.success('已拒绝')
    fetchList()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.audit-manage {
  padding: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
}
.audit-manage > :last-child { flex: 1; }
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
</style>
