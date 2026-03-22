import { useState } from "react"
import { Trash2, AlertTriangle } from "lucide-react"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { Button } from "@/components/ui/button"
import { factoryReset } from "@/api/provisioning"

interface FactoryResetDialogProps {
  onReset?: () => void
  trigger?: React.ReactNode
}

export function FactoryResetDialog({ onReset, trigger }: FactoryResetDialogProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [open, setOpen] = useState(false)

  const handleReset = async () => {
    setIsLoading(true)
    try {
      await factoryReset()
      setOpen(false)
      onReset?.()
      // The device will reset and the page will likely become unavailable
      // as the hotspot is re-enabled
    } catch (error) {
      console.error("Factory reset failed:", error)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <AlertDialog open={open} onOpenChange={setOpen}>
      <AlertDialogTrigger asChild>
        {trigger || (
          <Button variant="destructive" size="sm">
            <Trash2 className="size-4" />
            Factory Reset
          </Button>
        )}
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <div className="flex items-center gap-3 mb-2">
            <div className="flex size-10 items-center justify-center rounded-full bg-destructive/10">
              <AlertTriangle className="size-5 text-destructive" />
            </div>
            <AlertDialogTitle>Factory Reset</AlertDialogTitle>
          </div>
          <AlertDialogDescription className="text-left">
            This will <strong>permanently reset</strong> the device to its initial state:
          </AlertDialogDescription>
          <ul className="text-sm text-muted-foreground list-disc list-inside space-y-1 mt-2">
            <li>All saved WiFi networks will be forgotten</li>
            <li>All device configuration will be cleared</li>
            <li>The device will return to provisioning mode</li>
            <li>You will need to set up the device again</li>
          </ul>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel disabled={isLoading}>Cancel</AlertDialogCancel>
          <AlertDialogAction
            onClick={handleReset}
            disabled={isLoading}
            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
          >
            {isLoading ? "Resetting..." : "Reset Device"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}
