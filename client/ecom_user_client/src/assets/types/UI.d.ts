import { ImageProps } from 'primereact/image';
import { ReactNode } from 'react';

interface PageProps {
    params: any;
    searchParams?: any;
    children?: ReactNode;
}

interface BreadcrumbProps {
    label: string;
    url: string;
    icon?: string;
}

interface LoaderProps {
    label?: string;
    show?: boolean;
}

interface DividerProps {
    height?: number;
    color?: string;
    mx?: number;
    my?: number;
    px?: number;
    py?: number;
}

interface CustomImageProps extends ImageProps {}

export type { PageProps, BreadcrumbProps, LoaderProps, DividerProps, CustomImageProps };
