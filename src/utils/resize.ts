// 调整 iframe 高度
function adjustIframeHeight() {
  try {
    const iframe = window.frameElement as HTMLIFrameElement | null
    if (!iframe) return

    // 统一获取元素真实高度
    const { scrollHeight, offsetHeight } = document.documentElement
    const height = Math.max(scrollHeight, offsetHeight)

    // 避免重复赋值
    if (iframe.style.height !== `${height}px`) {
      iframe.style.height = `${height}px`
    }
  } catch (e) {
    console.error('调整 iframe 高度失败:', e)
  }
}

// 初始化调整函数
const resizeIframeHeight = function () {
  // 初始化
  window.addEventListener('load', adjustIframeHeight)

  // 监听 DOM 变化
  const observer = new MutationObserver(adjustIframeHeight)
  observer.observe(document.body, {
    attributes: true,
    childList: true,
    subtree: true,
    characterData: true,
  })

  // 监听窗口尺寸变化
  window.addEventListener('resize', adjustIframeHeight)
}

export default resizeIframeHeight
