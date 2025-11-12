interface MetaType {
    currentPage: number;
    hasNextPage?: boolean;
    hasPreviousPage?: boolean;
    messages?: string[];
    limit: number;
    totalCount?: number;
    totalPages: number;
}
interface ParamType {
    page: number,
    limit: number,
    orderBy: string,
    orderDirection: string
}

interface ResponseType<dataType = any> {
    result: dataType;
    paging: {
        totalPages: number,
        currentPages: number,
        limitPages: number
    }
    filters: any;
    status: string;
    message: string;
    code: number;
}

export type { MetaType, ParamType, ResponseType };
