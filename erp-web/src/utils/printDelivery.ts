/**
 * 打印出库单（送货单）
 * 使用浏览器 window.print() 实现，样式完全还原模板
 */

export interface DeliveryItem {
  seq: number
  name: string
  spec: string
  brand: string
  qty: string
  unit: string
  remark1?: string
  remark2?: string
  remark3?: string
}

export interface DeliveryData {
  customerName: string   // 客户名称
  contractNo: string     // 合同号
  projectName: string    // 项目名称
  deliveryNo: string     // 送货单号
  deliveryAddr: string   // 送货地址
  printDate: string      // 打印日期
  items: DeliveryItem[]
}

export function printDelivery(data: DeliveryData) {
  // 确保至少有6行（不足补空行）
  const rows = [...data.items]
  while (rows.length < 6) {
    rows.push({ seq: rows.length + 1, name: '', spec: '', brand: '', qty: '', unit: '' })
  }

  const rowsHtml = rows.map(r => `
    <tr>
      <td class="center red">${r.seq}</td>
      <td>${r.name || ''}</td>
      <td>${r.spec || ''}</td>
      <td>${r.brand || ''}</td>
      <td class="center">${r.qty || ''}</td>
      <td class="center">${r.unit || ''}</td>
      <td>${r.remark1 || ''}</td>
      <td>${r.remark2 || ''}</td>
      <td>${r.remark3 || ''}</td>
    </tr>
  `).join('')

  const html = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <title>送货单</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body {
      font-family: '宋体', SimSun, serif;
      font-size: 12px;
      color: #c00;
      padding: 20px 30px;
      background: #fff;
    }
    .title-wrap {
      text-align: center;
      margin-bottom: 6px;
    }
    .company-name {
      font-size: 22px;
      font-weight: bold;
      color: #c00;
      letter-spacing: 2px;
    }
    .doc-title {
      font-size: 18px;
      font-weight: bold;
      color: #c00;
      margin-top: 4px;
    }
    .info-section {
      display: flex;
      justify-content: space-between;
      margin: 10px 0 8px;
      font-size: 12px;
      color: #c00;
    }
    .info-left, .info-right {
      display: flex;
      flex-direction: column;
      gap: 4px;
    }
    .info-right { text-align: left; }
    .info-row { display: flex; gap: 4px; }
    .info-label { white-space: nowrap; }
    .info-value {
      border-bottom: 1px solid #c00;
      min-width: 160px;
      flex: 1;
    }
    .info-value-sm { min-width: 120px; }

    table {
      width: 100%;
      border-collapse: collapse;
      margin-top: 4px;
    }
    th, td {
      border: 1px solid #c00;
      padding: 5px 4px;
      color: #c00;
      font-size: 12px;
    }
    th {
      background: #fff;
      font-weight: bold;
      text-align: center;
    }
    td { height: 28px; }
    .center { text-align: center; }
    .red { color: #c00; }

    .footer {
      display: flex;
      justify-content: space-between;
      margin-top: 20px;
      font-size: 12px;
      color: #c00;
    }
    .footer-block { display: flex; flex-direction: column; gap: 8px; }
    .footer-line { display: flex; gap: 4px; align-items: flex-end; }
    .footer-underline {
      border-bottom: 1px solid #c00;
      min-width: 140px;
      height: 18px;
    }

    @media print {
      body { padding: 10px 20px; }
      @page { margin: 10mm; size: A4; }
    }
  </style>
</head>
<body>
  <div class="title-wrap">
    <div class="company-name">湖南盛宏新材料有限公司</div>
    <div class="doc-title">送货单</div>
  </div>

  <div class="info-section">
    <div class="info-left">
      <div class="info-row">
        <span class="info-label">客户名称:</span>
        <span class="info-value">${data.customerName}</span>
      </div>
      <div class="info-row">
        <span class="info-label">项目名称:</span>
        <span class="info-value">${data.projectName}</span>
      </div>
      <div class="info-row">
        <span class="info-label">送货地址:</span>
        <span class="info-value">${data.deliveryAddr}</span>
      </div>
    </div>
    <div class="info-right">
      <div class="info-row">
        <span class="info-label">合同号：</span>
        <span class="info-value info-value-sm">${data.contractNo}</span>
      </div>
      <div class="info-row">
        <span class="info-label">送货单号:</span>
        <span class="info-value info-value-sm">${data.deliveryNo}</span>
      </div>
      <div class="info-row">
        <span class="info-label">打印日期:</span>
        <span class="info-value info-value-sm">${data.printDate}</span>
      </div>
    </div>
  </div>

  <table>
    <thead>
      <tr>
        <th style="width:40px">序号</th>
        <th style="width:18%">产品名称</th>
        <th style="width:16%">规格</th>
        <th style="width:10%">特征</th>
        <th style="width:7%">数量</th>
        <th style="width:6%">单位</th>
        <th style="width:10%">备注1</th>
        <th style="width:10%">备注2</th>
        <th style="width:10%">备注3</th>
      </tr>
    </thead>
    <tbody>
      ${rowsHtml}
    </tbody>
  </table>

  <div class="footer">
    <div class="footer-block">
      <div class="footer-line">
        <span>送货单位签字或盖章：</span>
        <span class="footer-underline"></span>
      </div>
      <div class="footer-line">
        <span>日期:</span>
        <span class="footer-underline"></span>
      </div>
    </div>
    <div class="footer-block">
      <div class="footer-line">
        <span>收货单位签字或盖章：</span>
        <span class="footer-underline"></span>
      </div>
      <div class="footer-line">
        <span>日期：</span>
        <span class="footer-underline"></span>
      </div>
    </div>
  </div>

  <script>
    window.onload = function() { window.print(); }
  <\/script>
</body>
</html>`

  const win = window.open('', '_blank', 'width=900,height=700')
  if (!win) {
    alert('请允许弹出窗口以打印送货单')
    return
  }
  win.document.write(html)
  win.document.close()
}
