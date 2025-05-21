/**
 * 格式化时间戳为可读的日期时间格式
 * @param {number} timestamp - 时间戳（毫秒）
 * @param {boolean} [showSeconds=false] - 是否显示秒
 * @returns {string} 格式化的日期时间字符串
 */
export function formatDateTime(timestamp, showSeconds = false) {
  if (typeof timestamp !== 'number') {
    console.warn('Invalid timestamp:', timestamp);
    return '';
  }
  
  const date = new Date(timestamp);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  
  return showSeconds 
    ? `${year}-${month}-${day} ${hours}:${minutes}:${seconds}` 
    : `${year}-${month}-${day} ${hours}:${minutes}`;
}

/**
 * 生成唯一ID
 * @returns {string} 唯一ID
 */
export function generateUniqueId() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c === 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
}

/**
 * 检查字符串是否为空
 * @param {string} str - 待检查的字符串
 * @returns {boolean} 如果字符串为空或仅包含空格，返回true，否则返回false
 */
export function isEmpty(str) {
  return str === undefined || str === null || (typeof str === 'string' && str.trim() === '');
}

/**
 * 对象深拷贝
 * @param {Object} obj - 待拷贝的对象
 * @returns {Object} 拷贝后的新对象
 */
export function deepClone(obj) {
  if (obj === null || typeof obj !== 'object') {
    return obj;
  }
  
  let clone;
  
  if (Array.isArray(obj)) {
    clone = [];
    for (let i = 0; i < obj.length; i++) {
      clone[i] = deepClone(obj[i]);
    }
  } else {
    clone = {};
    for (const key in obj) {
      if (obj.hasOwnProperty(key)) {
        clone[key] = deepClone(obj[key]);
      }
    }
  }
  
  return clone;
}

/**
 * 检查是否为移动设备
 * @returns {boolean} 如果是移动设备返回true，否则返回false
 */
export function isMobileDevice() {
  return /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent);
}

/**
 * 获取URL参数
 * @param {string} paramName - 参数名
 * @returns {string|null} 参数值，如果不存在则返回null
 */
export function getUrlParam(paramName) {
  try {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get(paramName);
  } catch (error) {
    console.error('Failed to parse URL parameters:', error);
    return null;
  }
}

/**
 * 防抖函数
 * @param {Function} func - 要执行的函数
 * @param {number} delay - 延迟时间（毫秒）
 * @returns {Function} 防抖处理后的函数
 */
export function debounce(func, delay) {
  if (typeof func !== 'function') {
    throw new Error('debounce: first argument must be a function');
  }
  
  let timer;
  return function() {
    const context = this;
    const args = arguments;
    clearTimeout(timer);
    timer = setTimeout(() => func.apply(context, args), delay);
  };
}

/**
 * 节流函数
 * @param {Function} func - 要执行的函数
 * @param {number} limit - 限制时间（毫秒）
 * @returns {Function} 节流处理后的函数
 */
export function throttle(func, limit) {
  if (typeof func !== 'function') {
    throw new Error('throttle: first argument must be a function');
  }
  
  let inThrottle;
  return function() {
    const context = this;
    const args = arguments;
    if (!inThrottle) {
      func.apply(context, args);
      inThrottle = true;
      setTimeout(() => inThrottle = false, limit);
    }
  };
}

/**
 * 发送HTTP请求
 * @param {string} url - 请求URL
 * @param {Object} [options={}] - 请求选项
 * @param {string} [options.method='GET'] - 请求方法
 * @param {Object} [options.headers={}] - 请求头
 * @param {Object} [options.body=null] - 请求体
 * @returns {Promise<Object>} 响应数据
 */
export async function fetchData(url, options = {}) {
  try {
    const response = await fetch(url, {
      method: options.method || 'GET',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers
      },
      body: options.body ? JSON.stringify(options.body) : null
    });
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.includes('application/json')) {
      return await response.json();
    } else {
      return await response.text();
    }
  } catch (error) {
    console.error('Fetch error:', error);
    throw error;
  }
}    