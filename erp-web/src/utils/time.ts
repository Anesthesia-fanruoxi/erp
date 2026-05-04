/**
 * 时间工具函数
 * 所有时间传递使用毫秒时间戳，显示使用中国标准时间 (UTC+8)
 */

/**
 * 毫秒时间戳格式化为中国标准时间字符串
 * @param ms 毫秒时间戳
 * @param format 格式: 'datetime' | 'date' | 'time'
 * @returns 格式化后的时间字符串，如 "2026-04-30 14:30:00"
 */
export function formatTime(ms: number, format: 'datetime' | 'date' | 'time' = 'datetime'): string {
  if (!ms || ms === 0) return '--'

  const date = new Date(ms)

  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  const h = String(date.getHours()).padStart(2, '0')
  const min = String(date.getMinutes()).padStart(2, '0')
  const s = String(date.getSeconds()).padStart(2, '0')

  if (format === 'date') return `${y}-${m}-${d}`
  if (format === 'time') return `${h}:${min}:${s}`
  return `${y}-${m}-${d} ${h}:${min}:${s}`
}
