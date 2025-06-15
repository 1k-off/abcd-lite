import { Project } from "@/types/project"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Plus, X } from "lucide-react"
import { useState, useEffect } from "react"

interface ProjectFormProps {
  project?: Project
  onSubmit: (data: Omit<Project, "id" | "createdAt" | "updatedAt">) => void
  onCancel: () => void
}

export function ProjectForm({ project, onSubmit, onCancel }: ProjectFormProps) {
  const [formData, setFormData] = useState({
    name: "",
    iisSites: [""],
    apiKeys: [""],
  })

  useEffect(() => {
    if (project) {
      setFormData({
        name: project.name,
        iisSites: project.iisSites.length > 0 ? project.iisSites : [""],
        apiKeys: project.apiKeys.length > 0 ? project.apiKeys : [""],
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

  const handleArrayInputChange = (index: number, value: string, field: "iisSites" | "apiKeys") => {
    setFormData((prev) => ({
      ...prev,
      [field]: prev[field].map((item, i) => (i === index ? value : item)),
    }))
  }

  const handleAddItem = (field: "iisSites" | "apiKeys") => {
    setFormData((prev) => ({
      ...prev,
      [field]: [...prev[field], ""],
    }))
  }

  const handleRemoveItem = (index: number, field: "iisSites" | "apiKeys") => {
    setFormData((prev) => ({
      ...prev,
      [field]: prev[field].filter((_, i) => i !== index),
    }))
  }

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    const data = {
      ...formData,
      iisSites: formData.iisSites.filter((site) => site.trim() !== ""),
      apiKeys: formData.apiKeys.filter((key) => key.trim() !== ""),
    }
    onSubmit(data)
  }

  return (
    <Card>
      <form onSubmit={handleSubmit}>
        <CardHeader>
          <CardTitle>{project ? "Edit Project" : "Create Project"}</CardTitle>
          <CardDescription>
            {project ? "Update your project details" : "Add a new project to your list"}
          </CardDescription>
        </CardHeader>
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
              <div key={index} className="flex gap-2 mb-2">
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
            {formData.apiKeys.map((key, index) => (
              <div key={index} className="flex gap-2 mb-2">
                <Input
                  value={key}
                  onChange={(e) => handleArrayInputChange(index, e.target.value, "apiKeys")}
                  placeholder="Enter API key"
                />
                {formData.apiKeys.length > 1 && (
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={() => handleRemoveItem(index, "apiKeys")}
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
              onClick={() => handleAddItem("apiKeys")}
              className="mt-2"
            >
              <Plus className="h-4 w-4 mr-1" />
              Add API Key
            </Button>
          </div>
        </CardContent>
        <CardFooter className="flex justify-end space-x-3">
          <Button type="button" variant="outline" onClick={onCancel}>
            Cancel
          </Button>
          <Button type="submit">{project ? "Update Project" : "Create Project"}</Button>
        </CardFooter>
      </form>
    </Card>
  )
} 