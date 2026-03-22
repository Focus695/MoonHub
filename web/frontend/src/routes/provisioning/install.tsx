import { createFileRoute } from "@tanstack/react-router"
import { useState } from "react"

import { Button } from "@/components/ui/button"

export const Route = createFileRoute("/provisioning/install")({
  component: ProvisioningInstall,
})

function ProvisioningInstall() {
  const [isChecking, setIsChecking] = useState(false)
  const [isInstalled, setIsInstalled] = useState(false)

  const handleCheckInstall = async () => {
    setIsChecking(true)
    // Simulate install check
    await new Promise((resolve) => setTimeout(resolve, 1500))
    // Check if running as PWA
    const isStandalone = window.matchMedia("(display-mode: standalone)").matches
    setIsInstalled(isStandalone)
    setIsChecking(false)
  }

  const handleStart = () => {
    // In a real implementation, this would launch the PWA
    window.location.reload()
  }

  return (
    <main className="flex-grow flex flex-col items-center px-6 pt-8 pb-32 max-w-lg mx-auto w-full">
      {/* Hero Section */}
      <section className="text-center mb-12">
        <div className="mb-6 inline-flex p-4 rounded-full bg-[#d4e4f7]/30">
          <span className="material-symbols-outlined text-[#506070] text-4xl" style={{ fontVariationSettings: "'FILL' 0" }}>install_mobile</span>
        </div>
        <h2 className="text-3xl font-light tracking-tight text-[#506070] mb-4">添加月枢到桌面</h2>
        <p className="text-[#586064] text-sm leading-relaxed max-w-xs mx-auto">
          安装后即可解锁完整的设备发现与控制功能，享受更纯净的交互体验。
        </p>
      </section>

      {/* Bento Instructions Grid */}
      <div className="grid grid-cols-1 gap-6 w-full">
        {/* iOS Instructions */}
        <div className="bg-white/40 backdrop-blur-xl border border-white/40 rounded-3xl p-6 shadow-[0_0_40px_5px_rgba(80,96,112,0.05)] flex flex-col gap-6">
          <div className="flex items-center gap-3">
            <span className="material-symbols-outlined text-[#506070]/60 text-xl">phone_iphone</span>
            <span className="text-xs uppercase tracking-[0.2em] font-medium text-[#737c7f]">iOS 指引</span>
          </div>
          <div className="flex items-start gap-4">
            <div className="flex-grow space-y-4">
              <div className="flex items-center gap-4">
                <span className="flex-shrink-0 w-6 h-6 rounded-full bg-[#506070] text-white text-[10px] flex items-center justify-center font-bold">01</span>
                <p className="text-sm text-[#2b3437]">
                  点击浏览器底部的{" "}
                  <span className="inline-flex items-center px-1.5 py-0.5 rounded bg-[#e3e9ec] border border-[#abb3b7]/20">
                    <span className="material-symbols-outlined text-sm">ios_share</span>
                  </span>
                  {" 分享图标"}
                </p>
              </div>
              <div className="flex items-center gap-4">
                <span className="flex-shrink-0 w-6 h-6 rounded-full bg-[#506070] text-white text-[10px] flex items-center justify-center font-bold">02</span>
                <p className="text-sm text-[#2b3437]">
                  向上滑动并选择 <span className="font-medium text-[#506070]">"添加到主屏幕"</span>
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Android Instructions */}
        <div className="bg-white/40 backdrop-blur-xl border border-white/40 rounded-3xl p-6 shadow-[0_0_40px_5px_rgba(80,96,112,0.05)] flex flex-col gap-6">
          <div className="flex items-center gap-3">
            <span className="material-symbols-outlined text-[#506070]/60 text-xl">phone_android</span>
            <span className="text-xs uppercase tracking-[0.2em] font-medium text-[#737c7f]">Android 指引</span>
          </div>
          <div className="flex items-start gap-4">
            <div className="flex-grow space-y-4">
              <div className="flex items-center gap-4">
                <span className="flex-shrink-0 w-6 h-6 rounded-full bg-[#506070] text-white text-[10px] flex items-center justify-center font-bold">01</span>
                <p className="text-sm text-[#2b3437]">
                  点击右上角或底部的{" "}
                  <span className="inline-flex items-center px-1.5 py-0.5 rounded bg-[#e3e9ec] border border-[#abb3b7]/20">
                    <span className="material-symbols-outlined text-sm">more_vert</span>
                  </span>
                  {" 菜单图标"}
                </p>
              </div>
              <div className="flex items-center gap-4">
                <span className="flex-shrink-0 w-6 h-6 rounded-full bg-[#506070] text-white text-[10px] flex items-center justify-center font-bold">02</span>
                <p className="text-sm text-[#2b3437]">
                  在菜单中点击 <span className="font-medium text-[#506070]">"安装应用"</span>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Status Check Action */}
      <button
        onClick={handleCheckInstall}
        disabled={isChecking}
        className="mt-10 group flex items-center gap-3 px-6 py-3 rounded-full border border-[#506070]/20 hover:bg-[#506070]/5 transition-all duration-300 active:scale-95 disabled:opacity-50"
      >
        <span className={`material-symbols-outlined text-[#506070] text-xl ${isChecking ? "animate-spin" : ""}`}>
          {isChecking ? "sync" : "search_check"}
        </span>
        <span className="text-sm font-medium text-[#506070]">
          {isChecking ? "检测中..." : isInstalled ? "已安装 ✓" : "检测安装状态"}
        </span>
      </button>

      {/* Ambient Glow Background Decor */}
      <div className="fixed top-0 left-0 w-full h-full pointer-events-none -z-10 overflow-hidden">
        <div className="absolute top-[-10%] right-[-10%] w-[50%] h-[50%] rounded-full bg-[#506070]/5 blur-[120px]" />
        <div className="absolute bottom-[-5%] left-[-5%] w-[40%] h-[40%] rounded-full bg-[#d4e4f7]/10 blur-[100px]" />
      </div>

      {/* Bottom Actions */}
      <div className="fixed bottom-0 left-0 w-full z-50 flex flex-col items-center">
        <div className="px-6 pb-6 w-full max-w-lg">
          <Button
            onClick={handleStart}
            className="w-full h-14 bg-gradient-to-br from-[#506070] to-[#708090] text-white rounded-2xl shadow-lg shadow-[#506070]/20 flex items-center justify-center gap-2 hover:opacity-90 active:scale-[0.98] transition-all duration-300"
          >
            <span className="font-medium tracking-wide">我已安装，开始使用</span>
            <span className="material-symbols-outlined text-xl">arrow_forward</span>
          </Button>
        </div>
      </div>
    </main>
  )
}
