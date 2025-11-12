import { MetaType, ResponseType } from '@/assets/types/request';
import axios, { AxiosRequestConfig, AxiosResponse } from 'axios';
import API from '../configs/api';

export const request = axios.create({
    baseURL: API.base_vinh,
    headers: {
        'Content-Type': 'application/json',
    },
});

request.interceptors.request.use(
    (config) => {
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
        throw error
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
    try {
        const response = request.post(path, data, configs);
        return response;
    } catch (error) {
        return Promise.reject(error);
    }
    
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
