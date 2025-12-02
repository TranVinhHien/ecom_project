"use client";

import React, { useState, useEffect, useRef } from 'react';
import { X, ChevronLeft, ChevronRight, ZoomIn, ZoomOut, Maximize2 } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { getImageUrl } from '@/assets/helpers/convert_tool';
import { cn } from '@/lib/utils';

interface ImageGalleryModalProps {
  images: string[];
  initialIndex: number;
  isOpen: boolean;
  onClose: () => void;
  notes?: (string | null)[];
}

export default function ImageGalleryModal({
  images,
  initialIndex,
  isOpen,
  onClose,
  notes = []
}: ImageGalleryModalProps) {
  const [currentIndex, setCurrentIndex] = useState(initialIndex);
  const [scale, setScale] = useState(1);
  const [position, setPosition] = useState({ x: 0, y: 0 });
  const [isDragging, setIsDragging] = useState(false);
  const [dragStart, setDragStart] = useState({ x: 0, y: 0 });
  const imageContainerRef = useRef<HTMLDivElement>(null);
  const imageRef = useRef<HTMLImageElement>(null);

  const MIN_SCALE = 1;
  const MAX_SCALE = 5;
  const SCALE_STEP = 0.5;

  // Reset index when modal opens
  useEffect(() => {
    setCurrentIndex(initialIndex);
    resetZoom();
  }, [initialIndex, isOpen]);

  // Reset zoom when changing images
  useEffect(() => {
    resetZoom();
  }, [currentIndex]);

  // Reset zoom function
  const resetZoom = () => {
    setScale(1);
    setPosition({ x: 0, y: 0 });
  };

  // Keyboard navigation
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape') {
        onClose();
      } else if (e.key === 'ArrowLeft' && scale === 1) {
        handlePrevious();
      } else if (e.key === 'ArrowRight' && scale === 1) {
        handleNext();
      } else if (e.key === '+' || e.key === '=') {
        handleZoomIn();
      } else if (e.key === '-' || e.key === '_') {
        handleZoomOut();
      } else if (e.key === '0') {
        resetZoom();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, currentIndex, images.length, scale]);

  // Prevent body scroll when modal is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'unset';
    }
    return () => {
      document.body.style.overflow = 'unset';
    };
  }, [isOpen]);

  const handleNext = () => {
    setCurrentIndex((prev) => (prev + 1) % images.length);
  };

  const handlePrevious = () => {
    setCurrentIndex((prev) => (prev - 1 + images.length) % images.length);
  };

  // Zoom functions
  const handleZoomIn = () => {
    setScale((prev) => Math.min(prev + SCALE_STEP, MAX_SCALE));
  };

  const handleZoomOut = () => {
    if (scale <= MIN_SCALE) {
      resetZoom();
    } else {
      setScale((prev) => Math.max(prev - SCALE_STEP, MIN_SCALE));
    }
  };

  // Mouse wheel zoom
  const handleWheel = (e: React.WheelEvent) => {
    e.preventDefault();
    if (e.deltaY < 0) {
      handleZoomIn();
    } else {
      handleZoomOut();
    }
  };

  // Touch/Pinch zoom
  useEffect(() => {
    const container = imageContainerRef.current;
    if (!container || !isOpen) return;

    let initialDistance = 0;
    let initialScale = 1;

    const handleTouchStart = (e: TouchEvent) => {
      if (e.touches.length === 2) {
        e.preventDefault();
        initialDistance = Math.hypot(
          e.touches[0].pageX - e.touches[1].pageX,
          e.touches[0].pageY - e.touches[1].pageY
        );
        initialScale = scale;
      }
    };

    const handleTouchMove = (e: TouchEvent) => {
      if (e.touches.length === 2) {
        e.preventDefault();
        const currentDistance = Math.hypot(
          e.touches[0].pageX - e.touches[1].pageX,
          e.touches[0].pageY - e.touches[1].pageY
        );
        const newScale = (currentDistance / initialDistance) * initialScale;
        setScale(Math.min(Math.max(newScale, MIN_SCALE), MAX_SCALE));
      }
    };

    container.addEventListener('touchstart', handleTouchStart, { passive: false });
    container.addEventListener('touchmove', handleTouchMove, { passive: false });

    return () => {
      container.removeEventListener('touchstart', handleTouchStart);
      container.removeEventListener('touchmove', handleTouchMove);
    };
  }, [isOpen, scale]);

  // Drag to pan when zoomed
  const handleMouseDown = (e: React.MouseEvent) => {
    if (scale > 1) {
      setIsDragging(true);
      setDragStart({ x: e.clientX - position.x, y: e.clientY - position.y });
    }
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (isDragging && scale > 1) {
      setPosition({
        x: e.clientX - dragStart.x,
        y: e.clientY - dragStart.y,
      });
    }
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  const handleMouseLeave = () => {
    setIsDragging(false);
  };

  const handleBackdropClick = (e: React.MouseEvent<HTMLDivElement>) => {
    // Only close if clicking on backdrop and not zoomed in
    if (e.target === e.currentTarget && scale === 1) {
      onClose();
    }
  };

  if (!isOpen) return null;

  const currentImage = images[currentIndex];
  const currentNote = notes[currentIndex] || null;
  const isVideo = currentImage?.endsWith('.mp4');

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/90 backdrop-blur-sm"
      onClick={handleBackdropClick}
    >
      {/* Close Button */}
      <Button
        variant="ghost"
        size="icon"
        className="absolute top-4 right-4 z-50 text-white hover:bg-white/20 h-12 w-12"
        onClick={onClose}
      >
        <X className="h-6 w-6" />
      </Button>

      {/* Zoom Controls */}
      <div className="absolute top-4 right-20 z-50 flex gap-2">
        <Button
          variant="ghost"
          size="icon"
          className="text-white hover:bg-white/20 h-10 w-10"
          onClick={handleZoomIn}
          disabled={scale >= MAX_SCALE}
          title="Phóng to (+)"
        >
          <ZoomIn className="h-5 w-5" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="text-white hover:bg-white/20 h-10 w-10"
          onClick={handleZoomOut}
          disabled={scale <= MIN_SCALE}
          title="Thu nhỏ (-)"
        >
          <ZoomOut className="h-5 w-5" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="text-white hover:bg-white/20 h-10 w-10"
          onClick={resetZoom}
          disabled={scale === MIN_SCALE}
          title="Reset zoom (0)"
        >
          <Maximize2 className="h-5 w-5" />
        </Button>
      </div>

      {/* Image Counter & Zoom Level */}
      <div className="absolute top-4 left-4 z-50 flex flex-col gap-2">
        <div className="bg-black/60 text-white px-4 py-2 rounded-lg font-medium">
          {currentIndex + 1} / {images.length}
        </div>
        {scale > 1 && (
          <div className="bg-black/60 text-white px-4 py-2 rounded-lg font-medium text-sm">
            {Math.round(scale * 100)}%
          </div>
        )}
      </div>

      {/* Previous Button */}
      {images.length > 1 && scale === 1 && (
        <Button
          variant="ghost"
          size="icon"
          className="absolute left-4 top-1/2 -translate-y-1/2 z-50 text-white hover:bg-white/20 h-16 w-16"
          onClick={handlePrevious}
        >
          <ChevronLeft className="h-8 w-8" />
        </Button>
      )}

      {/* Main Content Container */}
      <div 
        ref={imageContainerRef}
        className="relative w-full h-full flex flex-col items-center justify-center p-4"
        onMouseDown={handleMouseDown}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseLeave}
      >
        {/* Image/Video Container */}
        <div 
          className="relative max-w-[90vw] max-h-[80vh] flex items-center justify-center overflow-hidden"
          onWheel={handleWheel}
        >
          {isVideo ? (
            <video
              src={getImageUrl(currentImage)}
              className="max-w-full max-h-[80vh] object-contain rounded-lg"
              style={{
                transform: `scale(${scale}) translate(${position.x / scale}px, ${position.y / scale}px)`,
                transition: isDragging ? 'none' : 'transform 0.2s ease-out',
                cursor: scale > 1 ? (isDragging ? 'grabbing' : 'grab') : 'default',
              }}
              controls
              autoPlay
              playsInline
              onError={(e) => {
                console.error('Video load error');
              }}
            />
          ) : (
            <img
              ref={imageRef}
              src={getImageUrl(currentImage)}
              alt={`Image ${currentIndex + 1}`}
              className="max-w-full max-h-[80vh] object-contain rounded-lg shadow-2xl select-none"
              style={{
                transform: `scale(${scale}) translate(${position.x / scale}px, ${position.y / scale}px)`,
                transition: isDragging ? 'none' : 'transform 0.2s ease-out',
                cursor: scale > 1 ? (isDragging ? 'grabbing' : 'grab') : 'zoom-in',
              }}
              onError={(e) => {
                (e.target as HTMLImageElement).src = '/placeholder.png';
              }}
              onClick={(e) => {
                e.stopPropagation();
                if (scale === 1) {
                  handleZoomIn();
                }
              }}
              draggable={false}
            />
          )}
        </div>

        {/* Note Section */}
        {currentNote && (
          <div className="absolute bottom-20 left-1/2 -translate-x-1/2 max-w-2xl w-full mx-auto px-4">
            <div className="bg-black/80 text-white px-6 py-4 rounded-lg shadow-xl backdrop-blur-sm">
              <p className="text-sm leading-relaxed text-center">{currentNote}</p>
            </div>
          </div>
        )}

        {/* Thumbnail Navigation */}
        {images.length > 1 && scale === 1 && (
          <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2 bg-black/60 px-4 py-3 rounded-lg max-w-[90vw] overflow-x-auto">
            {images.map((img, idx) => {
              const isThumbVideo = img.endsWith('.mp4');
              return (
                <button
                  key={idx}
                  className={cn(
                    "flex-shrink-0 w-16 h-16 rounded border-2 overflow-hidden transition-all",
                    currentIndex === idx
                      ? "border-white ring-2 ring-white"
                      : "border-gray-400 hover:border-white opacity-60 hover:opacity-100"
                  )}
                  onClick={() => setCurrentIndex(idx)}
                >
                  {isThumbVideo ? (
                    <video
                      src={getImageUrl(img)}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <img
                      src={getImageUrl(img)}
                      alt={`Thumbnail ${idx + 1}`}
                      className="w-full h-full object-cover"
                      onError={(e) => {
                        (e.target as HTMLImageElement).src = '/placeholder.png';
                      }}
                    />
                  )}
                </button>
              );
            })}
          </div>
        )}

        {/* Zoom Hint */}
        {scale === 1 && !isVideo && (
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 pointer-events-none">
            <div className="bg-black/40 text-white px-4 py-2 rounded-lg text-sm opacity-0 hover:opacity-100 transition-opacity">
              Click để phóng to • Scroll để zoom • Kéo khi đã zoom
            </div>
          </div>
        )}
      </div>

      {/* Next Button */}
      {images.length > 1 && scale === 1 && (
        <Button
          variant="ghost"
          size="icon"
          className="absolute right-4 top-1/2 -translate-y-1/2 z-50 text-white hover:bg-white/20 h-16 w-16"
          onClick={handleNext}
        >
          <ChevronRight className="h-8 w-8" />
        </Button>
      )}
    </div>
  );
}
