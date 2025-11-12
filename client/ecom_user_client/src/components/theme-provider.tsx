"use client"

import * as React from "react"
import { ThemeProvider as NextThemesProvider } from "next-themes"
import ThemeDataProvider from "./theme-color-provider"

export function ThemeProvider({
    children,
    ...props
}: React.ComponentProps<typeof NextThemesProvider>) {
    return <NextThemesProvider {...props}>
        <ThemeDataProvider>
            {children}
        </ThemeDataProvider>
    </NextThemesProvider>
}