const LoginComponent = {
    render() {
        return `
            <!-- 登录页面 -->
            <div id="login-page" class="flex items-center justify-center min-h-screen bg-gradient-to-br from-primary to-secondary">
                <div class="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md mx-4 transform transition-all duration-500 hover:scale-[1.02]">
                    <div class="text-center mb-8">
                        <h1 class="text-[clamp(1.8rem,5vw,2.5rem)] font-bold text-primary" data-i18n="systemName">智能客服系统</h1>
                        <p class="text-gray-500 mt-2" data-i18n="loginPrompt">请登录您的账户</p>
                    </div>
                    
                    <form id="login-form" class="space-y-5">
                        <div class="space-y-2">
                            <label for="email" class="block text-sm font-medium text-gray-700" data-i18n="email">邮箱</label>
                            <div class="relative">
                                <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                                    <i class="fa-solid fa-envelope"></i>
                                </span>
                                <input type="email" id="email" name="email" required
                                    class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all duration-300"
                                    placeholder="请输入您的邮箱">
                            </div>
                        </div>
                        
                        <div class="space-y-2">
                            <label for="password" class="block text-sm font-medium text-gray-700" data-i18n="password">密码</label>
                            <div class="relative">
                                <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-400">
                                    <i class="fa-solid fa-lock"></i>
                                </span>
                                <input type="password" id="password" name="password" required
                                    class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all duration-300"
                                    placeholder="请输入您的密码">
                            </div>
                        </div>
                        
                        <div class="flex items-center justify-between">
                            <div class="flex items-center">
                                <input id="remember-me" name="remember-me" type="checkbox"
                                    class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded">
                                <label for="remember-me" class="ml-2 block text-sm text-gray-700" data-i18n="rememberMe">记住我</label>
                            </div>
                            <a href="#" class="text-sm font-medium text-primary hover:text-primary/80 transition-colors duration-200" data-i18n="forgotPassword">忘记密码?</a>
                        </div>
                        
                        <button type="submit"
                            class="w-full bg-primary hover:bg-primary/90 text-white font-medium py-3 px-4 rounded-lg transition-all duration-300 transform hover:scale-[1.02] hover:shadow-lg flex items-center justify-center">
                            <span data-i18n="login">登录</span>
                            <i class="fa-solid fa-arrow-right ml-2"></i>
                        </button>
                    </form>
                    
                    <div class="mt-8 text-center">
                        <p class="text-gray-500" data-i18n="noAccount">还没有账户? <a href="#" class="font-medium text-primary hover:text-primary/80 transition-colors duration-200" data-i18n="register">注册</a></p>
                    </div>
                </div>
                
                <!-- 语言选择器 -->
                <div class="absolute top-4 right-4">
                    <select id="lang-selector" class="bg-white/20 backdrop-blur-sm text-white border border-white/30 rounded-lg px-3 py-1 focus:outline-none focus:ring-2 focus:ring-white/50">
                        <option value="zh-CN" selected>中文</option>
                        <option value="en-US">English</option>
                    </select>
                </div>
            </div>
        `;
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
    }
};    