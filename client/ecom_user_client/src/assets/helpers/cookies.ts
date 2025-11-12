import { deleteCookie, getCookie, setCookie } from 'cookies-next';
import { ACCESS_TOKEN, INFO_USER, REFERSH_TOKEN, ROLE_USER } from '../configs/request';
import axios from 'axios';

// Import dynamically to avoid circular dependency
let tokenRefreshService: any = null;
if (typeof window !== 'undefined') {
    import('@/lib/tokenRefreshService').then(module => {
        tokenRefreshService = module.tokenRefreshService;
    });
}

const getCookieValues = <T>(key: string): T | undefined => {
    const value = getCookie(key);
    console.log("Get Cookie Value:", key, value);
    if (!value) {
        return undefined;
    }
    try {
        return JSON.parse(value) as T;
    } catch (e) {
        return JSON.parse(JSON.stringify(value));
    }
};


const setCookieValues = (key: string, value: any, expiresAt?: number) => {
    const options: any = {};

    if (expiresAt) {
        options.expires = new Date(expiresAt * 1000);
    }

    if (typeof value === 'string') {
        setCookie(key, value, options);
    } else {
        setCookie(key, JSON.stringify(value), options);
    }
};

const removeCookies = (key: string) => {
    deleteCookie(key);
};

const logOut = () => {
    const token = getCookie(ACCESS_TOKEN);
    
    // Stop token refresh service
    if (tokenRefreshService) {
        console.log('🛑 Stopping token refresh service on logout...');
        tokenRefreshService.stopScheduler();
    }
    
    // axios.post("http://localhost:8888/auth/log-out", {
    //     token: token
    // })
    removeCookies(ACCESS_TOKEN);
    removeCookies(REFERSH_TOKEN);
    localStorage.removeItem(INFO_USER);
    
    // Clear chat history on logout
    localStorage.removeItem('chat_session_id');
    localStorage.removeItem('chat_messages');
    
    window.location.reload();
};

export { getCookieValues, logOut, removeCookies, setCookieValues };
