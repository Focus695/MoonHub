import { useState } from "react"

import { cn } from "@/lib/utils"

interface AuthCodeProps {
  code: string
  deviceId?: string
  onRegenerate?: () => void
  isRegenerating?: boolean
}

export function AuthCode({
  code,
  deviceId,
  onRegenerate,
  isRegenerating = false,
}: AuthCodeProps) {
  const [copied, setCopied] = useState(false)

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(code.replace(/\s/g, ""))
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {
      // Fallback for older browsers
      const textArea = document.createElement("textarea")
      textArea.value = code.replace(/\s/g, "")
      document.body.appendChild(textArea)
      textArea.select()
      document.execCommand("copy")
      document.body.removeChild(textArea)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    }
  }

  // Format code as "XXX XXX" if it's 6 digits
  const formattedCode = code.length === 6
    ? `${code.slice(0, 3)} ${code.slice(3)}`
    : code

  return (
    <div className="space-y-8">
      {/* 6-Digit Auth Code */}
      <div className="py-10 bg-[#ffffff]/40 backdrop-blur-md rounded-[2.5rem] border border-white/40 shadow-sm relative overflow-hidden group">
        <div className="relative z-10">
          <span className="block text-6xl md:text-7xl font-extralight tracking-[0.15em] text-[#506070] text-center">
            {formattedCode}
          </span>
        </div>
        {/* Glass Reflection */}
        <div className="absolute inset-0 bg-gradient-to-tr from-transparent via-white/5 to-transparent pointer-events-none" />
      </div>

      <div className="space-y-6">
        <p className="text-[#586064] font-light text-sm leading-relaxed text-center px-8">
          请妥善保存此授权码，用于连接 PWA
        </p>

        {/* Copy Button */}
        <div className="flex justify-center gap-3">
          <button
            onClick={handleCopy}
            className="inline-flex items-center space-x-2 px-6 py-3 rounded-full bg-[#e3e9ec] text-[#506070] text-sm font-medium hover:bg-[#dbe4e7] transition-colors active:scale-95 duration-200"
          >
            <span className="material-symbols-outlined text-sm">
              {copied ? "check" : "content_copy"}
            </span>
            <span>{copied ? "已复制" : "复制授权码"}</span>
          </button>

          {onRegenerate && (
            <button
              onClick={onRegenerate}
              disabled={isRegenerating}
              className={cn(
                "inline-flex items-center space-x-2 px-6 py-3 rounded-full text-sm font-medium transition-colors active:scale-95 duration-200",
                isRegenerating
                  ? "bg-[#e3e9ec] text-[#abb3b7] cursor-not-allowed"
                  : "bg-[#e3e9ec] text-[#506070] hover:bg-[#dbe4e7]"
              )}
            >
              <span
                className={cn(
                  "material-symbols-outlined text-sm",
                  isRegenerating && "animate-spin"
                )}
              >
                refresh
              </span>
              <span>{isRegenerating ? "重新生成中..." : "重新生成"}</span>
            </button>
          )}
        </div>
      </div>

      {/* Device ID */}
      {deviceId && (
        <div className="flex flex-col items-center space-y-1">
          <span className="text-[10px] text-[#abb3b7] tracking-widest uppercase">
            Device Identity
          </span>
          <span className="text-[#586064] text-xs font-light">{deviceId}</span>
        </div>
      )}
    </div>
  )
}
