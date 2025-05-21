// 初始化应用
document.addEventListener('DOMContentLoaded', () => {
    // 加载语言文件
    loadLanguage('zh-CN').then(() => {
        // 渲染登录页面
        renderLoginPage();
        
        // 初始化事件监听
        initEventListeners();
    });
});

// 渲染登录页面
function renderLoginPage() {
    const appContainer = document.getElementById('app-container');
    appContainer.innerHTML = LoginComponent.render();
    LoginComponent.afterRender();
}

// 初始化事件监听
function initEventListeners() {
    // 登录表单提交
    document.getElementById('login-form').addEventListener('submit', function(e) {
        e.preventDefault();
        // 模拟登录
        authenticateUser();
    });
    
    // 语言切换
    document.getElementById('lang-selector').addEventListener('change', function() {
        loadLanguage(this.value).then(() => {
            // 重新渲染当前页面
            if (document.getElementById('login-page')) {
                renderLoginPage();
            } else {
                renderDashboard();
            }
        });
    });
}

// 模拟用户认证
function authenticateUser() {
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    
    // 简单验证
    if (email && password) {
        // 隐藏登录页面，显示主应用
        document.getElementById('login-page').classList.add('hidden');
        document.getElementById('app').classList.remove('hidden');
        document.getElementById('chat-widget').classList.remove('hidden');
        
        // 渲染仪表盘
        renderDashboard();
    } else {
        showToast('请输入邮箱和密码', 'error');
    }
}

// 渲染仪表盘
function renderDashboard() {
    const appContainer = document.getElementById('app-container');
    appContainer.innerHTML = DashboardComponent.render();
    DashboardComponent.afterRender();
}    