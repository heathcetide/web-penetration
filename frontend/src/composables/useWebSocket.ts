import { ref, onMounted, onUnmounted } from 'vue'

export function useWebSocket(url: string) {
  const data = ref(null)
  const error = ref(null)
  let ws: WebSocket | null = null

  const connect = () => {
    ws = new WebSocket(`ws://${window.location.host}${url}`)

    ws.onmessage = (event) => {
      try {
        data.value = JSON.parse(event.data)
      } catch (e) {
        error.value = e
      }
    }

    ws.onerror = (e) => {
      error.value = e
    }

    ws.onclose = () => {
      // 重连逻辑
      setTimeout(connect, 1000)
    }
  }

  onMounted(() => {
    connect()
  })

  onUnmounted(() => {
    if (ws) {
      ws.close()
    }
  })

  return {
    data,
    error
  }
} 