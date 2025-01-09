import request from '@/utils/request'

// 获取任务指标
export function getTaskMetrics(taskId: number) {
  return request({
    url: `/api/v1/dirscan/tasks/${taskId}/metrics`,
    method: 'get'
  })
}

// 获取任务性能报告
export function getTaskPerformance(taskId: number) {
  return request({
    url: `/api/v1/dirscan/tasks/${taskId}/performance`,
    method: 'get'
  })
}

// 获取目录树
export function getDirectoryTree(taskId: number) {
  return request({
    url: `/api/v1/dirscan/tasks/${taskId}/tree`,
    method: 'get'
  })
}

// 获取漏洞信息
export function getVulnerabilities(taskId: number) {
  return request({
    url: `/api/v1/dirscan/tasks/${taskId}/vulnerabilities`,
    method: 'get'
  })
}

// 导出报告
export function exportReport(taskId: number, format: string) {
  return request({
    url: `/api/v1/dirscan/tasks/${taskId}/report/${format}`,
    method: 'get',
    responseType: 'blob'
  })
} 