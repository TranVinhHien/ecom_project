"use client"
import React, { useEffect, useRef, useState } from 'react'
import { Button } from '@/components/ui/button'
// dark mode 
import { useTheme } from "next-themes";
import { AlignJustify, Moon, Sun, ChevronRight, ChevronDown, Settings, Search, ShoppingCart, Bell, User, LogOut, Camera, Menu, X } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { Sheet, SheetContent, SheetTrigger, SheetHeader, SheetTitle } from "@/components/ui/sheet";

import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuGroup,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuPortal,
    DropdownMenuSeparator,
    DropdownMenuShortcut,
    DropdownMenuSub,
    DropdownMenuSubContent,
    DropdownMenuSubTrigger,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select"
import { useLocale, useTranslations } from 'next-intl'
import { useRouter,usePathname } from "@/i18n/routing"
import { useSearchParams } from 'next/navigation'
import Image from 'next/image'
import { ModeToggle } from './theme-mode.toggle';
import LanguageToggle from './language-toggle';
import ThemeColorToggle from './theme-color-toggle';
import { INFO_USER } from '@/assets/configs/request';
import ROUTER from '@/assets/configs/routers';
import { Link } from '@/i18n/routing';
import { cookies, request, requestNoToken } from '@/assets/helpers';
import { dataTagErrorSymbol, useMutation } from '@tanstack/react-query';
import { Input } from "@/components/ui/input"
import API from '@/assets/configs/api';
import { AxiosError } from 'axios';
import { MetaType, ParamType } from '@/assets/types/request';
import { useGetCategories } from '@/services/apiService';
import { useCartStore } from '@/store/cartStore';
import { Category } from '@/types/category.types';
import logo from "../../../public/logo_doan.png"
import { useGetCartCount } from '@/services/apiService';

export default function Header({ onClick }: { onClick: () => void }) {
    const t = useTranslations("System")
    const locale = useLocale();
    const pathname = usePathname();
    const searchParams = useSearchParams();
    const router = useRouter()
    const [searchQuery, setSearchQuery] = useState("")
    const [openSettings, setOpenSettings] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);
    const [isHydrated, setIsHydrated] = useState(false);
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
    const [info, setInfo] = useState<UserLoginType | null>(null);

    // Sử dụng React Query để lấy categories
    const { data: categories = [], isLoading: categoriesLoading } = useGetCategories();
    
    // Lấy số lượng sản phẩm trong giỏ hàng
    // Nếu đã đăng nhập: lấy từ API
    // Nếu chưa đăng nhập: lấy từ localStorage
    const localCartItemsCount = useCartStore((state) => state.getTotalItems());
    const { data: apiCartCount, isLoading: apiCartLoading } = useGetCartCount();
    
    // Determine cart count based on login status
    const cartItemsCount = isHydrated && info ? (apiCartCount || 0) : localCartItemsCount;

    // Hydrate cart store
    useEffect(() => {
        useCartStore.persist.rehydrate();
        setIsHydrated(true);
    }, []);

    useEffect(() => {
        const jinfo = localStorage.getItem(INFO_USER)
        if (jinfo === null) return
        const info: UserLoginType = JSON.parse(jinfo)
        setInfo(info)
    }, [])

    const handleImageUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
        // const file = event.target.files?.[0];
        // if (file) {
        //     const reader = new FileReader();
        //     reader.onload = (e) => {
        //         // const base64Image = e.target?.result as string;
        //         // Store the image in localStorage
        //         // localStorage.setItem('searchImage', base64Image);
        //         // if ( pathname=== ROUTER.timkiem.image) {
        //         //     window.location.reload();
        //         // } else {
        //         //     router.push(ROUTER.timkiem.image);
        //         // }
    
        // };
        //     reader.readAsDataURL(file);
        // }
    };

    // Function to check if a category is selected
    const isCategorySelected = (categoryPath: string) => {
        if (pathname === '/tim-kiem') {
            const currentCategoryPath = searchParams.get('cate_path');
            return currentCategoryPath === categoryPath;
        }
        return false;
    };

    // Function to check if a category is a parent of the selected category
    const isParentOfSelected = (category: Category): boolean => {
        if (!category.child.valid || !category.child.data) return false;
        
        const currentCategoryPath = searchParams.get('cate_path');
        if (!currentCategoryPath) return false;

        // Check if current category path is parent of selected path
        return currentCategoryPath.startsWith(category.path);
    };

    const handleChangeLanguage = (lang: string) => {
        router.push(pathname, { locale: lang });
    };

    const renderCategoryTree = (category: Category, isTopLevel: boolean = false) => {
        const isSelected = isCategorySelected(category.path);
        const isParent = isParentOfSelected(category);
        const shouldHighlight = isSelected || isParent;
        
        // Nếu không có child, render như link thông thường
        if (!category.child.valid || !category.child.data || category.child.data.length === 0) {
            return (
                <DropdownMenuItem 
                    key={category.category_id}
                    className={`p-3 transition-all duration-200 rounded-lg ${
                        shouldHighlight 
                            ? 'bg-[hsl(var(--primary)/0.15)] hover:bg-[hsl(var(--primary)/0.2)]' 
                            : 'hover:bg-[hsl(var(--primary)/0.08)]'
                    }`}
                    asChild
                >
                    <Link 
                        href={`/tim-kiem?cate_path=${encodeURIComponent(category.path)}`}
                        className={`font-medium w-full ${
                            shouldHighlight 
                                ? 'text-[hsl(var(--primary))]' 
                                : 'text-[hsl(var(--primary)/0.7)] hover:text-[hsl(var(--primary))]'
                        }`}
                    >
                        {category.name}
                    </Link>
                </DropdownMenuItem>
            );
        }
        
        // Tính số column dựa vào số lượng children
        const childCount = category.child.data.length;
        let columns = 1;
        if (childCount > 20) columns = 4;
        else if (childCount > 12) columns = 3;
        else if (childCount > 6) columns = 2;

        // Nếu có child, render với submenu
        return (
            <DropdownMenuSub key={category.category_id}>
                <div className="flex items-center w-full group">
                    <DropdownMenuSubTrigger 
                        className={`flex-1 p-3 transition-all duration-200 rounded-lg w-full ${
                            shouldHighlight 
                                ? 'bg-[hsl(var(--primary)/0.15)] hover:bg-[hsl(var(--primary)/0.2)]' 
                                : 'hover:bg-[hsl(var(--primary)/0.08)]'
                        }`}
                        showArrow={false}
                    >
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-2">
                                <Link 
                                    href={`/tim-kiem?cate_path=${encodeURIComponent(category.path)}`}
                                    className={`font-medium ${
                                        shouldHighlight 
                                            ? 'text-[hsl(var(--primary))]' 
                                            : 'text-[hsl(var(--primary)/0.7)] group-hover:text-[hsl(var(--primary))]'
                                    }`}
                                >
                                    {category.name}
                                </Link>
                                <span className={`text-xs ${
                                    shouldHighlight 
                                        ? 'text-[hsl(var(--primary))]' 
                                        : 'text-[hsl(var(--primary)/0.7)] group-hover:text-[hsl(var(--primary))]'
                                }`}>
                                    ({category.child.data.length})
                                </span>
                            </div>
                            <ChevronRight className={`h-4 w-4 transition-transform duration-200 group-hover:translate-x-1 ${
                                shouldHighlight 
                                    ? 'text-[hsl(var(--primary))]' 
                                    : 'text-[hsl(var(--primary)/0.7)] group-hover:text-[hsl(var(--primary))]'
                            }`} />
                        </div>
                    </DropdownMenuSubTrigger>
                </div>
                <DropdownMenuSubContent 
                    className={`bg-white shadow-lg rounded-lg p-4 border border-[hsl(var(--primary)/0.15)] ${
                        columns === 1 ? 'min-w-[250px]' :
                        columns === 2 ? 'min-w-[500px] max-w-[600px]' :
                        columns === 3 ? 'min-w-[750px] max-w-[900px]' :
                        'min-w-[1000px] max-w-[70vw]'
                    }`}
                    asChild
                    sideOffset={5}
                >
                    <motion.div
                        initial={{ opacity: 0, x: -10 }}
                        animate={{ opacity: 1, x: 0 }}
                        exit={{ opacity: 0, x: -10 }}
                        transition={{ 
                            duration: 0.2,
                            ease: "easeOut"
                        }}
                    >
                        <div 
                            className="gap-4"
                            style={{
                                display: 'grid',
                                gridTemplateColumns: `repeat(${columns}, 1fr)`,
                            }}
                        >
                            {category.child.data.map((childCategory) =>
                                <div key={childCategory.category_id} className="space-y-1">
                                    {renderCategoryTree(childCategory, false)}
                                </div>
                            )}
                        </div>
                    </motion.div>
                </DropdownMenuSubContent>
            </DropdownMenuSub>
        );
    };

    return (
        <div className="w-full shadow-sm">
            {/* Top Bar - Hidden on mobile */}
            <div className="hidden md:block bg-[hsl(var(--primary))] text-white py-2 px-4 lg:px-28">
                <div className="container mx-auto flex justify-between items-center">
                    <div className="flex items-center space-x-4">
                        <Link href={ROUTER.home} className="text-sm hover:text-gray-200">
                            {t('trang-chu')}
                        </Link>
                        <Link href="#" className="text-sm hover:text-gray-200">
                            {t('gioi-thieu')}
                        </Link>
                        <Link href={ROUTER.khieunai} className="text-sm hover:text-gray-200">
                            {t('khieu-nai')}
                        </Link>
                    </div>
                    {!info ? (
                        <Button variant="outline" className="flex items-center gap-2 border-[hsl(var(--primary))] text-[hsl(var(--primary))] justify-center">
                            <User className="h-4 w-4" />
                            <Link href={ROUTER.auth.login}>
                                {t('dang-nhap')}
                            </Link>
                            
                        </Button>
                    ) : (
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button variant="ghost" className='flex items-center gap-2 hover:bg-white/10'>
                                    <img 
                                        src='/default-avatar.png'
                                        alt="avatar" 
                                        className="w-8 h-8 rounded-full object-cover border-2 border-white" 
                                        onError={()=>{
                                            return "/default-avatar.png";
                                        }}
                                    />
                                    <span className="font-medium hidden lg:inline">{info.name}</span>
                                    <ChevronDown className="h-4 w-4" />
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent className="w-56" align="end">
                                <DropdownMenuLabel>{t('tai-khoan-cua-toi')}</DropdownMenuLabel>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem onClick={() => router.push(ROUTER.profile)}>
                                    <User className="h-4 w-4 mr-2" />
                                    {t('thong-tin-tai-khoan')}
                                </DropdownMenuItem>
                                <DropdownMenuItem onClick={() => router.push(ROUTER.donhang)}>
                                    <ShoppingCart className="h-4 w-4 mr-2" />
                                    {t('don-hang')}
                                </DropdownMenuItem>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem onClick={() => {
                                    cookies.logOut()
                                    router.push(ROUTER.home)
                                }} className="text-red-600">
                                    <LogOut className="h-4 w-4 mr-2" />
                                    {t('dang-xuat')}
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    )}
                </div>
            </div>

            {/* Main Header */}
            <div className="container mx-auto py-3 px-4 md:px-8 lg:px-14">
                <div className="flex items-center justify-between gap-2 md:gap-4 lg:gap-8 w-full">
                    {/* Mobile Menu Button */}
                    <Sheet open={isMobileMenuOpen} onOpenChange={setIsMobileMenuOpen}>
                        <SheetTrigger asChild className="md:hidden">
                            <Button variant="ghost" size="icon">
                                <Menu className="h-6 w-6 text-[hsl(var(--primary))]" />
                            </Button>
                        </SheetTrigger>
                        <SheetContent side="left" className="w-[280px] sm:w-[350px]">
                            <SheetHeader>
                                <SheetTitle className="text-[hsl(var(--primary))]">Menu</SheetTitle>
                            </SheetHeader>
                            <div className="mt-6 space-y-4">
                                {/* User Info Mobile */}
                                {!info ? (
                                    <Link href={ROUTER.auth.login} onClick={() => setIsMobileMenuOpen(false)}>
                                        <Button className="w-full bg-[hsl(var(--primary))]">
                                            <User className="h-4 w-4 mr-2" />
                                            {t('dang-nhap')}
                                        </Button>
                                    </Link>
                                ) : (
                                    <div className="space-y-2">
                                        <div className="flex items-center gap-3 p-3 bg-[hsl(var(--primary)/0.1)] rounded-lg">
                                            <img 
                                                src='/default-avatar.png'
                                                alt="avatar" 
                                                className="w-12 h-12 rounded-full object-cover" 
                                            />
                                            <div>
                                                <p className="font-semibold text-sm">{info.name}</p>
                                            </div>
                                        </div>
                                        <Button 
                                            variant="outline" 
                                            className="w-full justify-start"
                                            onClick={() => {
                                                router.push(ROUTER.profile);
                                                setIsMobileMenuOpen(false);
                                            }}
                                        >
                                            <User className="h-4 w-4 mr-2" />
                                            {t('thong-tin-tai-khoan')}
                                        </Button>
                                        <Button 
                                            variant="outline" 
                                            className="w-full justify-start"
                                            onClick={() => {
                                                router.push(ROUTER.donhang);
                                                setIsMobileMenuOpen(false);
                                            }}
                                        >
                                            <ShoppingCart className="h-4 w-4 mr-2" />
                                            {t('don-hang')}
                                        </Button>
                                        <Button 
                                            variant="outline" 
                                            className="w-full justify-start text-red-600"
                                            onClick={() => {
                                                cookies.logOut();
                                                router.push(ROUTER.home);
                                                setIsMobileMenuOpen(false);
                                            }}
                                        >
                                            <LogOut className="h-4 w-4 mr-2" />
                                            {t('dang-xuat')}
                                        </Button>
                                    </div>
                                )}

                                {/* Navigation Links Mobile */}
                                <div className="border-t pt-4 space-y-2">
                                    <Link 
                                        href={ROUTER.home} 
                                        className="block py-2 hover:bg-[hsl(var(--primary)/0.1)] px-3 rounded-lg"
                                        onClick={() => setIsMobileMenuOpen(false)}
                                    >
                                        {t('trang-chu')}
                                    </Link>
                                    <Link 
                                        href="#" 
                                        className="block py-2 hover:bg-[hsl(var(--primary)/0.1)] px-3 rounded-lg"
                                        onClick={() => setIsMobileMenuOpen(false)}
                                    >
                                        {t('gioi-thieu')}
                                    </Link>
                                    <Link 
                                        href="#" 
                                        className="block py-2 hover:bg-[hsl(var(--primary)/0.1)] px-3 rounded-lg"
                                        onClick={() => setIsMobileMenuOpen(false)}
                                    >
                                        {t('khuyen-mai')}
                                    </Link>
                                    <Link 
                                        href="#" 
                                        className="block py-2 hover:bg-[hsl(var(--primary)/0.1)] px-3 rounded-lg"
                                        onClick={() => setIsMobileMenuOpen(false)}
                                    >
                                        {t('san-pham-moi')}
                                    </Link>
                                    <Link 
                                        href="#" 
                                        className="block py-2 hover:bg-[hsl(var(--primary)/0.1)] px-3 rounded-lg"
                                        onClick={() => setIsMobileMenuOpen(false)}
                                    >
                                        {t('thuong-hieu')}
                                    </Link>
                                </div>

                                {/* Settings Mobile */}
                                <div className="border-t pt-4 space-y-3">
                                    <div className="flex items-center justify-between px-3">
                                        <span className="text-sm font-medium">{t('mau-nen')}</span>
                                        <ThemeColorToggle />
                                    </div>
                                    <div className="flex items-center justify-between px-3">
                                        <span className="text-sm font-medium">{t('che-do-toi')}</span>
                                        <ModeToggle />
                                    </div>
                                    <div className="flex items-center justify-between px-3">
                                        <span className="text-sm font-medium">{t('ngon-ngu')}</span>
                                        <LanguageToggle handleChangeLanguage={handleChangeLanguage} locale={locale} />
                                    </div>
                                </div>
                            </div>
                        </SheetContent>
                    </Sheet>

                    {/* Left: Logo + Category (Desktop) */}
                    <div className="hidden md:flex items-center gap-4 min-w-[180px] lg:min-w-[220px] justify-center">
                        <Link href={ROUTER.home} className="flex-shrink-0 flex items-center justify-center">
                            <div className="text-xl lg:text-2xl font-bold text-[hsl(var(--primary))]">
                                {/* E-Shop */}
                                <Image src={logo} alt='' width={100} height={100} />

                            </div>
                        </Link>
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <Button className="bg-[hsl(var(--primary))] text-white hover:bg-[hsl(var(--primary)/.9)] px-3 lg:px-4 py-2 rounded-lg flex items-center justify-center">
                                    <AlignJustify className="h-4 w-4 lg:h-5 lg:w-5" />
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent className="w-64 max-h-[70vh] overflow-y-auto">
                                {categories.map((category) => renderCategoryTree(category, true))}
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </div>

                    {/* Mobile Logo */}
                    <Link href={ROUTER.home} className="md:hidden flex-shrink-0">
                        <div className="text-xl font-bold text-[hsl(var(--primary))]">
                            {/* E-Shop */}
                            <Image src={logo} alt='' width={100} height={100} />

                        </div>
                    </Link>

                    {/* Center: Search */}
                    <div className="flex-1 max-w-2xl mx-auto">
                        <div className="relative w-full flex items-center gap-1">
                            <Input
                                type="text"
                                placeholder={t('tim-kiem')}
                                className="w-full pl-4 pr-20 md:pr-24 py-2 border-2 border-[hsl(var(--primary))] rounded-full focus:ring-2 focus:ring-[hsl(var(--primary))] text-sm"
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                onKeyPress={(e) => {
                                    if (e.key === 'Enter' && searchQuery.trim() !== "") {
                                        router.push(ROUTER.timkiem.query + "?query=" + searchQuery);
                                    }
                                }}
                            />
                            <input
                                type="file"
                                ref={fileInputRef}
                                className="hidden"
                                accept="image/*"
                                onChange={handleImageUpload}
                            />
                            <Button
                                className="absolute right-12 md:right-14 top-1/2 transform -translate-y-1/2 bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/.9)] text-white rounded-full p-1.5 h-8 w-8 flex items-center justify-center"
                                size="icon"
                                onClick={() => fileInputRef.current?.click()}
                                title={t('tim-kiem-bang-hinh-anh')}
                            >
                                <Camera className="h-3 w-3 md:h-4 md:w-4" />
                            </Button>
                            <Button
                                className="absolute right-1 top-1/2 transform -translate-y-1/2 bg-[hsl(var(--primary))] hover:bg-[hsl(var(--primary)/.9)] text-white rounded-full h-8 w-8 md:h-9 md:w-9"
                                size="icon"
                                onClick={() => {
                                    if (searchQuery.trim() !== "") {
                                        router.push(ROUTER.timkiem.query + "?query=" + searchQuery);
                                    }
                                }}
                            >
                                <Search className="h-3 w-3 md:h-4 md:w-4" />
                            </Button>
                        </div>
                    </div>

                    {/* Right: Actions */}
                    
                    <div className="flex items-center gap-1 md:gap-2">
                                  
                        {/* Desktop Settings */}
                        
                        <DropdownMenu open={openSettings} onOpenChange={setOpenSettings}>
                            <DropdownMenuTrigger asChild className="hidden lg:flex">
                                <Button variant="ghost" size="icon">
                                    <Settings className="h-5 w-5 text-[hsl(var(--primary))]" />
                                </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent className="w-56">
                                <DropdownMenuLabel>{t('cai-dat')}</DropdownMenuLabel>
                                <DropdownMenuSeparator />
                                <DropdownMenuGroup>
                                    <DropdownMenuItem asChild>
                                        <div className="flex items-center justify-between w-full px-2 py-1.5">
                                            <span>{t('mau-nen')}</span>
                                            <ThemeColorToggle />
                                        </div>
                                    </DropdownMenuItem>
                                    <DropdownMenuItem asChild>
                                        <div className="flex items-center justify-between w-full px-2 py-1.5">
                                            <span>{t('che-do-toi')}</span>
                                            <ModeToggle />
                                        </div>
                                    </DropdownMenuItem>
                                    <DropdownMenuItem asChild>
                                        <div className="flex items-center justify-between w-full px-2 py-1.5">
                                            <span>{t('ngon-ngu')}</span>
                                            <LanguageToggle handleChangeLanguage={handleChangeLanguage} locale={locale} />
                                        </div>
                                    </DropdownMenuItem>
                                </DropdownMenuGroup>
                            </DropdownMenuContent>
                        </DropdownMenu>
                        
                        {/* Notifications & Cart - Only show when logged in */}
                        {info && (
                            <>
                                <Button variant="ghost" size="icon" className="relative hidden md:flex items-center justify-center">
                                    <Bell className="h-4 w-4 lg:h-5 lg:w-5 text-[hsl(var(--primary))]" />
                                    <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center">
                                        3
                                    </span>
                                </Button>
                            
                                <Button 
                                    variant="ghost" 
                                    size="icon" 
                                    className="relative flex items-center justify-center"
                                    onClick={() => router.push(ROUTER.giohang)}
                                >
                                    <ShoppingCart className="h-4 w-4 lg:h-5 lg:w-5 text-[hsl(var(--primary))]" />
                                    {isHydrated && cartItemsCount > 0 && (
                                        <span className="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-4 w-4 flex items-center justify-center">
                                            {cartItemsCount}
                                        </span>
                                    )}
                                </Button>
                            </>
                        )}
                    </div>
                </div>

                {/* Bottom Links - Hidden on mobile, visible on tablet+ */}
                <div className="hidden md:flex justify-center space-x-4 lg:space-x-8 mt-2">
                    <Link href="#" className="text-xs lg:text-sm font-medium text-[hsl(var(--primary))] hover:underline">
                        {t('khuyen-mai')}
                    </Link>
                    <Link href="#" className="text-xs lg:text-sm font-medium text-[hsl(var(--primary))] hover:underline">
                        {t('san-pham-moi')}
                    </Link>
                    <Link href="#" className="text-xs lg:text-sm font-medium text-[hsl(var(--primary))] hover:underline">
                        {t('thuong-hieu')}
                    </Link>
                </div>
            </div>

            {/* Đường kẻ phân cách */}
            <div className="w-full border-b border-[hsl(var(--primary)/0.15)] rounded-b-lg shadow-sm"></div>
        </div>
    )
}




