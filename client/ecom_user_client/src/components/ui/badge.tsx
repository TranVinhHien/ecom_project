import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const badgeVariants = cva(
  "inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2",
  {
    variants: {
      variant: {
        default:
          "border-transparent bg-gradient-to-r from-[hsl(var(--button-primary))] to-[hsl(var(--button-secondary))] text-primary-foreground shadow-sm hover:shadow-md hover:scale-105",
        secondary:
          "border-transparent bg-gradient-to-r from-[hsl(var(--button-secondary))] to-[hsl(var(--accent))] text-primary-foreground shadow-sm hover:shadow-md hover:scale-105",
        destructive:
          "border-transparent bg-gradient-to-r from-destructive to-red-600 text-destructive-foreground shadow-sm hover:shadow-md hover:scale-105",
        outline: "text-foreground border-[hsl(var(--border))] hover:bg-gradient-to-r hover:from-[hsl(var(--accent))]/10 hover:to-[hsl(var(--button-primary))]/5 hover:border-[hsl(var(--accent))]",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

export interface BadgeProps
  extends React.HTMLAttributes<HTMLDivElement>,
    VariantProps<typeof badgeVariants> {}

function Badge({ className, variant, ...props }: BadgeProps) {
  return (
    <div className={cn(badgeVariants({ variant }), className)} {...props} />
  )
}

export { Badge, badgeVariants }
