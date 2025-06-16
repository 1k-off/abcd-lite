import { Project } from "@/types/project"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Pencil, Trash2 } from "lucide-react"

interface ProjectCardProps {
  project: Project
  onEdit: (project: Project) => void
  onDelete: (id: string) => void
}

export function ProjectCard({ project, onEdit, onDelete }: ProjectCardProps) {
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
                <li key={`${project.id}-iis-${site}-${index}`} className="text-sm text-muted-foreground">
                  {site}
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
                <li key={`${project.id}-api-${key}-${index}`} className="text-sm text-muted-foreground">
                  {key}
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
    </Card>
  )
} 