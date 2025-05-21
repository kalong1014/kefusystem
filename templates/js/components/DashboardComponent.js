const DashboardComponent = {
    render() {
        return `
            <!-- 主应用 -->
            <div id="app" class="flex-1 flex flex-col md:flex-row">
                <!-- 侧边栏 -->
                ${SidebarComponent.render()}
                
                <!-- 主内容区 -->
                <main class="flex-1 flex flex-col bg-gray-50">
                    <!-- 顶部导航 -->
                    ${TopNavComponent.render()}
                    
                    <!-- 内容区域 -->
                    <div class="flex-1 overflow-y-auto p-4">
                        ${DashboardStatsComponent.render()}
                        ${ActiveSessionsComponent.render()}
                        ${FAQsComponent.render()}
                        ${OnlineStaffComponent.render()}
                    </div>
                </main>
            </div>
            
            <!-- 客服聊天组件 -->
            ${ChatWidgetComponent.render()}
        `;
    },
    
    afterRender() {
        // 初始化组件事件
        SidebarComponent.afterRender();
        ChatWidgetComponent.afterRender();
        TopNavComponent.afterRender();
    }
};    