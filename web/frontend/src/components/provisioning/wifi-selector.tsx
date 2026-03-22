import { useState } from "react"

import type { WifiNetwork } from "@/api/provisioning"
import { cn } from "@/lib/utils"

interface WifiSelectorProps {
  networks: WifiNetwork[]
  selectedSsid: string | null
  onSelect: (network: WifiNetwork) => void
  onRefresh: () => void
  isLoading?: boolean
}

function getSignalIcon(signal: number): string {
  if (signal >= 80) return "signal_wifi_4_bar"
  if (signal >= 60) return "signal_wifi_3_bar"
  if (signal >= 40) return "signal_wifi_2_bar"
  if (signal >= 20) return "signal_wifi_1_bar"
  return "signal_wifi_0_bar"
}

export function WifiSelector({
  networks,
  selectedSsid,
  onSelect,
  onRefresh,
  isLoading = false,
}: WifiSelectorProps) {
  const [isOpen, setIsOpen] = useState(false)

  const selectedNetwork = networks.find((n) => n.ssid === selectedSsid)

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-end px-1">
        <label className="text-[11px] uppercase tracking-widest text-[#737c7f]">
          选择网络
        </label>
        <button
          type="button"
          onClick={onRefresh}
          disabled={isLoading}
          className="flex items-center gap-1 text-[#506070] text-xs font-medium hover:opacity-70 transition-opacity disabled:opacity-50"
        >
          <span
            className={cn(
              "material-symbols-outlined text-sm",
              isLoading && "animate-spin"
            )}
          >
            refresh
          </span>
          <span>刷新列表</span>
        </button>
      </div>

      {/* Dropdown Trigger */}
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="w-full group relative bg-[#f1f4f6] hover:bg-[#e3e9ec] transition-colors duration-500 rounded-xl p-4 flex items-center justify-between cursor-pointer border border-white/20"
      >
        <div className="flex items-center gap-3">
          <span className="material-symbols-outlined text-[#455463]">
            {selectedNetwork ? "wifi" : "wifi_tethering"}
          </span>
          <span className="text-[#2b3437] font-medium">
            {selectedNetwork?.ssid || "选择 WiFi 网络"}
          </span>
        </div>
        <div className="flex items-center gap-2">
          {selectedNetwork && (
            <span className="material-symbols-outlined text-[#abb3b7] text-lg">
              {getSignalIcon(selectedNetwork.signal)}
            </span>
          )}
          <span
            className={cn(
              "material-symbols-outlined text-[#abb3b7] text-lg transition-transform",
              isOpen && "rotate-180"
            )}
          >
            expand_more
          </span>
        </div>
      </button>

      {/* Dropdown List */}
      {isOpen && (
        <div className="bg-[#ffffff] rounded-xl border border-white/40 shadow-lg overflow-hidden">
          {networks.length === 0 ? (
            <div className="p-4 text-center text-[#737c7f] text-sm">
              {isLoading ? "扫描中..." : "未找到网络"}
            </div>
          ) : (
            <ul className="max-h-64 overflow-y-auto">
              {networks.map((network) => (
                <li key={network.ssid}>
                  <button
                    type="button"
                    onClick={() => {
                      onSelect(network)
                      setIsOpen(false)
                    }}
                    className={cn(
                      "w-full px-4 py-3 flex items-center justify-between hover:bg-[#f1f4f6] transition-colors",
                      selectedSsid === network.ssid && "bg-[#d4e4f7]/30"
                    )}
                  >
                    <div className="flex items-center gap-3">
                      <span className="material-symbols-outlined text-[#506070]">
                        {getSignalIcon(network.signal)}
                      </span>
                      <span className="text-[#2b3437]">{network.ssid}</span>
                      {network.inUse && (
                        <span className="text-[10px] px-2 py-0.5 bg-[#d4e4f7] text-[#506070] rounded-full">
                          已连接
                        </span>
                      )}
                    </div>
                    {network.security && network.security !== "open" && (
                      <span className="material-symbols-outlined text-[#abb3b7] text-lg">
                        lock
                      </span>
                    )}
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}
    </div>
  )
}
