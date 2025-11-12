import type { Config } from "tailwindcss";

const config: Config = {
	darkMode: ["class"],
	content: [
		"./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
		"./src/components/**/*.{js,ts,jsx,tsx,mdx}",
		"./src/app/**/*.{js,ts,jsx,tsx,mdx}",
		"./src/resources/**/*.{js,ts,jsx,tsx,mdx}",
	],
	theme: {
		extend: {
			/* ============================================
			   ARISTOCRATIC EUROPEAN LUXURY THEME
			   Optimized for Vietnamese readability
			   ============================================ */
			
			fontFamily: {
				// Serif font for aristocratic headings
				serif: ['Playfair Display', 'Georgia', 'Times New Roman', 'serif'],
				// Sans-serif for body text (better Vietnamese readability)
				sans: ['Inter', '-apple-system', 'BlinkMacSystemFont', 'Segoe UI', 'sans-serif'],
			},
			
			colors: {
				background: 'hsl(var(--background))',
				foreground: 'hsl(var(--foreground))',
				card: {
					DEFAULT: 'hsl(var(--card))',
					foreground: 'hsl(var(--card-foreground))'
				},
				popover: {
					DEFAULT: 'hsl(var(--popover))',
					foreground: 'hsl(var(--popover-foreground))'
				},
				primary: {
					DEFAULT: 'hsl(var(--primary))',
					foreground: 'hsl(var(--primary-foreground))'
				},
				secondary: {
					DEFAULT: 'hsl(var(--secondary))',
					foreground: 'hsl(var(--secondary-foreground))'
				},
				muted: {
					DEFAULT: 'hsl(var(--muted))',
					foreground: 'hsl(var(--muted-foreground))'
				},
				accent: {
					DEFAULT: 'hsl(var(--accent))',
					foreground: 'hsl(var(--accent-foreground))'
				},
				destructive: {
					DEFAULT: 'hsl(var(--destructive))',
					foreground: 'hsl(var(--destructive-foreground))'
				},
				border: 'hsl(var(--border))',
				input: 'hsl(var(--input))',
				ring: 'hsl(var(--ring))',
				chart: {
					'1': 'hsl(var(--chart-1))',
					'2': 'hsl(var(--chart-2))',
					'3': 'hsl(var(--chart-3))',
					'4': 'hsl(var(--chart-4))',
					'5': 'hsl(var(--chart-5))'
				},
				
				// Additional Aristocratic Colors
				'antique-gold': 'hsl(var(--antique-gold))',
				'royal-blue': 'hsl(var(--royal-blue))',
				'burgundy': 'hsl(var(--burgundy))',
				'emerald': 'hsl(var(--emerald))',
				'charcoal': 'hsl(var(--charcoal))',
			},
			
			borderRadius: {
				// Sharp corners for aristocratic look
				lg: 'var(--radius)',
				md: 'calc(var(--radius) - 2px)',
				sm: 'calc(var(--radius) - 4px)',
				none: '0',
			},
			
			spacing: {
				// Additional spacing for generous layouts
				'18': '4.5rem',
				'22': '5.5rem',
				'26': '6.5rem',
				'30': '7.5rem',
			},
			
			letterSpacing: {
				// Refined letter spacing
				tighter: '-0.02em',
				tight: '-0.01em',
				normal: '0',
				wide: '0.01em',
				wider: '0.02em',
				widest: '0.05em',
			},
			
			boxShadow: {
				// Aristocratic shadow effects
				'frame': 'inset 0 0 0 1px hsl(var(--antique-gold) / 0.3), inset 0 0 0 3px hsl(var(--background)), inset 0 0 0 4px hsl(var(--antique-gold) / 0.2)',
				'double': 'inset 0 0 0 3px hsl(var(--background)), inset 0 0 0 4px hsl(var(--border))',
				'elegant': '0 4px 6px -1px hsl(var(--foreground) / 0.05), 0 2px 4px -2px hsl(var(--foreground) / 0.05)',
				'elegant-lg': '0 10px 15px -3px hsl(var(--foreground) / 0.08), 0 4px 6px -4px hsl(var(--foreground) / 0.08)',
			},
			
			keyframes: {
				'accordion-down': {
					from: {
						height: '0'
					},
					to: {
						height: 'var(--radix-accordion-content-height)'
					}
				},
				'accordion-up': {
					from: {
						height: 'var(--radix-accordion-content-height)'
					},
					to: {
						height: '0'
					}
				},
				'fade-in': {
					from: {
						opacity: '0',
						transform: 'translateY(10px)'
					},
					to: {
						opacity: '1',
						transform: 'translateY(0)'
					}
				},
			},
			
			animation: {
				'accordion-down': 'accordion-down 0.2s ease-out',
				'accordion-up': 'accordion-up 0.2s ease-out',
				'fade-in': 'fade-in 0.5s ease-out',
			}
		}
	},
	plugins: [require("tailwindcss-animate")],
};
export default config;
