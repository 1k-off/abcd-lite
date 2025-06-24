import { Project, APIKey } from "@/types/project"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardFooter } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Plus, X } from "lucide-react"
import { useState, useEffect } from "react"
import { generateApiKey, deleteApiKey } from "@/services/api"
import { Dialog } from "@/components/Dialog"

interface ProjectFormProps {
  project?: Project
  onSubmit: (data: Omit<Project, "id" | "createdAt" | "updatedAt">) => void
  onCancel: () => void
  onApiKeyCreated?: (projectId: string, apiKeyMeta: APIKey) => void
}

export function ProjectForm({ project, onSubmit, onCancel, onApiKeyCreated }: ProjectFormProps) {
  const [formData, setFormData] = useState({
    name: "",
    iisSites: [""],
    apiKeys: [] as APIKey[],
  })
  const [showApiKey, setShowApiKey] = useState<string | null>(null)
  const [apiKeyWarningOpen, setApiKeyWarningOpen] = useState(false)
  const [loadingApiKey, setLoadingApiKey] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [deletingKeyId, setDeletingKeyId] = useState<string | null>(null)

  useEffect(() => {
    if (project) {
      setFormData({
        name: project.name,
        iisSites: project.iisSites?.length > 0 ? project.iisSites : [""],
        apiKeys: project.apiKeys || [],
      })
    }
  }, [project])

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }))
  }

  const handleArrayInputChange = (index: number, value: string, field: "iisSites") => {
    setFormData((prev) => ({
      ...prev,
      [field]: prev[field].map((item, i) => (i === index ? value : item)),
    }))
  }

  const handleAddItem = (field: "iisSites") => {
    setFormData((prev) => ({
      ...prev,
      [field]: [...prev[field], ""],
    }))
  }

  const handleRemoveItem = (index: number, field: "iisSites") => {
    setFormData((prev) => ({
      ...prev,
      [field]: prev[field].filter((_, i) => i !== index),
    }))
  }

  const handleGenerateApiKey = async () => {
    if (!project?.id) {
      setError("Project must be created before generating API keys.")
      return
    }
    setLoadingApiKey(true)
    setError(null)
    try {
      const { apiKey, apiKeyMeta } = await generateApiKey(project.id)
      setShowApiKey(apiKey)
      setApiKeyWarningOpen(true)
      setFormData((prev) => ({
        ...prev,
        apiKeys: [...prev.apiKeys, apiKeyMeta],
      }))
      if (typeof onApiKeyCreated === "function") {
        onApiKeyCreated(project.id, apiKeyMeta)
      }
    } catch (err: any) {
      setError(err.message || "Failed to generate API key.")
    } finally {
      setLoadingApiKey(false)
    }
  }

  const handleDeleteApiKey = async (keyId: string) => {
    if (!project?.id) return
    setDeletingKeyId(keyId)
    setError(null)
    try {
      await deleteApiKey(project.id, keyId)
      setFormData((prev) => ({
        ...prev,
        apiKeys: prev.apiKeys.filter((k) => k.id !== keyId),
      }))
    } catch (err: any) {
      setError(err.message || "Failed to delete API key.")
    } finally {
      setDeletingKeyId(null)
    }
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const data = {
      ...formData,
      iisSites: formData.iisSites.filter((site) => site.trim() !== ""),
    }
    onSubmit(data)
  }

  return (
    <Card>
      <form onSubmit={handleSubmit}>
        <CardContent className="space-y-6">
          <div className="space-y-2">
            <Label htmlFor="name">Project Name</Label>
            <Input
              id="name"
              name="name"
              value={formData.name}
              onChange={handleInputChange}
              required
            />
          </div>

          <div>
            <Label className="mb-2 block">IIS Sites</Label>
            {formData.iisSites.map((site, index) => (
              <div key={`iis-input-${index}`} className="flex gap-2 mb-2">
                <Input
                  value={site}
                  onChange={(e) => handleArrayInputChange(index, e.target.value, "iisSites")}
                  placeholder="Enter IIS site"
                />
                {formData.iisSites.length > 1 && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={() => handleRemoveItem(index, "iisSites")}
                    className="text-destructive hover:text-destructive"
                  >
                    <X className="h-4 w-4" />
                  </Button>
                )}
              </div>
            ))}
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={() => handleAddItem("iisSites")}
              className="mt-2"
            >
              <Plus className="h-4 w-4 mr-1" />
              Add IIS Site
            </Button>
          </div>

          <div>
            <Label className="mb-2 block">API Keys</Label>
            <div className="mb-2">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={handleGenerateApiKey}
                disabled={loadingApiKey || !project?.id}
              >
                {loadingApiKey ? "Generating..." : "Create API Key"}
              </Button>
            </div>
            {formData.apiKeys.length === 0 && <div className="text-muted-foreground">No API keys yet.</div>}
            <ul className="space-y-1">
              {formData.apiKeys.map((key) => (
                <li key={key.id} className="font-mono bg-muted rounded px-2 py-1 text-sm flex items-center justify-between">
                  <span>
                    {key.prefix + "****" + key.suffix}
                    <span className="ml-2 text-xs text-muted-foreground">({new Date(key.createdAt).toLocaleString()})</span>
                  </span>
                  {project?.id && (
                    <Button
                      type="button"
                      size="icon"
                      variant="ghost"
                      className="ml-2 text-destructive hover:text-destructive"
                      disabled={deletingKeyId === key.id}
                      onClick={() => handleDeleteApiKey(key.id)}
                      title="Delete API Key"
                    >
                      {deletingKeyId === key.id ? <span className="animate-spin">⏳</span> : <span>🗑️</span>}
                    </Button>
                  )}
                </li>
              ))}
            </ul>
          </div>
          {error && <div className="text-destructive text-sm mt-2">{error}</div>}
        </CardContent>
        <CardFooter className="flex justify-end space-x-3">
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancel
          </Button>
          <Button type="submit">{project ? "Update Project" : "Create Project"}</Button>
        </CardFooter>
      </form>
      <Dialog open={apiKeyWarningOpen} onOpenChange={setApiKeyWarningOpen} title="New API Key" description="Copy this API key now. It will be shown only once!">
        <div className="space-y-2">
          <div className="font-mono bg-muted rounded px-2 py-2 text-lg text-center select-all">
            {showApiKey}
          </div>
          <div className="text-destructive text-sm text-center">This API key will not be shown again. Please copy and store it securely.</div>
          <Button className="w-full mt-2" onClick={() => setApiKeyWarningOpen(false)}>I have copied the key</Button>
        </div>
      </Dialog>
    </Card>
  )
} 