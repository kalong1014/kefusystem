// 确保LoginComponent被正确定义
export default {
    render() {
        return `
            <div id="login-page" class="min-h-screen flex items-center justify-center bg-gray-50">
                <div class="max-w-md w-full p-6 bg-white rounded-lg shadow-xl">
                    <div class="text-center mb-8">
                        <h1 class="text-3xl font-bold text-gray-800" data-i18n="welcome">欢迎使用智能客服</h1>
                    </div>
                    <form id="login-form" class="space-y-4">
                        <div>
                            <label for="email" class="block text-sm font-medium text-gray-700" data-i18n="email">邮箱</label>
                            <input type="email" id="email" class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" required>
                        </div>
                        <div>
                            <label for="password" class="block text-sm font-medium text-gray-700" data-i18n="password">密码</label>
                            <input type="password" id="password" class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent" required>
                        </div>
                        <div class="flex items-center justify-between">
                            <div class="flex items-center">
                                <input id="remember-me" name="remember-me" type="checkbox" class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded">
                                <label for="remember-me" class="ml-2 block text-sm text-gray-900" data-i18n="remember">记住我</label>
                            </div>
                            <div class="text-sm">
                                <a href="#" class="font-medium text-blue-600 hover:text-blue-500" data-i18n="forgot">忘记密码?</a>
                            </div>
                        </div>
                        <div>
                            <button type="submit" class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500" data-i18n="login">登录</button>
                        </div>
                    </form>
                    <div class="mt-8">
                        <div class="relative">
                            <div class="absolute inset-0 flex items-center">
                                <div class="w-full border-t border-gray-300"></div>
                            </div>
                            <div class="relative flex justify-center text-sm">
                                <span class="px-2 bg-white text-gray-500" data-i18n="or">或者</span>
                            </div>
                        </div>
                        <div class="mt-6 grid grid-cols-3 gap-3">
                            <div>
                                <a href="#" class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                                    <i class="fa-brands fa-weixin text-green-600 mr-2"></i>
                                    <span data-i18n="wechat">微信</span>
                                </a>
                            </div>
                            <div>
                                <a href="#" class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                                    <i class="fa-brands fa-qq text-blue-500 mr-2"></i>
                                    <span data-i18n="qq">QQ</span>
                                </a>
                            </div>
                            <div>
                                <a href="#" class="w-full inline-flex justify-center py-2 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50">
                                    <i class="fa-brands fa-google text-red-500 mr-2"></i>
                                    <span data-i18n="google">Google</span>
                                </a>
                            </div>
                        </div>
                    </div>
                    <div class="mt-6 text-center text-sm text-gray-500">
                        <span data-i18n="need-help">需要帮助?</span>
                        <a href="#" class="font-medium text-blue-600 hover:text-blue-500" data-i18n="contact">联系我们</a>
                    </div>
                </div>
            </div>
        `;
    },
    async afterRender() {
        // 绑定登录表单提交事件
        const loginForm = document.getElementById('login-form');
        loginForm.addEventListener('submit', async (event) => {
            event.preventDefault();
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;

            // 简单验证
            if (!email || !password) {
                showToast('请输入邮箱和密码', 'error');
                return;
            }

            // 模拟登录验证
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
        });

        this.addLoginFormAnimation();
    },   
    afterRender() {
        // 添加登录表单动画
        const formInputs = document.querySelectorAll('#login-form input');
        formInputs.forEach(input => {
            input.addEventListener('focus', () => {
                input.parentElement.parentElement.classList.add('scale-[1.02]');
                input.parentElement.parentElement.style.transition = 'transform 0.3s ease';
            });
            
            input.addEventListener('blur', () => {
                input.parentElement.parentElement.classList.remove('scale-[1.02]');
            });
        });
            console.log('Login form rendered');
    } 
}   