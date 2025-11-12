"use client";

import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

/**
 * ColorSystemDemo Component
 * 
 * Component demo để hiển thị tất cả các màu và hiệu ứng hover mới
 * Sử dụng component này để test và xem preview các màu
 */
export default function ColorSystemDemo() {
    return (
        <div className="container mx-auto p-8 space-y-12">
            {/* Header */}
            <div className="text-center space-y-4">
                <h1 className="text-4xl font-bold text-[hsl(var(--foreground))]">
                    🎨 Enhanced Color System
                </h1>
                <p className="text-[hsl(var(--muted-foreground))]">
                    Beautiful colors and smooth hover effects without white overlay
                </p>
            </div>

            {/* Button Variants */}
            <section className="space-y-6">
                <h2 className="text-2xl font-semibold text-[hsl(var(--foreground))]">
                    Button Variants
                </h2>
                <div className="flex flex-wrap gap-4">
                    <Button variant="default">Primary Button</Button>
                    <Button variant="secondary">Secondary Button</Button>
                    <Button variant="outline">Outline Button</Button>
                    <Button variant="ghost">Ghost Button</Button>
                    <Button variant="gold">Gold Button</Button>
                    <Button variant="royal">Royal Button</Button>
                    <Button variant="destructive">Destructive</Button>
                </div>
            </section>

            {/* Hover Effects */}
            <section className="space-y-6">
                <h2 className="text-2xl font-semibold text-[hsl(var(--foreground))]">
                    ✨ New Gradient Hover Effects
                </h2>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <Button className="hover-glow-gradient">
                        Glow Gradient
                    </Button>
                    <Button className="hover-lift-smooth">
                        Lift Smooth
                    </Button>
                    <Button className="hover-shimmer-gradient">
                        Shimmer Gradient
                    </Button>
                    <Button className="hover-gradient-smooth">
                        Gradient Smooth
                    </Button>
                    <Button className="hover-border-glow">
                        Border Glow
                    </Button>
                    <Button variant="secondary">
                        Secondary Gradient
                    </Button>
                </div>
            </section>

            {/* Badges */}
            <section className="space-y-6">
                <h2 className="text-2xl font-semibold text-[hsl(var(--foreground))]">
                    Badges
                </h2>
                <div className="flex flex-wrap gap-4">
                    <Badge variant="default">Default Badge</Badge>
                    <Badge variant="secondary">Secondary Badge</Badge>
                    <Badge variant="destructive">Destructive Badge</Badge>
                    <Badge variant="outline">Outline Badge</Badge>
                </div>
            </section>

            {/* Text Colors */}
            <section className="space-y-6">
                <h2 className="text-2xl font-semibold text-[hsl(var(--foreground))]">
                    Text Colors
                </h2>
                <div className="space-y-3">
                    <p className="text-[hsl(var(--foreground))]">
                        Foreground Text - Main text color
                    </p>
                    <p className="text-[hsl(var(--primary))]">
                        Primary Text - Primary theme color
                    </p>
                    <p className="text-[hsl(var(--accent))]">
                        Accent Text - Accent highlights
                    </p>
                    <p className="text-[hsl(var(--muted-foreground))]">
                        Muted Text - Secondary information
                    </p>
                </div>
            </section>

            {/* Links with Hover */}
            <section className="space-y-6">
                <h2 className="text-2xl font-semibold text-[hsl(var(--foreground))]">
                    🔗 Links - Màu Tương Phản Cao
                </h2>
                <div className="space-y-3">
                    <a 
                        href="#" 
                        className="block link-gradient-underline font-medium"
                    >
                        Link with gradient underline (màu rõ ràng, không bị mất)
                    </a>
                    <a 
                        href="#" 
                        className="block text-[hsl(var(--link-text))] hover:text-[hsl(var(--link-hover))] transition-colors duration-300 font-medium"
                    >
                        Link with color transition (tương phản cao)
                    </a>
                    <a 
                        href="#" 
                        className="block text-[hsl(var(--link-text))] hover:text-[hsl(var(--link-hover))] transition-all duration-300 hover:scale-105 font-medium"
                    >
                        Link with scale effect (hover mượt mà)
                    </a>
                </div>
            </section>

            {/* Cards */}
            <section className="space-y-6">
                <h2 className="text-2xl font-semibold text-[hsl(var(--foreground))]">
                    🎴 Cards - Bo Tròn & Gradient
                </h2>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <div className="card-hover-smooth p-6 border border-[hsl(var(--border))]">
                        <h3 className="font-semibold mb-2 text-[hsl(var(--foreground))]">Card 1</h3>
                        <p className="text-sm text-[hsl(var(--muted-foreground))]">
                            Card with smooth lift & gradient background
                        </p>
                    </div>
                    <div className="hover-lift-smooth p-6 rounded-xl bg-gradient-to-br from-[hsl(var(--card))] to-[hsl(var(--popover))] border border-[hsl(var(--border))]">
                        <h3 className="font-semibold mb-2 text-[hsl(var(--foreground))]">Card 2</h3>
                        <p className="text-sm text-[hsl(var(--muted-foreground))]">
                            Card with lift effect & gradient
                        </p>
                    </div>
                    <div className="hover-glow-gradient p-6 rounded-xl">
                        <h3 className="font-semibold mb-2 text-white">Card 3</h3>
                        <p className="text-sm text-white/90">
                            Card with glow gradient effect
                        </p>
                    </div>
                </div>
            </section>

            {/* Color Palette Display */}
            <section className="space-y-6">
                <h2 className="text-2xl font-semibold text-[hsl(var(--foreground))]">
                    Current Theme Palette
                </h2>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div className="space-y-2">
                        <div className="h-20 rounded-lg bg-[hsl(var(--button-primary))] flex items-center justify-center text-white font-semibold">
                            Primary
                        </div>
                        <p className="text-sm text-center text-[hsl(var(--muted-foreground))]">
                            Button Primary
                        </p>
                    </div>
                    <div className="space-y-2">
                        <div className="h-20 rounded-lg bg-[hsl(var(--button-primary-hover))] flex items-center justify-center text-white font-semibold">
                            Primary Hover
                        </div>
                        <p className="text-sm text-center text-[hsl(var(--muted-foreground))]">
                            On Hover
                        </p>
                    </div>
                    <div className="space-y-2">
                        <div className="h-20 rounded-lg bg-[hsl(var(--button-secondary))] flex items-center justify-center text-white font-semibold">
                            Secondary
                        </div>
                        <p className="text-sm text-center text-[hsl(var(--muted-foreground))]">
                            Button Secondary
                        </p>
                    </div>
                    <div className="space-y-2">
                        <div className="h-20 rounded-lg bg-[hsl(var(--accent))] flex items-center justify-center text-white font-semibold">
                            Accent
                        </div>
                        <p className="text-sm text-center text-[hsl(var(--muted-foreground))]">
                            Highlights
                        </p>
                    </div>
                </div>
            </section>

            {/* Instructions */}
            <section className="bg-[hsl(var(--card))] p-6 rounded-lg border border-[hsl(var(--border))]">
                <h2 className="text-2xl font-semibold text-[hsl(var(--foreground))] mb-4">
                    💡 How to Use
                </h2>
                <ul className="space-y-2 text-[hsl(var(--muted-foreground))]">
                    <li>✅ Chuyển đổi theme color để xem các bảng màu khác nhau</li>
                    <li>✅ Tất cả hover effects không có white overlay</li>
                    <li>✅ Màu sắc hài hòa trong cả Light và Dark mode</li>
                    <li>✅ Mỗi theme có màu riêng cho buttons, text, headers</li>
                    <li>✅ Dễ dàng thêm màu mới vào hệ thống</li>
                </ul>
            </section>
        </div>
    );
}
