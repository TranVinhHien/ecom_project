import React from 'react';
import ROUTER from './routers';
import { MenuItemType } from '../types/menu';

const USER_MENU = (): MenuItemType[] => {

    return [
        {
            code: 'info',
            parent: 'info',
            to: ROUTER.information,
            label: 'Thông tin',

        },
        {
            code: 'logout',
            parent: 'logout',
            label: 'Đăng xuất',
            to: ROUTER.auth.login,
        },
    ];
};

export default USER_MENU;
