@echo off
echo 正在移除Windows网络转发规则...

REM 检查管理员权限
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo 错误：需要管理员权限运行此脚本
    echo 请右键点击脚本并选择"以管理员身份运行"
    pause
    exit /b 1
)

echo 移除tun0接口转发规则...

REM 删除Windows防火墙规则
echo 删除防火墙规则...
netsh advfirewall firewall delete rule name="TUN0 Accept All In" 2>nul
netsh advfirewall firewall delete rule name="TUN0 Accept All Out" 2>nul
netsh advfirewall firewall delete rule name="TUN0 Forward In" 2>nul
netsh advfirewall firewall delete rule name="TUN0 Forward Out" 2>nul

REM 禁用tun0接口转发
netsh interface ipv4 set interface "tun0" forwarding=disabled 2>nul

REM 禁用tun0接口
netsh interface set interface "tun0" disable 2>nul

REM 查找主网络适配器并禁用其转发
for /f "tokens=2 delims=:" %%i in ('netsh interface show interface ^| findstr "已连接"') do (
    set "MAIN_INTERFACE=%%i"
    goto :found_main_down
)

:found_main_down
if defined MAIN_INTERFACE (
    echo 禁用主网卡转发: %MAIN_INTERFACE%
    netsh interface ipv4 set interface "%MAIN_INTERFACE%" forwarding=disabled 2>nul
)

REM 禁用IP转发
netsh interface ipv4 set global forwarding=disabled

echo.
echo 网络转发规则已移除！
echo - 已禁用IP转发
echo - 已移除tun0接口配置
echo - 已删除路由规则
echo - 已删除防火墙规则

REM 显示当前状态
echo.
echo 当前tun0接口状态：
netsh interface show interface | findstr "tun0"

echo.
echo 当前防火墙规则（tun0相关）：
netsh advfirewall firewall show rule name="TUN0 Accept All In" 2>nul
netsh advfirewall firewall show rule name="TUN0 Accept All Out" 2>nul

pause
