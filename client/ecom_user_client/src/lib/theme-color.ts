import { themes } from "@/assets/configs/theme-color";

export default function setGlobalColorTheme(
    themeMode: "light" | "dark",
    color: ThemeColors,
) {
    const theme = themes[color][themeMode] as {
        [key: string]: string;
    };

    // Apply all theme colors to CSS variables
    for (const key in theme) {
        const cssVarName = key.replace(/([A-Z])/g, '-$1').toLowerCase();
        document.documentElement.style.setProperty(`--${cssVarName}`, theme[key]);
    }
}

// Helper function to get current theme color values
export function getThemeColor(colorKey: string): string {
    return getComputedStyle(document.documentElement)
        .getPropertyValue(`--${colorKey}`)
        .trim();
}

// Helper to create harmonious color combinations
export function getColorPalette(baseColor: ThemeColors, themeMode: "light" | "dark") {
    return themes[baseColor][themeMode];
}