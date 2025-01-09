<template>
  <div class="dir-scan-monitor">
    <el-card class="monitor-card">
      <div slot="header">
        <span>实时监控面板</span>
        <el-button style="float: right" type="text" @click="refresh">刷新</el-button>
      </div>
      
      <!-- 性能指标 -->
      <el-row :gutter="20">
        <el-col :span="6" v-for="metric in metrics" :key="metric.name">
          <el-card shadow="hover" class="metric-card">
            <div class="metric-value">{{ metric.value }}</div>
            <div class="metric-name">{{ metric.name }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 实时图表 -->
      <div class="charts-container">
        <div class="chart">
          <v-chart :options="requestChart" />
        </div>
        <div class="chart">
          <v-chart :options="responseTimeChart" />
        </div>
      </div>

      <!-- 最新结果 -->
      <el-table :data="latestResults" style="width: 100%">
        <el-table-column prop="url" label="URL" />
        <el-table-column prop="status" label="状态码" width="100" />
        <el-table-column prop="type" label="类型" width="100" />
        <el-table-column prop="found_time" label="发现时间" width="180" />
      </el-table>
    </el-card>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { useWebSocket } from '@/composables/useWebSocket'
import { getTaskMetrics } from '@/api/dirScan'

export default {
  name: 'DirScanMonitor',
  props: {
    taskId: {
      type: Number,
      required: true
    }
  },
  setup(props) {
    const metrics = ref([
      { name: '请求速率', value: '0/s' },
      { name: '平均响应时间', value: '0ms' },
      { name: '错误率', value: '0%' },
      { name: '发现目录数', value: 0 }
    ])

    const latestResults = ref([])
    const requestChart = ref({})
    const responseTimeChart = ref({})

    // WebSocket连接
    const { data: wsData } = useWebSocket(`/ws/dirscan/${props.taskId}`)

    // 处理WebSocket消息
    const handleWSMessage = (message) => {
      switch (message.type) {
        case 'progress':
          updateMetrics(message.data)
          break
        case 'result':
          updateResults(message.data)
          break
        case 'alert':
          handleAlert(message.data)
          break
      }
    }

    // 更新指标
    const updateMetrics = (data) => {
      metrics.value[0].value = `${data.request_rate}/s`
      metrics.value[1].value = `${data.avg_response_time}ms`
      metrics.value[2].value = `${(data.error_rate * 100).toFixed(2)}%`
      metrics.value[3].value = data.directory_count

      updateCharts(data)
    }

    // 更新图表
    const updateCharts = (data) => {
      // 更新请求图表
      requestChart.value.series[0].data.push([
        new Date().getTime(),
        data.request_rate
      ])
      if (requestChart.value.series[0].data.length > 100) {
        requestChart.value.series[0].data.shift()
      }

      // 更新响应时间图表
      responseTimeChart.value.series[0].data.push([
        new Date().getTime(),
        data.avg_response_time
      ])
      if (responseTimeChart.value.series[0].data.length > 100) {
        responseTimeChart.value.series[0].data.shift()
      }
    }

    // 初始化图表配置
    const initCharts = () => {
      requestChart.value = {
        title: { text: '请求速率' },
        xAxis: { type: 'time' },
        yAxis: { type: 'value' },
        series: [{
          name: '请求/秒',
          type: 'line',
          data: []
        }]
      }

      responseTimeChart.value = {
        title: { text: '响应时间' },
        xAxis: { type: 'time' },
        yAxis: { type: 'value' },
        series: [{
          name: '毫秒',
          type: 'line',
          data: []
        }]
      }
    }

    // 刷新数据
    const refresh = async () => {
      try {
        const data = await getTaskMetrics(props.taskId)
        updateMetrics(data)
      } catch (error) {
        console.error('Failed to refresh metrics:', error)
      }
    }

    onMounted(() => {
      initCharts()
      refresh()
    })

    onUnmounted(() => {
      // 清理工作
    })

    return {
      metrics,
      latestResults,
      requestChart,
      responseTimeChart,
      refresh
    }
  }
}
</script>

<style scoped>
.dir-scan-monitor {
  padding: 20px;
}

.monitor-card {
  margin-bottom: 20px;
}

.metric-card {
  text-align: center;
  padding: 20px;
}

.metric-value {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
}

.metric-name {
  margin-top: 10px;
  color: #666;
}

.charts-container {
  display: flex;
  margin: 20px 0;
}

.chart {
  flex: 1;
  height: 300px;
  margin: 0 10px;
}
</style> 