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
	// 核心逻辑：通用危险操作关键词（增删改操作）
	// 只要命令中包含这些关键词，就触发二次确认
	// =====================================================================

	// ========== 删除类操作 ==========
	regexp.MustCompile(`(?i)\b(delete|del|remove|rm|drop|destroy|kill|stop|terminate|wipe|erase)\b`),
	regexp.MustCompile(`(?i)\b(uninstall|purge|clean|clear|flush|truncate|prune)\b`),

	// ========== 新增类操作 ==========
	regexp.MustCompile(`(?i)\b(create|add|insert|install|write|append|put|make|new|push|upload|import)\b`),

	// ========== 修改类操作 ==========
	regexp.MustCompile(`(?i)\b(update|modify|alter|change|replace|patch|edit|rename|set|apply|commit|pull|fetch|clone|init|build|run|exec)\b`),
	regexp.MustCompile(`(?i)\b(enable|disable|mask|unmask|reload|restart|start|reboot|shutdown|poweroff|halt|reset|rollback|revert|restore|migrate)\b`),

	// =====================================================================
	// 补充：高危系统命令（命令本身就危险，或通用关键词无法覆盖的）
	// =====================================================================

	// ========== 文件操作（mv/cp 不包含通用关键词，需单独匹配） ==========
	regexp.MustCompile(`(?i)^\s*mv\s+`),    // mv 移动/重命名
	regexp.MustCompile(`(?i)^\s*cp\s+`),    // cp 复制（可能覆盖）
	regexp.MustCompile(`(?i)^\s*ln\s+`),    // ln 链接
	regexp.MustCompile(`(?i)^\s*dd\s+`),    // dd 磁盘操作
	regexp.MustCompile(`(?i)^\s*shred\s+`), // shred 安全删除
	regexp.MustCompile(`(?i)^\s*sed\s+-i`), // sed 原地修改
	regexp.MustCompile(`(?i)^\s*tee\s+`),   // tee 写文件

	// ========== 文件内容重定向 ==========
	regexp.MustCompile(`(?i)>\s*\S+`), // 重定向写入文件

	// ========== 系统电源（halt 等已在通用关键词，init 需单独匹配） ==========
	regexp.MustCompile(`(?i)\binit\s+[06]\b`), // init 0/6

	// ========== 磁盘/分区/文件系统（专有命令，通用关键词无法覆盖） ==========
	regexp.MustCompile(`(?i)\b(wipefs|mkfs|parted|fdisk|gdisk|cfdisk|sfdisk)\b`),
	regexp.MustCompile(`(?i)\b(lvcreate|lvremove|lvresize|lvextend|vgremove|vgcreate|pvcreate|pvremove)\b`),
	regexp.MustCompile(`(?i)\b(resize2fs|e2fsck|tune2fs|xfs_growfs|xfs_repair|fsck)\b`),
	regexp.MustCompile(`(?i)^\s*(mount|umount)\s+`),
	regexp.MustCompile(`(?i)^\s*(mkswap|swapoff|swapon)\s+`),
	regexp.MustCompile(`(?i)^\s*(blkdiscard|hdparm)\s+`),
	regexp.MustCompile(`(?i)\b(cryptsetup|mdadm|zfs|btrfs)\b`),

	// ========== 用户/权限（专有命令） ==========
	regexp.MustCompile(`(?i)^\s*(useradd|userdel|usermod|groupadd|groupdel|groupmod)\s+`),
	regexp.MustCompile(`(?i)^\s*passwd\b`),
	regexp.MustCompile(`(?i)^\s*(chown|chmod|chattr|setfacl)\s+`),
	regexp.MustCompile(`(?i)^\s*visudo\b`),
	regexp.MustCompile(`(?i)^\s*(setenforce|aa-enforce|aa-complain|aa-disable)\b`),

	// ========== 网络配置（专有命令） ==========
	regexp.MustCompile(`(?i)^\s*(iptables|ip6tables|nft)\b`),
	regexp.MustCompile(`(?i)^\s*ip\s+(route|addr|link|rule|neigh)\s+`),

	// ========== 内核/系统参数 ==========
	regexp.MustCompile(`(?i)^\s*sysctl\s+-w\b`), // sysctl 写参数
	regexp.MustCompile(`(?i)^\s*(modprobe|rmmod|insmod)\b`),

	// ========== 定时任务 ==========
	regexp.MustCompile(`(?i)^\s*crontab\s+-[re]\b`), // crontab 编辑/删除
	regexp.MustCompile(`(?i)^\s*at\s+`),             // at 定时任务

	// ========== 提权操作 ==========
	regexp.MustCompile(`(?i)\bsudo\b`),
	regexp.MustCompile(`(?i)^\s*su\s+`),
}

func IsDangerousCommandRegex(cmd string) bool {
	for _, pattern := range dangerousPatterns {
		if pattern.MatchString(cmd) {
			return true
		}
	}
	return false
}
