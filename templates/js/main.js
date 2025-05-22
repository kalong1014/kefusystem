import { loadLanguage } from './i18n.js';
import LoginComponent from './components/LoginComponent.js';
import DashboardComponent from './components/DashboardComponent.js';

// 初始化应用
document.addEventListener('DOMContentLoaded', async () => {
    try {
        console.log('Starting application initialization...');
        await loadLanguage('zh-CN');
        await renderLoginPage();
        initEventListeners(); // 内部调用不需要导出
        console.log('Application initialized successfully');
    } catch (error) {
        console.error('应用初始化失败:', error);
        showToast('应用加载失败，请刷新页面', 'error');
    }
});

// 渲染登录页面
async function renderLoginPage() {
    try {
        console.log('Rendering login page...');
        
        const appContainer = document.getElementById('app-container');
        appContainer.innerHTML = await LoginComponent.render();
        
        // 执行渲染后操作
        if (typeof LoginComponent.afterRender === 'function') {
            await LoginComponent.afterRender();
        }
        
        console.log('Login page rendered successfully');
    } catch (error) {
        console.error('渲染登录页面失败:', error);
        showToast('登录页面加载失败', 'error');
    }
}
// 初始化事件监听
export function initEventListeners() {
    // 使用事件委托处理动态元素
    document.addEventListener('submit', async (e) => {
        if (e.target.id === 'login-form') {
            e.preventDefault();
            await authenticateUser();
        }
    });
    
    // 语言切换
    document.getElementById('lang-selector')?.addEventListener('change', async function() {
        try {
            await loadLanguage(this.value);
            if (document.getElementById('login-page')) {
                await renderLoginPage();
            } else {
                await renderDashboard();
            }
        } catch (error) {
            console.error('语言切换失败:', error);
            showToast('语言切换失败', 'error');
        }
    });
}


// 用户认证
async function authenticateUser() {
    try {
        const email = document.getElementById('email')?.value;
        const password = document.getElementById('password')?.value;
        
        // 简单验证
        if (!email || !password) {
            showToast('请输入邮箱和密码', 'error');
            return;
        }
        
        // 模拟登录验证
        async function simulateLogin(email, password) {
            // ... 保持不变 ...
        }
        
        const isAuthenticated = await simulateLogin(email, password);
        
        if (isAuthenticated) {
            // 隐藏登录页面，显示主应用
            document.getElementById('login-page')?.classList.add('hidden');
            document.getElementById('app')?.classList.remove('hidden');
            document.getElementById('chat-widget')?.classList.remove('hidden');
            
            // 渲染仪表盘
            await renderDashboard();
        } else {
            showToast('认证失败，请检查凭证', 'error');
        }
    } catch (error) {
        console.error('登录过程出错:', error);
        showToast('登录过程发生错误', 'error');
    }
}

// 模拟登录验证
async function simulateLogin(email, password) {
    // 实际项目中应替换为真实的API调用
    return new Promise(resolve => {
        setTimeout(() => resolve(true), 500); // 模拟API延迟
    });
}

// 渲染仪表盘
async function renderDashboard() {
    try {
        // 不需要再次动态导入，已经在文件顶部导入
        const appContainer = document.getElementById('app-container');
        appContainer.innerHTML = await DashboardComponent.render();
        
        // 执行渲染后操作
        if (typeof DashboardComponent.afterRender === 'function') {
            await DashboardComponent.afterRender();
        }
    } catch (error) {
        console.error('渲染仪表盘失败:', error);
        showToast('仪表盘加载失败', 'error');
    }
}  
// 显示通知提示
function showToast(message, type = 'info') {
    const toast = document.createElement('div');
    toast.className = `fixed top-4 right-4 px-6 py-3 rounded-lg shadow-lg transform transition-all duration-300 ${
        type === 'error' ? 'bg-red-500 text-white' : 'bg-blue-500 text-white'
    }`;
    toast.textContent = message;
    
    document.body.appendChild(toast);
    
    // 自动消失
    setTimeout(() => {
        toast.classList.add('opacity-0', 'translate-y-[-20px]');
        setTimeout(() => document.body.removeChild(toast), 300);
    }, 3000);
}    

// 导出 renderLoginPage 函数
export { renderLoginPage };