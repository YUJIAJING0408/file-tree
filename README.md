# file-tree
This is a Go library that allows for quick traversal of folders. By inputting a valid file path, you will receive a file tree. Naturally, the final file tree supports output in Json and Yaml formats.
***
1.
增量快照 & diff（实用★★★★☆ 成本★★★☆☆ 炫技★★★☆☆）
把每次扫描结果写进 bolt/badger 里做“快照”，下次再扫同一路径时直接 diff，输出“新增 / 删除 / 大小变化”列表。
关键词：文件指纹（xxhash64+mtime）、快照 ID、boltdb bucket。
2.
重复文件猎手（实用★★★★★ 成本★★★☆☆ 炫技★★☆☆☆）
用“大小→首512B→全文件 xxHash”三级过滤，把重复项聚合成组，支持交互式或 dry-run 删除/硬链接。
关键词：并发 hash、map[uint64][]string、hardlink syscall。
3.
热图可视化（实用★★★☆☆ 成本★★☆☆☆ 炫技★★★★☆）
扫描完把结果直接吐成 HTML，用 d3.js 画“可缩放旭日图”，颜色按文件类型/大小梯度。
关键词：template/html、embed FS、WebSocket 实时推送进度条。
4.
忽略语法 & .gitignore 兼容（实用★★★★☆ 成本★☆☆☆☆ 炫技☆☆☆☆☆）
把 git 的 glob 解析抄过来（github.com/go-git/go-git/v5/plumbing/gitignore），默认读取 .gitignore，支持 CLI 额外 -i 参数。
关键词：双栈匹配、预编译 glob。
5.
并发复制 / 镜像（实用★★★★☆ 成本★★☆☆☆ 炫技★★☆☆☆）
在现有 walker 基础上加“copier”插件，支持 --mirror src dst，用 go-chans 做生产者-消费者，带进度条、限速、断点续传（记录已复制文件列表）。
关键词：io.CopyBuffer、pb/v3、sparse file detect。
6.
大文件 TopN + 老化策略（实用★★★★☆ 成本★☆☆☆☆ 炫技☆☆☆☆☆）
扫完直接出“最大的 50 个文件”列表，支持 --older-than 30d --delete-dry-run 一键清理老化大文件。
关键词：最小堆、os.Stat、humanize.Time。
7.
插件系统（实用★★★☆☆ 成本★★★★☆ 炫技★★★★★）
用 go-plugin 或 yaegi 做脚本钩子，让用户写 lua/starlight 脚本自定义“文件名→权重”逻辑，实时影响树排序/过滤。
关键词：plugin.Open、reflect、lua.State。
8.
远程扫描 agent（实用★★★☆☆ 成本★★★☆☆ 炫技★★★☆☆）
把 walker 拆成 gRPC service，本地 CLI 只负责展示，agent 跑在 NAS／容器里；内网一次性扫 20 台机器。
关键词：grpc+protobuf、stream、tls mutual。
9.
权限修复助手（实用★★★☆☆ 成本★☆☆☆☆ 炫技☆☆☆☆☆）
扫树时记录 os.FileInfo.Mode()，一键把“其他用户可写”或“无执行位”的目录批量修正成 0755/0644，支持 dry-run。
关键词：unix.Perm、walk symlinks。
10.
交互式 TUI（实用★★★☆☆ 成本★★☆☆☆ 炫技★★★☆☆）
用 charm/bubbletea 写全屏 TUI，左右分栏：左边树，右边实时预览 / hex/文本；支持 vim 键位删除、复制路径到剪贴板。
关键词：bubbletea、lipgloss、clipboard。
11.
一键生成 .tar.zst / .zip（实用★★★☆☆ 成本★☆☆☆☆ 炫技☆☆☆☆☆）
扫完直接把过滤后的树打包，支持 zstd 级别 3，自动把时间戳、权限写回。
关键词：archive/tar、github.com/klauspost/compress/zstd。
12.
文件类型识别（实用★★☆☆☆ 成本★☆☆☆☆ 炫技★★☆☆☆）
不看后缀，读文件头 32 B，用 http.DetectContentType+mimetype 库，输出“真实类型 vs 后缀不符”列表。
关键词：magic number、sync.Pool 复用 buf。
13.
硬链 & 软链图谱（实用★★☆☆☆ 成本★★☆☆☆ 炫技★★☆☆☆）
统计 inode nlink>1 的硬链组，画“同一 inode 文件列表”；软链则检测是否 dangling。
关键词：sys.Stat_t.Ino、map[uint64][]string。
14.
监控模式（实用★★★☆☆ 成本★★★☆☆ 炫技★★☆☆☆）
用 fsnotify 监听根目录，树变动时增量更新内存结构，WebSocket 推给前端热图实时刷新。
关键词：fsnotify.Watcher、debounce、hub。
15.
许可证扫描（实用★★☆☆☆ 成本★☆☆☆☆ 炫技★☆☆☆☆）
遍历 node_modules、vendor，把 LICENSE* 文件摘出来，汇总 GPLv3 / BSD / MIT 计数，生成风险报告。
关键词：filepath.Glob、spdx 标识符映射表。