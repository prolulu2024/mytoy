package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

// 启动子进程
func run(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()
}

// 生成 sing-box 配置
func generateConfig(templatePath, outPath string) {
	data, err := ioutil.ReadFile(templatePath)
	if err != nil {
		fmt.Println("读取模板失败:", err)
		return
	}

	var cfg map[string]interface{}
	json.Unmarshal(data, &cfg)

	port := os.Getenv("SERVER_PORT")
	if port != "" {
		inbounds := cfg["inbounds"].([]interface{})
		inbounds[0].(map[string]interface{})["listen"] = "0.0.0.0"
		inbounds[0].(map[string]interface{})["port"] = port
	}

	newData, _ := json.MarshalIndent(cfg, "", "  ")
	ioutil.WriteFile(outPath, newData, 0644)
}

// 发送 Telegram 消息
func sendTG(msg string) {
	token := os.Getenv("TG_TOKEN")
	chat := os.Getenv("TG_CHAT")
	if token == "" || chat == "" {
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s",
		token, chat, msg)
	http.Get(url)
}

func main() {
	fmt.Println("🎉 启动玩具级节点系统...")

	// 1️⃣ 生成 sing-box 配置
	generateConfig("config-template.json", "config.json")

	// 2️⃣ 启动 sing-box 节点核心
	run("./sing-box", "run", "-c", "config.json")

	// 3️⃣ 启动哪吒 Agent
	nezhaServer := os.Getenv("NEZHA_SERVER")
	nezhaKey := os.Getenv("NEZHA_KEY")
	if nezhaServer != "" && nezhaKey != "" {
		run("./nezha-agent", "-s", nezhaServer, "-p", nezhaKey)
	}

	// 4️⃣ 可选启动 CF Tunnel
	cfToken := os.Getenv("CF_TOKEN")
	if cfToken != "" {
		run("./cloudflared", "tunnel", "--token", cfToken)
	}

	// 5️⃣ 发送 Telegram 通知
	sendTG("✅ 节点系统已启动成功！")

	select {} // 阻塞
}
