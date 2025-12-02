"use client"

import { useRouter } from "@/i18n/routing";
import { Banner } from "@/types/shop.types";
import { getImageUrl } from "@/assets/helpers/convert_tool";
import { useEffect, useMemo, useState } from "react";

interface HomeBannerSliderProps {
  banners: Banner[];
  className?: string;
  height?: string;
}

export default function HomeBannerSlider({ 
  banners, 
  className = "w-full px-4 md:px-8 lg:px-16 mb-10",
  height = "h-48 md:h-64 lg:h-80"
}: HomeBannerSliderProps) {
  const router = useRouter();
  const sortedBanners = useMemo(() => {
    return [...banners].sort((a, b) => (a.bannerOrder || 0) - (b.bannerOrder || 0));
  }, [banners]);
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    if (sortedBanners.length <= 1) return;
    const interval = setInterval(() => {
      setCurrentIndex((prev) => (prev + 1) % sortedBanners.length);
    }, 5000);
    return () => clearInterval(interval);
  }, [sortedBanners.length]);

  if (sortedBanners.length === 0) {
    return null;
  }

  const handleNavigate = (rawUrl?: string) => {
    if (!rawUrl) return;
    const url = rawUrl.trim();
    if (url.startsWith("http")) {
      window.open(url, "_blank", "noopener noreferrer");
      return;
    }
    router.push(url);
  };

  return (
    <div className={className}>
      <div className={`relative w-full ${height} rounded-2xl overflow-hidden shadow-lg`}>
        {sortedBanners.map((banner, index) => (
          <button
            key={banner.id}
            type="button"
            onClick={() => handleNavigate(banner.bannerUrl)}
            className={`absolute inset-0 transition-opacity duration-700 ${currentIndex === index ? "opacity-100 z-10" : "opacity-0"}`}
            aria-label={banner.bannerName}
          >
            <img
              src={getImageUrl(banner.bannerImage)}
              alt={banner.bannerName}
              className="object-cover w-[70vw] block mx-auto h-[250px] md:h-[400px] lg:h-[500px] xl:h-[650px]"
              sizes="70vw"
            />
          </button>
        ))}

        {sortedBanners.length > 1 && (
          <>
            <button
              type="button"
              className="absolute left-4 top-1/2 -translate-y-1/2 bg-black/40 text-white rounded-full p-2 z-20"
              onClick={() => setCurrentIndex((prev) => (prev - 1 + sortedBanners.length) % sortedBanners.length)}
              aria-label="Previous banner"
            >
              ‹
            </button>
            <button
              type="button"
              className="absolute right-4 top-1/2 -translate-y-1/2 bg-black/40 text-white rounded-full p-2 z-20"
              onClick={() => setCurrentIndex((prev) => (prev + 1) % sortedBanners.length)}
              aria-label="Next banner"
            >
              ›
            </button>
            <div className="absolute bottom-4 inset-x-0 flex justify-center gap-2 z-20">
              {sortedBanners.map((_, index) => (
                <span
                  key={`dot-${index}`}
                  className={`h-2 w-2 rounded-full ${currentIndex === index ? "bg-white" : "bg-white/50"}`}
                />
              ))}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
