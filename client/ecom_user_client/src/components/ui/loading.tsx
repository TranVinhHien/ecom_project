import { cn } from "@/lib/utils"

interface LoadingProps extends React.HTMLAttributes<HTMLDivElement> {
  size?: "default" | "sm" | "lg"
  variant?: "default" | "primary"
}

export function Loading({
  className,
  size = "default",
  variant = "default",
  ...props
}: LoadingProps) {
  return (
    <div
      className={cn(
        "animate-spin rounded-full border-2 border-current border-t-transparent",
        {
          "h-4 w-4": size === "sm",
          "h-6 w-6": size === "default",
          "h-8 w-8": size === "lg",
          "text-muted-foreground": variant === "default",
          "text-primary": variant === "primary",
        },
        className
      )}
      {...props}
    >
      <span className="sr-only">Loading...</span>
    </div>
  )
} 