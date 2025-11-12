import { MenuItemType } from '@/assets/types/menu';
import ROUTER from './routers';
import { roleE } from './general';
import { Ribeye } from 'next/font/google';


const SIDEBAR_MENU: MenuItemType[] = [
    {
        code: 'Home',
        label: 'Trang chủ',
        parent: 'Home',
        to: ROUTER.home
    },
    {
        code: 'products',
        label: 'Sản phẩm',
        parent: 'products',
        to: ROUTER.product
    },
    {
        code: 'cart',
        label: 'Giỏ hàng',
        parent: 'cart',
        to: ROUTER.giohang
    },
    {
        code: 'orders',
        label: 'Đơn hàng',
        parent: 'orders',
        to: ROUTER.donhang
    },
    {
        code: 'profile',
        label: 'Thông tin cá nhân',
        parent: 'profile',
        to: ROUTER.profile
    }
];

export { SIDEBAR_MENU };
