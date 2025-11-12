import { Button } from "@/components/ui/button"
import { Moon, Sun } from "lucide-react"
import { useTheme } from "next-themes"
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Toggle } from "@/components/ui/toggle"
export function ModeToggle() {
    const { theme, setTheme } = useTheme()
    return (
        <>
            {theme === "light" ? <Toggle className="border-none rounded-full" onClick={() => setTheme("dark")}>
                <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
            </Toggle> :
                <Toggle className="border-none rounded-full" onClick={() => setTheme("light")}>  <Moon className="h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
                </Toggle>

            }
        </>

    )
}
