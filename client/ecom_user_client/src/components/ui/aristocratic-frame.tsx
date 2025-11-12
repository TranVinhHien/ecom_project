import * as React from "react"
import { cn } from "@/lib/utils"

/**
 * AristocraticFrame Component
 * 
 * A decorative wrapper component that adds elegant aristocratic styling
 * Perfect for wrapping important content sections, product cards, or featured elements
 * 
 * Usage:
 * <AristocraticFrame variant="gold">
 *   <YourContent />
 * </AristocraticFrame>
 */

interface AristocraticFrameProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: "default" | "gold" | "double" | "minimal"
  children: React.ReactNode
}

const AristocraticFrame = React.forwardRef<HTMLDivElement, AristocraticFrameProps>(
  ({ className, variant = "default", children, ...props }, ref) => {
    const variants = {
      default: "border-2 border-border shadow-elegant before:border-2 before:border-antique-gold/20",
      gold: "border-2 border-antique-gold/40 shadow-elegant-lg before:border-2 before:border-antique-gold/30",
      double: "border-2 border-border shadow-double before:border before:border-antique-gold/15",
      minimal: "border border-border/50 shadow-elegant before:border before:border-antique-gold/10",
    }

    return (
      <div
        ref={ref}
        className={cn(
          "relative overflow-hidden rounded-none p-1",
          "before:absolute before:inset-0 before:pointer-events-none",
          variants[variant],
          className
        )}
        {...props}
      >
        <div className="relative z-10 h-full w-full">
          {children}
        </div>
      </div>
    )
  }
)

AristocraticFrame.displayName = "AristocraticFrame"

export { AristocraticFrame }
export type { AristocraticFrameProps }
