<template>
  <div class="menu-manage">
    <!-- 工具栏 -->
    <el-card shadow="never" class="toolbar-card">
      <div class="toolbar">
        <div class="left">
          <el-button type="success" :icon="Plus" @click="handleAdd(0)">新增根菜单</el-button>
          <el-button :icon="Refresh" @click="loadData">刷新</el-button>
          <el-button :icon="Expand" @click="toggleExpandAll">{{ expandAll ? '全部折叠' : '全部展开' }}</el-button>
        </div>
        <div class="right">
          <el-input
            v-model="keyword"
            placeholder="按名称/编码筛选"
            clearable
            style="width: 240px"
            :prefix-icon="Search"
          />
        </div>
      </div>
    </el-card>

    <!-- 树形表格 -->
    <el-card shadow="never" style="margin-top: 16px">
      <el-table
        :data="filteredTree"
        v-loading="loading"
        row-key="id"
        :tree-props="{ children: 'children' }"
        :expand-row-keys="expandKeys"
        default-expand-all
        stripe
        style="width: 100%"
      >
        <el-table-column prop="name" label="名称" min-width="200">
          <template #default="{ row }">
            <el-icon v-if="row.icon" style="vertical-align: middle; margin-right: 6px">
              <component :is="resolveIcon(row.icon)" />
            </el-icon>
            <span>{{ row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="code" label="权限编码" min-width="180">
          <template #default="{ row }">
            <code class="code-tag">{{ row.code }}</code>
          </template>
        </el-table-column>
        <el-table-column label="类型" width="80">
          <template #default="{ row }">
            <el-tag :type="typeTagType(row.type)" size="small">{{ typeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="path" label="路由路径" min-width="150" />
        <el-table-column prop="sort" label="排序" width="70" align="center" />
        <el-table-column label="可见" width="70" align="center">
          <template #default="{ row }">
            <el-tag v-if="row.visible === 1" type="success" size="small">显示</el-tag>
            <el-tag v-else type="info" size="small">隐藏</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <div class="action-cell">
              <!-- 第一列：详情 或 新增子项，固定宽度占位 -->
              <span class="action-slot">
                <el-button
                  v-if="row.type === 2"
                  type="info" link size="small"
                  @click="handleDetail(row)"
                >详情</el-button>
                <el-button
                  v-else-if="row.type === 1"
                  type="primary" link size="small"
                  @click="handleAdd(row.id)"
                >新增子项</el-button>
              </span>
              <!-- 第二列：编辑 -->
              <span class="action-slot">
                <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
              </span>
              <!-- 第三列：删除 -->
              <span class="action-slot">
                <el-popconfirm title="确定删除该菜单项？" @confirm="handleDelete(row.id)">
                  <template #reference>
                    <el-button type="danger" link size="small">删除</el-button>
                  </template>
                </el-popconfirm>
              </span>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 目录详情抽屉 -->
    <el-drawer
      v-model="detailVisible"
      :title="detailDir?.name + ' · 子菜单'"
      direction="rtl"
      size="480px"
    >
      <div v-if="detailDir" class="detail-wrap">
        <!-- 目录基本信息 -->
        <div class="detail-dir-info">
          <div class="dir-info-row">
            <el-icon class="dir-info-icon"><component :is="resolveIcon(detailDir.icon) || 'FolderOpened'" /></el-icon>
            <div class="dir-info-main">
              <div class="dir-info-name">{{ detailDir.name }}</div>
              <code class="code-tag">{{ detailDir.code }}</code>
            </div>
            <el-tag type="info" size="small">目录</el-tag>
          </div>
        </div>

        <!-- 子菜单列表 -->
        <div class="detail-section-title">
          <span>子菜单</span>
          <el-button type="success" size="small" :icon="Plus" @click="handleAdd(detailDir.id)">新增</el-button>
        </div>

        <div v-if="detailChildren.length === 0" class="detail-empty">
          <el-empty description="暂无子菜单" :image-size="80" />
        </div>

        <div v-else class="detail-menu-list">
          <div
            v-for="child in detailChildren"
            :key="child.id"
            class="detail-menu-card"
          >
            <div class="menu-card-left">
              <div class="menu-card-icon">
                <el-icon :size="18"><component :is="resolveIcon(child.icon) || 'Document'" /></el-icon>
              </div>
              <div class="menu-card-info">
                <div class="menu-card-name">{{ child.name }}</div>
                <div class="menu-card-meta">
                  <code class="code-tag">{{ child.code }}</code>
                  <span v-if="child.path" class="path-tag">{{ child.path }}</span>
                </div>
              </div>
            </div>
            <div class="menu-card-right">
              <el-tag
                :type="child.visible === 1 ? 'success' : 'info'"
                size="small"
                style="margin-right: 8px"
              >{{ child.visible === 1 ? '显示' : '隐藏' }}</el-tag>
              <el-button type="primary" link size="small" @click="handleEdit(child)">编辑</el-button>
              <el-popconfirm title="确定删除？" @confirm="handleDelete(child.id)">
                <template #reference>
                  <el-button type="danger" link size="small">删除</el-button>
                </template>
              </el-popconfirm>
            </div>
          </div>
        </div>
      </div>
    </el-drawer>

    <!-- 新增/编辑弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="560px"
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio :value="1">菜单</el-radio>
            <el-radio :value="2">目录</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="父级" prop="parentId">
          <el-tree-select
            v-model="form.parentId"
            :data="parentOptions"
            :props="{ label: 'name', value: 'id' }"
            node-key="id"
            check-strictly
            clearable
            placeholder="顶级菜单"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="如: 菜单管理" />
        </el-form-item>
        <el-form-item label="权限编码" prop="code">
          <el-input v-model="form.code" placeholder="如: system:menu" />
        </el-form-item>
        <el-form-item v-if="form.type === 1" label="路由路径" prop="path">
          <el-input v-model="form.path" placeholder="如: /system/menu" />
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" placeholder="Element Plus 图标名, 如: Menu" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.sort" :min="0" :max="999" />
        </el-form-item>
        <el-form-item v-if="form.type === 1" label="是否显示">
          <el-switch v-model="form.visible" :active-value="1" :inactive-value="0" />
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
import { ref, computed, onMounted, reactive } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh, Search, Expand } from '@element-plus/icons-vue'
import * as ElIcons from '@element-plus/icons-vue'
import { getMenuList, createMenu, updateMenu, deleteMenu, type MenuItem } from '@/api/menu'

const loading = ref(false)
const submitLoading = ref(false)
const flatList = ref<MenuItem[]>([])
const keyword = ref('')
const expandAll = ref(true)
const expandKeys = ref<number[]>([])

// 详情抽屉
const detailVisible = ref(false)
const detailDir = ref<MenuItem | null>(null)
const detailChildren = computed(() =>
  detailDir.value
    ? flatList.value
        .filter(m => m.parentId === detailDir.value!.id)
        .sort((a, b) => a.sort - b.sort || a.id - b.id)
    : []
)

const dialogVisible = ref(false)
const dialogTitle = ref('新增菜单')
const isEdit = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const form = reactive<Partial<MenuItem>>({
  type: 1,
  parentId: 0,
  name: '',
  code: '',
  path: '',
  icon: '',
  sort: 0,
  visible: 1,
})

const rules: FormRules = {
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入权限编码', trigger: 'blur' }],
}

function buildTree(list: MenuItem[]): MenuItem[] {
  const map = new Map<number, MenuItem>()
  const roots: MenuItem[] = []
  list.forEach(item => map.set(item.id, { ...item, children: [] }))
  map.forEach(node => {
    if (node.parentId && map.has(node.parentId)) {
      map.get(node.parentId)!.children!.push(node)
    } else {
      roots.push(node)
    }
  })
  return roots
}

const tree = computed(() => buildTree(flatList.value))

const filteredTree = computed(() => {
  if (!keyword.value.trim()) return tree.value
  const kw = keyword.value.toLowerCase()
  const filter = (nodes: MenuItem[]): MenuItem[] => {
    const result: MenuItem[] = []
    for (const n of nodes) {
      const matched = n.name.toLowerCase().includes(kw) || n.code.toLowerCase().includes(kw)
      const children = n.children ? filter(n.children) : []
      if (matched || children.length > 0) {
        result.push({ ...n, children })
      }
    }
    return result
  }
  return filter(tree.value)
})

const parentOptions = computed(() => {
  const filterMenu = (nodes: MenuItem[]): MenuItem[] => {
    return nodes
      .filter(n => n.type === 2)
      .map(n => ({ ...n, children: n.children ? filterMenu(n.children) : [] }))
  }
  return [{ id: 0, name: '顶级目录', children: filterMenu(tree.value) } as any]
})

function typeLabel(t: number) { return t === 1 ? '菜单' : '目录' }
function typeTagType(t: number): 'primary' | 'success' {
  return t === 1 ? 'primary' : 'success'
}
function resolveIcon(name: string) {
  return name ? (ElIcons as any)[name] || null : null
}

function handleDetail(row: MenuItem) {
  detailDir.value = row
  detailVisible.value = true
}

async function loadData() {
  loading.value = true
  try {
    const res: any = await getMenuList()
    flatList.value = res.data || []
    expandKeys.value = flatList.value.map(i => i.id)
    // 如果抽屉打开，同步刷新
    if (detailDir.value) {
      const updated = flatList.value.find(m => m.id === detailDir.value!.id)
      if (updated) detailDir.value = updated
    }
  } finally {
    loading.value = false
  }
}

function toggleExpandAll() {
  expandAll.value = !expandAll.value
  expandKeys.value = expandAll.value ? flatList.value.map(i => i.id) : []
}

function resetForm() {
  Object.assign(form, { type: 1, parentId: 0, name: '', code: '', path: '', icon: '', sort: 0, visible: 1 })
  formRef.value?.clearValidate()
}

function handleAdd(parentId: number) {
  isEdit.value = false
  editingId.value = null
  resetForm()
  form.parentId = parentId
  const parent = flatList.value.find(i => i.id === parentId)
  form.type = parent?.type === 2 ? 1 : 1
  dialogTitle.value = '新增菜单'
  dialogVisible.value = true
}

function handleEdit(row: MenuItem) {
  isEdit.value = true
  editingId.value = row.id
  resetForm()
  Object.assign(form, {
    type: row.type, parentId: row.parentId, name: row.name,
    code: row.code, path: row.path, icon: row.icon, sort: row.sort, visible: row.visible,
  })
  dialogTitle.value = '编辑菜单'
  dialogVisible.value = true
}

async function handleDelete(id: number) {
  try {
    await deleteMenu(id)
    ElMessage.success('删除成功')
    await loadData()
  } catch (e: any) {
    ElMessage.error(e?.message || '删除失败')
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      const payload = { ...form, parentId: form.parentId || 0 }
      if (isEdit.value && editingId.value != null) {
        await updateMenu(editingId.value, payload)
        ElMessage.success('更新成功')
      } else {
        await createMenu(payload)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      await loadData()
    } catch (e: any) {
      ElMessage.error(e?.message || '操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

onMounted(loadData)
</script>

<style scoped>
.menu-manage {
  padding: 0;
  min-height: calc(100vh - 64px - 40px);
  display: flex;
  flex-direction: column;
}
.menu-manage > :last-child { flex: 1; }

.toolbar-card :deep(.el-card__body) { padding: 14px 18px; }
.toolbar { display: flex; justify-content: space-between; align-items: center; }

.code-tag {
  display: inline-block;
  padding: 2px 8px;
  background: #f5f7fa;
  border-radius: 4px;
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #606266;
}

/* 详情抽屉 */
.detail-wrap { padding: 0 4px; }

.detail-dir-info {
  background: linear-gradient(135deg, #f0f4ff, #f8fafc);
  border-radius: 12px;
  padding: 16px;
  margin-bottom: 20px;
  border: 1px solid #e8eeff;
}

.dir-info-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.dir-info-icon {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 20px;
  flex-shrink: 0;
}

.dir-info-main { flex: 1; }
.dir-info-name {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a2e;
  margin-bottom: 4px;
}

.detail-section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f5;
}

.detail-empty { padding: 20px 0; }

.detail-menu-list { display: flex; flex-direction: column; gap: 10px; }

.detail-menu-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 10px;
  transition: all 0.2s ease;
}

.detail-menu-card:hover {
  border-color: #c0ccff;
  box-shadow: 0 2px 8px rgba(102, 126, 234, 0.1);
}

.menu-card-left { display: flex; align-items: center; gap: 12px; flex: 1; min-width: 0; }

.menu-card-icon {
  width: 36px;
  height: 36px;
  background: #f0f4ff;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #667eea;
  flex-shrink: 0;
}

.menu-card-info { min-width: 0; }
.menu-card-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}
.menu-card-meta { display: flex; align-items: center; gap: 8px; flex-wrap: wrap; }

.path-tag {
  font-size: 11px;
  color: #909399;
  font-family: 'JetBrains Mono', Consolas, monospace;
}

.menu-card-right { display: flex; align-items: center; flex-shrink: 0; }

/* 操作列对齐 */
.action-cell {
  display: flex;
  align-items: center;
}
.action-slot {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 64px;
}
</style>
