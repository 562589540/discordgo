package discordgo

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath" // 用于安全地处理文件名
	"sync"
	"time" // 用于生成唯一文件名

	crand "crypto/rand" // 用于生成唯一文件名中的随机部分
	"math/rand"         // 使用 crypto/rand 初始化 math/rand
)

// MirrorZlibAdapter 使用连续流方法处理镜像站场景的 zlib 解压缩。
// 它维护一个用于压缩数据的内部缓冲区和一个用于处理解压缩流的持久 zlib.Reader。
// 它实现了 io.Reader 接口，允许 json.Decoder 直接从中读取解压缩后的字节流。
// 修改：使用唯一日志文件，HandleError尝试从日志恢复，Reset/Close删除日志。
type MirrorZlibAdapter struct {
	buffer      *bytes.Buffer // 用于累积压缩数据的缓冲区。创建时初始化。
	reader      io.ReadCloser // 用于缓冲区的持久 zlib 读取器。在 Read 方法中惰性初始化。
	mutex       sync.Mutex    // 保护缓冲区的写入以及读取器的初始化/访问。
	logFile     *os.File      // 当前打开的日志文件句柄
	logFilename string        // 当前日志文件的完整路径名
}

// 生成一个伪随机字符串
func pseudoUuid() string {
	var seed [8]byte
	_, err := crand.Read(seed[:])
	if err != nil {
		// 如果读取 crypto/rand 失败，回退到时间戳
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	// 使用加密安全的种子初始化 math/rand
	r := rand.New(rand.NewSource(int64(seed[0]) | int64(seed[1])<<8 | int64(seed[2])<<16 | int64(seed[3])<<24 |
		int64(seed[4])<<32 | int64(seed[5])<<40 | int64(seed[6])<<48 | int64(seed[7])<<56))
	b := make([]byte, 16)
	r.Read(b)
	return fmt.Sprintf("%x", b)
}

func getLogFileName() string {
	// 生成唯一文件名
	filename := fmt.Sprintf("ws_log_%s_%s.bin",
		time.Now().Format("20060102_150405"),
		pseudoUuid(),
	)
	// 获取绝对路径可能更安全
	absFilename, err := filepath.Abs(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 无法获取日志文件 %s 的绝对路径: %v\n", filename, err)
		absFilename = filename // 回退到相对路径
	}
	return absFilename
}

// NewMirrorZlibAdapter 创建一个新的、空的适配器，并创建唯一的日志文件。
func NewMirrorZlibAdapter() *MirrorZlibAdapter {
	absFilename := getLogFileName()
	// 创建并打开唯一日志文件（覆盖模式）
	f, err := os.OpenFile(absFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 无法创建或打开 websocket 消息日志文件 %s: %v\n", absFilename, err)
		f = nil
	}

	return &MirrorZlibAdapter{
		buffer:      new(bytes.Buffer),
		logFile:     f,
		logFilename: absFilename, // 存储完整路径名
	}
}

// AppendMessage 将新的压缩 websocket 消息安全地添加到内部缓冲区，并直接写入日志文件。
func (a *MirrorZlibAdapter) AppendMessage(message []byte) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.logFile == nil {
		return fmt.Errorf("适配器日志文件 %s 已关闭，无法处理新消息", a.logFilename)
	}

	// 1. 写入主累积缓冲区 (用于解压)
	_, err := a.buffer.Write(message)
	if err != nil {
		return fmt.Errorf("写入消息到主缓冲区失败: %w", err)
	}

	// 2. 直接将当前消息写入日志文件
	_, writeErr := a.logFile.Write(message)
	if writeErr != nil {
		fmt.Fprintf(os.Stderr, "警告: 写入 websocket 消息到日志文件 %s 失败: %v\n", a.logFilename, writeErr)
		_ = a.logFile.Close() // 尝试关闭文件
		a.logFile = nil
		return fmt.Errorf("写入日志文件 %s 失败: %w", a.logFilename, writeErr)
	}

	return nil
}

// Read 实现了 io.Reader 接口。
func (a *MirrorZlibAdapter) Read(p []byte) (n int, err error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.reader == nil {
		if a.buffer.Len() == 0 {
			return 0, io.EOF
		}
		zr, initErr := zlib.NewReader(a.buffer)
		if initErr != nil {
			return 0, fmt.Errorf("初始化 zlib 读取器失败: %w", initErr)
		}
		a.reader = zr
	}

	n, err = a.reader.Read(p)

	return n, err
}

// closeAndRemoveLogFile 是一个内部辅助函数，用于关闭并删除日志文件。
func (a *MirrorZlibAdapter) closeAndRemoveLogFile() error {
	var combinedErr error
	if a.logFile != nil {
		filename := a.logFilename // 先保存文件名
		if err := a.logFile.Close(); err != nil {
			combinedErr = fmt.Errorf("关闭日志文件 %s 失败: %w", filename, err)
		}
		// 尝试删除文件
		if removeErr := os.Remove(filename); removeErr != nil {
			if combinedErr != nil {
				combinedErr = fmt.Errorf("%v; 同时删除日志文件 %s 失败: %w", combinedErr, filename, removeErr)
			} else {
				combinedErr = fmt.Errorf("删除日志文件 %s 失败: %w", filename, removeErr)
			}
		}
	} else if a.logFilename != "" { // 文件句柄是 nil，但文件名存在，尝试删除
		if removeErr := os.Remove(a.logFilename); removeErr != nil {
			combinedErr = fmt.Errorf("删除日志文件 %s 失败: %w", a.logFilename, removeErr)
		}
	}
	a.logFile = nil    // 无论关闭是否成功，都标记为 nil
	a.logFilename = "" // 清除文件名
	return combinedErr
}

// Close 关闭 reader，关闭并删除日志文件，重置 buffer。
func (a *MirrorZlibAdapter) Close() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	var closeErrs []error

	if a.reader != nil {
		if err := a.reader.Close(); err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			closeErrs = append(closeErrs, fmt.Errorf("关闭 zlib reader 失败: %w", err))
		}
		a.reader = nil
	}

	if err := a.closeAndRemoveLogFile(); err != nil {
		closeErrs = append(closeErrs, err)
	}

	a.buffer.Reset()

	if len(closeErrs) > 0 {
		return fmt.Errorf("关闭适配器时发生错误: %v", closeErrs)
	}

	return nil
}

// Reset 关闭 reader，关闭并删除日志文件，重置 buffer。
func (a *MirrorZlibAdapter) Reset() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	var resetErrs []error

	if a.reader != nil {
		if err := a.reader.Close(); err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			resetErrs = append(resetErrs, fmt.Errorf("关闭 zlib reader 失败: %w", err))
		}
		a.reader = nil
	}

	a.buffer.Reset()

	if err := a.closeAndRemoveLogFile(); err != nil {
		resetErrs = append(resetErrs, err)
	}

	a.logFilename = getLogFileName()
	//需要创建一个新的logfile
	newLogFile, openErr := os.OpenFile(a.logFilename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if openErr != nil {
		resetErrs = append(resetErrs, fmt.Errorf("无法重新打开日志文件 %s 进行追加: %w", a.logFilename, openErr))
	}

	if len(resetErrs) > 0 {
		return fmt.Errorf("重置适配器时发生错误: %v", resetErrs)
	}

	a.logFile = newLogFile
	fmt.Println("Reset: successfully created new log file", a.logFilename)
	return nil
}

// HandleError 尝试从日志文件恢复状态，并重新打开日志文件以继续处理。
func (a *MirrorZlibAdapter) HandleError() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	filename := a.logFilename // 保存文件名
	fmt.Println("HandleError called: attempting to restore state from log file", filename)

	if filename == "" {
		fmt.Println("HandleError: log filename is empty, cannot restore.")
		// 仅重置内存状态
		if a.reader != nil {
			_ = a.reader.Close()
			a.reader = nil
		}
		a.buffer.Reset()
		return errors.New("日志文件名为空，无法从日志恢复")
	}

	// 1. 关闭当前日志文件句柄 (如果还打开着)，确保数据刷盘
	if a.logFile != nil {
		if err := a.logFile.Close(); err != nil {
			// 关闭失败也继续尝试恢复，但记录警告
			fmt.Fprintf(os.Stderr, "警告: HandleError 关闭旧日志文件句柄 %s 时出错: %v\n", filename, err)
		}
		a.logFile = nil // 标记为已关闭，以便后续重新打开
	}

	// 2. 重置内部状态 (reader 和 buffer)
	if a.reader != nil {
		_ = a.reader.Close() // 忽略关闭错误
		a.reader = nil
	}
	a.buffer.Reset()

	// 3. 读取已关闭的日志文件内容
	loggedData, readErr := os.ReadFile(filename)
	if readErr != nil {
		fmt.Fprintf(os.Stderr, "错误: HandleError 无法读取日志文件 %s: %v\n", filename, readErr)
		// 读取失败，buffer 保持为空，无法恢复，后续 AppendMessage 也会因为 logFile=nil 而失败
		return fmt.Errorf("无法读取日志文件 %s 以恢复状态: %w", filename, readErr)
	}

	// 4. 将日志数据写回 buffer
	_, writeErr := a.buffer.Write(loggedData)
	if writeErr != nil {
		fmt.Fprintf(os.Stderr, "错误: HandleError 无法将日志数据写回 buffer: %v\n", writeErr)
		// 写入失败，buffer 状态不确定，后续 AppendMessage 也会因为 logFile=nil 而失败
		return fmt.Errorf("无法将日志数据写回 buffer: %w", writeErr)
	}

	// 5. 重新初始化 zlib 读取器
	if a.buffer.Len() == 0 {
		return io.EOF
	}
	zr, initErr := zlib.NewReader(a.buffer)
	if initErr != nil {
		return fmt.Errorf("初始化 zlib 读取器失败: %w", initErr)
	}
	a.reader = zr

	// 6. 重新打开同一个日志文件，用于追加后续消息
	newLogFile, openErr := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if openErr != nil {
		fmt.Fprintf(os.Stderr, "错误: HandleError 无法重新打开日志文件 %s 进行追加: %v\n", filename, openErr)
		// 重新打开失败，logFile 保持为 nil，后续 AppendMessage 会失败
		return fmt.Errorf("无法重新打开日志文件 %s 进行追加: %w", filename, openErr)
	}

	// 恢复成功，更新 logFile 句柄
	a.logFile = newLogFile

	fmt.Printf("HandleError: successfully restored %d bytes from log file %s into buffer and re-opened log for appending.\n", len(loggedData), filename)
	return nil // 表示尝试恢复并准备好继续处理
}
