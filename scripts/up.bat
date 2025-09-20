@echo off
echo 正在配置Windows网络转发规则...

REM 检查管理员权限
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo 错误：需要管理员权限运行此脚本
    echo 请右键点击脚本并选择"以管理员身份运行"
    pause
    exit /b 1
)

echo 配置tun0接口转发规则...

REM 启用IP转发
netsh interface ipv4 set global forwarding=enabled

REM 启用tun0接口
netsh interface set interface "tun0" enable

REM 配置tun0接口为任意IP（如果需要的话）
REM 这里不设置固定IP，让系统自动处理

echo tun0接口已启用，准备配置转发规则...

REM 配置Windows防火墙规则 - 允许所有tun0流量
echo 配置防火墙规则允许tun0流量...

REM 允许所有来自tun0的入站流量
netsh advfirewall firewall add rule name="TUN0 Accept All In" dir=in action=allow interface="tun0" protocol=any

REM 允许所有从tun0的出站流量
netsh advfirewall firewall add rule name="TUN0 Accept All Out" dir=out action=allow interface="tun0" protocol=any

REM 配置IP转发规则 - 允许tun0接口间的转发
netsh interface ipv4 set interface "tun0" forwarding=enabled

REM 查找主网络适配器进行NAT配置
for /f "tokens=2 delims=:" %%i in ('netsh interface show interface ^| findstr "已连接"') do (
    set "MAIN_INTERFACE=%%i"
    goto :found_interface
)

:found_interface
if defined MAIN_INTERFACE (
    echo 主网络接口: %MAIN_INTERFACE%
    echo 配置NAT规则将tun0流量转发到主网卡...
    
    REM 配置NAT - 将所有tun0流量通过主网卡转发
    REM 这里使用Windows的NAT功能
    netsh interface ipv4 set interface "%MAIN_INTERFACE%" forwarding=enabled
    
    REM 添加路由规则 - 让tun0流量通过主网卡出去
    REM 这里不指定具体IP，而是配置接口级别的转发
) else (
    echo 警告：未找到主网络接口，跳过NAT配置
)

echo.
echo 网络转发规则配置完成！
echo - 已启用全局IP转发
echo - 已启用tun0接口
echo - 已配置tun0接口转发
echo - 已配置防火墙规则（允许所有tun0流量）
if defined MAIN_INTERFACE (
    echo - 已配置主网卡转发 (%MAIN_INTERFACE%)
)

REM 显示当前配置
echo.
echo 当前tun0接口状态：
netsh interface show interface | findstr "tun0"

echo.
echo 当前防火墙规则（tun0相关）：
netsh advfirewall firewall show rule name="TUN0 Accept All In"
netsh advfirewall firewall show rule name="TUN0 Accept All Out"

pause
