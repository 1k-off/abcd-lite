import { Project } from "@/types/project"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Pencil, Trash2, Github, Terminal } from "lucide-react"
import { useState } from "react"
import { CodeSnippetModal } from "./CodeSnippetModal"

interface ProjectCardProps {
  project: Project
  onEdit: (project: Project) => void
  onDelete: (id: string) => void
}

export function ProjectCard({ project, onEdit, onDelete }: ProjectCardProps) {
  const [snippet, setSnippet] = useState<null | { type: string; code?: string }>(null)
  const currentUrl = window.location.origin

  const handleShowSnippet = async (type: "github" | "azure" | "bash" | "powershell", site: string) => {
    if (type === "github") {
      const res = await fetch("/snippets/github-actions.yml")
      let template = await res.text()
      template = template
        .replace(/{{ABCDLITE_URL}}/g, currentUrl)
        .replace(/{{PROJECT_ID}}/g, project.id)
        .replace(/{{SITE_NAME}}/g, site)
      setSnippet({
        type: "GitHub Actions",
        code: template
      })
    } else if (type === "azure") {
      const res = await fetch("/snippets/azure-devops.yml")
      let template = await res.text()
      template = template
        .replace(/{{ABCDLITE_URL}}/g, currentUrl)
        .replace(/{{PROJECT_ID}}/g, project.id)
        .replace(/{{SITE_NAME}}/g, site)
      setSnippet({
        type: "Azure DevOps",
        code: template
      })
    } else if (type === "bash") {
      const res = await fetch("/snippets/generic-bash.sh")
      let template = await res.text()
      template = template
        .replace(/{{ABCDLITE_URL}}/g, currentUrl)
        .replace(/{{PROJECT_ID}}/g, project.id)
        .replace(/{{SITE_NAME}}/g, site)
        .replace(/{{DEPLOY_KEY}}/g, "")
        .replace(/{{PACKAGE_REF}}/g, "")
      setSnippet({ type: "Generic (bash)", code: template })
    } else if (type === "powershell") {
      const res = await fetch("/snippets/generic-powershell.ps1")
      let template = await res.text()
      template = template
        .replace(/{{ABCDLITE_URL}}/g, currentUrl)
        .replace(/{{PROJECT_ID}}/g, project.id)
        .replace(/{{SITE_NAME}}/g, site)
        .replace(/{{DEPLOY_KEY}}/g, "")
        .replace(/{{PACKAGE_REF}}/g, "")
      setSnippet({ type: "Generic (powershell)", code: template })
    }
  }

  return (
    <Card>
      <CardHeader>
        <div className="flex items-start justify-between">
          <div>
            <CardTitle>{project.name}</CardTitle>
            <CardDescription>ID: {project.id}</CardDescription>
          </div>
          <div className="flex space-x-2">
            <Button
              variant="ghost"
              size="icon"
              onClick={() => onEdit(project)}
              title="Edit project"
            >
              <Pencil className="h-4 w-4" />
            </Button>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => onDelete(project.id)}
              title="Delete project"
              className="text-destructive hover:text-destructive"
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <h4 className="text-sm font-medium mb-2">IIS Sites</h4>
          {project.iisSites?.length > 0 ? (
            <ul className="space-y-1">
              {project.iisSites.map((site, index) => (
                <li key={`${project.id}-iis-${site}-${index}`} className="text-sm text-muted-foreground flex items-center gap-2">
                  <span>{site}</span>
                  <span className="flex items-center gap-1 ml-2">
                    {/* GitHub Icon */}
                    <span title="GitHub Actions" onClick={() => handleShowSnippet("github", site)}>
                      <Github className="h-4 w-4 cursor-pointer hover:text-primary" />
                    </span>
                    {/* Azure Icon (inline SVG) */}
                    <span title="Azure DevOps" onClick={() => handleShowSnippet("azure", site)}>
                      <svg className="h-4 w-4 cursor-pointer hover:text-primary" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <polygon points="2,26 16,2 30,26 16,30" fill="#0078D4" />
                        <polygon points="16,2 30,26 16,22" fill="#50E6FF" />
                      </svg>
                    </span>
                    {/* Bash Icon (Terminal) */}
                    <span title="Generic (bash)" onClick={() => handleShowSnippet("bash", site)}>
                      <Terminal className="h-4 w-4 cursor-pointer hover:text-primary" />
                    </span>
                    {/* PowerShell Icon (inline SVG) */}
                    <span title="Generic (powershell)" onClick={() => handleShowSnippet("powershell", site)}>
                      <svg className="h-4 w-4 cursor-pointer hover:text-primary" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
                        <rect x="4" y="6" width="24" height="20" rx="2" fill="#012456" />
                        <path d="M10 20L18 12" stroke="#39A9DC" strokeWidth="2" strokeLinecap="round" />
                        <path d="M14 22H20" stroke="#39A9DC" strokeWidth="2" strokeLinecap="round" />
                      </svg>
                    </span>
                  </span>
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-sm text-muted-foreground">No IIS sites configured</p>
          )}
        </div>
        <div>
          <h4 className="text-sm font-medium mb-2">API Keys</h4>
          {project.apiKeys?.length > 0 ? (
            <ul className="space-y-1">
              {project.apiKeys.map((key, index) => (
                <li key={`${project.id}-api-${key.hash}-${index}`} className="text-sm text-muted-foreground font-mono">
                  {key.prefix + '****' + key.suffix}
                  <span className="ml-2 text-xs text-muted-foreground">({new Date(key.createdAt).toLocaleString()})</span>
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-sm text-muted-foreground">No API keys configured</p>
          )}
        </div>
      </CardContent>
      <CardFooter className="flex flex-col items-start space-y-1 border-t pt-4">
        <p className="text-xs text-muted-foreground">
          Created: {new Date(project.createdAt).toLocaleString()}
        </p>
        <p className="text-xs text-muted-foreground">
          Updated: {new Date(project.updatedAt).toLocaleString()}
        </p>
      </CardFooter>
      <CodeSnippetModal open={!!snippet} onOpenChange={() => setSnippet(null)} type={snippet?.type} code={snippet?.code} />
    </Card>
  )
} 