// 语言资源
let translations = {};
let currentLanguage = 'zh-CN';

// 加载语言文件
async function loadLanguage(lang) {
    try {
        const response = await fetch(`locales/${lang}.json`);
        translations = await response.json();
        currentLanguage = lang;
        document.documentElement.lang = lang;
        applyTranslations();
        return true;
    } catch (error) {
        console.error('加载语言文件失败:', error);
        return false;
    }
}

// 应用翻译到DOM
function applyTranslations() {
    document.querySelectorAll('[data-i18n]').forEach(element => {
        const key = element.getAttribute('data-i18n');
        if (translations[key]) {
            element.textContent = translations[key];
        }
    });
    
    document.querySelectorAll('[data-i18n-placeholder]').forEach(element => {
        const key = element.getAttribute('data-i18n-placeholder');
        if (translations[key]) {
            element.placeholder = translations[key];
        }
    });
    
    // 更新页面标题
    const pageTitleKey = document.getElementById('page-title').getAttribute('data-i18n');
    if (pageTitleKey && translations[pageTitleKey]) {
        document.title = translations[pageTitleKey];
    }
}

// 获取翻译文本
function __(key) {
    return translations[key] || key;
}

// 注册翻译函数到全局
window.__ = __;    