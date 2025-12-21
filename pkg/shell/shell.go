package shell

import (
	"regexp"
	"unicode"
)
import "strings"

// IsCommand 判断是否为命令行表达式
func IsCommand(input string) bool {
	input = strings.TrimSpace(input)

	// 如果包含中文，则认为是自然语言
	for _, r := range input {
		if unicode.Is(unicode.Han, r) {
			return false
		}
	}

	// 匹配命令结构（英文+符号组合）
	// 允许的字符集：字母数字空格和典型命令行符号
	validCmdPattern := regexp.MustCompile(`^[a-zA-Z0-9\s/\-_.|&><=~:"'\\]+$`)
	return validCmdPattern.MatchString(input)
}

var dangerousPatterns = []*regexp.Regexp{
	// =====================================================================
	// 核心逻辑：只有"删除"和"修改"是危险操作
	// "查询"和"新增"不需要确认
	// =====================================================================

	// ========== 删除类操作（危险）==========
	regexp.MustCompile(`(?i)\b(delete|del|remove|rm|drop|destroy|kill|terminate|wipe|erase)\b`),
	regexp.MustCompile(`(?i)\b(uninstall|purge|clean|clear|flush|truncate|prune)\b`),

	// ========== 修改类操作（危险）==========
	regexp.MustCompile(`(?i)\b(update|modify|alter|change|replace|patch|edit|rename)\b`),
	regexp.MustCompile(`(?i)\b(disable|mask|unmask|reload|restart|reboot|shutdown|poweroff|halt)\b`),
	regexp.MustCompile(`(?i)\b(reset|rollback|revert|restore|migrate|format)\b`),

	// ========== 停止服务（危险）==========
	regexp.MustCompile(`(?i)\bstop\b`),

	// =====================================================================
	// 补充：高危系统命令（命令本身就危险）
	// =====================================================================

	// ========== 文件删除/覆盖操作 ==========
	regexp.MustCompile(`(?i)^\s*rm\s+`),    // rm 删除
	regexp.MustCompile(`(?i)^\s*rmdir\s+`), // rmdir 删除目录
	regexp.MustCompile(`(?i)^\s*mv\s+`),    // mv 移动/重命名（可能覆盖）
	regexp.MustCompile(`(?i)^\s*dd\s+`),    // dd 磁盘操作
	regexp.MustCompile(`(?i)^\s*shred\s+`), // shred 安全删除
	regexp.MustCompile(`(?i)^\s*sed\s+-i`), // sed 原地修改
	regexp.MustCompile(`(?i)>\s*/`),        // 重定向到绝对路径（危险）

	// ========== 系统电源 ==========
	regexp.MustCompile(`(?i)\binit\s+[06]\b`), // init 0/6

	// ========== 磁盘/分区操作（修改类）==========
	regexp.MustCompile(`(?i)\b(wipefs|mkfs|parted|fdisk|gdisk)\b`),
	regexp.MustCompile(`(?i)\b(lvremove|lvresize|vgremove|pvremove)\b`),
	regexp.MustCompile(`(?i)\b(resize2fs|e2fsck|tune2fs|xfs_repair|fsck)\b`),
	regexp.MustCompile(`(?i)^\s*umount\s+`),     // 卸载（修改）
	regexp.MustCompile(`(?i)^\s*swapoff\s+`),    // 关闭swap（修改）
	regexp.MustCompile(`(?i)^\s*blkdiscard\s+`), // 磁盘擦除

	// ========== 用户/权限修改 ==========
	regexp.MustCompile(`(?i)^\s*userdel\s+`),              // 删除用户
	regexp.MustCompile(`(?i)^\s*groupdel\s+`),             // 删除组
	regexp.MustCompile(`(?i)^\s*usermod\s+`),              // 修改用户
	regexp.MustCompile(`(?i)^\s*passwd\b`),                // 修改密码
	regexp.MustCompile(`(?i)^\s*(chown|chmod)\s+.*-[rR]`), // 递归修改权限

	// ========== 网络配置修改 ==========
	regexp.MustCompile(`(?i)^\s*iptables\s+.*(-[ADIRF]|--delete|--flush)`),   // 修改防火墙规则
	regexp.MustCompile(`(?i)^\s*ip\s+(route|addr|link)\s+(del|flush|set)\b`), // 删除/修改网络

	// ========== 内核参数修改 ==========
	regexp.MustCompile(`(?i)^\s*sysctl\s+-w\b`), // sysctl 写参数
	regexp.MustCompile(`(?i)^\s*rmmod\b`),       // 卸载内核模块

	// ========== 定时任务删除/修改 ==========
	regexp.MustCompile(`(?i)^\s*crontab\s+-r\b`), // crontab 删除

	// ========== Git 危险操作（会丢失数据）==========
	regexp.MustCompile(`(?i)\bgit\s+reset\s+--hard\b`),        // 重置丢弃更改
	regexp.MustCompile(`(?i)\bgit\s+push\s+.*(-f|--force)\b`), // 强制推送
	regexp.MustCompile(`(?i)\bgit\s+clean\s+.*-f`),            // 清理文件

	// ========== 数据库删除/修改 ==========
	regexp.MustCompile(`(?i)\b(DROP|DELETE|TRUNCATE|ALTER)\s+`), // SQL 危险操作
}

func IsDangerousCommandRegex(cmd string) bool {
	for _, pattern := range dangerousPatterns {
		if pattern.MatchString(cmd) {
			return true
		}
	}
	return false
}
