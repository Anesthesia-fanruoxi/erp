<template>
  <div class="user-manage">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="用户名/姓名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
          <el-button type="success" :icon="Plus" @click="handleAdd">新增用户</el-button>
          <el-button :icon="RefreshRight" @click="fetchList" title="刷新">刷新</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" style="margin-top: 16px">
      <el-table :data="tableData" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="userName" label="用户名" width="120" />
        <el-table-column prop="realName" label="真实姓名" width="120" />
        <el-table-column prop="email" label="邮箱" />
        <el-table-column prop="phone" label="手机" width="130" />
        <el-table-column label="角色" width="130">
          <template #default="{ row }">
            <el-tag v-if="row.roleName" size="small" type="info">{{ row.roleName }}</el-tag>
            <span v-else style="color: #c0c4cc; font-size: 12px">暂无角色</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90" align="center">
          <template #default="{ row }">
            <el-switch
              :model-value="row.status === 1"
              :loading="switchLoadingId === row.id"
              active-text=""
              inactive-text=""
              @change="(val: boolean) => handleToggleStatus(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="175">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="130" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该用户？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button type="danger" link size="small">删除</el-button>
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

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑用户' : '新增用户'"
      width="500px"
      @closed="resetForm"
    >
      <el-form ref="dialogFormRef" :model="dialogForm" :rules="dialogRules" label-width="80px">
        <el-form-item label="用户名" prop="userName">
          <el-input v-model="dialogForm.userName" :disabled="isEdit" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item v-if="!isEdit" label="密码" prop="password">
          <el-input v-model="dialogForm.password" type="password" show-password placeholder="请输入密码" />
        </el-form-item>
        <el-form-item v-if="isEdit" label="密码">
          <el-input v-model="dialogForm.password" type="password" show-password placeholder="不修改请留空" />
        </el-form-item>
        <el-form-item label="真实姓名" prop="realName">
          <el-input v-model="dialogForm.realName" placeholder="请输入真实姓名" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="dialogForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="手机">
          <el-input v-model="dialogForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select
            v-model="dialogForm.roleId"
            placeholder="请选择角色（可不选）"
            clearable
            style="width: 100%"
          >
            <el-option
              v-for="role in roleOptions"
              :key="role.id"
              :label="role.name"
              :value="role.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getUserList, createUser, updateUser, deleteUser, updateUserStatus } from '@/api/user'
import { getRoleList } from '@/api/role'
import { ElMessage } from 'element-plus'
import { Search, Refresh, Plus, RefreshRight } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/time'
import type { FormInstance, FormRules } from 'element-plus'

const loading = ref(false)
const submitLoading = ref(false)
const switchLoadingId = ref<number | null>(null)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const dialogFormRef = ref<FormInstance>()

const statusMap: Record<number, { label: string; type: 'success' | 'danger' }> = {
  1: { label: '正常', type: 'success' },
  2: { label: '禁用', type: 'danger' }
}

const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const tableData = ref<any[]>([])
const roleOptions = ref<{ id: number; name: string }[]>([])

const dialogForm = reactive({
  userName: '',
  password: '',
  realName: '',
  email: '',
  phone: '',
  roleId: undefined as number | undefined
})

const dialogRules: FormRules = {
  userName: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不少于6位', trigger: 'blur' }
  ],
  realName: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }]
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getUserList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      ...searchForm
    })
    tableData.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (err: any) {
    ElMessage.error(err.message || '获取列表失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  fetchList()
}

function handleReset() {
  searchForm.keyword = ''
  searchForm.status = undefined
  pagination.page = 1
  fetchList()
}

function handleAdd() {
  isEdit.value = false
  dialogVisible.value = true
  fetchRoles()
}

function handleEdit(row: any) {
  isEdit.value = true
  editId.value = row.id
  dialogForm.userName = row.userName
  dialogForm.realName = row.realName
  dialogForm.email = row.email || ''
  dialogForm.phone = row.phone || ''
  dialogForm.password = ''
  dialogForm.roleId = row.roleId ?? undefined
  dialogVisible.value = true
  fetchRoles()
}

async function handleSubmit() {
  if (!dialogFormRef.value) return
  await dialogFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (isEdit.value) {
        const payload: any = {
          realName: dialogForm.realName,
          email: dialogForm.email,
          phone: dialogForm.phone,
          // roleId 传 null 表示清除，传数字表示设置，不传表示不修改
          roleId: dialogForm.roleId ?? 0
        }
        if (dialogForm.password) payload.password = dialogForm.password
        await updateUser(editId.value, payload)
        ElMessage.success('更新成功')
      } else {
        const payload: any = {
          userName: dialogForm.userName,
          password: dialogForm.password,
          realName: dialogForm.realName,
          email: dialogForm.email,
          phone: dialogForm.phone,
        }
        if (dialogForm.roleId) payload.roleId = dialogForm.roleId
        await createUser(payload)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchList()
    } catch (err: any) {
      ElMessage.error(err.message || '操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

async function handleToggleStatus(row: any, val: boolean) {
  if (switchLoadingId.value !== null) return
  const newStatus = val ? 1 : 2
  switchLoadingId.value = row.id
  try {
    await updateUserStatus(row.id, newStatus)
    row.status = newStatus
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  } finally {
    switchLoadingId.value = null
  }
}

async function handleDelete(id: number) {
  try {
    await deleteUser(id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (err: any) {
    ElMessage.error(err.message || '删除失败')
  }
}

function resetForm() {
  dialogForm.userName = ''
  dialogForm.password = ''
  dialogForm.realName = ''
  dialogForm.email = ''
  dialogForm.phone = ''
  dialogForm.roleId = undefined
  dialogFormRef.value?.resetFields()
}

async function fetchRoles() {
  try {
    const res = await getRoleList()
    roleOptions.value = (res.data || []).map((r: any) => ({ id: r.id, name: r.name }))
  } catch {
    roleOptions.value = []
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped>
.user-manage {
  padding: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
}
.user-manage > :last-child { flex: 1; }
.search-card :deep(.el-form-item) { margin-bottom: 0; }
:deep(.el-table) { --el-table-header-bg-color: #f8f9fc; }
:deep(.el-table .el-table__header th) { font-size: 13px; color: #606266; font-weight: 600; }
:deep(.el-table .el-table__body td) { font-size: 13px; }
</style>
