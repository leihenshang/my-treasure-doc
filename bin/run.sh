#!/usr/bin/env bash
# treasure-doc 前后端服务管理脚本
# 用法:
#   ./run.sh start [config_path]   启动后端（默认使用 module/user/config.toml）
#   ./run.sh stop                  停止后端
#   ./run.sh status                查看后端状态
#   ./run.sh restart [config_path] 重启后端
#   ./run.sh front-stop           停止前端 vite 服务（pnpm dev，端口 2024）
#   ./run.sh help                  帮助
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
USER_DIR="$ROOT_DIR/module/user"

# 前端项目目录（与后端仓库同级）及其默认 dev 端口
FRONT_DIR="$(cd "$ROOT_DIR/.." && pwd)/$(basename "$ROOT_DIR")-front"
FRONT_PORT=2024

RUNTIME_DIR="$ROOT_DIR/bin/_runtime"
ACT_BIN="$RUNTIME_DIR/treasure_user"
PID_FILE="$RUNTIME_DIR/treasure-doc.pid"
LOG_FILE="$RUNTIME_DIR/treasure-doc.log"

DEFAULT_CONFIG="$USER_DIR/config.toml"

resolve_config() {
    local cfg="${1:-$DEFAULT_CONFIG}"
    # 允许相对路径，基于仓库根解析
    if [[ "$cfg" != /* ]]; then
        cfg="$ROOT_DIR/$cfg"
    fi
    if [[ ! -f "$cfg" ]]; then
        echo "[ERROR] 配置文件不存在: $cfg" >&2
        exit 1
    fi
    CONFIG="$(cd "$(dirname "$cfg")" && pwd)/$(basename "$cfg")"
}

build() {
    mkdir -p "$RUNTIME_DIR"
    echo "[build] 编译后端二进制..."
    (cd "$USER_DIR" && go build -trimpath -o "$ACT_BIN" .)
}

get_pid() {
    [[ -f "$PID_FILE" ]] && cat "$PID_FILE" || echo ""
}

is_running() {
    local pid; pid="$(get_pid)"
    [[ -n "$pid" ]] && kill -0 "$pid" 2>/dev/null
}

# 从配置文件中解析 app.port（兼容带/不带引号的写法）
config_port() {
    local cfg="${1:-$DEFAULT_CONFIG}"
    [[ -f "$cfg" ]] || { echo ""; return; }
    grep -E '^[[:space:]]*port[[:space:]]*=' "$cfg" 2>/dev/null | head -n1 | grep -oE '[0-9]+' | head -n1
}

do_start() {
    local config_path="${1:-}"
    resolve_config "$config_path"
    if is_running; then
        echo "[start] 服务已在运行 (PID $(get_pid))，如需重启请用 restart"
        return 0
    fi
    build
    echo "[start] 使用配置: $CONFIG"
    # 必须从 module/user 目录运行（配置 dsn/files 等用相对路径）
    # 直接后台启动二进制并用 disown 脱离作业控制；不经过 nohup，避免某些
    # nohup 实现把二进制 fork 成子进程，导致记录的 PID 不是真正监听端口的进程。
    ( cd "$USER_DIR" || exit 1; "$ACT_BIN" -c "$CONFIG" >>"$LOG_FILE" 2>&1 & echo $! > "$PID_FILE"; disown )
    sleep 1
    if is_running; then
        echo "[start] 已启动，PID $(get_pid)"
    else
        echo "[start] 启动后进程未存活，请查看日志: $LOG_FILE"
        tail -n 30 "$LOG_FILE" 2>/dev/null || true
        return 1
    fi
    echo "[start] 日志: $LOG_FILE"
}

do_stop() {
    local pid; pid="$(get_pid)"
    if [[ -z "$pid" ]] || ! kill -0 "$pid" 2>/dev/null; then
        echo "[stop] 服务未在运行"
        rm -f "$PID_FILE"
        return 0
    fi
    echo "[stop] 正在停止 PID $pid ..."
    kill "$pid" 2>/dev/null || true
    for _ in $(seq 1 25); do
        kill -0 "$pid" 2>/dev/null || break
        sleep 0.2
    done
    if kill -0 "$pid" 2>/dev/null; then
        echo "[stop] 强制终止 PID $pid"
        kill -9 "$pid" 2>/dev/null || true
    fi
    rm -f "$PID_FILE"
    echo "[stop] 已停止"
}

do_status() {
    local pid; pid="$(get_pid)"
    if is_running; then
        local port; port="$(config_port "$DEFAULT_CONFIG")"
        echo "[status] 运行中"
        echo "  进程 PID : $pid"
        echo "  运行时长 : $(ps -o etime= -p "$pid" 2>/dev/null | tr -d ' ')"
        echo "  监听端口 : ${port:-未知} (HTTP 探测: $(curl -s -o /dev/null -m 2 "http://127.0.0.1:${port:-0}" && echo 可达 || echo 未响应))"
        echo "  日志文件 : $LOG_FILE"
    else
        echo "[status] 服务未运行"
    fi
}

# 解析前端配置：提取 vite.config.(ts|js) 里的 server.port
front_config_port() {
    local cfg
    cfg="$FRONT_DIR/vite.config.ts"
    [[ -f "$cfg" ]] || cfg="$FRONT_DIR/vite.config.js"
    if [[ -f "$cfg" ]]; then
        grep -oE 'port[[:space:]]*:[[:space:]]*[0-9]+' "$cfg" | grep -oE '[0-9]+' | head -n1
    fi
}

# 定位监听指定端口的 PID（兼容 ss / lsof / fuser）
port_pids() {
    local port="$1"
    if command -v ss >/dev/null 2>&1; then
        ss -tlnp 2>/dev/null | awk -v p=":$port$" '$4 ~ p { for (i=1;i<=NF;i++) if ($i ~ /pid=/){ split($i,a,","); for(x in a) if (a[x] ~ /^pid=/){ gsub(/pid=/,"",a[x]); print a[x] } } }'
    fi
}

do_front_stop() {
    local port
    port="$(front_config_port)"
    [[ -n "$port" ]] || port="$FRONT_PORT"
    if [[ ! -d "$FRONT_DIR" ]]; then
        echo "[front-stop] 前端目录不存在: $FRONT_DIR" >&2
        exit 1
    fi
    local pids pid
    pids="$(port_pids "$port")"
    if [[ -z "$pids" ]]; then
        echo "[front-stop] 前端服务未在运行（端口 $port 无监听进程）"
        return 0
    fi
    echo "[front-stop] 正在停止前端服务（端口 $port）..."
    for pid in $pids; do
        echo "[front-stop] 终止 PID $pid"
        kill "$pid" 2>/dev/null || true
    done
    # 等待进程退出，超时后强制终止
    for _ in $(seq 1 25); do
        if [[ -z "$(port_pids "$port")" ]]; then break; fi
        sleep 0.2
    done
    local remaining
    remaining="$(port_pids "$port")"
    if [[ -n "$remaining" ]]; then
        echo "[front-stop] 强制终止残留进程..."
        for pid in $remaining; do
            kill -9 "$pid" 2>/dev/null || true
        done
    fi
    echo "[front-stop] 前端服务已停止"
}

case "${1:-help}" in
    start)   do_start "${2:-}";;
    stop)    do_stop;;
    status)  do_status;;
    restart) do_stop; do_start "${2:-}";;
    front-stop) do_front_stop;;
    help|-h|--help)
        sed -n '2,10p' "${BASH_SOURCE[0]}"
        ;;
    *)
        echo "用法错误，可选命令: start | stop | status | restart | front-stop" >&2
        exit 1
        ;;
esac
