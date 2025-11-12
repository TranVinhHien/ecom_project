import * as React from "react"

import { cn } from "@/lib/utils"

const Separator = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      "shrink-0 h-[2px] w-full relative",
      "bg-gradient-to-r from-transparent via-border to-transparent",
      "before:content-['◆'] before:absolute before:left-0 before:top-1/2 before:-translate-y-1/2 before:text-antique-gold/50 before:text-xs",
      "after:content-['◆'] after:absolute after:right-0 after:top-1/2 after:-translate-y-1/2 after:text-antique-gold/50 after:text-xs",
      className
    )}
    {...props}
  />
))
Separator.displayName = "Separator"

export { Separator }
