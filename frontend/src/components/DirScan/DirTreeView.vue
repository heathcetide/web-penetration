<template>
  <div class="dir-tree-view">
    <el-card>
      <div slot="header">
        <span>目录结构</span>
        <el-button-group style="float: right">
          <el-button size="small" @click="expandAll">展开全部</el-button>
          <el-button size="small" @click="collapseAll">收起全部</el-button>
          <el-button size="small" @click="refresh">刷新</el-button>
        </el-button-group>
      </div>

      <!-- 搜索过滤 -->
      <el-input
        v-model="filterText"
        placeholder="输入关键字过滤"
        prefix-icon="el-icon-search"
        clearable
      />

      <!-- 目录树 -->
      <el-tree
        ref="tree"
        :data="treeData"
        :props="defaultProps"
        :filter-node-method="filterNode"
        node-key="path"
        :expand-on-click-node="false"
        :render-content="renderContent"
        @node-click="handleNodeClick"
      >
      </el-tree>
    </el-card>

    <!-- 节点详情对话框 -->
    <el-dialog
      title="节点详情"
      :visible.sync="dialogVisible"
      width="60%"
    >
      <div v-if="selectedNode">
        <el-descriptions border>
          <el-descriptions-item label="路径">{{ selectedNode.path }}</el-descriptions-item>
          <el-descriptions-item label="类型">{{ selectedNode.type }}</el-descriptions-item>
          <el-descriptions-item label="大小">{{ formatSize(selectedNode.size) }}</el-descriptions-item>
          <el-descriptions-item label="状态码">{{ selectedNode.metadata?.status_code }}</el-descriptions-item>
          <el-descriptions-item label="内容类型">{{ selectedNode.metadata?.content_type }}</el-descriptions-item>
          <el-descriptions-item label="发现时间">{{ formatTime(selectedNode.metadata?.found_time) }}</el-descriptions-item>
        </el-descriptions>

        <!-- 漏洞信息 -->
        <div v-if="selectedNode.vulnerabilities?.length" class="vuln-info">
          <h3>发现的漏洞</h3>
          <el-table :data="selectedNode.vulnerabilities">
            <el-table-column prop="type" label="类型" width="120" />
            <el-table-column prop="severity" label="严重程度" width="100">
              <template #default="{ row }">
                <el-tag :type="getSeverityType(row.severity)">{{ row.severity }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="description" label="描述" />
          </el-table>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { ref, watch, onMounted } from 'vue'
import { getDirectoryTree } from '@/api/dirScan'
import { formatBytes, formatDateTime } from '@/utils/format'

export default {
  name: 'DirTreeView',
  props: {
    taskId: {
      type: Number,
      required: true
    }
  },
  setup(props) {
    const tree = ref(null)
    const treeData = ref([])
    const filterText = ref('')
    const dialogVisible = ref(false)
    const selectedNode = ref(null)

    const defaultProps = {
      children: 'children',
      label: 'name'
    }

    // 过滤节点
    const filterNode = (value: string, data: any) => {
      if (!value) return true
      return data.path.toLowerCase().includes(value.toLowerCase())
    }

    // 监听过滤文本变化
    watch(filterText, (val) => {
      tree.value?.filter(val)
    })

    // 渲染节点内容
    const renderContent = (h: any, { node, data }: any) => {
      return h('span', { class: 'custom-tree-node' }, [
        h('span', { class: `node-type-${data.type}` }, data.name),
        h('span', { class: 'node-meta' }, [
          data.type === 'file' && h('span', { class: 'node-size' }, formatBytes(data.size)),
          data.vulnerabilities?.length > 0 && h('el-badge', {
            props: {
              value: data.vulnerabilities.length,
              type: 'danger'
            }
          })
        ])
      ])
    }

    // 处理节点点击
    const handleNodeClick = (data: any) => {
      selectedNode.value = data
      dialogVisible.value = true
    }

    // 展开所有节点
    const expandAll = () => {
      const keys = getAllKeys(treeData.value)
      keys.forEach(key => tree.value?.store.nodesMap[key].expand())
    }

    // 收起所有节点
    const collapseAll = () => {
      const keys = getAllKeys(treeData.value)
      keys.forEach(key => tree.value?.store.nodesMap[key].collapse())
    }

    // 获取所有节点的key
    const getAllKeys = (nodes: any[]): string[] => {
      let keys: string[] = []
      nodes.forEach(node => {
        keys.push(node.path)
        if (node.children?.length) {
          keys = keys.concat(getAllKeys(node.children))
        }
      })
      return keys
    }

    // 获取严重程度类型
    const getSeverityType = (severity: string): string => {
      const types: { [key: string]: string } = {
        high: 'danger',
        medium: 'warning',
        low: 'info'
      }
      return types[severity] || 'info'
    }

    // 刷新数据
    const refresh = async () => {
      try {
        const data = await getDirectoryTree(props.taskId)
        treeData.value = [data] // 根节点
      } catch (error) {
        console.error('Failed to fetch directory tree:', error)
      }
    }

    onMounted(() => {
      refresh()
    })

    return {
      tree,
      treeData,
      filterText,
      defaultProps,
      dialogVisible,
      selectedNode,
      filterNode,
      renderContent,
      handleNodeClick,
      expandAll,
      collapseAll,
      refresh,
      getSeverityType,
      formatSize: formatBytes,
      formatTime: formatDateTime
    }
  }
}
</script>

<style scoped>
.dir-tree-view {
  padding: 20px;
}

.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-right: 8px;
}

.node-type-directory {
  color: #409EFF;
}

.node-type-file {
  color: #67C23A;
}

.node-meta {
  font-size: 12px;
  color: #909399;
}

.node-size {
  margin-right: 8px;
}

.vuln-info {
  margin-top: 20px;
}
</style> 