package test

import (
	// Re-added standard library zlib
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json" // 添加 errors 包导入
	"io"
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestDiscordClient(t *testing.T) {
	token := "XDINTg1NjU0MDMzNDE2MTky.GYxlGL"
	dg, err := discordgo.New(token)
	if err != nil {
		t.Fatalf("创建22223333322222 discord session 失败: %v", err)
	}
	dg.Debug = true
	dg.LogLevel = 10

	adapter := discordgo.NewMirrorZlibAdapter()
	dg.SetMirrorZlibAdapter(adapter)
	dg.SetBaseURL("https://discord-d-com-s-mjc.aiwentu.net/")
	dg.SetCookie("___SID=p9_xIHlgY4QHI6_rZ915d; ___AID_mjc=1; ___AID_mjg2=1; ___MJ4=%7B%22mj%22%3A%7B%22s%22%3A%22%22%2C%22style%22%3A%22%22%2C%22v%22%3A%226.1%22%2C%22w%22%3A%22%22%2C%22params%22%3A%22%22%2C%22desc%22%3A%22%22%7D%2C%22niji%22%3A%7B%22s%22%3A%22%22%2C%22style%22%3A%22%22%2C%22v%22%3A%226%22%2C%22w%22%3A%22%22%2C%22params%22%3A%22%22%2C%22desc%22%3A%22%22%7D%2C%22_v%22%3A%222%22%2C%22toEn%22%3Afalse%2C%22relax%22%3Afalse%7D; ___FP2=fbo8heyJsYW5ndWFnZXMiOiJ6aC1DTix6aCIsInVzZXJBZ2VudCI6Ik1vemlsbGEvNS4wIChNYWNpbnRvc2g7IEludGVsIE1hYyBPUyBYIDEwXzE1XzcpIEFwcGxlV2ViS2l0LzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZS8xMzUuMC4wLjAgU2FmYXJpLzUzNy4zNiIsInRpbWVab25lIjoiQXNpYS9TaGFuZ2hhaSIsInRpbWVab25lT2Zmc2V0Ijo4LCJkZXZpY2VNZW1vcnkiOjgsIm1heFRvdWNoUG9pbnRzIjowLCJoYXJkd2FyZUNvbmN1cnJlbmN5Ijo4LCJwbGF0Zm9ybSI6Im1hY09TIiwic2NyZWVuIjoiMTQ0MCw4MDAiLCJ3ZWJkcml2ZXIiOmZhbHNlLCJ2aWRlb0NhcmQiOiJBTkdMRSAoQXBwbGUsIEFOR0xFIE1ldGFsIFJlbmRlcmVyOiBBcHBsZSBNMSwgVW5zcGVjaWZpZWQgVmVyc2lvbikiLCJjYW52YXNJZCI6MjExNzM2MDQ5LCJkZXNrdG9wIjpmYWxzZSwiaXNQcml2YXRlIjpmYWxzZSwiZnAiOjQ1OTI3NDQwNCwiaGFzR2Z3Ijp0cnVlLCJpcENuIjoiMTA2LjExNy4yMjYuMzEiLCJpcEZyb20iOiJxcSIsImNkcCI6dHJ1ZX0.1y8h")

	//dg.SetGateway("wss://gateway-d-discord-d-gg-s-mjc.aiwentu.net")
	err = dg.Open()
	if err != nil {
		t.Fatalf("打开 discord session 失败: %v", err)
	}
	t.Log("打开 discord session 成功")
	dg.Close()
}

var messagerB64 = "eF6qVipRssorzcnRUSqGMfILlKwMDXSUUpSsqpUyUhOLSpJSE0viM/NKUovKEnOUrEwMjUwNdJTiS4oSk1OVrKJja2sBAAAA//8="
var messagerB642 = "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"

// 模式discodews 流 ws只有第一个有78头 必须使用一个全局解包器！
func TestDecodeSpecificMessage(t *testing.T) {
	// 解码两个 Base64 字符串
	msgBytes1, err1 := base64.StdEncoding.DecodeString(messagerB64)
	if err1 != nil {
		t.Fatalf("Failed to decode base64 message 1: %v", err1)
	}
	msgBytes2, err2 := base64.StdEncoding.DecodeString(messagerB642)
	if err2 != nil {
		t.Fatalf("Failed to decode base64 message 2: %v", err2)
	}
	t.Logf("Decoded message 1 (%d bytes) and message 2 (%d bytes).", len(msgBytes1), len(msgBytes2))

	// --- 测试: 镜像站适配器作为 io.Reader ---
	t.Run("Test Mirror Adapter as io.Reader", func(t *testing.T) {
		t.Logf("--- 测试: 镜像站适配器作为 io.Reader ---")

		adapter := discordgo.NewMirrorZlibAdapter()
		defer func() {
			if err := adapter.Close(); err != nil {
				// Close 现在内部处理预期的 EOF
				t.Errorf("adapter.Close() failed unexpectedly: %v", err)
			} else {
				t.Logf("Mirror adapter closed successfully.")
			}
		}()

		// 使用适配器作为 io.Reader 创建 json.Decoder
		decoder := json.NewDecoder(adapter)

		// 1. 追加第一个消息
		t.Logf("Appending message 1...")
		if err := adapter.AppendMessage(msgBytes1); err != nil {
			t.Fatalf("adapter.AppendMessage(msgBytes1) failed: %v", err)
		}
		t.Logf("Appended message 1 (%d bytes).", len(msgBytes1))

		// 2. 尝试解码第一个 JSON 对象
		t.Logf("Attempting to decode JSON object 1...")
		var event1 discordgo.Event // 假设解码目标是 Event 类型
		if err := decoder.Decode(&event1); err != nil {
			t.Fatalf("decoder.Decode() for event 1 failed: %v", err)
		}
		t.Logf("Successfully decoded JSON object 1. Op: %d, Type: %s", event1.Operation, event1.Type)
		// (可选) 打印解码出的 JSON
		prettyJSON1, _ := json.MarshalIndent(event1, "", "  ")
		t.Logf("Decoded JSON 1:\n%s", string(prettyJSON1))

		// 3. 追加第二个消息
		t.Logf("Appending message 2...")
		if err := adapter.AppendMessage(msgBytes2); err != nil {
			t.Fatalf("adapter.AppendMessage(msgBytes2) failed: %v", err)
		}
		t.Logf("Appended message 2 (%d bytes).", len(msgBytes2))

		// 4. 尝试解码第二个 JSON 对象 (使用同一个 decoder)
		t.Logf("Attempting to decode JSON object 2...")
		var event2 discordgo.Event // 假设解码目标是 Event 类型
		if err := decoder.Decode(&event2); err != nil {
			// 如果这里出现 EOF 错误，可能表示流意外结束，或者解码器状态问题
			t.Fatalf("decoder.Decode() for event 2 failed: %v", err)
		}
		t.Logf("Successfully decoded JSON object 2. Op: %d, Type: %s", event2.Operation, event2.Type)
		// (可选) 打印解码出的 JSON
		prettyJSON2, _ := json.MarshalIndent(event2, "", "  ")
		t.Logf("Decoded JSON 2:\n%s", string(prettyJSON2))

		// 5. (可选) 检查流是否还有更多数据
		t.Logf("Checking if decoder has more data...")
		if decoder.More() {
			t.Errorf("decoder.More() returned true unexpectedly after decoding both objects.")
			// 尝试读取并打印剩余内容
			var remaining interface{}
			if err := decoder.Decode(&remaining); err != nil {
				t.Logf("Error decoding remaining data: %v", err)
			} else {
				t.Logf("Unexpected remaining data: %#v", remaining)
			}
		} else {
			t.Logf("Decoder correctly reports no more data.")
		}
	})
}

// parseJSONStream 辅助函数，尝试从 reader 读取并解析多个 JSON 对象
func parseJSONStream(t *testing.T, description string, r io.Reader) {
	t.Helper()
	decoder := json.NewDecoder(r)
	count := 0
	for {
		var msg json.RawMessage // 使用 RawMessage 避免具体结构 unmarshal 失败
		err := decoder.Decode(&msg)
		if err == io.EOF {
			t.Logf("[%s] Successfully decoded %d JSON object(s). End of stream.", description, count)
			break // 正常结束
		}
		if err != nil {
			// 记录解码错误，包括部分读取的数据（如果有）
			// 注意：NewDecoder 在出错后可能无法准确提供剩余数据
			buf := new(bytes.Buffer)
			// 尝试从 decoder 的 buffered reader 读取剩余数据（可能不准确）
			if br, ok := decoder.Buffered().(*bytes.Reader); ok {
				io.Copy(buf, br)
			} else if br, ok := decoder.Buffered().(io.Reader); ok {
				// 尝试读取，但不保证能拿到准确的出错点数据
				io.Copy(buf, br)
			}

			t.Errorf("[%s] Error decoding JSON object after %d successful decodes: %v. Partially buffered data (approx): %s", description, count, err, buf.String())
			break // 出错结束
		}
		count++
		t.Logf("[%s] Decoded JSON object %d: %s", description, count, string(msg))
	}
}

// 测试验证是否丢包或数据损坏
func TestVerifyPacketLoss(t *testing.T) {
	binPath := "/Users/zhaojian/框架设计/go/Go-Midjourney-Pool/websocket_messages.bin" // 注意：路径需要根据实际情况修改或做成相对路径/参数
	// lastMassageB64 := "7J0xDsIw" // 不再需要 lastMassageB64

	t.Logf("--- 测试: 验证 websocket_messages.bin --- ")
	t.Logf("Log file path: %s", binPath)
	// t.Logf("Last message base64: %s", lastMassageB64) // 移除日志

	// 1. 读取日志文件内容
	logBytes, err := os.ReadFile(binPath)
	if err != nil {
		// 如果文件不存在，根据情况决定是失败还是跳过
		if os.IsNotExist(err) {
			t.Skipf("Log file %s not found, skipping verification.", binPath)
			return
		}
		t.Fatalf("Failed to read log file %s: %v", binPath, err)
	}
	t.Logf("Read %d bytes from log file.", len(logBytes))

	// 2. (移除) 解码最后一个消息块
	// lastMsgBytes, err := base64.StdEncoding.DecodeString(lastMassageB64)
	// ... (相关逻辑已移除)

	// 3. (移除) 合并数据
	// combinedData := append(logBytes, lastMsgBytes...)
	// t.Logf("Combined data length: %d bytes.", len(combinedData))

	// 4. 创建 zlib reader (直接使用 logBytes)
	logReader := bytes.NewReader(logBytes)
	zr, err := zlib.NewReader(logReader)
	if err != nil {
		// 如果这里报错，通常意味着流的开头（第一个块）就有问题
		t.Fatalf("Failed to create zlib reader from log file data: %v", err)
	}
	defer func() {
		if err := zr.Close(); err != nil {
			// zlib 关闭时返回 EOF 是正常的
			if err != io.EOF {
				t.Errorf("zr.Close() failed: %v", err)
			}
		}
	}()
	t.Logf("Created zlib reader successfully.")

	// 5. 读取所有解压数据
	decompressedData, err := io.ReadAll(zr)
	if err != nil {
		// 如果在ReadAll过程中出错，可能是 zlib 流本身损坏 (e.g., checksum error)
		t.Errorf("Failed to read all decompressed data from zlib reader: %v", err)
		// 即使出错，仍然尝试打印已解压的部分数据
		t.Logf("Partially decompressed data before error: %s", string(decompressedData))
		t.Fatalf("Stopping test due to zlib decompression error.") // 强制失败
	}
	t.Logf("Successfully read %d bytes of decompressed data.", len(decompressedData))
	t.Logf("Full decompressed data:\n%s", string(decompressedData)) // 打印完整的解压后数据

	// 6. 尝试解析解压后的数据作为 JSON 流
	t.Logf("Attempting to parse decompressed data as JSON stream...")
	decompressedReader := bytes.NewReader(decompressedData)
	parseJSONStream(t, "Log File Stream JSON Parse", decompressedReader)

	t.Logf("--- 测试完成 ---")
}
