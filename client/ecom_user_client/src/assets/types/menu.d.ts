import React from 'react';
interface MenuItemType {
    parent?: string;
    code?: string;
    label?: string;
    items?: MenuItemType[];
    icon?: React.JSX.Element;
    to?: string;
    itemClassName?: string;
    labelClassName?: string;
    iconClassName?: string;
    permission?: string[];
    checkPermission?: string[];
    onItemClick?: (_item: MenuItemType) => void;
    onSubItemClick?: (_item: MenuItemType) => void;
}


export type { MenuItemType };
