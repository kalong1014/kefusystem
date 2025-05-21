/**
 * 国际化支持模块
 */
let currentLocale = 'zh-CN';
let translations = {};

/**
 * 加载语言文件
 * @param {string} locale - 语言代码，如 'zh-CN'
 * @returns {Promise<void>}
 */
export async function loadLanguage(locale) {
  try {
    const response = await fetch(`/locales/${locale}.json`);
    if (!response.ok) {
      throw new Error(`Failed to load locale: ${locale}`);
    }
    
    translations = await response.json();
    currentLocale = locale;
    console.log(`Loaded locale: ${locale}`);
    
    // 更新页面上的所有翻译
    updateTranslations();
  } catch (error) {
    console.error('Error loading language:', error);
    throw error;
  }
}

/**
 * 获取翻译文本
 * @param {string} key - 翻译键
 * @param {Object} [placeholders={}] - 占位符对象
 * @returns {string} 翻译后的文本
 */
export function translate(key, placeholders = {}) {
  let text = translations[key] || key;
  
  // 替换占位符
  Object.keys(placeholders).forEach(placeholder => {
    text = text.replace(`{${placeholder}}`, placeholders[placeholder]);
  });
  
  return text;
}

/**
 * 更新页面上的所有翻译元素
 */
function updateTranslations() {
  document.querySelectorAll('[data-i18n]').forEach(element => {
    const key = element.getAttribute('data-i18n');
    element.textContent = translate(key);
  });
  
  document.querySelectorAll('[data-i18n-placeholder]').forEach(element => {
    const key = element.getAttribute('data-i18n-placeholder');
    element.setAttribute('placeholder', translate(key));
  });
}

/**
 * 获取当前语言
 * @returns {string} 当前语言代码
 */
export function getCurrentLocale() {
  return currentLocale;
}    