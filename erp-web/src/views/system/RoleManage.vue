<template>
  <div class="role-manage">
    <el-card shadow="never" class="search-card">
      <el-button type="success" :icon="Plus" @click="handleAdd">新增角色</el-button>
      <el-button :icon="RefreshRight" @click="fetchList" style="margin-left: 8px">刷新</el-button>
    </el-card>

    <el-card shadow="never" style="margin-top: 16px">
      <el-table :data="tableData" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="角色名" width="150" />
        <el-table-column prop="description" label="描述" />
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
      width="760px"
      top="6vh"
      @closed="resetForm"
    >
      <el-form ref="dialogFormRef" :model="dialogForm" :rules="dialogRules" label-width="80px">
        <el-form-item label="角色名" prop="name">
          <el-input v-model="dialogForm.name" placeholder="请输入角色名" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="dialogForm.description" type="textarea" :rows="2" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="权限分配">
          <div class="perm-panel">
            <!-- 工具栏 -->
            <div class="perm-toolbar">
              <el-button size="small" @click="handleSelectAll">全选</el-button>
              <el-button size="small" @click="handleClearAll">清空</el-button>
              <div class="perm-legend">
                <span class="legend-item"><el-tag size="small" type="info">查看</el-tag> 菜单可见</span>
                <span class="legend-item"><el-tag size="small" type="success">只读</el-tag> GET接口</span>
                <span class="legend-item"><el-tag size="small" type="warning">操作</el-tag> 写接口</span>
              </div>
            </div>

            <!-- 权限表格 -->
            <div class="perm-table-wrap">
              <!-- 表头 -->
              <div class="perm-row perm-header">
                <div class="perm-col-name">菜单</div>
                <div class="perm-col-check">查看</div>
                <div class="perm-col-check">只读</div>
                <div class="perm-col-check">操作</div>
              </div>

              <!-- 目录分组 -->
              <template v-for="dir in menuTree" :key="dir.id">
                <!-- 目录行 -->
                <div class="perm-row perm-dir-row">
                  <div class="perm-col-name">
                    <el-icon v-if="dir.icon" class="dir-icon"><component :is="dir.icon" /></el-icon>
                    <span class="dir-name">{{ dir.name }}</span>
                  </div>
                  <div class="perm-col-check">
                    <el-checkbox
                      :model-value="isDirChecked(dir, 1)"
                      :indeterminate="isDirIndeterminate(dir, 1)"
                      @change="(v: boolean) => toggleDir(dir, 1, v)"
                    />
                  </div>
                  <div class="perm-col-check">
                    <el-checkbox
                      :model-value="isDirChecked(dir, 2)"
                      :indeterminate="isDirIndeterminate(dir, 2)"
                      @change="(v: boolean) => toggleDir(dir, 2, v)"
                    />
                  </div>
                  <div class="perm-col-check">
                    <el-checkbox
                      :model-value="isDirChecked(dir, 3)"
                      :indeterminate="isDirIndeterminate(dir, 3)"
                      @change="(v: boolean) => toggleDir(dir, 3, v)"
                    />
                  </div>
                </div>

                <!-- 菜单行 -->
                <div
                  v-for="menu in dir.children || []"
                  :key="menu.id"
                  class="perm-row perm-menu-row"
                >
                  <div class="perm-col-name">
                    <span class="menu-indent">└</span>
                    <span class="menu-name">{{ menu.name }}</span>
                  </div>
                  <!-- 查看 (perm_type=1) -->
                  <div class="perm-col-check">
                    <el-checkbox
                      :model-value="hasPermType(menu.id, 1)"
                      @change="(v: boolean) => toggleMenuPerm(menu.id, 1, v)"
                    />
                  </div>
                  <!-- 只读 (perm_type=2) -->
                  <div class="perm-col-check">
                    <el-checkbox
                      :model-value="hasPermType(menu.id, 2)"
                      @change="(v: boolean) => toggleMenuPerm(menu.id, 2, v)"
                    />
                  </div>
                  <!-- 操作 (perm_type=3) -->
                  <div class="perm-col-check">
                    <el-checkbox
                      :model-value="hasPermType(menu.id, 3)"
                      @change="(v: boolean) => toggleMenuPerm(menu.id, 3, v)"
                    />
                  </div>
                </div>
              </template>

              <div v-if="!menuTree.length" class="perm-empty">暂无菜单数据</div>
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
import { getRoleList, createRole, updateRole, deleteRole } from '@/api/role'
import { getMenuList, type MenuItem } from '@/api/menu'
import { ElMessage } from 'element-plus'
import { Plus, RefreshRight } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/time'
import type { FormInstance, FormRules } from 'element-plus'

// menuId -> Set<permType>
// permType: 1-查看 2-只读 3-操作
const permMap = ref<Map<number, Set<number>>>(new Map())

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const dialogFormRef = ref<FormInstance>()
const tableData = ref<any[]>([])
const flatMenus = ref<MenuItem[]>([])

const dialogForm = reactive({ name: '', description: '' })
const dialogRules: FormRules = {
  name: [{ required: true, message: '请输入角色名', trigger: 'blur' }]
}

// 构建目录树: 顶级目录(type=2) -> 子菜单(type=1)
const menuTree = computed(() => {
  const dirs = flatMenus.value.filter(m => m.type === 2 && m.parentId === 0)
    .sort((a, b) => a.sort - b.sort || a.id - b.id)
  return dirs.map(dir => ({
    ...dir,
    children: flatMenus.value
      .filter(m => m.parentId === dir.id)
      .sort((a, b) => a.sort - b.sort || a.id - b.id)
  }))
})

// 判断某菜单是否拥有某权限类型
function hasPermType(menuId: number, permType: number): boolean {
  return permMap.value.get(menuId)?.has(permType) ?? false
}

// 切换单个菜单的某权限类型（三列完全独立，无联动）
function toggleMenuPerm(menuId: number, permType: number, checked: boolean) {
  const next = new Map(permMap.value)
  const set = new Set(next.get(menuId) ?? [])
  if (checked) {
    set.add(permType)
  } else {
    set.delete(permType)
  }
  if (set.size === 0) next.delete(menuId)
  else next.set(menuId, set)
  permMap.value = next
}

// 目录行：某列是否全选
function isDirChecked(dir: any, permType: number): boolean {
  const children: MenuItem[] = dir.children || []
  if (!children.length) return false
  return children.every(m => hasPermType(m.id, permType))
}

// 目录行：某列是否半选
function isDirIndeterminate(dir: any, permType: number): boolean {
  const children: MenuItem[] = dir.children || []
  if (!children.length) return false
  const count = children.filter(m => hasPermType(m.id, permType)).length
  return count > 0 && count < children.length
}

// 切换目录某列（全选/全取消）
function toggleDir(dir: any, permType: number, checked: boolean) {
  const children: MenuItem[] = dir.children || []
  children.forEach(m => toggleMenuPerm(m.id, permType, checked))
}

function handleSelectAll() {
  flatMenus.value.filter(m => m.type === 1).forEach(m => {
    const set = new Set([1, 2, 3])
    const next = new Map(permMap.value)
    next.set(m.id, set)
    permMap.value = next
  })
}

function handleClearAll() {
  permMap.value = new Map()
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

async function fetchMenus() {
  try {
    const res: any = await getMenuList()
    flatMenus.value = res.data || []
  } catch (err: any) {
    ElMessage.error(err.message || '获取菜单失败')
  }
}

function handleAdd() {
  isEdit.value = false
  permMap.value = new Map()
  dialogVisible.value = true
}

function handleEdit(row: any) {
  isEdit.value = true
  editId.value = row.id
  dialogForm.name = row.name
  dialogForm.description = row.description || ''

  // 还原 menuPerms -> permMap
  const next = new Map<number, Set<number>>()
  for (const mp of (row.menuPerms || [])) {
    next.set(mp.menuId, new Set(mp.permTypes))
  }
  permMap.value = next
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!dialogFormRef.value) return
  await dialogFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      // 将 permMap 转为 menuPerms 数组
      const menuPerms: { menuId: number; permTypes: number[] }[] = []
      permMap.value.forEach((types, menuId) => {
        if (types.size > 0) {
          menuPerms.push({ menuId, permTypes: Array.from(types) })
        }
      })

      const data = {
        name: dialogForm.name,
        description: dialogForm.description,
        menuPerms
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
  permMap.value = new Map()
  dialogFormRef.value?.resetFields()
}

onMounted(() => {
  fetchList()
  fetchMenus()
})
</script>

<style scoped>
.role-manage {
  padding: 0;
  flex: 1;
  display: flex;
  flex-direction: column;
}
.role-manage > :last-child { flex: 1; }

:deep(.el-table) { --el-table-header-bg-color: #f8f9fc; }
:deep(.el-table .el-table__header th) { font-size: 13px; color: #606266; font-weight: 600; }
:deep(.el-table .el-table__body td) { font-size: 13px; }

/* 权限面板 */
.perm-panel {
  width: 100%;
  border: 1px solid #ebeef5;
  border-radius: 10px;
  background: #fafbfd;
  overflow: hidden;
}

.perm-toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border-bottom: 1px dashed #e4e7ed;
  background: #f5f7fa;
}

.perm-legend {
  margin-left: auto;
  display: flex;
  gap: 14px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 12px;
  color: #909399;
}

.perm-table-wrap {
  max-height: 420px;
  overflow-y: auto;
}

/* 行通用 */
.perm-row {
  display: grid;
  grid-template-columns: 1fr 80px 80px 80px;
  align-items: center;
  border-bottom: 1px solid #f0f0f5;
}
.perm-row:last-child { border-bottom: none; }

/* 表头 */
.perm-header {
  background: #f5f7fa;
  font-size: 13px;
  font-weight: 600;
  color: #606266;
  position: sticky;
  top: 0;
  z-index: 1;
}
.perm-header .perm-col-name { padding: 10px 14px; }
.perm-header .perm-col-check { text-align: center; padding: 10px 0; }

/* 目录行 */
.perm-dir-row {
  background: linear-gradient(90deg, #f0f4ff, #f8fafc);
  padding: 8px 0;
}
.perm-dir-row .perm-col-name {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 14px;
  font-weight: 600;
  color: #303133;
}
.dir-icon { color: #667eea; font-size: 16px; }

/* 菜单行 */
.perm-menu-row {
  background: #fff;
  padding: 6px 0;
}
.perm-menu-row:hover { background: #fafbff; }
.perm-menu-row .perm-col-name {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0 14px;
  color: #606266;
  font-size: 13px;
}
.menu-indent { color: #c0c4cc; font-size: 14px; }
.menu-name { color: #303133; }

/* checkbox 列 */
.perm-col-check {
  display: flex;
  justify-content: center;
  align-items: center;
}

.perm-empty {
  text-align: center;
  padding: 30px 0;
  color: #c0c4cc;
  font-size: 13px;
}
</style>
