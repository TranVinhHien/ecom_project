type ThemeColors = "Zinc" | "Rose" | "Blue" | "Green" | "Orange" | "Purple" | "Cyan" | "Yellow" | "Teal"
interface ThemeColorStateParams {
    themeColor: ThemeColors;
    setThemeColor: React.Dispatch<React.SetStateAction<ThemeColors>>
}