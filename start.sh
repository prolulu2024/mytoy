#!/bin/bash
set -e

echo "⬇️ 下载核心程序..."
curl -L https://github.com/prolulu2024/mytoy/releases/download/v1.0/main-amd -o main-amd
chmod +x main-amd

echo "⬇️ 准备哪吒 Agent..."
chmod +x nezha-agent

if [ ! -z "$CF_TOKEN" ]; then
    echo "⬇️ 准备 Cloudflare 隧道..."
    chmod +x cloudflared
fi

echo "🚀 启动主程序..."
./main-amd
