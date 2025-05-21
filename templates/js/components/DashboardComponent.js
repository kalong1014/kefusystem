// 方案1：使用默认导出（推荐）
export default {
    render() {
        return `
            <div id="dashboard" class="min-h-screen bg-gray-50">
                <!-- 仪表盘内容 -->
                <header class="bg-white shadow-sm">
                    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                        <div class="flex justify-between h-16">
                            <div class="flex">
                                <div class="flex-shrink-0 flex items-center">
                                    <h1 class="text-xl font-bold text-gray-900" data-i18n="dashboard">仪表盘</h1>
                                </div>
                            </div>
                            <div class="flex items-center">
                                <button id="logout-btn" class="ml-3 inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-md text-white bg-red-600 hover:bg-red-500 focus:outline-none focus:border-red-700 focus:shadow-outline-red transition ease-in-out duration-150">
                                    <i class="fa-solid fa-sign-out-alt mr-2"></i>
                                    <span data-i18n="logout">退出</span>
                                </button>
                            </div>
                        </div>
                    </div>
                </header>
                
                <main>
                    <div class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                        <!-- 主内容区 -->
                        <div class="bg-white overflow-hidden shadow-sm sm:rounded-lg">
                            <div class="p-6 bg-white border-b border-gray-200">
                                <h2 class="text-lg font-medium text-gray-900 mb-4" data-i18n="welcome-user">欢迎回来，用户</h2>
                                
                                <!-- 聊天区域 -->
                                <div id="chat-widget" class="hidden">
                                    <div class="border rounded-lg overflow-hidden shadow-md max-w-2xl mx-auto">
                                        <div class="bg-blue-600 text-white p-3 flex justify-between items-center">
                                            <h3 class="font-semibold" data-i18n="chat-with-us">与客服聊天</h3>
                                            <button id="toggle-chat" class="text-white hover:text-gray-200">
                                                <i class="fa fa-chevron-down"></i>
                                            </button>
                                        </div>
                                        <div id="chat-messages" class="p-4 h-64 overflow-y-auto bg-white">
                                            <!-- 聊天消息将在这里动态添加 -->
                                        </div>
                                        <div class="p-3 bg-gray-100">
                                            <form id="chat-form" class="flex">
                                                <input type="text" id="chat-input" class="flex-1 border rounded-l-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500" placeholder="输入消息..." required>
                                                <button type="submit" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-r-md">
                                                    <i class="fa fa-paper-plane"></i>
                                                </button>
                                            </form>
                                        </div>
                                    </div>
                                </div>
                                
                                <!-- 其他仪表盘内容 -->
                                <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
                                    <div class="bg-white p-6 rounded-lg shadow-md">
                                        <div class="flex items-center">
                                            <div class="flex-shrink-0 bg-blue-100 rounded-full p-3">
                                                <i class="fa fa-comments text-blue-600 text-xl"></i>
                                            </div>
                                            <div class="ml-4">
                                                <h3 class="text-sm font-medium text-gray-500" data-i18n="open-tickets">开放工单</h3>
                                                <div class="text-2xl font-bold text-gray-900">5</div>
                                            </div>
                                        </div>
                                    </div>
                                    
                                    <div class="bg-white p-6 rounded-lg shadow-md">
                                        <div class="flex items-center">
                                            <div class="flex-shrink-0 bg-green-100 rounded-full p-3">
                                                <i class="fa fa-check-circle text-green-600 text-xl"></i>
                                            </div>
                                            <div class="ml-4">
                                                <h3 class="text-sm font-medium text-gray-500" data-i18n="resolved-tickets">已解决工单</h3>
                                                <div class="text-2xl font-bold text-gray-900">24</div>
                                            </div>
                                        </div>
                                    </div>
                                    
                                    <div class="bg-white p-6 rounded-lg shadow-md">
                                        <div class="flex items-center">
                                            <div class="flex-shrink-0 bg-yellow-100 rounded-full p-3">
                                                <i class="fa fa-clock-o text-yellow-600 text-xl"></i>
                                            </div>
                                            <div class="ml-4">
                                                <h3 class="text-sm font-medium text-gray-500" data-i18n="avg-response">平均响应时间</h3>
                                                <div class="text-2xl font-bold text-gray-900">12分钟</div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </main>
            </div>
        `;
    },
    
    afterRender() {
        // 初始化聊天功能
        initChatWidget();
        
        // 添加登出事件监听
        document.getElementById('logout-btn')?.addEventListener('click', function() {
            // 清除用户会话
            localStorage.removeItem('auth_token');
            
            // 重新渲染登录页面
            document.getElementById('app').classList.add('hidden');
            document.getElementById('chat-widget').classList.add('hidden');
            renderLoginPage();
        });
        
        console.log('Dashboard rendered successfully');
    }
};

// 聊天功能初始化
function initChatWidget() {
    // 聊天表单提交
    document.getElementById('chat-form')?.addEventListener('submit', function(e) {
        e.preventDefault();
        const input = document.getElementById('chat-input');
        const message = input.value.trim();
        
        if (message) {
            addMessage('user', message);
            input.value = '';
            
            // 模拟客服回复
            setTimeout(() => {
                addMessage('agent', '感谢您的咨询，我们将尽快回复您的问题。');
            }, 1000);
        }
    });
    
    // 切换聊天窗口
    document.getElementById('toggle-chat')?.addEventListener('click', function() {
        const chatMessages = document.getElementById('chat-messages');
        chatMessages.classList.toggle('hidden');
        this.innerHTML = chatMessages.classList.contains('hidden') 
            ? '<i class="fa fa-chevron-up"></i>' 
            : '<i class="fa fa-chevron-down"></i>';
    });
}

// 添加聊天消息
function addMessage(sender, message) {
    const chatMessages = document.getElementById('chat-messages');
    const msgDiv = document.createElement('div');
    
    if (sender === 'user') {
        msgDiv.className = 'flex justify-end mb-2';
        msgDiv.innerHTML = `
            <div class="bg-blue-600 text-white rounded-lg p-3 max-w-xs break-words">
                ${message}
            </div>
        `;
    } else {
        msgDiv.className = 'flex justify-start mb-2';
        msgDiv.innerHTML = `
            <div class="bg-gray-200 text-gray-800 rounded-lg p-3 max-w-xs break-words">
                ${message}
            </div>
        `;
    }
    
    chatMessages.appendChild(msgDiv);
    chatMessages.scrollTop = chatMessages.scrollHeight;
}    