"use client"

import { Button } from "@/components/ui/button"
import { PlusCircle } from "lucide-react"
import { cn } from "@/lib/utils" // Optional if you're using class merging

interface NewCertificateButtonProps {
  onClick?: () => void
  className?: string
}

export function NewCertificateButton({
  onClick,
  className,
}: NewCertificateButtonProps) {
  return (
    <Button
      onClick={onClick}
      variant="outline"
      className={cn(
        "gap-2 text-sm font-medium p-2.5", // Added proper padding and fixed font-weight
        "border-border/60 hover:border-primary/50", // Subtle border color change
        "hover:bg-primary/5 hover:text-primary", // Gentler background hover
        "transition-all duration-200 ease-in-out", // Smooth transition for all properties
        "hover:shadow-sm hover:scale-[1.02]", // Subtle lift effect
        "active:scale-[0.98] active:transition-none", // Quick press feedback
        "focus-visible:ring-2 focus-visible:ring-primary/20", // Better focus state,
        "hover:bg-gradient-to-r hover:from-primary/5 hover:to-primary/10", // Gradient hover
        "shadow-[0_1px_3px_rgba(0,0,0,0.1)] hover:shadow-[0_4px_12px_rgba(0,0,0,0.15)]" // Dynamic shadow
        ,className
      )}
    >
      <PlusCircle className="h-4 w-4 transition-transform duration-200 group-hover:rotate-90" />
      New Certificate
    </Button>
  )
}

