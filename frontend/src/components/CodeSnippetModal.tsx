import { Dialog } from "./Dialog"
import { Button } from "@/components/ui/button"
import { useState } from "react"

interface CodeSnippetModalProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  type?: string
  code?: string
}

export function CodeSnippetModal({ open, onOpenChange, type, code }: CodeSnippetModalProps) {
  const [copied, setCopied] = useState(false)

  const handleCopy = () => {
    if (code) {
      navigator.clipboard.writeText(code)
      setCopied(true)
      setTimeout(() => setCopied(false), 1200)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange} title={type || "Code snippet"}>
      <div className="space-y-4">
        {code ? (
          <div className="relative">
            <pre className="bg-muted rounded p-4 overflow-x-auto text-sm font-mono whitespace-pre-wrap"><code>{code}</code></pre>
            <Button
              type="button"
              size="sm"
              variant="outline"
              className="absolute top-2 right-2"
              onClick={handleCopy}
            >
              {copied ? "Copied!" : "Copy"}
            </Button>
          </div>
        ) : (
          <div className="text-center text-muted-foreground py-8">Under development</div>
        )}
      </div>
    </Dialog>
  )
} 