<template>
  <div class="role-manage">
    <el-card shadow="never" class="search-card">
      <el-button type="success" :icon="Plus" @click="handleAdd">新增角色</el-button>
    </el-card>

    <el-card shadow="never" style="margin-top: 16px">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="角色名" width="150" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="权限数" width="100">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ (row.permissionIds || []).length }} 项</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除该角色？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button type="danger" link size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑角色' : '新增角色'"
      width="820px"
      top="6vh"
      @closed="resetForm"
    >
      <el-form ref="dialogFormRef" :model="dialogForm" :rules="dialogRules" label-width="80px">
        <el-form-item label="角色名" prop="name">
          <el-input v-model="dialogForm.name" placeholder="请输入角色名" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="dialogForm.description" type="textarea" :rows="2" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="权限分配">
          <div class="perm-panel">
            <div class="perm-toolbar">
              <el-button size="small" @click="handleSelectAll">全选</el-button>
              <el-button size="small" @click="handleClearAll">清空</el-button>
              <span class="perm-count">已选 {{ selectedIds.size }} 项</span>
            </div>

            <div class="perm-groups">
              <div
                v-for="top in permissionGroups"
                :key="top.id"
                class="perm-group"
              >
                <div class="group-header">
                  <el-checkbox
                    :model-value="isChecked(top.id)"
                    :indeterminate="isGroupIndeterminate(top)"
                    @change="(v: boolean) => toggleGroup(top, v)"
                  >
                    <span class="group-name">{{ top.name }}</span>
                    <el-tag size="small" type="info" class="group-tag">菜单</el-tag>
                  </el-checkbox>
                </div>

                <div class="group-body">
                  <div
                    v-for="sub in top.children || []"
                    :key="sub.id"
                    class="perm-row"
                  >
                    <div class="perm-row-label">
                      <el-checkbox
                        :model-value="isChecked(sub.id)"
                        :indeterminate="isSubIndeterminate(sub)"
                        @change="(v: boolean) => toggleSub(sub, v)"
                      >
                        <span class="sub-name">{{ sub.name }}</span>
                        <el-tag size="small" type="success" class="sub-tag">显示</el-tag>
                      </el-checkbox>
                    </div>
                    <div class="perm-row-actions">
                      <el-checkbox
                        v-for="btn in sub.children || []"
                        :key="btn.id"
                        :model-value="isChecked(btn.id)"
                        @change="(v: boolean) => toggleButton(sub, btn, v)"
                      >
                        <span class="btn-name">{{ btn.name }}</span>
                      </el-checkbox>
                      <span
                        v-if="!(sub.children || []).length"
                        class="no-action-tip"
                      >暂无操作按钮</span>
                    </div>
                  </div>
                  <div v-if="!(top.children || []).length" class="no-sub-tip">
                    该菜单下无子菜单
                  </div>
                </div>
              </div>
              <div v-if="!permissionGroups.length" class="empty-tip">
                暂无权限数据
              </div>
            </div>
          </div>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { getRoleList, createRole, updateRole, deleteRole, getPermissionList } from '@/api/role'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/time'
import type { FormInstance, FormRules } from 'element-plus'

interface Permission {
  id: number
  code: string
  name: string
  type: number
  parentId: number
  path: string
  icon: string
  sort: number
  visible: number
  children?: Permission[]
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const dialogFormRef = ref<FormInstance>()

const tableData = ref<any[]>([])
const allPermissions = ref<Permission[]>([])
const selectedIds = ref<Set<number>>(new Set())

const dialogForm = reactive({
  name: '',
  description: ''
})

const dialogRules: FormRules = {
  name: [{ required: true, message: '请输入角色名', trigger: 'blur' }]
}

const permissionGroups = computed<Permission[]>(() => {
  const list = allPermissions.value
  const byParent = new Map<number, Permission[]>()
  for (const p of list) {
    const arr = byParent.get(p.parentId) || []
    arr.push({ ...p })
    byParent.set(p.parentId, arr)
  }
  function build(parentId: number): Permission[] {
    const children = byParent.get(parentId) || []
    return children
      .sort((a, b) => a.sort - b.sort || a.id - b.id)
      .map((c) => ({ ...c, children: build(c.id) }))
  }
  return build(0)
})

function isChecked(id: number): boolean {
  return selectedIds.value.has(id)
}

function collectIds(node: Permission): number[] {
  const ids: number[] = [node.id]
  for (const c of node.children || []) {
    ids.push(...collectIds(c))
  }
  return ids
}

function isGroupIndeterminate(top: Permission): boolean {
  const all = collectIds(top)
  const selectedCount = all.filter((id) => selectedIds.value.has(id)).length
  return selectedCount > 0 && selectedCount < all.length
}

function isSubIndeterminate(sub: Permission): boolean {
  const all = collectIds(sub)
  const selectedCount = all.filter((id) => selectedIds.value.has(id)).length
  return selectedCount > 0 && selectedCount < all.length
}

function setBatch(ids: number[], checked: boolean) {
  const next = new Set(selectedIds.value)
  ids.forEach((id) => {
    if (checked) next.add(id)
    else next.delete(id)
  })
  selectedIds.value = next
}

function toggleGroup(top: Permission, checked: boolean) {
  setBatch(collectIds(top), checked)
}

function toggleSub(sub: Permission, checked: boolean) {
  const ids = collectIds(sub)
  if (checked) {
    const parent = allPermissions.value.find((p) => p.id === sub.parentId)
    if (parent) ids.push(parent.id)
  }
  setBatch(ids, checked)
}

function toggleButton(sub: Permission, btn: Permission, checked: boolean) {
  const ids = [btn.id]
  if (checked) {
    ids.push(sub.id)
    const top = allPermissions.value.find((p) => p.id === sub.parentId)
    if (top) ids.push(top.id)
  }
  setBatch(ids, checked)
}

function handleSelectAll() {
  const all: number[] = []
  permissionGroups.value.forEach((top) => all.push(...collectIds(top)))
  setBatch(all, true)
}

function handleClearAll() {
  selectedIds.value = new Set()
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getRoleList()
    tableData.value = res.data?.list || res.data || []
  } catch (err: any) {
    ElMessage.error(err.message || '获取角色列表失败')
  } finally {
    loading.value = false
  }
}

async function fetchPermissions() {
  try {
    const res = await getPermissionList()
    allPermissions.value = res.data || []
  } catch (err: any) {
    ElMessage.error(err.message || '获取权限失败')
  }
}

function handleAdd() {
  isEdit.value = false
  dialogVisible.value = true
  selectedIds.value = new Set()
}

function handleEdit(row: any) {
  isEdit.value = true
  editId.value = row.id
  dialogForm.name = row.name
  dialogForm.description = row.description || ''
  selectedIds.value = new Set<number>(row.permissionIds || [])
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!dialogFormRef.value) return
  await dialogFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const data = {
        name: dialogForm.name,
        description: dialogForm.description,
        permissionIds: Array.from(selectedIds.value)
      }
      if (isEdit.value) {
        await updateRole(editId.value, data)
        ElMessage.success('更新成功')
      } else {
        await createRole(data)
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

async function handleDelete(id: number) {
  try {
    await deleteRole(id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (err: any) {
    ElMessage.error(err.message || '删除失败')
  }
}

function resetForm() {
  dialogForm.name = ''
  dialogForm.description = ''
  selectedIds.value = new Set()
  dialogFormRef.value?.resetFields()
}

onMounted(() => {
  fetchList()
  fetchPermissions()
})
</script>

<style scoped>
.role-manage {
  padding: 0;
}
.perm-panel {
  width: 100%;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  background: #fafbfd;
  padding: 12px;
}
.perm-toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding-bottom: 10px;
  border-bottom: 1px dashed #e4e7ed;
  margin-bottom: 12px;
}
.perm-count {
  margin-left: auto;
  font-size: 13px;
  color: #909399;
}
.perm-groups {
  max-height: 440px;
  overflow-y: auto;
  padding-right: 4px;
}
.perm-group {
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  margin-bottom: 10px;
  overflow: hidden;
}
.group-header {
  background: #f5f7fa;
  padding: 10px 14px;
  border-bottom: 1px solid #ebeef5;
}
.group-name {
  font-weight: 600;
  color: #303133;
  margin-right: 8px;
}
.group-tag,
.sub-tag {
  margin-left: 6px;
}
.group-body {
  padding: 6px 14px;
}
.perm-row {
  display: flex;
  align-items: flex-start;
  padding: 10px 0;
  border-bottom: 1px dashed #f0f0f5;
  gap: 12px;
}
.perm-row:last-child {
  border-bottom: none;
}
.perm-row-label {
  width: 180px;
  flex-shrink: 0;
}
.sub-name {
  color: #303133;
  font-weight: 500;
}
.perm-row-actions {
  flex: 1;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px 18px;
}
.btn-name {
  color: #606266;
  font-size: 13px;
}
.no-action-tip,
.no-sub-tip,
.empty-tip {
  color: #c0c4cc;
  font-size: 12px;
  font-style: italic;
}
.no-sub-tip {
  padding: 10px 0;
}
.empty-tip {
  text-align: center;
  padding: 24px 0;
}
</style>
