import { createNavigation } from "next-intl/navigation"
import { locales } from "@/assets/configs/language"
export const { Link, redirect, usePathname, useRouter } = createNavigation({ locales })