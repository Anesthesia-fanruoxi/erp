<template>
  <div class="contract-view">
    <!-- 工具栏 -->
    <el-card shadow="never" class="toolbar-card">
      <div class="toolbar">
        <div class="left">
          <el-input v-model="searchForm.keyword" placeholder="项目名称/公司" clearable
            style="width: 200px" @keyup.enter="handleSearch" />
          <el-input v-model="searchForm.orderNo" placeholder="订单号" clearable
            style="width: 160px" @keyup.enter="handleSearch" />
          <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </div>
        <div class="right">
          <el-button :icon="RefreshRight" @click="fetchList">刷新</el-button>
          <el-button type="success" :icon="Plus" @click="handleCreate">新建合同</el-button>
        </div>
      </div>
    </el-card>

    <!-- 合同列表 -->
    <el-card shadow="never" style="margin-top: 16px">
      <el-table :data="tableData" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="orderNo" label="订单号" width="140">
          <template #default="{ row }">
            <span class="order-no-text">{{ row.orderNo || '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="projectName" label="项目名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="fromCompany" label="发货方" min-width="140" show-overflow-tooltip />
        <el-table-column prop="toCompany" label="收货方" min-width="140" show-overflow-tooltip />
        <el-table-column prop="buyer" label="采购员" width="90" />
        <el-table-column prop="orderDate" label="订单日期" width="110" />
        <el-table-column prop="totalAmount" label="合计金额" width="110" align="right">
          <template #default="{ row }">
            <span class="amount-text">{{ row.totalAmount ? '¥ ' + row.totalAmount : '—' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="175">
          <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row.id)">编辑</el-button>
            <el-button type="info" link size="small" @click="handleView(row.id)">查看</el-button>
            <el-popconfirm title="确定删除该合同？" @confirm="handleDelete(row.id)">
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

    <!-- 新建/编辑弹窗 -->
    <el-dialog
      v-model="editVisible"
      :title="isEdit ? '编辑合同' : '新建合同'"
      width="900px"
      top="4vh"
      :close-on-click-modal="false"
      @closed="resetForm"
    >
      <!-- 弹窗内工具栏 -->
      <div class="dialog-toolbar">
        <el-button :icon="DocumentCopy" size="small" @click="showPasteDialog = true">
          粘贴导入
        </el-button>
      </div>

      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" style="margin-top: 12px">
        <!-- 基本信息 -->
        <div class="section-title">基本信息</div>
        <el-row :gutter="16">
          <el-col :span="24">
            <el-form-item label="项目名称" prop="projectName">
              <el-input v-model="form.projectName" placeholder="请输入项目名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="订单号">
              <el-input v-model="form.orderNo" placeholder="请输入订单号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="订单日期">
              <el-date-picker v-model="form.orderDate" type="date" placeholder="选择日期"
                value-format="YYYY-MM-DD" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="发货方">
              <el-input v-model="form.fromCompany" placeholder="请输入发货方" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="收货方">
              <el-input v-model="form.toCompany" placeholder="请输入收货方" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="采购员">
              <el-input v-model="form.buyer" placeholder="请输入采购员" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系人">
              <el-input v-model="form.attn" placeholder="请输入联系人" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="采购员电话">
              <el-input v-model="form.buyerTel" placeholder="请输入电话" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系人电话">
              <el-input v-model="form.attnTel" placeholder="请输入电话" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="收货地址">
              <el-input v-model="form.deliveryAddr" placeholder="请输入收货地址" />
            </el-form-item>
          </el-col>
          <el-col :span="24">
            <el-form-item label="备注">
              <el-input v-model="form.remark" type="textarea" :rows="2" placeholder="请输入备注" />
            </el-form-item>
          </el-col>
        </el-row>

        <!-- 采购明细 -->
        <div class="section-title" style="margin-top: 8px">
          采购明细
          <el-button type="primary" size="small" :icon="Plus" @click="handleAddItem">新增行</el-button>
        </div>
        <el-table :data="form.items" border size="small" style="width: 100%" max-height="300">
          <el-table-column label="序号" width="55" align="center">
            <template #default="{ $index }">{{ $index + 1 }}</template>
          </el-table-column>
          <el-table-column label="名称" min-width="120">
            <template #default="{ row }">
              <el-input v-model="row.name" size="small" placeholder="名称" />
            </template>
          </el-table-column>
          <el-table-column label="规格" min-width="110">
            <template #default="{ row }">
              <el-input v-model="row.spec" size="small" placeholder="规格" />
            </template>
          </el-table-column>
          <el-table-column label="品牌" width="90">
            <template #default="{ row }">
              <el-input v-model="row.brand" size="small" placeholder="品牌" />
            </template>
          </el-table-column>
          <el-table-column label="数量" width="80">
            <template #default="{ row }">
              <el-input v-model="row.qty" size="small" @input="calcAmount(row)" />
            </template>
          </el-table-column>
          <el-table-column label="单位" width="65">
            <template #default="{ row }">
              <el-input v-model="row.unit" size="small" />
            </template>
          </el-table-column>
          <el-table-column label="单价" width="90">
            <template #default="{ row }">
              <el-input v-model="row.unitPrice" size="small" @input="calcAmount(row)" />
            </template>
          </el-table-column>
          <el-table-column label="金额" width="90">
            <template #default="{ row }">
              <el-input v-model="row.amount" size="small" readonly />
            </template>
          </el-table-column>
          <el-table-column label="下单人" width="75">
            <template #default="{ row }">
              <el-input v-model="row.operator" size="small" />
            </template>
          </el-table-column>
          <el-table-column label="位置" width="80">
            <template #default="{ row }">
              <el-input v-model="row.location" size="small" />
            </template>
          </el-table-column>
          <el-table-column label="备注" width="90">
            <template #default="{ row }">
              <el-input v-model="row.remark" size="small" />
            </template>
          </el-table-column>
          <el-table-column label="" width="50" fixed="right" align="center">
            <template #default="{ $index }">
              <el-button type="danger" link size="small" @click="handleRemoveItem($index)">删</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div class="total-row">
          合计金额：<span class="total-amount">¥ {{ totalAmount }}</span>
        </div>
      </el-form>

      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">保存</el-button>
      </template>
    </el-dialog>

    <!-- 粘贴导入弹窗 -->
    <el-dialog v-model="showPasteDialog" title="粘贴 Excel 内容" width="500px" append-to-body>
      <div class="paste-tip">
        <el-icon color="#667eea"><InfoFilled /></el-icon>
        在 Excel 中全选（Ctrl+A）→ 复制（Ctrl+C）→ 点击下方区域后粘贴（Ctrl+V），自动解析填充
      </div>
      <div
        class="paste-zone"
        tabindex="0"
        @paste.prevent="handlePasteEvent"
        @click="focusPasteZone"
        ref="pasteZoneRef"
      >
        <div class="paste-zone-placeholder">
          <el-icon :size="36" color="#c0c4cc"><DocumentCopy /></el-icon>
          <p>点击此处激活，然后按 Ctrl+V 粘贴</p>
          <p style="font-size: 12px; color: #c0c4cc">粘贴后自动解析并关闭弹窗</p>
        </div>
      </div>
      <template #footer>
        <el-button @click="showPasteDialog = false; rawHtml = ''">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 查看详情抽屉 -->
    <el-drawer v-model="viewVisible" title="合同详情" size="75%" direction="rtl">
      <template #header>
        <div style="display:flex;align-items:center;justify-content:space-between;width:100%">
          <span style="font-size:16px;font-weight:600">合同详情</span>
          <el-button type="warning" :icon="Printer" @click="handleOpenDelivery" style="margin-right:16px">
            导出出库单
          </el-button>
        </div>
      </template>
      <div v-if="viewData" class="view-wrap">
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="项目名称" :span="2">{{ viewData.projectName }}</el-descriptions-item>
          <el-descriptions-item label="订单号">
            <span class="order-no-text">{{ viewData.orderNo || '—' }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="订单日期">{{ viewData.orderDate || '—' }}</el-descriptions-item>
          <el-descriptions-item label="发货方">{{ viewData.fromCompany || '—' }}</el-descriptions-item>
          <el-descriptions-item label="收货方">{{ viewData.toCompany || '—' }}</el-descriptions-item>
          <el-descriptions-item label="采购员">{{ viewData.buyer || '—' }}</el-descriptions-item>
          <el-descriptions-item label="联系人">{{ viewData.attn || '—' }}</el-descriptions-item>
          <el-descriptions-item label="采购员电话">{{ viewData.buyerTel || '—' }}</el-descriptions-item>
          <el-descriptions-item label="联系人电话">{{ viewData.attnTel || '—' }}</el-descriptions-item>
          <el-descriptions-item label="合计金额">
            <span class="amount-text">{{ viewData.totalAmount ? '¥ ' + viewData.totalAmount : '—' }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="收货地址" :span="2">{{ viewData.deliveryAddr || '—' }}</el-descriptions-item>
          <el-descriptions-item label="备注" :span="2">{{ viewData.remark || '—' }}</el-descriptions-item>
        </el-descriptions>

        <div class="view-section-title" style="margin-top: 20px">
          采购明细
          <span class="item-count">共 {{ viewData.items?.length || 0 }} 项</span>
        </div>
        <el-table :data="viewData.items || []" border size="small" style="width: 100%">
          <el-table-column prop="seq" label="序号" width="55" align="center" />
          <el-table-column prop="name" label="名称" min-width="120" show-overflow-tooltip />
          <el-table-column prop="spec" label="规格" min-width="110" show-overflow-tooltip />
          <el-table-column prop="brand" label="品牌" width="90" />
          <el-table-column prop="qty" label="数量" width="75" align="right" />
          <el-table-column prop="unit" label="单位" width="60" align="center" />
          <el-table-column prop="unitPrice" label="单价" width="90" align="right" />
          <el-table-column prop="amount" label="金额" width="90" align="right">
            <template #default="{ row }">
              <span class="amount-text">{{ row.amount }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="operator" label="下单人" width="75" align="center" />
          <el-table-column prop="location" label="安装位置" width="80" align="center" />
          <el-table-column prop="remark" label="备注" min-width="80" show-overflow-tooltip />
        </el-table>
      </div>
    </el-drawer>

    <!-- 导出出库单：勾选产品弹窗 -->
    <el-dialog v-model="deliveryVisible" title="选择导出产品" width="700px" append-to-body>
      <div class="delivery-tip">勾选需要导出到出库单的产品，未勾选的行将留空</div>
      <el-table
        ref="deliveryTableRef"
        :data="viewData?.items || []"
        border size="small"
        style="width:100%"
        @selection-change="handleDeliverySelection"
      >
        <el-table-column type="selection" width="46" />
        <el-table-column prop="seq" label="序号" width="55" align="center" />
        <el-table-column prop="name" label="产品名称" min-width="130" show-overflow-tooltip />
        <el-table-column prop="spec" label="规格" min-width="110" show-overflow-tooltip />
        <el-table-column prop="brand" label="特征" width="90" />
        <el-table-column prop="qty" label="数量" width="75" align="right" />
        <el-table-column prop="unit" label="单位" width="60" align="center" />
      </el-table>
      <div style="margin-top:10px;font-size:13px;color:#909399">
        已选 <strong>{{ selectedDeliveryItems.length }}</strong> 项
      </div>
      <template #footer>
        <el-button @click="deliveryVisible = false">取消</el-button>
        <el-button type="warning" :icon="Printer" :disabled="!selectedDeliveryItems.length" @click="handlePrintDelivery">
          打印出库单
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import {
  getContractList, getContractDetail, createContract, updateContract, deleteContract,
  type ContractItem, type ContractForm
} from '@/api/contract'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Refresh, RefreshRight, Plus, DocumentCopy, InfoFilled, Printer } from '@element-plus/icons-vue'
import { formatTime } from '@/utils/time'
import { printDelivery } from '@/utils/printDelivery'

// ---- 列表 ----
const loading = ref(false)
const tableData = ref<any[]>([])
const searchForm = reactive({ keyword: '', orderNo: '' })
const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

async function fetchList() {
  loading.value = true
  try {
    const res = await getContractList({
      page: pagination.page, pageSize: pagination.pageSize,
      keyword: searchForm.keyword || undefined,
      orderNo: searchForm.orderNo || undefined
    })
    tableData.value = res.data?.list || []
    pagination.total = res.data?.total || 0
  } catch (err: any) {
    ElMessage.error(err.message || '获取列表失败')
  } finally {
    loading.value = false
  }
}

function handleSearch() { pagination.page = 1; fetchList() }
function handleReset() { searchForm.keyword = ''; searchForm.orderNo = ''; pagination.page = 1; fetchList() }

// ---- 新建/编辑弹窗 ----
const editVisible = ref(false)
const isEdit = ref(false)
const editId = ref<number | null>(null)
const submitLoading = ref(false)
const formRef = ref<FormInstance>()
const showPasteDialog = ref(false)
const rawText = ref('')
const rawHtml = ref('')
const previewRows = ref(0)
const pasteZoneRef = ref<HTMLElement>()

const emptyForm = (): ContractForm => ({
  projectName: '', orderNo: '', orderDate: '', fromCompany: '', toCompany: '',
  buyer: '', attn: '', buyerEmail: '', buyerTel: '', attnTel: '',
  totalAmount: '', deliveryAddr: '', remark: '', items: []
})

const form = reactive<ContractForm>(emptyForm())

const rules: FormRules = {
  projectName: [{ required: true, message: '请输入项目名称', trigger: 'blur' }]
}

const totalAmount = computed(() => {
  const sum = form.items.reduce((acc, item) => {
    return acc + (parseFloat(item.amount.replace(/,/g, '')) || 0)
  }, 0)
  form.totalAmount = sum.toFixed(2)
  return sum.toFixed(2)
})

function calcAmount(row: ContractItem) {
  const qty = parseFloat(row.qty.replace(/,/g, '')) || 0
  const price = parseFloat(row.unitPrice.replace(/,/g, '')) || 0
  row.amount = (qty * price).toFixed(2)
}

function handleAddItem() {
  form.items.push({
    seq: form.items.length + 1, name: '', spec: '', brand: '',
    qty: '', unit: '', unitPrice: '', amount: '', operator: '', location: '', remark: ''
  })
}

function handleRemoveItem(idx: number) {
  form.items.splice(idx, 1)
  form.items.forEach((item, i) => { item.seq = i + 1 })
}

function resetForm() {
  Object.assign(form, emptyForm())
  rawHtml.value = ''
  formRef.value?.resetFields()
}

function handleCreate() {
  isEdit.value = false
  editId.value = null
  resetForm()
  editVisible.value = true
}

async function handleEdit(id: number) {
  isEdit.value = true
  editId.value = id
  resetForm()
  try {
    const res = await getContractDetail(id)
    const d = res.data
    Object.assign(form, {
      projectName: d.projectName, orderNo: d.orderNo, orderDate: d.orderDate,
      fromCompany: d.fromCompany, toCompany: d.toCompany, buyer: d.buyer,
      attn: d.attn, buyerEmail: d.buyerEmail, buyerTel: d.buyerTel,
      attnTel: d.attnTel, totalAmount: d.totalAmount,
      deliveryAddr: d.deliveryAddr, remark: d.remark,
      items: d.items || []
    })
    editVisible.value = true
  } catch (err: any) {
    ElMessage.error(err.message || '加载失败')
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (isEdit.value && editId.value) {
        await updateContract(editId.value, form)
        ElMessage.success('更新成功')
      } else {
        await createContract(form)
        ElMessage.success('创建成功')
      }
      editVisible.value = false
      fetchList()
    } catch (err: any) {
      ElMessage.error(err.message || '操作失败')
    } finally {
      submitLoading.value = false
    }
  })
}

// ---- 粘贴解析 ----
function focusPasteZone() {
  pasteZoneRef.value?.focus()
}

function handlePasteEvent(e: ClipboardEvent) {
  const html = e.clipboardData?.getData('text/html') || ''
  const text = e.clipboardData?.getData('text/plain') || ''

  const hasData = form.projectName || form.items.length > 0

  const doParse = () => {
    try {
      let result
      if (html) {
        rawHtml.value = html
        previewRows.value = (html.match(/<tr/gi) || []).length
        result = parseFromHtml(html)
      } else {
        rawHtml.value = '__TEXT__' + text
        previewRows.value = text.split('\n').filter(l => l.trim()).length
        result = parseFromText(text)
      }
      Object.assign(form, result.header)
      form.items = result.items
      showPasteDialog.value = false
      rawHtml.value = ''
      ElMessage.success(`解析成功，共 ${result.items.length} 条明细`)
    } catch (e: any) {
      ElMessage.error('解析失败：' + (e.message || '格式不匹配'))
    }
  }

  if (hasData) {
    ElMessageBox.confirm('当前合同已有数据，粘贴将覆盖所有内容，是否继续？', '提示', {
      confirmButtonText: '覆盖',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(doParse).catch(() => {})
  } else {
    doParse()
  }
}

// 从 HTML 解析（主路径）
function parseFromHtml(html: string): { header: any; items: ContractItem[] } {
  const parser = new DOMParser()
  const doc = parser.parseFromString(html, 'text/html')
  const rows = Array.from(doc.querySelectorAll('tr'))
  console.log('[HTML解析] 总行数:', rows.length)

  const h: any = {}
  const itemList: ContractItem[] = []
  let inTable = false

  // 提取每行的单元格文本（去除空白）
  const getRowCells = (row: Element): string[] =>
    Array.from(row.querySelectorAll('td,th')).map(td => td.textContent?.trim() || '')

  for (const row of rows) {
    const cols = getRowCells(row)
    // 过滤全空行
    if (cols.every(c => !c)) continue
    console.log('[HTML行]', cols)

    // 找第一个非空列作为 key
    const firstNonEmpty = cols.find(c => c) || ''

    if (firstNonEmpty.includes('项目名称')) {
      h.projectName = cols.filter(c => c && !c.includes('项目名称')).join(' ').trim()
      console.log('[匹配] 项目名称:', h.projectName)
      continue
    }
    if (firstNonEmpty.toLowerCase().startsWith('from')) {
      h.fromCompany = cols.filter(c => c && !c.toLowerCase().startsWith('from') && !c.toLowerCase().startsWith('to')).find(Boolean) || ''
      const toIdx = cols.findIndex(c => c.toLowerCase() === 'to:' || c.toLowerCase() === 'to')
      if (toIdx !== -1) h.toCompany = cols.slice(toIdx + 1).find(c => c) || ''
      console.log('[匹配] From:', h.fromCompany, 'To:', h.toCompany)
      continue
    }
    if (firstNonEmpty.toLowerCase().startsWith('buyer')) {
      h.buyer = cols.filter(c => c && !c.toLowerCase().startsWith('buyer') && !c.toLowerCase().startsWith('attn')).find(Boolean) || ''
      const attnIdx = cols.findIndex(c => c.toLowerCase().startsWith('attn'))
      if (attnIdx !== -1) h.attn = cols.slice(attnIdx + 1).find(c => c) || ''
      console.log('[匹配] Buyer:', h.buyer, 'Attn:', h.attn)
      continue
    }
    if (firstNonEmpty.toLowerCase().startsWith('tel')) {
      const nonLabel = cols.filter(c => c && !c.toLowerCase().startsWith('tel') && !c.toLowerCase().startsWith('fax') && !c.toLowerCase().startsWith('e-mail'))
      h.buyerTel = nonLabel[0] || ''
      h.attnTel = nonLabel[1] || ''
      console.log('[匹配] Tel:', h.buyerTel, 'AttnTel:', h.attnTel)
      continue
    }
    if (firstNonEmpty.toLowerCase().includes('order date') || firstNonEmpty.includes('订单日期')) {
      h.orderDate = cols.filter(c => c && c !== firstNonEmpty).find(Boolean) || ''
      console.log('[匹配] OrderDate:', h.orderDate)
      continue
    }
    if (firstNonEmpty.toLowerCase().includes('order no') || firstNonEmpty.includes('订单号')) {
      h.orderNo = cols.filter(c => c && c !== firstNonEmpty).find(Boolean) || ''
      console.log('[匹配] OrderNo:', h.orderNo)
      continue
    }
    // 收货地址（格式：3.收货地址：xxx 或 收货地址：xxx）
    if (firstNonEmpty.includes('收货地址')) {
      let val = cols.filter(c => c && !c.includes('收货地址')).join('').trim()
        || firstNonEmpty.replace(/^[\d.]*收货地址[：:]\s*/, '').trim()
      console.log('[收货地址] 截断前原始值:', JSON.stringify(val))
      console.log('[收货地址] cols:', JSON.stringify(cols))
      console.log('[收货地址] firstNonEmpty:', JSON.stringify(firstNonEmpty))
      // 修复：使用负向后瞻确保数字序列不是其他数字的一部分
      // (?<!\d) 确保前面没有数字，\d{1,3} 匹配1-3位数字（章节号）
      // 这样可以避免匹配手机号中的数字（如"199185094664"中的"664"）
      const sectionPattern = /(?<!\d)\d{1,3}[.．]\s*[\u4e00-\u9fa5]/
      const matchIndex = val.search(sectionPattern)
      if (matchIndex > 0) {
        val = val.substring(0, matchIndex).trim()
      }
      console.log('[收货地址] 截断后:', JSON.stringify(val))
      h.deliveryAddr = val
      continue
    }
    // 明细表头
    if (!inTable && (firstNonEmpty === '序号' || cols.some(c => c.includes('材料') || c.includes('设备名称')))) {
      inTable = true
      console.log('[匹配] 进入明细表格')
      continue
    }
    // 明细数据行
    if (inTable && /^\d+$/.test(firstNonEmpty)) {
      // 找到序号列的位置
      const seqIdx = cols.findIndex(c => /^\d+$/.test(c))
      const get = (offset: number) => cols[seqIdx + offset] || ''
      const item: ContractItem = {
        seq: parseInt(firstNonEmpty),
        name: get(1), spec: get(2), brand: get(3),
        qty: get(4), unit: get(5), unitPrice: get(6), amount: get(7),
        operator: get(8), location: get(9), remark: get(10)
      }
      console.log('[明细行]', item)
      itemList.push(item)
    }
  }

  console.log('[解析结果] header:', h, 'items:', itemList.length)
  return { header: h, items: itemList }
}

// 纯文本降级解析
function parseFromText(text: string): { header: any; items: ContractItem[] } {
  const lines = text.split('\n').map(l => l.trim()).filter(l => l)
  const h: any = {}
  const itemList: ContractItem[] = []
  const cells = (line: string) => line.split('\t').map(c => c.trim())
  let inTable = false

  for (const line of lines) {
    const cols = cells(line)
    const nonEmpty = cols.filter(c => c)
    if (!nonEmpty.length) continue

    if (cols[0].includes('项目名称')) { h.projectName = nonEmpty.slice(1).join(' ').trim(); continue }
    if (cols[0].toLowerCase().startsWith('from')) {
      h.fromCompany = nonEmpty[1] || ''
      const toIdx = nonEmpty.findIndex(c => c.toLowerCase().startsWith('to'))
      if (toIdx !== -1) h.toCompany = nonEmpty[toIdx + 1] || ''
      continue
    }
    if (cols[0].toLowerCase().startsWith('buyer')) {
      h.buyer = nonEmpty[1] || ''
      const attnIdx = nonEmpty.findIndex(c => c.toLowerCase().startsWith('attn'))
      if (attnIdx !== -1) h.attn = nonEmpty[attnIdx + 1] || ''
      continue
    }
    if (cols[0].toLowerCase().startsWith('tel')) {
      const vals = nonEmpty.filter(c => !c.toLowerCase().startsWith('tel'))
      h.buyerTel = vals[0] || ''; h.attnTel = vals[1] || ''
      continue
    }
    if (cols[0].toLowerCase().includes('order date') || cols[0].includes('订单日期')) {
      h.orderDate = nonEmpty[1] || ''; continue
    }
    if (cols[0].toLowerCase().includes('order no') || cols[0].includes('订单号')) {
      h.orderNo = nonEmpty[1] || ''; continue
    }
    if (cols[0].includes('收货地址')) {
      let val = nonEmpty.filter(c => !c.includes('收货地址')).join('').trim()
        || cols[0].replace(/^[\d.]*收货地址[：:]\s*/, '').trim()
      // 修复：使用负向后瞻确保数字序列不是其他数字的一部分
      // (?<!\d) 确保前面没有数字，\d{1,3} 匹配1-3位数字（章节号）
      // 这样可以避免匹配手机号中的数字（如"199185094664"中的"664"）
      const sectionPattern = /(?<!\d)\d{1,3}[.．]\s*[\u4e00-\u9fa5]/
      const matchIndex = val.search(sectionPattern)
      if (matchIndex > 0) {
        val = val.substring(0, matchIndex).trim()
      }
      h.deliveryAddr = val
      continue
    }
    if (!inTable && (cols[0] === '序号' || cols.some(c => c.includes('材料')))) {
      inTable = true; continue
    }
    if (inTable && /^\d+$/.test(cols[0])) {
      const get = (idx: number) => cols[idx] || ''
      itemList.push({
        seq: parseInt(cols[0]), name: get(1), spec: get(2), brand: get(3),
        qty: get(4), unit: get(5), unitPrice: get(6), amount: get(7),
        operator: get(8), location: get(9), remark: get(10)
      })
    }
  }
  return { header: h, items: itemList }
}

// ---- 查看详情 ----
const viewVisible = ref(false)
const viewData = ref<any>(null)

// ---- 导出出库单 ----
const deliveryVisible = ref(false)
const deliveryTableRef = ref()
const selectedDeliveryItems = ref<any[]>([])

function handleDeliverySelection(rows: any[]) {
  selectedDeliveryItems.value = rows
}

function handleOpenDelivery() {
  selectedDeliveryItems.value = []
  deliveryVisible.value = true
}

function handlePrintDelivery() {
  if (!viewData.value) return
  const today = new Date()
  const printDate = `${today.getFullYear()}-${String(today.getMonth()+1).padStart(2,'0')}-${String(today.getDate()).padStart(2,'0')}`

  printDelivery({
    customerName: viewData.value.fromCompany || '',
    contractNo: viewData.value.orderNo || '',
    projectName: viewData.value.projectName || '',
    deliveryNo: '',
    deliveryAddr: viewData.value.deliveryAddr || '',
    printDate,
    items: selectedDeliveryItems.value.map((item, idx) => ({
      seq: idx + 1,
      name: item.name || '',
      spec: item.spec || '',
      brand: item.brand || '',
      qty: item.qty || '',
      unit: item.unit || '',
    }))
  })
  deliveryVisible.value = false
}

async function handleView(id: number) {
  try {
    const res = await getContractDetail(id)
    viewData.value = res.data
    viewVisible.value = true
  } catch (err: any) {
    ElMessage.error(err.message || '获取详情失败')
  }
}

async function handleDelete(id: number) {
  try {
    await deleteContract(id)
    ElMessage.success('删除成功')
    fetchList()
  } catch (err: any) {
    ElMessage.error(err.message || '删除失败')
  }
}

onMounted(fetchList)
</script>

<style scoped>
.contract-view { padding: 0; }
.toolbar-card :deep(.el-card__body) { padding: 14px 18px; }
.toolbar { display: flex; align-items: center; justify-content: space-between; }
.left, .right { display: flex; align-items: center; gap: 8px; }

.order-no-text {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  color: #667eea;
}
.amount-text { color: #f56c6c; font-weight: 500; }

:deep(.el-table) { --el-table-header-bg-color: #f8f9fc; }
:deep(.el-table .el-table__header th) { font-size: 13px; color: #606266; font-weight: 600; }

/* 弹窗内 */
.dialog-toolbar {
  display: flex;
  justify-content: flex-end;
  padding-bottom: 4px;
  border-bottom: 1px solid #f0f0f5;
}
.section-title {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #f0f0f5;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.total-row {
  margin-top: 10px;
  text-align: right;
  font-size: 13px;
  color: #606266;
}
.total-amount {
  font-size: 16px;
  font-weight: 700;
  color: #f56c6c;
  margin-left: 6px;
}

/* 粘贴弹窗 */
.paste-tip {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #606266;
  background: #f0f4ff;
  border-radius: 8px;
  padding: 10px 14px;
  margin-bottom: 12px;
}
.paste-zone {
  border: 2px dashed #d0d3e0;
  border-radius: 10px;
  height: 140px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  outline: none;
  transition: border-color 0.2s, background 0.2s;
}
.paste-zone:focus, .paste-zone:hover {
  border-color: #667eea;
  background: #f5f7ff;
}
.paste-zone-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  color: #909399;
  font-size: 14px;
}
.paste-zone-preview {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #67c23a;
  font-weight: 500;
}
.paste-area :deep(.el-textarea__inner) {
  font-family: 'JetBrains Mono', Consolas, monospace;
  font-size: 12px;
  line-height: 1.6;
}

/* 查看抽屉 */
.view-wrap { padding: 0 4px; }
.view-section-title {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.item-count { font-size: 12px; color: #909399; font-weight: normal; }

/* 导出出库单 */
.delivery-tip {
  font-size: 13px;
  color: #909399;
  margin-bottom: 10px;
  padding: 8px 12px;
  background: #fdf6ec;
  border-radius: 6px;
  border-left: 3px solid #e6a23c;
}
</style>
