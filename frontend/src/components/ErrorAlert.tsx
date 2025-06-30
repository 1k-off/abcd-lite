import { AlertCircle } from "lucide-react"
import { Button } from "@/components/ui/button"

interface ErrorAlertProps {
  error: string | null
  onDismiss: () => void
}

export function ErrorAlert({ error, onDismiss }: ErrorAlertProps) {
  if (!error) return null
  return (
    <div className="mb-6 p-4 bg-destructive/10 border border-destructive rounded-lg flex items-center gap-2 text-destructive">
      <AlertCircle className="h-5 w-5" />
      <p>{error}</p>
      <Button
        variant="ghost"
        size="sm"
        className="ml-auto"
        onClick={onDismiss}
      >
        Dismiss
      </Button>
    </div>
  )
} 