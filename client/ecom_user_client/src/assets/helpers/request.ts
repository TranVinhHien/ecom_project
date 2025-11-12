import { MetaType, ResponseType } from '@/assets/types/request';
import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import API from '../configs/api';
import * as cookies from "./cookies";
import { ACCESS_TOKEN, REFERSH_TOKEN } from '../configs/request';
import { jwtDecode } from 'jwt-decode';
const request = axios.create({
    baseURL: API.base,
    headers: {
        'Content-Type': 'application/x-www-form-urlencoded'
    },
});

request.interceptors.request.use(
    async (config) => {
        const token = cookies.getCookieValues(ACCESS_TOKEN);
        if (!config.headers.Authorization) {
            config.headers.Authorization = '';
        }
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        else {
            const refresh_token = cookies.getCookieValues(REFERSH_TOKEN)
            if (refresh_token) {
                try {
                    const res = await axios.post(API.base + API.user.new_access_token, { refresh_token })
                    const accessToken = res?.data?.result.access_token
                    if (accessToken) {
                        const decodedACC: any = jwtDecode(accessToken);
                        cookies.setCookieValues(ACCESS_TOKEN, accessToken, decodedACC?.exp);
                        config.headers.Authorization = `Bearer ${accessToken}`;
                    }
                } catch (error) {
                    cookies.logOut()
                    throw new Error("You need login!")
                }
            } else {
                cookies.logOut()
                throw new Error("You need login!")
            }
        }
        return config;
    },
    (error) => {
        return Promise.reject(error);
    },
);

request.interceptors.response.use(
    (response) => {
        return response;
    },
    async (error: any) => {
        // if (error.response.data.code === 9993) { // token  hết hạn
        //     // const token = cookies.getCookieValues(ACCESS_TOKEN);
        //     // if (token === null) {
        //     //     throw new Error("you need login !!!")
        //     // }
        //     // const req = await axios.post(API.base + API.auth.refresh, {
        //     //     token: token
        //     // })
        //     // const decoded: any = jwtDecode(req.data.result.accessToken);
        //     // const arr = decoded.scope.split(" ");
        //     // cookies.setCookieValues(ACCESS_TOKEN, req.data.result.accessToken, { expires: new Date(decoded.exp * 1000 + 5 * 1000) })
        //     // //cookies.set(REFERSH_TOKEN, response.data.result.refreshToken)
        //     // cookies.setCookieValues(ROLE_USER, arr, { expires: new Date(decoded.exp * 1000) })
        //     return
        // }

        throw error
        // return Promise.reject(error);
    },
);


const get = <T = any>(path: string, configs?: AxiosRequestConfig): Promise<AxiosResponse<ResponseType<T>, any>> => {
    const response = request.get(path, configs);
    return response;
};


const post = <T = any>(
    path: string,
    data: any,
    configs?: AxiosRequestConfig,
): Promise<AxiosResponse<ResponseType<T>, any>> => {
    const response = request.post(path, data, configs);
    return response;
};

const update = <T = any>(
    path: string,
    data: any,
    configs?: AxiosRequestConfig,
): Promise<AxiosResponse<ResponseType<T>, any>> => {
    const response = request.put(path, data, configs);
    return response;
};

const remove = <T = any>(path: string, configs?: AxiosRequestConfig
): Promise<AxiosResponse<ResponseType<T>, any>> => {
    const response = request.delete(path, configs);
    return response;
};
const defaultMeta: MetaType = {
    currentPage: 1,
    hasNextPage: false,
    hasPreviousPage: false,
    messages: [],
    limit: 10,
    totalCount: 1,
    totalPages: 1,
};
const ROW_PER_PAGE = [10, 20, 30]
export { get, post, remove, update, defaultMeta, ROW_PER_PAGE };
