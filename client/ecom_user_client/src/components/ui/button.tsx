import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap text-base font-medium tracking-wide ring-offset-background transition-all duration-300 ease-in-out focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0",
  {
    variants: {
      variant: {
        default: "rounded-xl bg-gradient-to-br from-[hsl(var(--button-primary))] to-[hsl(var(--button-secondary))] text-primary-foreground shadow-lg hover:shadow-xl hover:scale-[1.02] hover:-translate-y-0.5 active:scale-[0.98]",
        destructive:
          "rounded-xl bg-gradient-to-br from-destructive to-red-600 text-destructive-foreground shadow-lg hover:shadow-xl hover:scale-[1.02] hover:-translate-y-0.5 active:scale-[0.98]",
        outline:
          "rounded-xl border-2 border-[hsl(var(--border))] bg-background/50 backdrop-blur-sm text-foreground hover:bg-gradient-to-br hover:from-[hsl(var(--accent))]/10 hover:to-[hsl(var(--button-primary))]/5 hover:border-[hsl(var(--accent))] hover:scale-[1.02] active:scale-[0.98]",
        secondary:
          "rounded-xl bg-gradient-to-br from-[hsl(var(--button-secondary))] to-[hsl(var(--accent))] text-primary-foreground shadow-md hover:shadow-lg hover:scale-[1.02] hover:-translate-y-0.5 active:scale-[0.98]",
        ghost: "rounded-xl text-[hsl(var(--link-text))] hover:bg-gradient-to-br hover:from-[hsl(var(--accent))]/15 hover:to-[hsl(var(--button-primary))]/10 hover:text-[hsl(var(--link-hover))] hover:scale-[1.02] active:scale-[0.98]",
        link: "rounded-lg text-[hsl(var(--link-text))] underline-offset-4 hover:underline hover:text-[hsl(var(--link-hover))] transition-colors",
        gold: "rounded-xl bg-gradient-to-br from-[hsl(var(--accent))] to-[hsl(var(--button-primary))] text-primary-foreground shadow-lg hover:shadow-xl hover:scale-[1.02] hover:-translate-y-0.5 active:scale-[0.98]",
        royal: "rounded-xl bg-gradient-to-br from-primary to-[hsl(var(--accent))] text-primary-foreground shadow-lg hover:shadow-xl hover:scale-[1.02] hover:-translate-y-0.5 active:scale-[0.98]",
      },
      size: {
        default: "h-11 px-6 py-2.5",
        sm: "h-9 rounded-lg px-4 text-sm",
        lg: "h-13 rounded-xl px-8 text-lg",
        icon: "h-11 w-11 rounded-xl",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
)

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button"
    return (
      <Comp
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        {...props}
      />
    )
  }
)
Button.displayName = "Button"

export { Button, buttonVariants }
