"use client";
import * as React from "react";

import { cn } from "@/lib/utils";

import { Button } from "@/components/ui/button"
import {
    Command,
    CommandGroup,
    CommandItem,
    CommandList,
} from "@/components/ui/command"
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from "@/components/ui/popover"
import { useTheme } from "next-themes";
import { useThemeContext } from "@/components/theme-color-provider";

const availableThemeColors = [
    { name: "Zinc", light: "bg-zinc-900", dark: "bg-zinc-700" },
    { name: "Rose", light: "bg-rose-600", dark: "bg-rose-700" },
    { name: "Blue", light: "bg-blue-600", dark: "bg-blue-700" },
    { name: "Green", light: "bg-green-600", dark: "bg-green-500" },
    { name: "Orange", light: "bg-orange-500", dark: "bg-orange-700" },
    { name: "Purple", light: "bg-purple-600", dark: "bg-purple-500" },
    { name: "Cyan", light: "bg-cyan-500", dark: "bg-cyan-600" },
    { name: "Yellow", light: "bg-yellow-500", dark: "bg-yellow-400" },
    { name: "Teal", light: "bg-teal-600", dark: "bg-teal-500" },
];
export default function ThemeColorToggle() {
    const [open, setOpen] = React.useState(false)
    const { themeColor, setThemeColor } = useThemeContext();
    const { theme } = useTheme();
    const refColor = React.useRef<HTMLDivElement>(null)

    React.useEffect(() => {
        // when hover on Popover, set open to true
        const handleMouseEnter = () => {
            setOpen(true)
        }
        const handleMouseLeave = () => {
            setOpen(false)
        }
        refColor.current?.addEventListener('mouseenter', handleMouseEnter)
        // refColor.current?.addEventListener('mouseleave', handleMouseLeave)
    }, [])

    const createSelectItems = () => {
        return <div className="flex w-max">
            {
                availableThemeColors.map(({ name, light, dark }) => (
                    <CommandItem
                        onSelect={(currentValue: any) => {
                            setThemeColor(currentValue as ThemeColors)
                        }}
                        className="rounded-full p-3 w-auto h-auto justify-center"
                        key={name} value={name}>
                        <div className="flex item-center space-x-3 rounded-full">
                            <div
                                className={cn(
                                    "rounded-full",
                                    "w-[20px]",
                                    "h-[20px]",
                                    theme === "light" ? light : dark,
                                )}
                            ></div>
                        </div>
                    </CommandItem>
                ))

            }
        </div>
    };  
    return (
        <Popover open={open}>
            <PopoverTrigger asChild>
                <Button
                    variant="outline"
                    role="combobox"
                    aria-expanded={open}
                    className="w-[40px] justify-center border-none rounded-full p-3"
                >
                    {themeColor
                        ? <div className="flex items-center ">
                            <div
                                ref={refColor}
                                className={cn(
                                    "items-center",
                                    "rounded-full",
                                    "w-[20px]",
                                    "h-[20px]",
                                    theme === "light" ? availableThemeColors.find(cl => cl.name === themeColor)?.light :
                                        availableThemeColors.find(cl => cl.name === themeColor)?.dark,
                                )}
                            ></div>
                        </div>
                        : "Select color..."}
                </Button>
            </PopoverTrigger>
            <PopoverContent className="p-0">
                <Command>
                    <CommandList>
                        <CommandGroup>
                            {createSelectItems()}
                        </CommandGroup>
                    </CommandList>
                </Command>
            </PopoverContent>
        </Popover>
    )
}
